package payment

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

// GetLastpayment compares all payments's CreatAt properties, then returns the most recent payment
func (a *ServicePayment) GetLastpayment(payments []*model.Payment) *model.Payment {
	if len(payments) == 0 {
		return nil
	}

	if len(payments) == 1 {
		return payments[0]
	}

	res := payments[0]
	for _, pm := range payments[1:] {
		if pm != nil && pm.CreateAt > res.CreateAt {
			res = pm
		}
	}

	return res
}

func (a *ServicePayment) GetTotalAuthorized(payments []*model.Payment, fallbackCurrency string) (*goprices.Money, *model_helper.AppError) {
	zeroMoney, err := util.ZeroMoney(fallbackCurrency)
	if err != nil {
		return nil, model_helper.NewAppError("GetTotalAuthorized", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "fallbackCurrency"}, err.Error(), http.StatusBadRequest)
	}

	lastPayment := a.GetLastpayment(payments)
	if lastPayment != nil && *lastPayment.IsActive {
		paymentAuthorizedAmount, appErr := a.PaymentGetAuthorizedAmount(lastPayment)
		if appErr != nil {
			return nil, appErr
		}

		return paymentAuthorizedAmount, nil
	}

	return zeroMoney, nil
}

// GetSubTotal adds up all Total prices of given order lines
func (a *ServicePayment) GetSubTotal(orderLines []*model.OrderLine, fallbackCurrency string) (*goprices.TaxedMoney, *model_helper.AppError) {
	total, err := util.ZeroTaxedMoney(fallbackCurrency)
	if err != nil {
		return nil, model_helper.NewAppError("GetSubTotal", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "fallbackCurrency"}, err.Error(), http.StatusBadRequest)
	}

	for _, line := range orderLines {
		line.PopulateNonDbFields()

		total, err = total.Add(line.TotalPrice)
		if err != nil {
			return nil, model_helper.NewAppError("GetSubTotal", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "fallbackCurrency"}, err.Error(), http.StatusBadRequest)
		}
	}

	return total, nil
}
