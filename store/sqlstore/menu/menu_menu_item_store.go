package menu

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlMenuItemStore struct {
	store.Store
}

func NewSqlMenuItemStore(s store.Store) store.MenuItemStore {
	return &SqlMenuItemStore{s}
}

func (is *SqlMenuItemStore) Save(item *model.MenuItem) (*model.MenuItem, error) {
	if err := is.GetMaster().Save(item).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to save menu item with id=%s", item.Id)
	}
	return item, nil
}

func (is *SqlMenuItemStore) GetByOptions(options *model.MenuItemFilterOptions) (*model.MenuItem, error) {
	args, err := store.BuildSqlizer(options.Conditions, "MenuItem_GetByOptions")
	if err != nil {
		return nil, err
	}

	var menuItem model.MenuItem
	err = is.GetReplica().First(&menuItem, args...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.MenuItemTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find menu item with given options")
	}

	return &menuItem, nil
}

func (is *SqlMenuItemStore) FilterByOptions(options *model.MenuItemFilterOptions) ([]*model.MenuItem, error) {
	args, err := store.BuildSqlizer(options.Conditions, "MenuItem_FilterByOptions")
	if err != nil {
		return nil, err
	}
	var res []*model.MenuItem
	err = is.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find menu items by given options")
	}

	return res, nil
}

func (s *SqlMenuItemStore) Delete(ids []string) (int64, *model_helper.AppError) {
	result := s.GetMaster().Raw("DELETE FROM "+model.MenuItemTableName+" WHERE Id IN ?", ids)
	if result.Error != nil {
		return 0, model_helper.NewAppError("DeleteMenu", "app.menu.delete_menu_items.app_error", nil, result.Error.Error(), http.StatusInternalServerError)
	}
	return result.RowsAffected, nil
}
