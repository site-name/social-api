package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CheckoutLineDelete(ctx context.Context, checkoutID *string, lineID *string, token *uuid.UUID) (*gqlmodel.CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesAdd(ctx context.Context, checkoutID *string, lines []*gqlmodel.CheckoutLineInput, token *uuid.UUID) (*gqlmodel.CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesUpdate(ctx context.Context, checkoutID *string, lines []*gqlmodel.CheckoutLineInput, token *uuid.UUID) (*gqlmodel.CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLines(ctx context.Context, before *string, after *string, first *int, last *int) (*gqlmodel.CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
