package storetest

import (
	"github.com/sitename/sitename/store/sqlstore"
	"gorm.io/gorm"
)

var _ SqlStore = (*sqlstore.SqlStore)(nil)

type SqlStore interface {
	GetMaster(noTimeout ...bool) *gorm.DB
	DriverName() string
}
