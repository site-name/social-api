package model

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

// A Shopping checkout.
// Ordering by CreateAt ASC
type Checkout struct {
	Token                  UUID             `json:"token" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Token"` // uuid4, primary_key, NO EDITABLE
	CreateAt               int64            `json:"create_at" gorm:"column:CreateAt;autoCreateTime:milli"`
	UpdateAt               int64            `json:"update_at" gorm:"column:UpdateAt;autoUpdateTime:milli"`
	UserID                 *UUID            `json:"user_id" gorm:"type:uuid;column:UserID"`
	Email                  string           `json:"email" gorm:"type:varchar(128);column:Email"`
	Quantity               int              `json:"quantity" gorm:"column:Quantity"`
	ChannelID              UUID             `json:"channel_id" gorm:"column:ChannelID;type:uuid"`
	BillingAddressID       *UUID            `json:"billing_address_id,omitempty" gorm:"column:BillingAddressID;type:uuid"`   // NOT EDITABLE
	ShippingAddressID      *UUID            `json:"shipping_address_id,omitempty" gorm:"column:ShippingAddressID;type:uuid"` // NOT EDITABLE
	ShippingMethodID       *UUID            `json:"shipping_method_id,omitempty" gorm:"type:uuid;column:ShippingMethodID"`
	CollectionPointID      *UUID            `json:"collection_point_id" gorm:"type:uuid;column:CollectionPointID"` // foreign key *Warehouse
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
	if c.UserID != nil && !IsValidId(*c.UserID) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(c.ChannelID) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}
	if c.BillingAddressID != nil && !IsValidId(*c.BillingAddressID) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.billing_address_id.app_error", nil, "please provide valid billing address id", http.StatusBadRequest)
	}
	if c.ShippingAddressID != nil && !IsValidId(*c.ShippingAddressID) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.shipping_address_id.app_error", nil, "please provide valid shipping address id", http.StatusBadRequest)
	}
	if c.ShippingMethodID != nil && !IsValidId(*c.ShippingMethodID) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.shipping_method_id.app_error", nil, "please provide valid shipping method id", http.StatusBadRequest)
	}
	if c.CollectionPointID != nil && !IsValidId(*c.CollectionPointID) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.collection_point_id.app_error", nil, "please provide valid collection point id", http.StatusBadRequest)
	}
	if un, err := currency.ParseISO(c.Currency); err != nil || !strings.EqualFold(un.String(), c.Currency) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}
	if !IsValidEmail(c.Email) {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.email.app_error", nil, "please provide valid email", http.StatusBadRequest)
	}
	if !c.LanguageCode.IsValid() {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}
	if !c.Country.IsValid() {
		return NewAppError("Checkout.IsValid", "model.checkout.is_valid.country.app_error", nil, "please provide valid country", http.StatusBadRequest)
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
		c.DiscountAmount = &decimal.Zero
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
