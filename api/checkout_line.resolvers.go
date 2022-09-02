package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CheckoutLineDelete(ctx context.Context, args struct {
	CheckoutID *string
	LineID     *string
	Token      *string
}) (*gqlmodel.CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesAdd(ctx context.Context, args struct {
	CheckoutID *string
	Lines      []*gqlmodel.CheckoutLineInput
	Token      *string
}) (*gqlmodel.CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesUpdate(ctx context.Context, args struct {
	CheckoutID *string
	Lines      []*gqlmodel.CheckoutLineInput
	Token      *string
}) (*gqlmodel.CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLines(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
