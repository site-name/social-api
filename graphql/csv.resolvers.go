package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/store"
)

func (r *exportEventResolver) User(ctx context.Context, obj *gqlmodel.ExportEvent) (*gqlmodel.User, error) {
	session, appErr := CheckUserAuthenticated("exportEventResolver.User", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if obj.UserID == nil || (session.UserId != *obj.UserID && !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageStaff)) {
		return nil, nil
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.SystemUserToGraphqlUser(user), nil
}

func (r *exportFileResolver) User(ctx context.Context, obj *gqlmodel.ExportFile) (*gqlmodel.User, error) {
	session, appErr := CheckUserAuthenticated("exportFileResolver.User", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if obj.UserID == nil || (session.UserId != *obj.UserID && !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageStaff)) {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.SystemUserToGraphqlUser(user), nil
}

func (r *exportFileResolver) URL(ctx context.Context, obj *gqlmodel.ExportFile) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *exportFileResolver) Events(ctx context.Context, obj *gqlmodel.ExportFile) ([]*gqlmodel.ExportEvent, error) {
	events, appErr := r.Srv().CsvService().ExportEventsByOption(&csv.ExportEventFilterOption{
		ExportFileID: squirrel.Eq{store.CsvExportEventTablename + ".ExportFileID": obj.ID},
	})
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.SystemExportEventsToGraphqlExportEvents(events), nil
}

func (r *mutationResolver) ExportProducts(ctx context.Context, input gqlmodel.ExportProductsInput) (*gqlmodel.ExportProducts, error) {
	panic("not implemented")
}

func (r *queryResolver) ExportFile(ctx context.Context, id string) (*gqlmodel.ExportFile, error) {
	panic(fmt.Errorf("not implemented"))
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
