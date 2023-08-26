package menu

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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
	var res model.Menu
	err := ms.GetReplica().First(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.MenuTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find menu with given options")
	}

	return &res, nil
}

func (ms *SqlMenuStore) FilterByOptions(options *model.MenuFilterOptions) ([]*model.Menu, error) {
	var res []*model.Menu
	err := ms.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find menus with given options")
	}

	return res, nil
}

func (s *SqlMenuStore) Delete(ids []string) (int64, *model.AppError) {
	result := s.GetMaster().Raw("DELETE FROM "+model.MenuTableName+" WHERE Id IN ?", ids)
	if result.Error != nil {
		return 0, model.NewAppError("DeleteMenu", "app.menu.delete_menus.app_error", nil, result.Error.Error(), http.StatusInternalServerError)
	}
	return result.RowsAffected, nil
}
