package app

import (
	"github.com/sitename/sitename/model"
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

func (as *SqlAppTokenStore) Save(appToken *model.AppToken) (*model.AppToken, error) {
	panic("not implemented")
}
