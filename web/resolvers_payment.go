package web

import (
	"context"
	"time"

	"github.com/sitename/sitename/web/model"
)

func (r *queryResolver) payment(ctx context.Context, id string) (*model.Payment, error) {

	// _, ok := ctx.Value(ApiContextKey).(*Context)
	// fmt.Println("It is", ok)

	now := time.Now()
	return &model.Payment{
		ID:                     id,
		Gateway:                "ahihi",
		IsActive:               true,
		Created:                now,
		Modified:               now,
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
