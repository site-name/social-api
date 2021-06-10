package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) OrderSettingsUpdate(ctx context.Context, input gqlmodel.OrderSettingsUpdateInput) (*gqlmodel.OrderSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderAddNote(ctx context.Context, order string, input gqlmodel.OrderAddNoteInput) (*gqlmodel.OrderAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderCancel(ctx context.Context, id string) (*gqlmodel.OrderCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderCapture(ctx context.Context, amount string, id string) (*gqlmodel.OrderCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderConfirm(ctx context.Context, id string) (*gqlmodel.OrderConfirm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfill(ctx context.Context, input gqlmodel.OrderFulfillInput, order *string) (*gqlmodel.OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentCancel(ctx context.Context, id string, input gqlmodel.FulfillmentCancelInput) (*gqlmodel.FulfillmentCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentUpdateTracking(ctx context.Context, id string, input gqlmodel.FulfillmentUpdateTrackingInput) (*gqlmodel.FulfillmentUpdateTracking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentRefundProducts(ctx context.Context, input gqlmodel.OrderRefundProductsInput, order string) (*gqlmodel.FulfillmentRefundProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentReturnProducts(ctx context.Context, input gqlmodel.OrderReturnProductsInput, order string) (*gqlmodel.FulfillmentReturnProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderMarkAsPaid(ctx context.Context, id string, transactionReference *string) (*gqlmodel.OrderMarkAsPaid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderRefund(ctx context.Context, amount string, id string) (*gqlmodel.OrderRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderUpdate(ctx context.Context, id string, input gqlmodel.OrderUpdateInput) (*gqlmodel.OrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderUpdateShipping(ctx context.Context, order string, input *gqlmodel.OrderUpdateShippingInput) (*gqlmodel.OrderUpdateShipping, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderVoid(ctx context.Context, id string) (*gqlmodel.OrderVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderBulkCancel(ctx context.Context, ids []*string) (*gqlmodel.OrderBulkCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrderSettings(ctx context.Context) (*gqlmodel.OrderSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Order(ctx context.Context, id string) (*gqlmodel.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Orders(ctx context.Context, sortBy *gqlmodel.OrderSortingInput, filter *gqlmodel.OrderFilterInput, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DraftOrders(ctx context.Context, sortBy *gqlmodel.OrderSortingInput, filter *gqlmodel.OrderDraftFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrdersTotal(ctx context.Context, period *gqlmodel.ReportingPeriod, channel *string) (*gqlmodel.TaxedMoney, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrderByToken(ctx context.Context, token uuid.UUID) (*gqlmodel.Order, error) {
	panic(fmt.Errorf("not implemented"))
}
