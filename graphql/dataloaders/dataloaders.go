//go:generate go run github.com/vektah/dataloaden AddressByIDLoader string github.com/sitename/sitename/graphql/gqlmodel.Address
package dataloaders

import (
	"sync"
	"time"

	"github.com/graph-gophers/dataloader"
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
	AttributeLoader *dataloader.Loader
	AddressLoader   *dataloader.Loader
}

// NewLoaders returns new pointer to a DataLoaders
func NewLoaders(srv *app.Server) *DataLoaders {
	once.Do(func() {
		var (
			attributeReader = &attributeReader{srv}
			addressReader   = &addressReader{srv}
		)

		loaders = &DataLoaders{
			AddressLoader:   dataloader.NewBatchedLoader(addressReader.getAddresses),
			AttributeLoader: dataloader.NewBatchedLoader(attributeReader.getAttributes),
		}
	})

	return loaders
}
