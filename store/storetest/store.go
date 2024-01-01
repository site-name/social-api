package storetest

import (
	"sync"
	"testing"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/sync/errgroup"
)

type storeType struct {
	Name        string
	SqlSettings *model_helper.SqlSettings
	SqlStore    *sqlstore.SqlStore
	Store       store.Store
}

var storeTypes []*storeType

// NOTE: name must be either "MySQL" or "PostgreSQL"
func newStoreType(name, driver string) *storeType {
	return &storeType{
		Name:        name,
		SqlSettings: MakeSqlSettings(driver, false),
	}
}

func StoreTest(t *testing.T, f func(*testing.T, store.Store)) {
	defer func() {
		if err := recover(); err != nil {
			tearDownStores()
			panic(err)
		}
	}()
	for _, st := range storeTypes {
		st := st
		t.Run(st.Name, func(t *testing.T) {
			if testing.Short() {
				t.SkipNow()
			}
			f(t, st.Store)
		})
	}
}

func StoreTestWithSqlStore(t *testing.T, f func(*testing.T, store.Store, SqlStore)) {
	defer func() {
		if err := recover(); err != nil {
			tearDownStores()
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

func initStores() {
	if testing.Short() {
		return
	}

	storeTypes = append(storeTypes, newStoreType("PostgreSQL", model_helper.DATABASE_DRIVER_POSTGRES))
	defer func() {
		if err := recover(); err != nil {
			tearDownStores()
			panic(err)
		}
	}()

	var eg errgroup.Group
	for _, st := range storeTypes {
		st := st
		eg.Go(func() error {
			var err error
			st.SqlStore = sqlstore.New(*st.SqlSettings, nil)
			if err != nil {
				return err
			}
			st.Store = st.SqlStore
			st.Store.DropAllTables()
			st.Store.MarkSystemRanUnitTests()

			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		panic(err)
	}
}

var tearDownStoresOnce sync.Once

func tearDownStores() {
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

func TestGetReplica(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Description                string
		DataSourceReplicaNum       int
		DataSourceSearchReplicaNum int
	}{
		{
			"no replicas",
			0,
			0,
		},
		{
			"one source replica",
			1,
			0,
		},
		{
			"multiple source replicas",
			3,
			0,
		},
		{
			"one source search replica",
			0,
			1,
		},
		{
			"multiple source search replicas",
			0,
			3,
		},
		{
			"one source replica, one source search replica",
			1,
			1,
		},
		{
			"one source replica, multiple source search replicas",
			1,
			3,
		},
		{
			"multiple source replica, one source search replica",
			3,
			1,
		},
		{
			"multiple source replica, multiple source search replicas",
			3,
			3,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			settings := MakeSqlSettings(model_helper.DATABASE_DRIVER_POSTGRES, false)

			dataSourceReplicas := []string{}
			dataSourceSearchReplicas := []string{}

			for i := 0; i < testCase.DataSourceReplicaNum; i++ {
				dataSourceReplicas = append(dataSourceReplicas, *settings.DataSource)
			}
			for i := 0; i < testCase.DataSourceSearchReplicaNum; i++ {
				dataSourceSearchReplicas = append(dataSourceSearchReplicas, *settings.DataSource)
			}

			settings.DataSourceReplicas = dataSourceReplicas
			settings.DataSourceSearchReplicas = dataSourceSearchReplicas

			store := sqlstore.New(*settings, nil)
			defer func() {
				store.Close()
				CleanupSqlSettings(settings)
			}()

			replicas := map[boil.ContextExecutor]bool{}
			for i := 0; i < 5; i++ {
				replicas[store.GetReplica()] = true
			}

			if testCase.DataSourceReplicaNum > 0 {
				// If replicas were defined, ensure none are the master.
				assert.Len(t, replicas, testCase.DataSourceReplicaNum)

				for repica := range replicas {
					assert.NotSame(t, store.GetMaster(), repica)
				}

			} else if assert.Len(t, replicas, 1) {
				// Otherwise ensure the replicas contains only the master.
				for replica := range replicas {
					assert.Same(t, store.GetMaster(), replica)
				}
			}

		})
	}
}
