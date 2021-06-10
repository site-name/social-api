package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) InvoiceRequest(ctx context.Context, number *string, orderID string) (*InvoiceRequest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceRequestDelete(ctx context.Context, id string) (*InvoiceRequestDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceCreate(ctx context.Context, input InvoiceCreateInput, orderID string) (*InvoiceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceDelete(ctx context.Context, id string) (*InvoiceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceUpdate(ctx context.Context, id string, input UpdateInvoiceInput) (*InvoiceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceSendNotification(ctx context.Context, id string) (*InvoiceSendNotification, error) {
	panic(fmt.Errorf("not implemented"))
}
