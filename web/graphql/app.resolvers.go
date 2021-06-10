package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) AppCreate(ctx context.Context, input AppInput) (*AppCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppUpdate(ctx context.Context, id string, input AppInput) (*AppUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDelete(ctx context.Context, id string) (*AppDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenCreate(ctx context.Context, input AppTokenInput) (*AppTokenCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenDelete(ctx context.Context, id string) (*AppTokenDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenVerify(ctx context.Context, token string) (*AppTokenVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppInstall(ctx context.Context, input AppInstallInput) (*AppInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppRetryInstall(ctx context.Context, activateAfterInstallation *bool, id string) (*AppRetryInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeleteFailedInstallation(ctx context.Context, id string) (*AppDeleteFailedInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppFetchManifest(ctx context.Context, manifestURL string) (*AppFetchManifest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppActivate(ctx context.Context, id string) (*AppActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeactivate(ctx context.Context, id string) (*AppDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AppsInstallations(ctx context.Context) ([]AppInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Apps(ctx context.Context, filter *AppFilterInput, sortBy *AppSortingInput, before *string, after *string, first *int, last *int) (*AppCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) App(ctx context.Context, id *string) (*App, error) {
	panic(fmt.Errorf("not implemented"))
}
