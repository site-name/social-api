package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) InvoiceRequest(ctx context.Context, args struct {
	Number  *string
	OrderID string
}) (*InvoiceRequest, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateOrder)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	if !model.IsValidId(args.OrderID) {
		return nil, model.NewAppError("InvoiceRequest", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderID"}, "invalid id provided", http.StatusBadRequest)
	}

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
	if args.Number != nil && *args.Number != "" {
		shallowInvoice.Number = *args.Number
	}

	shallowInvoice, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoice(shallowInvoice)
	if appErr != nil {
		return nil, appErr
	}

	panic("not implemented") // todo: complete plugin
	var invoice *model.Invoice

	newOrderEventOptions := &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &embedCtx.AppContext.Session().UserId,
	}
	if invoice.Status == model.JobStatusSuccess {
		newOrderEventOptions.Parameters["invoice_number"] = invoice.Number
	}

	_, appErr = embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(nil, newOrderEventOptions)
	if appErr != nil {
		return nil, appErr
	}

	invoiceEventOpts := &model.InvoiceEventOption{
		UserID:  &embedCtx.AppContext.Session().UserId,
		OrderID: &order.Id,
	}
	if args.Number != nil {
		invoiceEventOpts.Parameters["number"] = *args.Number
	}
	_, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(invoiceEventOpts)
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceRequest{
		Invoice: SystemInvoiceToGraphqlInvoice(invoice),
		Order:   SystemOrderToGraphqlOrder(order),
	}, nil
}

func (r *Resolver) InvoiceRequestDelete(ctx context.Context, args struct{ Id string }) (*InvoiceRequestDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("InvoiceRequestDelete", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("InvoiceRequestDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid id provided", http.StatusBadRequest)
	}
	invoices, appErr := embedCtx.App.Srv().InvoiceService().FilterInvoicesByOptions(&model.InvoiceFilterOptions{
		Id: squirrel.Eq{store.InvoiceTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(invoices) == 0 {
		return nil, model.NewAppError("InvoiceRequestDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid id provided", http.StatusBadRequest)
	}

	invoice := invoices[0]
	invoice.Status = model.JobStatusPending

	updatedInvoice, appErr := embedCtx.App.Srv().InvoiceService().UpsertInvoice(invoice)
	if appErr != nil {
		return nil, appErr
	}

	panic("not implemented") // complete plugin module

	_, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(&model.InvoiceEventOption{
		UserID:    &embedCtx.AppContext.Session().UserId,
		InvoiceID: &updatedInvoice.Id,
		OrderID:   updatedInvoice.OrderID,
		Type:      model.REQUESTED_DELETION,
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
func commonCleanOrder(order *model.Order) *model.AppError {
	if order.IsDraft() || order.IsUnconfirmed() {
		return model.NewAppError("commonCleanOrder", "app.order.order_draft_unconfirmed.app_error", nil, "cannot create an invoice for draft or unconfirmed order", http.StatusNotAcceptable)
	}
	if order.BillingAddressID == nil {
		return model.NewAppError("commonCleanOrder", "app.order.order_billing_address_not_set.app_error", nil, "cannot create an invoice for order without biling address", http.StatusNotAcceptable)
	}
	return nil
}

func (r *Resolver) InvoiceCreate(ctx context.Context, args struct {
	Input   InvoiceCreateInput
	OrderID string
}) (*InvoiceCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("InvoiceCreate", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

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
		return nil, model.NewAppError("InvoiceCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "url or number"}, "url and number input cannot be empty", http.StatusBadRequest)
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
	_, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(&model.InvoiceEventOption{
		UserID:    &embedCtx.AppContext.Session().UserId,
		InvoiceID: &savedInvoice.Id,
		Type:      model.CREATED,
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
		Type:    model.INVOICE_GENERATED,
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

func (r *Resolver) InvoiceDelete(ctx context.Context, args struct{ Id string }) (*InvoiceDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("InvoiceDelete", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("InvoiceDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	invoices, appErr := embedCtx.App.Srv().InvoiceService().FilterInvoicesByOptions(&model.InvoiceFilterOptions{
		Id:    squirrel.Eq{store.InvoiceTableName + ".Id": args.Id},
		Limit: 1,
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(invoices) == 0 {
		return nil, model.NewAppError("InvoiceDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invoice with given id does not exist", http.StatusBadRequest)
	}
	invoice := invoices[0]

	// safely delete invoice
	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("InvoiceDelete", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer r.srv.Store.FinalizeTransaction(transaction)

	err = embedCtx.App.Srv().Store.Invoice().Delete(transaction, []string{invoice.Id})
	if err != nil {
		return nil, model.NewAppError("InvoiceDelete", "app.invoice.delete_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = embedCtx.App.Srv().InvoiceService().UpsertInvoiceEvent(&model.InvoiceEventOption{
		Type:   model.DELETED,
		UserID: &embedCtx.AppContext.Session().UserId,
		Parameters: model.StringMAP{
			"invoice_id": invoice.Id,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceDelete{
		Invoice: SystemInvoiceToGraphqlInvoice(invoice),
	}, nil
}

func (r *Resolver) InvoiceUpdate(ctx context.Context, args struct {
	Id    string
	Input UpdateInvoiceInput
}) (*InvoiceUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("InvoiceUpdate", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("InvoiceUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	invoices, appErr := embedCtx.App.Srv().InvoiceService().FilterInvoicesByOptions(&model.InvoiceFilterOptions{
		Id:    squirrel.Eq{store.InvoiceTableName + ".Id": args.Id},
		Limit: 1,
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(invoices) == 0 {
		return nil, model.NewAppError("InvoiceUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invoice with given id does not exist", http.StatusBadRequest)
	}
	invoice := invoices[0]

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
		return nil, model.NewAppError("InvoiceUpdate", "app.invoice.value_not_set.app_error", nil, "number and url must be set after update operation", http.StatusNotAcceptable)
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
		Type:   model.INVOICE_UPDATED,
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

func (r *Resolver) InvoiceSendNotification(ctx context.Context, args struct{ Id string }) (*InvoiceSendNotification, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("InvoiceSendNotification", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("InvoiceSendNotification", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	invoices, appErr := embedCtx.App.Srv().InvoiceService().FilterInvoicesByOptions(&model.InvoiceFilterOptions{
		Id:                 squirrel.Eq{store.InvoiceTableName + ".Id": args.Id},
		Limit:              1,
		SelectRelatedOrder: true,
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(invoices) == 0 {
		return nil, model.NewAppError("InvoiceSendNotification", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invoice with given id does not exist", http.StatusBadRequest)
	}
	invoice := invoices[0]

	// clean instance
	switch {
	case invoice.Status != model.JobStatusSuccess:
		return nil, model.NewAppError("InvoiceSendNotification", "app.invoice.status_not_success.app_error", nil, "Provided invoice is not ready to be sent", http.StatusNotAcceptable)
	case invoice.ExternalUrl == "":
		return nil, model.NewAppError("InvoiceSendNotification", "app.invoice.url_not_set.app_error", nil, "Provided invoice needs to have an URL", http.StatusNotAcceptable)
	case invoice.Number == "":
		return nil, model.NewAppError("InvoiceSendNotification", "app.invoice.number_not_set.app_error", nil, "Provided invoice needs to have an invoice number", http.StatusNotAcceptable)
	case invoice.OrderID == nil:
		return nil, model.NewAppError("InvoiceSendNotification", "app.invoice.order_not_set.app_error", nil, "provided invoice needs an associated order", http.StatusNotAcceptable)
	}

	orderEmail, appErr := embedCtx.App.Srv().OrderService().CustomerEmail(invoice.GetOrder())
	if appErr != nil {
		return nil, appErr
	}
	if orderEmail == "" {
		return nil, model.NewAppError("InvoiceSendNotification", "app.order.order_email_not_set.app_error", nil, "provided invoice order needs an email address", http.StatusNotAcceptable)
	}

	panic("not implemented") // TODO: compete plugin manager first
	appErr = embedCtx.App.Srv().InvoiceService().SendInvoice(*invoice, &model.User{Id: embedCtx.AppContext.Session().UserId}, nil, nil)
	if appErr != nil {
		return nil, appErr
	}

	return &InvoiceSendNotification{
		Invoice: SystemInvoiceToGraphqlInvoice(invoice),
	}, nil
}
