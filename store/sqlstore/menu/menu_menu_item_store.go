package menu

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/page"
	"github.com/sitename/sitename/store/sqlstore/product"
)

const (
	MenuItemTableName = "MenuItems"
)

type SqlMenuItemStore struct {
	store.Store
}

func NewSqlMenuItemStore(s store.Store) store.MenuItemStore {
	is := &SqlMenuItemStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(menu.MenuItem{}, MenuItemTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("MenuID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ParentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PageID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(menu.MENU_ITEM_NAME_MAX_LENGTH)
		table.ColMap("Url").SetMaxSize(menu.MENU_ITEM_URL_MAX_LENGTH)
	}
	return is
}

func (is *SqlMenuItemStore) CreateIndexesIfNotExists() {
	is.CreateForeignKeyIfNotExists(MenuItemTableName, "MenuID", MenuTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(MenuItemTableName, "ParentID", MenuItemTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(MenuItemTableName, "CategoryID", store.ProductCategoryTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(MenuItemTableName, "CollectionID", product.CollectionTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(MenuItemTableName, "PageID", page.PageTableName, "Id", true)
}

func (is *SqlMenuItemStore) Save(item *menu.MenuItem) (*menu.MenuItem, error) {
	item.PreSave()
	if err := item.IsValid(); err != nil {
		return nil, err
	}

	if err := is.GetMaster().Insert(item); err != nil {
		return nil, errors.Wrapf(err, "failed to save menu item with id=%s", item.Id)
	}
	return item, nil
}

func (is *SqlMenuItemStore) GetById(id string) (*menu.MenuItem, error) {
	res, err := is.GetReplica().Get(menu.MenuItem{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(MenuItemTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find menu item with id=%s", id)
	}

	return res.(*menu.MenuItem), nil
}

func (is *SqlMenuItemStore) GetByName(name string) (*menu.MenuItem, error) {
	var menuItem *menu.MenuItem

	err := is.GetReplica().SelectOne(&menuItem, "SELECT * FROM "+MenuItemTableName+" WHERE Name = :Name", map[string]interface{}{"Name": name})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(MenuItemTableName, "name="+name)
		}
		return nil, errors.Wrapf(err, "failed to find menu item with name=%s", name)
	}

	return menuItem, nil
}
