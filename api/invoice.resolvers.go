package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) InvoiceRequest(ctx context.Context, args struct {
	Number  *string
	OrderID string
}) (*gqlmodel.InvoiceRequest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceRequestDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.InvoiceRequestDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceCreate(ctx context.Context, args struct {
	Input   gqlmodel.InvoiceCreateInput
	OrderID string
}) (*gqlmodel.InvoiceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.InvoiceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.UpdateInvoiceInput
}) (*gqlmodel.InvoiceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceSendNotification(ctx context.Context, args struct{ Id string }) (*gqlmodel.InvoiceSendNotification, error) {
	panic(fmt.Errorf("not implemented"))
}
