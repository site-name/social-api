package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/model"
)

func (r *Resolver) OrderLinesCreate(ctx context.Context, args struct {
	Id    string
	Input []*OrderLineCreateInput
}) (*OrderLinesCreate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("OrderLinesCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid order id", http.StatusBadRequest)
	}

	panic("not implemented")
}

func (r *Resolver) OrderLineDelete(ctx context.Context, args struct{ Id string }) (*OrderLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderLineUpdate(ctx context.Context, args struct {
	Id    string
	Input OrderLineInput
}) (*OrderLineUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderDiscountAdd(ctx context.Context, args struct {
	Input   OrderDiscountCommonInput
	OrderID string
}) (*OrderDiscountAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderDiscountUpdate(ctx context.Context, args struct {
	DiscountID string
	Input      OrderDiscountCommonInput
}) (*OrderDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderDiscountDelete(ctx context.Context, args struct{ DiscountID string }) (*OrderDiscountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderLineDiscountUpdate(ctx context.Context, args struct {
	Input       OrderDiscountCommonInput
	OrderLineID string
}) (*OrderLineDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderLineDiscountRemove(ctx context.Context, args struct{ OrderLineID string }) (*OrderLineDiscountRemove, error) {
	panic(fmt.Errorf("not implemented"))
}
