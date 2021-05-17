package graph

import (
	"context"
	"time"

	"github.com/sitename/sitename/graph/model"
)

func (r *queryResolver) payment(ctx context.Context, id string) (*model.Payment, error) {
	now := time.Now()
	return &model.Payment{
		ID:                     id,
		Gateway:                "ahihi",
		IsActive:               true,
		Created:                now.String(),
		Modified:               now.String(),
		Token:                  "ThisistheToken",
		Checkout:               nil,
		Order:                  nil,
		PaymentMethodType:      "hihi",
		CustomerIPAddress:      nil,
		ChargeStatus:           model.PaymentChargeStatusEnumNotCharged,
		Actions:                []*model.OrderAction{},
		Total:                  nil,
		CapturedAmount:         nil,
		Transactions:           []*model.Transaction{},
		AvailableCaptureAmount: nil,
		AvailableRefundAmount:  nil,
		CreditCard:             nil,
	}, nil
}
