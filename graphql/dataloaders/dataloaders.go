//go:generate go run github.com/vektah/dataloaden AddressByIDLoader string github.com/sitename/sitename/graphql/gqlmodel.Address
package dataloaders

import (
	"sync"
	"time"

	"github.com/sitename/sitename/app"
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
	OrdersByUser         *OrdersByUser
	CustomerEventsByUser *CustomerEventsByUserLoader
	AddressByID          *AddressByIDLoader
}

// NewLoaders returns new pointer to a DataLoaders
func NewLoaders(a app.AppIface) *DataLoaders {
	once.Do(func() {
		loaders = &DataLoaders{
			OrdersByUser:         ordersByUserLoader(a.Srv()),
			CustomerEventsByUser: customerEventsByUserLoader(a.Srv()),
			AddressByID:          addressByIDLoader(a.Srv()),
		}
	})

	return loaders
}
