package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

// max lengths for Checkout table
const (
	CHECKOUT_DISCOUNT_NAME_MAX_LENGTH            = 255
	CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH = 255
	CHECKOUT_TRACKING_CODE_MAX_LENGTH            = 255
	CHECKOUT_VOUCHER_CODE_MAX_LENGTH             = 12
	CHECKOUT_LANGUAGE_CODE_MAX_LENGTH            = 35
)

// A Shopping checkout.
// Ordering by CreateAt ASC
type Checkout struct {
	Token                  string           `json:"token" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Token"` // uuid4, primary_key, NO EDITABLE
	CreateAt               int64            `json:"create_at" gorm:"column:CreateAt;autoCreateTime:milli"`
	UpdateAt               int64            `json:"update_at" gorm:"column:UpdateAt;autoUpdateTime:milli"`
	UserID                 *string          `json:"user_id" gorm:"type:uuid;column:UserID"`
	Email                  string           `json:"email" gorm:"type:varchar(128);column:Email"`
	Quantity               int              `json:"quantity" gorm:"column:Quantity"`
	ChannelID              string           `json:"channel_id" gorm:"column:ChannelID;type:uuid"`
	BillingAddressID       *string          `json:"billing_address_id,omitempty" gorm:"column:BillingAddressID;type:uuid"`   // NOT EDITABLE
	ShippingAddressID      *string          `json:"shipping_address_id,omitempty" gorm:"column:ShippingAddressID;type:uuid"` // NOT EDITABLE
	ShippingMethodID       *string          `json:"shipping_method_id,omitempty" gorm:"type:uuid;column:ShippingMethodID"`
	CollectionPointID      *string          `json:"collection_point_id" gorm:"type:uuid;column:CollectionPointID"` // foreign key *Warehouse
	Note                   string           `json:"note" gorm:"column:Note"`
	Currency               string           `json:"currency" gorm:"type:varchar(3);column:Currency"`        // default "USD"
	Country                CountryCode      `json:"country" gorm:"column:Country;type:varchar(20)"`         // one country only
	DiscountAmount         *decimal.Decimal `json:"discount_amount" gorm:"column:DiscountAmount;default:0"` // default decimal(0)
	DiscountName           *string          `json:"discount_name" gorm:"column:DiscountName;type:varchar(255)"`
	TranslatedDiscountName *string          `json:"translated_discount_name" gorm:"type:varchar(255);column:TranslatedDiscountName"`
	VoucherCode            *string          `json:"voucher_code" gorm:"column:VoucherCode;type:varchar(12)"`
	RedirectURL            *string          `json:"redirect_url" gorm:"type:varchar(200);column:RedirectURL"`
	TrackingCode           *string          `json:"tracking_code" gorm:"type:varchar(255);column:TrackingCode"`
	LanguageCode           LanguageCodeEnum `json:"language_code" gorm:"type:varchar(35);column:LanguageCode"` // default to DEFAULT_LOCALE
	ModelMetadata

	channel        *Channel        `gorm:"-"`
	billingAddress *Address        `gorm:"-"`
	user           *User           `gorm:"-"`
	Discount       *goprices.Money `gorm:"-" json:"discount,omitempty"`
	Giftcards      Giftcards       `json:"-" gorm:"many2many:GiftcardCheckouts"`
}

// column names of table checkout
const (
	CheckoutColumnToken                  = "Token"
	CheckoutColumnCreateAt               = "CreateAt"
	CheckoutColumnUpdateAt               = "UpdateAt"
	CheckoutColumnUserID                 = "UserID"
	CheckoutColumnEmail                  = "Email"
	CheckoutColumnQuantity               = "Quantity"
	CheckoutColumnChannelID              = "ChannelID"
	CheckoutColumnBillingAddressID       = "BillingAddressID"
	CheckoutColumnShippingAddressID      = "ShippingAddressID"
	CheckoutColumnShippingMethodID       = "ShippingMethodID"
	CheckoutColumnCollectionPointID      = "CollectionPointID"
	CheckoutColumnNote                   = "Note"
	CheckoutColumnCurrency               = "Currency"
	CheckoutColumnCountry                = "Country"
	CheckoutColumnDiscountAmount         = "DiscountAmount"
	CheckoutColumnDiscountName           = "DiscountName"
	CheckoutColumnTranslatedDiscountName = "TranslatedDiscountName"
	CheckoutColumnVoucherCode            = "VoucherCode"
	CheckoutColumnRedirectURL            = "RedirectURL"
	CheckoutColumnTrackingCode           = "TrackingCode"
	CheckoutColumnLanguageCode           = "LanguageCode"
)

func (c *Checkout) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *Checkout) BeforeUpdate(_ *gorm.DB) error {
	c.PreUpdate()
	return c.IsValid()
}
func (c *Checkout) TableName() string              { return CheckoutTableName }
func (c *Checkout) SetChannel(ch *Channel)         { c.channel = ch }
func (c *Checkout) GetChannel() *Channel           { return c.channel }
func (c *Checkout) SetBilingAddress(addr *Address) { c.billingAddress = addr }
func (c *Checkout) GetBilingAddress() *Address     { return c.billingAddress }
func (c *Checkout) SetUser(u *User)                { c.user = u }
func (c *Checkout) GetUser() *User                 { return c.user }

// CheckoutFilterOption is used for bulding sql queries
type CheckoutFilterOption struct {
	Conditions squirrel.Sqlizer

	ChannelIsActive squirrel.Sqlizer // INNER JOIN Channels ON ... WHERE Channels.IsActive = ?

	SelectRelatedChannel        bool // this will populate the field `channel`
	SelectRelatedBillingAddress bool // this will populate the field 'billingAddress'
	SelectRelatedUser           bool

	GraphqlPaginationValues GraphqlPaginationValues
	CountTotal              bool
}

func (c *Checkout) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.checkout.is_valid.%s.app_error",
		"checkout_token=",
		"Checkout.IsValid",
	)
	if c.UserID != nil && !IsValidId(*c.UserID) {
		return outer("user_id", &c.Token)
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
	if !IsValidEmail(c.Email) {
		return outer("email", &c.Token)
	}
	if !c.LanguageCode.IsValid() {
		return outer("language_code", &c.Token)
	}
	if !c.Country.IsValid() {
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

func (c *Checkout) PreSave() {
	c.commonPre()
}

func (c *Checkout) commonPre() {
	if c.Discount != nil {
		c.DiscountAmount = &c.Discount.Amount
	} else {
		c.DiscountAmount = GetPointerOfValue(decimal.Zero)
	}

	c.Note = SanitizeUnicode(c.Note)
	c.Email = NormalizeEmail(c.Email)
	if c.LanguageCode == "" {
		c.LanguageCode = DEFAULT_LOCALE
	}
	if c.Currency == "" {
		c.Currency = DEFAULT_CURRENCY
	}
	if c.Country == "" {
		c.Country = DEFAULT_COUNTRY
	}
}

func (c *Checkout) PreUpdate() {
	c.commonPre()
}

func (c *Checkout) DeepCopy() *Checkout {
	if c == nil {
		return nil
	}

	res := *c
	if c.UserID != nil {
		res.UserID = GetPointerOfValue(*c.UserID)
	}
	if c.BillingAddressID != nil {
		res.BillingAddressID = GetPointerOfValue(*c.BillingAddressID)
	}
	if c.ShippingAddressID != nil {
		res.ShippingAddressID = GetPointerOfValue(*c.ShippingAddressID)
	}
	if c.ShippingMethodID != nil {
		res.ShippingMethodID = GetPointerOfValue(*c.ShippingMethodID)
	}
	if c.CollectionPointID != nil {
		res.CollectionPointID = GetPointerOfValue(*c.CollectionPointID)
	}
	if c.DiscountName != nil {
		res.DiscountName = GetPointerOfValue(*c.DiscountName)
	}
	if c.TranslatedDiscountName != nil {
		res.TranslatedDiscountName = GetPointerOfValue(*c.TranslatedDiscountName)
	}
	if c.VoucherCode != nil {
		res.VoucherCode = GetPointerOfValue(*c.VoucherCode)
	}
	if c.RedirectURL != nil {
		res.RedirectURL = GetPointerOfValue(*c.RedirectURL)
	}
	if c.TrackingCode != nil {
		res.TrackingCode = GetPointerOfValue(*c.TrackingCode)
	}
	if c.RedirectURL != nil {
		res.RedirectURL = GetPointerOfValue(*c.RedirectURL)
	}

	if c.DiscountAmount != nil {
		res.DiscountAmount = GetPointerOfValue(*c.DiscountAmount)
	}

	if c.channel != nil {
		res.channel = c.channel.DeepCopy()
	}
	return &res
}
