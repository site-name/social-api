package sqlstore

import (
	"testing"

	"github.com/sitename/sitename/store/storetest"
)

func Test_TestUserStore(t *testing.T) {
	StoreTestWithSqlStore(t, storetest.TestUserStore)
}
