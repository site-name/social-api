package dataloaders

import (
	"sync"
	"time"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graphql/gqlmodel"
)

type DataloaderContextKeyType string

const (
	// DataloaderContextKey is used to get embedded `DataLoaders`
	DataloaderContextKey DataloaderContextKeyType = "dataloader_key"
	maxBatch                                      = 100
	wait                                          = 250 * time.Microsecond
)

var (
	once    sync.Once
	loaders *DataLoaders
)

// DataLoaders contains all data loaders for the project
type DataLoaders struct {
	OrdersByUser         *OrdersByUserLoader
	CustomerEventsByUser *CustomerEventsByUserLoader
}

// NewLoaders returns new pointer to a DataLoaders
func NewLoaders(a app.AppIface) *DataLoaders {
	// initialize only once
	once.Do(func() {
		loaders = &DataLoaders{
			OrdersByUser:         ordersByUserLoader(a.Srv()),
			CustomerEventsByUser: customerEventsByUserLoader(a.Srv()),
		}
	})

	return loaders
}

func ordersByUserLoader(server *app.Server) *OrdersByUserLoader {
	return &OrdersByUserLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) ([][]*gqlmodel.Order, []error) {
			panic("not implemented")
		},
	}
}

func customerEventsByUserLoader(server *app.Server) *CustomerEventsByUserLoader {
	return &CustomerEventsByUserLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) ([][]*gqlmodel.CustomerEvent, []error) {
			panic("not implemented")
		},
	}
}
