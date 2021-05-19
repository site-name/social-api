package sqlstore

// import (
// "github.com/sitename/sitename/store"
// )

type SqlOrderStore struct {
	*SqlStore
}

// func newSqlOrderStore(sqlStore *SqlStore) store.OrderStore {
// 	os := &SqlOrderStore{sqlStore}

// 	for _, db := range sqlStore.GetAllConns() {
// 		table := db.AddTableWithName(or)
// 	}

// 	return os
// }
