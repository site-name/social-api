package checkout

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// max lengths for Checkout table
const (
	CHECKOUT_DISCOUNT_NAME_MAX_LENGTH            = 255
	CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH = 255
	CHECKOUT_TRACKING_CODE_MAX_LENGTH            = 255
	CHECKOUT_VOUCHER_CODE_MAX_LENGTH             = 12
	CHECKOUT_LANGUAGE_CODE_MAX_LENGTH            = 35
)

// A Shopping checkout
type Checkout struct {
	Token                  string           `json:"token"` // uuid4, primary_key
	CreateAt               int64            `json:"create_at"`
	UpdateAt               int64            `json:"update_at"`
	UserID                 *string          `json:"user_id"`
	Email                  string           `json:"email"`
	Quantity               uint             `json:"quantity"`
	ChannelID              string           `json:"channel_id"`
	BillingAddressID       *string          `json:"billing_address_id,omitempty"`
	ShippingAddressID      *string          `json:"shipping_address_id,omitempty"`
	ShippingMethodID       *string          `json:"shipping_method_id,omitempty"`
	Note                   string           `json:"note"`
	Currency               string           `json:"currency"`
	Country                string           `json:"country"` // one country only
	DiscountAmount         *decimal.Decimal `json:"discount_amount"`
	Discount               *goprices.Money  `db:"-" json:"discount,omitempty"`
	DiscountName           *string          `json:"discount_name"`
	TranslatedDiscountName *string          `json:"translated_discount_name"`
	VoucherCode            *string          `json:"voucher_code"`
	RedirectURL            *string          `json:"redirect_url"`
	TrackingCode           *string          `json:"tracking_code"`
	LanguageCode           string           `json:"language_code"`
	model.ModelMetadata
}

func (c *Checkout) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.checkout.is_valid.%s.app_error",
		"checkout_token=",
		"Checkout.IsValid",
	)
	if c.UserID != nil && !model.IsValidId(*c.UserID) {
		return outer("user_id", &c.Token)
	}
	if c.BillingAddressID != nil && !model.IsValidId(*c.BillingAddressID) {
		return outer("billing_address", &c.Token)
	}
	if c.ShippingAddressID != nil && !model.IsValidId(*c.ShippingAddressID) {
		return outer("shipping_address", &c.Token)
	}
	if un, err := currency.ParseISO(c.Currency); err != nil || !strings.EqualFold(un.String(), c.Currency) {
		return outer("currency", &c.Token)
	}
	if c.CreateAt == 0 {
		return outer("create_at", &c.Token)
	}
	if c.UpdateAt == 0 {
		return outer("update_at", &c.Token)
	}
	if !model.IsValidEmail(c.Email) || len(c.Email) > model.USER_EMAIL_MAX_LENGTH {
		return outer("email", &c.Token)
	}
	if c.DiscountName != nil && utf8.RuneCountInString(*c.DiscountName) > CHECKOUT_DISCOUNT_NAME_MAX_LENGTH {
		return outer("discount_name", &c.Token)
	}
	if c.TranslatedDiscountName != nil && utf8.RuneCountInString(*c.TranslatedDiscountName) > CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH {
		return outer("translated_discount_name", &c.Token)
	}
	if c.VoucherCode != nil && len(*c.VoucherCode) > CHECKOUT_VOUCHER_CODE_MAX_LENGTH || *c.VoucherCode == "" {
		return outer("voucher_code", &c.Token)
	}
	if c.RedirectURL != nil && len(*c.RedirectURL) > model.URL_LINK_MAX_LENGTH {
		return outer("redirect_url", &c.Token)
	}
	if tag, err := language.Parse(c.LanguageCode); err != nil || !strings.EqualFold(tag.String(), c.LanguageCode) {
		return outer("language_code", &c.Token)
	}
	if c.TrackingCode != nil && len(*c.TrackingCode) > CHECKOUT_TRACKING_CODE_MAX_LENGTH || *c.TrackingCode == "" {
		return outer("tracking_code", &c.Token)
	}
	if _, ok := model.Countries[c.Country]; !ok {
		return outer("country", &c.Token)
	}

	return nil
}

func (c *Checkout) ToJson() string {
	c.Discount = &goprices.Money{
		Amount:   c.DiscountAmount,
		Currency: c.Currency,
	}
	return model.ModelToJson(c)
}

func CheckoutFromJson(data io.Reader) *Checkout {
	var checkout Checkout
	model.ModelFromJson(&checkout, data)
	return &checkout
}

func (c *Checkout) PreSave() {
	if c.Token == "" {
		c.Token = model.NewId()
	}
	if c.LanguageCode == "" {
		c.LanguageCode = model.DEFAULT_LOCALE
	}
	if c.DiscountAmount == nil {
		c.DiscountAmount = &decimal.Zero
	}
	if c.Country == "" {
		c.Country = model.DEFAULT_COUNTRY
	} else {
		c.Country = strings.ToUpper(strings.TrimSpace(c.Country))
	}
	if c.Currency == "" {
		c.Currency = model.DEFAULT_CURRENCY
	} else {
		c.Currency = strings.ToUpper(strings.TrimSpace(c.Currency))
	}
	c.Note = model.SanitizeUnicode(c.Note)

	c.Email = model.NormalizeEmail(c.Email)

	c.CreateAt = model.GetMillis()
	c.UpdateAt = c.CreateAt
}

func (c *Checkout) PreUpdate() {
	c.Note = model.SanitizeUnicode(c.Note)
	c.Email = model.NormalizeEmail(c.Email)
	c.UpdateAt = model.GetMillis()
	if c.Country == "" {
		c.Country = model.DEFAULT_COUNTRY
	} else {
		c.Country = strings.ToUpper(strings.TrimSpace(c.Country))
	}
	if c.Currency == "" {
		c.Currency = model.DEFAULT_CURRENCY
	} else {
		c.Currency = strings.ToUpper(strings.TrimSpace(c.Currency))
	}
}
