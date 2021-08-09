package menu

import (
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

type SqlMenuItemTranslationStore struct {
	store.Store
}

func NewSqlMenuItemTranslationStore(sqlStore store.Store) store.MenuItemTranslationStore {
	mits := &SqlMenuItemTranslationStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(menu.MenuItemTranslation{}, store.MenuItemTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(10)
		table.ColMap("MenuItemID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(menu.MENU_ITEM_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "MenuItemID")
	}

	return mits
}

func (mits *SqlMenuItemTranslationStore) CreateIndexesIfNotExists() {
	mits.CreateForeignKeyIfNotExists(store.MenuItemTranslationTableName, "MenuItemID", store.MenuItemTableName, "Id", true)
}
