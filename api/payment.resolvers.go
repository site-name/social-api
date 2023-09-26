package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/payment.graphqls for details on directives used.
func (r *Resolver) PaymentCapture(ctx context.Context, args struct {
	Amount    *decimal.Decimal
	PaymentID UUID
}) (*PaymentCapture, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("PaymentCapture", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	payment, appErr := embedCtx.App.Srv().PaymentService().PaymentByID(tx, args.PaymentID.String(), true)
	if appErr != nil {
		return nil, appErr
	}

	var channelID string
	switch {
	case payment.OrderID != nil:
		order, appErr := embedCtx.App.Srv().OrderService().OrderById(*payment.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = order.ChannelID

	case payment.CheckoutID != nil:
		checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
			Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", *payment.CheckoutID),
		})
		if appErr != nil {
			return nil, appErr
		}
		channelID = checkout.ChannelID
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, paymentErr, appErr := embedCtx.App.Srv().PaymentService().Capture(tx, *payment, pluginMng, channelID, args.Amount, nil, false)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, model.NewAppError("PaymentCapture", model.ErrPayment, map[string]interface{}{"Code": paymentErr.Code}, paymentErr.Error(), http.StatusInternalServerError)
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("PaymentCapture", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &PaymentCapture{
		Payment: SystemPaymentToGraphqlPayment(payment),
	}, nil
}

// NOTE: Refer to ./schemas/payment.graphqls for details on directives used.
func (r *Resolver) PaymentRefund(ctx context.Context, args struct {
	Amount    *decimal.Decimal
	PaymentID UUID
}) (*PaymentRefund, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("PaymentRefund", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	payment, appErr := embedCtx.App.Srv().PaymentService().PaymentByID(tx, args.PaymentID.String(), true)
	if appErr != nil {
		return nil, appErr
	}

	var channelID string
	switch {
	case payment.OrderID != nil:
		order, appErr := embedCtx.App.Srv().OrderService().OrderById(*payment.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = order.ChannelID

	case payment.CheckoutID != nil:
		checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
			Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", *payment.CheckoutID),
		})
		if appErr != nil {
			return nil, appErr
		}
		channelID = checkout.ChannelID
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, paymentErr, appErr := embedCtx.App.Srv().PaymentService().Refund(tx, *payment, pluginMng, channelID, args.Amount)
	if appErr != nil {
		return nil, appErr
	}

	if paymentErr != nil {
		return nil, model.NewAppError("PaymentRefund", model.ErrPayment, map[string]interface{}{"Code": paymentErr.Code}, paymentErr.Error(), http.StatusInternalServerError)
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("PaymentRefund", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &PaymentRefund{
		Payment: SystemPaymentToGraphqlPayment(payment),
	}, nil
}

// NOTE: Refer to ./schemas/payment.graphqls for details on directives used.
func (r *Resolver) PaymentVoid(ctx context.Context, args struct{ PaymentID UUID }) (*PaymentVoid, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("PaymentVoid", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	payment, appErr := embedCtx.App.Srv().PaymentService().PaymentByID(tx, args.PaymentID.String(), true)
	if appErr != nil {
		return nil, appErr
	}

	var channelID string
	switch {
	case payment.OrderID != nil:
		order, appErr := embedCtx.App.Srv().OrderService().OrderById(*payment.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = order.ChannelID

	case payment.CheckoutID != nil:
		checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
			Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", *payment.CheckoutID),
		})
		if appErr != nil {
			return nil, appErr
		}
		channelID = checkout.ChannelID
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, paymentErr, appErr := embedCtx.App.Srv().PaymentService().Void(tx, *payment, pluginMng, channelID)
	if appErr != nil {
		return nil, appErr
	}

	if paymentErr != nil {
		return nil, model.NewAppError("PaymentVoid", model.ErrPayment, map[string]interface{}{"Code": paymentErr.Code}, paymentErr.Error(), http.StatusInternalServerError)
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("PaymentVoid", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &PaymentVoid{
		Payment: SystemPaymentToGraphqlPayment(payment),
	}, nil
}

// NOTE: Refer to ./schemas/payment.graphqls for details on directives used.
func (r *Resolver) PaymentInitialize(ctx context.Context, args struct {
	ChannelID   UUID
	Gateway     string
	PaymentData model.StringInterface
}) (*PaymentInitialize, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": args.ChannelID},
	})
	if appErr != nil {
		return nil, appErr
	}

	if !channel.IsActive {
		return nil, model.NewAppError("PaymentInitialize", "app.channel.channel_not_active", nil, fmt.Sprintf("Channel with id=%s is inactive", args.ChannelID), http.StatusNotAcceptable)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	response := pluginMng.InitializePayment(args.Gateway, args.PaymentData, channel.Id)

	data, err := json.Marshal(response.Data)
	if err != nil {
		return nil, model.NewAppError("PaymentInitialize", model.ErrorMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	var mapData = JSONString{}
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		return nil, model.NewAppError("PaymentInitialize", model.ErrorUnMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &PaymentInitialize{
		InitializedPayment: &PaymentInitialized{
			Gateway: args.Gateway,
			Name:    response.Name,
			Data:    mapData,
		},
	}, nil
}

// NOTE: Refer to ./schemas/payment.graphqls for details on directives used.
func (r *Resolver) Payment(ctx context.Context, args struct{ Id UUID }) (*Payment, error) {
	payment, err := PaymentByIdLoader.Load(ctx, args.Id.String())()
	if err != nil {
		return nil, err
	}

	return SystemPaymentToGraphqlPayment(payment), nil
}

// NOTE: Refer to ./schemas/payment.graphqls for details on directives used.
func (r *Resolver) Payments(ctx context.Context, args struct {
	Filter *PaymentFilterInput
	GraphqlParams
}) (*PaymentCountableConnection, error) {
	paginValues, appErr := args.GraphqlParams.Parse("Payments")
	if appErr != nil {
		return nil, appErr
	}

	paymentFilterOpts := &model.PaymentFilterOption{
		PaginationValues: *paginValues,
		CountTotal:       true,
	}

	if paymentFilterOpts.PaginationValues.OrderBy == "" {
		// default ordering by gateway and createAt
		ordering := args.GraphqlParams.orderDirection()
		paymentFilterOpts.PaginationValues.OrderBy = fmt.Sprintf("%[1]s.GateWay %[2]s, %[1]s.CreateAt %[2]s", model.PaymentTableName, ordering)
	}

	if args.Filter != nil && len(args.Filter.Checkouts) > 0 {
		paymentFilterOpts.Conditions = squirrel.Eq{model.PaymentTableName + ".CheckoutID": args.Filter.Checkouts}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalCount, payments, appErr := embedCtx.App.Srv().PaymentService().PaymentsByOption(paymentFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	hasNextPage, hasPrevPage := args.GraphqlParams.checkNextPageAndPreviousPage(len(payments))
	keyFunc := func(p *model.Payment) []any {
		return []any{
			model.PaymentTableName + ".GateWay", p.GateWay,
			model.PaymentTableName + ".CreateAt", p.CreateAt,
		}
	}
	connection := constructCountableConnection(payments, totalCount, hasNextPage, hasPrevPage, keyFunc, SystemPaymentToGraphqlPayment)
	return (*PaymentCountableConnection)(unsafe.Pointer(connection)), nil
}
