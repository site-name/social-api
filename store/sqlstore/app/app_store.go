package app

import (
	"github.com/sitename/sitename/store"
)

type SqlAppStore struct {
	store.Store
}

func NewSqlAppStore(sqlStore store.Store) store.AppStore {
	return &SqlAppStore{sqlStore}
}
