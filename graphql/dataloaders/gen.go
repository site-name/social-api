package dataloaders

import (
	"time"

	"github.com/sitename/sitename/app"
)

type DataloaderContextKeyType string

const (
	DataloaderContextKey DataloaderContextKeyType = "dataloader_key"
)

type Loaders struct {
	OrdersByUser *OrdersByUserLoader
}

func NewLoaders(server *app.Server) *Loaders {
	var (
		wait     = 250 * time.Microsecond
		maxBatch = 100
	)

	return &Loaders{
		OrdersByUser: &OrdersByUserLoader{
			wait:     wait,
			maxBatch: maxBatch,
			fetch:    ordersByUserFetchCreator(server),
		},
	}
}
