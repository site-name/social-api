package product_and_discount

import (
	"fmt"
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/language"
)

// Applicable values for Voucher type
const (
	SHIPPING         = "shipping"
	ENTIRE_ORDER     = "entire_order"
	SPECIFIC_PRODUCT = "specific_product"
)

// Applicable values for voucher's discount value type
const (
	FIXED      = "fixed"
	PERCENTAGE = "percentage"
)

var SALE_TYPES = model.StringArray([]string{FIXED, PERCENTAGE})

// max length values for some fields of voucher
const (
	VOUCHER_TYPE_MAX_LENGTH                = 20
	VOUCHER_NAME_MAX_LENGTH                = 255
	VOUCHER_CODE_MAX_LENGTH                = 12
	VOUCHER_DISCOUNT_VALUE_TYPE_MAX_LENGTH = 10
)

var (
	VOUCHER_DISCOUNT_COUNTRIES_MAX_LENGTH = len(model.Countries)*3 - 1 // hihi
)

// NotApplicable represent error to some discount (vouchers || sales)
type NotApplicable struct {
	MinSpent                 *goprices.Money
	MinCheckoutItemsQuantity *uint
	Msg                      string
}

func NewNotApplicable(msg string, minSpent *goprices.Money, minCheckoutItemsQuantity *uint) *NotApplicable {
	return &NotApplicable{
		MinSpent:                 minSpent,
		Msg:                      msg,
		MinCheckoutItemsQuantity: minCheckoutItemsQuantity,
	}
}

// Error implements error interface
func (n *NotApplicable) Error() string {
	return n.Msg
}

type Voucher struct {
	Id                       string `json:"id"`
	ShopID                   string `json:"shop_id"` // the shop which issued this voucher
	Type                     string `json:"type"`
	Name                     string `json:"name"`
	Code                     string `json:"code"`
	UsageLimit               uint   `json:"usage_limit"`
	Used                     uint   `json:"used"` // not editable
	StartDate                int64  `json:"start_date"`
	EndDate                  *int64 `json:"end_date"`
	ApplyOncePerOrder        bool   `json:"apply_once_per_order"`
	ApplyOncePerCustomer     bool   `json:"apply_once_per_customer"`
	OnlyForStaff             *bool  `json:"only_for_staff"` // default false
	DiscountValueType        string `json:"discount_value_type"`
	Countries                string `json:"countries"` // multiple. E.g: "Vietnam America China"
	MinCheckoutItemsQuantity uint   `json:"min_checkout_items_quantity"`
	model.ModelMetadata
}

// VoucherValidateMinCheckoutItemsQuantity validates the quantity >= minimum requirement
func (voucher *Voucher) VoucherValidateMinCheckoutItemsQuantity(quantity uint) *model.AppError {
	if voucher.MinCheckoutItemsQuantity > quantity {
		return model.NewAppError("VoucherValidateMinCheckoutItemsQuantity", "app.discount.voucher_not_applicable_for_quantity_below", map[string]interface{}{"MinQuantity": voucher.MinCheckoutItemsQuantity}, "", http.StatusNotAcceptable)
	}
	return nil
}

// ValidateMinCheckoutItemsQuantity checks if given `quantity` satisfies min items quantity required
func (v *Voucher) ValidateMinCheckoutItemsQuantity(quantity uint) *NotApplicable {
	if v.MinCheckoutItemsQuantity > 0 && v.MinCheckoutItemsQuantity > quantity {
		return &NotApplicable{
			Msg:                      fmt.Sprintf("This offer is only valid for orders with minimum of %d items", v.MinCheckoutItemsQuantity),
			MinCheckoutItemsQuantity: &v.MinCheckoutItemsQuantity,
		}
	}

	return nil
}

func (v *Voucher) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher.is_valid.%s.app_error",
		"voucher_id=",
		"Voucher.IsValid",
	)
	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.ShopID) {
		return outer("shop_id", &v.Id)
	}
	if len(v.Type) > VOUCHER_TYPE_MAX_LENGTH || !model.StringArray([]string{SHIPPING, ENTIRE_ORDER, SPECIFIC_PRODUCT}).Contains(v.Type) {
		return outer("type", &v.Id)
	}
	if len(v.Name) > VOUCHER_NAME_MAX_LENGTH {
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
		if model.Countries[strings.ToUpper(country)] == "" { // does not exist in map
			return outer("countries", &v.Id)
		}
	}

	return nil
}

func (v *Voucher) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
	if v.Type == "" {
		v.Type = ENTIRE_ORDER
	}
	if v.StartDate == 0 {
		v.StartDate = model.GetMillis()
	}
	if v.DiscountValueType == "" {
		v.DiscountValueType = FIXED
	}
	if v.OnlyForStaff == nil {
		v.OnlyForStaff = model.NewBool(false)
	}
	v.Name = model.SanitizeUnicode(v.Name)
}

func (v *Voucher) PreUpdate() {
	if v.Type == "" {
		v.Type = ENTIRE_ORDER
	}
	if v.DiscountValueType == "" {
		v.DiscountValueType = FIXED
	}
	if v.OnlyForStaff == nil {
		v.OnlyForStaff = model.NewBool(false)
	}
	v.Name = model.SanitizeUnicode(v.Name)
}

// VoucherTranslation represents translation for a voucher
type VoucherTranslation struct {
	Id           string `json:"id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	VoucherID    string `json:"voucher_id"`
}

func (v *VoucherTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher_translation.is_valid.%s.app_error",
		"voucher_trabslation_id=",
		"VoucherTranslation.IsValid",
	)
	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.VoucherID) {
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
		v.Id = model.NewId()
	}
}
