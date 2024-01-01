package account

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type SqlRoleStore struct {
	store.Store
}

func NewSqlRoleStore(sqlStore store.Store) store.RoleStore {
	return &SqlRoleStore{sqlStore}
}

// Save can be used to both save and update roles
func (s *SqlRoleStore) Upsert(role model.Role) (*model.Role, error) {
	if !model_helper.RoleIsValidWithoutId(role) {
		return nil, store.NewErrInvalidInput(model.TableNames.Roles, "any", nil)
	}

	if role.ID == "" {
		tx, err := s.GetMaster().BeginTx(s.Context(), &sql.TxOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "begin_transaction")
		}
		defer s.FinalizeTransaction(tx)

		model_helper.RolePreSave(&role)

		err = role.Insert(tx, boil.Infer())
		if err != nil {
			return nil, err
		} else if err := tx.Commit(); err != nil {
			return nil, errors.Wrap(err, "commit_transaction")
		}

		return &role, nil
	}

	_, err := role.Update(s.GetMaster(), boil.Blacklist(model.RoleColumns.CreatedAt))
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *SqlRoleStore) Get(roleId string) (*model.Role, error) {
	role, err := model.FindRole(s.GetReplica(), roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Roles, roleId)
		}
		return nil, err
	}
	return role, nil
}

func (s *SqlRoleStore) GetAll() (model.RoleSlice, error) {
	return model.Roles().All(s.GetReplica())
}

func (s *SqlRoleStore) GetByName(ctx context.Context, name string) (*model.Role, error) {
	role, err := model.Roles(model.RoleWhere.Name.EQ(name)).One(s.DBXFromContext(ctx))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Roles, name)
		}
		return nil, err
	}

	return role, nil
}

func (s *SqlRoleStore) GetByNames(names []string) (model.RoleSlice, error) {
	return model.Roles(model.RoleWhere.Name.IN(names)).All(s.GetReplica())
}

func (s *SqlRoleStore) Delete(roleId string) (*model.Role, error) {
	_, err := model.
		Roles(model.RoleWhere.ID.EQ(roleId)).
		UpdateAll(s.GetMaster(), model.M{
			model.RoleColumns.DeleteAt: model_helper.GetMillis(),
		})
	if err != nil {
		return nil, err
	}

	return &model.Role{
		ID: roleId,
	}, nil
}

func (s *SqlRoleStore) PermanentDeleteAll() error {
	_, err := model.Roles().DeleteAll(s.GetMaster())
	return err
}
