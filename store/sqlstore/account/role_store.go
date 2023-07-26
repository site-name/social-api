package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type Role struct {
	Id            string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name          string `gorm:"type:varchar(64);column:Name;unique"`
	DisplayName   string `gorm:"type:varchar(128);column:DisplayName"`
	Description   string `gorm:"type:varchar(1024);column:Description"`
	CreateAt      int64  `gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	UpdateAt      int64  `gorm:"type:bigint;autoCreateTime:milli;autoUpdateTime:milli;column:UpdateAt"`
	DeleteAt      int64  `gorm:"type:bigint;column:DeleteAt"`
	Permissions   string `gorm:"type:varchar(5000);column:Permissions"`
	SchemeManaged bool   `gorm:"column:SchemeManaged"`
	BuiltIn       bool   `gorm:"column:Builtin"`
}

func (c *Role) BeforeCreate(_ *gorm.DB) error { return nil }
func (c *Role) BeforeUpdate(_ *gorm.DB) error { return nil }
func (c *Role) TableName() string             { return model.RoleTableName }

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
		transaction := s.GetMaster().Begin()
		defer transaction.Rollback()

		createdRole, err := s.createRole(role, transaction)
		if err != nil {
			_ = transaction.Rollback()
			return nil, errors.Wrap(err, "unable to create Role")
		} else if err := transaction.Commit().Error; err != nil {
			return nil, errors.Wrap(err, "commit_transaction")
		}
		return createdRole, nil
	}

	// update

	dbRole := NewRoleFromModel(role)
	dbRole.CreateAt = 0 // prevent update

	err := s.GetMaster().Model(&dbRole).Updates(dbRole).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	}

	return dbRole.ToModel(), nil
}

func (s *SqlRoleStore) createRole(role *model.Role, transaction *gorm.DB) (*model.Role, error) {
	// Check the role is valid before proceeding.
	if !role.IsValidWithoutId() {
		return nil, store.NewErrInvalidInput("Role", "<any>", fmt.Sprintf("%v", role))
	}

	dbRole := NewRoleFromModel(role)
	if err := transaction.Create(dbRole).Error; err != nil {
		return nil, errors.Wrap(err, "failed to save Role")
	}

	return dbRole.ToModel(), nil
}

func (s *SqlRoleStore) Get(roleId string) (*model.Role, error) {
	var role Role
	if err := s.GetReplica().First(&role, "Id = ?", roleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.RoleTableName, roleId)
		}
		return nil, errors.Wrap(err, "failed to get Role")
	}

	return role.ToModel(), nil
}

func (s *SqlRoleStore) GetAll() ([]*model.Role, error) {
	dbRoles := []*Role{}
	if err := s.GetReplica().Find(&dbRoles).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find Roles")
	}
	return lo.Map(dbRoles, func(item *Role, _ int) *model.Role { return item.ToModel() }), nil
}

func (s *SqlRoleStore) GetByName(ctx context.Context, name string) (*model.Role, error) {
	dbRole := Role{}
	if err := s.DBXFromContext(ctx).First(&dbRole, "Name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("Role", fmt.Sprintf("name=%s", name))
		}
		return nil, errors.Wrapf(err, "failed to find Roles with name=%s", name)
	}

	return dbRole.ToModel(), nil
}

func (s *SqlRoleStore) GetByNames(names []string) ([]*model.Role, error) {
	var roles []*Role
	err := s.GetReplica().Find(&roles, "Name IN ?", names).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find roles by names")
	}
	return lo.Map(roles, func(item *Role, _ int) *model.Role { return item.ToModel() }), nil
}

func (s *SqlRoleStore) Delete(roleId string) (*model.Role, error) {
	// Get the role.
	var role Role
	if err := s.GetReplica().First(&role, "Id = ?", roleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("Role", roleId)
		}
		return nil, errors.Wrapf(err, "failed to get Role with id=%s", roleId)
	}

	time := model.GetMillis()
	role.CreateAt = 0 // prevent update
	role.DeleteAt = time

	err := s.GetMaster().Model(role).Updates(role).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	}

	return role.ToModel(), nil
}

func (s *SqlRoleStore) PermanentDeleteAll() error {
	if err := s.GetMaster().Raw("DELETE FROM Roles").Error; err != nil {
		return errors.Wrap(err, "failed to delete Roles")
	}

	return nil
}
