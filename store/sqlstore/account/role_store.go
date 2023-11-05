package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

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
		tx := s.GetMaster().Begin()
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "begin_transaction")
		}
		defer s.FinalizeTransaction(tx)

		createdRole, err := s.createRole(role, tx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create Role")
		} else if err := tx.Commit().Error; err != nil {
			return nil, errors.Wrap(err, "commit_transaction")
		}
		return createdRole, nil
	}

	// update

	role.CreateAt = 0 // prevent update

	err := s.GetMaster().Model(&role).Updates(role).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Role")
	}
	role.Permissions = strings.Fields(role.Permmissions_)
	return role, nil
}

func (s *SqlRoleStore) createRole(role *model.Role, transaction *gorm.DB) (*model.Role, error) {
	// Check the role is valid before proceeding.
	if !role.IsValidWithoutId() {
		return nil, store.NewErrInvalidInput("Role", "<any>", fmt.Sprintf("%v", role))
	}

	if err := transaction.Create(role).Error; err != nil {
		return nil, errors.Wrap(err, "failed to save Role")
	}

	return role, nil
}

func (s *SqlRoleStore) Get(roleId string) (*model.Role, error) {
	var role model.Role
	if err := s.GetReplica().First(&role, "Id = ?", roleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.RoleTableName, roleId)
		}
		return nil, errors.Wrap(err, "failed to get Role")
	}

	role.Permissions = strings.Fields(role.Permmissions_)
	return &role, nil
}

func (s *SqlRoleStore) GetAll() ([]*model.Role, error) {
	dbRoles := []*model.Role{}
	if err := s.GetReplica().Find(&dbRoles).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find Roles")
	}

	for _, role := range dbRoles {
		role.Permissions = strings.Fields(role.Permmissions_)
	}
	return dbRoles, nil
}

func (s *SqlRoleStore) GetByName(ctx context.Context, name string) (*model.Role, error) {
	dbRole := model.Role{}
	if err := s.DBXFromContext(ctx).First(&dbRole, "Name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("Role", fmt.Sprintf("name=%s", name))
		}
		return nil, errors.Wrapf(err, "failed to find Roles with name=%s", name)
	}

	dbRole.Permissions = strings.Fields(dbRole.Permmissions_)
	return &dbRole, nil
}

func (s *SqlRoleStore) GetByNames(names []string) ([]*model.Role, error) {
	var roles []*model.Role
	err := s.GetReplica().Find(&roles, "Name IN ?", names).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find roles by names")
	}
	for _, role := range roles {
		role.Permissions = strings.Fields(role.Permmissions_)
	}
	return roles, nil
}

func (s *SqlRoleStore) Delete(roleId string) (*model.Role, error) {
	// Get the role.
	var role model.Role
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

	return &role, nil
}

func (s *SqlRoleStore) PermanentDeleteAll() error {
	if err := s.GetMaster().Raw("DELETE FROM Roles").Error; err != nil {
		return errors.Wrap(err, "failed to delete Roles")
	}

	return nil
}
