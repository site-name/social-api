package menu

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlMenuStore struct {
	store.Store
}

func NewSqlMenuStore(sqlStore store.Store) store.MenuStore {
	return &SqlMenuStore{sqlStore}
}

func (ms *SqlMenuStore) Save(menu *model.Menu) (*model.Menu, error) {
	if err := ms.GetMaster().Save(menu).Error; err != nil {
		if ms.IsUniqueConstraintError(err, []string{"Name", "name_unique_key"}) {
			return nil, store.NewErrInvalidInput(model.MenuTableName, "Name", menu.Name)
		}
		if ms.IsUniqueConstraintError(err, []string{"Slug", "slug_unique_key"}) {
			return nil, store.NewErrInvalidInput(model.MenuTableName, "Slug", menu.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save menu with id=%s", menu.Id)
	}

	return menu, nil
}

func (ms *SqlMenuStore) GetByOptions(options *model.MenuFilterOptions) (*model.Menu, error) {
	args, err := store.BuildSqlizer(options.Conditions, "Menu_GetByOptions")
	if err != nil {
		return nil, err
	}

	var res model.Menu
	err = ms.GetReplica().First(&res, args...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.MenuTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find menu with given options")
	}

	return &res, nil
}

func (ms *SqlMenuStore) FilterByOptions(options *model.MenuFilterOptions) ([]*model.Menu, error) {
	args, err := store.BuildSqlizer(options.Conditions, "Menu_FilterByOptions")
	if err != nil {
		return nil, err
	}

	var res []*model.Menu
	err = ms.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find menus with given options")
	}

	return res, nil
}

func (s *SqlMenuStore) Delete(ids []string) (int64, *model_helper.AppError) {
	result := s.GetMaster().Raw("DELETE FROM "+model.MenuTableName+" WHERE Id IN ?", ids)
	if result.Error != nil {
		return 0, model_helper.NewAppError("DeleteMenu", "app.menu.delete_menus.app_error", nil, result.Error.Error(), http.StatusInternalServerError)
	}
	return result.RowsAffected, nil
}
