package menu

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

type SqlMenuItemStore struct {
	store.Store
}

func NewSqlMenuItemStore(s store.Store) store.MenuItemStore {
	return &SqlMenuItemStore{s}
}

func (s *SqlMenuItemStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"MenuID",
		"Name",
		"ParentID",
		"Url",
		"CategoryID",
		"CollectionID",
		"PageID",
		"Metadata",
		"PrivateMetadata",
		"SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (is *SqlMenuItemStore) Save(item *menu.MenuItem) (*menu.MenuItem, error) {
	item.PreSave()
	if err := item.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.MenuItemTableName + "(" + is.ModelFields("").Join(",") + ") VALUES (" + is.ModelFields(":").Join(",") + ")"
	if _, err := is.GetMasterX().NamedExec(query, item); err != nil {
		return nil, errors.Wrapf(err, "failed to save menu item with id=%s", item.Id)
	}
	return item, nil
}

func (is *SqlMenuItemStore) GetByOptions(options *menu.MenuItemFilterOptions) (*menu.MenuItem, error) {
	query := is.GetQueryBuilder().
		Select("*").
		From(store.MenuItemTableName)

	// parse options
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Name != nil {
		query = query.Where(options.Name)
	}
	if options.MenuID != nil {
		query = query.Where(options.MenuID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOptions_ToSql")
	}

	var menuItem menu.MenuItem

	err = is.GetReplicaX().Get(&menuItem, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.MenuItemTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find menu item with given options")
	}

	return &menuItem, nil
}
