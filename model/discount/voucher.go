package discount

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product"
	"github.com/sitename/sitename/modules/json"
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
	MinSpent                 *model.Money
	MinCheckoutItemsQuantity *uint
	Msg                      string
}

func NewNotApplicable(msg string, minSpent *model.Money, minCheckoutItemsQuantity *uint) *NotApplicable {
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
	Id                       string                `json:"id"`
	Type                     string                `json:"type"`
	Name                     string                `json:"name"`
	Code                     string                `json:"code"`
	UsageLimit               uint                  `json:"usage_limit"`
	Used                     uint                  `json:"used"`
	StartDate                int64                 `json:"start_date"`
	EndDate                  *int64                `json:"end_date"`
	ApplyOncePerOrder        bool                  `json:"apply_once_per_order"`
	ApplyOncePerCustomer     bool                  `json:"apply_once_per_customer"`
	DiscountValueType        string                `json:"discount_value_type"`
	Countries                string                `json:"countries"` // multiple. E.g: "Vietnam America China"
	MinCheckoutItemsQuantity uint                  `json:"min_checkout_items_quantity"`
	Products                 []*product.Product    `json:"products,omitempty" db:"-"`
	Collections              []*product.Collection `json:"collections,omitempty" db:"-"`
	Categories               []*product.Category   `json:"categories,omitempty" db:"-"`
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

func (v *Voucher) ToJson() string {
	b, _ := json.JSON.Marshal(v)
	return string(b)
}

func VoucherFromJson(data io.Reader) *Voucher {
	var v Voucher
	err := json.JSON.NewDecoder(data).Decode(&v)
	if err != nil {
		return nil
	}
	return &v
}

func (v *Voucher) createAppError(field string) *model.AppError {
	id := fmt.Sprintf("model.voucher.is_valid.%s.app_error", field)
	var details string
	if !strings.EqualFold(field, "id") {
		details = "voucher_id=" + v.Id
	}

	return model.NewAppError("Voucher.IsValid", id, nil, details, http.StatusBadRequest)
}

func (v *Voucher) IsValid() *model.AppError {
	if v.Id == "" {
		return v.createAppError("id")
	}
	if len(v.Type) > VOUCHER_TYPE_MAX_LENGTH || !model.StringArray([]string{SHIPPING, ENTIRE_ORDER, SPECIFIC_PRODUCT}).Contains(v.Type) {
		return v.createAppError("type")
	}
	if len(v.Name) > VOUCHER_NAME_MAX_LENGTH {
		return v.createAppError("name")
	}
	if len(v.Code) > VOUCHER_CODE_MAX_LENGTH {
		return v.createAppError("code")
	}
	if v.StartDate == 0 {
		return v.createAppError("start_date")
	}
	if len(v.DiscountValueType) > VOUCHER_DISCOUNT_VALUE_TYPE_MAX_LENGTH || !SALE_TYPES.Contains(v.DiscountValueType) {
		return v.createAppError("discount_value_type")
	}
	if len(v.Countries) > VOUCHER_DISCOUNT_COUNTRIES_MAX_LENGTH {
		return v.createAppError("countries")
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
	v.Name = model.SanitizeUnicode(v.Name)
}

func (v *Voucher) PreUpdate() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
	if v.Type == "" {
		v.Type = ENTIRE_ORDER
	}
	if v.DiscountValueType == "" {
		v.DiscountValueType = FIXED
	}
	v.Name = model.SanitizeUnicode(v.Name)
}

// ------------------------------
type VoucherTranslation struct {
	Id           string `json:"id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	VoucherID    string `json:"voucher_id"`
}

func (v *VoucherTranslation) ToJson() string {
	b, _ := json.JSON.Marshal(v)
	return string(b)
}

func VoucherTranslationFromJson(data io.Reader) *VoucherTranslation {
	var vt VoucherTranslation
	err := json.JSON.NewDecoder(data).Decode(&vt)
	if err != nil {
		return nil
	}
	return &vt
}

func (v *VoucherTranslation) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.voucher_translation.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "voucher_trabslation_id=" + v.Id
	}

	return model.NewAppError("VoucherTranslation.IsValid", id, nil, details, http.StatusBadRequest)
}

func (v *VoucherTranslation) IsValid() *model.AppError {
	if v.Id == "" {
		return v.createAppError("id")
	}
	if v.VoucherID == "" {
		return v.createAppError("voucher_id")
	}
	if len(v.Name) > VOUCHER_NAME_MAX_LENGTH {
		return v.createAppError("name")
	}
	if tag, err := language.Parse(v.LanguageCode); err != nil || !strings.EqualFold(tag.String(), v.LanguageCode) {
		return v.createAppError("language_code")
	}

	return nil
}

func (v *VoucherTranslation) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
}
