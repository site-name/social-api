package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) OrderLinesCreate(ctx context.Context, args struct {
	id    string
	input []*gqlmodel.OrderLineCreateInput
}) (*gqlmodel.OrderLinesCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderLineDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.OrderLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderLineUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.OrderLineInput
}) (*gqlmodel.OrderLineUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderDiscountAdd(ctx context.Context, args struct {
	input   gqlmodel.OrderDiscountCommonInput
	orderID string
}) (*gqlmodel.OrderDiscountAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderDiscountUpdate(ctx context.Context, args struct {
	discountID string
	input      gqlmodel.OrderDiscountCommonInput
}) (*gqlmodel.OrderDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderDiscountDelete(ctx context.Context, args struct{ discountID string }) (*gqlmodel.OrderDiscountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderLineDiscountUpdate(ctx context.Context, args struct {
	input       gqlmodel.OrderDiscountCommonInput
	orderLineID string
}) (*gqlmodel.OrderLineDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderLineDiscountRemove(ctx context.Context, args struct{ orderLineID string }) (*gqlmodel.OrderLineDiscountRemove, error) {
	panic(fmt.Errorf("not implemented"))
}
