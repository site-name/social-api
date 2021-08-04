package menu

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

const (
	MenuTableName = "Menus"
)

type SqlMenuStore struct {
	store.Store
}

func NewSqlMenuStore(sqlStore store.Store) store.MenuStore {
	ms := &SqlMenuStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(menu.Menu{}, MenuTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(menu.MENU_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(menu.MENU_SLUG_MAX_LENGTH).SetUnique(true)
	}

	return ms
}

func (ms *SqlMenuStore) CreateIndexesIfNotExists() {
	ms.CreateIndexIfNotExists("idx_menus_name", MenuTableName, "Name")
	ms.CreateIndexIfNotExists("idx_menus_slug", MenuTableName, "Slug")
	ms.CreateIndexIfNotExists("idx_menus_name_lower_textpattern", MenuTableName, "lower(Name) text_pattern_ops")
}

func (ms *SqlMenuStore) Save(mnu *menu.Menu) (*menu.Menu, error) {
	mnu.PreSave()
	if err := mnu.IsValid(); err != nil {
		return nil, err
	}

	if err := ms.GetMaster().Insert(mnu); err != nil {
		if ms.IsUniqueConstraintError(err, []string{"Name", "menus_name_key", "idx_menus_name_unique"}) {
			return nil, store.NewErrInvalidInput(MenuTableName, "Name", mnu.Name)
		}
		if ms.IsUniqueConstraintError(err, []string{"Slug", "menus_slug_key", "idx_menus_slug_unique"}) {
			return nil, store.NewErrInvalidInput(MenuTableName, "Slug", mnu.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save menu with id=%s", mnu.Id)
	}

	return mnu, nil
}

func (ms *SqlMenuStore) GetById(id string) (*menu.Menu, error) {
	var res menu.Menu
	err := ms.GetReplica().SelectOne(&res, "SELECT * FROM "+store.MenuTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(MenuTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find menu with id=&s", id)
	}

	return &res, nil
}

func (ms *SqlMenuStore) GetByName(name string) (*menu.Menu, error) {
	var res *menu.Menu
	if err := ms.GetReplica().SelectOne(&res, "SELECT * FROM "+MenuTableName+" WHERE Name = :Name", map[string]interface{}{"Name": name}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(MenuTableName, "name="+name)
		}
		return nil, errors.Wrapf(err, "failed to find menu with name=%s", name)
	}

	return res, nil
}

func (ms *SqlMenuStore) GetBySlug(slug string) (*menu.Menu, error) {
	var res *menu.Menu
	if err := ms.GetReplica().SelectOne(&res, "SELECT * FROM "+MenuTableName+" WHERE Slug = :Slug", map[string]interface{}{"Slug": slug}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(MenuTableName, "name="+slug)
		}
		return nil, errors.Wrapf(err, "failed to find menu with slug=%s", slug)
	}

	return res, nil
}
