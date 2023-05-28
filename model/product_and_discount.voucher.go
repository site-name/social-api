package model

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/language"
)

// Applicable values for Voucher type
const (
	SHIPPING         = "shipping"
	ENTIRE_ORDER     = "entire_order"
	SPECIFIC_PRODUCT = "specific_product"
)

var AllVoucherTypes = util.AnyArray[string]{SHIPPING, ENTIRE_ORDER, SPECIFIC_PRODUCT}

// Applicable values for voucher's discount value type
const (
	FIXED      = "fixed"
	PERCENTAGE = "percentage"
)

var SALE_TYPES = util.AnyArray[string]{FIXED, PERCENTAGE}

// max length values for some fields of voucher
const (
	VOUCHER_TYPE_MAX_LENGTH                = 20
	VOUCHER_NAME_MAX_LENGTH                = 255
	VOUCHER_CODE_MAX_LENGTH                = 16
	VOUCHER_DISCOUNT_VALUE_TYPE_MAX_LENGTH = 10
)

var (
	VOUCHER_DISCOUNT_COUNTRIES_MAX_LENGTH = len(Countries)*3 - 1 // hihi
)

type Voucher struct {
	Id                       string  `json:"id"`
	Type                     string  `json:"type"` // default to "entire_order"
	Name                     *string `json:"name"`
	Code                     string  `json:"code"` // UNIQUE
	UsageLimit               *int    `json:"usage_limit"`
	Used                     int     `json:"used"` // not editable
	StartDate                int64   `json:"start_date"`
	EndDate                  *int64  `json:"end_date"`
	ApplyOncePerOrder        bool    `json:"apply_once_per_order"`
	ApplyOncePerCustomer     bool    `json:"apply_once_per_customer"`
	OnlyForStaff             *bool   `json:"only_for_staff"` // default false
	DiscountValueType        string  `json:"discount_value_type"`
	Countries                string  `json:"countries"` // multiple. E.g: "VN US CN"
	MinCheckoutItemsQuantity int     `json:"min_checkout_items_quantity"`
	CreateAt                 int64   `json:"create_at"` // this field is for ordering
	UpdateAt                 int64   `json:"update_at"`
	ModelMetadata
}

// VoucherFilterOption
type VoucherFilterOption struct {
	Id                   squirrel.Sqlizer
	UsageLimit           squirrel.Sqlizer
	EndDate              squirrel.Sqlizer
	StartDate            squirrel.Sqlizer
	ChannelListingSlug   squirrel.Sqlizer // INNER JOIN ChannelListings ON ... INNER JOIN Channels ON ... WHERE Channels.Slug ...
	Code                 squirrel.Sqlizer
	ChannelListingActive *bool
	WithLook             bool // this add FOR UPDATE to sql queries
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

func (v *Voucher) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.voucher.is_valid.%s.app_error",
		"voucher_id=",
		"Voucher.IsValid",
	)
	if !IsValidId(v.Id) {
		return outer("id", nil)
	}

	if len(v.Type) > VOUCHER_TYPE_MAX_LENGTH || !AllVoucherTypes.Contains(v.Type) {
		return outer("type", &v.Id)
	}
	if v.Name != nil && len(*v.Name) > VOUCHER_NAME_MAX_LENGTH {
		return outer("name", &v.Id)
	}
	if len(v.Code) > VOUCHER_CODE_MAX_LENGTH {
		return outer("code", &v.Id)
	}
	if v.StartDate == 0 {
		return outer("start_date", &v.Id)
	}
	if len(v.DiscountValueType) > VOUCHER_DISCOUNT_VALUE_TYPE_MAX_LENGTH || !SALE_TYPES.Contains(v.DiscountValueType) {
		return outer("discount_value_type", &v.Id)
	}
	if len(v.Countries) > VOUCHER_DISCOUNT_COUNTRIES_MAX_LENGTH {
		return outer("countries", &v.Id)
	}
	for _, country := range strings.Fields(v.Countries) {
		if !CountryCode(country).IsValid() {
			return outer("countries", &v.Id)
		}
	}
	if v.CreateAt == 0 {
		return outer("create_at", &v.Id)
	}
	if v.UpdateAt == 0 {
		return outer("update_at", &v.Id)
	}

	return nil
}

func (v *Voucher) PreSave() {
	if v.Id == "" {
		v.Id = NewId()
	}
	v.CreateAt = GetMillis()

	v.UpdateAt = v.CreateAt
	if v.StartDate == 0 {
		v.StartDate = GetMillis()
	}
	v.commonPre()
}

func (v *Voucher) commonPre() {
	if v.OnlyForStaff == nil {
		v.OnlyForStaff = NewPrimitive(false)
	}
	if v.Name != nil {
		v.Name = NewPrimitive(SanitizeUnicode(*v.Name))
	}
	if v.DiscountValueType == "" {
		v.DiscountValueType = FIXED
	}
	if v.Type == "" {
		v.Type = ENTIRE_ORDER
	}
	if v.UsageLimit != nil && *v.UsageLimit < 0 {
		v.UsageLimit = NewPrimitive(0)
	}
	v.Countries = strings.ToUpper(v.Countries)
}

func (v *Voucher) PreUpdate() {
	v.UpdateAt = GetMillis()
	v.commonPre()
}

// VoucherTranslation represents translation for a voucher
type VoucherTranslation struct {
	Id           string `json:"id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	VoucherID    string `json:"voucher_id"`
	CreateAt     int64  `json:"create_at"`
}

// VoucherTranslationFilterOption is used to build squirrel queries
type VoucherTranslationFilterOption struct {
	Id           squirrel.Sqlizer
	LanguageCode squirrel.Sqlizer
	VoucherID    squirrel.Sqlizer
}

func (v *VoucherTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.voucher_translation.is_valid.%s.app_error",
		"voucher_trabslation_id=",
		"VoucherTranslation.IsValid",
	)
	if !IsValidId(v.Id) {
		return outer("id", nil)
	}
	if v.CreateAt == 0 {
		return outer("create_at", &v.Id)
	}
	if !IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if len(v.Name) > VOUCHER_NAME_MAX_LENGTH {
		return outer("name", &v.Id)
	}
	if tag, err := language.Parse(v.LanguageCode); err != nil || !strings.EqualFold(tag.String(), v.LanguageCode) {
		return outer("language_code", &v.Id)
	}

	return nil
}

func (v *VoucherTranslation) PreSave() {
	if v.Id == "" {
		v.Id = NewId()
	}
	v.CreateAt = GetMillis()
}
