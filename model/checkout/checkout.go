package checkout

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
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
	Discount               *Money               `db:"-" json:"discount,omitempty"`
	DiscountName           *string              `json:"discount_name"`
	TranslatedDiscountName *string              `json:"translated_discount_name"`
	VoucherCode            *string              `json:"voucher_code"`
	GiftCards              []*giftcard.GiftCard `json:"gift_cards,omitempty" db:"-"`
	RedirectURL            *string              `json:"redirect_url"`
	TrackingCode           *string              `json:"tracking_code"`
	LanguageCode           string               `json:"language_code"`
}

func (c *Checkout) checkoutAppErr(field string) *model.AppError {
	var details string
	if strings.ToLower(field) != "id" {
		details += "checkout_id=" + c.Id
	}
	id := fmt.Sprintf("model.checkout.is_valid.%s.app_error", field)

	return model.NewAppError("Checkout.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *Checkout) IsValid() *model.AppError {
	if c.Id == "" {
		return c.checkoutAppErr("id")
	}
	if c.UserID != nil && *c.UserID == "" {
		return c.checkoutAppErr("user_id")
	}
	if c.BillingAddressID != nil && *c.BillingAddressID == "" {
		return c.checkoutAppErr("billing_address")
	}
	if c.ShippingAddressID != nil && *c.ShippingAddressID == "" {
		return c.checkoutAppErr("shipping_address")
	}
	if c.Currency == "" || len(c.Currency) > model.MAX_LENGTH_CURRENCY_CODE {
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
	if !model.IsValidEmail(c.Email) || len(c.Email) > model.USER_EMAIL_MAX_LENGTH {
		return c.checkoutAppErr("email")
	}
	if c.DiscountName != nil && utf8.RuneCountInString(*c.DiscountName) > CHECKOUT_DISCOUNT_NAME_MAX_LENGTH {
		return c.checkoutAppErr("discount_name")
	}
	if c.TranslatedDiscountName != nil && utf8.RuneCountInString(*c.TranslatedDiscountName) > CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH {
		return c.checkoutAppErr("translated_discount_name")
	}
	if c.VoucherCode != nil && len(*c.VoucherCode) > CHECKOUT_VOUCHER_CODE_MAX_LENGTH || *c.VoucherCode == "" {
		return c.checkoutAppErr("voucher_code")
	}
	if c.LanguageCode == "" {
		return c.checkoutAppErr("language_code")
	}
	if tag, err := language.Parse(c.LanguageCode); err != nil || !strings.EqualFold(tag.String(), c.LanguageCode) {
		return c.checkoutAppErr("language_code")
	}
	if c.TrackingCode != nil && len(*c.TrackingCode) > CHECKOUT_TRACKING_CODE_MAX_LENGTH || *c.TrackingCode == "" {
		return c.checkoutAppErr("tracking_code")
	}
	if _, ok := model.Countries[strings.ToUpper(c.Country)]; !ok {
		return c.checkoutAppErr("country")
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

func (c *Checkout) GetTotalGiftCardsBalance() *Money {
	panic("not impl")
}

func (c *Checkout) GetTotalWeight() {
	panic("not impl")
}
