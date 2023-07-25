package menu

import (
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
	if err := ms.GetMaster().Create(menu).Error; err != nil {
		if ms.IsUniqueConstraintError(err, []string{"Name", "menus_name_key", "idx_menus_name_unique"}) {
			return nil, store.NewErrInvalidInput(model.MenuTableName, "Name", menu.Name)
		}
		if ms.IsUniqueConstraintError(err, []string{"Slug", "menus_slug_key", "idx_menus_slug_unique"}) {
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
