package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

// validateOrderCanBeEdited checks if given order's status is not `draft` nor `unconfirmed`.
// If yes, return non-nil error.
func validateOrderCanBeEdited(order *model.Order, where string) *model.AppError {
	if order.Status != model.ORDER_STATUS_DRAFT && order.Status != model.ORDER_STATUS_UNCONFIRMED {
		return model.NewAppError(where, "app.order.order_cant_update.app_error", nil, "only draft and unconfirmed orders can be edited", http.StatusNotAcceptable)
	}
	return nil
}

func validateOrderDiscountInput(where string, maxTotal *goprices.Money, input OrderDiscountCommonInput) *model.AppError {
	if input.ValueType == model.DISCOUNT_VALUE_TYPE_FIXED {
		if decimal.Decimal(input.Value).GreaterThan(maxTotal.Amount) {
			return model.NewAppError(where, "app.order.value_greater_than_max_total.app_error", nil, "The value cannot be higher than max total amount", http.StatusNotAcceptable)
		}
	} else if decimal.Decimal(input.Value).GreaterThan(decimal.NewFromInt(100)) {
		return model.NewAppError(where, "app.order.discount_greater_than_100.app_error", nil, "The value cannot be higher than 100", http.StatusNotAcceptable)
	}

	return nil
}

// NOTE: please refer to ./graphql/schemas/order_line.graphqls for details on directives used.
func (r *Resolver) OrderLinesCreate(ctx context.Context, args struct {
	Id    string // id of an order to which new order lines are added
	Input []*OrderLineCreateInput
}) (*OrderLinesCreate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("OrderLinesCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid order id", http.StatusBadRequest)
	}
	args.Input = lo.Filter(args.Input, func(item *OrderLineCreateInput, _ int) bool { return item != nil })
	if len(args.Input) == 0 {
		return nil, model.NewAppError("OrderLinesCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Input"}, "please provide variants to add", http.StatusBadRequest)
	}

	variantIds := make([]string, len(args.Input))
	for idx, input := range args.Input {
		appErr := input.validate("OrderLinesCreate")
		if appErr != nil {
			return nil, appErr
		}

		variantIds[idx] = input.VariantID
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	// validate order
	// NOTE: only draft and uconfirmed orders can be edited
	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id)
	if appErr != nil {
		return nil, appErr
	}
	appErr = validateOrderCanBeEdited(order, "OrderLinesCreate")
	if appErr != nil {
		return nil, appErr
	}

	// validate variants
	variants, appErr := embedCtx.App.Srv().ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": variantIds},
	})
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().ValidateProductIsPublishedInChannel(variants, order.ChannelID)
	if appErr != nil {
		return nil, appErr
	}
	appErr = embedCtx.App.Srv().ProductService().ValidateVariantsAvailableInChannel(variants.IDs(), order.ChannelID)
	if appErr != nil {
		return nil, appErr
	}

	requester, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	// add order lines
	var addedOrderLines = make(model.OrderLines, 0, len(args.Input))
	var linesToAdd = make(model.QuantityOrderLines, 0, len(args.Input))
	var variantsMap = lo.SliceToMap(variants, func(v *model.ProductVariant) (string, *model.ProductVariant) { return v.Id, v })

	for _, input := range args.Input {
		variant, ok := variantsMap[input.VariantID]
		if ok && variant != nil {
			orderLine, insufStockErr, appErr := embedCtx.App.Srv().OrderService().AddVariantToOrder(*order, *variant, int(input.Quantity), requester, nil, pluginMng, []*model.DiscountInfo{}, order.IsUnconfirmed())
			if appErr != nil {
				return nil, appErr
			}
			if insufStockErr != nil {
				return nil, insufStockErr.ToAppError("OrderLinesCreate")
			}

			addedOrderLines = append(addedOrderLines, orderLine)
			linesToAdd = append(linesToAdd, &model.QuantityOrderLine{
				Quantity:  int(input.Quantity),
				OrderLine: orderLine,
			})
		}
	}

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model.NewAppError("OrderLinesCreate", model.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tran)

	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(tran, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &requester.Id,
		Type:    model.ORDER_EVENT_TYPE_ADDED_PRODUCTS,
		Parameters: model.StringInterface{
			"lines": embedCtx.App.Srv().OrderService().LinesPerQuantityToLineObjectList(linesToAdd),
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(tran, order, map[string]interface{}{})
	if appErr != nil {
		return nil, appErr
	}

	// commit tran
	if err := tran.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderLinesCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusBadRequest)
	}

	return &OrderLinesCreate{
		Order:      SystemOrderToGraphqlOrder(order),
		OrderLines: systemRecordsToGraphql(addedOrderLines, SystemOrderLineToGraphqlOrderLine),
	}, nil
}

// NOTE: please refer to ./graphql/schemas/order_line.graphqls for details on directives used.
func (r *Resolver) OrderLineDelete(ctx context.Context, args struct{ Id string }) (*OrderLineDelete, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("OrderLineDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid order line id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".Id = ?", args.Id),
		Preload:    []string{"Order", "ProductVariant"},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orderLines) == 0 {
		return nil, nil
	}

	orderLine := orderLines[0]
	order := orderLine.Order

	// NOTE: can update only orders that are "draft" or "unconfirmed"
	appErr = validateOrderCanBeEdited(orderLine.Order, "OrderLineDelete")
	if appErr != nil {
		return nil, appErr
	}

	var warehouseId *string = nil
	if order.IsUnconfirmed() {
		allocations, appErr := embedCtx.App.Srv().WarehouseService().AllocationsByOption(&model.AllocationFilterOption{
			Conditions:           squirrel.Expr(model.AllocationTableName+".OrderLineID = ?", orderLine.Id),
			SelectedRelatedStock: true,
		})
		if appErr != nil {
			return nil, appErr
		}
		if len(allocations) > 0 {
			warehouseId = &allocations[0].Stock.WarehouseID
		}
	}

	lineInfo := &model.OrderLineData{
		Line:        *orderLine,
		Quantity:    orderLine.Quantity,
		Variant:     orderLine.ProductVariant,
		WarehouseID: warehouseId,
	}

	// begin transaction
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderLineDelete", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	insufStockErr, appErr := embedCtx.App.Srv().OrderService().DeleteOrderLine(tx, lineInfo, pluginMng)
	if appErr != nil {
		return nil, appErr
	}
	if insufStockErr != nil {
		return nil, insufStockErr.ToAppError("OrderLineDelete")
	}

	orderRequiresShipping, appErr := embedCtx.App.Srv().OrderService().OrderShippingIsRequired(orderLine.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	if !orderRequiresShipping {
		order.ShippingMethodID = nil
		order.ShippingMethodName = nil
		order.ShippingPrice, _ = util.ZeroTaxedMoney(order.Currency)

		_, appErr := embedCtx.App.Srv().OrderService().UpsertOrder(tx, order)
		if appErr != nil {
			return nil, appErr
		}
	}

	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(tx, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &embedCtx.AppContext.Session().UserId,
		Type:    model.ORDER_EVENT_TYPE_REMOVED_PRODUCTS,
		Parameters: model.StringInterface{
			"lines": embedCtx.App.Srv().OrderService().LinesPerQuantityToLineObjectList([]*model.QuantityOrderLine{
				{
					Quantity:  orderLine.Quantity,
					OrderLine: orderLine,
				},
			}),
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(tx, order, map[string]interface{}{})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderLineDelete", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderLineDelete{
		Order:     SystemOrderToGraphqlOrder(order),
		OrderLine: SystemOrderLineToGraphqlOrderLine(orderLine),
	}, nil
}

// NOTE: please refer to ./graphql/schemas/order_line.graphqls for details on directives used.
func (r *Resolver) OrderLineUpdate(ctx context.Context, args struct {
	Id    string
	Input OrderLineInput
}) (*OrderLineUpdate, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("OrderLineUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid order line id", http.StatusBadRequest)
	}
	if args.Input.Quantity < 0 {
		return nil, model.NewAppError("OrderLineUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Quantity"}, "quantity must be greater than or equal to 0", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".Id = ?", args.Id),
		Preload:    []string{"Order", "ProductVariant"},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orderLines) == 0 {
		return nil, nil
	}

	orderLine := orderLines[0]
	order := orderLine.Order

	appErr = validateOrderCanBeEdited(order, "OrderLineUpdate")
	if appErr != nil {
		return nil, appErr
	}

	var warehouseId *string = nil
	if order.IsUnconfirmed() {
		allocations, appErr := embedCtx.App.Srv().WarehouseService().AllocationsByOption(&model.AllocationFilterOption{
			Conditions:           squirrel.Expr(model.AllocationTableName+".OrderLineID = ?", args.Id),
			SelectedRelatedStock: true,
		})
		if appErr != nil {
			return nil, appErr
		}
		if len(allocations) > 0 {
			warehouseId = &allocations[0].Stock.WarehouseID
		}
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderLineUpdate", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Expr(model.ChannelTableName+".Id = ?", order.ChannelID),
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	lineInfo := &model.OrderLineData{
		Line:        *orderLine,
		Quantity:    int(args.Input.Quantity),
		Variant:     orderLine.ProductVariant,
		WarehouseID: warehouseId,
	}
	userId := embedCtx.AppContext.Session().UserId
	inSufStockErr, appErr := embedCtx.App.Srv().OrderService().ChangeOrderLineQuantity(tx, userId, nil, lineInfo, orderLine.Quantity, int(args.Input.Quantity), channel.Slug, pluginMng, true)
	if appErr != nil {
		return nil, appErr
	}
	if inSufStockErr != nil {
		return nil, inSufStockErr.ToAppError("OrderLineUpdate")
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(tx, order, map[string]interface{}{})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderLineUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderLineUpdate{
		Order:     SystemOrderToGraphqlOrder(order),
		OrderLine: SystemOrderLineToGraphqlOrderLine(orderLine),
	}, nil
}

// NOTE: please refer to ./graphql/schemas/order_line.graphqls for details on directives used.
func (r *Resolver) OrderDiscountDelete(ctx context.Context, args struct{ DiscountID string }) (*OrderDiscountDelete, error) {
	if !model.IsValidId(args.DiscountID) {
		return nil, model.NewAppError("OrderDiscountDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "DiscountID"}, "please provide valid discount id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	orderDiscounts, appErr := embedCtx.App.Srv().DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		Conditions:   squirrel.Expr(model.OrderDiscountTableName+".Id = ?", args.DiscountID),
		PreloadOrder: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	if len(orderDiscounts) == 0 {
		return nil, nil
	}

	orderDiscount := orderDiscounts[0]
	order := orderDiscount.Order

	appErr = validateOrderCanBeEdited(order, "OrderDiscountDelete")
	if appErr != nil {
		return nil, appErr
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderDiscountDelete", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	appErr = embedCtx.App.Srv().OrderService().RemoveOrderDiscountFromOrder(tx, order, orderDiscount)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(tx, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &embedCtx.AppContext.Session().UserId,
		Type:    model.ORDER_EVENT_TYPE_ORDER_DISCOUNT_DELETED,
		Parameters: model.StringInterface{
			"discount": embedCtx.App.Srv().OrderService().PrepareDiscountObject(orderDiscount, nil),
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(tx, order, map[string]interface{}{})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderDiscountDelete", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderDiscountDelete{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: please refer to ./graphql/schemas/order_line.graphqls for details on directives used.
func (r *Resolver) OrderLineDiscountUpdate(ctx context.Context, args struct {
	Input       OrderDiscountCommonInput
	OrderLineID string
}) (*OrderLineDiscountUpdate, error) {
	if !model.IsValidId(args.OrderLineID) {
		return nil, model.NewAppError("OrderLineDiscountUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "OrderLineID"}, "please provide valid order line id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".Id = ?", args.OrderLineID),
		Preload:    []string{"Order"},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orderLines) == 0 {
		return nil, nil
	}

	orderLine := orderLines[0]
	order := orderLine.Order

	value := (*decimal.Decimal)(unsafe.Pointer(&args.Input.Value))
	valueType := args.Input.ValueType
	if value.Equal(decimal.Zero) {
		value = orderLine.UnitDiscountValue
	}
	if !valueType.IsValid() {
		valueType = orderLine.UnitDiscountType
	}

	appErr = validateOrderCanBeEdited(order, "OrderLineDiscountUpdate")
	if appErr != nil {
		return nil, appErr
	}

	orderLine.PopulateNonDbFields() // NOTE: this code is needed to use "UnDiscountedUnitPrice" below

	appErr = validateOrderDiscountInput("OrderLineDiscountUpdate", orderLine.UnDiscountedUnitPrice.Gross, args.Input)
	if appErr != nil {
		return nil, appErr
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderLineDiscountUpdate", model.ErrorCommittingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	orderLineBeforeUpdate := orderLine.DeepCopy()

	var reason string
	if args.Input.Reason != nil {
		reason = *args.Input.Reason
	}

	shopSettings := embedCtx.App.Config().ShopSettings
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	appErr = embedCtx.App.Srv().OrderService().UpdateDiscountForOrderLine(tx, *orderLine, *order, reason, valueType, value, pluginMng, *shopSettings.IncludeTaxesInPrice)
	if appErr != nil {
		return nil, appErr
	}

	if !orderLineBeforeUpdate.UnitDiscountValue.Equal(*value) ||
		orderLineBeforeUpdate.UnitDiscountType != valueType {

		requester, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
		if appErr != nil {
			return nil, appErr
		}

		_, appErr = embedCtx.App.Srv().OrderService().OrderLineDiscountEvent(
			model.ORDER_EVENT_TYPE_ORDER_LINE_DISCOUNT_UPDATED,
			order,
			requester,
			orderLine,
			orderLineBeforeUpdate,
		)
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(tx, order, map[string]interface{}{})
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderLineDiscountUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderLineDiscountUpdate{
		OrderLine: SystemOrderLineToGraphqlOrderLine(orderLine),
		Order:     SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: please refer to ./graphql/schemas/order_line.graphqls for details on directives used.
func (r *Resolver) OrderLineDiscountRemove(ctx context.Context, args struct{ OrderLineID string }) (*OrderLineDiscountRemove, error) {
	if !model.IsValidId(args.OrderLineID) {
		return nil, model.NewAppError("OrderLineDiscountRemove", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "OrderLineID"}, "please provide valid order line id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".Id = ?", args.OrderLineID),
		Preload:    []string{"Order"},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(orderLines) == 0 {
		return nil, nil
	}

	orderLine := orderLines[0]
	order := orderLine.Order

	appErr = validateOrderCanBeEdited(order, "OrderLineDiscountRemove")
	if appErr != nil {
		return nil, appErr
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("OrderLineDiscountRemove", model.ErrorCommittingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	shopSettings := embedCtx.App.Config().ShopSettings

	appErr = embedCtx.App.Srv().OrderService().RemoveDiscountFromOrderLine(tx, *orderLine, *order, pluginMng, *shopSettings.IncludeTaxesInPrice)
	if appErr != nil {
		return nil, appErr
	}

	requester, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = embedCtx.App.Srv().OrderService().OrderLineDiscountEvent(model.ORDER_EVENT_TYPE_ORDER_LINE_DISCOUNT_REMOVED, order, requester, orderLine, nil)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(tx, order, map[string]interface{}{})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("OrderLineDiscountRemove", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderLineDiscountRemove{
		OrderLine: SystemOrderLineToGraphqlOrderLine(orderLine),
		Order:     SystemOrderToGraphqlOrder(order),
	}, nil
}
