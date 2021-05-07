package checkout

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
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

// A Shopping checkout
type Checkout struct {
	Id                     string               `json:"id"`
	CreateAt               int64                `json:"create_at"`
	UpdateAt               int64                `json:"update_at"`
	UserID                 *string              `json:"user_id"`
	Email                  string               `json:"email"`
	Token                  string               `json:"token"`
	Quantity               uint                 `json:"quantity"`
	ChannelID              string               `json:"channel_id"`
	BillingAddressID       *string              `json:"billing_address_id,omitempty"`
	ShippingAddressID      *string              `json:"shipping_address_id,omitempty"`
	ShippingMethodID       *string              `json:"shipping_method_id,omitempty"`
	Note                   string               `json:"note"`
	Currency               string               `json:"currency"`
	Country                string               `json:"country"`
	DiscountAmount         *decimal.Decimal     `json:"discount_amount"`
	Discount               *model.Money         `db:"-" json:"discount,omitempty"`
	DiscountName           *string              `json:"discount_name"`
	TranslatedDiscountName *string              `json:"translated_discount_name"`
	VoucherCode            *string              `json:"voucher_code"`
	GiftCards              []*giftcard.GiftCard `json:"gift_cards,omitempty" db:"-"`
	RedirectURL            *string              `json:"redirect_url"`
	TrackingCode           *string              `json:"tracking_code"`
	LanguageCode           string               `json:"language_code"`
}

func (c *Checkout) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.checkout.is_valid.%s.app_error",
		"checkout_id=",
		"Checkout.IsValid",
	)
	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.UserID != nil && !model.IsValidId(*c.UserID) {
		return outer("user_id", &c.Id)
	}
	if c.BillingAddressID != nil && !model.IsValidId(*c.BillingAddressID) {
		return outer("billing_address", &c.Id)
	}
	if c.ShippingAddressID != nil && !model.IsValidId(*c.ShippingAddressID) {
		return outer("shipping_address", &c.Id)
	}
	if c.Currency == "" || len(c.Currency) > model.MAX_LENGTH_CURRENCY_CODE {
		return outer("currency", &c.Id)
	}
	if un, err := currency.ParseISO(c.Currency); err != nil || !strings.EqualFold(un.String(), c.Currency) {
		return outer("currency", &c.Id)
	}
	if c.CreateAt == 0 {
		return outer("create_at", &c.Id)
	}
	if c.UpdateAt == 0 {
		return outer("update_at", &c.Id)
	}
	if !model.IsValidEmail(c.Email) || len(c.Email) > model.USER_EMAIL_MAX_LENGTH {
		return outer("email", &c.Id)
	}
	if c.DiscountName != nil && utf8.RuneCountInString(*c.DiscountName) > CHECKOUT_DISCOUNT_NAME_MAX_LENGTH {
		return outer("discount_name", &c.Id)
	}
	if c.TranslatedDiscountName != nil && utf8.RuneCountInString(*c.TranslatedDiscountName) > CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH {
		return outer("translated_discount_name", &c.Id)
	}
	if c.VoucherCode != nil && len(*c.VoucherCode) > CHECKOUT_VOUCHER_CODE_MAX_LENGTH || *c.VoucherCode == "" {
		return outer("voucher_code", &c.Id)
	}
	if tag, err := language.Parse(c.LanguageCode); err != nil || !strings.EqualFold(tag.String(), c.LanguageCode) {
		return outer("language_code", &c.Id)
	}
	if c.TrackingCode != nil && len(*c.TrackingCode) > CHECKOUT_TRACKING_CODE_MAX_LENGTH || *c.TrackingCode == "" {
		return outer("tracking_code", &c.Id)
	}
	if _, ok := model.Countries[strings.ToUpper(c.Country)]; !ok {
		return outer("country", &c.Id)
	}

	return nil
}

func (c *Checkout) ToJson() string {
	c.Discount = &model.Money{
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
	if c.Id == "" {
		c.Id = model.NewId()
	}
	if c.LanguageCode == "" {
		c.LanguageCode = "en"
	}
	if c.DiscountAmount == nil {
		c.DiscountAmount = &decimal.Zero
	}
	if c.Country == "" {
		c.Country = "US"
	}
	if c.Token == "" {
		c.Token = model.NewId()
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
}

func (c *Checkout) GetCustomerEmail() string {
	panic("not implemented")
}

func (c *Checkout) IsShippingRequired() bool {
	panic("not implemented")
}

func (c *Checkout) GetTotalGiftCardsBalance() *model.Money {
	panic("not impl")
}

func (c *Checkout) GetTotalWeight() {
	panic("not impl")
}
