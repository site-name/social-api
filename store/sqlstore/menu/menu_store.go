package menu

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

type SqlMenuStore struct {
	store.Store
}

func NewSqlMenuStore(sqlStore store.Store) store.MenuStore {
	return &SqlMenuStore{sqlStore}
}

func (s *SqlMenuStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"Name",
		"Slug",
		"CreateAt",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ms *SqlMenuStore) Save(mnu *menu.Menu) (*menu.Menu, error) {
	mnu.PreSave()
	if err := mnu.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.MenuTableName + "(" + ms.ModelFields("").Join(",") + ") VALUES (" + ms.ModelFields(":").Join(",") + ")"

	if _, err := ms.GetMasterX().NamedExec(query, mnu); err != nil {
		if ms.IsUniqueConstraintError(err, []string{"Name", "menus_name_key", "idx_menus_name_unique"}) {
			return nil, store.NewErrInvalidInput(store.MenuTableName, "Name", mnu.Name)
		}
		if ms.IsUniqueConstraintError(err, []string{"Slug", "menus_slug_key", "idx_menus_slug_unique"}) {
			return nil, store.NewErrInvalidInput(store.MenuTableName, "Slug", mnu.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save menu with id=%s", mnu.Id)
	}

	return mnu, nil
}

func (ms *SqlMenuStore) GetByOptions(options *menu.MenuFilterOptions) (*menu.Menu, error) {
	query := ms.GetQueryBuilder().
		Select("*").
		From(store.MenuTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Name != nil {
		query = query.Where(options.Name)
	}
	if options.Slug != nil {
		query = query.Where(options.Slug)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOptions_ToSql")
	}

	var res menu.Menu
	err = ms.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.MenuTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find menu with given options")
	}

	return &res, nil
}
