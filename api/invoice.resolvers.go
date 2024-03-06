package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/invoice.graphqls for details on directives used.
func (r *Resolver) InvoiceRequest(ctx context.Context, args struct {
	Number  string
	OrderID string
}) (*InvoiceRequest, error) {
	// validate params
	if !model_helper.IsValidId(args.OrderID) {
		return nil, model_helper.NewAppError("InvoiceRequest", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderID"}, "invalid id provided", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// clean order
	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.OrderID)
	if appErr != nil {
		return nil, appErr
	}
	if appErr := commonCleanOrder(order); appErr != nil {
		return nil, appErr
	}

	shallowInvoice := &model.Invoice{
		OrderID: &order.Id,
	}
	shallowInvoice.Number = args.Number

	shallowInvoice, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoice(shallowInvoice)
	if appErr != nil {
		return nil, appErr
	}

	var invoice = &model.Invoice{}

	// TODO: try making InvoiceRequest return *model.Invoice directly
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	invoiceIface, appErr := pluginMng.InvoiceRequest(*order, *shallowInvoice, args.Number)
	if appErr != nil {
		return nil, appErr
	}
	if invoiceIface != nil {
		switch t := invoiceIface.(type) {
		case *model.Invoice:
			invoice = t
		case model.Invoice:
			invoice = &t
		}
	}

	newOrderEventOptions := &model.OrderEventOption{
		OrderID:    order.Id,
		UserID:     &embedCtx.AppContext.Session().UserId,
		Parameters: model.StringInterface{},
	}
	if invoice.Status == model.JobStatusSuccess {
		newOrderEventOptions.Parameters["invoice_number"] = invoice.Number
	}

	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(nil, newOrderEventOptions)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(&model.InvoiceEventCreationOptions{
		UserID:     &embedCtx.AppContext.Session().UserId,
		OrderID:    &order.Id,
		Parameters: model.StringMAP{"number": args.Number},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceRequest{
		Invoice: SystemInvoiceToGraphqlInvoice(invoice),
		Order:   SystemOrderToGraphqlOrder(order),
	}, nil
}

// NOTE: Refer to ./schemas/invoice.graphqls for details on directives used.
func (r *Resolver) InvoiceRequestDelete(ctx context.Context, args struct{ Id string }) (*InvoiceRequestDelete, error) {
	// validate params
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("InvoiceRequestDelete", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid id provided", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	invoice, appErr := embedCtx.App.Srv().InvoiceService().GetInvoiceByOptions(&model.InvoiceFilterOptions{
		Conditions: squirrel.Eq{model.InvoiceTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	invoice.Status = model.JobStatusPending

	updatedInvoice, appErr := embedCtx.App.Srv().InvoiceService().UpsertInvoice(invoice)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.InvoiceDelete(*invoice)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(&model.InvoiceEventCreationOptions{
		UserID:    &embedCtx.AppContext.Session().UserId,
		InvoiceID: &updatedInvoice.Id,
		OrderID:   updatedInvoice.OrderID,
		Type:      model.INVOICE_EVENT_TYPE_REQUESTED_DELETION,
	})
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceRequestDelete{
		Invoice: SystemInvoiceToGraphqlInvoice(updatedInvoice),
	}, nil
}

// commonCleanOrder checks:
//
// If order is draft or unconfirmed, return error.
//
// If order has no BillingAddressID, return error.
func commonCleanOrder(order *model.Order) *model_helper.AppError {
	if order.IsDraft() || order.IsUnconfirmed() {
		return model_helper.NewAppError("commonCleanOrder", "app.order.order_draft_unconfirmed.app_error", nil, "cannot create an invoice for draft or unconfirmed order", http.StatusNotAcceptable)
	}
	if order.BillingAddressID == nil {
		return model_helper.NewAppError("commonCleanOrder", "app.order.order_billing_address_not_set.app_error", nil, "cannot create an invoice for order without biling address", http.StatusNotAcceptable)
	}
	return nil
}

// NOTE: Refer to ./schemas/invoice.graphqls for details on directives used.
func (r *Resolver) InvoiceCreate(ctx context.Context, args struct {
	Input   InvoiceCreateInput
	OrderID string
}) (*InvoiceCreate, error) {
	// validate args
	if !model_helper.IsValidId(args.OrderID) {
		return nil, model_helper.NewAppError("InvoiceCreate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "order id"}, "please provide valid order id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// try finding order in db
	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	// clean order
	if appErr := commonCleanOrder(order); appErr != nil {
		return nil, appErr
	}

	// clean input
	if args.Input.Number == "" || args.Input.URL == "" {
		return nil, model_helper.NewAppError("InvoiceCreate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "url or number"}, "url and number input cannot be empty", http.StatusBadRequest)
	}

	newInvoice := &model.Invoice{
		Number:      args.Input.Number,
		ExternalUrl: args.Input.URL,
		OrderID:     &order.Id,
		Status:      model.JobStatusSuccess,
	}
	savedInvoice, appErr := embedCtx.App.Srv().InvoiceService().UpsertInvoice(newInvoice)
	if appErr != nil {
		return nil, appErr
	}

	// upsert invoice event
	_, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(&model.InvoiceEventCreationOptions{
		UserID:    &embedCtx.AppContext.Session().UserId,
		InvoiceID: &savedInvoice.Id,
		Type:      model.INVOICE_EVENT_TYPE_CREATED,
		Parameters: model.StringMAP{
			"number": args.Input.Number,
			"url":    args.Input.URL,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &embedCtx.AppContext.Session().UserId,
		Type:    model.ORDER_EVENT_TYPE_INVOICE_GENERATED,
		Parameters: model.StringInterface{
			"invoice_number": args.Input.Number,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceCreate{
		Invoice: SystemInvoiceToGraphqlInvoice(savedInvoice),
	}, nil
}

// NOTE: Refer to ./schemas/invoice.graphqls for details on directives used.
func (r *Resolver) InvoiceDelete(ctx context.Context, args struct{ Id string }) (*InvoiceDelete, error) {
	// validate params
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("InvoiceDelete", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	err := embedCtx.App.Srv().Store.Invoice().Delete(nil, args.Id)
	if err != nil {
		return nil, model_helper.NewAppError("InvoiceDelete", "app.invoice.delete_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr := embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(&model.InvoiceEventCreationOptions{
		Type:   model.INVOICE_EVENT_TYPE_DELETED,
		UserID: &embedCtx.AppContext.Session().UserId,
		Parameters: model.StringMAP{
			"invoice_id": args.Id,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceDelete{
		Invoice: &Invoice{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/invoice.graphqls for details on directives used.
func (r *Resolver) InvoiceUpdate(ctx context.Context, args struct {
	Id    string
	Input UpdateInvoiceInput
}) (*InvoiceUpdate, error) {

	// validate params
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("InvoiceUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	invoice, appErr := embedCtx.App.Srv().InvoiceService().GetInvoiceByOptions(&model.InvoiceFilterOptions{
		Conditions: squirrel.Eq{model.InvoiceTableName + ".Id": args.Id},
		Limit:      1,
	})
	if appErr != nil {
		return nil, appErr
	}

	// clean input
	var number = invoice.Number
	var anUrl = invoice.ExternalUrl
	if numInput := args.Input.Number; numInput != nil && *numInput != "" {
		number = *numInput
	}
	if url := args.Input.URL; url != nil && *url != "" {
		anUrl = *url
	}

	if number == "" || anUrl == "" {
		return nil, model_helper.NewAppError("InvoiceUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "number/url"}, "number and url must be set after update operation", http.StatusNotAcceptable)
	}

	invoice.Number = number
	invoice.ExternalUrl = anUrl
	invoice.Status = model.JobStatusSuccess
	updatedInvoice, appErr := embedCtx.App.Srv().InvoiceService().UpsertInvoice(invoice)
	if appErr != nil {
		return nil, appErr
	}

	orderEventOptions := &model.OrderEventOption{
		UserID: &embedCtx.AppContext.Session().UserId,
		Type:   model.ORDER_EVENT_TYPE_INVOICE_UPDATED,
		Parameters: model.StringInterface{
			"invoice_number": updatedInvoice.Number,
			"url":            updatedInvoice.ExternalUrl,
			"status":         updatedInvoice.Status,
		},
	}
	if oID := updatedInvoice.OrderID; oID != nil {
		orderEventOptions.OrderID = *oID
	}
	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(nil, orderEventOptions)
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceUpdate{
		Invoice: SystemInvoiceToGraphqlInvoice(updatedInvoice),
	}, nil
}

// NOTE: Refer to ./schemas/invoice.graphqls for details on directives used.
func (r *Resolver) InvoiceSendNotification(ctx context.Context, args struct{ Id string }) (*InvoiceSendNotification, error) {
	// validate params
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("InvoiceSendNotification", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	invoice, appErr := embedCtx.App.Srv().InvoiceService().GetInvoiceByOptions(&model.InvoiceFilterOptions{
		Conditions:         squirrel.Eq{model.InvoiceTableName + ".Id": args.Id},
		SelectRelatedOrder: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	// clean instance
	switch {
	case invoice.Status != model.JobStatusSuccess:
		return nil, model_helper.NewAppError("InvoiceSendNotification", "app.invoice.status_not_success.app_error", nil, "Provided invoice is not ready to be sent", http.StatusNotAcceptable)
	case invoice.ExternalUrl == "":
		return nil, model_helper.NewAppError("InvoiceSendNotification", "app.invoice.url_not_set.app_error", nil, "Provided invoice needs to have an URL", http.StatusNotAcceptable)
	case invoice.Number == "":
		return nil, model_helper.NewAppError("InvoiceSendNotification", "app.invoice.number_not_set.app_error", nil, "Provided invoice needs to have an invoice number", http.StatusNotAcceptable)
	case invoice.OrderID == nil:
		return nil, model_helper.NewAppError("InvoiceSendNotification", "app.invoice.order_not_set.app_error", nil, "provided invoice needs an associated order", http.StatusNotAcceptable)
	}

	orderEmail, appErr := embedCtx.App.Srv().OrderService().CustomerEmail(invoice.GetOrder())
	if appErr != nil {
		return nil, appErr
	}
	if orderEmail == "" {
		return nil, model_helper.NewAppError("InvoiceSendNotification", "app.order.order_email_not_set.app_error", nil, "provided invoice order needs an email address", http.StatusNotAcceptable)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	appErr = embedCtx.App.Srv().InvoiceService().SendInvoice(*invoice, &model.User{Id: embedCtx.AppContext.Session().UserId}, nil, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceSendNotification{
		Invoice: SystemInvoiceToGraphqlInvoice(invoice),
	}, nil
}
