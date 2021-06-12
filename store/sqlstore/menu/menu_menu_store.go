package menu

import (
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

type SqlMenuStore struct {
	store.Store
}

func NewSqlMenuStore(sqlStore store.Store) store.MenuStore {
	ms := &SqlMenuStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(menu.Menu{}, "Menus").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(menu.MENU_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(menu.MENU_SLUG_MAX_LENGTH).SetUnique(true)
	}

	return ms
}

func (ms *SqlMenuStore) CreateIndexesIfNotExists() {
	ms.CreateIndexIfNotExists("idx_menus_name", "Menus", "Name")
	ms.CreateIndexIfNotExists("idx_menus_slug", "Menus", "Slug")
	ms.CreateIndexIfNotExists("idx_menus_name_lower_textpattern", "Menus", "lower(Name) text_pattern_ops")
}
