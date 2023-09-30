package storetest

import (
	"github.com/sitename/sitename/store/sqlstore"
	"gorm.io/gorm"
)

type StoreTestWrapper struct {
	orig *sqlstore.SqlStore
}

func NewStoreTestWrapper(orig *sqlstore.SqlStore) *StoreTestWrapper {
	return &StoreTestWrapper{orig}
}

func (w *StoreTestWrapper) GetMaster(noTimeout ...bool) *gorm.DB {
	return w.orig.GetMaster(noTimeout...)
}

func (w *StoreTestWrapper) DriverName() string {
	return w.orig.DriverName()
}
