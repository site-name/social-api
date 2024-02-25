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
	return goprices.Money{
		Amount:   c.DiscountAmount,
		Currency: c.Currency.String(),
	}
}

func CheckoutAddDiscountAmount(c *model.Checkout, amount decimal.Decimal) {
	if c == nil {
		return
	}
	c.DiscountAmount = c.DiscountAmount.Add(amount)
}

func CheckoutPreSave(checkout *model.Checkout) {
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

	return nil
}

func CheckoutLinePreSave(cl *model.CheckoutLine) {
	if cl.CreatedAt == 0 {
		cl.CreatedAt = GetMillis()
	}
}

func CheckoutLineIsValid(cl model.CheckoutLine) *AppError {
	if cl.CreatedAt <= 0 {
		return NewAppError("CheckoutLineIsValid", "model.checkout_line.is_valid.created_at.app_error", nil, "please provide valid created at", http.StatusBadRequest)
	}
	if cl.Quantity < 0 {
		return NewAppError("CheckoutLineIsValid", "model.checkout_line.is_valid.quantity.app_error", nil, "please provide valid quantity value", http.StatusBadRequest)
	}
	return nil
}
