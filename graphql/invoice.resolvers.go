package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) InvoiceRequest(ctx context.Context, number *string, orderID string) (*gqlmodel.InvoiceRequest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceRequestDelete(ctx context.Context, id string) (*gqlmodel.InvoiceRequestDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceCreate(ctx context.Context, input gqlmodel.InvoiceCreateInput, orderID string) (*gqlmodel.InvoiceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceDelete(ctx context.Context, id string) (*gqlmodel.InvoiceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceUpdate(ctx context.Context, id string, input gqlmodel.UpdateInvoiceInput) (*gqlmodel.InvoiceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceSendNotification(ctx context.Context, id string) (*gqlmodel.InvoiceSendNotification, error) {
	panic(fmt.Errorf("not implemented"))
}
