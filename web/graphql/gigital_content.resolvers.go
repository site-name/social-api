package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) DigitalContentCreate(ctx context.Context, input DigitalContentUploadInput, variantID string) (*DigitalContentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentDelete(ctx context.Context, variantID string) (*DigitalContentDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentUpdate(ctx context.Context, input DigitalContentInput, variantID string) (*DigitalContentUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentURLCreate(ctx context.Context, input DigitalContentURLCreateInput) (*DigitalContentURLCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DigitalContent(ctx context.Context, id string) (*DigitalContent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DigitalContents(ctx context.Context, before *string, after *string, first *int, last *int) (*DigitalContentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
