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

type Role struct {
	Id            string
	Name          string
	DisplayName   string
	Description   string
	CreateAt      int64
	UpdateAt      int64
	DeleteAt      int64
	Permissions   string //
	SchemeManaged bool
	BuiltIn       bool
}

func NewRoleFromModel(role *model.Role) *Role {
	permissionsMap := make(map[string]bool)
	var permissionBuilder strings.Builder

	for _, permission := range role.Permissions {
		if !permissionsMap[permission] {
			permissionBuilder.WriteString(" " + permission)
			permissionsMap[permission] = true
		}
	}

	return &Role{
		Id:            role.Id,
		Name:          role.Name,
		DisplayName:   role.DisplayName,
		Description:   role.Description,
		CreateAt:      role.CreateAt,
		UpdateAt:      role.UpdateAt,
		DeleteAt:      role.DeleteAt,
		Permissions:   permissionBuilder.String(),
		SchemeManaged: role.SchemeManaged,
		BuiltIn:       role.BuiltIn,
	}
}

func (role Role) ToModel() *model.Role {
	return &model.Role{
		Id:            role.Id,
		Name:          role.Name,
		DisplayName:   role.DisplayName,
		Description:   role.Description,
		CreateAt:      role.CreateAt,
		UpdateAt:      role.UpdateAt,
		DeleteAt:      role.DeleteAt,
		Permissions:   strings.Fields(role.Permissions),
		SchemeManaged: role.SchemeManaged,
		BuiltIn:       role.BuiltIn,
	}
}

type SqlRoleStore struct {
	store.Store
}

func NewSqlRoleStore(sqlStore store.Store) store.RoleStore {
	return &SqlRoleStore{sqlStore}
}

// Save can be used to both save and update roles
func (s *SqlRoleStore) Save(role *model.Role) (*model.Role, error) {
	// Check the role is valid before proceeding.
	if !role.IsValidWithoutId() {
		return nil, store.NewErrInvalidInput("Role", "<any>", fmt.Sprintf("%v", role))
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

	dbRole := NewRoleFromModel(role)
	dbRole.UpdateAt = model.GetMillis()

	res, err := s.GetMasterX().NamedExec(`UPDATE Roles
		SET UpdateAt=:UpdateAt, DeleteAt=:DeleteAt, CreateAt=:CreateAt,  Name=:Name, DisplayName=:DisplayName,
		Description=:Description, Permissions=:Permissions, SchemeManaged=:SchemeManaged, BuiltIn=:BuiltIn
		 WHERE Id=:Id`, &dbRole)

	if err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	}

	rowsChanged, err := res.RowsAffected()
	if err != nil {
		return nil, errors.Wrap(err, "error while getting rows_affected")
	}

	if rowsChanged != 1 {
		return nil, fmt.Errorf("invalid number of updated rows, expected 1 but got %d", rowsChanged)
	}

	return dbRole.ToModel(), nil
}

func (s *SqlRoleStore) createRole(role *model.Role, transaction store_iface.SqlxTxExecutor) (*model.Role, error) {
	// Check the role is valid before proceeding.
	if !role.IsValidWithoutId() {
		return nil, store.NewErrInvalidInput("Role", "<any>", fmt.Sprintf("%v", role))
	}

	dbRole := NewRoleFromModel(role)

	dbRole.Id = model.NewId()
	dbRole.CreateAt = model.GetMillis()
	dbRole.UpdateAt = dbRole.CreateAt

	if _, err := transaction.NamedExec(`INSERT INTO Roles
		(Id, Name, DisplayName, Description, Permissions, CreateAt, UpdateAt, DeleteAt, SchemeManaged, BuiltIn)
		VALUES
		(:Id, :Name, :DisplayName, :Description, :Permissions, :CreateAt, :UpdateAt, :DeleteAt, :SchemeManaged, :BuiltIn)`, dbRole); err != nil {
		return nil, errors.Wrap(err, "failed to save Role")
	}

	return dbRole.ToModel(), nil
}

func (s *SqlRoleStore) Get(roleId string) (*model.Role, error) {
	var role Role
	if err := s.GetReplicaX().Get(&role, "SELECT * from Roles WHERE Id = ?", roleId); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.RoleTableName, roleId)
		}
		return nil, errors.Wrap(err, "failed to get Role")
	}

	return role.ToModel(), nil
}

func (s *SqlRoleStore) GetAll() ([]*model.Role, error) {
	dbRoles := []Role{}

	if err := s.GetReplicaX().Select(&dbRoles, "SELECT * from Roles"); err != nil {
		return nil, errors.Wrap(err, "failed to find Roles")
	}

	roles := []*model.Role{}
	for _, dbRole := range dbRoles {
		roles = append(roles, dbRole.ToModel())
	}
	return roles, nil
}

func (s *SqlRoleStore) GetByName(ctx context.Context, name string) (*model.Role, error) {
	dbRole := Role{}
	if err := s.DBXFromContext(ctx).Get(&dbRole, "SELECT * from Roles WHERE Name = ?", name); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Role", fmt.Sprintf("name=%s", name))
		}
		return nil, errors.Wrapf(err, "failed to find Roles with name=%s", name)
	}

	return dbRole.ToModel(), nil
}

func (s *SqlRoleStore) GetByNames(names []string) ([]*model.Role, error) {
	if len(names) == 0 {
		return []*model.Role{}, nil
	}

	query := s.GetQueryBuilder().
		Select("Id, Name, DisplayName, Description, CreateAt, UpdateAt, DeleteAt, Permissions, SchemeManaged, BuiltIn").
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
	defer rows.Close()

	roles := []*model.Role{}
	for rows.Next() {
		var role Role
		err = rows.Scan(
			&role.Id, &role.Name, &role.DisplayName, &role.Description,
			&role.CreateAt, &role.UpdateAt, &role.DeleteAt, &role.Permissions,
			&role.SchemeManaged, &role.BuiltIn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan values")
		}
		roles = append(roles, role.ToModel())
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "unable to iterate over rows")
	}

	return roles, nil
}

func (s *SqlRoleStore) Delete(roleId string) (*model.Role, error) {
	// Get the role.
	var role Role
	if err := s.GetReplicaX().Get(&role, "SELECT * from Roles WHERE Id = ?", roleId); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Role", roleId)
		}
		return nil, errors.Wrapf(err, "failed to get Role with id=%s", roleId)
	}

	time := model.GetMillis()
	role.DeleteAt = time
	role.UpdateAt = time

	res, err := s.GetMasterX().NamedExec(`UPDATE Roles
		SET UpdateAt=:UpdateAt, DeleteAt=:DeleteAt, CreateAt=:CreateAt,  Name=:Name, DisplayName=:DisplayName,
		Description=:Description, Permissions=:Permissions, SchemeManaged=:SchemeManaged, BuiltIn=:BuiltIn
		 WHERE Id=:Id`, &role)

	if err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	}

	rowsChanged, err := res.RowsAffected()
	if err != nil {
		return nil, errors.Wrap(err, "error while getting rows_affected")
	}

	if rowsChanged != 1 {
		return nil, fmt.Errorf("invalid number of updated rows, expected 1 but got %d", rowsChanged)
	}

	return role.ToModel(), nil
}

func (s *SqlRoleStore) PermanentDeleteAll() error {
	if _, err := s.GetMasterX().Exec("DELETE FROM Roles"); err != nil {
		return errors.Wrap(err, "failed to delete Roles")
	}

	return nil
}
