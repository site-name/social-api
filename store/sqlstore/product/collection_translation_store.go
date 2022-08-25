package product

import (
	"github.com/sitename/sitename/store"
)

type SqlCollectionTranslationStore struct {
	store.Store
}

func NewSqlCollectionTranslationStore(s store.Store) store.CollectionTranslationStore {
	return &SqlCollectionTranslationStore{s}
}
