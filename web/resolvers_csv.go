package web

import (
	"context"

	dbmodel "github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/web/model"
)

func (m *mutationResolver) exportProducts(ctx context.Context, input model.ExportProductsInput) (*model.ExportProducts, error) {

	// check export scope:
	scope := make(map[string]interface{})
	switch input.Scope {
	case model.ExportScopeIDS:
		if len(input.Ids) == 0 {
			return &model.ExportProducts{
				Errors: []model.ExportError{
					{
						Field:   dbmodel.NewString("ids"),
						Message: dbmodel.NewString("You must provide at least one product id."),
						Code:    model.ExportErrorCodeRequired,
					},
				},
			}, nil
		}
		scope["ids"] = input.Ids
	case model.ExportScopeFilter:
		if input.Filter == nil {
			return &model.ExportProducts{
				Errors: []model.ExportError{
					{
						Field:   dbmodel.NewString("filter"),
						Message: dbmodel.NewString("You must provide at least one product id."),
						Code:    model.ExportErrorCodeRequired,
					},
				},
			}, nil
		}
		scope["filter"] = input.Filter
	case model.ExportScopeAll:
		scope["all"] = nil
	}

	// check export info
	exportInfo := make(map[string]interface{})
	if input.ExportInfo != nil {
		if len(input.ExportInfo.Fields) > 0 {
			exportInfo["fields"] = input.ExportInfo.Fields
		}
	}

	if len(input.ExportInfo.Attributes) > 0 {
		exportInfo["attributes"] = input.ExportInfo.Attributes
	}
	if len(input.ExportInfo.Warehouses) > 0 {
		exportInfo["warehouses"] = input.ExportInfo.Warehouses
	}
	if len(input.ExportInfo.Channels) > 0 {
		exportInfo["channels"] = input.ExportInfo.Channels
	}

	// create exfport file in database
	exportFile := &csv.ExportFile{}
}
