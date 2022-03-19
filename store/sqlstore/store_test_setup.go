package sqlstore

import (
	"os"
	"sync"
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/storetest"
)

var Name string

type storeType struct {
	Name        string
	SqlSettings *model.SqlSettings
	SqlStore    *SqlStore
	Store       store.Store
}

var storeTypes []*storeType

func newStoreType(name string) *storeType {
	return &storeType{
		Name:        name,
		SqlSettings: storetest.MakeSqlSettings(model.DATABASE_DRIVER_POSTGRES, false),
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
		t.Run(st.Name, func(t *testing.T) {
			if testing.Short() {
				t.SkipNow()
			}
			f(t, st.Store)
		})
	}
}

func StoreTestWithSqlStore(t *testing.T, f func(*testing.T, store.Store, storetest.SqlStore)) {
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

	if os.Getenv("IS_CI") == "true" {
		storeTypes = append(storeTypes, newStoreType("PostgreSQL"))
	} else {
		storeTypes = append(storeTypes, newStoreType("PostgreSQL"))
	}

	defer func() {
		if err := recover(); err != nil {
			tearDownStores()
			panic(err)
		}
	}()
	var wg sync.WaitGroup
	for _, st := range storeTypes {
		st := st
		wg.Add(1)
		go func() {
			defer wg.Done()
			st.SqlStore = New(*st.SqlSettings, nil)
			st.Store = st.SqlStore
			st.Store.DropAllTables()
			st.Store.MarkSystemRanUnitTests()
		}()
	}
	wg.Wait()
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
					storetest.CleanupSqlSettings(st.SqlSettings)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
