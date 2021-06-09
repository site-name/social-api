package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) PaymentCapture(ctx context.Context, amount *string, paymentID string) (*PaymentCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentRefund(ctx context.Context, amount *string, paymentID string) (*PaymentRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentVoid(ctx context.Context, paymentID string) (*PaymentVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentInitialize(ctx context.Context, channel *string, gateway string, paymentData *string) (*PaymentInitialize, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Payment(ctx context.Context, id string) (*Payment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Payments(ctx context.Context, filter *PaymentFilterInput, before *string, after *string, first *int, last *int) (*PaymentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
