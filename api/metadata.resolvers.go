package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) DeleteMetadata(ctx context.Context, args struct {
	Id   string
	Keys []string
}) (*DeleteMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DeletePrivateMetadata(ctx context.Context, args struct {
	Id   string
	Keys []string
}) (*DeletePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UpdateMetadata(ctx context.Context, args struct {
	Id    string
	Input []MetadataInput
}) (*UpdateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UpdatePrivateMetadata(ctx context.Context, args struct {
	Id    string
	Input []MetadataInput
}) (*UpdatePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}
