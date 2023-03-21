package testlib

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/services/searchengine"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/searchlayer"
	"github.com/sitename/sitename/store/sqlstore"
	"github.com/sitename/sitename/store/storetest"
)

type MainHelper struct {
	Settings         *model.SqlSettings
	Store            store.Store
	SearchEngine     *searchengine.Broker
	SQLStore         *sqlstore.SqlStore
	ClusterInterface *FakeClusterInterface

	status           int
	testResourcePath string
	replicas         []string
}

type HelperOptions struct {
	EnableStore     bool
	EnableResources bool
	WithReadReplica bool
}

func NewMainHelper() *MainHelper {
	return NewMainHelperWithOptions(&HelperOptions{
		EnableStore:     true,
		EnableResources: true,
	})
}

func NewMainHelperWithOptions(options *HelperOptions) *MainHelper {
	var mainHelper MainHelper
	flag.Parse()

	util.TranslationsPreInit()

	if options != nil {
		if options.EnableStore && !testing.Short() {
			mainHelper.setupStore(options.WithReadReplica)
		}

		if options.EnableResources {
			mainHelper.setupResources()
		}
	}

	return &mainHelper
}

func (h *MainHelper) Main(m *testing.M) {
	if h.testResourcePath != "" {
		prevDir, err := os.Getwd()
		if err != nil {
			panic("Failed to get current working directory: " + err.Error())
		}

		err = os.Chdir(h.testResourcePath)
		if err != nil {
			panic(fmt.Sprintf("Failed to set current working directory to %s: %s", h.testResourcePath, err.Error()))
		}

		defer func() {
			err := os.Chdir(prevDir)
			if err != nil {
				panic(fmt.Sprintf("Failed to restore current working directory to %s: %s", prevDir, err.Error()))
			}
		}()
	}

	h.status = m.Run()
}

func (h *MainHelper) Close() error {
	if h.SQLStore != nil {
		h.SQLStore.Close()
	}
	if h.Settings != nil {
		storetest.CleanupSqlSettings(h.Settings)
	}
	if h.testResourcePath != "" {
		os.RemoveAll(h.testResourcePath)
	}

	if r := recover(); r != nil {
		log.Fatalln(r)
	}

	os.Exit(h.status)

	return nil
}

func (h *MainHelper) setupStore(withReadReplica bool) {
	driverName := os.Getenv("MM_SQLSETTINGS_DRIVERNAME")
	if driverName == "" {
		driverName = model.DATABASE_DRIVER_POSTGRES
	}

	h.Settings = storetest.MakeSqlSettings(driverName, withReadReplica)
	h.replicas = h.Settings.DataSourceReplicas

	config := &model.Config{}
	config.SetDefaults()

	h.SearchEngine = searchengine.NewBroker(config)
	h.ClusterInterface = &FakeClusterInterface{}
	h.SQLStore = sqlstore.New(*h.Settings, nil)
	h.Store = searchlayer.NewSearchLayer(&TestStore{
		h.SQLStore,
	}, h.SearchEngine, config)
}

func (h *MainHelper) setupResources() {
	var err error
	h.testResourcePath, err = SetupTestResources()
	if err != nil {
		panic("failed to setup test resources: " + err.Error())
	}
}
