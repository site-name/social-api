package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
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
	Token                  string           `json:"token"` // uuid4, primary_key, NO EDITABLE
	CreateAt               int64            `json:"create_at"`
	UpdateAt               int64            `json:"update_at"`
	UserID                 *string          `json:"user_id"`
	ShopID                 string           `json:"shop_id"` // shop in which this checkout is placed
	Email                  string           `json:"email"`
	Quantity               int              `json:"quantity"`
	ChannelID              string           `json:"channel_id"`
	BillingAddressID       *string          `json:"billing_address_id,omitempty"`  // NO EDITABLE
	ShippingAddressID      *string          `json:"shipping_address_id,omitempty"` // NO EDITABLE
	ShippingMethodID       *string          `json:"shipping_method_id,omitempty"`
	CollectionPointID      *string          `json:"collection_point_id"` // foreign key *Warehouse
	Note                   string           `json:"note"`
	Currency               string           `json:"currency"`        // default "USD"
	Country                string           `json:"country"`         // one country only
	DiscountAmount         *decimal.Decimal `json:"discount_amount"` // default decimal(0)
	Discount               *goprices.Money  `db:"-" json:"discount,omitempty"`
	DiscountName           *string          `json:"discount_name"`
	TranslatedDiscountName *string          `json:"translated_discount_name"`
	VoucherCode            *string          `json:"voucher_code"`
	RedirectURL            *string          `json:"redirect_url"`
	TrackingCode           *string          `json:"tracking_code"`
	LanguageCode           string           `json:"language_code"`
	ModelMetadata

	channel *Channel `db:"-"`
}

// CheckoutFilterOption is used for bulding sql queries
type CheckoutFilterOption struct {
	Token           squirrel.Sqlizer
	UserID          squirrel.Sqlizer
	ChannelID       squirrel.Sqlizer
	Extra           squirrel.Sqlizer
	ChannelIsActive *bool

	SelectRelatedChannel bool
	Limit                int
}

func (c *Checkout) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"checkout.is_valid.%s.app_error",
		"checkout_token=",
		"Checkout.IsValid",
	)
	if c.UserID != nil && !IsValidId(*c.UserID) {
		return outer("user_id", &c.Token)
	}
	if !IsValidId(c.ShopID) {
		return outer("shop_id", &c.Token)
	}
	if !IsValidId(c.ChannelID) {
		return outer("channel_id", &c.Token)
	}
	if c.BillingAddressID != nil && !IsValidId(*c.BillingAddressID) {
		return outer("billing_address", &c.Token)
	}
	if c.ShippingAddressID != nil && !IsValidId(*c.ShippingAddressID) {
		return outer("shipping_address", &c.Token)
	}
	if c.ShippingMethodID != nil && !IsValidId(*c.ShippingMethodID) {
		return outer("shipping_method", &c.Token)
	}
	if c.CollectionPointID != nil && !IsValidId(*c.CollectionPointID) {
		return outer("collection_point_id", &c.Token)
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
	if !IsValidEmail(c.Email) || len(c.Email) > USER_EMAIL_MAX_LENGTH {
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
	if c.RedirectURL != nil && len(*c.RedirectURL) > URL_LINK_MAX_LENGTH {
		return outer("redirect_url", &c.Token)
	}
	if tag, err := language.Parse(c.LanguageCode); err != nil || !strings.EqualFold(tag.String(), c.LanguageCode) {
		return outer("language_code", &c.Token)
	}
	if c.TrackingCode != nil && len(*c.TrackingCode) > CHECKOUT_TRACKING_CODE_MAX_LENGTH || *c.TrackingCode == "" {
		return outer("tracking_code", &c.Token)
	}
	if _, ok := Countries[c.Country]; !ok {
		return outer("country", &c.Token)
	}

	return nil
}

// PopulateNonDbFields populates fields that are not saved to database.
// But are made of other fields belong to this struct.
func (c *Checkout) PopulateNonDbFields() {
	if c.DiscountAmount != nil && c.Currency != "" {
		c.Discount = &goprices.Money{
			Amount:   *c.DiscountAmount,
			Currency: c.Currency,
		}
	}
}

func (s *Checkout) SetChannel(c *Channel) {
	s.channel = c
}

func (s *Checkout) GetChannel() *Channel {
	return s.channel
}

func (c *Checkout) PreSave() {
	if c.Token == "" {
		c.Token = NewId()
	}

	c.CreateAt = GetMillis()
	c.UpdateAt = c.CreateAt
	c.commonPre()
}

func (c *Checkout) commonPre() {
	if c.Discount != nil {
		c.DiscountAmount = &c.Discount.Amount
	} else {
		c.DiscountAmount = &decimal.Zero
	}

	c.Note = SanitizeUnicode(c.Note)
	c.Email = NormalizeEmail(c.Email)
	if c.LanguageCode == "" {
		c.LanguageCode = DEFAULT_LOCALE
	}
	if c.Currency == "" {
		c.Currency = DEFAULT_CURRENCY
	} else {
		c.Currency = strings.ToUpper(c.Currency)
	}
	if c.Country == "" {
		c.Country = DEFAULT_COUNTRY
	} else {
		c.Country = strings.ToUpper(c.Country)
	}
}

func (c *Checkout) PreUpdate() {
	c.UpdateAt = GetMillis()
	c.commonPre()
}

func (c *Checkout) DeepCopy() *Checkout {
	if c == nil {
		return nil
	}

	res := *c
	if c.UserID != nil {
		res.UserID = NewPrimitive(*c.UserID)
	}
	if c.BillingAddressID != nil {
		res.BillingAddressID = NewPrimitive(*c.BillingAddressID)
	}
	if c.ShippingAddressID != nil {
		res.ShippingAddressID = NewPrimitive(*c.ShippingAddressID)
	}
	if c.ShippingMethodID != nil {
		res.ShippingMethodID = NewPrimitive(*c.ShippingMethodID)
	}
	if c.CollectionPointID != nil {
		res.CollectionPointID = NewPrimitive(*c.CollectionPointID)
	}
	if c.DiscountName != nil {
		res.DiscountName = NewPrimitive(*c.DiscountName)
	}
	if c.TranslatedDiscountName != nil {
		res.TranslatedDiscountName = NewPrimitive(*c.TranslatedDiscountName)
	}
	if c.VoucherCode != nil {
		res.VoucherCode = NewPrimitive(*c.VoucherCode)
	}
	if c.RedirectURL != nil {
		res.RedirectURL = NewPrimitive(*c.RedirectURL)
	}
	if c.TrackingCode != nil {
		res.TrackingCode = NewPrimitive(*c.TrackingCode)
	}
	if c.RedirectURL != nil {
		res.RedirectURL = NewPrimitive(*c.RedirectURL)
	}

	if c.DiscountAmount != nil {
		res.DiscountAmount = NewPrimitive(*c.DiscountAmount)
	}

	if c.channel != nil {
		res.channel = c.channel.DeepCopy()
	}
	return &res
}
