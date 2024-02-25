package model_helper

import (
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
)

const (
	MIN_CC_EXP_MONTH = 1
	MAX_CC_EXP_MONTH = 12
	MIN_CC_EXP_YEAR  = 1000

	// Payment Gateways
	GATE_WAY_MANUAL = "manual"
)

func PaymentPreSave(p *model.Payment) {

}

func PaymentGetChargeAmount(p model.Payment) decimal.Decimal {
	return p.Total.Sub(p.CapturedAmount)
}

func PaymentIsNotCharged(p model.Payment) bool {
	return p.ChargeStatus == model.PaymentChargeStatusNotCharged
}

func PaymentCanAuthorize(p model.Payment) bool {
	return model_types.PrimitiveIsNotNilAndEqual(p.IsActive.Bool, true) && PaymentIsNotCharged(p)
}

var PaymentCanCapture = PaymentCanAuthorize
var PaymentCanConfirm = PaymentCanAuthorize

var CanRefundPaymentChargeStatusMap = map[model.PaymentChargeStatus]bool{
	model.PaymentChargeStatusPartiallyCharged:  true,
	model.PaymentChargeStatusFullyCharged:      true,
	model.PaymentChargeStatusPartiallyRefunded: true,
}

func PaymentCanRefund(p model.Payment) bool {
	return CanRefundPaymentChargeStatusMap[p.ChargeStatus] && model_types.PrimitiveIsNotNilAndEqual(p.IsActive.Bool, true)
}

func PaymentIsManual(p model.Payment) bool {
	return p.Gateway == GATE_WAY_MANUAL
}

func PaymentGetTotalPrice(p model.Payment) goprices.Money {
	return goprices.Money{
		Amount:   p.Total,
		Currency: p.Currency.String(),
	}
}

func PaymentGetCapturedAmount(p model.Payment) goprices.Money {
	return goprices.Money{
		Amount:   p.CapturedAmount,
		Currency: p.Currency.String(),
	}
}
