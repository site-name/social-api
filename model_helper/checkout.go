package model_helper

import (
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
)

type CheckoutLineFilterOptions struct {
	CommonQueryOptions
}

type CheckoutFilterOptions struct {
	CommonQueryOptions
}

func CheckoutGetDiscountMoney(c model.Checkout) goprices.Money {
	money, _ := goprices.NewMoney(c.DiscountAmount.InexactFloat64(), string(c.Currency))
	return *money
}

func CheckoutSetDiscountAmount(c *model.Checkout, money goprices.Money) {
	c.DiscountAmount = money.GetAmount()
	c.Currency = model.Currency(money.GetCurrency())
}

func CheckoutAddDiscountAmount(c *model.Checkout, amount decimal.Decimal) {
	if c == nil {
		return
	}
	c.DiscountAmount = c.DiscountAmount.Add(amount)
}

func CheckoutPreSave(checkout *model.Checkout) {
	if checkout.Token == "" {
		checkout.Token = NewId()
	}
	if checkout.CreatedAt == 0 {
		checkout.CreatedAt = GetMillis()
	}
	checkout.UpdatedAt = checkout.CreatedAt
	checkoutCommonPre(checkout)
}

func CheckoutPreUpdate(checkout *model.Checkout) {
	checkoutCommonPre(checkout)
	checkout.UpdatedAt = GetMillis()
}

func checkoutCommonPre(c *model.Checkout) {
	c.Note = SanitizeUnicode(c.Note)
	c.Email = NormalizeEmail(c.Email)
	if c.LanguageCode.IsValid() != nil {
		c.LanguageCode = DEFAULT_LOCALE
	}
	if c.Currency.IsValid() != nil {
		c.Currency = DEFAULT_CURRENCY
	}
	if c.Country.IsValid() != nil {
		c.Country = DEFAULT_COUNTRY
	}
}

func CheckoutIsValid(c model.Checkout) *AppError {
	if !IsValidId(c.Token) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.token.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !IsValidEmail(c.Email) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.email.app_error", nil, "please provide valid email", http.StatusBadRequest)
	}
	if c.CreatedAt <= 0 {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.created_at.app_error", nil, "please provide valid created at", http.StatusBadRequest)
	}
	if c.UpdatedAt <= 0 {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.updated_at.app_error", nil, "please provide valid updated at", http.StatusBadRequest)
	}
	if c.Quantity < 0 {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.quantity.app_error", nil, "please provide valid quantity value", http.StatusBadRequest)
	}
	if !c.UserID.IsNil() && !IsValidId(*c.UserID.String) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(c.ChannelID) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}
	if !c.BillingAddressID.IsNil() && !IsValidId(*c.BillingAddressID.String) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.billing_address_id.app_error", nil, "please provide valid billing address id", http.StatusBadRequest)
	}
	if !c.ShippingAddressID.IsNil() && !IsValidId(*c.ShippingAddressID.String) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.shipping_address_id.app_error", nil, "please provide valid shipping address id", http.StatusBadRequest)
	}
	if !c.ShippingMethodID.IsNil() && !IsValidId(*c.ShippingMethodID.String) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.shipping_method_id.app_error", nil, "please provide valid shipping method id", http.StatusBadRequest)
	}
	if !c.CollectionPointID.IsNil() && !IsValidId(*c.CollectionPointID.String) {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.collection_point_id.app_error", nil, "please provide valid payment method id", http.StatusBadRequest)
	}
	if c.Currency.IsValid() != nil {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}
	if c.LanguageCode.IsValid() != nil {
		return NewAppError("CheckoutIsValid", "model.checkout.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func CheckoutLinePreSave(cl *model.CheckoutLine) {
	if cl.CreatedAt == 0 {
		cl.CreatedAt = GetMillis()
	}
	if cl.ID == "" {
		cl.ID = NewId()
	}
}

func CheckoutLineIsValid(cl model.CheckoutLine) *AppError {
	if cl.CreatedAt <= 0 {
		return NewAppError("CheckoutLineIsValid", "model.checkout_line.is_valid.created_at.app_error", nil, "please provide valid created at", http.StatusBadRequest)
	}
	if cl.Quantity < 0 {
		return NewAppError("CheckoutLineIsValid", "model.checkout_line.is_valid.quantity.app_error", nil, "please provide valid quantity value", http.StatusBadRequest)
	}
	if !IsValidId(cl.CheckoutID) {
		return NewAppError("CheckoutLineIsValid", "model.checkout_line.is_valid.checkout_id.app_error", nil, "please provide valid checkout id", http.StatusBadRequest)
	}
	if !IsValidId(cl.VariantID) {
		return NewAppError("CheckoutLineIsValid", "model.checkout_line.is_valid.variant_id.app_error", nil, "please provide valid variant id", http.StatusBadRequest)
	}

	return nil
}
