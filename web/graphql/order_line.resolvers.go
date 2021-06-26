package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) OrderLinesCreate(ctx context.Context, id string, input []*gqlmodel.OrderLineCreateInput) (*gqlmodel.OrderLinesCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDelete(ctx context.Context, id string) (*gqlmodel.OrderLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineUpdate(ctx context.Context, id string, input gqlmodel.OrderLineInput) (*gqlmodel.OrderLineUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountAdd(ctx context.Context, input gqlmodel.OrderDiscountCommonInput, orderID string) (*gqlmodel.OrderDiscountAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountUpdate(ctx context.Context, discountID string, input gqlmodel.OrderDiscountCommonInput) (*gqlmodel.OrderDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountDelete(ctx context.Context, discountID string) (*gqlmodel.OrderDiscountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDiscountUpdate(ctx context.Context, input gqlmodel.OrderDiscountCommonInput, orderLineID string) (*gqlmodel.OrderLineDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDiscountRemove(ctx context.Context, orderLineID string) (*gqlmodel.OrderLineDiscountRemove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderLineResolver) DigitalContentURL(ctx context.Context, obj *gqlmodel.OrderLine) (*gqlmodel.DigitalContentURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderLineResolver) Thumbnail(ctx context.Context, obj *gqlmodel.OrderLine, size *int) (*gqlmodel.Image, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderLineResolver) Variant(ctx context.Context, obj *gqlmodel.OrderLine) (*gqlmodel.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderLineResolver) Allocations(ctx context.Context, obj *gqlmodel.OrderLine) ([]gqlmodel.Allocation, error) {
	panic(fmt.Errorf("not implemented"))
}

// OrderLine returns OrderLineResolver implementation.
func (r *Resolver) OrderLine() OrderLineResolver { return &orderLineResolver{r} }

type orderLineResolver struct{ *Resolver }
