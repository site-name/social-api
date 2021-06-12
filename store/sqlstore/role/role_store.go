package role

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlRoleStore struct {
	store.Store
}

type channelRolesPermissions struct {
	GuestRoleName                string
	UserRoleName                 string
	AdminRoleName                string
	HigherScopedGuestPermissions string
	HigherScopedUserPermissions  string
	HigherScopedAdminPermissions string
}

func NewSqlRoleStore(sqlStore store.Store) store.RoleStore {
	s := &SqlRoleStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Role{}, "Roles").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(64).SetUnique(true)
		table.ColMap("DisplayName").SetMaxSize(128)
		table.ColMap("Description").SetMaxSize(1024)
		table.ColMap("PermissionsStr") // NOTE: don's save Permissions slice
	}
	return s
}

func (s *SqlRoleStore) CreateIndexesIfNotExists() {}

// Save can be used to both save and update roles
func (s *SqlRoleStore) Save(role *model.Role) (*model.Role, error) {
	// Check the role is valid before proceeding.
	if ok, field, errValue := role.IsValidWithoutId(); !ok {
		return nil, store.NewErrInvalidInput("Role", field, fmt.Sprintf("%v", errValue))
	}

	if role.Id == "" { // this means create new Role
		transaction, err := s.GetMaster().Begin()
		if err != nil {
			return nil, errors.Wrap(err, "begin_transaction")
		}
		defer store.FinalizeTransaction(transaction)
		createdRole, err := s.createRole(role, transaction)
		if err != nil {
			_ = transaction.Rollback()
			return nil, errors.Wrap(err, "unable to create Role")
		} else if err := transaction.Commit(); err != nil {
			return nil, errors.Wrap(err, "commit_transaction")
		}
		return createdRole, nil
	}

	role.PreUpdate()
	if rowsChanged, err := s.GetMaster().Update(role); err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	} else if rowsChanged != 1 {
		return nil, fmt.Errorf("invalid number of updated rows, expected 1 but got %d", rowsChanged)
	}

	role.PopulatePermissionSlice()
	return role, nil
}

func (s *SqlRoleStore) createRole(role *model.Role, transaction *gorp.Transaction) (*model.Role, error) {
	role.PreSave()
	if err := transaction.Insert(role); err != nil {
		return nil, errors.Wrap(err, "failed to save Role")
	}

	role.PopulatePermissionSlice()
	return role, nil
}

func (s *SqlRoleStore) Get(roleId string) (*model.Role, error) {
	var role model.Role
	if err := s.GetReplica().SelectOne(&role, "SELECT * from Roles WHERE Id = :Id", map[string]interface{}{"Id": roleId}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Role", roleId)
		}
		return nil, errors.Wrap(err, "failed to get Role")
	}

	role.PopulatePermissionSlice()
	return &role, nil
}

func (s *SqlRoleStore) GetAll() ([]*model.Role, error) {
	var dbRoles []*model.Role

	if _, err := s.GetReplica().Select(&dbRoles, "SELECT * from Roles", map[string]interface{}{}); err != nil {
		return nil, errors.Wrap(err, "failed to find Roles")
	}

	for _, role := range dbRoles {
		role.Permissions = strings.Fields(role.PermissionsStr)
		role.PermissionsStr = ""
	}
	return dbRoles, nil
}

func (s *SqlRoleStore) GetByName(ctx context.Context, name string) (*model.Role, error) {
	var rolw model.Role
	if err := s.DBFromContext(ctx).SelectOne(&rolw, "SELECT * from Roles WHERE Name = :Name", map[string]interface{}{"Name": name}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Role", fmt.Sprintf("name=%s", name))
		}
		return nil, errors.Wrapf(err, "failed to find Roles with name=%s", name)
	}

	rolw.PopulatePermissionSlice()
	return &rolw, nil
}

func (s *SqlRoleStore) GetByNames(names []string) ([]*model.Role, error) {
	if len(names) == 0 {
		return []*model.Role{}, nil
	}

	query := s.GetQueryBuilder().
		Select("Id, Name, DisplayName, Description, CreateAt, UpdateAt, DeleteAt, PermissionsStr, SchemeManaged, BuiltIn").
		From("Roles").
		Where(sq.Eq{"Name": names})
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "role_tosql")
	}

	rows, err := s.GetReplica().Db.Query(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find Roles")
	}

	var roles []*model.Role
	defer rows.Close()
	for rows.Next() {
		var role model.Role
		err = rows.Scan(
			&role.Id, &role.Name, &role.DisplayName, &role.Description,
			&role.CreateAt, &role.UpdateAt, &role.DeleteAt, &role.PermissionsStr,
			&role.SchemeManaged, &role.BuiltIn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan values")
		}
		role.PopulatePermissionSlice()
		roles = append(roles, &role)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "unable to iterate over rows")
	}

	return roles, nil
}

func (s *SqlRoleStore) Delete(roleId string) (*model.Role, error) {
	// Get the role.
	var role model.Role
	if err := s.GetReplica().SelectOne(&role, "SELECT * from Roles WHERE Id = :Id", map[string]interface{}{"Id": roleId}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Role", roleId)
		}
		return nil, errors.Wrapf(err, "failed to get Role with id=%s", roleId)
	}

	role.PreUpdate()
	role.DeleteAt = role.UpdateAt

	if rowsChanged, err := s.GetMaster().Update(role); err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	} else if rowsChanged != 1 {
		return nil, errors.Wrapf(err, "invalid number of updated rows, expected 1 but got %d", rowsChanged)
	}

	role.PopulatePermissionSlice()
	return &role, nil
}

func (s *SqlRoleStore) PermanentDeleteAll() error {
	if _, err := s.GetMaster().Exec("DELETE FROM Roles"); err != nil {
		return errors.Wrap(err, "failed to delete Roles")
	}

	return nil
}

func (s *SqlRoleStore) channelHigherScopedPermissionsQuery(roleNames []string) string {
	sqlTmpl := `
		SELECT
			'' AS GuestRoleName,
			RoleSchemes.DefaultChannelUserRole AS UserRoleName,
			RoleSchemes.DefaultChannelAdminRole AS AdminRoleName,
			'' AS HigherScopedGuestPermissions,
			UserRoles.Permissions AS HigherScopedUserPermissions,
			AdminRoles.Permissions AS HigherScopedAdminPermissions
		FROM
			Schemes AS RoleSchemes
			JOIN Channels ON Channels.SchemeId = RoleSchemes.Id
			JOIN Teams ON Teams.Id = Channels.TeamId
			JOIN Schemes ON Schemes.Id = Teams.SchemeId
			RIGHT JOIN Roles AS UserRoles ON UserRoles.Name = Schemes.DefaultChannelUserRole
			RIGHT JOIN Roles AS AdminRoles ON AdminRoles.Name = Schemes.DefaultChannelAdminRole
		WHERE
			RoleSchemes.DefaultChannelUserRole IN ('%[1]s')
			OR RoleSchemes.DefaultChannelAdminRole IN ('%[1]s')

		UNION

		SELECT
			RoleSchemes.DefaultChannelGuestRole AS GuestRoleName,
			'' AS UserRoleName,
			'' AS AdminRoleName,
			GuestRoles.Permissions AS HigherScopedGuestPermissions,
			'' AS HigherScopedUserPermissions,
			'' AS HigherScopedAdminPermissions
		FROM
			Schemes AS RoleSchemes
			JOIN Channels ON Channels.SchemeId = RoleSchemes.Id
			JOIN Teams ON Teams.Id = Channels.TeamId
			JOIN Schemes ON Schemes.Id = Teams.SchemeId
			RIGHT JOIN Roles AS GuestRoles ON GuestRoles.Name = Schemes.DefaultChannelGuestRole
		WHERE
			RoleSchemes.DefaultChannelGuestRole IN ('%[1]s')

		UNION

		SELECT
			Schemes.DefaultChannelGuestRole AS GuestRoleName,
			Schemes.DefaultChannelUserRole AS UserRoleName,
			Schemes.DefaultChannelAdminRole AS AdminRoleName,
			GuestRoles.Permissions AS HigherScopedGuestPermissions,
			UserRoles.Permissions AS HigherScopedUserPermissions,
			AdminRoles.Permissions AS HigherScopedAdminPermissions
		FROM
			Schemes
			JOIN Channels ON Channels.SchemeId = Schemes.Id
			JOIN Teams ON Teams.Id = Channels.TeamId
			JOIN Roles AS GuestRoles ON GuestRoles.Name = '%[2]s'
			JOIN Roles AS UserRoles ON UserRoles.Name = '%[3]s'
			JOIN Roles AS AdminRoles ON AdminRoles.Name = '%[4]s'
		WHERE
			(Schemes.DefaultChannelGuestRole IN ('%[1]s')
			OR Schemes.DefaultChannelUserRole IN ('%[1]s')
			OR Schemes.DefaultChannelAdminRole IN ('%[1]s'))
		AND (Teams.SchemeId = ''
			OR Teams.SchemeId IS NULL)
	`

	// The below three channel role names are referenced by their name value because there is no system scheme
	// record that ships with Mattermost, otherwise the system scheme would be referenced by name and the channel
	// roles would be referenced by their column names.
	return fmt.Sprintf(
		sqlTmpl,
		strings.Join(roleNames, "', '"),
		model.CHANNEL_GUEST_ROLE_ID,
		model.CHANNEL_USER_ROLE_ID,
		model.CHANNEL_ADMIN_ROLE_ID,
	)
}

func (s *SqlRoleStore) ChannelHigherScopedPermissions(roleNames []string) (map[string]*model.RolePermissions, error) {
	query := s.channelHigherScopedPermissionsQuery(roleNames)

	var rolesPermissions []*channelRolesPermissions
	if _, err := s.GetReplica().Select(&rolesPermissions, query); err != nil {
		return nil, errors.Wrap(err, "failed to find RolePermissions")
	}

	roleNameHigherScopedPermissions := map[string]*model.RolePermissions{}

	for _, rp := range rolesPermissions {
		roleNameHigherScopedPermissions[rp.GuestRoleName] = &model.RolePermissions{RoleID: model.CHANNEL_GUEST_ROLE_ID, Permissions: strings.Split(rp.HigherScopedGuestPermissions, " ")}
		roleNameHigherScopedPermissions[rp.UserRoleName] = &model.RolePermissions{RoleID: model.CHANNEL_USER_ROLE_ID, Permissions: strings.Split(rp.HigherScopedUserPermissions, " ")}
		roleNameHigherScopedPermissions[rp.AdminRoleName] = &model.RolePermissions{RoleID: model.CHANNEL_ADMIN_ROLE_ID, Permissions: strings.Split(rp.HigherScopedAdminPermissions, " ")}
	}

	return roleNameHigherScopedPermissions, nil
}

// func (s *SqlRoleStore) AllChannelSchemeRoles() ([]*model.Role, error) {
// 	query := s.getQueryBuilder().
// 		Select("Roles.*").
// 		From("Schemes").
// 		Join("Roles ON Schemes.DefaultChannelGuestRole = Roles.Name OR Schemes.DefaultChannelUserRole = Roles.Name OR Schemes.DefaultChannelAdminRole = Roles.Name").
// 		Where(sq.Eq{"Schemes.Scope": model.SCHEME_SCOPE_CHANNEL}).
// 		Where(sq.Eq{"Roles.DeleteAt": 0}).
// 		Where(sq.Eq{"Schemes.DeleteAt": 0})

// 	queryString, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "role_tosql")
// 	}

// 	var dbRoles []*Role
// 	if _, err = s.GetReplica().Select(&dbRoles, queryString, args...); err != nil {
// 		return nil, errors.Wrap(err, "failed to find Roles")
// 	}

// 	var roles []*model.Role
// 	for _, dbRole := range dbRoles {
// 		roles = append(roles, dbRole.ToModel())
// 	}

// 	return roles, nil
// }

// ChannelRolesUnderTeamRole finds all of the channel-scheme roles under the team of the given team-scheme role.
// func (s *SqlRoleStore) ChannelRolesUnderTeamRole(roleName string) ([]*model.Role, error) {
// 	query := s.getQueryBuilder().
// 		Select("ChannelSchemeRoles.*").
// 		From("Roles AS HigherScopedRoles").
// 		Join("Schemes AS HigherScopedSchemes ON (HigherScopedRoles.Name = HigherScopedSchemes.DefaultChannelGuestRole OR HigherScopedRoles.Name = HigherScopedSchemes.DefaultChannelUserRole OR HigherScopedRoles.Name = HigherScopedSchemes.DefaultChannelAdminRole)").
// 		Join("Teams ON Teams.SchemeId = HigherScopedSchemes.Id").
// 		Join("Channels ON Channels.TeamId = Teams.Id").
// 		Join("Schemes AS ChannelSchemes ON Channels.SchemeId = ChannelSchemes.Id").
// 		Join("Roles AS ChannelSchemeRoles ON (ChannelSchemeRoles.Name = ChannelSchemes.DefaultChannelGuestRole OR ChannelSchemeRoles.Name = ChannelSchemes.DefaultChannelUserRole OR ChannelSchemeRoles.Name = ChannelSchemes.DefaultChannelAdminRole)").
// 		Where(sq.Eq{"HigherScopedSchemes.Scope": model.SCHEME_SCOPE_TEAM}).
// 		Where(sq.Eq{"HigherScopedRoles.Name": roleName}).
// 		Where(sq.Eq{"HigherScopedRoles.DeleteAt": 0}).
// 		Where(sq.Eq{"HigherScopedSchemes.DeleteAt": 0}).
// 		Where(sq.Eq{"Teams.DeleteAt": 0}).
// 		Where(sq.Eq{"Channels.DeleteAt": 0}).
// 		Where(sq.Eq{"ChannelSchemes.DeleteAt": 0}).
// 		Where(sq.Eq{"ChannelSchemeRoles.DeleteAt": 0})

// 	queryString, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "role_tosql")
// 	}

// 	var dbRoles []*Role
// 	if _, err = s.GetReplica().Select(&dbRoles, queryString, args...); err != nil {
// 		return nil, errors.Wrap(err, "failed to find Roles")
// 	}

// 	var roles []*model.Role
// 	for _, dbRole := range dbRoles {
// 		roles = append(roles, dbRole.ToModel())
// 	}

// 	return roles, nil
// }
