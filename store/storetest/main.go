package storetest

import (
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore"
	"gorm.io/gorm"
)

const (
	defaultPostgresqlDSN = "postgres://mmuser:mostest@localhost:5432/mattermost_test?sslmode=disable&connect_timeout=10"
)

type StoreTestWrapper struct {
	orig store.Store
}

func NewStoreTestWrapper(orig store.Store) *StoreTestWrapper {
	return &StoreTestWrapper{orig}
}

func (w *StoreTestWrapper) GetMaster() *gorm.DB {
	return w.orig.GetMaster()
}

func (w *StoreTestWrapper) DriverName() string {
	return model.DATABASE_DRIVER_POSTGRES
}

type storeType struct {
	Name        string
	SqlSettings *model.SqlSettings
	SqlStore    store.Store
	Store       store.Store
}

var storeTypes []*storeType

func newStoreType(name, driver string) *storeType {
	return &storeType{
		Name:        name,
		SqlSettings: MakeSqlSettings(driver, false),
	}
}

func getEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

// PostgresSQLSettings returns the database settings to connect to the PostgreSQL unittesting database.
// The database name is generated randomly and must be created before use.
func PostgreSQLSettings() *model.SqlSettings {
	dsn := getEnv("TEST_DATABASE_POSTGRESQL_DSN", defaultPostgresqlDSN)
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	// Generate a random database name
	dsnURL.Path = "db" + model.NewId()

	return databaseSettings("postgres", dsnURL.String())
}

func databaseSettings(driver, dataSource string) *model.SqlSettings {
	settings := &model.SqlSettings{
		DriverName:                        &driver,
		DataSource:                        &dataSource,
		DataSourceReplicas:                []string{},
		DataSourceSearchReplicas:          []string{},
		MaxIdleConns:                      new(int),
		ConnMaxLifetimeMilliseconds:       new(int),
		ConnMaxIdleTimeMilliseconds:       new(int),
		MaxOpenConns:                      new(int),
		Trace:                             model.NewPrimitive(false),
		AtRestEncryptKey:                  model.NewPrimitive(model.NewRandomString(32)),
		QueryTimeout:                      new(int),
		MigrationsStatementTimeoutSeconds: new(int),
	}
	*settings.MaxIdleConns = 10
	*settings.ConnMaxLifetimeMilliseconds = 3600000
	*settings.ConnMaxIdleTimeMilliseconds = 300000
	*settings.MaxOpenConns = 100
	*settings.QueryTimeout = 60
	*settings.MigrationsStatementTimeoutSeconds = 10

	return settings
}

// MakeSqlSettings creates a randomly named database and returns the corresponding sql settings
func MakeSqlSettings(driver string, withReplica bool) *model.SqlSettings {
	settings := PostgreSQLSettings()
	dbName := postgreSQLDSNDatabase(*settings.DataSource)

	if err := execAsRoot(settings, "CREATE DATABASE "+dbName); err != nil {
		panic("failed to create temporary database " + dbName + ": " + err.Error())
	}

	if err := execAsRoot(settings, "GRANT ALL PRIVILEGES ON DATABASE \""+dbName+"\" TO mmuser"); err != nil {
		panic("failed to grant mmuser permission to " + dbName + ":" + err.Error())
	}

	log("Created temporary " + driver + " database " + dbName)

	return settings
}

var tearDownStoresOnce sync.Once

func TearDownStores() {
	if testing.Short() {
		return
	}
	tearDownStoresOnce.Do(func() {
		var wg sync.WaitGroup
		wg.Add(len(storeTypes))
		for _, st := range storeTypes {
			st := st
			go func() {
				if st.Store != nil {
					st.Store.Close()
				}
				if st.SqlSettings != nil {
					CleanupSqlSettings(st.SqlSettings)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func StoreTestWithSqlStore(t *testing.T, f func(*testing.T, store.Store, SqlStore)) {
	defer func() {
		if err := recover(); err != nil {
			TearDownStores()
			panic(err)
		}
	}()
	for _, st := range storeTypes {
		st := st
		t.Run(st.Name, func(t *testing.T) {
			if testing.Short() {
				t.SkipNow()
			}
			f(t, st.Store, &StoreTestWrapper{st.SqlStore})
		})
	}
}

func InitStores() {
	if testing.Short() {
		return
	}

	storeTypes = append(storeTypes,
		newStoreType("PostgreSQL", model.DATABASE_DRIVER_POSTGRES),
	)

	defer func() {
		if err := recover(); err != nil {
			TearDownStores()
			panic(err)
		}
	}()
	var wg sync.WaitGroup
	for _, st := range storeTypes {
		st := st
		wg.Add(1)
		go func() {
			defer wg.Done()
			st.SqlStore = sqlstore.New(*st.SqlSettings, nil)
			st.Store = st.SqlStore
			st.Store.DropAllTables()
			st.Store.MarkSystemRanUnitTests()
		}()
	}
	wg.Wait()
}
