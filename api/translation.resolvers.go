package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) Translations(ctx context.Context, args struct {
	kind   gqlmodel.TranslatableKinds
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.TranslatableItemConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Translation(ctx context.Context, args struct {
	id   string
	kind gqlmodel.TranslatableKinds
}) (gqlmodel.TranslatableItem, error) {
	panic(fmt.Errorf("not implemented"))
}
