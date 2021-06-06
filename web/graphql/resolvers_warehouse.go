package graphql

import (
	"context"
	// "fmt"
	// "github.com/sitename/sitename/web/consts"
)

func createWarehouse(ctx context.Context, input WarehouseCreateInput) (*WarehouseCreate, error) {
	return &WarehouseCreate{
		Errors: []WarehouseError{},
		Warehouse: &Warehouse{
			ID:          "sdhuiher988er-dsfdfdg",
			Name:        "This is the name",
			Slug:        "this-is-the-name",
			CompanyName: "Sitename",
			Email:       "example@gmail.com",
		},
	}, nil
}
