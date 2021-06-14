package app

import (
	"flag"
	"testing"

	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/testlib"
)

var mainHelper *testlib.MainHelper
var replicaFlag bool

func TestMain(m *testing.M) {
	if f := flag.Lookup("mysql-replica"); f == nil {
		flag.BoolVar(&replicaFlag, "mysql-replica", false, "")
		flag.Parse()
	}

	var options = testlib.HelperOptions{
		EnableStore:     true,
		EnableResources: true,
		WithReadReplica: replicaFlag,
	}

	slog.DisableZap()

	mainHelper = testlib.NewMainHelperWithOptions(&options)
	defer mainHelper.Close()

	mainHelper.Main(m)
}
