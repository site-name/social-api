package shipping

import (
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodTranslationStore struct {
	store.Store
}

func NewSqlShippingMethodTranslationStore(s store.Store) store.ShippingMethodTranslationStore {
	return &SqlShippingMethodTranslationStore{s}
}
