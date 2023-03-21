package storetest_test

import (
	"testing"

	"github.com/sitename/sitename/store/storetest"
	"github.com/sitename/sitename/testlib"
)

func TestMain(m *testing.M) {

	mainHelper := testlib.NewMainHelperWithOptions(nil)
	defer mainHelper.Close()
	storetest.InitStores()

	mainHelper.Main(m)

	storetest.TearDownStores()
}
