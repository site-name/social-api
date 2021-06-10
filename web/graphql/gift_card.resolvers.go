package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) GiftCardActivate(ctx context.Context, id string) (*gqlmodel.GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardCreate(ctx context.Context, input gqlmodel.GiftCardCreateInput) (*gqlmodel.GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardDeactivate(ctx context.Context, id string) (*gqlmodel.GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardUpdate(ctx context.Context, id string, input gqlmodel.GiftCardUpdateInput) (*gqlmodel.GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCard(ctx context.Context, id string) (*gqlmodel.GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCards(ctx context.Context, before *string, after *string, first *int, last *int) (*gqlmodel.GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
