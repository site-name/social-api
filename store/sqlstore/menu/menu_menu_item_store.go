package menu

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

type SqlMenuItemStore struct {
	store.Store
}

func NewSqlMenuItemStore(s store.Store) store.MenuItemStore {
	is := &SqlMenuItemStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(menu.MenuItem{}, store.MenuItemTableName).SetKeys(false, "Id")
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
	is.CreateForeignKeyIfNotExists(store.MenuItemTableName, "MenuID", store.MenuTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(store.MenuItemTableName, "ParentID", store.MenuItemTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(store.MenuItemTableName, "CategoryID", store.ProductCategoryTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(store.MenuItemTableName, "CollectionID", store.ProductCollectionTableName, "Id", true)
	is.CreateForeignKeyIfNotExists(store.MenuItemTableName, "PageID", store.PageTableName, "Id", true)
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
	var res menu.MenuItem
	err := is.GetReplica().SelectOne(&res, "SELECT * FROM "+store.MenuItemTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.MenuItemTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find menu item with id=%s", id)
	}

	return &res, nil
}

func (is *SqlMenuItemStore) GetByName(name string) (*menu.MenuItem, error) {
	var menuItem *menu.MenuItem

	err := is.GetReplica().SelectOne(&menuItem, "SELECT * FROM "+store.MenuItemTableName+" WHERE Name = :Name", map[string]interface{}{"Name": name})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.MenuItemTableName, "name="+name)
		}
		return nil, errors.Wrapf(err, "failed to find menu item with name=%s", name)
	}

	return menuItem, nil
}
