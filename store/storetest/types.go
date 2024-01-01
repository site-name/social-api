package storetest

import (
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore"
)

var _ SqlStore = (*sqlstore.SqlStore)(nil)

type SqlStore interface {
	GetMaster() store.ContextRunner
	DriverName() string
}
