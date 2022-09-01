package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) Translations(ctx context.Context, args struct {
	Kind   gqlmodel.TranslatableKinds
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.TranslatableItemConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Translation(ctx context.Context, args struct {
	Id   string
	Kind gqlmodel.TranslatableKinds
}) (gqlmodel.TranslatableItem, error) {
	panic(fmt.Errorf("not implemented"))
}
