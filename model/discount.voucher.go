package model

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	"gorm.io/gorm"
)

type VoucherType string

func (t *VoucherType) IsValid() bool {
	switch *t {
	case VOUCHER_TYPE_SHIPPING, VOUCHER_TYPE_ENTIRE_ORDER, VOUCHER_TYPE_SPECIFIC_PRODUCT:
		return true
	default:
		return false
	}
}

// Applicable values for Voucher type
const (
	VOUCHER_TYPE_SHIPPING         VoucherType = "shipping"
	VOUCHER_TYPE_ENTIRE_ORDER     VoucherType = "entire_order"
	VOUCHER_TYPE_SPECIFIC_PRODUCT VoucherType = "specific_product"
)

type DiscountValueType string

// Applicable values for voucher's discount value type
const (
	DISCOUNT_VALUE_TYPE_FIXED      DiscountValueType = "fixed"
	DISCOUNT_VALUE_TYPE_PERCENTAGE DiscountValueType = "percentage"
)

func (e DiscountValueType) IsValid() bool {
	switch e {
	case DISCOUNT_VALUE_TYPE_FIXED, DISCOUNT_VALUE_TYPE_PERCENTAGE:
		return true
	}
	return false
}

type Voucher struct {
	Id                       string            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Type                     VoucherType       `json:"type" gorm:"type:varchar(20);column:Type"` // default to "entire_order"
	Name                     *string           `json:"name" gorm:"type:varchar(255);column:Name"`
	Code                     string            `json:"code" gorm:"type:varchar(16);column:Code"` // UNIQUE, has format of XXXX-XXXX-XXXX
	UsageLimit               *int              `json:"usage_limit" gorm:"column:UsageLimit"`
	Used                     int               `json:"used" gorm:"column:Used"` // not editable
	StartDate                time.Time         `json:"start_date" gorm:"column:StartDate;autoCreateTime:milli"`
	EndDate                  *time.Time        `json:"end_date" gorm:"column:EndDate"`
	ApplyOncePerOrder        bool              `json:"apply_once_per_order" gorm:"column:ApplyOncePerOrder"`
	ApplyOncePerCustomer     bool              `json:"apply_once_per_customer" gorm:"column:ApplyOncePerCustomer"`
	OnlyForStaff             *bool             `json:"only_for_staff" gorm:"default:false;column:OnlyForStaff"` // default false
	DiscountValueType        DiscountValueType `json:"discount_value_type" gorm:"type:varchar(10);column:DiscountValueType"`
	Countries                string            `json:"countries" gorm:"type:varchar(1000);column:Countries"` // multiple. E.g: "VN US CN"
	MinCheckoutItemsQuantity int               `json:"min_checkout_items_quantity" gorm:"column:MinCheckoutItemsQuantity"`
	CreateAt                 int64             `json:"create_at" gorm:"autoCreateTime:milli;column:CreateAt"` // this field is for ordering
	UpdateAt                 int64             `json:"update_at" gorm:"autoCreateTime:milli;autoUpdateTime:milli;column:UpdateAt"`
	ModelMetadata

	Products               Products                `json:"-" gorm:"many2many:VoucherProducts"`
	Categories             Categories              `json:"-" gorm:"many2many:VoucherCategories"`
	ProductVariants        ProductVariants         `json:"-" gorm:"many2many:VoucherVariants"`
	Collections            Collections             `json:"-" gorm:"many2many:VoucherCollections"`
	VoucherChannelListings []VoucherChannelListing `json:"-" gorm:"foreignKey:VoucherID;constraint:OnDelete:CASCADE"`

	MinSpentAmount *decimal.Decimal `json:"-" gorm:"-"` // this field is used for sorting vouchers.
	DiscountValue  *decimal.Decimal `json:"-" gorm:"-"` // this field is used for sorting vouchers.
}

func (c *Voucher) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Voucher) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Voucher) TableName() string             { return VoucherTableName }

// VoucherFilterOption
type VoucherFilterOption struct {
	Conditions squirrel.Sqlizer

	CountTotal bool // if true, store counts all vouchers that satisfy options

	VoucherChannelListing_ChannelIsActive squirrel.Sqlizer // INNER JOIN VoucherChannelListings ON ... INNER JOIN Channels ON ... WHERE Channels.IsActive ...
	VoucherChannelListing_ChannelSlug     squirrel.Sqlizer // INNER JOIN VoucherChannelListings ON ... INNER JOIN Channels ON ... WHERE Channels.Slug ...
	ForUpdate                             bool             // this add FOR UPDATE to sql queries, NOTE: Only applied if Transaction field is non-nil
	Transaction                           *gorm.DB

	// if true, the field `MinSpentAmount` will be populated.
	// Please visit store_voucher.go for details.
	Annotate_MinimumSpentAmount bool
	// if true, the field `DiscountValue` will be populated.
	// Please visit store_voucher.go for details.
	Annotate_DiscountValue bool
	// NOTE: if `Annotate_MinimumSpentAmount` or `Annotate_DiscountValue`, this field must provided,
	// otherwise store will return *store.ErrInvalidInput error
	ChannelSlug string

	GraphqlPaginationValues GraphqlPaginationValues
}

// ValidateMinCheckoutItemsQuantity validates the quantity >= minimum requirement
func (voucher *Voucher) ValidateMinCheckoutItemsQuantity(quantity int) *NotApplicable {
	if voucher.MinCheckoutItemsQuantity > quantity {
		return &NotApplicable{
			Where:                    "ValidateMinCheckoutItemsQuantity",
			Message:                  fmt.Sprintf("This offer is onlyvalid for orders with a minimum of %d in quantity", voucher.MinCheckoutItemsQuantity),
			MinCheckoutItemsQuantity: voucher.MinCheckoutItemsQuantity,
		}
	}
	return nil
}

type Vouchers []*Voucher

func (v *Voucher) IsValid() *AppError {
	if !v.Type.IsValid() {
		return NewAppError("Voucher.IsValid", "model.voucher.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	if !v.DiscountValueType.IsValid() {
		return NewAppError("Voucher.IsValid", "model.voucher.is_valid.discount_value_type.app_error", nil, "please provide valid discount value type", http.StatusBadRequest)
	}
	for _, country := range strings.Fields(v.Countries) {
		if !CountryCode(country).IsValid() {
			return NewAppError("Voucher.IsValid", "model.voucher.is_valid.countries.app_error", nil, "please provide valid countries", http.StatusBadRequest)
		}
	}
	if !PromoCodeRegex.MatchString(v.Code) {
		return NewAppError("Voucher.IsValid", "model.voucher.is_valid.code.app_error", nil, "code must look like 78GH-UJKI-90RD", http.StatusBadRequest)
	}
	if v.EndDate != nil && v.StartDate.After(*v.EndDate) {
		return NewAppError("Voucher.IsValid", "model.voucher.is_valid.date.app_error", nil, "end date must be after start date", http.StatusBadRequest)
	}
	if v.StartDate.IsZero() {
		return NewAppError("Sale.IsValid", "model.voucher.is_valid.start_date.app_error", nil, "please provide valid start date", http.StatusBadRequest)
	}
	if v.EndDate != nil && v.EndDate.IsZero() {
		return NewAppError("Sale.IsValid", "model.voucher.is_valid.end_date.app_error", nil, "please provide valid end date", http.StatusBadRequest)
	}

	return nil
}

func (v *Voucher) commonPre() {
	if v.OnlyForStaff == nil {
		*v.OnlyForStaff = false
	}
	if v.Name != nil {
		*v.Name = SanitizeUnicode(*v.Name)
	}
	if v.DiscountValueType == "" {
		v.DiscountValueType = DISCOUNT_VALUE_TYPE_FIXED
	}
	if !v.Type.IsValid() {
		v.Type = VOUCHER_TYPE_ENTIRE_ORDER
	}
	if v.UsageLimit != nil && *v.UsageLimit < 0 {
		*v.UsageLimit = 0
	}
	v.Countries = strings.ToUpper(v.Countries)
	if v.Code == "" {
		v.Code = NewPromoCode()
	}
}

// VoucherTranslation represents translation for a voucher
type VoucherTranslation struct {
	Id           string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode"`
	Name         string           `json:"name" gorm:"type:varchar(255);column:Name"`
	VoucherID    string           `json:"voucher_id" gorm:"type:uuid;column:VoucherID"`
	CreateAt     int64            `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
}

func (c *VoucherTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *VoucherTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *VoucherTranslation) TableName() string             { return VoucherTranslationTableName }

// VoucherTranslationFilterOption is used to build squirrel queries
type VoucherTranslationFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (v *VoucherTranslation) commonPre() {
	v.Name = SanitizeUnicode(v.Name)
}

func (v *VoucherTranslation) IsValid() *AppError {
	if !IsValidId(v.VoucherID) {
		return NewAppError("VoucherTranslation.IsValid", "model.voucher_translation.is_valid.voucher_id.app_error", nil, "please provide valid voucher id", http.StatusBadRequest)
	}
	if !v.LanguageCode.IsValid() {
		return NewAppError("VoucherTranslation.IsValid", "model.voucher_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}
