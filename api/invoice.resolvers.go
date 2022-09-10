package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
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
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) InvoiceDelete(ctx context.Context, args struct{ Id string }) (*InvoiceDelete, error) {
	panic(fmt.Errorf("not implemented"))
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
