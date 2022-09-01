package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) OrderSettingsUpdate(ctx context.Context, args struct {
	input gqlmodel.OrderSettingsUpdateInput
}) (*gqlmodel.OrderSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardSettingsUpdate(ctx context.Context, args struct {
	input gqlmodel.GiftCardSettingsUpdateInput
}) (*gqlmodel.GiftCardSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderAddNote(ctx context.Context, args struct {
	order string
	input gqlmodel.OrderAddNoteInput
}) (*gqlmodel.OrderAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderCancel(ctx context.Context, args struct{ id string }) (*gqlmodel.OrderCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderCapture(ctx context.Context, args struct {
	amount *decimal.Decimal
	id     string
}) (*gqlmodel.OrderCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderConfirm(ctx context.Context, args struct{ id string }) (*gqlmodel.OrderConfirm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfill(ctx context.Context, args struct {
	input gqlmodel.OrderFulfillInput
	order *string
}) (*gqlmodel.OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentCancel(ctx context.Context, args struct {
	id    string
	input *gqlmodel.FulfillmentCancelInput
}) (*gqlmodel.FulfillmentCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentApprove(ctx context.Context, args struct {
	allowStockToBeExceeded *bool
	id                     string
	notifyCustomer         bool
}) (*gqlmodel.FulfillmentApprove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentUpdateTracking(ctx context.Context, args struct {
	id    string
	input gqlmodel.FulfillmentUpdateTrackingInput
}) (*gqlmodel.FulfillmentUpdateTracking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentRefundProducts(ctx context.Context, args struct {
	input gqlmodel.OrderRefundProductsInput
	order string
}) (*gqlmodel.FulfillmentRefundProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentReturnProducts(ctx context.Context, args struct {
	input gqlmodel.OrderReturnProductsInput
	order string
}) (*gqlmodel.FulfillmentReturnProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderMarkAsPaid(ctx context.Context, args struct {
	id                   string
	transactionReference *string
}) (*gqlmodel.OrderMarkAsPaid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderRefund(ctx context.Context, args struct {
	amount *decimal.Decimal
	id     string
}) (*gqlmodel.OrderRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.OrderUpdateInput
}) (*gqlmodel.OrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdateShipping(ctx context.Context, args struct {
	order string
	input gqlmodel.OrderUpdateShippingInput
}) (*gqlmodel.OrderUpdateShipping, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderVoid(ctx context.Context, args struct{ id string }) (*gqlmodel.OrderVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderBulkCancel(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.OrderBulkCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderSettings(ctx context.Context) (*gqlmodel.OrderSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Order(ctx context.Context, args struct{ id string }) (*gqlmodel.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Orders(ctx context.Context, args struct {
	sortBy  *gqlmodel.OrderSortingInput
	filter  *gqlmodel.OrderFilterInput
	channel *string
	before  *string
	after   *string
	first   *int
	last    *int
}) (*gqlmodel.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrders(ctx context.Context, args struct {
	sortBy *gqlmodel.OrderSortingInput
	filter *gqlmodel.OrderDraftFilterInput
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrdersTotal(ctx context.Context, args struct {
	period  *gqlmodel.ReportingPeriod
	channel *string
}) (*gqlmodel.TaxedMoney, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderByToken(ctx context.Context, args struct{ token uuid.UUID }) (*gqlmodel.Order, error) {
	panic(fmt.Errorf("not implemented"))
}
