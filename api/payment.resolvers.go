package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/api/gqlmodel"
	"github.com/sitename/sitename/model"
)

func (r *Resolver) PaymentCapture(ctx context.Context, args struct {
	Amount    *decimal.Decimal
	PaymentID string
}) (*gqlmodel.PaymentCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PaymentRefund(ctx context.Context, args struct {
	Amount    *decimal.Decimal
	PaymentID string
}) (*gqlmodel.PaymentRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PaymentVoid(ctx context.Context, args struct{ PaymentID string }) (*gqlmodel.PaymentVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PaymentInitialize(ctx context.Context, args struct {
	Channel     *string
	Gateway     string
	PaymentData model.StringInterface
}) (*gqlmodel.PaymentInitialize, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Payment(ctx context.Context, args struct{ Id string }) (*gqlmodel.Payment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Payments(ctx context.Context, args struct {
	Filter *gqlmodel.PaymentFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.PaymentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}