package attribute

import (
	"github.com/sitename/sitename/store"
)

type SqlAttributeTranslationStore struct {
	store.Store
}

func NewSqlAttributeTranslationStore(s store.Store) store.AttributeTranslationStore {
	return &SqlAttributeTranslationStore{s}
}
