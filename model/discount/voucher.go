package discount

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
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

type NotApplicable struct {
	MinSpent                 uint8
	MinCheckoutItemsQuantity uint8
}

type Voucher struct {
	Id                       string     `json:"id"`
	Type                     string     `json:"type"`
	Name                     string     `json:"name"`
	Code                     string     `json:"code"`
	UsageLimit               uint       `json:"usage_limit"`
	Used                     uint       `json:"used"`
	StartDate                *time.Time `json:"start_date"`
	EndDate                  *time.Time `json:"end_date"`
	ApplyOncePerOrder        bool       `json:"apply_once_per_order"`
	ApplyOncePerCustomer     bool       `json:"apply_once_per_customer"`
	DiscountValueType        string     `json:"discount_value_type"`
	Countries                string     `json:"countries"`
	MinCheckoutItemsQuantity uint       `json:"min_checkout_items_quantity"`
	// NOTE; there are oter fields
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
	if v.StartDate == nil {
		return v.createAppError("start_date")
	}
	if len(v.DiscountValueType) > VOUCHER_DISCOUNT_VALUE_TYPE_MAX_LENGTH || !model.StringArray([]string{FIXED, PERCENTAGE}).Contains(v.DiscountValueType) {
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
	if v.StartDate == nil {
		now := time.Now()
		v.StartDate = &now
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
	if v.StartDate == nil {
		now := time.Now()
		v.StartDate = &now
	}
	if v.DiscountValueType == "" {
		v.DiscountValueType = FIXED
	}
	v.Name = model.SanitizeUnicode(v.Name)
}
