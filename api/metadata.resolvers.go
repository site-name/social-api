package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) DeleteMetadata(ctx context.Context, args struct {
	id   string
	keys []string
}) (*gqlmodel.DeleteMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DeletePrivateMetadata(ctx context.Context, args struct {
	id   string
	keys []string
}) (*gqlmodel.DeletePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UpdateMetadata(ctx context.Context, args struct {
	id    string
	input []gqlmodel.MetadataInput
}) (*gqlmodel.UpdateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UpdatePrivateMetadata(ctx context.Context, args struct {
	id    string
	input []gqlmodel.MetadataInput
}) (*gqlmodel.UpdatePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}
