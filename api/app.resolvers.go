package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) AppCreate(ctx context.Context, args struct{ Input gqlmodel.AppInput }) (*gqlmodel.AppCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.AppInput
}) (*gqlmodel.AppUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.AppDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppTokenCreate(ctx context.Context, args struct{ Input gqlmodel.AppTokenInput }) (*gqlmodel.AppTokenCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppTokenDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.AppTokenDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppTokenVerify(ctx context.Context, args struct{ Token string }) (*gqlmodel.AppTokenVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppInstall(ctx context.Context, args struct{ Input gqlmodel.AppInstallInput }) (*gqlmodel.AppInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppRetryInstall(ctx context.Context, args struct {
	ActivateAfterInstallation *bool
	Id                        string
}) (*gqlmodel.AppRetryInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppDeleteFailedInstallation(ctx context.Context, args struct{ Id string }) (*gqlmodel.AppDeleteFailedInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppFetchManifest(ctx context.Context, args struct{ ManifestURL string }) (*gqlmodel.AppFetchManifest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppActivate(ctx context.Context, args struct{ Id string }) (*gqlmodel.AppActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppDeactivate(ctx context.Context, args struct{ Id string }) (*gqlmodel.AppDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppsInstallations(ctx context.Context) ([]gqlmodel.AppInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Apps(ctx context.Context, args struct {
	Filter *gqlmodel.AppFilterInput
	SortBy *gqlmodel.AppSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.AppCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) App(ctx context.Context, args struct{ Id *string }) (*gqlmodel.App, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppExtensions(ctx context.Context, args struct {
	Filter *gqlmodel.AppExtensionFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.AppExtensionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppExtension(ctx context.Context, args struct{ Id string }) (*gqlmodel.AppExtension, error) {
	panic(fmt.Errorf("not implemented"))
}
