package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
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
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceRequestDelete(ctx context.Context, args struct{ Id string }) (*InvoiceRequestDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceCreate(ctx context.Context, args struct {
	Input   InvoiceCreateInput
	OrderID string
}) (*InvoiceCreate, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("InvoiceCreate", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

	// try finding order in db
	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	// clean order
	if order.IsDraft() || order.IsUnconfirmed() {
		return nil, model.NewAppError("InvoiceCreate", "app.order.order_draft_unconfirmed.app_error", nil, "cannot create an invoice for draft or unconfirmed order", http.StatusNotAcceptable)
	}
	if order.BillingAddressID == nil {
		return nil, model.NewAppError("InvoiceCreate", "app.order.order_billing_address_not_set.app_error", nil, "cannot create an invoice for order without biling address", http.StatusNotAcceptable)
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
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("InvoiceDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("InvoiceDelete", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
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
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceSendNotification(ctx context.Context, args struct{ Id string }) (*InvoiceSendNotification, error) {
	panic(fmt.Errorf("not implemented"))
}
