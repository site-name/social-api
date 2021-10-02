package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) AppCreate(ctx context.Context, input gqlmodel.AppInput) (*gqlmodel.AppCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppUpdate(ctx context.Context, id string, input gqlmodel.AppInput) (*gqlmodel.AppUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDelete(ctx context.Context, id string) (*gqlmodel.AppDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenCreate(ctx context.Context, input gqlmodel.AppTokenInput) (*gqlmodel.AppTokenCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenDelete(ctx context.Context, id string) (*gqlmodel.AppTokenDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenVerify(ctx context.Context, token string) (*gqlmodel.AppTokenVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppInstall(ctx context.Context, input gqlmodel.AppInstallInput) (*gqlmodel.AppInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppRetryInstall(ctx context.Context, activateAfterInstallation *bool, id string) (*gqlmodel.AppRetryInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeleteFailedInstallation(ctx context.Context, id string) (*gqlmodel.AppDeleteFailedInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppFetchManifest(ctx context.Context, manifestURL string) (*gqlmodel.AppFetchManifest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppActivate(ctx context.Context, id string) (*gqlmodel.AppActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeactivate(ctx context.Context, id string) (*gqlmodel.AppDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AppsInstallations(ctx context.Context) ([]*gqlmodel.AppInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Apps(ctx context.Context, filter *gqlmodel.AppFilterInput, sortBy *gqlmodel.AppSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.AppCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) App(ctx context.Context, id *string) (*gqlmodel.App, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AppExtensions(ctx context.Context, filter *gqlmodel.AppExtensionFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.AppExtensionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AppExtension(ctx context.Context, id string) (*gqlmodel.AppExtension, error) {
	panic(fmt.Errorf("not implemented"))
}
