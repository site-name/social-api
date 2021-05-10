package sqlstore

import (
	"github.com/sitename/sitename/store"
)

type SqlCheckoutStore struct {
	*SqlStore
}

func newSqlCheckoutStore(sqlStore *SqlStore) store.CheckoutStore {
	cs := &SqlCheckoutStore{
		SqlStore: sqlStore,
	}
}
