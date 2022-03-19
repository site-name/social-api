package sqlstore_test

import (
	"testing"

	"github.com/sitename/sitename/modules/testlib"
	"github.com/sitename/sitename/store/sqlstore"
)

func TestMain(m *testing.M) {
	mainHelper := testlib.NewMainHelperWithOptions(nil)
	defer mainHelper.Close()

	sqlstore.InitTest()
	mainHelper.Main(m)
	sqlstore.TearDownTest()
}
