package api

import (
	"context"
	"net/http"
	"unsafe"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/draft_order.graphqls for details on directives used.
func (r *Resolver) DraftOrderComplete(ctx context.Context, args struct{ Id string }) (*DraftOrderComplete, error) {
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("DraftOrderComplete", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, args.Id+" is not a valid order id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, orders, appErr := embedCtx.App.Srv().
		OrderService().
		FilterOrdersByOptions(&model.OrderFilterOption{
			Conditions: squirrel.Expr(model.OrderTableName+".Id = ?", args.Id),
			Preload:    []string{"OrderLines.ProductVariant", "Channel"},
		})
	if appErr != nil {
		return nil, appErr
	}
	order := orders[0]

	country, appErr := embedCtx.App.Srv().OrderService().GetOrderCountry(order)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().ValidateDraftOrder(order)
	if appErr != nil {
		return nil, appErr
	}

	// update user fields
	if order.UserID != nil {
		user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, *order.UserID)
		if appErr != nil {
			return nil, appErr
		}
		order.UserEmail = user.Email
	} else if order.UserEmail != "" {
		user, appErr := embedCtx.App.Srv().AccountService().GetUserByOptions(ctx, &model.UserFilterOptions{
			Conditions: squirrel.Expr(model.UserTableName+".Email = ?", order.UserEmail),
		})
		if appErr != nil {
			return nil, appErr
		}
		order.UserID = &user.Id
	}
	order.Status = model.ORDER_STATUS_UNFULFILLED

	orderRequireShipping := lo.SomeBy(order.OrderLines, func(item *model.OrderLine) bool { return item != nil && item.IsShippingRequired })

	if !orderRequireShipping {
		order.ShippingMethodName = nil
		order.ShippingPrice, _ = util.ZeroTaxedMoney(order.Currency)

		if order.ShippingAddressID != nil {
			appErr = embedCtx.App.Srv().Store.Address().DeleteAddresses(nil, []string{*order.ShippingAddressID})
			if appErr != nil {
				return nil, appErr
			}
			order.ShippingAddressID = nil
		}
	}

	// save order
	savedOrder, appErr := embedCtx.App.Srv().OrderService().UpsertOrder(nil, order)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	for _, line := range order.OrderLines {
		if line == nil {
			continue
		}
		if line.ProductVariant != nil && *line.ProductVariant.TrackInventory && line.ProductVariant.IsPreorderActive() {
			lineData := &model.OrderLineData{
				Line:     *line,
				Quantity: line.Quantity,
				Variant:  line.ProductVariant,
			}
			var channelSlug string
			if savedOrder.Channel != nil {
				channelSlug = savedOrder.Channel.Slug
			} else {
				channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
					Conditions: squirrel.Expr(model.ChannelTableName+".Id = ?", savedOrder.ChannelID),
				})
				if appErr != nil {
					return nil, appErr
				}
				channelSlug = channel.Slug
			}

			inSufStockErr, appErr := embedCtx.App.Srv().WarehouseService().AllocateStocks(model.OrderLineDatas{lineData}, country, channelSlug, pluginMng, model.StringInterface{})
			if appErr != nil {
				return nil, appErr
			}
			if inSufStockErr != nil {
				return nil, inSufStockErr.ToAppError("DraftOrderComplete.AllocateStocks")
			}

			// allocate pre order
			inSufStockErr, appErr = embedCtx.App.Srv().WarehouseService().AllocatePreOrders(model.OrderLineDatas{lineData}, channelSlug)
			if appErr != nil {
				return nil, appErr
			}
			if inSufStockErr != nil {
				return nil, inSufStockErr.ToAppError("DraftOrderComplete.AllocatePreOrders")
			}
		}
	}

	requester, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	insufficientStockErr, appErr := embedCtx.App.Srv().OrderService().OrderCreated(nil, *savedOrder, requester, nil, pluginMng, true)
	if appErr != nil {
		return nil, appErr
	}
	if insufficientStockErr != nil {
		return nil, insufficientStockErr.ToAppError("DraftOrderComplete.AllocatePreOrders")
	}

	return &DraftOrderComplete{
		Order: SystemOrderToGraphqlOrder(savedOrder),
	}, nil
}

// NOTE: Refer to ./schemas/draft_order.graphqls for details on directives used.
func (r *Resolver) DraftOrderDelete(ctx context.Context, args struct{ Id string }) (*DraftOrderDelete, error) {
	// validate params
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("DraftOrderDelete", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid draft order id", http.StatusBadRequest)
	}

	// find order:
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	if order.Status != model.ORDER_STATUS_DRAFT {
		return nil, model_helper.NewAppError("DraftOrderDelete", "api.order.delete_non_draft_order.app_error", nil, "cannot delete non-draft order", http.StatusNotAcceptable)
	}

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model_helper.NewAppError("DraftOrderDelete", model.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tran)

	_, appErr = embedCtx.App.Srv().OrderService().DeleteOrders(tran, []string{args.Id})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tran.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("DraftOrderDelete", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.DraftOrderDeleted(*order)
	if appErr != nil {
		return nil, appErr
	}

	return &DraftOrderDelete{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Refer to ./schemas/draft_order.graphqls for details on directives used.
func (r *Resolver) DraftOrderBulkDelete(ctx context.Context, args struct{ Ids []string }) (*DraftOrderBulkDelete, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model_helper.IsValidId) {
		return nil, model_helper.NewAppError("DraftOrderBulkDelete", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Ids"}, "please provide valid draft order ids", http.StatusBadRequest)
	}

	// validate all orders are draft
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
		Conditions: squirrel.Eq{model.OrderTableName + ".Id": args.Ids},
	})
	if appErr != nil {
		return nil, appErr
	}

	for _, order := range orders {
		if order != nil && order.Status != model.ORDER_STATUS_DRAFT {
			return nil, model_helper.NewAppError("DraftOrderBulkDelete", "api.order.delete_non_draft_order.app_error", nil, "order with id="+order.Id+" is not draft order", http.StatusNotAcceptable)
		}
	}

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model_helper.NewAppError("DraftOrderBulkDelete", model.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tran)

	totalCount, appErr := embedCtx.App.Srv().OrderService().DeleteOrders(tran, args.Ids)
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tran.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("DraftOrderBulkDelete", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &DraftOrderBulkDelete{
		Count: *(*int32)(unsafe.Pointer(&totalCount)),
	}, nil
}

// NOTE: Refer to ./schemas/draft_order.graphqls for details on directives used.
func (r *Resolver) DraftOrderCreate(ctx context.Context, args struct {
	Input DraftOrderCreateInput
}) (*DraftOrderCreate, error) {
	// validate params
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	appErr := args.Input.validate("DraftOrderCreate", embedCtx)
	if appErr != nil {
		return nil, appErr
	}
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("DraftOrderCreate", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(transaction)

	var order model.Order
	appErr = args.Input.patchOrder(embedCtx, &order, transaction, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(transaction, &order, model.StringInterface{})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	err := transaction.Commit().Error
	if err != nil {
		return nil, model_helper.NewAppError("DraftOrderCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = pluginMng.DraftOrderCreated(order)
	if appErr != nil {
		return nil, appErr
	}

	return &DraftOrderCreate{
		Order: SystemOrderToGraphqlOrder(&order),
	}, nil
}

// NOTE: Refer to ./schemas/draft_order.graphqls for details on directives used.
func (r *Resolver) DraftOrderUpdate(ctx context.Context, args struct {
	Id    string
	Input DraftOrderInput
}) (*DraftOrderUpdate, error) {
	// validate params
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("DraftOrderUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid draft order id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = args.Input.validate(embedCtx, "DraftOrderUpdate")
	if appErr != nil {
		return nil, appErr
	}
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	// create transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("DraftOrderUpdate", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(transaction)

	appErr = args.Input.patchOrder(embedCtx, order, transaction, pluginMng, true)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().OrderService().RecalculateOrder(transaction, order, model.StringInterface{})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	err := transaction.Commit().Error
	if err != nil {
		return nil, model_helper.NewAppError("DraftOrderCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = pluginMng.DraftOrderUpdated(*order)
	if appErr != nil {
		return nil, appErr
	}

	return &DraftOrderUpdate{
		Order: SystemOrderToGraphqlOrder(order),
	}, nil
}
