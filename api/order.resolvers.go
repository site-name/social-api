package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
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

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderConfirm(ctx context.Context, args struct{ Id UUID }) (*OrderConfirm, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	if !order.IsUnconfirmed() {
		return nil, model.NewAppError("OrderConfirm", "app.order.order_status_different_than_unconfirmed.app_error", nil, "given order has status different than unconfirmed", http.StatusNotAcceptable)
	}
	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".OrderID = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orderLines) == 0 {
		return nil, model.NewAppError("OrderConfirm", "app.order.order_has_no_lines.app_error", nil, "given order cotains no product", http.StatusNotAcceptable)
	}

	// begin transaction
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderConfirm", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	// update order
	order.Status = model.ORDER_STATUS_UNFULFILLED

	order, appErr = embedCtx.App.Srv().OrderService().UpsertOrder(tx, order)
	if appErr != nil {
		return nil, appErr
	}

	lastPayment, appErr := embedCtx.App.Srv().PaymentService().GetLastOrderPayment(order.Id)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	paymentAuthorized, appErr := embedCtx.App.Srv().PaymentService().PaymentIsAuthorized(lastPayment.Id)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	if paymentAuthorized && lastPayment.CanCapture() {
		_, pmError, appErr := embedCtx.App.Srv().PaymentService().Capture(tx, *lastPayment, pluginMng, order.ChannelID, nil, nil, false)
		if appErr != nil {
			return nil, appErr
		}
		if pmError != nil {
			return nil, model.NewAppError("OrderConfirm", "app.order.payment_error.app_error", nil, pmError.Error(), http.StatusInternalServerError)
		}

		inSufStockErr, appErr := embedCtx.App.Srv().OrderService().OrderCaptured(*order, user, nil, nil, *lastPayment, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
		if inSufStockErr != nil {
			return nil, embedCtx.App.Srv().OrderService().PrepareInsufficientStockOrderValidationAppError("OrderConfirm", inSufStockErr)
		}
	}

	appErr = embedCtx.App.Srv().OrderService().OrderConfirmed(tx, *order, user, nil, pluginMng, true)
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderConfirm", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderConfirm{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderFulfillmentCancel(ctx context.Context, args struct {
	Id    UUID
	Input *FulfillmentCancelInput
}) (*FulfillmentCancel, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	fulfillment, appErr := embedCtx.App.Srv().OrderService().FulfillmentByOption(&model.FulfillmentFilterOption{
		Conditions:         squirrel.Expr(model.FulfillmentTableName+".Id = ?", args.Id),
		SelectRelatedOrder: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	orderHasGiftcards, appErr := embedCtx.App.Srv().GiftcardService().OrderHasGiftcardLines(fulfillment.GetOrder())
	if appErr != nil {
		return nil, appErr
	}

	if orderHasGiftcards {
		return nil, model.NewAppError("OrderFulfillmentCancel", "app.order.cancel_fulfillment_with_giftcards.app_error", nil, "cannot cancel fulfillment with giftcard lines", http.StatusNotAcceptable)
	}

	// validate fulfillment
	if !fulfillment.CanEdit() {
		return nil, model.NewAppError("OrderFulfillmentCancel", "app.order.fulfillment_cannot_cancel.app_error", nil, "this fulfillment can not be canceled", http.StatusNotAcceptable)
	}

	var warehouse *model.WareHouse = nil
	if args.Input != nil && args.Input.WarehouseID != nil {
		warehouse, appErr = embedCtx.App.Srv().WarehouseService().WarehouseByOption(&model.WarehouseFilterOption{
			Conditions: squirrel.Expr(model.WarehouseTableName+".Id = ?", *args.Input.WarehouseID),
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	if fulfillment.Status != model.FULFILLMENT_WAITING_FOR_APPROVAL && warehouse == nil {
		return nil, model.NewAppError("OrderFulfillmentCancel", "app.order.fulfillment_require_warehouse.app_error", nil, "warehouse is required for this fulfillment", http.StatusNotAcceptable)
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	if fulfillment.Status == model.FULFILLMENT_WAITING_FOR_APPROVAL {
		appErr = embedCtx.App.Srv().OrderService().CancelWaitingFulfillment(*fulfillment, user, nil, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
	} else {
		fulfillment, appErr = embedCtx.App.Srv().OrderService().CancelFulfillment(*fulfillment, user, nil, warehouse, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &FulfillmentCancel{
		Order:       SystemOrderToGraphqlOrder(fulfillment.GetOrder()),
		Fulfillment: SystemFulfillmentToGraphqlFulfillment(fulfillment),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderFulfillmentApprove(ctx context.Context, args struct {
	AllowStockToBeExceeded *bool
	Id                     UUID
	NotifyCustomer         bool
}) (*FulfillmentApprove, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	fulfillment, appErr := embedCtx.App.Srv().OrderService().FulfillmentByOption(&model.FulfillmentFilterOption{
		Conditions: squirrel.Expr(model.FulfillmentTableName+".Id = ?", args.Id.String()),
	})
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	shopSettings := embedCtx.App.Config().ShopSettings

	_, insufStockErr, appErr := embedCtx.App.Srv().OrderService().ApproveFulfillment(fulfillment, user, nil, pluginMng, shopSettings, args.NotifyCustomer, *args.AllowStockToBeExceeded)
	if appErr != nil {
		return nil, appErr
	}
	if insufStockErr != nil {
		return nil, embedCtx.App.Srv().OrderService().PrepareInsufficientStockOrderValidationAppError("OrderFulfillmentApprove", insufStockErr)
	}

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	return &FulfillmentApprove{
		Fulfillment: SystemFulfillmentToGraphqlFulfillment(fulfillment),
		Order:       SystemOrderToGraphqlOrder(order),
	}, nil
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

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderByToken(ctx context.Context, args struct{ Token UUID }) (*Order, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
		Conditions: squirrel.And{
			squirrel.Expr(model.OrderTableName+".Status != ?", model.ORDER_STATUS_DRAFT),
			squirrel.Expr(model.OrderTableName+".Token = ?", args.Token),
		},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orders) == 0 {
		return nil, nil
	}

	return SystemOrderToGraphqlOrder(orders[0]), nil
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

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderFulfill(ctx context.Context, args struct {
	Input OrderFulfillInput
	Order *UUID
}) (*OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}
