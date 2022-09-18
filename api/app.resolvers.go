package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) AppCreate(ctx context.Context, args struct{ Input AppInput }) (*AppCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppUpdate(ctx context.Context, args struct {
	Id    string
	Input AppInput
}) (*AppUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppDelete(ctx context.Context, args struct{ Id string }) (*AppDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppTokenCreate(ctx context.Context, args struct{ Input AppTokenInput }) (*AppTokenCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppTokenDelete(ctx context.Context, args struct{ Id string }) (*AppTokenDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppTokenVerify(ctx context.Context, args struct{ Token string }) (*AppTokenVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppInstall(ctx context.Context, args struct{ Input AppInstallInput }) (*AppInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppRetryInstall(ctx context.Context, args struct {
	ActivateAfterInstallation bool
	Id                        string
}) (*AppRetryInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppDeleteFailedInstallation(ctx context.Context, args struct{ Id string }) (*AppDeleteFailedInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppFetchManifest(ctx context.Context, args struct{ ManifestURL string }) (*AppFetchManifest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppActivate(ctx context.Context, args struct{ Id string }) (*AppActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppDeactivate(ctx context.Context, args struct{ Id string }) (*AppDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppsInstallations(ctx context.Context) ([]AppInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Apps(ctx context.Context, args struct {
	Filter *AppFilterInput
	SortBy *AppSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*AppCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) App(ctx context.Context, args struct{ Id *string }) (*App, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppExtensions(ctx context.Context, args struct {
	Filter *AppExtensionFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*AppExtensionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AppExtension(ctx context.Context, args struct{ Id string }) (*AppExtension, error) {
	panic(fmt.Errorf("not implemented"))
}
