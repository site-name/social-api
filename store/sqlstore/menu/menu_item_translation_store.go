package menu

import (
	"github.com/sitename/sitename/store"
)

type SqlMenuItemTranslationStore struct {
	store.Store
}

func NewSqlMenuItemTranslationStore(sqlStore store.Store) store.MenuItemTranslationStore {
	return &SqlMenuItemTranslationStore{sqlStore}
}
