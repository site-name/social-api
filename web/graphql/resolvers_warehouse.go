package graphql

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/shared"
)

func createWarehouse(ctx context.Context, input WarehouseCreateInput) (*WarehouseCreate, error) {
	embeddedCtx := ctx.Value(shared.APIContextKey).(*shared.Context)

	fmt.Println(embeddedCtx.App != nil)

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
