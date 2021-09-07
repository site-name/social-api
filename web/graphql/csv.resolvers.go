package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
	"github.com/sitename/sitename/web/shared"
)

func (r *mutationResolver) ExportProducts(ctx context.Context, input gqlmodel.ExportProductsInput) (*gqlmodel.ExportProducts, error) {
	// authentication and permissions checks are already done, thank to directive.
	embedContext := ctx.Value(shared.APIContextKey).(*shared.Context)

	// parse export scope
	scope := map[string]interface{}{}
	switch input.Scope {
	case gqlmodel.ExportScopeIDS:
		if input.Ids == nil || len(input.Ids) == 0 {
			return nil, model.NewAppError("graphql.ExportProducts", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "IDs"}, "", http.StatusBadRequest)
		}
		scope["ids"] = input.Ids
	case gqlmodel.ExportScopeFilter:
		if input.Filter == nil {
			return nil, model.NewAppError("graphql.ExportProducts", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Filter"}, "", http.StatusBadRequest)
		}
		scope["filter"] = input.Filter
	default:
		scope["all"] = true
	}

	// parse export info
	exportInfo := map[string]interface{}{}
	if fields := len(input.ExportInfo.Fields); fields > 0 {
		exportInfo["fields"] = fields
	}

	newExportFile, appErr := r.Srv().CsvService().CreateExportFile(&csv.ExportFile{
		UserID: &embedContext.AppContext.Session().UserId, // embedContext is usable since current user is authenticated
	})
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = r.Srv().CsvService().CommonCreateExportEvent(&csv.ExportEvent{
		ExportFileID: newExportFile.Id,
		UserID:       &embedContext.AppContext.Session().UserId,
		Type:         csv.EXPORT_PENDING,
	})
	if appErr != nil {
		return nil, appErr
	}

}

func (r *queryResolver) ExportFile(ctx context.Context, id string) (*gqlmodel.ExportFile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFiles(ctx context.Context, filter *gqlmodel.ExportFileFilterInput, sortBy *gqlmodel.ExportFileSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.ExportFileCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
