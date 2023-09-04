package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderAddNote(ctx context.Context, args struct {
	Order UUID
	Input OrderAddNoteInput
}) (*OrderAddNote, error) {
	args.Input.Message = strings.TrimSpace(args.Input.Message)
	if args.Input.Message == "" {
		return nil, model.NewAppError("OrderAddNote", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Message"}, "please provide non empty message", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	tx := embedCtx.App.Srv().Store.GetMaster()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderAddNote", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Order.String())
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.Srv().Account.UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	orderEvent, appErr := embedCtx.App.Srv().OrderService().OrderNoteAddedEvent(tx, order, user, args.Input.Message)
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderAddNote", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderAddNote{
		Order: SystemOrderToGraphqlOrder(order),
		Event: SystemOrderEventToGraphqlOrderEvent(orderEvent),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderCancel(ctx context.Context, args struct{ Id UUID }) (*OrderCancel, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	tx := embedCtx.App.Srv().Store.GetMaster()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderCancel", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.Srv().Account.UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	appErr = embedCtx.App.Srv().OrderService().CancelOrder(tx, order, user, nil, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().GiftcardService().DeactivateOrderGiftcards(tx, order.Id, user, nil)
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderCancel", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderCancel{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderCapture(ctx context.Context, args struct {
	Amount PositiveDecimal
	Id     UUID
}) (*OrderCapture, error) {
	if args.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, model.NewAppError("OrderCapture", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Amount"}, "amount should be a positive number.", http.StatusBadRequest)
	}
	decimalAmount := (*decimal.Decimal)(unsafe.Pointer(&args.Amount))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	// NOTE: lastPayment can possibly be nil
	lastPayment, appErr := embedCtx.App.Srv().PaymentService().GetLastOrderPayment(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	appErr = cleanOrderCapture("api.OrderCapture", lastPayment)
	if appErr != nil {
		return nil, appErr
	}

	// begin
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderCapture", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	paymentTransaction, pmErr, appErr := embedCtx.App.Srv().PaymentService().Capture(tx, *lastPayment, pluginMng, order.ChannelID, decimalAmount, nil, false)
	if appErr != nil {
		return nil, appErr
	}
	if pmErr != nil {
		return nil, logAndReturnPaymentFailedAppError("OrderCapture", embedCtx, tx, pmErr, order, lastPayment)
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	// Confirm that we changed the status to capture. Some payment can receive
	// asynchronous webhook with update status
	if paymentTransaction.Kind == model.CAPTURE {
		insufStockErr, appErr := embedCtx.App.Srv().OrderService().OrderCaptured(*order, user, nil, decimalAmount, *lastPayment, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
		if insufStockErr != nil {
			return nil, embedCtx.App.Srv().OrderService().PrepareInsufficientStockOrderValidationAppError("OrderCapture", insufStockErr)
		}
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderCapture", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderCapture{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
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
	Amount PositiveDecimal
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
