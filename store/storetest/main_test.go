package storetest_test

import (
	"testing"

	"github.com/sitename/sitename/modules/testlib"
	"github.com/sitename/sitename/store/storetest"
)

var mainHelper *testlib.MainHelper

func TestMain(m *testing.M) {
	mainHelper = testlib.NewMainHelperWithOptions(nil)
	defer mainHelper.Close()

	storetest.InitTest()
	mainHelper.Main(m)
	storetest.TearDownTest()
}
