package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) OrderLinesCreate(ctx context.Context, id string, input []*OrderLineCreateInput) (*OrderLinesCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDelete(ctx context.Context, id string) (*OrderLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineUpdate(ctx context.Context, id string, input OrderLineInput) (*OrderLineUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountAdd(ctx context.Context, input OrderDiscountCommonInput, orderID string) (*OrderDiscountAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountUpdate(ctx context.Context, discountID string, input OrderDiscountCommonInput) (*OrderDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountDelete(ctx context.Context, discountID string) (*OrderDiscountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDiscountUpdate(ctx context.Context, input OrderDiscountCommonInput, orderLineID string) (*OrderLineDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDiscountRemove(ctx context.Context, orderLineID string) (*OrderLineDiscountRemove, error) {
	panic(fmt.Errorf("not implemented"))
}
