package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"strings"
	"unsafe"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderAddNote(ctx context.Context, args struct {
	Order UUID
	Input OrderAddNoteInput
}) (*OrderAddNote, error) {
	args.Input.Message = strings.TrimSpace(args.Input.Message)
	if args.Input.Message == "" {
		return nil, model_helper.NewAppError("OrderAddNote", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Message"}, "please provide non empty message", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderAddNote", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

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
		return nil, model_helper.NewAppError("OrderAddNote", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderCancel", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

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
		return nil, model_helper.NewAppError("OrderCancel", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
	if args.Amount.ToDecimal().LessThanOrEqual(decimal.Zero) {
		return nil, model_helper.NewAppError("OrderCapture", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Amount"}, "amount should be a positive number.", http.StatusBadRequest)
	}
	decimalAmount := args.Amount.ToDecimal()

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
		return nil, model_helper.NewAppError("OrderCapture", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	paymentTransaction, paymentErr, appErr := embedCtx.App.Srv().PaymentService().Capture(tx, *lastPayment, pluginMng, order.ChannelID, &decimalAmount, nil, false)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, logAndReturnPaymentFailedAppError("OrderCapture", embedCtx, tx, paymentErr, order, lastPayment)
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	// Confirm that we changed the status to capture. Some payment can receive
	// asynchronous webhook with update status
	if paymentTransaction.Kind == model.TRANSACTION_KIND_CAPTURE {
		insufStockErr, appErr := embedCtx.App.Srv().OrderService().OrderCaptured(*order, user, nil, &decimalAmount, *lastPayment, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
		if insufStockErr != nil {
			return nil, insufStockErr.ToAppError("OrderCapture")
		}
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("OrderCapture", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model_helper.NewAppError("OrderConfirm", "app.order.order_status_different_than_unconfirmed.app_error", nil, "given order has status different than unconfirmed", http.StatusNotAcceptable)
	}
	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".OrderID = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orderLines) == 0 {
		return nil, model_helper.NewAppError("OrderConfirm", "app.order.order_has_no_lines.app_error", nil, "given order cotains no product", http.StatusNotAcceptable)
	}

	// begin transaction
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderConfirm", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

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
			return nil, model_helper.NewAppError("OrderConfirm", "app.order.payment_error.app_error", nil, pmError.Error(), http.StatusInternalServerError)
		}

		inSufStockErr, appErr := embedCtx.App.Srv().OrderService().OrderCaptured(*order, user, nil, nil, *lastPayment, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
		if inSufStockErr != nil {
			return nil, inSufStockErr.ToAppError("OrderConfirm")
		}
	}

	appErr = embedCtx.App.Srv().OrderService().OrderConfirmed(tx, *order, user, nil, pluginMng, true)
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("OrderConfirm", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model_helper.NewAppError("OrderFulfillmentCancel", "app.order.cancel_fulfillment_with_giftcards.app_error", nil, "cannot cancel fulfillment with giftcard lines", http.StatusNotAcceptable)
	}

	// validate fulfillment
	if !fulfillment.CanEdit() {
		return nil, model_helper.NewAppError("OrderFulfillmentCancel", "app.order.fulfillment_cannot_cancel.app_error", nil, "this fulfillment can not be canceled", http.StatusNotAcceptable)
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
		return nil, model_helper.NewAppError("OrderFulfillmentCancel", "app.order.fulfillment_require_warehouse.app_error", nil, "warehouse is required for this fulfillment", http.StatusNotAcceptable)
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
		return nil, insufStockErr.ToAppError("OrderFulfillmentApprove")
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

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderFulfillmentUpdateTracking(ctx context.Context, args struct {
	Id    UUID
	Input FulfillmentUpdateTrackingInput
}) (*FulfillmentUpdateTracking, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	fulfillment, appErr := embedCtx.App.Srv().OrderService().FulfillmentByOption(&model.FulfillmentFilterOption{
		Conditions:         squirrel.Expr(model.FulfillmentTableName+".Id = ?", args.Id),
		SelectRelatedOrder: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	if args.Input.TrackingNumber != nil {
		fulfillment.TrackingNumber = *args.Input.TrackingNumber
	}

	fulfillment, appErr = embedCtx.App.Srv().OrderService().UpsertFulfillment(nil, fulfillment)
	if appErr != nil {
		return nil, appErr
	}
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	var trackingNumber string
	if args.Input.TrackingNumber != nil {
		trackingNumber = *args.Input.TrackingNumber
	}

	appErr = embedCtx.App.Srv().OrderService().FulfillmentTrackingUpdated(fulfillment, user, nil, trackingNumber, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	if args.Input.NotifyCustomer != nil && *args.Input.NotifyCustomer {
		appErr = embedCtx.App.Srv().OrderService().SendFulfillmentUpdate(fulfillment.GetOrder(), fulfillment, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &FulfillmentUpdateTracking{
		Fulfillment: SystemFulfillmentToGraphqlFulfillment(fulfillment),
		Order:       SystemOrderToGraphqlOrder(fulfillment.GetOrder()),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderFulfillmentRefundProducts(ctx context.Context, args struct {
	Input OrderRefundProductsInput
	Order UUID
}) (*FulfillmentRefundProducts, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Order.String())
	if appErr != nil {
		return nil, appErr
	}

	lastPayment, appErr := embedCtx.App.Srv().PaymentService().GetLastOrderPayment(order.Id)
	if appErr != nil {
		return nil, appErr
	}

	// clean order payment
	appErr = cleanOrderPayment("OrderFulfillmentRefundProducts", lastPayment)
	if appErr != nil {
		return nil, appErr
	}

	amountToRefund := (*decimal.Decimal)(unsafe.Pointer(args.Input.AmountToRefund))
	appErr = cleanAmountToRefund(embedCtx, "OrderFulfillmentRefundProducts", order, lastPayment, amountToRefund)
	if appErr != nil {
		return nil, appErr
	}

	var (
		cleanedOrderLines      model.OrderLineDatas
		cleanedFulfillmentLins []*model.FulfillmentLineData
	)
	if len(args.Input.OrderLines) > 0 {
		orderLineRefundIfaces := lo.Map(args.Input.OrderLines, func(item *OrderRefundLineInput, _ int) orderLineReturnRefundLineCommon { return item })
		cleanedOrderLines, appErr = cleanLines(embedCtx, "OrderFulfillmentRefundProducts", orderLineRefundIfaces)
		if appErr != nil {
			return nil, appErr
		}
	}
	if len(args.Input.FulfillmentLines) > 0 {
		fulfillmentLineRefundIfaces := lo.Map(args.Input.FulfillmentLines, func(item *OrderRefundFulfillmentLineInput, _ int) orderRefundReturnFulfillmentLineCommon { return item })
		cleanedFulfillmentLins, appErr = cleanFulfillmentLines(embedCtx, "OrderFulfillmentRefundProducts", fulfillmentLineRefundIfaces, []model.FulfillmentStatus{
			model.FULFILLMENT_FULFILLED,
			model.FULFILLMENT_RETURNED,
			model.FULFILLMENT_WAITING_FOR_APPROVAL,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	requester, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	var includingShippingCosts bool
	if args.Input.IncludeShippingCosts != nil {
		includingShippingCosts = *args.Input.IncludeShippingCosts
	}

	fulfillment, paymentErr, appErr := embedCtx.App.Srv().OrderService().CreateRefundFulfillment(
		requester,
		nil,
		*order,
		*lastPayment,
		cleanedOrderLines,
		cleanedFulfillmentLins,
		embedCtx.App.Srv().Plugin.GetPluginManager(),
		amountToRefund,
		includingShippingCosts,
	)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, model_helper.NewAppError("OrderFulfillmentRefundProducts", model.ErrPayment, map[string]any{"Code": paymentErr.Code}, paymentErr.Error(), http.StatusInternalServerError)
	}

	return &FulfillmentRefundProducts{
		Order:       SystemOrderToGraphqlOrder(order),
		Fulfillment: SystemFulfillmentToGraphqlFulfillment(fulfillment),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderFulfillmentReturnProducts(ctx context.Context, args struct {
	Input OrderReturnProductsInput
	Order UUID
}) (*FulfillmentReturnProducts, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Order.String())
	if appErr != nil {
		return nil, appErr
	}

	lastPaymentOfOrder, appErr := embedCtx.App.Srv().PaymentService().GetLastOrderPayment(order.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = cleanOrderPayment("OrderFulfillmentReturnProducts", lastPaymentOfOrder)
	if appErr != nil {
		return nil, appErr
	}

	amountToRefund := (*decimal.Decimal)(unsafe.Pointer(args.Input.AmountToRefund))
	appErr = cleanAmountToRefund(embedCtx, "OrderFulfillmentReturnProducts", order, lastPaymentOfOrder, amountToRefund)
	if appErr != nil {
		return nil, appErr
	}

	var (
		cleanedOrderLines       model.OrderLineDatas
		cleanedFulfillmentLines []*model.FulfillmentLineData
	)
	if len(args.Input.OrderLines) > 0 {
		orderLineReturnIfaces := lo.Map(args.Input.OrderLines, func(item *OrderReturnLineInput, _ int) orderLineReturnRefundLineCommon { return item })
		cleanedOrderLines, appErr = cleanLines(embedCtx, "OrderFulfillmentReturnProducts", orderLineReturnIfaces)
		if appErr != nil {
			return nil, appErr
		}
	}
	if len(args.Input.FulfillmentLines) > 0 {
		fulfillmenLineReturnIfaces := lo.Map(args.Input.FulfillmentLines, func(item *OrderReturnFulfillmentLineInput, _ int) orderRefundReturnFulfillmentLineCommon { return item })
		cleanedFulfillmentLines, appErr = cleanFulfillmentLines(embedCtx, "OrderFulfillmentReturnProducts", fulfillmenLineReturnIfaces, []model.FulfillmentStatus{
			model.FULFILLMENT_FULFILLED,
			model.FULFILLMENT_REFUNDED,
			model.FULFILLMENT_WAITING_FOR_APPROVAL,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	// perform mutation
	requester, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	var refund, refundShipingCosts bool
	if args.Input.Refund != nil {
		refund = *args.Input.Refund
	}
	if args.Input.IncludeShippingCosts != nil {
		refundShipingCosts = *args.Input.Refund
	}
	returnFulfillment, replaceFulfillment, replaceOrder, paymentErr, appErr := embedCtx.App.Srv().OrderService().CreateFulfillmentsForReturnedProducts(
		requester,
		nil,
		*order,
		lastPaymentOfOrder,
		cleanedOrderLines,
		cleanedFulfillmentLines,
		embedCtx.App.Srv().Plugin.GetPluginManager(),
		refund,
		amountToRefund,
		refundShipingCosts,
	)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, model_helper.NewAppError("OrderFulfillmentReturnProducts", model.ErrPayment, map[string]any{"Code": paymentErr.Code}, paymentErr.Error(), http.StatusInternalServerError)
	}

	return &FulfillmentReturnProducts{
		Order:              SystemOrderToGraphqlOrder(order),
		ReplaceOrder:       SystemOrderToGraphqlOrder(replaceOrder),
		ReturnFulfillment:  SystemFulfillmentToGraphqlFulfillment(returnFulfillment),
		ReplaceFulfillment: SystemFulfillmentToGraphqlFulfillment(replaceFulfillment),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderMarkAsPaid(ctx context.Context, args struct {
	Id                   UUID
	TransactionReference *string
}) (*OrderMarkAsPaid, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}
	if order.BillingAddressID == nil {
		return nil, model_helper.NewAppError("OrderMarkAsPaid", "app.order.order_no_billing_address.app_error", nil, "order billing address is required to mark order as paid", http.StatusNotAcceptable)
	}

	appErr = embedCtx.App.Srv().OrderService().CleanMarkOrderAsPaid(order)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return nil, logAndReturnPaymentFailedAppError("OrderMarkAsPaid", embedCtx, nil, nil, order, nil)
	}

	requester, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	var externalReference string
	if args.TransactionReference != nil {
		externalReference = *args.TransactionReference
	}

	paymentErr, appErr := embedCtx.App.Srv().OrderService().MarkOrderAsPaid(*order, requester, nil, pluginMng, externalReference)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, model_helper.NewAppError("OrderMarkAsPaid", model.ErrPayment, map[string]any{"Code": paymentErr.Code}, paymentErr.Error(), http.StatusInternalServerError)
	}

	return &OrderMarkAsPaid{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderRefund(ctx context.Context, args struct {
	Amount PositiveDecimal
	Id     UUID
}) (*OrderRefund, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	amount := args.Amount.ToDecimal()
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, model_helper.NewAppError("OrderRefund", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Amount"}, "amount must be positive", http.StatusBadRequest)
	}

	appErr = cleanOrderRefund("OrderRefund", embedCtx.App, order)
	if appErr != nil {
		return nil, appErr
	}

	lastOrderPayment, appErr := embedCtx.App.Srv().PaymentService().GetLastOrderPayment(order.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = cleanRefundPayment("OrderRefund", lastOrderPayment)
	if appErr != nil {
		return nil, appErr
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderRefund", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	transaction, paymentErr, appErr := embedCtx.App.Srv().PaymentService().Refund(tx, *lastOrderPayment, pluginMng, order.ChannelID, &amount)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, logAndReturnPaymentFailedAppError("OrderRefund", embedCtx, tx, paymentErr, order, lastOrderPayment)
	}

	// create fufillment
	_, appErr = embedCtx.App.Srv().OrderService().UpsertFulfillment(tx, &model.Fulfillment{
		Status:            model.FULFILLMENT_REFUNDED,
		OrderID:           order.Id,
		TotalRefundAmount: &amount,
	})
	if appErr != nil {
		return nil, appErr
	}

	// Confirm that we changed the status to refund. Some payment can receive
	// asynchronous webhook with update status
	if transaction.Kind == model.TRANSACTION_KIND_REFUND {
		user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().OrderService().OrderRefunded(*order, user, nil, amount, *lastOrderPayment, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("OrderRefund", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderRefund{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: currently only shop staffs can update order.
// TODO: check if we can let orders' owners update themself.
// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderUpdate(ctx context.Context, args struct {
	Id    UUID
	Input OrderUpdateInput
}) (*OrderUpdate, error) {
	// validate params
	if args.Input.UserEmail != nil && !model.IsValidEmail(*args.Input.UserEmail) {
		return nil, model_helper.NewAppError("OrderUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "UserEmail"}, "please proide valid user email", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	if order.Status == model.ORDER_STATUS_DRAFT {
		return nil, model_helper.NewAppError("OrderUpdate", "api.order.wrong_method.app_error", nil, "use DraftOrderUpdate method instead", http.StatusBadRequest)
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderUpdate", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	// update addresses
	for _, addressInput := range []*AddressInput{
		args.Input.BillingAddress,
		args.Input.ShippingAddress,
	} {
		if addressInput != nil {
			appErr = addressInput.validate("OrderUpdate")
			if appErr != nil {
				return nil, appErr
			}

			var newAddress model.Address
			addressInput.PatchAddress(&newAddress)

			savedAddress, appErr := embedCtx.App.Srv().AccountService().UpsertAddress(tx, &newAddress)
			if appErr != nil {
				return nil, appErr
			}

			switch addressInput {
			case args.Input.BillingAddress:
				order.BillingAddressID = &savedAddress.Id

			case args.Input.ShippingAddress:
				order.ShippingAddressID = &savedAddress.Id
			}
		}
	}

	// update user
	if args.Input.UserEmail != nil {
		user, appErr := embedCtx.App.Srv().AccountService().GetUserByOptions(ctx, &model.UserFilterOptions{
			Conditions: squirrel.Expr(model.UserTableName+".Email = ?", *args.Input.UserEmail),
		})
		if appErr != nil {
			return nil, appErr
		}

		order.UserID = &user.Id
		order.UserEmail = *args.Input.UserEmail
	}

	// update order
	updatedOrder, appErr := embedCtx.App.Srv().OrderService().UpsertOrder(tx, order)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	shopSettings := embedCtx.App.Config().ShopSettings
	appErr = embedCtx.App.Srv().OrderService().UpdateOrderPrices(tx, *updatedOrder, pluginMng, *shopSettings.IncludeTaxesInPrice)
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("OrderUpdate", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = pluginMng.OrderUpdated(*updatedOrder)
	if appErr != nil {
		return nil, appErr
	}

	return &OrderUpdate{
		Order: SystemOrderToGraphqlOrder(updatedOrder),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderUpdateShipping(ctx context.Context, args struct {
	Order UUID
	Input OrderUpdateShippingInput
}) (*OrderUpdateShipping, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Order.String())
	if appErr != nil {
		return nil, appErr
	}

	if args.Input.ShippingMethod == nil {
		orderRequiresShipping, appErr := embedCtx.App.Srv().OrderService().OrderShippingIsRequired(order.Id)
		if appErr != nil {
			return nil, appErr
		}

		if !order.IsDraft() && orderRequiresShipping {
			return nil, model_helper.NewAppError("OrderUpdateShipping", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "ShippingMethod"}, "shipping method is required for this order", http.StatusBadRequest)
		}

		order.ShippingMethodID = nil
		order.ShippingPrice, _ = util.ZeroTaxedMoney(order.Currency)
		order.ShippingMethodName = nil
		updatedOrder, appErr := embedCtx.App.Srv().OrderService().UpsertOrder(nil, order)
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(nil, updatedOrder, map[string]any{})
		if appErr != nil {
			return nil, appErr
		}

		return &OrderUpdateShipping{
			Order: SystemOrderToGraphqlOrder(updatedOrder),
		}, nil
	}

	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Conditions: squirrel.Expr(model.ShippingMethodTableName+".Id = ?", args.Input.ShippingMethod),
	})
	if appErr != nil {
		return nil, appErr
	}
	appErr = cleanOrderUpdateShipping("", embedCtx.App, order, shippingMethod)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	order.ShippingMethodID = &shippingMethod.Id
	shippingPrice, appErr := pluginMng.CalculateOrderShipping(*order)
	if appErr != nil {
		return nil, appErr
	}
	shippingTaxRate, appErr := pluginMng.GetOrderShippingTaxRate(*order, *shippingPrice)
	if appErr != nil {
		return nil, appErr
	}

	order.ShippingTaxRate = shippingTaxRate
	order.ShippingPrice = shippingPrice
	order.ShippingMethodName = &shippingMethod.Name

	updatedOrder, appErr := embedCtx.App.Srv().OrderService().UpsertOrder(nil, order)
	if appErr != nil {
		return nil, appErr
	}

	shopSettings := embedCtx.App.Config().ShopSettings
	appErr = embedCtx.App.Srv().OrderService().UpdateOrderPrices(nil, *updatedOrder, pluginMng, *shopSettings.IncludeTaxesInPrice)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().OrderShippingUpdated(*updatedOrder, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	return &OrderUpdateShipping{
		Order: SystemOrderToGraphqlOrder(updatedOrder),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderVoid(ctx context.Context, args struct{ Id UUID }) (*OrderVoid, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	lastOrderPayment, appErr := embedCtx.App.Srv().PaymentService().GetLastOrderPayment(order.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = cleanVoidPayment("OrderVoid", lastOrderPayment)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	transaction, paymentErr, appErr := embedCtx.App.Srv().PaymentService().Void(nil, *lastOrderPayment, pluginMng, order.ChannelID)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, logAndReturnPaymentFailedAppError("OrderVoid", embedCtx, nil, paymentErr, order, lastOrderPayment)
	}

	// Confirm that we changed the status to void. Some payment can receive
	// asynchronous webhook with update status
	if transaction.Kind == model.TRANSACTION_KIND_VOID {
		user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().OrderService().OrderVoided(*order, user, nil, lastOrderPayment, pluginMng)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &OrderVoid{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderBulkCancel(ctx context.Context, args struct{ Ids []UUID }) (*OrderBulkCancel, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
		Conditions: squirrel.Eq{model.OrderTableName + ".Id": args.Ids},
	})
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	var totalOrdersCanceled int32 = 0

	for _, order := range orders {
		appErr = cleanOrderCancel("OrderBulkCancel", embedCtx.App, order)
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().OrderService().CancelOrder(nil, order, user, nil, pluginMng)
		if appErr != nil {
			return nil, appErr
		}

		totalOrdersCanceled++
	}

	return &OrderBulkCancel{
		Count: totalOrdersCanceled,
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderSettings(ctx context.Context) (*OrderSettings, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	shopSettings := embedCtx.App.Config().ShopSettings

	return &OrderSettings{
		AutomaticallyConfirmAllNewOrders:         *shopSettings.AutomaticallyConfirmAllNewOrders,
		AutomaticallyFulfillNonShippableGiftCard: *shopSettings.AutomaticallyFulfillNonShippableGiftcard,
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) Order(ctx context.Context, args struct{ Id UUID }) (*Order, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	return SystemOrderToGraphqlOrder(order), nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderByToken(ctx context.Context, args struct{ Token UUID }) (*Order, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
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

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderDiscountAdd(ctx context.Context, args struct {
	Input   OrderDiscountCommonInput
	OrderID UUID
}) (*OrderDiscountAdd, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.OrderID.String())
	if appErr != nil {
		return nil, appErr
	}

	if !(order.IsDraft() || order.IsUnconfirmed()) {
		return nil, model_helper.NewAppError("OrderDiscountAdd", "app.order.only_draft_and_unconfirmed_order_can_update.app_error", nil, "only draft and unconfirmed order can be modified", http.StatusBadRequest)
	}

	orderDiscounts, appErr := embedCtx.App.Srv().OrderService().GetOrderDiscounts(order)
	if appErr != nil {
		return nil, appErr
	}
	if len(orderDiscounts) > 0 {
		return nil, model_helper.NewAppError("OrderDiscountAdd", "app.order.order_already_has_discount.app_error", nil, "order already has discounts", http.StatusBadRequest)
	}

	order.PopulateNonDbFields()
	appErr = validateOrderDiscountInput("OrderDiscountAdd", order.UnDiscountedTotal.Gross, args.Input)
	if appErr != nil {
		return nil, appErr
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderDiscountAdd", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	var reason string
	if args.Input.Reason != nil {
		reason = *args.Input.Reason
	}
	value := args.Input.Value.ToDecimal()

	orderDiscount, appErr := embedCtx.App.Srv().OrderService().CreateOrderDiscountForOrder(tx, order, reason, args.Input.ValueType, &value)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(tx, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &embedCtx.AppContext.Session().UserId,
		Type:    model.ORDER_EVENT_TYPE_ORDER_DISCOUNT_ADDED,
		Parameters: model_types.JSONString{
			"discount": embedCtx.App.Srv().OrderService().PrepareDiscountObject(orderDiscount, nil),
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("OrderDiscountAdd", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderDiscountAdd{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderDiscountUpdate(ctx context.Context, args struct {
	DiscountID UUID
	Input      OrderDiscountCommonInput
}) (*OrderDiscountUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	orderDiscouts, appErr := embedCtx.App.Srv().DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		Conditions:   squirrel.Expr(model.OrderDiscountTableName+".Id = ?", args.DiscountID),
		PreloadOrder: true,
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orderDiscouts) == 0 {
		return nil, model_helper.NewAppError("OrderDiscountUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "DiscountID"}, "please provide valid order discount id", http.StatusBadRequest)
	}

	// validate order
	orderDiscount := orderDiscouts[0]
	order := orderDiscount.Order

	if args.Input.Reason == nil || *args.Input.Reason == "" {
		args.Input.Reason = orderDiscount.Reason
	}
	if args.Input.Value.ToDecimal().Equal(decimal.Zero) && orderDiscount.Value != nil {
		args.Input.Value = PositiveDecimal(*orderDiscount.Value)
	}
	if !args.Input.ValueType.IsValid() {
		args.Input.ValueType = orderDiscount.ValueType
	}

	if !(order.IsDraft() || order.IsUnconfirmed()) {
		return nil, model_helper.NewAppError("OrderDiscountUpdate", "app.order.only_draft_and_unconfirmed_order_can_update.app_error", nil, "only draft and unconfirmed orders can be updated", http.StatusBadRequest)
	}

	order.PopulateNonDbFields()
	appErr = validateOrderDiscountInput("OrderDiscountUpdate", order.UnDiscountedTotal.Gross, args.Input)
	if appErr != nil {
		return nil, appErr
	}

	orderDiscountBeforeUpdate := orderDiscount.DeepCopy()

	orderDiscount.Reason = args.Input.Reason
	orderDiscount.Value = model_helper.GetPointerOfValue(args.Input.Value.ToDecimal())
	orderDiscount.ValueType = args.Input.ValueType

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderDiscountUpdate", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	orderDiscount, appErr = embedCtx.App.Srv().DiscountService().UpsertOrderDiscount(tx, orderDiscount)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(tx, order, map[string]any{})
	if appErr != nil {
		return nil, appErr
	}

	if orderDiscountBeforeUpdate.ValueType != args.Input.ValueType ||
		(orderDiscountBeforeUpdate.Value != nil &&
			!orderDiscountBeforeUpdate.Value.Equal(args.Input.Value.ToDecimal())) {

		_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(tx, &model.OrderEventOption{
			OrderID: order.Id,
			UserID:  &embedCtx.AppContext.Session().UserId,
			Type:    model.ORDER_EVENT_TYPE_ORDER_DISCOUNT_UPDATED,
			Parameters: model_types.JSONString{
				"discount": embedCtx.App.Srv().OrderService().PrepareDiscountObject(orderDiscount, orderDiscountBeforeUpdate),
			},
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("OrderDiscountUpdate", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderDiscountUpdate{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrderFulfill(ctx context.Context, args struct {
	Input OrderFulfillInput
	Order UUID
}) (*OrderFulfill, error) {
	var notifyCustomer, allowStockToBeExceed bool
	if args.Input.NotifyCustomer != nil {
		notifyCustomer = *args.Input.NotifyCustomer
	}
	if args.Input.AllowStockToBeExceeded != nil {
		allowStockToBeExceed = *args.Input.AllowStockToBeExceeded
	}

	// validate duplicates
	lineIdsMeetMap := map[UUID]bool{}           // keys are order line ids
	totalQuantityForLineMap := map[string]int{} // keys are order line ids

	for _, lineInput := range args.Input.Lines {
		if lineIdsMeetMap[lineInput.OrderLineID] {
			return nil, model_helper.NewAppError("OrderFulfill", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Input"}, "duplicate order line ids detected", http.StatusBadRequest)
		}
		lineIdsMeetMap[lineInput.OrderLineID] = true

		warehouseIdsOfLineMeetMap := map[UUID]bool{} // keys are warehouse ids
		for _, stockInput := range lineInput.Stocks {
			if warehouseIdsOfLineMeetMap[stockInput.Warehouse] {
				return nil, model_helper.NewAppError("OrderFulfill", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Input"}, "duplicate warehouse ids detected", http.StatusBadRequest)
			}
			warehouseIdsOfLineMeetMap[stockInput.Warehouse] = true

			totalQuantityForLineMap[lineInput.OrderLineID.String()] += int(stockInput.Quantity)
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": lo.Keys(totalQuantityForLineMap)},
		Preload: []string{
			"ProductVariant",
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	for _, orderLine := range orderLines {
		if totalQuantityForLineMap[orderLine.Id] > orderLine.QuantityUnFulfilled() {
			return nil, model_helper.NewAppError("OrderFulfill", "app.order.quantity_to_fulfill_greater_than_quantity_unfulfilled.app_error", map[string]any{"OrderLine": orderLine.Id, "RemainToFulfill": orderLine.QuantityUnFulfilled()}, "required quantity to fulfill greater than remaining quantity to fulfill", http.StatusBadRequest)
		}
	}

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Order.String())
	if appErr != nil {
		return nil, appErr
	}

	// clean input
	shopSettings := embedCtx.App.Config().ShopSettings
	if !order.IsFullyPaid() && *shopSettings.FulfillmentAutoApprove && *shopSettings.FulfillmentAllowUnPaid {
		return nil, model_helper.NewAppError("OrderFulfill", "app.order.cannot_fulfill_unpaid_order.app_error", nil, "cannot fulfill unpaid order", http.StatusNotAcceptable)
	}

	// check lines for preorder
	if *shopSettings.FulfillmentAutoApprove {
		for _, orderLine := range orderLines {
			if orderLine.ProductVariant != nil && orderLine.ProductVariant.IsPreorderActive() {
				return nil, model_helper.NewAppError("OrderFulfill", "app.order.cannot_fulfill_preorder_variant.app_error", nil, "cannot fulfill preorder variant", http.StatusNotAcceptable)
			}
		}
	}

	// check total quantity of item
	if lo.Sum(lo.Values(totalQuantityForLineMap)) <= 0 {
		return nil, model_helper.NewAppError("OrderFulfill", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Input"}, "total fulfill quantity must be positive", http.StatusBadRequest)
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("OrderFulfill", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	if *shopSettings.FulfillmentAutoApprove {
		giftardLines := lo.Filter(orderLines, func(item *model.OrderLine, _ int) bool { return item.IsGiftcard })

		_, appErr = embedCtx.App.Srv().GiftcardService().GiftcardsCreate(
			tx,
			order,
			giftardLines,
			totalQuantityForLineMap,
			shopSettings,
			user,
			nil,
			pluginMng,
		)
		if appErr != nil {
			return nil, appErr
		}
	}

	orderLinesMap := lo.SliceToMap(orderLines, func(item *model.OrderLine) (string, *model.OrderLine) { return item.Id, item })
	linesForWarehouse := map[string][]*model.QuantityOrderLine{}

	for i := 0; i < min(len(args.Input.Lines), len(orderLinesMap)); i++ {
		line := args.Input.Lines[i]
		orderLine := orderLinesMap[line.OrderLineID.String()]

		for _, stockInput := range line.Stocks {
			if stockInput.Quantity > 0 {
				linesForWarehouse[stockInput.Warehouse.String()] = append(linesForWarehouse[stockInput.Warehouse.String()], &model.QuantityOrderLine{
					Quantity:  int(stockInput.Quantity),
					OrderLine: orderLine,
				})
			}
		}
	}

	fulfillments, insufficientStockErr, appErr := embedCtx.App.Srv().OrderService().CreateFulfillments(
		user,
		nil,
		order,
		linesForWarehouse,
		pluginMng,
		notifyCustomer,
		*shopSettings.FulfillmentAutoApprove,
		allowStockToBeExceed,
	)
	if appErr != nil {
		return nil, appErr
	}
	if insufficientStockErr != nil {
		return nil, insufficientStockErr.ToAppError("OrderFulfill")
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("OrderFulfill", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderFulfill{
		Order:        SystemOrderToGraphqlOrder(order),
		Fulfillments: systemRecordsToGraphql(fulfillments, SystemFulfillmentToGraphqlFulfillment),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) Orders(ctx context.Context, args struct {
	SortBy    *OrderSortingInput
	Filter    *OrderFilterInput
	ChannelID *UUID
	GraphqlParams
}) (*OrderCountableConnection, error) {
	// validate params
	var orderFilterOpts = new(model.OrderFilterOption)
	var appErr *model_helper.AppError

	if args.Filter != nil {
		orderFilterOpts, appErr = args.Filter.parse("Orders")
		if appErr != nil {
			return nil, appErr
		}
	}

	paginValues, appErr := args.GraphqlParams.Parse("Orders")
	if appErr != nil {
		return nil, appErr
	}

	orderFilterOpts.GraphqlPaginationValues = *paginValues
	orderFilterOpts.CountTotal = true // ask store to count total too

	// add filter for non-draft orders
	orderFilterOpts.Conditions = append(
		orderFilterOpts.Conditions.(squirrel.And),
		squirrel.Expr(model.OrderTableName+".Status != ?", model.ORDER_STATUS_DRAFT),
	)

	// check if filter by channel id too
	if args.ChannelID != nil {
		orderFilterOpts.Conditions = append(
			orderFilterOpts.Conditions.(squirrel.And),
			squirrel.Expr(model.OrderTableName+".ChannelID = ?", *args.ChannelID),
		)
	}

	if orderFilterOpts.GraphqlPaginationValues.OrderBy == "" {
		// default sort to order numbers
		orderSortFields := orderSortFieldsMap[OrderSortFieldNumber].fields

		if args.SortBy != nil {
			orderSortFields = orderSortFieldsMap[args.SortBy.Field].fields

			switch args.SortBy.Field {
			case OrderSortFieldCustomer:
				orderFilterOpts.AnnotateBillingAddressNames = true

			case OrderSortFieldPayment:
				orderFilterOpts.AnnotateLastPaymentChargeStatus = true
			}
		}

		ordering := args.GraphqlParams.orderDirection().String()
		orderFilterOpts.GraphqlPaginationValues.OrderBy = orderSortFields.Map(func(_ int, item string) string { return item + " " + ordering }).Join(", ")
	}

	// find orders
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalCount, orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(orderFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := orderSortFieldsMap[OrderSortFieldNumber].keyFunc
	if args.SortBy != nil {
		keyFunc = orderSortFieldsMap[args.SortBy.Field].keyFunc
	}
	connection := constructCountableConnection(orders, totalCount, args.GraphqlParams, keyFunc, SystemOrderToGraphqlOrder)

	return (*OrderCountableConnection)(unsafe.Pointer(connection)), nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) DraftOrders(ctx context.Context, args struct {
	SortBy *OrderSortingInput
	Filter *OrderDraftFilterInput
	GraphqlParams
}) (*OrderCountableConnection, error) {
	// validate params
	var orderFilterOpts = new(model.OrderFilterOption)
	var appErr *model_helper.AppError

	if args.Filter != nil {
		orderFilterOpts, appErr = args.Filter.parse("Orders")
		if appErr != nil {
			return nil, appErr
		}
	}

	paginValues, appErr := args.GraphqlParams.Parse("Orders")
	if appErr != nil {
		return nil, appErr
	}

	orderFilterOpts.GraphqlPaginationValues = *paginValues
	orderFilterOpts.CountTotal = true // ask store to count total too

	// add filter for draft orders only:
	orderFilterOpts.Conditions = append(
		orderFilterOpts.Conditions.(squirrel.And),
		squirrel.Expr(model.OrderTableName+".Status = ?", model.ORDER_STATUS_DRAFT),
	)

	if orderFilterOpts.GraphqlPaginationValues.OrderBy == "" {
		// default sort to order numbers
		orderSortFields := orderSortFieldsMap[OrderSortFieldNumber].fields

		if args.SortBy != nil {
			orderSortFields = orderSortFieldsMap[args.SortBy.Field].fields

			switch args.SortBy.Field {
			case OrderSortFieldCustomer:
				orderFilterOpts.AnnotateBillingAddressNames = true

			case OrderSortFieldPayment:
				orderFilterOpts.AnnotateLastPaymentChargeStatus = true
			}
		}

		ordering := args.GraphqlParams.orderDirection().String()
		orderFilterOpts.GraphqlPaginationValues.OrderBy = orderSortFields.Map(func(_ int, item string) string { return item + " " + ordering }).Join(", ")
	}

	// find orders
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalCount, orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(orderFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := orderSortFieldsMap[OrderSortFieldNumber].keyFunc
	if args.SortBy != nil {
		keyFunc = orderSortFieldsMap[args.SortBy.Field].keyFunc
	}
	connection := constructCountableConnection(orders, totalCount, args.GraphqlParams, keyFunc, SystemOrderToGraphqlOrder)

	return (*OrderCountableConnection)(unsafe.Pointer(connection)), nil
}

// NOTE: Please refer to ./graphql/schemas/order.graphqls for details on directives used
func (r *Resolver) OrdersTotal(ctx context.Context, args struct {
	Period    ReportingPeriod
	ChannelID UUID
}) (*TaxedMoney, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Expr(model.ChannelTableName+".Id = ?", args.ChannelID),
	})
	if appErr != nil {
		return nil, appErr
	}

	createTimeFilterOpt := reportingPeriodToDate(args.Period).UnixNano() / 1000
	orderFilterOpts := model.OrderFilterOption{
		Conditions: squirrel.And{
			squirrel.NotEq{
				model.OrderTableName + ".Status": []string{
					string(model.ORDER_STATUS_CANCELED),
					string(model.ORDER_STATUS_DRAFT),
				},
			},
			squirrel.Expr(model.OrderTableName+".ChannelID = ?", args.ChannelID),
			squirrel.Expr(model.OrderTableName+".CreateAt >= ?", createTimeFilterOpt),
		},
	}

	_, orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(&orderFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	orderTotal, _ := util.ZeroTaxedMoney(orders[0].Currency)

	for _, order := range orders {
		order.PopulateNonDbFields()
		orderTotal, _ = orderTotal.Add(order.Total)
	}

	return SystemTaxedMoneyToGraphqlTaxedMoney(orderTotal), nil
}
