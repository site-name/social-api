package attribute

import (
	"github.com/sitename/sitename/store"
)

type SqlAttributeValueTranslationStore struct {
	store.Store
}

func NewSqlAttributeValueTranslationStore(s store.Store) store.AttributeValueTranslationStore {
	return &SqlAttributeValueTranslationStore{s}
}
