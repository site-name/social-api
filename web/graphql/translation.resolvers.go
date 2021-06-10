package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) Translations(ctx context.Context, kind TranslatableKinds, before *string, after *string, first *int, last *int) (*TranslatableItemConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Translation(ctx context.Context, id string, kind TranslatableKinds) (TranslatableItem, error) {
	panic(fmt.Errorf("not implemented"))
}
