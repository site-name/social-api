package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) CheckoutLineDelete(ctx context.Context, args struct {
	CheckoutID *string
	LineID     *string
	Token      *string
}) (*CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesAdd(ctx context.Context, args struct {
	CheckoutID *string
	Lines      []*CheckoutLineInput
	Token      *string
}) (*CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesUpdate(ctx context.Context, args struct {
	CheckoutID *string
	Lines      []*CheckoutLineInput
	Token      *string
}) (*CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLines(ctx context.Context, args GraphqlParams) (*CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
