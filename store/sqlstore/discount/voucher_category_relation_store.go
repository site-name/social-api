package discount

import "github.com/sitename/sitename/store"

type SqlVoucherCategoryStore struct {
	store.Store
}

func NewSqlVoucherCategoryStore(s store.Store) store.VoucherCategoryStore {
	// vcs := &SqlVoucherCategoryStore{s}

	// for _, db := range s.GetAllConns() {
	// 	table := db.AddTableWithName()
	// }
}
