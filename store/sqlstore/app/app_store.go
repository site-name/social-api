package app

import (
	"github.com/sitename/sitename/model/app"
	"github.com/sitename/sitename/store"
)

type SqlAppStore struct {
	store.Store
}

func NewSqlAppStore(sqlStore store.Store) store.AppStore {
	return &SqlAppStore{sqlStore}
}

func (as *SqlAppStore) Save(app *app.App) (*app.App, error) {
	panic("not implemented") // NOTE: fixme
}
