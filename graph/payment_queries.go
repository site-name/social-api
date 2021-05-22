package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/sitename/sitename/graph/model"
	"github.com/sitename/sitename/modules/consts"
	"github.com/sitename/sitename/web"
)

func (r *queryResolver) payment(ctx context.Context, id string) (*model.Payment, error) {

	fmt.Println(ctx.Value(consts.ApiContextKey).(*web.Context))

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
