package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Payment struct {
	ID                     string                  `json:"id"`
	Gateway                string                  `json:"gateway"`
	IsActive               bool                    `json:"isActive"`
	Created                DateTime                `json:"created"`
	Modified               DateTime                `json:"modified"`
	Token                  string                  `json:"token"`
	Checkout               *Checkout               `json:"checkout"`
	Order                  *Order                  `json:"order"`
	PaymentMethodType      string                  `json:"paymentMethodType"`
	CustomerIPAddress      *string                 `json:"customerIpAddress"`
	PrivateMetadata        []*MetadataItem         `json:"privateMetadata"`
	Metadata               []*MetadataItem         `json:"metadata"`
	ChargeStatus           PaymentChargeStatusEnum `json:"chargeStatus"`
	Actions                []*OrderAction          `json:"actions"`
	Total                  *Money                  `json:"total"`
	CapturedAmount         *Money                  `json:"capturedAmount"`
	Transactions           []*Transaction          `json:"transactions"`
	AvailableCaptureAmount *Money                  `json:"availableCaptureAmount"`
	AvailableRefundAmount  *Money                  `json:"availableRefundAmount"`
	CreditCard             *CreditCard             `json:"creditCard"`
}

func SystemPaymentToGraphqlPayment(p *model.Payment) *Payment {
	if p == nil {
		return nil
	}

	res := &Payment{}
	panic("not implemented")
	return res
}

func paymentsByOrderIdLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.Payment] {
	var (
		res        = make([]*dataloader.Result[[]*model.Payment], len(orderIDs))
		payments   []*model.Payment
		appErr     *model.AppError
		paymentMap = map[string][]*model.Payment{} // keys are order ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	payments, appErr = embedCtx.App.Srv().
		PaymentService().
		PaymentsByOption(&model.PaymentFilterOption{
			OrderID: squirrel.Eq{store.PaymentTableName + ".OrderID": orderIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, payment := range payments {
		if payment.OrderID == nil {
			continue
		}
		paymentMap[*payment.OrderID] = append(paymentMap[*payment.OrderID], payment)
	}

	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.Payment]{Data: paymentMap[id]}
	}
	return res

errorLabel:
	for idx := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.Payment]{Error: err}
	}
	return res
}
