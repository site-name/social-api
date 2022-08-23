package app

import (
	"github.com/sitename/sitename/model/app"
	"github.com/sitename/sitename/store"
)

type SqlAppTokenStore struct {
	store.Store
}

func NewSqlAppTokenStore(sqlStore store.Store) store.AppTokenStore {
	return &SqlAppTokenStore{
		Store: sqlStore,
	}

}

func (as *SqlAppTokenStore) Save(appToken *app.AppToken) (*app.AppToken, error) {
	panic("not implemented")
}
