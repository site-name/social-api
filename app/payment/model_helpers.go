package payment

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/util"
)

// GetLastpayment compares all payments's CreatAt properties, then returns the most recent payment
func (a *AppPayment) GetLastpayment(payments []*payment.Payment) *payment.Payment {
	if len(payments) == 0 {
		return nil
	}

	var latestPayment *payment.Payment
	for _, payMent := range payments {
		if latestPayment == nil || (latestPayment != nil && payMent.CreateAt > latestPayment.CreateAt) {
			latestPayment = payMent
		}
	}

	return latestPayment
}

func (a *AppPayment) GetTotalAuthorized(payments []*payment.Payment, fallbackCurrency string) (*goprices.Money, *model.AppError) {
	zeroMoney, err := util.ZeroMoney(fallbackCurrency)
	if err != nil {
		return nil, model.NewAppError("GetTotalAuthorized", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "fallbackCurrency"}, err.Error(), http.StatusBadRequest)
	}

	lastPayment := a.GetLastpayment(payments)
	if lastPayment != nil {
		if *lastPayment.IsActive {
			paymentAuthorizedAmount, appErr := a.PaymentGetAuthorizedAmount(lastPayment)
			if appErr != nil {
				return nil, appErr
			}

			return paymentAuthorizedAmount, nil
		}
	}

	return zeroMoney, nil
}

// GetSubTotal adds up all Total prices of given order lines
func (a *AppPayment) GetSubTotal(orderLines []*order.OrderLine, fallbackCurrency string) (*goprices.TaxedMoney, *model.AppError) {
	total, err := util.ZeroTaxedMoney(fallbackCurrency)
	if err != nil {
		return nil, model.NewAppError("GetSubTotal", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "fallbackCurrency"}, err.Error(), http.StatusBadRequest)
	}

	for _, line := range orderLines {
		line.PopulateNonDbFields()

		total, _ = total.Add(line.TotalPrice)
	}

	return total, nil
}