package dataloaders

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graphql/gqlmodel"
)

func ordersByUserLoader(server *app.Server) *OrdersByUser {
	return &OrdersByUser{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) ([]*gqlmodel.Order, []error) {
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

func addressByIDLoader(server *app.Server) *AddressByIDLoader {
	return &AddressByIDLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) ([]gqlmodel.Address, []error) {
			panic("not implemented")
		},
	}
}
