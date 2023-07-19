package menu

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlMenuItemStore struct {
	store.Store
}

func NewSqlMenuItemStore(s store.Store) store.MenuItemStore {
	return &SqlMenuItemStore{s}
}

func (s *SqlMenuItemStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

func (is *SqlMenuItemStore) Save(item *model.MenuItem) (*model.MenuItem, error) {
	item.PreSave()
	if err := item.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.MenuItemTableName + "(" + is.ModelFields("").Join(",") + ") VALUES (" + is.ModelFields(":").Join(",") + ")"
	if _, err := is.GetMaster().NamedExec(query, item); err != nil {
		return nil, errors.Wrapf(err, "failed to save menu item with id=%s", item.Id)
	}
	return item, nil
}

func (is *SqlMenuItemStore) commonQueryBuilder(options *model.MenuItemFilterOptions) squirrel.SelectBuilder {
	query := is.GetQueryBuilder().
		Select(is.ModelFields(model.MenuItemTableName + ".")...).
		From(model.MenuItemTableName)

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

	return query
}

func (is *SqlMenuItemStore) GetByOptions(options *model.MenuItemFilterOptions) (*model.MenuItem, error) {
	queryString, args, err := is.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOptions_ToSql")
	}

	var menuItem model.MenuItem

	err = is.GetReplica().Get(&menuItem, queryString, args...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.MenuItemTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find menu item with given options")
	}

	return &menuItem, nil
}

func (is *SqlMenuItemStore) FilterByOptions(options *model.MenuItemFilterOptions) ([]*model.MenuItem, error) {
	queryString, args, err := is.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.MenuItem
	err = is.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find menu items by given options")
	}

	return res, nil
}
