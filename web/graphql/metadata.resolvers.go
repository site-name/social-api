package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) DeleteMetadata(ctx context.Context, id string, keys []string) (*DeleteMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePrivateMetadata(ctx context.Context, id string, keys []string) (*DeletePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateMetadata(ctx context.Context, id string, input []MetadataInput) (*UpdateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePrivateMetadata(ctx context.Context, id string, input []MetadataInput) (*UpdatePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}
