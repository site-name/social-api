package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/modules/json"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

const (
	CHECKOUT_DISCOUNT_NAME_MAX_LENGTH            = 255
	CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH = 255
	CHECKOUT_TRACKING_CODE_MAX_LENGTH            = 255
	CHECKOUT_VOUCHER_CODE_MAX_LENGTH             = 12
	CHECKOUT_LANGUAGE_CODE_MAX_LENGTH            = 35
)

type Money struct {
	Amount   *decimal.Decimal
	Currency string
}

// A Shopping checkout
type Checkout struct {
	Id                     string           `json:"id"`
	CreateAt               int64            `json:"create_at"`
	UpdateAt               int64            `json:"update_at"`
	UserID                 string           `json:"user_id"`
	Email                  string           `json:"email"`
	Token                  string           `json:"token"`
	Quantity               uint             `json:"quantity"`
	ChannelID              string           `json:"channel_id"`
	BillingAddressID       string           `json:"billing_address_id"`
	ShippingAddressID      string           `json:"shipping_address_id"`
	ShippingMethodID       string           `json:"shipping_method_id"`
	Note                   string           `json:"note"`
	Currency               string           `json:"currency"`
	Country                string           `json:"country"`
	DiscountAmount         *decimal.Decimal `json:"discount_amount"`
	Discount               *Money           `db:"-" json:"discount"`
	DiscountName           string           `json:"discount_name"`
	TranslatedDiscountName string           `json:"translated_discount_name"`
	VoucherCode            string           `json:"voucher_code"`
	GiftCards              []*GiftCard      `json:"gift_cards" db:"-"`
	RedirectURL            string           `json:"redirect_url"`
	TrackingCode           string           `json:"tracking_code"`
	LanguageCode           string           `json:"language_code"`
}

func (c *Checkout) checkoutAppErr(field string) *AppError {
	var details string
	if strings.ToLower(field) != "id" {
		details += "checkout_id=" + c.Id
	}
	id := fmt.Sprintf("model.checkout.is_valid.%s.app_error", field)

	return NewAppError("Checkout.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *Checkout) IsValid() *AppError {
	if c.Id == "" {
		return c.checkoutAppErr("id")
	}
	if c.UserID == "" {
		return c.checkoutAppErr("user_id")
	}
	if c.Currency == "" || len(c.Currency) > MAX_LENGTH_CURRENCY_CODE {
		return c.checkoutAppErr("currency")
	}
	if un, err := currency.ParseISO(c.Currency); err != nil || !strings.EqualFold(un.String(), c.Currency) {
		return c.checkoutAppErr("currency")
	}
	if c.CreateAt == 0 {
		return c.checkoutAppErr("create_at")
	}
	if c.UpdateAt == 0 {
		return c.checkoutAppErr("update_at")
	}
	if !IsValidEmail(c.Email) || len(c.Email) > USER_EMAIL_MAX_LENGTH {
		return c.checkoutAppErr("email")
	}
	if len(c.DiscountName) > CHECKOUT_DISCOUNT_NAME_MAX_LENGTH {
		return c.checkoutAppErr("discount_name")
	}
	if len(c.TranslatedDiscountName) > CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH {
		return c.checkoutAppErr("translated_discount_name")
	}
	if len(c.VoucherCode) > CHECKOUT_VOUCHER_CODE_MAX_LENGTH {
		return c.checkoutAppErr("voucher_code")
	}
	if c.LanguageCode == "" {
		return c.checkoutAppErr("language_code")
	}
	if tag, err := language.Parse(c.LanguageCode); err != nil || !strings.EqualFold(tag.String(), c.LanguageCode) {
		return c.checkoutAppErr("language_code")
	}
	if len(c.TrackingCode) > CHECKOUT_TRACKING_CODE_MAX_LENGTH || c.TrackingCode == "" {
		return c.checkoutAppErr("tracking_code")
	}

	return nil
}

func (c *Checkout) ToJson() string {
	b, _ := json.JSON.Marshal(c)
	return string(b)
}

func CheckoutFromJson(data io.Reader) *Checkout {
	var checkout Checkout
	err := json.JSON.NewDecoder(data).Decode(&checkout)
	if err != nil {
		return nil
	}
	return &checkout
}

func (c *Checkout) GetCustomerEmail() string {
	panic("not implemented")
}

func (c *Checkout) IsShippingRequired() bool {
	panic("not implemented")
}

func (c *Checkout) GetTotalGiftCardsBalance() *Money {
	panic("not impl")
}

func (c *Checkout) GetTotalWeight() {
	panic("not impl")
}
