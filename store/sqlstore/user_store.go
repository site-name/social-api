package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/store"
)

// const (
// 	MaxGroupChannelsForProfiles = 50
// )

var (
	UserSearchTypeNames_NO_FULL_NAME = []string{"Username", "Nickname"}
	UserSearchTypeNames              = []string{"Username", "FirstName", "LastName", "Nickname"}
	UserSearchTypeAll_NO_FULL_NAME   = []string{"Username", "Nickname", "Email"}
	UserSearchTypeAll                = []string{"Username", "FirstName", "LastName", "Nickname", "Email"}
)

type SqlUserStore struct {
	*SqlStore
	metrics einterfaces.MetricsInterface

	// usersQuery is a starting point for all queries that return one or more Users.
	usersQuery squirrel.SelectBuilder
}

func (us *SqlUserStore) ClearCaches() {}

func newSqlUserStore(sqlStore *SqlStore, metrics einterfaces.MetricsInterface) store.UserStore {
	us := &SqlUserStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	// note: we are providing field names explicitly here to maintain order of columns (needed when using raw queries)
	us.usersQuery = us.getQueryBuilder().
		Select(
			"u.Id",
			"u.CreateAt",
			"u.UpdateAt",
			"u.DeleteAt",
			"u.Username",
			"u.Password",
			"u.AuthData",
			"u.AuthService",
			"u.Email",
			"u.EmailVerified",
			"u.Nickname",
			"u.FirstName",
			"u.LastName",
			"u.Roles",
			"u.Props",
			"u.NotifyProps",
			"u.LastPasswordUpdate",
			"u.LastPictureUpdate",
			"u.FailedAttempts",
			"u.Locale",
			"u.Timezone",
			"u.MfaActive",
			"u.MfaSecret",
			"u.DefaultShippingAddressID",
			"u.DefaultBillingAddressID",
		).
		From("Users u")

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(account.User{}, "Users").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH) // 16
		table.ColMap("Username").SetMaxSize(account.USER_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Password").SetMaxSize(account.USER_HASH_PASSWORD_MAX_LENGTH)
		table.ColMap("AuthData").SetMaxSize(account.USER_AUTH_DATA_MAX_LENGTH).SetUnique(true)
		table.ColMap("AuthService").SetMaxSize(32)
		table.ColMap("Email").SetMaxSize(account.USER_EMAIL_MAX_LENGTH).SetUnique(true)
		table.ColMap("Nickname").SetMaxSize(account.USER_NAME_MAX_LENGTH)
		table.ColMap("FirstName").SetMaxSize(account.USER_FIRST_NAME_MAX_RUNES)
		table.ColMap("LastName").SetMaxSize(account.USER_LAST_NAME_MAX_RUNES)
		table.ColMap("Roles").SetMaxSize(256)
		table.ColMap("Locale").SetMaxSize(account.USER_LOCALE_MAX_LENGTH)
		table.ColMap("MfaSecret").SetMaxSize(128)
		table.ColMap("Timezone").SetMaxSize(account.USER_TIMEZONE_MAX_RUNES)
		table.ColMap("NotifyProps").SetMaxSize(2000)
		table.ColMap("Props").SetMaxSize(4000)
	}

	return us
}

func (us *SqlUserStore) createIndexesIfNotExists() {
	us.CreateIndexIfNotExists("idx_users_email", "Users", "Email")
	us.CreateIndexIfNotExists("idx_users_update_at", "Users", "UpdateAt")
	us.CreateIndexIfNotExists("idx_users_create_at", "Users", "CreateAt")
	us.CreateIndexIfNotExists("idx_users_delete_at", "Users", "DeleteAt")
	us.CreateIndexIfNotExists("idx_users_email_lower_textpattern", "Users", "lower(Email) text_pattern_ops")
	us.CreateIndexIfNotExists("idx_users_username_lower_textpattern", "Users", "lower(Username) text_pattern_ops")
	us.CreateIndexIfNotExists("idx_users_nickname_lower_textpattern", "Users", "lower(Nickname) text_pattern_ops")
	us.CreateIndexIfNotExists("idx_users_firstname_lower_textpattern", "Users", "lower(FirstName) text_pattern_ops")
	us.CreateIndexIfNotExists("idx_users_lastname_lower_textpattern", "Users", "lower(LastName) text_pattern_ops")
	us.CreateFullTextIndexIfNotExists("idx_users_all_txt", "Users", strings.Join(UserSearchTypeAll, ", "))
	us.CreateFullTextIndexIfNotExists("idx_users_all_no_full_name_txt", "Users", strings.Join(UserSearchTypeAll_NO_FULL_NAME, ", "))
	us.CreateFullTextIndexIfNotExists("idx_users_names_txt", "Users", strings.Join(UserSearchTypeNames, ", "))
	us.CreateFullTextIndexIfNotExists("idx_users_names_no_full_name_txt", "Users", strings.Join(UserSearchTypeNames_NO_FULL_NAME, ", "))
}

// DeactivateGuests
func (us *SqlUserStore) DeactivateGuests() ([]string, error) {
	curTime := model.GetMillis()
	updateQuery := us.
		getQueryBuilder().
		Update("Users").
		Set("UpdateAt", curTime).
		Set("DeleteAt", curTime).
		Where(squirrel.Eq{"Roles": "system_guest"}).
		Where(squirrel.Eq{"DeleteAt": 0})

	queryString, args, err := updateQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "deactivate_guests_tosql")
	}

	_, err = us.GetMaster().Exec(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Users with roles=system_guest")
	}

	selectQuery := us.
		getQueryBuilder().
		Select("Id").
		From("Users").
		Where(squirrel.Eq{"DeleteAt": curTime})

	queryString, args, err = selectQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "deactivate_guests_tosql")
	}

	userIds := []string{}
	_, err = us.GetMaster().Select(&userIds, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	return userIds, nil
}

// ResetAuthDataToEmailForUsers resets the AuthData of users whose AuthService
// is |service| to their Email. If userIDs is non-empty, only the users whose
// IDs are in userIDs will be affected. If dryRun is true, only the number
// of users who *would* be affected is returned; otherwise, the number of
// users who actually were affected is returned.
func (us *SqlUserStore) ResetAuthDataToEmailForUsers(service string, userIDs []string, includeDeleted bool, dryRun bool) (int, error) {
	whereEquals := squirrel.Eq{"AuthService": service}
	if len(userIDs) > 0 {
		whereEquals["Id"] = userIDs
	}
	if !includeDeleted {
		whereEquals["DeleteAt"] = 0
	}

	if dryRun {
		builder := us.
			getQueryBuilder().
			Select("COUNT(*)").
			From("Users").
			Where(whereEquals)
		query, args, err := builder.ToSql()
		if err != nil {
			return 0, errors.Wrap(err, "select_count_users_tosql")
		}
		numAffected, err := us.GetReplica().SelectInt(query, args...)
		return int(numAffected), err
	}
	builder := us.
		getQueryBuilder().
		Update("Users").
		Set("AuthData", squirrel.Expr("Email")).
		Where(whereEquals)
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "update_users_tosql")
	}
	result, err := us.GetMaster().Exec(query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to update user's AuthData")
	}
	numAffected, err := result.RowsAffected()
	return int(numAffected), err
}

func (us *SqlUserStore) InvalidateProfileCacheForUser(userId string) {}

func (us *SqlUserStore) GetEtagForProfiles(teamId string) string {
	// updateAt, err := us.GetReplica().SelectInt("SELECT UpdateAt FROM User")
	// if err != nil {
	// 	return fmt.Sprintf("%v.%v", model.CurrentVersion, updateAt)
	// }
	// return fmt.Sprintf("%v.%v", model.CurrentVersion, updateAt)
	panic("not implemented")
}

func (us *SqlUserStore) GetEtagForAllProfiles() string {
	updateAt, err := us.GetReplica().SelectInt("SELECT UpdateAt FROM Users ORDER BY UpdateAt DESC LIMIT 1")
	if err != nil {
		return fmt.Sprintf("%v.%v", model.CurrentVersion, model.GetMillis())
	}
	return fmt.Sprintf("%v.%v", model.CurrentVersion, updateAt)
}

func (us *SqlUserStore) Save(user *account.User) (*account.User, error) {
	// if user.Id != "" && !user
	user.PreSave()
	if err := user.IsValid(); err != nil {
		return nil, err
	}

	if err := us.GetMaster().Insert(user); err != nil {
		if IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique"}) {
			return nil, store.NewErrInvalidInput("User", "email", user.Email)
		}
		if IsUniqueConstraintError(err, []string{"Username", "users_username_key", "idx_users_username_unique"}) {
			return nil, store.NewErrInvalidInput("User", "username", user.Username)
		}
		return nil, errors.Wrapf(err, "failed to save User with userId=%s", user.Id)
	}

	return user, nil
}

// Update updates user
func (us *SqlUserStore) Update(user *account.User, trustedUpdateData bool) (*account.UserUpdate, error) {
	user.PreUpdate()

	if err := user.IsValid(); err != nil {
		return nil, err
	}

	oldUserResult, err := us.GetMaster().Get(account.User{}, user.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user with userId=%s", user.Id)
	}

	if oldUserResult == nil {
		return nil, store.NewErrInvalidInput("User", "id", user.Id)
	}

	oldUser := oldUserResult.(*account.User)
	user.CreateAt = oldUser.CreateAt
	user.AuthData = oldUser.AuthData
	user.AuthService = oldUser.AuthService
	user.Password = oldUser.Password
	user.LastPasswordUpdate = oldUser.LastPasswordUpdate
	user.LastPictureUpdate = oldUser.LastPictureUpdate
	user.EmailVerified = oldUser.EmailVerified
	user.FailedAttempts = oldUser.FailedAttempts
	user.MfaSecret = oldUser.MfaSecret
	user.MfaActive = oldUser.MfaActive

	if !trustedUpdateData {
		user.Roles = oldUser.Roles
		user.DeleteAt = oldUser.DeleteAt
	}

	if user.IsOAuthUser() {
		if !trustedUpdateData {
			user.Email = oldUser.Email
		}
	} else if user.IsLDAPUser() && !trustedUpdateData {
		if user.Username != oldUser.Username || user.Email != oldUser.Email {
			return nil, store.NewErrInvalidInput("User", "id", user.Id)
		}
	} else if user.Email != oldUser.Email {
		user.EmailVerified = false
	}

	if user.Username != oldUser.Username {
		user.UpdateMentionKeysFromUsername(oldUser.Username)
	}

	count, err := us.GetMaster().Update(user)
	if err != nil {
		if IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique"}) {
			return nil, store.NewErrInvalidInput("User", "id", user.Id)
		}
		if IsUniqueConstraintError(err, []string{"Username", "users_username_key", "idx_users_username_unique"}) {
			return nil, store.NewErrInvalidInput("User", "id", user.Id)
		}
		return nil, errors.Wrapf(err, "failed to update User with userId=%s", user.Id)
	}

	if count > 1 {
		return nil, fmt.Errorf("multiple users were update: userId=%s, count=%d", user.Id, count)
	}

	user.Sanitize(map[string]bool{})
	oldUser.Sanitize(map[string]bool{})
	return &account.UserUpdate{New: user, Old: oldUser}, nil
}

func (us *SqlUserStore) UpdateLastPictureUpdate(userId string) error {
	now := model.GetMillis()
	if _, err := us.GetMaster().Exec("UPDATE Users SET LastPictureUpdate = :Time, UpdateAt = :Time WHERE Id = :UserId", map[string]interface{}{"Time": now, "UserId": userId}); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) ResetLastPictureUpdate(userId string) error {
	now := model.GetMillis()
	if _, err := us.GetMaster().Exec("UPDATE Users SET LastPictureUpdate = :PictureUpdateTime, UpdateAt = :UpdateTime WHERE Id = :UserId", map[string]interface{}{"PictureUpdateTime": 0, "UpdateTime": now, "UserId": userId}); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) UpdateUpdateAt(userId string) (int64, error) {
	now := model.GetMillis()
	if _, err := us.GetMaster().Exec("UPDATE Users SET UpdateAt = :Time WHERE Id = :UserId", map[string]interface{}{"Time": now, "UserId": userId}); err != nil {
		return now, errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return now, nil
}

func (us *SqlUserStore) UpdatePassword(userId, hashedPassword string) error {
	now := model.GetMillis()
	if _, err := us.GetMaster().Exec("UPDATE Users SET Password = :Password, LastPasswordUpdate = :LastPasswordUpdate, UpdateAt = :UpdateAt, AuthData = NULL, AuthService = '', FailedAttempts = 0 WHERE Id = :UserId", map[string]interface{}{
		"Password":           hashedPassword,
		"LastPasswordUpdate": now,
		"UpdateAt":           now,
		"UserId":             userId,
	}); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) UpdateFailedPasswordAttempts(userId string, attempts int) error {
	if _, err := us.GetMaster().Exec(
		"UPDATE Users SET FailedAttempts = :FailedAttempts WHERE Id = :UserId",
		map[string]interface{}{
			"FailedAttempts": attempts,
			"UserId":         userId,
		},
	); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}
	return nil
}

// UpdateAuthData updates auth data of user
func (us *SqlUserStore) UpdateAuthData(userId string, service string, authData *string, email string, resetMfa bool) (string, error) {
	updateAt := model.GetMillis()
	query := `
		UPDATE
			Users
		SET
			Password = '',
			LastPasswordUpdate = :LastPasswordUpdate,
			UpdateAt = :UpdateAt,
			FailedAttempts = 0,
			AuthService = :AuthService,
			AuthData = :AuthData
	`
	if email != "" {
		query += ", Email = lower(:Email)"
	}
	if resetMfa {
		query += ", MfaActive = false, MfaSecret = ''"
	}
	query += " WHERE Id = :UserId"

	if _, err := us.
		GetMaster().
		Exec(query, map[string]interface{}{
			"LastPasswordUpdate": updateAt,
			"UpdateAt":           updateAt,
			"UserId":             userId,
			"AuthService":        service,
			"AuthData":           authData,
			"Email":              email,
		}); err != nil {
		if IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique", "AuthData", "users_authdata_key"}) {
			return "", store.NewErrInvalidInput("User", "id", userId)
		}
		return "", errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}
	return userId, nil
}

// UpdateMfaSecret updates mfa secret for current user
func (us *SqlUserStore) UpdateMfaSecret(userId, secret string) error {
	updateAt := model.GetMillis()

	if _, err := us.
		GetMaster().
		Exec("UPDATE Users SET MfaSecret = :Secret, UpdateAt = :UpdateAt WHERE Id = :UserId",
			map[string]interface{}{"Secret": secret, "UpdateAt": updateAt, "userId": userId}); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) UpdateMfaActive(userId string, active bool) error {
	updateAt := model.GetMillis()
	if _, err := us.GetMaster().Exec("UPDATE Users SET MfaActive = :Active, UpdateAt = :UpdateAt WHERE Id = :UserId",
		map[string]interface{}{"Active": active, "UpdateAt": updateAt, "UserId": userId}); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

// GetMany returns a list of users for the provided list of ids
func (us *SqlUserStore) GetMany(ctx context.Context, ids []string) ([]*account.User, error) {
	query := us.usersQuery.Where(squirrel.Eq{"Id": ids})
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "users_get_many_tosql")
	}

	var users []*account.User
	if _, err := us.SqlStore.DBFromContext(ctx).Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "users_get_many_select")
	}

	return users, nil
}

// Get returns single user that has Id matches given id.
// If an user with given id does not exist, return nil with according error
func (us *SqlUserStore) Get(ctx context.Context, id string) (*account.User, error) {
	query := us.usersQuery.Where("Id = ?", id)
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "users_get_tosql")
	}
	row := us.SqlStore.DBFromContext(ctx).Db.QueryRow(queryString, args...)

	var user account.User
	var props, notifyProps, timezone []byte
	err = row.Scan(
		&user.Id,
		&user.CreateAt,
		&user.UpdateAt,
		&user.DeleteAt,
		&user.Username,
		&user.Password,
		&user.AuthData,
		&user.AuthService,
		&user.Email,
		&user.EmailVerified,
		&user.Nickname,
		&user.FirstName,
		&user.LastName,
		&user.Roles,
		&props,
		&notifyProps,
		&user.LastPasswordUpdate,
		&user.LastPictureUpdate,
		&user.FailedAttempts,
		&user.Locale,
		&timezone,
		&user.MfaActive,
		&user.MfaSecret,
		&user.DefaultBillingAddressID,
		&user.DefaultShippingAddressID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("User", id)
		}
		return nil, errors.Wrapf(err, "failed to get User with userId=%s", id)
	}
	if err = json.JSON.Unmarshal(props, &user.Props); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user props")
	}
	if err = json.JSON.Unmarshal(notifyProps, &user.NotifyProps); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user notify props")
	}
	if err = json.JSON.Unmarshal(timezone, &user.Timezone); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user timezone")
	}

	return &user, nil
}

// GetAll fetches all users from database and returns to the caller
func (us *SqlUserStore) GetAll() ([]*account.User, error) {
	query := us.usersQuery.OrderBy("Username ASC")

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_all_users_tosql")
	}

	var data []*account.User
	if _, err := us.GetReplica().Select(&data, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	return data, nil
}

// GetAllAfter get users that have id less than given id
func (us *SqlUserStore) GetAllAfter(limit int, afterId string) ([]*account.User, error) {
	query := us.usersQuery.Where("Id > ?", afterId).OrderBy("Id ASC").Limit(uint64(limit))

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_all_after_tosql")
	}

	var users []*account.User
	if _, err := us.GetReplica().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	return users, nil
}

func (us *SqlUserStore) GetProfiles(options *account.UserGetOptions) ([]*account.User, error) {
	// query := us.
	// 	usersQuery.
	// 	Join("TeamMembers tm ON ( tm.UserId = u.Id AND tm.DeleteAt = 0 )").
	// 	Where("tm.TeamId = ?", options.InTeamId)
	panic("not implemented")
}

func (us *SqlUserStore) GetProfilesByUsernames(usernames []string /*, viewRestrictions *account.ViewUsersRestrictions */) ([]*account.User, error) {
	query := us.usersQuery.Where(squirrel.Eq{"Username": usernames}).OrderBy("u.Username ASC")

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_profiles_by_usernames")
	}

	var users []*account.User
	if _, err := us.GetReplica().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	return users, nil
}

func (us *SqlUserStore) GetProfileByIds(ctx context.Context, userIds []string, options *store.UserGetByIdsOpts, allowFromCache bool) ([]*account.User, error) {
	if options == nil {
		options = new(store.UserGetByIdsOpts)
	}

	users := []*account.User{}
	query := us.usersQuery.Where(squirrel.Eq{"u.Id": userIds}).OrderBy("u.Username ASC")

	if options.Since > 0 {
		query = query.Where(squirrel.Gt{"u.UpdateAt": options.Since})
	}

	query = applyViewRestrictionsFilter(query, true)
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_profile_by_ids_tosql")
	}

	if _, err := us.SqlStore.DBFromContext(ctx).Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	return users, nil
}

func (us *SqlUserStore) GetSystemAdminProfiles() (map[string]*account.User, error) {
	query := us.usersQuery.Where("Roles LIKE ?", "%system_admin%").OrderBy("u.Username ASC")
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_system_admin_profiles_tosql")
	}

	var users []*account.User
	if _, err := us.GetReplica().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	userMap := make(map[string]*account.User)
	for _, u := range users {
		u.Sanitize(map[string]bool{})
		userMap[u.Id] = u
	}

	return userMap, nil
}

func (us *SqlUserStore) GetByEmail(email string) (*account.User, error) {
	query := us.usersQuery.Where("Email = lower(?)", email)
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_by_email_tosql")
	}

	user := account.User{}
	if err := us.GetReplica().SelectOne(&user, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.NewErrNotFound("User", fmt.Sprintf("email=%s", email)), "failed to find User")
		}

		return nil, errors.Wrapf(err, "failed to get User with email=%s", email)
	}

	return &user, nil
}

func (us *SqlUserStore) GetByAuth(authData *string, authService string) (*account.User, error) {
	if authData == nil || *authData == "" {
		return nil, store.NewErrInvalidInput("User", "<authData>", "empty or nil")
	}

	query := us.usersQuery.
		Where("u.AuthData = ?", *authData).
		Where("u.AuthService = ?", authService)
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_by_auth_tosql")
	}

	user := account.User{}
	if err := us.GetReplica().SelectOne(&user, queryString, args...); err == sql.ErrNoRows {
		return nil, store.NewErrNotFound("User", fmt.Sprintf("authData=%s, authService=%s", *authData, authService))
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to find User with authData=%s")
	}
	return &user, nil
}

func (us *SqlUserStore) GetAllUsingAuthService(authService string) ([]*account.User, error) {
	query := us.
		usersQuery.
		Where("u.AuthService = ?", authService).
		OrderBy("u.Username ASC")
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_all_using_auth_service_tosql")
	}

	var users []*account.User
	if _, err := us.GetReplica().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Users with authService=%s", authService)
	}

	return users, nil
}

func (us *SqlUserStore) GetAllNotInAuthService(authServices []string) ([]*account.User, error) {
	query := us.
		usersQuery.
		Where(squirrel.NotEq{"u.AuthService": authServices}).
		OrderBy("u.Username ASC")

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_all_not_in_auth_service_tosql")
	}
	var users []*account.User
	if _, err := us.GetReplica().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Users with authService in %v", authServices)
	}

	return users, nil
}

func (us *SqlUserStore) GetByUsername(username string) (*account.User, error) {
	query := us.usersQuery.Where("u.Username = lower(?)", username)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_by_username_tosql")
	}

	var user *account.User
	if err := us.GetReplica().SelectOne(&user, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.NewErrNotFound("User", fmt.Sprintf("username=%s", username)), "failed to find User")
		}

		return nil, errors.Wrapf(err, "failed to find User with username=%s", username)
	}

	return user, nil
}

func (us *SqlUserStore) GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail bool) (*account.User, error) {
	query := us.usersQuery
	if allowSignInWithUsername && allowSignInWithEmail {
		query = query.Where("Username = lower(?) OR Email = lower(?)", loginId, loginId)
	} else if allowSignInWithUsername {
		query = query.Where("Username = lower(?)", loginId)
	} else if allowSignInWithEmail {
		query = query.Where("Email = lower(?)", loginId)
	} else {
		return nil, errors.New("sign in with username and email are disabled")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_for_login_tosql")
	}

	users := []*account.User{}
	if _, err := us.GetReplica().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	if len(users) == 0 {
		return nil, errors.New("user not found")
	}

	if len(users) > 1 {
		return nil, errors.New("multiple users found")
	}

	return users[0], nil
}

func (us *SqlUserStore) VerifyEmail(userId, email string) (string, error) {
	now := model.GetMillis()
	if _, err := us.
		GetMaster().
		Exec("UPDATE Users SET Email = lower(:email), EmailVerified = true, UpdateAt = :Time WHERE Id = :UserId",
			map[string]interface{}{"email": email, "Time": now, "UserId": userId}); err != nil {
		return "", errors.Wrapf(err, "failed to update Users with userId=%s and email=%s")
	}

	return userId, nil
}

func (us *SqlUserStore) PermanentDelete(userId string) error {
	if _, err := us.GetMaster().Exec("DELETE FROM Users WHERE Id = :UserId", map[string]interface{}{"UserId": userId}); err != nil {
		return errors.Wrapf(err, "failed to delete User with userId=%s", userId)
	}
	return nil
}

func applyViewRestrictionsFilter(query squirrel.SelectBuilder, distinct bool) squirrel.SelectBuilder {
	if distinct {
		return query.Distinct()
	}

	return query
}

func (us *SqlUserStore) Count(options account.UserCountOptions) (int64, error) {
	query := us.getQueryBuilder().Select("COUNT(DISTINCT u.Id)").From("Users AS u")

	if !options.IncludeDeleted {
		query = query.Where("u.DeleteAt = 0")
	}

	if options.IncludeBotAccounts {
		if options.ExcludeRegularUsers {
			query = query.Join("Bots ON u.Id = Bots.UserId")
		}
	} else {
		query = query.LeftJoin("Bots ON u.Id = Bots.UserId").Where("Bots.UserId IS NULL")
		if options.ExcludeRegularUsers {
			// Currently this doesn't make sense because it will always return 0
			return int64(0), errors.New("query with IncludeBotAccounts=false and excludeRegularUsers=true always return 0")
		}
	}

	query = applyViewRestrictionsFilter(query, true)
	query = applyMultiRoleFilters(query, options.Roles, true)

	query = query.PlaceholderFormat(squirrel.Dollar)
	queryString, args, err := query.ToSql()
	if err != nil {
		return int64(0), errors.Wrap(err, "count_tosql")
	}

	count, err := us.GetReplica().SelectInt(queryString, args...)
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count Users")
	}
	return count, nil
}

func (us *SqlUserStore) AnalyticsActiveCount(timePeriod int64, options account.UserCountOptions) (int64, error) {
	time := model.GetMillis() - timePeriod
	query := us.getQueryBuilder().Select("COUNT(*)").From("Status AS s").Where("LastActivityAt > :Time", map[string]interface{}{"Time": time})
	if !options.IncludeBotAccounts {
		query = query.LeftJoin("Bots ON s.UserId = Bots.UserId").Where("Bots.UserId IS NULL")
	}

	if !options.IncludeDeleted {
		query = query.LeftJoin("Users ON s.UserId = Users.Id").Where("Users.DeleteAt = 0")
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "analytics_active_count_tosql")
	}

	v, err := us.GetReplica().SelectInt(queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count Users")
	}
	return v, nil
}

func (us *SqlUserStore) AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options account.UserCountOptions) (int64, error) {
	query := us.getQueryBuilder().Select("COUNT(*)").From("Status AS s").Where("LastActivityAt > :StartTime AND LastActivityAt <= :EndTime", map[string]interface{}{"StartTime": startTime, "EndTime": endTime})
	if !options.IncludeBotAccounts {
		query = query.LeftJoin("Bots ON s.UserId = Bots.UserId").Where("Bots.UserId IS NULL")
	}
	if !options.IncludeDeleted {
		query = query.LeftJoin("Users ON s.UserId = Users.Id").Where("Users.DeleteAt = 0")
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to build query.")
	}

	v, err := us.GetReplica().SelectInt(queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to get the active users during the requested period")
	}
	return v, nil
}

func (us *SqlUserStore) GetUnreadCount(userId string) (int64, error) {
	panic("not implemented")
}

func applyMultiRoleFilters(query squirrel.SelectBuilder, systemRoles []string, isPostgreSQL bool) squirrel.SelectBuilder {
	sqOr := squirrel.Or{}

	if len(systemRoles) > 0 && systemRoles[0] != "" {
		for _, role := range systemRoles {
			queryRole := wildcardSearchTerm(role)
			switch role {
			case model.SYSTEM_USER_ROLE_ID:
				// If querying for a `system_user` ensure that the user is only a system_user.
				sqOr = append(sqOr, squirrel.Eq{"u.Roles": role})
			case model.SYSTEM_GUEST_ROLE_ID, model.SYSTEM_ADMIN_ROLE_ID, model.SYSTEM_USER_MANAGER_ROLE_ID, model.SYSTEM_READ_ONLY_ADMIN_ROLE_ID, model.SYSTEM_MANAGER_ROLE_ID:
				// If querying for any other roles search using a wildcard
				if isPostgreSQL {
					sqOr = append(sqOr, squirrel.ILike{"u.Roles": queryRole})
				}
			}
		}
	}

	if len(sqOr) > 0 {
		return query.Where(sqOr)
	}

	return query
}

func generateSearchQuery(query squirrel.SelectBuilder, terms []string, fields []string) squirrel.SelectBuilder {
	for _, term := range terms {
		searchFields := []string{}
		termArgs := []interface{}{}
		for _, field := range fields {
			searchFields = append(searchFields, fmt.Sprintf("lower(%s) LIKE lower(?) escape '*' ", field))
			termArgs = append(termArgs, fmt.Sprintf("%s%%", strings.TrimLeft(term, "@")))
		}
		query = query.Where(fmt.Sprintf("(%s)", strings.Join(searchFields, " OR ")), termArgs...)
	}

	return query
}

func (us *SqlUserStore) Search(teamId string, term string, options *account.UserSearchOptions) ([]*account.User, error) {
	query := us.usersQuery.
		OrderBy("Username ASC").
		Limit(uint64(options.Limit))

	// if teamId != "" {
	// 	query = query.Join("TeamMembers tm ON ( tm.UserId = u.Id AND tm.DeleteAt = 0 AND tm.TeamId = ? )", teamId)
	// }
	return us.performSearch(query, term, options)
}

func applyRoleFilter(query squirrel.SelectBuilder, role string, isPostgreSQL bool) squirrel.SelectBuilder {
	if role == "" {
		return query
	}

	if isPostgreSQL {
		roleParam := fmt.Sprintf("%%%s%%", sanitizeSearchTerm(role, "\\"))
		return query.Where("u.Roles LIKE LOWER(?)", roleParam)
	}

	roleParam := fmt.Sprintf("%%%s%%", sanitizeSearchTerm(role, "*"))
	return query.Where("u.Roles LIKE ? ESCAPE '*'", roleParam)
}

func (us *SqlUserStore) performSearch(query squirrel.SelectBuilder, term string, options *account.UserSearchOptions) ([]*account.User, error) {
	term = sanitizeSearchTerm(term, "*")

	var searchType []string
	if options.AllowEmails {
		if options.AllowFullNames {
			searchType = UserSearchTypeAll
		} else {
			searchType = UserSearchTypeAll_NO_FULL_NAME
		}
	} else {
		if options.AllowFullNames {
			searchType = UserSearchTypeNames
		} else {
			searchType = UserSearchTypeNames_NO_FULL_NAME
		}
	}

	query = applyRoleFilter(query, options.Role, true)
	query = applyMultiRoleFilters(query, options.Roles, true)

	if !options.AllowInactive {
		query = query.Where("u.DeleteAt = 0")
	}

	if strings.TrimSpace(term) != "" {
		query = generateSearchQuery(query, strings.Fields(term), searchType)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "perform_search_tosql")
	}

	var users []*account.User
	if _, err := us.GetReplica().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Users with term=%s and searchType=%v", term, searchType)
	}
	for _, u := range users {
		u.Sanitize(map[string]bool{})
	}

	return users, nil
}

func (us *SqlUserStore) AnalyticsGetInactiveUsersCount() (int64, error) {
	count, err := us.GetReplica().SelectInt("SELECT COUNT(Id) FROM Users WHERE DeletedAt > 0")
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count inactive Users")
	}
	return count, nil
}

func (us *SqlUserStore) AnalyticsGetExternalUsers(hostDomain string) (bool, error) {
	count, err := us.GetReplica().SelectInt("SELECT COUNT(Id) FROM Users WHERE LOWER(Email) NOT LIKE :HostDomain", map[string]interface{}{"HostDomain": "%@" + strings.ToLower(hostDomain)})
	if err != nil {
		return false, errors.Wrap(err, "failed to count inactive Users")
	}
	return count > 0, nil
}

func (us *SqlUserStore) AnalyticsGetGuestCount() (int64, error) {
	count, err := us.GetReplica().SelectInt("SELECT count(*) FROM Users WHERE Roles LIKE :Roles and DeleteAt = 0", map[string]interface{}{"Roles": "%system_guest%"})
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count guest Users")
	}
	return count, nil
}

func (us *SqlUserStore) AnalyticsGetSystemAdminCount() (int64, error) {
	count, err := us.GetReplica().SelectInt("SELECT count(*) FROM Users WHERE Roles LIKE :Roles and DeleteAt = 0", map[string]interface{}{"Roles": "%system_admin%"})
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count system admin Users")
	}
	return count, nil
}

func (us *SqlUserStore) ClearAllCustomRoleAssignments() error {
	builtinRoles := model.MakeDefaultRoles()
	lastUserId := strings.Repeat("0", len(uuid.Nil))

	for {
		var transaction *gorp.Transaction
		var err error

		if transaction, err = us.GetMaster().Begin(); err != nil {
			return errors.Wrap(err, "begin_transaction")
		}
		defer finalizeTransaction(transaction)

		var users []*account.User
		if _, err := transaction.Select(&users, "SELECT * FROM Users WHERE Id > :Id ORDER BY Id LIMIT 1000", map[string]interface{}{"Id": lastUserId}); err != nil {
			return errors.Wrapf(err, "failed to find Users with id > %s", lastUserId)
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			lastUserId = user.Id
			var newRoles []string

			for _, role := range strings.Fields(user.Roles) {
				for name := range builtinRoles {
					if name == role {
						newRoles = append(newRoles, role)
						break
					}
				}
			}

			newRolesString := strings.Join(newRoles, " ")
			if newRolesString != user.Roles {
				if _, err := transaction.Exec("UPDATE Users SET Roles = :Roles WHERE Id = :Id", map[string]interface{}{"Roles": newRolesString, "Id": user.Id}); err != nil {
					return errors.Wrap(err, "failed to update Users")
				}
			}
		}

		if err := transaction.Commit(); err != nil {
			return errors.Wrap(err, "commit_transaction")
		}
	}

	return nil
}

func (us *SqlUserStore) InferSystemInstallDate() (int64, error) {
	createAt, err := us.GetReplica().SelectInt("SELECT CreateAt FROM Users WHERE CreateAt IS NOT NULL ORDER BY CreateAt ASC LIMIT 1")
	if err != nil {
		return 0, errors.Wrap(err, "failed to infer system install date")
	}

	return createAt, nil
}

func (us *SqlUserStore) GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*account.UserForIndexing, error) {
	// var users []*account.User
	// userQuery, args, _ := us.usersQuery.
	// 	Where(squirrel.GtOrEq{"u.CreateAt": startTime}).
	// 	Where(squirrel.Lt{"u.CreateAt": endTime}).
	// 	OrderBy("u.CreateAt").
	// 	Limit(uint64(limit)).
	// 	ToSql()

	// _, err := us.GetSearchReplica().Select(&users, userQuery, args...)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to find Users")
	// }

	// userIds := []string{}
	// for _, user := range users {
	// 	userIds = append(userIds, user.Id)
	// }

	// var channelMembers []*account.
	panic("not implemented")
}

func (us *SqlUserStore) PromoteGuestToUser(userId string) error {
	panic("not implemented")
}

func (us *SqlUserStore) DemoteUserToGuest(userID string) (*account.User, error) {
	panic("not implemented")

}

func (us *SqlUserStore) GetKnownUsers(userId string) ([]string, error) {
	panic("not implemented")

}

func (us *SqlUserStore) GetAllProfiles(options *account.UserGetOptions) ([]*account.User, error) {
	panic("not implemented")

}
