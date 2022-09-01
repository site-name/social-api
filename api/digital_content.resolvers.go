package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) DigitalContentCreate(ctx context.Context, input gqlmodel.DigitalContentUploadInput, variantID string) (*gqlmodel.DigitalContentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DigitalContentDelete(ctx context.Context, variantID string) (*gqlmodel.DigitalContentDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DigitalContentUpdate(ctx context.Context, input gqlmodel.DigitalContentInput, variantID string) (*gqlmodel.DigitalContentUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DigitalContentURLCreate(ctx context.Context, input gqlmodel.DigitalContentURLCreateInput) (*gqlmodel.DigitalContentURLCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DigitalContent(ctx context.Context, id string) (*gqlmodel.DigitalContent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DigitalContents(ctx context.Context, before *string, after *string, first *int, last *int) (*gqlmodel.DigitalContentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
