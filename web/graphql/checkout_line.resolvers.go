package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CheckoutLineDelete(ctx context.Context, checkoutID string, lineID *string) (*CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesAdd(ctx context.Context, checkoutID string, lines []*CheckoutLineInput) (*CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesUpdate(ctx context.Context, checkoutID string, lines []*CheckoutLineInput) (*CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CheckoutLine(ctx context.Context, id *string) (*CheckoutLine, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CheckoutLines(ctx context.Context, before *string, after *string, first *int, last *int) (*CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
