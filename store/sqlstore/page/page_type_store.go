package page

import (
	"github.com/sitename/sitename/store"
)

type SqlPageTypeStore struct {
	store.Store
}

func NewSqlPageTypeStore(s store.Store) store.PageTypeStore {
	return &SqlPageTypeStore{s}
}
