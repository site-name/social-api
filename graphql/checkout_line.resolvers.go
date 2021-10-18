package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) CheckoutLineDelete(ctx context.Context, checkoutID *string, lineID *string, token *uuid.UUID) (*gqlmodel.CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesAdd(ctx context.Context, checkoutID *string, lines []*gqlmodel.CheckoutLineInput, token *uuid.UUID) (*gqlmodel.CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesUpdate(ctx context.Context, checkoutID *string, lines []*gqlmodel.CheckoutLineInput, token *uuid.UUID) (*gqlmodel.CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CheckoutLines(ctx context.Context, before *string, after *string, first *int, last *int) (*gqlmodel.CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
