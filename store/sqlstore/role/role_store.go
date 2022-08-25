package role

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
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
	return &SqlRoleStore{sqlStore}
}

// Save can be used to both save and update roles
func (s *SqlRoleStore) Save(role *model.Role) (*model.Role, error) {
	// Check the role is valid before proceeding.
	if ok, field, errValue := role.IsValidWithoutId(); !ok {
		return nil, store.NewErrInvalidInput("Role", field, fmt.Sprintf("%v", errValue))
	}

	if role.Id == "" { // this means create new Role
		transaction, err := s.GetMasterX().Beginx()
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
	if rowsChanged, err := s.GetMasterX().Update(role); err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	} else if rowsChanged != 1 {
		return nil, fmt.Errorf("invalid number of updated rows, expected 1 but got %d", rowsChanged)
	}

	role.PopulatePermissionSlice()
	return role, nil
}

func (s *SqlRoleStore) createRole(role *model.Role, transaction store_iface.SqlxTxExecutor) (*model.Role, error) {
	role.PreSave()
	if err := transaction.Insert(role); err != nil {
		return nil, errors.Wrap(err, "failed to save Role")
	}

	role.PopulatePermissionSlice()
	return role, nil
}

func (s *SqlRoleStore) Get(roleId string) (*model.Role, error) {
	var role model.Role
	if err := s.GetReplicaX().Get(&role, "SELECT * from Roles WHERE Id = ?", roleId); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.RoleTableName, roleId)
		}
		return nil, errors.Wrap(err, "failed to get Role")
	}

	role.PopulatePermissionSlice()
	return &role, nil
}

func (s *SqlRoleStore) GetAll() ([]*model.Role, error) {
	var dbRoles []*model.Role

	if err := s.GetReplicaX().Select(&dbRoles, "SELECT * from Roles"); err != nil {
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
	if err := s.DBXFromContext(ctx).Get(&rolw, "SELECT * from Roles WHERE Name = ?", name); err != nil {
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

	rows, err := s.GetReplicaX().QueryX(queryString, args...)
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
	if err := s.GetReplicaX().Get(&role, "SELECT * from Roles WHERE Id = ?", roleId); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Role", roleId)
		}
		return nil, errors.Wrapf(err, "failed to get Role with id=%s", roleId)
	}

	role.PreUpdate()
	role.DeleteAt = role.UpdateAt

	if rowsChanged, err := s.GetMasterX().Update(role); err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	} else if rowsChanged != 1 {
		return nil, errors.Wrapf(err, "invalid number of updated rows, expected 1 but got %d", rowsChanged)
	}

	role.PopulatePermissionSlice()

	return &role, nil
}

func (s *SqlRoleStore) PermanentDeleteAll() error {
	if _, err := s.GetMasterX().Exec("DELETE FROM Roles"); err != nil {
		return errors.Wrap(err, "failed to delete Roles")
	}

	return nil
}
