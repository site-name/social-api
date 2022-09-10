package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) Translations(ctx context.Context, args struct {
	Kind   TranslatableKinds
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*TranslatableItemConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Translation(ctx context.Context, args struct {
	Id   string
	Kind TranslatableKinds
}) (TranslatableItem, error) {
	panic(fmt.Errorf("not implemented"))
}
