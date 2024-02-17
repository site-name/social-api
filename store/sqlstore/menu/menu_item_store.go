package menu

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlMenuItemStore struct {
	store.Store
}

func NewSqlMenuItemStore(s store.Store) store.MenuItemStore {
	return &SqlMenuItemStore{s}
}

func (mis *SqlMenuItemStore) Upsert(item model.MenuItem) (*model.MenuItem, error) {
	isSaving := item.ID == ""
	model_helper.MenuItemCommonPre(&item)

	if err := model_helper.MenuItemIsValid(item); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = item.Insert(mis.GetMaster(), boil.Infer())
	} else {
		_, err = item.Update(mis.GetMaster(), boil.Infer())
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (mis *SqlMenuItemStore) GetByOptions(options model_helper.MenuItemFilterOptions) (*model.MenuItem, error) {
	item, err := model.MenuItems(options.Conditions...).One(mis.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.MenuItems, "options")
		}
		return nil, err
	}
	return item, nil
}

func (mis *SqlMenuItemStore) FilterByOptions(options model_helper.MenuItemFilterOptions) (model.MenuItemSlice, error) {
	return model.MenuItems(options.Conditions...).All(mis.GetReplica())
}

func (s *SqlMenuItemStore) Delete(ids []string) (int64, error) {
	return model.MenuItems(model.MenuItemWhere.ID.IN(ids)).DeleteAll(s.GetMaster())
}
