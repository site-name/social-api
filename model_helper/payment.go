package model_helper

import (
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	MIN_CC_EXP_MONTH = 1
	MAX_CC_EXP_MONTH = 12
	MIN_CC_EXP_YEAR  = 1000

	// Payment Gateways
	GATE_WAY_MANUAL = "manual"
)

func PaymentPreSave(p *model.Payment) {
	if p.ID == "" {
		p.ID = NewId()
	}
	if p.CreatedAt == 0 {
		p.CreatedAt = GetMillis()
	}
	p.UpdatedAt = p.CreatedAt
	paymentCommonPre(p)
}

func PaymentPreUpdate(p *model.Payment) {
	p.UpdatedAt = GetMillis()
	paymentCommonPre(p)
}

func paymentCommonPre(p *model.Payment) {
	p.BillingEmail = NormalizeEmail(p.BillingEmail)
	p.BillingFirstName = SanitizeUnicode(CleanNamePart(p.BillingFirstName))
	p.BillingLastName = SanitizeUnicode(CleanNamePart(p.BillingLastName))
	if p.Total.LessThanOrEqual(decimal.Zero) {
		p.Total = decimal.NewFromInt(0)
	}
	if p.CapturedAmount.LessThanOrEqual(decimal.Zero) {
		p.CapturedAmount = decimal.NewFromInt(0)
	}
	if p.ChargeStatus.IsValid() != nil {
		p.ChargeStatus = model.PaymentChargeStatusNotCharged
	}
	if p.Currency.IsValid() != nil {
		p.Currency = DEFAULT_CURRENCY
	}
	if p.StorePaymentMethod.IsValid() != nil {
		p.StorePaymentMethod = model.StorePaymentMethodNone
	}
}

type PaymentFilterOptions struct {
	CommonQueryOptions
	TransactionCondition qm.QueryMod // INNER JOIN payment_transactions ON payment_transactions.payment_id = payments.id WHERE ...
}

// PaymentPatch is used to update payments
type PaymentPatch struct {
	CheckoutID   string
	OrderID      string
	BillingEmail string
}

func PaymentGetChargeAmount(p model.Payment) decimal.Decimal {
	return p.Total.Sub(p.CapturedAmount)
}

func PaymentIsNotCharged(p model.Payment) bool {
	return p.ChargeStatus == model.PaymentChargeStatusNotCharged
}

func PaymentCanAuthorize(p model.Payment) bool {
	return p.IsActive && PaymentIsNotCharged(p)
}

var PaymentCanCapture = PaymentCanAuthorize
var PaymentCanConfirm = PaymentCanAuthorize

var CanRefundPaymentChargeStatusMap = map[model.PaymentChargeStatus]bool{
	model.PaymentChargeStatusPartiallyCharged:  true,
	model.PaymentChargeStatusFullyCharged:      true,
	model.PaymentChargeStatusPartiallyRefunded: true,
}

func PaymentCanRefund(p model.Payment) bool {
	return CanRefundPaymentChargeStatusMap[p.ChargeStatus] && p.IsActive
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

func PaymentIsValid(p model.Payment) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.CheckoutID.IsNil() && !IsValidId(*p.CheckoutID.String) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.checkout_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Total.LessThanOrEqual(decimal.Zero) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.total.app_error", nil, "", http.StatusBadRequest)
	}
	if p.CapturedAmount.LessThan(decimal.Zero) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.captured_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidEmail(p.BillingEmail) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.billing_email.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidNamePart(p.BillingFirstName) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.billing_first_name.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidNamePart(p.BillingLastName) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.billing_last_name.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Currency.IsValid() != nil {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if p.StorePaymentMethod.IsValid() != nil {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.store_payment_method.app_error", nil, "", http.StatusBadRequest)
	}
	if p.ChargeStatus.IsValid() != nil {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.charge_status.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.CCExpMonth.IsNil() && *p.CCExpMonth.Int < MIN_CC_EXP_MONTH || *p.CCExpMonth.Int > MAX_CC_EXP_MONTH {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.cc_exp_month.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.CCExpYear.IsNil() && *p.CCExpYear.Int < MIN_CC_EXP_YEAR {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.cc_exp_year.app_error", nil, "", http.StatusBadRequest)
	}
	if p.CreatedAt <= 0 {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if p.UpdatedAt <= 0 {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.updated_at.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func PaymentTransactionPreSave(pt *model.PaymentTransaction) {
	if pt.ID == "" {
		pt.ID = NewId()
	}
	if pt.CreatedAt == 0 {
		pt.CreatedAt = GetMillis()
	}
	PaymentTransactionCommonPre(pt)
}

func PaymentTransactionCommonPre(pt *model.PaymentTransaction) {
	if !pt.Error.IsNil() {
		*pt.Error.String = SanitizeUnicode(*pt.Error.String)
	}
	if pt.Amount.LessThanOrEqual(decimal.Zero) {
		pt.Amount = decimal.NewFromInt(0)
	}
	if pt.ActionRequiredData == nil {
		pt.ActionRequiredData = make(model_types.JSONString)
	}
	if pt.GatewayResponse == nil {
		pt.GatewayResponse = make(model_types.JSONString)
	}
}

type PaymentTransactionFilterOpts struct {
	CommonQueryOptions
}

func PaymentTransactionIsValid(pt model.PaymentTransaction) *AppError {
	if !IsValidId(pt.ID) {
		return NewAppError("PaymentTransaction.IsValid", "model.payment_transaction.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(pt.PaymentID) {
		return NewAppError("PaymentTransaction.IsValid", "model.payment_transaction.is_valid.payment_id.app_error", nil, "", http.StatusBadRequest)
	}
	if pt.Amount.LessThanOrEqual(decimal.Zero) {
		return NewAppError("PaymentTransaction.IsValid", "model.payment_transaction.is_valid.amount.app_error", nil, "", http.StatusBadRequest)
	}
	if pt.Currency.IsValid() != nil {
		return NewAppError("PaymentTransaction.IsValid", "model.payment_transaction.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if pt.Kind.IsValid() != nil {
		return NewAppError("PaymentTransaction.IsValid", "model.payment_transaction.is_valid.kind.app_error", nil, "", http.StatusBadRequest)
	}
	if pt.CreatedAt <= 0 {
		return NewAppError("PaymentTransaction.IsValid", "model.payment_transaction.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func PaymentTransactionGetAmount(pt model.PaymentTransaction) goprices.Money {
	return goprices.Money{
		Amount:   pt.Amount,
		Currency: pt.Currency.String(),
	}
}
