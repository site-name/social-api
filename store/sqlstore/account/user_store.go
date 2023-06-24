package account

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

var (
	UserSearchTypeNames_NO_FULL_NAME = []string{"Username", "Nickname"}
	UserSearchTypeNames              = []string{"Username", "FirstName", "LastName", "Nickname"}
	UserSearchTypeAll_NO_FULL_NAME   = []string{"Username", "Nickname", "Email"}
	UserSearchTypeAll                = []string{"Username", "FirstName", "LastName", "Nickname", "Email"}
)

type SqlUserStore struct {
	store.Store
	metrics einterfaces.MetricsInterface

	// usersQuery is a starting point for all queries that return one or more Users.
	usersQuery squirrel.SelectBuilder
}

func (us *SqlUserStore) ClearCaches() {}

func NewSqlUserStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.UserStore {
	us := &SqlUserStore{
		Store:   sqlStore,
		metrics: metrics,
	}

	// note: we are providing field names explicitly here to maintain order of columns (needed when using raw queries)
	us.usersQuery = us.
		GetQueryBuilder().
		Select(us.ModelFields("u.")...).
		From(store.UserTableName + " u")

	return us
}

func (us *SqlUserStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Email",
		"Username",
		"FirstName",
		"LastName",
		"DefaultShippingAddressID",
		"DefaultBillingAddressID",
		"Password",
		"AuthData",
		"AuthService",
		"EmailVerified",
		"Nickname",
		"Roles",
		"Props",
		"NotifyProps",
		"LastPasswordUpdate",
		"LastPictureUpdate",
		"FailedAttempts",
		"Locale",
		"Timezone",
		"MfaActive",
		"MfaSecret",
		"CreateAt",
		"UpdateAt",
		"DeleteAt",
		"IsActive",
		"Note",
		"JwtTokenKey",
		"LastActivityAt",
		"TermsOfServiceId",
		"TermsOfServiceCreateAt",
		"DisableWelcomeEmail",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (us *SqlUserStore) ScanFields(user *model.User) []interface{} {
	return []interface{}{
		&user.Id,
		&user.Email,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.DefaultShippingAddressID,
		&user.DefaultBillingAddressID,
		&user.Password,
		&user.AuthData,
		&user.AuthService,
		&user.EmailVerified,
		&user.Nickname,
		&user.Roles,
		&user.Props,
		&user.NotifyProps,
		&user.LastPasswordUpdate,
		&user.LastPictureUpdate,
		&user.FailedAttempts,
		&user.Locale,
		&user.Timezone,
		&user.MfaActive,
		&user.MfaSecret,
		&user.CreateAt,
		&user.UpdateAt,
		&user.DeleteAt,
		&user.IsActive,
		&user.Note,
		&user.JwtTokenKey,
		&user.LastActivityAt,
		&user.TermsOfServiceId,
		&user.TermsOfServiceCreateAt,
		&user.DisableWelcomeEmail,
		&user.Metadata,
		&user.PrivateMetadata,
	}
}

// TODO: remove this
func (us *SqlUserStore) GetUnreadCount(userID string) (int64, error) {
	panic("not implemented")
}

// DeactivateGuests
func (us *SqlUserStore) DeactivateGuests() ([]string, error) {
	curTime := model.GetMillis()
	updateQuery := us.
		GetQueryBuilder().
		Update(store.UserTableName).
		Set("UpdateAt", curTime).
		Set("DeleteAt", curTime).
		Where(squirrel.Eq{"Roles": "system_guest"}).
		Where(squirrel.Eq{"DeleteAt": 0})

	queryString, args, err := updateQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "deactivate_guests_tosql")
	}

	_, err = us.GetMasterX().Exec(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Users with roles=system_guest")
	}

	selectQuery := us.
		GetQueryBuilder().
		Select("Id").
		From(store.UserTableName).
		Where(squirrel.Eq{"DeleteAt": curTime})

	queryString, args, err = selectQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "deactivate_guests_tosql")
	}

	userIds := []string{}
	err = us.GetMasterX().Select(&userIds, queryString, args...)
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
		builder := us.GetQueryBuilder().
			Select("COUNT(*)").
			From("Users").
			Where(whereEquals)
		query, args, err := builder.ToSql()
		if err != nil {
			return 0, errors.Wrap(err, "select_count_users_tosql")
		}
		var numAffected int
		err = us.GetReplicaX().Get(&numAffected, query, args...)
		return numAffected, err
	}
	builder := us.GetQueryBuilder().
		Update("Users").
		Set("AuthData", "Email").
		Where(whereEquals)
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "update_users_tosql")
	}
	result, err := us.GetMasterX().Exec(query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to update users' AuthData")
	}
	numAffected, err := result.RowsAffected()
	return int(numAffected), err
}

func (us *SqlUserStore) InvalidateProfileCacheForUser(userId string) {}

func (us *SqlUserStore) GetEtagForProfiles(teamId string) string {
	panic("not implemented")
}

func (us *SqlUserStore) GetEtagForAllProfiles() string {
	var updateAt int64
	err := us.GetReplicaX().Get(&updateAt, "SELECT UpdateAt FROM "+store.UserTableName+" ORDER BY UpdateAt DESC LIMIT 1")
	if err != nil {
		return fmt.Sprintf("%v.%v", model.CurrentVersion, model.GetMillis())
	}
	return fmt.Sprintf("%v.%v", model.CurrentVersion, updateAt)
}

func (us *SqlUserStore) Save(user *model.User) (*model.User, error) {
	user.PreSave()
	if err := user.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.UserTableName + " (" + us.ModelFields("").Join(",") + ") VALUES (" + us.ModelFields(":").Join(",") + ")"
	if _, err := us.GetMasterX().NamedExec(query, user); err != nil {
		if us.IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique"}) {
			return nil, store.NewErrInvalidInput("User", "email", user.Email)
		}
		if us.IsUniqueConstraintError(err, []string{"Username", "users_username_key", "idx_users_username_unique"}) {
			return nil, store.NewErrInvalidInput("User", "username", user.Username)
		}
		return nil, errors.Wrapf(err, "failed to save User with userId=%s", user.Id)
	}

	return user, nil
}

// Update updates user
func (us *SqlUserStore) Update(user *model.User, trustedUpdateData bool) (*model.UserUpdate, error) {
	user.PreUpdate()
	if err := user.IsValid(); err != nil {
		return nil, err
	}
	var oldUser model.User

	err := us.GetMasterX().Get(&oldUser, "SELECT * FROM "+store.UserTableName+" WHERE Id = ?", user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.UserTableName, user.Id)
		}
		return nil, errors.Wrapf(err, "failed to get user with userId=%s", user.Id)
	}

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

	if user.IsOAuthUser() && !trustedUpdateData {
		user.Email = oldUser.Email
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

	query := "UPDATE " + store.UserTableName + " SET " +
		us.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

	result, err := us.GetMasterX().NamedExec(query, user)
	if err != nil {
		if us.IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique"}) {
			return nil, store.NewErrInvalidInput("User", "id", user.Id)
		}
		if us.IsUniqueConstraintError(err, []string{"Username", "users_username_key", "idx_users_username_unique"}) {
			return nil, store.NewErrInvalidInput("User", "id", user.Id)
		}
		return nil, errors.Wrapf(err, "failed to update User with userId=%s", user.Id)
	}

	count, _ := result.RowsAffected()
	if count > 1 {
		return nil, fmt.Errorf("multiple users were update: userId=%s, count=%d", user.Id, count)
	}

	user.Sanitize(map[string]bool{})
	oldUser.Sanitize(map[string]bool{})
	return &model.UserUpdate{New: user, Old: &oldUser}, nil
}

func (us *SqlUserStore) UpdateLastPictureUpdate(userId string) error {
	now := model.GetMillis()
	if _, err := us.GetMasterX().Exec("UPDATE "+store.UserTableName+" SET LastPictureUpdate = ?, UpdateAt = ? WHERE Id = ?", now, now, userId); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) ResetLastPictureUpdate(userId string) error {
	now := model.GetMillis()
	if _, err := us.GetMasterX().Exec("UPDATE "+store.UserTableName+" SET LastPictureUpdate = ?, UpdateAt = ? WHERE Id = ?", 0, now, userId); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) UpdateUpdateAt(userId string) (int64, error) {
	now := model.GetMillis()
	if _, err := us.GetMasterX().Exec("UPDATE "+store.UserTableName+" SET UpdateAt = ? WHERE Id = ?", now, userId); err != nil {
		return now, errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return now, nil
}

func (us *SqlUserStore) UpdatePassword(userId, hashedPassword string) error {
	now := model.GetMillis()
	if _, err := us.GetMasterX().
		Exec("UPDATE "+store.UserTableName+" SET Password = ?, LastPasswordUpdate = ?, UpdateAt = ?, AuthData = NULL, AuthService = '', FailedAttempts = 0 WHERE Id = ?",
			hashedPassword,
			now,
			now,
			userId,
		); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) UpdateFailedPasswordAttempts(userId string, attempts int) error {
	if _, err := us.GetMasterX().Exec(
		"UPDATE "+store.UserTableName+" SET FailedAttempts = ? WHERE Id = ?",
		attempts, userId,
	); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}
	return nil
}

// UpdateAuthData updates auth data of user
func (us *SqlUserStore) UpdateAuthData(userId string, service string, authData *string, email string, resetMfa bool) (string, error) {
	updateAt := model.GetMillis()
	query := `UPDATE ` + store.UserTableName +
		` SET
			Password = '',
			LastPasswordUpdate = ?,
			UpdateAt = ?,
			FailedAttempts = 0,
			AuthService = ?,
			AuthData = ?
	`
	if email != "" {
		query += ", Email = lower(?)"
	}
	if resetMfa {
		query += ", MfaActive = false, MfaSecret = ''"
	}
	query += " WHERE Id = ?"

	if _, err := us.
		GetMasterX().Exec(query, updateAt, updateAt, service, authData, email, userId); err != nil {
		if us.IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique", "AuthData", "users_authdata_key"}) {
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
		GetMasterX().Exec("UPDATE "+store.UserTableName+" SET MfaSecret = ?, UpdateAt = ? WHERE Id = ?", secret, updateAt, userId); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) UpdateMfaActive(userId string, active bool) error {
	updateAt := model.GetMillis()
	if _, err := us.GetMasterX().Exec("UPDATE "+store.UserTableName+" SET MfaActive = ?, UpdateAt = ? WHERE Id = ?", active, updateAt, userId); err != nil {
		return errors.Wrapf(err, "failed to update User with userId=%s", userId)
	}

	return nil
}

func (us *SqlUserStore) GetProfileByIds(ctx context.Context, userIds []string, options *store.UserGetByIdsOpts, allowFromCache bool) ([]*model.User, error) {
	if options == nil {
		options = new(store.UserGetByIdsOpts)
	}

	users := []*model.User{}
	query := us.usersQuery.Where(squirrel.Eq{"u.Id": userIds}).OrderBy("u.Username ASC")

	if options.Since > 0 {
		query = query.Where(squirrel.Gt{"u.UpdateAt": options.Since})
	}

	query = applyViewRestrictionsFilter(query, true)
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_profile_by_ids_tosql")
	}

	if err := us.DBXFromContext(ctx).Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	return users, nil
}

func (us *SqlUserStore) GetSystemAdminProfiles() (map[string]*model.User, error) {
	query := us.usersQuery.Where("Roles LIKE ?", "%system_admin%").OrderBy("u.Username ASC")
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_system_admin_profiles_tosql")
	}

	var users []*model.User
	if err := us.GetReplicaX().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Users")
	}

	userMap := make(map[string]*model.User)
	for _, u := range users {
		u.Sanitize(map[string]bool{})
		userMap[u.Id] = u
	}

	return userMap, nil
}

func (us *SqlUserStore) GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail bool) (*model.User, error) {
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

	users := []*model.User{}
	if err := us.GetReplicaX().Select(&users, queryString, args...); err != nil {
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
		GetMasterX().
		Exec("UPDATE "+store.UserTableName+" SET Email = lower(?), EmailVerified = true, UpdateAt = ? WHERE Id = ?",
			email, now, userId,
		); err != nil {
		return "", errors.Wrapf(err, "failed to update Users with userId=%s and email=%s", userId, email)
	}

	return userId, nil
}

func (us *SqlUserStore) PermanentDelete(userId string) error {
	if _, err := us.GetMasterX().Exec("DELETE FROM "+store.UserTableName+" WHERE Id = ?", userId); err != nil {
		return errors.Wrapf(err, "failed to delete User with userId=%s", userId)
	}
	return nil
}

// applyViewRestrictionsFilter add "DISTINCT" to query if given distinct is `true`
func applyViewRestrictionsFilter(query squirrel.SelectBuilder, distinct bool) squirrel.SelectBuilder {
	if distinct {
		return query.Distinct()
	}

	return query
}

func (us *SqlUserStore) Count(options model.UserCountOptions) (int64, error) {
	query := us.GetQueryBuilder().Select("COUNT(DISTINCT u.Id)").From(store.UserTableName + " AS u")

	if !options.IncludeDeleted {
		query = query.Where("u.DeleteAt = 0")
	}

	query = applyViewRestrictionsFilter(query, true)
	query = applyMultiRoleFilters(query, options.Roles)

	queryString, args, err := query.ToSql()
	if err != nil {
		return int64(0), errors.Wrap(err, "count_tosql")
	}

	var count int64
	err = us.GetReplicaX().Get(&count, queryString, args...)
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count Users")
	}
	return count, nil
}

func (us *SqlUserStore) AnalyticsActiveCount(timePeriod int64, options model.UserCountOptions) (int64, error) {
	time := model.GetMillis() - timePeriod
	query := us.GetQueryBuilder().Select("COUNT(*)").From("Status AS s").Where("LastActivityAt > ?", time)

	if !options.IncludeDeleted {
		query = query.LeftJoin("Users ON s.UserId = Users.Id").Where("Users.DeleteAt = 0")
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "analytics_active_count_tosql")
	}

	var v int64
	err = us.GetReplicaX().Get(&v, queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count Users")
	}
	return v, nil
}

func (us *SqlUserStore) AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options model.UserCountOptions) (int64, error) {
	query := us.GetQueryBuilder().Select("COUNT(*)").From("Status AS s").Where("LastActivityAt > ? AND LastActivityAt <= ?", startTime, endTime)
	if !options.IncludeDeleted {
		query = query.LeftJoin("Users ON s.UserId = Users.Id").Where("Users.DeleteAt = 0")
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to build query.")
	}

	var v int64
	err = us.GetReplicaX().Get(&v, queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to get the active users during the requested period")
	}
	return v, nil
}

func applyMultiRoleFilters(query squirrel.SelectBuilder, systemRoles []string) squirrel.SelectBuilder {
	sqOr := squirrel.Or{}

	if len(systemRoles) > 0 && systemRoles[0] != "" {
		for _, role := range systemRoles {
			queryRole := store.WildcardSearchTerm(role)
			switch role {
			case model.SystemUserRoleId:
				// If querying for a `system_user` ensure that the user is only a system_user.
				sqOr = append(sqOr, squirrel.Eq{"u.Roles": role})
			case model.SystemAdminRoleId,
				model.SystemUserManagerRoleId,
				model.SystemReadOnlyAdminRoleId,
				model.SystemManagerRoleId:
				// If querying for any other roles search using a wildcard
				sqOr = append(sqOr, squirrel.ILike{"u.Roles": queryRole})
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

func (us *SqlUserStore) Search(term string, options *model.UserSearchOptions) ([]*model.User, error) {
	query := us.usersQuery.
		OrderBy("Username ASC").
		Limit(uint64(options.Limit))
	return us.performSearch(query, term, options)
}

func applyRoleFilter(query squirrel.SelectBuilder, role string) squirrel.SelectBuilder {
	if role == "" {
		return query
	}

	roleParam := fmt.Sprintf("%%%s%%", store.SanitizeSearchTerm(role, "\\"))
	return query.Where("u.Roles LIKE LOWER(?)", roleParam)
}

func (us *SqlUserStore) performSearch(query squirrel.SelectBuilder, term string, options *model.UserSearchOptions) ([]*model.User, error) {
	term = store.SanitizeSearchTerm(term, "*")

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

	query = applyRoleFilter(query, options.Role)
	query = applyMultiRoleFilters(query, options.Roles)

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

	var users []*model.User
	if err := us.GetReplicaX().Select(&users, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Users with term=%s and searchType=%v", term, searchType)
	}
	for _, u := range users {
		u.Sanitize(map[string]bool{})
	}

	return users, nil
}

func (us *SqlUserStore) AnalyticsGetInactiveUsersCount() (int64, error) {
	var count int64
	err := us.GetReplicaX().Get(&count, "SELECT COUNT(Id) FROM Users WHERE DeleteAt > 0")
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count inactive Users")
	}
	return count, nil
}

func (us *SqlUserStore) AnalyticsGetExternalUsers(hostDomain string) (bool, error) {
	var count int64
	err := us.GetReplicaX().Get(&count, "SELECT COUNT(Id) FROM Users WHERE LOWER(Email) NOT LIKE ?", "%@"+strings.ToLower(hostDomain))
	if err != nil {
		return false, errors.Wrap(err, "failed to count inactive Users")
	}
	return count > 0, nil
}

func (us *SqlUserStore) AnalyticsGetGuestCount() (int64, error) {
	var count int64
	err := us.GetReplicaX().Get(&count, "SELECT count(*) FROM Users WHERE Roles LIKE ? and DeleteAt = 0", "%system_guest%")
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count guest Users")
	}
	return count, nil
}

func (us *SqlUserStore) AnalyticsGetSystemAdminCount() (int64, error) {
	var count int64
	err := us.GetReplicaX().Get(&count, "SELECT count(*) FROM Users WHERE Roles LIKE ? and DeleteAt = 0", "%system_admin%")
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count system admin Users")
	}
	return count, nil
}

func (us *SqlUserStore) ClearAllCustomRoleAssignments() error {
	builtinRoles := model.MakeDefaultRoles()
	lastUserId := strings.Repeat("0", len(model.NewId()))

	for {
		var transaction store_iface.SqlxTxExecutor
		var err error

		if transaction, err = us.GetMasterX().Beginx(); err != nil {
			return errors.Wrap(err, "begin_transaction")
		}
		defer store.FinalizeTransaction(transaction)

		var users []*model.User
		if err := transaction.Select(&users, "SELECT * FROM Users WHERE Id > ? ORDER BY Id LIMIT 1000", lastUserId); err != nil {
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
				if _, err := transaction.Exec("UPDATE Users SET Roles = ? WHERE Id = ?", newRolesString, user.Id); err != nil {
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
	var createAt int64
	err := us.GetReplicaX().Get(&createAt, "SELECT CreateAt FROM Users WHERE CreateAt IS NOT NULL ORDER BY CreateAt ASC LIMIT 1")
	if err != nil {
		return 0, errors.Wrap(err, "failed to infer system install date")
	}

	return createAt, nil
}

func (us *SqlUserStore) GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*model.UserForIndexing, error) {
	panic("not implemented")
}

func (us *SqlUserStore) PromoteGuestToUser(userId string) error {
	panic("not implemented")
}

func (us *SqlUserStore) DemoteUserToGuest(userID string) (*model.User, error) {
	panic("not implemented")
}

func (us *SqlUserStore) GetKnownUsers(userId string) ([]string, error) {
	panic("not implemented")
}

func (us *SqlUserStore) GetAllProfiles(options *model.UserGetOptions) ([]*model.User, error) {
	panic("not implemented")
}

func (s *SqlUserStore) commonSelectQueryBuilder(options *model.UserFilterOptions) squirrel.SelectBuilder {
	query := s.
		GetQueryBuilder().
		Select(s.ModelFields(store.UserTableName + ".")...).
		From(store.UserTableName)

	for _, opt := range []squirrel.Sqlizer{
		options.Id,
		options.Email,
		options.Username,
		options.FirstName,
		options.LastName,
		options.OrderID,
		options.AuthData,
		options.AuthService,
		options.Extra,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if options.HasNoOrder {
		query = query.
			LeftJoin(store.OrderTableName + " ON Orders.UserID = Users.Id").
			Where("Orders.UserID IS NULL")
	}
	if options.ExcludeBoardMembers {
		query = query.
			LeftJoin(store.ShopStaffTableName + " ON ShopStaffs.StaffID = Users.Id").
			Where("hopStaffs.StaffID IS NULL")
	}
	if options.Limit > 0 {
		query = query.Limit(uint64(options.Limit))
	}
	if options.OrderBy != "" {
		query = query.OrderBy(options.OrderBy)
	}

	return query
}

func (s *SqlUserStore) FilterByOptions(ctx context.Context, options *model.UserFilterOptions) ([]*model.User, error) {
	queryString, args, err := s.commonSelectQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var users []*model.User
	err = s.GetReplicaX().Select(&users, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find users with given options")
	}

	return users, nil
}

func (s *SqlUserStore) GetByOptions(ctx context.Context, options *model.UserFilterOptions) (*model.User, error) {
	queryString, args, err := s.commonSelectQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var user model.User
	err = s.DBXFromContext(ctx).Get(&user, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.UserTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find users with given options")
	}

	return &user, nil
}
