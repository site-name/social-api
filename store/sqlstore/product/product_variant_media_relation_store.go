package product

import (
	"github.com/sitename/sitename/store"
)

type SqlVariantMediaStore struct {
	store.Store
}

func NewSqlVariantMediaStore(s store.Store) store.VariantMediaStore {
	return &SqlVariantMediaStore{s}
}
