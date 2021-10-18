package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) PaymentCapture(ctx context.Context, amount *decimal.Decimal, paymentID string) (*gqlmodel.PaymentCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentRefund(ctx context.Context, amount *decimal.Decimal, paymentID string) (*gqlmodel.PaymentRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentVoid(ctx context.Context, paymentID string) (*gqlmodel.PaymentVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentInitialize(ctx context.Context, channel *string, gateway string, paymentData *string) (*gqlmodel.PaymentInitialize, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Payment(ctx context.Context, id string) (*gqlmodel.Payment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Payments(ctx context.Context, filter *gqlmodel.PaymentFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.PaymentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
