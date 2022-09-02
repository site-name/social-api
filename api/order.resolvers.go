package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) OrderSettingsUpdate(ctx context.Context, args struct {
	Input gqlmodel.OrderSettingsUpdateInput
}) (*gqlmodel.OrderSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardSettingsUpdate(ctx context.Context, args struct {
	Input gqlmodel.GiftCardSettingsUpdateInput
}) (*gqlmodel.GiftCardSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderAddNote(ctx context.Context, args struct {
	Order string
	Input gqlmodel.OrderAddNoteInput
}) (*gqlmodel.OrderAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderCancel(ctx context.Context, args struct{ Id string }) (*gqlmodel.OrderCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderCapture(ctx context.Context, args struct {
	Amount *decimal.Decimal
	Id     string
}) (*gqlmodel.OrderCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderConfirm(ctx context.Context, args struct{ Id string }) (*gqlmodel.OrderConfirm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfill(ctx context.Context, args struct {
	Input gqlmodel.OrderFulfillInput
	Order *string
}) (*gqlmodel.OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentCancel(ctx context.Context, args struct {
	Id    string
	Input *gqlmodel.FulfillmentCancelInput
}) (*gqlmodel.FulfillmentCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentApprove(ctx context.Context, args struct {
	AllowStockToBeExceeded *bool
	Id                     string
	NotifyCustomer         bool
}) (*gqlmodel.FulfillmentApprove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentUpdateTracking(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.FulfillmentUpdateTrackingInput
}) (*gqlmodel.FulfillmentUpdateTracking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentRefundProducts(ctx context.Context, args struct {
	Input gqlmodel.OrderRefundProductsInput
	Order string
}) (*gqlmodel.FulfillmentRefundProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentReturnProducts(ctx context.Context, args struct {
	Input gqlmodel.OrderReturnProductsInput
	Order string
}) (*gqlmodel.FulfillmentReturnProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderMarkAsPaid(ctx context.Context, args struct {
	Id                   string
	TransactionReference *string
}) (*gqlmodel.OrderMarkAsPaid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderRefund(ctx context.Context, args struct {
	Amount *decimal.Decimal
	Id     string
}) (*gqlmodel.OrderRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.OrderUpdateInput
}) (*gqlmodel.OrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdateShipping(ctx context.Context, args struct {
	Order string
	Input gqlmodel.OrderUpdateShippingInput
}) (*gqlmodel.OrderUpdateShipping, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderVoid(ctx context.Context, args struct{ Id string }) (*gqlmodel.OrderVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderBulkCancel(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.OrderBulkCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderSettings(ctx context.Context) (*gqlmodel.OrderSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Order(ctx context.Context, args struct{ Id string }) (*gqlmodel.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Orders(ctx context.Context, args struct {
	SortBy  *gqlmodel.OrderSortingInput
	Filter  *gqlmodel.OrderFilterInput
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrders(ctx context.Context, args struct {
	SortBy *gqlmodel.OrderSortingInput
	Filter *gqlmodel.OrderDraftFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrdersTotal(ctx context.Context, args struct {
	Period  *gqlmodel.ReportingPeriod
	Channel *string
}) (*gqlmodel.TaxedMoney, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderByToken(ctx context.Context, args struct{ Token string }) (*gqlmodel.Order, error) {
	panic(fmt.Errorf("not implemented"))
}
