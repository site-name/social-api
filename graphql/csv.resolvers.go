package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/web/shared"
)

func (r *exportEventResolver) User(ctx context.Context, obj *gqlmodel.ExportEvent) (*gqlmodel.User, error) {
	session, appErr := checkUserAuthenticated("exportEventResolver.User", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if obj.UserID == nil || (session.UserId != *obj.UserID && !r.Srv().AccountService().SessionHasPermissionTo(session, gqlmodel.SaleorGraphqlPermissionToSystemPermission(gqlmodel.PermissionEnumManageStaff))) {
		return nil, nil
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseUserToGraphqlUser(user), nil
}

func (r *exportFileResolver) User(ctx context.Context, obj *gqlmodel.ExportFile) (*gqlmodel.User, error) {
	session, appErr := checkUserAuthenticated("exportEventResolver.User", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if obj.UserID == nil || (session.UserId != *obj.UserID && !r.Srv().AccountService().SessionHasPermissionTo(session, gqlmodel.SaleorGraphqlPermissionToSystemPermission(gqlmodel.PermissionEnumManageStaff))) {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseUserToGraphqlUser(user), nil
}

func (r *exportFileResolver) URL(ctx context.Context, obj *gqlmodel.ExportFile) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *exportFileResolver) Events(ctx context.Context, obj *gqlmodel.ExportFile) ([]*gqlmodel.ExportEvent, error) {
	events, appErr := r.Srv().CsvService().ExportEventsByOption(&csv.ExportEventFilterOption{
		ExportFileID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: obj.ID,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.SystemExportEventsToGraphqlExportEvents(events), nil
}

func (r *mutationResolver) ExportProducts(ctx context.Context, input gqlmodel.ExportProductsInput) (*gqlmodel.ExportProducts, error) {
	// authentication and permissions checks are already done, thank to directive.
	embedContext := ctx.Value(shared.APIContextKey).(*shared.Context)

	// validate scope is provided properly
	switch input.Scope {
	case gqlmodel.ExportScopeIDS:
		if input.Ids == nil || len(input.Ids) == 0 {
			return nil, model.NewAppError("graphql.ExportProducts", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "IDs"}, "", http.StatusBadRequest)
		}
	case gqlmodel.ExportScopeFilter:
		if input.Filter == nil {
			return nil, model.NewAppError("graphql.ExportProducts", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Filter"}, "", http.StatusBadRequest)
		}
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

	// embed input to job's data
	exportInput, err := json.JSON.Marshal(input)
	if err != nil {
		return nil, model.NewAppError("graphql.ExportProducts", app.ErrorMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}
	_, appErr = r.Srv().Jobs.CreateJob(model.JOB_TYPE_EXPORT_CSV, map[string]string{
		"input":          string(exportInput),
		"export_file_id": newExportFile.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.ExportProducts{
		ExportFile: gqlmodel.SystemExportFileToGraphqlExportFile(newExportFile),
	}, nil
}

func (r *queryResolver) ExportFile(ctx context.Context, id string) (*gqlmodel.ExportFile, error) {
	// user authentication and permission checking are done in @directive already
	exportFile, appErr := r.Srv().CsvService().ExportFileById(id)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.SystemExportFileToGraphqlExportFile(exportFile), nil
}

func (r *queryResolver) ExportFiles(ctx context.Context, filter *gqlmodel.ExportFileFilterInput, sortBy *gqlmodel.ExportFileSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.ExportFileCountableConnection, error) {
	// user authentication and permission checking are done in @directive already
	panic("not implemented")
}

// ExportEvent returns generated.ExportEventResolver implementation.
func (r *Resolver) ExportEvent() generated.ExportEventResolver { return &exportEventResolver{r} }

// ExportFile returns generated.ExportFileResolver implementation.
func (r *Resolver) ExportFile() generated.ExportFileResolver { return &exportFileResolver{r} }

type exportEventResolver struct{ *Resolver }
type exportFileResolver struct{ *Resolver }
