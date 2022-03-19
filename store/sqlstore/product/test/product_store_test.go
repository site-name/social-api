package test

import (
	"testing"

	"github.com/sitename/sitename/store/sqlstore"
	"github.com/sitename/sitename/store/storetest"
)

func TestProductStore(t *testing.T) {
	sqlstore.StoreTestWithSqlStore(t, storetest.TestProductStore)

}
