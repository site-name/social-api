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
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	UserSearchTypeNames_NO_FULL_NAME = []string{"Username", "Nickname"}
	UserSearchTypeNames              = []string{"Username", "FirstName", "LastName", "Nickname"}
	UserSearchTypeAll_NO_FULL_NAME   = []string{"Username", "Nickname", "Email"}
	UserSearchTypeAll                = []string{"Username", "FirstName", "LastName", "Nickname", "Email"}
)

var (
	UserUniqueColumnConstraint = []string{
		"users_auth_data_key",
		"users_email_key",
		"users_username_key",
	}
)

type SqlUserStore struct {
	store.Store
	metrics einterfaces.MetricsInterface

	// usersQuery is a starting point for all queries that return one or more Users.
	// usersQuery squirrel.SelectBuilder
}

func NewSqlUserStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.UserStore {
	us := &SqlUserStore{
		Store:   sqlStore,
		metrics: metrics,
	}
	return us
}

func (us *SqlUserStore) ClearCaches() {}

func (us *SqlUserStore) Get(conds ...qm.QueryMod) (*model.User, error) {
	user, err := model.Users(conds...).One(us.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Users, "conds")
		}
		return nil, err
	}

	return user, nil
}

func (us *SqlUserStore) Find(conds ...qm.QueryMod) (model.UserSlice, error) {
	return model.Users(conds...).All(us.GetReplica())
}

// ResetAuthDataToEmailForUsers resets the AuthData of users whose AuthService
// is |service| to their Email. If userIDs is non-empty, only the users whose
// IDs are in userIDs will be affected. If dryRun is true, only the number
// of users who *would* be affected is returned; otherwise, the number of
// users who actually were affected is returned.
func (us *SqlUserStore) ResetAuthDataToEmailForUsers(service string, userIDs []string, includeDeleted bool, dryRun bool) (int, error) {
	queryMods := []qm.QueryMod{
		model.UserWhere.AuthService.EQ(service),
	}
	if len(userIDs) > 0 {
		queryMods = append(queryMods, model.UserWhere.ID.IN(userIDs))
	}
	if !includeDeleted {
		queryMods = append(queryMods, model.UserWhere.DeleteAt.EQ(0))
	}

	if dryRun {
		count, err := model.Users(queryMods...).Count(us.GetReplica())
		return int(count), err
	}

	numUpdated, err := model.Users(queryMods...).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.AuthData:  "Email",
			model.UserColumns.UpdatedAt: model_helper.GetMillis(),
		})
	return int(numUpdated), err
}

func (us *SqlUserStore) InvalidateProfileCacheForUser(userId string) {}

func (us *SqlUserStore) GetEtagForProfiles() string {
	var updatedAt int64
	err := model.
		Users(
			qm.Select(model.UserColumns.UpdatedAt),
			qm.OrderBy(fmt.Sprintf("%s DESC", model.UserColumns.UpdatedAt)),
			qm.Limit(1),
		).
		QueryRowContext(us.Context(), us.GetReplica()).
		Scan(&updatedAt)
	if err != nil {
		return fmt.Sprintf("%v.%v", model_helper.CurrentVersion, model_helper.GetMillis())
	}

	return fmt.Sprint("%v.%v", model_helper.CurrentVersion, updatedAt)
}

func (us *SqlUserStore) GetEtagForAllProfiles() string {
	user, err := model.
		Users(qm.OrderBy(model.UserColumns.UpdatedAt + " DESC")).
		One(us.GetReplica())
	if err != nil {
		return fmt.Sprintf("%v.%v", model_helper.CurrentVersion, model_helper.GetMillis())
	}
	return fmt.Sprintf("%v.%v", model_helper.CurrentVersion, user.UpdatedAt)
}

func (us *SqlUserStore) Save(user model.User) (*model.User, error) {
	model_helper.UserPreSave(&user)
	if err := model_helper.UserIsValid(user); err != nil {
		return nil, err
	}

	err := user.Insert(us.GetMaster(), boil.Infer())
	if err != nil {
		if us.IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique"}) {
			return nil, store.NewErrInvalidInput("User", "email", user.Email)
		}
		if us.IsUniqueConstraintError(err, []string{"Username", "users_username_key", "idx_users_username_unique"}) {
			return nil, store.NewErrInvalidInput("User", "username", user.Username)
		}
		return nil, err
	}

	return &user, nil
}

func (us *SqlUserStore) Update(user model.User, trustedUpdateData bool) (*model_helper.UserUpdate, error) {
	model_helper.UserPreUpdate(&user)
	if err := model_helper.UserIsValid(user); err != nil {
		return nil, err
	}

	oldUser, err := model.FindUser(us.GetReplica(), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Users, user.ID)
		}
		return nil, err
	}

	blackListedColumns := []string{
		model.UserColumns.CreatedAt,
		model.UserColumns.AuthData,
		model.UserColumns.AuthService,
		model.UserColumns.Password,
		model.UserColumns.LastPasswordUpdate,
		model.UserColumns.LastPictureUpdate,
		model.UserColumns.FailedAttempts,
		model.UserColumns.MfaSecret,
		model.UserColumns.MfaActive,
	}

	if !trustedUpdateData {
		blackListedColumns = append(blackListedColumns, model.UserColumns.Roles, model.UserColumns.DeleteAt)
	}

	if model_helper.UserIsOauth(user) && !trustedUpdateData {
		blackListedColumns = append(blackListedColumns, model.UserColumns.Email)
	} else if model_helper.UserIsLDAP(user) && !trustedUpdateData {
		if user.Username != oldUser.Username ||
			user.Email != oldUser.Email {
			return nil, store.NewErrInvalidInput(model.TableNames.Users, "id", user.ID)
		}
	}
	if user.Email == oldUser.Email {
		blackListedColumns = append(blackListedColumns, model.UserColumns.EmailVerified)
	} else {
		user.EmailVerified = false
	}

	if user.Username != oldUser.Username {
		model_helper.UserUpdateMentionKeysFromUsername(&user, oldUser.Username)
	}

	_, err = user.Update(us.GetMaster(), boil.Blacklist(blackListedColumns...))
	if err != nil {
		if us.IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Users, model.UserColumns.ID, user.Email)
		}
		if us.IsUniqueConstraintError(err, []string{"Username", "users_username_key", "idx_users_username_unique"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Users, model.UserColumns.ID, user.Username)
		}
	}

	model_helper.UserSanitize(oldUser, map[string]bool{})
	model_helper.UserSanitize(&user, map[string]bool{})

	return &model_helper.UserUpdate{Old: oldUser, New: &user}, nil
}

func (us *SqlUserStore) UpdateLastPictureUpdate(userId string, updateMillis int64) error {
	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.LastPictureUpdate: updateMillis,
		})
	return err
}

func (us *SqlUserStore) ResetLastPictureUpdate(userId string) error {
	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.LastPictureUpdate: model_helper.GetMillis(),
		})
	return err
}

func (us *SqlUserStore) UpdateUpdateAt(userId string) (int64, error) {
	now := model_helper.GetMillis()
	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.UpdatedAt: now,
		})
	return now, err
}

func (us *SqlUserStore) UpdatePassword(userId, hashedPassword string) error {
	now := model_helper.GetMillis()
	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.Password:           hashedPassword,
			model.UserColumns.LastPasswordUpdate: now,
			model.UserColumns.AuthData:           nil,
			model.UserColumns.AuthService:        "",
			model.UserColumns.FailedAttempts:     0,
		})
	return err
}

func (us *SqlUserStore) UpdateFailedPasswordAttempts(userId string, attempts int) error {
	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.FailedAttempts: attempts,
		})
	return err
}

func (us *SqlUserStore) UpdateAuthData(userId string, service string, authData *string, email string, resetMfa bool) (string, error) {
	updateAt := model_helper.GetMillis()
	var updateColumns = model.M{
		model.UserColumns.Password:           "",
		model.UserColumns.LastPasswordUpdate: updateAt,
		model.UserColumns.UpdatedAt:          updateAt,
		model.UserColumns.FailedAttempts:     0,
		model.UserColumns.AuthService:        service,
		model.UserColumns.AuthData:           authData,
	}
	if email != "" {
		updateColumns[model.UserColumns.Email] = fmt.Sprintf("lower(%s)", email)
	}
	if resetMfa {
		updateColumns[model.UserColumns.MfaActive] = false
		updateColumns[model.UserColumns.MfaSecret] = ""
	}

	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), updateColumns)
	if err != nil {
		if us.IsUniqueConstraintError(err, []string{"Email", "users_email_key", "idx_users_email_unique", "AuthData", "users_authdata_key"}) {
			return "", store.NewErrInvalidInput("User", "id", userId)
		}
		return "", err
	}
	return userId, nil
}

// UpdateMfaSecret updates mfa secret for current user
func (us *SqlUserStore) UpdateMfaSecret(userId, secret string) error {
	updateAt := model_helper.GetMillis()
	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.MfaSecret: secret,
			model.UserColumns.UpdatedAt: updateAt,
		})
	return err
}

func (us *SqlUserStore) UpdateMfaActive(userId string, active bool) error {
	updateAt := model_helper.GetMillis()
	_, err := model.
		Users(model.UserWhere.ID.EQ(userId)).
		UpdateAll(us.GetMaster(), model.M{
			model.UserColumns.MfaActive: active,
			model.UserColumns.UpdatedAt: updateAt,
		})
	return err
}

func (us *SqlUserStore) GetProfileByIds(ctx context.Context, userIds []string, options store.UserGetByIdsOpts, allowFromCache bool) (model.UserSlice, error) {
	queryMods := []qm.QueryMod{
		model.UserWhere.ID.IN(userIds),
		qm.OrderBy(model.UserColumns.Username + " ASC"),
		qm.Distinct(model.UserColumns.ID),
	}
	if options.Since > 0 {
		queryMods = append(queryMods, model.UserWhere.UpdatedAt.GT(options.Since))
	}

	return model.Users(queryMods...).All(us.GetReplica())
}

func (us *SqlUserStore) GetSystemAdminProfiles() (map[string]*model.User, error) {
	users, err := model.
		Users(model.UserWhere.Roles.LIKE("%system_admin%"), qm.OrderBy(model.UserColumns.Username+" ASC")).
		All(us.GetReplica())
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*model.User)
	for _, u := range users {
		model_helper.UserSanitize(u, map[string]bool{})
		userMap[u.ID] = u
	}

	return userMap, nil
}

func (us *SqlUserStore) GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail bool) (*model.User, error) {
	var queryMod qm.QueryMod
	loginId = strings.ToLower(loginId)

	if allowSignInWithUsername && allowSignInWithEmail {
		queryMod = qm.Where(fmt.Sprintf("%s = ? OR %s = ?", model.UserTableColumns.Email, model.UserTableColumns.Username), loginId, loginId)
	} else if allowSignInWithUsername {
		queryMod = model.UserWhere.Username.EQ(loginId)
	} else if allowSignInWithEmail {
		queryMod = model.UserWhere.Email.EQ(loginId)
	} else {
		return nil, errors.New("sign in with username and email are disabled")
	}

	user, err := model.Users(queryMod).One(us.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Users, loginId)
		}
		return nil, err
	}
	return user, nil
}

func (us *SqlUserStore) VerifyEmail(userId, email string) (string, error) {
	now := model_helper.GetMillis()

	_, err := model.Users(model.UserWhere.ID.EQ(userId)).UpdateAll(us.GetMaster(), model.M{
		model.UserColumns.Email:         strings.ToLower(email),
		model.UserColumns.EmailVerified: true,
		model.UserColumns.UpdatedAt:     now,
	})
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (us *SqlUserStore) PermanentDelete(userId string) error {
	_, err := model.Users(model.UserWhere.ID.EQ(userId)).DeleteAll(us.GetMaster())
	return err
}

// applyViewRestrictionsFilter add "DISTINCT" to query if given distinct is `true`
func applyViewRestrictionsFilter(query squirrel.SelectBuilder, distinct bool) squirrel.SelectBuilder {
	if distinct {
		return query.Distinct()
	}

	return query
}

func (us *SqlUserStore) Count(options model_helper.UserCountOptions) (int64, error) {
	query := us.GetQueryBuilder().
		Select(fmt.Sprintf("COUNT(DISTINCT %s)", model.UserTableColumns.ID)).
		From(model.TableNames.Users)

	if !options.IncludeDeleted {
		query = query.Where(squirrel.Eq{model.UserTableColumns.DeleteAt: 0})
	}

	query = applyViewRestrictionsFilter(query, true)
	query = applyMultiRoleFilters(query, options.Roles)

	queryString, args, err := query.ToSql()
	if err != nil {
		return int64(0), errors.Wrap(err, "count_tosql")
	}

	var count int64
	err = queries.Raw(queryString, args...).QueryRow(us.GetReplica()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (us *SqlUserStore) AnalyticsActiveCount(timePeriod int64, options model_helper.UserCountOptions) (int64, error) {
	period := model_helper.GetMillis() - timePeriod

	queryMods := []qm.QueryMod{
		model.StatusWhere.LastActivityAt.GT(period),
	}
	if !options.IncludeDeleted {
		queryMods = append(
			queryMods,
			qm.LeftOuterJoin(fmt.Sprintf(
				"%[1]s ON %[2]s = %[3]s",
				model.TableNames.Users,          // 1
				model.UserTableColumns.ID,       // 2
				model.StatusTableColumns.UserID, // 3
			)),
			model.UserWhere.DeleteAt.EQ(0),
		)
	}

	return model.Statuses(queryMods...).Count(us.GetReplica())
}

func (us *SqlUserStore) AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options model_helper.UserCountOptions) (int64, error) {
	queryMods := []qm.QueryMod{
		model.StatusWhere.LastActivityAt.GT(startTime),
		model.StatusWhere.LastActivityAt.LTE(endTime),
	}
	if !options.IncludeDeleted {
		queryMods = append(
			queryMods,
			qm.LeftOuterJoin(fmt.Sprintf(
				"%[1]s ON %[2]s = %[3]s",
				model.TableNames.Users,          // 1
				model.UserTableColumns.ID,       // 2
				model.StatusTableColumns.UserID, // 3
			)),
			model.UserWhere.DeleteAt.EQ(0),
		)
	}

	return model.Users(queryMods...).Count(us.GetReplica())
}

func applyMultiRoleFilters(query squirrel.SelectBuilder, systemRoles []string) squirrel.SelectBuilder {
	sqOr := squirrel.Or{}

	for _, role := range systemRoles {
		queryRole := store.WildcardSearchTerm(role)
		switch role {
		case model_helper.SystemUserRoleId:
			// If querying for a `system_user` ensure that the user is only a system_user.
			sqOr = append(sqOr, squirrel.Eq{model.UserTableColumns.Roles: role})
		case model_helper.SystemAdminRoleId,
			model_helper.SystemUserManagerRoleId,
			model_helper.SystemReadOnlyAdminRoleId,
			model_helper.SystemManagerRoleId:
			// If querying for any other roles search using a wildcard
			sqOr = append(sqOr, squirrel.ILike{model.UserTableColumns.Roles: queryRole})
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

func applyRoleFilter(query squirrel.SelectBuilder, role string) squirrel.SelectBuilder {
	if role == "" {
		return query
	}

	roleParam := fmt.Sprintf("%%%s%%", store.SanitizeSearchTerm(role, "\\"))
	return query.Where("u.Roles LIKE LOWER(?)", roleParam)
}

func (us *SqlUserStore) Search(term string, options *model_helper.UserSearchOptions) (model.UserSlice, error) {
	query := us.GetQueryBuilder().
		Select("*").
		From(model.TableNames.Users).
		OrderBy(fmt.Sprintf("%S ASC", model.UserColumns.Username)).
		Limit(uint64(options.Limit))
	return us.performSearch(query, term, options)
}

func (us *SqlUserStore) performSearch(query squirrel.SelectBuilder, term string, options *model_helper.UserSearchOptions) (model.UserSlice, error) {
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
		query = query.Where(squirrel.Eq{model.UserColumns.DeleteAt: 0})
	}

	if strings.TrimSpace(term) != "" {
		query = generateSearchQuery(query, strings.Fields(term), searchType)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "perform_search_tosql")
	}

	var users model.UserSlice
	err = queries.Raw(queryString, args...).Bind(us.Context(), us.GetReplica(), &users)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		model_helper.UserSanitize(user, map[string]bool{})
	}
	return users, nil
}

func (us *SqlUserStore) AnalyticsGetInactiveUsersCount() (int64, error) {
	return model.
		Users(model.UserWhere.DeleteAt.GT(0)).
		Count(us.GetReplica())
}

func (us *SqlUserStore) AnalyticsGetExternalUsers(hostDomain string) (bool, error) {
	count, err := model.
		Users(qm.Where("LOWER(email) NOT LIKE ?", "%@"+strings.ToLower(hostDomain))).
		Count(us.GetReplica())
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (us *SqlUserStore) AnalyticsGetGuestCount() (int64, error) {
	return model.
		Users(
			model.UserWhere.Roles.LIKE("%"+model_helper.SystemGuestRoleId+"%"),
			model.UserWhere.DeleteAt.EQ(0),
		).
		Count(us.GetReplica())
}

func (us *SqlUserStore) AnalyticsGetSystemAdminCount() (int64, error) {
	return model.
		Users(
			model.UserWhere.Roles.LIKE("%"+model_helper.SystemAdminRoleId+"%"),
			model.UserWhere.DeleteAt.EQ(0),
		).
		Count(us.GetReplica())
}

func (us *SqlUserStore) ClearAllCustomRoleAssignments() error {
	builtinRoles := model_helper.MakeDefaultRoles()
	lastUserId := strings.Repeat("0", len(model_helper.NewId()))

	for {
		tx, err := us.GetMaster().BeginTx(us.Context(), &sql.TxOptions{})
		if err != nil {
			return errors.Wrap(err, "begin_transaction")
		}
		defer us.FinalizeTransaction(tx)

		users, err := model.
			Users(
				model.UserWhere.ID.GT(lastUserId),
				qm.OrderBy(model.UserColumns.ID),
				qm.Limit(1000),
			).
			All(tx)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			lastUserId = user.ID
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
				_, err = user.Update(tx, boil.Whitelist(model.UserColumns.Roles))
				if err != nil {
					return err
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return errors.Wrap(err, "commit_transaction")
		}
	}

	return nil
}

func (us *SqlUserStore) InferSystemInstallDate() (int64, error) {
	var createAt int64
	return createAt, model.
		Users(
			qm.Select(model.UserColumns.CreatedAt),
			qm.OrderBy(model.UserColumns.CreatedAt),
			model.UserWhere.CreatedAt.NEQ(0),
			qm.Limit(1),
		).
		QueryRowContext(us.Context(), us.GetReplica()).
		Scan(&createAt)
}

func (us *SqlUserStore) GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*model_helper.UserForIndexing, error) {
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

func (us *SqlUserStore) GetAllProfiles(options model_helper.UserGetOptions) (model.UserSlice, error) {
	queryMods := []qm.QueryMod{}
	if options.Inactive {
		queryMods = append(queryMods, model.UserWhere.IsActive.EQ(false))
	} else if options.Active {
		queryMods = append(queryMods, model.UserWhere.IsActive.EQ(true))
	}

	if options.Role != "" {
		queryMods = append(queryMods, model.UserWhere.Roles.ILIKE(fmt.Sprintf("%LOWER(%s)%", options.Role)))
	}
	if options.Sort != "" {
		queryMods = append(queryMods, qm.OrderBy(options.Sort))
	} else {
		queryMods = append(queryMods, qm.OrderBy(model.UserTableColumns.Username+" ASC"))
	}

	queryMods = append(queryMods, qm.Offset(options.Page*options.PerPage), qm.Limit(options.PerPage))

	return model.Users(queryMods...).All(us.GetReplica())
}

func (s *SqlUserStore) IsEmpty() (bool, error) {
	exist, err := model.Users().Exists(s.GetReplica())
	if err != nil {
		return false, err
	}
	return !exist, nil
}

// func (s *SqlUserStore) commonSelectQueryBuilder(options *model.UserFilterOptions) squirrel.SelectBuilder {
// 	query := s.
// 		GetQueryBuilder().
// 		Select(model.UserTableName + ".*").
// 		From(model.UserTableName).
// 		Where(options.Conditions)

// 	if options.HasNoOrder {
// 		query = query.
// 			LeftJoin(model.OrderTableName + " ON Orders.UserID = Users.Id").
// 			Where("Orders.UserID IS NULL")

// 		goto orderBy
// 	}

// 	if options.OrderID != nil {
// 		query = query.
// 			InnerJoin(model.OrderTableName + " ON Orders.UserID = Users.Id").
// 			Where(options.OrderID)
// 	}

// 	if options.HasNoOrder || (options.AnnotateOrderCount && options.OrderCreatedDate == nil) {
// 		query = query.
// 			LeftJoin(model.OrderTableName + " ON Orders.UserID = Users.Id")

// 		if options.HasNoOrder {
// 			query = query.Where("Orders.UserID IS NULL")
// 		} else if options.AnnotateOrderCount {
// 			query = query.
// 				Column(`COUNT (Orders.Id) AS "Users.OrderCount"`).
// 				GroupBy(model.UserTableName + ".Id")
// 		}
// 	} else if options.OrderID != nil {
// 		query = query.
// 			InnerJoin(model.OrderTableName + " ON Orders.UserID = Users.Id").
// 			Where(options.OrderID)
// 	}
// 	if options.ExcludeBoardMembers {
// 		query = query.
// 			LeftJoin(model.ShopStaffTableName + " ON ShopStaffs.StaffID = Users.Id").
// 			Where("ShopStaffs.StaffID IS NULL")
// 	}

// 	if options.AnnotateOrderCount {
// 		query = query.
// 			LeftJoin(model.OrderTableName + " ON Orders.UserID = Users.Id").
// 			Column(`COUNT (Orders.Id) AS "Users.OrderCount"`).
// 			GroupBy(model.UserTableName + ".Id")
// 	} else if options.OrderCreatedDate != nil {

// 	}

// orderBy:
// 	if options.GraphqlPaginationValues.OrderBy != "" {
// 		query = query.OrderBy(options.GraphqlPaginationValues.OrderBy)
// 	}
// 	return query
// }

// func (s *SqlUserStore) FilterByOptions(ctx context.Context, options *model.UserFilterOptions) (int64, model.UserSlice, error) {
// 	query := s.commonSelectQueryBuilder(options)

// 	// count total if needed
// 	var totalUsersCount int64
// 	if options.CountTotal {
// 		countQuery, args, err := s.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
// 		if err != nil {
// 			return 0, nil, errors.Wrap(err, "FilterByOptions_CountTotal_ToSql")
// 		}

// 		err = s.GetReplica().Raw(countQuery, args...).Scan(&totalUsersCount).Error
// 		if err != nil {
// 			return 0, nil, errors.Wrap(err, "failed to count total number of users by given options")
// 		}
// 	}

// 	// apply pagination if needed
// 	// NOTE: we don't apply order by here since it's applied in commonQueryBuilder
// 	if options.GraphqlPaginationValues.PaginationApplicable() {
// 		query = query.
// 			Limit(options.GraphqlPaginationValues.Limit).
// 			Where(options.GraphqlPaginationValues.Condition)
// 	}

// 	queryString, args, err := query.ToSql()
// 	if err != nil {
// 		return 0, nil, errors.Wrap(err, "FilterByOptions_ToSql")
// 	}

// 	rows, err := s.GetReplica().Raw(queryString, args...).Rows()
// 	if err != nil {
// 		return 0, nil, errors.Wrap(err, "failed to find users by given options")
// 	}
// 	defer rows.Close()

// 	var users model.UserSlice
// 	for rows.Next() {
// 		var user model.User
// 		scanFields := s.ScanFields(&user)
// 		if options.AnnotateOrderCount {
// 			scanFields = append(scanFields, &user.OrderCount)
// 		}

// 		err := rows.Scan(scanFields...)
// 		if err != nil {
// 			return 0, nil, errors.Wrap(err, "failed to scan a row of user")
// 		}

// 		users = append(users, &user)
// 	}

// 	return totalUsersCount, users, nil
// }

// func (s *SqlUserStore) GetByOptions(ctx context.Context, options *model.UserFilterOptions) (*model.User, error) {
// 	queryString, args, err := s.commonSelectQueryBuilder(options).ToSql()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
// 	}

// 	var user model.User
// 	err = s.DBXFromContext(ctx).Raw(queryString, args...).First(&user).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, store.NewErrNotFound(model.UserTableName, "options")
// 		}
// 		return nil, errors.Wrap(err, "failed to find users with given options")
// 	}

// 	return &user, nil
// }

// func (s *SqlUserStore) AddRelations(transaction *gorm.DB, userID string, relations any, customerNoteOnUser bool) *model_helper.AppError {
// 	if transaction == nil {
// 		transaction = s.GetMaster()
// 	}
// 	var association string

// 	switch relations.(type) {
// 	case []*model.Address:
// 		association = "Addresses"
// 	case []*model.CustomerEvent:
// 		association = "CustomerEvents"
// 	case []*model.CustomerNote:
// 		if customerNoteOnUser {
// 			association = "NotesOnMe"
// 		} else {
// 			association = "CustomerNotes"
// 		}
// 	case []*model.StaffNotificationRecipient:
// 		association = "StaffNotificationRecipients"
// 	}

// 	err := transaction.Model(&model.User{Id: userID}).Association(association).Append(relations)
// 	if err != nil {
// 		return model_helper.NewAppError("UserStore.AddRelations", "app.account.add_user_relations.app_error", map[string]interface{}{"relation": "user-" + association}, err.Error(), http.StatusInternalServerError)
// 	}

// 	return nil
// }

// func (s *SqlUserStore) RemoveRelations(transaction *gorm.DB, userID string, relations any, customerNoteOnUser bool) *model_helper.AppError {
// 	if transaction == nil {
// 		transaction = s.GetMaster()
// 	}
// 	var association string

// 	switch relations.(type) {
// 	case []*model.Address:
// 		association = "Addresses"
// 	case []*model.CustomerEvent:
// 		association = "CustomerEvents"
// 	case []*model.CustomerNote:
// 		if customerNoteOnUser {
// 			association = "NotesOnMe"
// 		} else {
// 			association = "CustomerNotes"
// 		}
// 	case []*model.StaffNotificationRecipient:
// 		association = "StaffNotificationRecipients"
// 	}

// 	err := transaction.Model(&model.User{Id: userID}).Association(association).Delete(relations)
// 	if err != nil {
// 		return model_helper.NewAppError("UserStore.AddRelations", "app.account.remove_user_relations.app_error", map[string]interface{}{"relation": "user-" + association}, err.Error(), http.StatusInternalServerError)
// 	}

// 	return nil
// }
