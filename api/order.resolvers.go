package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/site-name/decimal"
)

func (r *Resolver) HomepageEvents(ctx context.Context, args GraphqlParams) (*OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderSettingsUpdate(ctx context.Context, args struct {
	Input OrderSettingsUpdateInput
}) (*OrderSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardSettingsUpdate(ctx context.Context, args struct {
	Input GiftCardSettingsUpdateInput
}) (*GiftCardSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderAddNote(ctx context.Context, args struct {
	Order string
	Input OrderAddNoteInput
}) (*OrderAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderCancel(ctx context.Context, args struct{ Id string }) (*OrderCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderCapture(ctx context.Context, args struct {
	Amount *decimal.Decimal
	Id     string
}) (*OrderCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderConfirm(ctx context.Context, args struct{ Id string }) (*OrderConfirm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfill(ctx context.Context, args struct {
	Input OrderFulfillInput
	Order *string
}) (*OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentCancel(ctx context.Context, args struct {
	Id    string
	Input *FulfillmentCancelInput
}) (*FulfillmentCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentApprove(ctx context.Context, args struct {
	AllowStockToBeExceeded *bool
	Id                     string
	NotifyCustomer         bool
}) (*FulfillmentApprove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentUpdateTracking(ctx context.Context, args struct {
	Id    string
	Input FulfillmentUpdateTrackingInput
}) (*FulfillmentUpdateTracking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentRefundProducts(ctx context.Context, args struct {
	Input OrderRefundProductsInput
	Order string
}) (*FulfillmentRefundProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentReturnProducts(ctx context.Context, args struct {
	Input OrderReturnProductsInput
	Order string
}) (*FulfillmentReturnProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderMarkAsPaid(ctx context.Context, args struct {
	Id                   string
	TransactionReference *string
}) (*OrderMarkAsPaid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderRefund(ctx context.Context, args struct {
	Amount *decimal.Decimal
	Id     string
}) (*OrderRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdate(ctx context.Context, args struct {
	Id    string
	Input OrderUpdateInput
}) (*OrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdateShipping(ctx context.Context, args struct {
	Order string
	Input OrderUpdateShippingInput
}) (*OrderUpdateShipping, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderVoid(ctx context.Context, args struct{ Id string }) (*OrderVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderBulkCancel(ctx context.Context, args struct{ Ids []string }) (*OrderBulkCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderSettings(ctx context.Context) (*OrderSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Order(ctx context.Context, args struct{ Id string }) (*Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Orders(ctx context.Context, args struct {
	SortBy  *OrderSortingInput
	Filter  *OrderFilterInput
	Channel *string
	GraphqlParams
}) (*OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrders(ctx context.Context, args struct {
	SortBy *OrderSortingInput
	Filter *OrderDraftFilterInput
	GraphqlParams
}) (*OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrdersTotal(ctx context.Context, args struct {
	Period  *ReportingPeriod
	Channel *string
}) (*TaxedMoney, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderByToken(ctx context.Context, args struct{ Token string }) (*Order, error) {
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
