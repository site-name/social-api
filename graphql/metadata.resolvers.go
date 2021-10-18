package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) DeleteMetadata(ctx context.Context, id string, keys []string) (*gqlmodel.DeleteMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePrivateMetadata(ctx context.Context, id string, keys []string) (*gqlmodel.DeletePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateMetadata(ctx context.Context, id string, input []*gqlmodel.MetadataInput) (*gqlmodel.UpdateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePrivateMetadata(ctx context.Context, id string, input []*gqlmodel.MetadataInput) (*gqlmodel.UpdatePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}
