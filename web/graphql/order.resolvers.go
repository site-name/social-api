package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func (r *mutationResolver) OrderCapture(ctx context.Context, amount *decimal.Decimal, id string) (*gqlmodel.OrderCapture, error) {
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

func (r *mutationResolver) OrderRefund(ctx context.Context, amount *decimal.Decimal, id string) (*gqlmodel.OrderRefund, error) {
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

func (r *orderResolver) User(ctx context.Context, obj *gqlmodel.Order) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) BillingAddress(ctx context.Context, obj *gqlmodel.Order) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) ShippingAddress(ctx context.Context, obj *gqlmodel.Order) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) ShippingMethod(ctx context.Context, obj *gqlmodel.Order) (*gqlmodel.ShippingMethod, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Channel(ctx context.Context, obj *gqlmodel.Order) (*gqlmodel.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Voucher(ctx context.Context, obj *gqlmodel.Order) (*gqlmodel.Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) GiftCards(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Fulfillments(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.Fulfillment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Lines(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.OrderLine, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) AvailableShippingMethods(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.ShippingMethod, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Invoices(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.Invoice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Payments(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.Payment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) CanFinalize(ctx context.Context, obj *gqlmodel.Order, _ *string) (bool, error) {
	// if obj.Status == gqlmodel.OrderStatusDraft {
	// 	countryCode, appErr := r.GetOrderCountryCode()
	// }
	// return true, nil
	panic("not implt")
}

func (r *orderResolver) Events(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.OrderEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Discounts(ctx context.Context, obj *gqlmodel.Order) ([]*gqlmodel.OrderDiscount, error) {
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

// Order returns OrderResolver implementation.
func (r *Resolver) Order() OrderResolver { return &orderResolver{r} }

type orderResolver struct{ *Resolver }
