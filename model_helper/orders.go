package model_helper

import (
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func OrderCommonPre(order *model.Order) {
	if order.Token == "" {
		order.Token = NewId()
	}
	order.CustomerNote = SanitizeUnicode(order.CustomerNote)
	if order.Currency.IsValid() != nil {
		order.Currency = DEFAULT_CURRENCY
	}
	if order.WeightAmount < 0 {
		order.WeightAmount = 0
	}
	if order.Status.IsValid() != nil {
		order.Status = model.OrderStatusUnfulfilled
	}
	if order.Currency.IsValid() != nil {
		order.Currency = DEFAULT_CURRENCY
	}
	if order.LanguageCode.IsValid() != nil {
		order.LanguageCode = DEFAULT_LOCALE
	}
	order.WeightUnit = strings.ToLower(order.WeightUnit)
}

func OrderPreSave(o *model.Order) {
	OrderCommonPre(o)
	if o.ID == "" {
		o.ID = NewId()
	}
	o.CreatedAt = GetMillis()

	// NOTE: those 2 fields are not placed inside commonPre because they are not editable
	if !o.ShippingMethodName.IsNil() {
		o.ShippingMethodName.String = GetPointerOfValue(SanitizeUnicode(*o.ShippingMethodName.String))
	}
	if !o.CollectionPointName.IsNil() {
		o.CollectionPointName.String = GetPointerOfValue(SanitizeUnicode(*o.CollectionPointName.String))
	}
}

func OrderIsValid(o model.Order) *AppError {
	if !IsValidId(o.ID) {
		return NewAppError("OrderIsValid", "model.order.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(o.Token) {
		return NewAppError("OrderIsValid", "model.order.is_valid.token.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Currency.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if o.TotalNetAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.total_net_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.TotalGrossAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.total_gross_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.UndiscountedTotalNetAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.undiscounted_total_net_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.UndiscountedTotalGrossAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.undiscounted_total_gross_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.WeightAmount < 0 {
		return NewAppError("OrderIsValid", "model.order.is_valid.weight_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.ShippingPriceNetAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_price_net_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.ShippingPriceGrossAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_price_gross_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.CreatedAt <= 0 {
		return NewAppError("OrderIsValid", "model.order.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if o.LanguageCode.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.language_code.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Status.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.status.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Currency.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Origin.Valid && o.Origin.Val.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.origin.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.ShippingAddressID.IsNil() && !IsValidId(*o.ShippingAddressID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_address_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.BillingAddressID.IsNil() && !IsValidId(*o.BillingAddressID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.billing_address_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.UserID.IsNil() && !IsValidId(*o.UserID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.VoucherID.IsNil() && !IsValidId(*o.VoucherID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.voucher_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.OriginalID.IsNil() && !IsValidId(*o.OriginalID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.original_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.ShippingMethodID.IsNil() && !IsValidId(*o.ShippingMethodID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_method_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.CollectionPointID.IsNil() && !IsValidId(*o.CollectionPointID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.collection_point_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(o.CheckoutToken) {
		return NewAppError("OrderIsValid", "model.order.is_valid.checkout_token.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(o.ChannelID) {
		return NewAppError("OrderIsValid", "model.order.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}
	_, ok := measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(o.WeightUnit)]
	if !ok {
		return NewAppError("OrderIsValid", "model.order.is_valid.weight_unit.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func OrderGetShippingPrice(order model.Order) goprices.TaxedMoney {
	shippingPriceNet := goprices.Money{
		Amount:   order.ShippingPriceNetAmount,
		Currency: order.Currency.String(),
	}

	shippingPriceGross := goprices.Money{
		Amount:   order.ShippingPriceGrossAmount,
		Currency: order.Currency.String(),
	}

	return goprices.TaxedMoney{
		Net:   shippingPriceNet,
		Gross: shippingPriceGross,
	}
}

func OrderGetTotalPrice(order model.Order) goprices.TaxedMoney {
	totalPriceNet := goprices.Money{
		Amount:   order.TotalNetAmount,
		Currency: order.Currency.String(),
	}

	totalPriceGross := goprices.Money{
		Amount:   order.TotalGrossAmount,
		Currency: order.Currency.String(),
	}

	return goprices.TaxedMoney{
		Net:   totalPriceNet,
		Gross: totalPriceGross,
	}
}

func OrderGetTotalPaidPrice(o model.Order) goprices.Money {
	return goprices.Money{
		Amount:   o.TotalPaidAmount,
		Currency: o.Currency.String(),
	}
}

func OrderGetUnDiscountedTotalPrice(order model.Order) goprices.TaxedMoney {
	totalPriceNet := goprices.Money{
		Amount:   order.UndiscountedTotalNetAmount,
		Currency: order.Currency.String(),
	}

	totalPriceGross := goprices.Money{
		Amount:   order.UndiscountedTotalGrossAmount,
		Currency: order.Currency.String(),
	}

	return goprices.TaxedMoney{
		Net:   totalPriceNet,
		Gross: totalPriceGross,
	}
}

func OrderIsDraft(o model.Order) bool {
	return o.Status == model.OrderStatusDraft
}

func OrderIsUnConfirmed(o model.Order) bool {
	return o.Status == model.OrderStatusUnconfirmed
}

func OrderIsOpen(o model.Order) bool {
	return o.Status == model.OrderStatusDraft || o.Status == model.OrderStatusPartiallyFulfilled
}

func OrderIsFullyPaid(o model.Order) bool {
	return o.TotalGrossAmount.LessThanOrEqual(o.TotalPaidAmount)
}

func OrderIsPartlyPaid(o model.Order) bool {
	return o.TotalPaidAmount.GreaterThan(decimal.Zero)
}

var OrderGetTotalCapturedPrice = OrderGetTotalPaidPrice

func OrderGetTotalBalance(o model.Order) goprices.Money {
	var capturedPrice = OrderGetTotalCapturedPrice(o)
	var totalPrice = OrderGetTotalPrice(o)
	result, _ := capturedPrice.Sub(totalPrice.Gross)
	return *result
}

func OrderGetTotalWeight(o model.Order) measurement.Weight {
	return measurement.Weight{
		Amount: float64(o.WeightAmount),
		Unit:   measurement.WeightUnit(o.WeightUnit),
	}
}

type OrderFilterOption struct {
	CommonQueryOptions

	// INNER JOIN user ON ... WHERE user.email % ... OR user.first_name % ... OR user.last_name % ...
	// NOTE: this use trigram search feature
	Customer string

	Search          string
	ChannelIdOrSlug string

	PaymentChargeStatus qm.QueryMod // INNER JOIN payment ON ... WHERE payment.charge_status = ...
	Statuses            util.AnyArray[model.OrderStatus]

	AnnotateBillingAddressNames     bool
	AnnotateLastPaymentChargeStatus bool
}

type CustomOrder struct {
	model.Order
	OrderBillingAddressLastName  *string                    `boil:"billing_address_last_name" json:"billing_address_last_name"`
	OrderBillingAddressFirstName *string                    `boil:"billing_address_first_name" json:"billing_address_first_name"`
	OrderLastPaymentChargeStatus *model.PaymentChargeStatus `boil:"last_payment_charge_status" json:"last_payment_charge_status"`
}

type CustomOrderSlice []*CustomOrder
