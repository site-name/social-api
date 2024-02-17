package menu

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlMenuStore struct {
	store.Store
}

func NewSqlMenuStore(sqlStore store.Store) store.MenuStore {
	return &SqlMenuStore{sqlStore}
}

func (ms *SqlMenuStore) Upsert(menu model.Menu) (*model.Menu, error) {
	isSaving := menu.ID == ""
	if isSaving {
		model_helper.MenuPreSave(&menu)
	} else {
		model_helper.MenuCommonPre(&menu)
	}

	if err := model_helper.MenuIsValid(menu); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = menu.Insert(ms.GetMaster(), boil.Infer())
	} else {
		_, err = menu.Update(ms.GetMaster(), boil.Blacklist(model.MenuColumns.CreatedAt))
	}
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func (ms *SqlMenuStore) GetByOptions(options model_helper.MenuFilterOptions) (*model.Menu, error) {
	menu, err := model.Menus(options.Conditions...).One(ms.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Menus, "options")
		}
		return nil, err
	}
	return menu, nil
}

func (ms *SqlMenuStore) FilterByOptions(options model_helper.MenuFilterOptions) (model.MenuSlice, error) {
	return model.Menus(options.Conditions...).All(ms.GetReplica())
}

func (s *SqlMenuStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = s.GetMaster()
	}
	return model.Menus(model.MenuWhere.ID.IN(ids)).DeleteAll(tx)
}
