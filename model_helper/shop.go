package model_helper

import (
	"net/http"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
)

type GiftCardSettingsExpiryType string

const (
	NEVER_EXPIRE  GiftCardSettingsExpiryType = "never_expire"
	EXPIRY_PERIOD GiftCardSettingsExpiryType = "expiry_period"
)

func (g GiftCardSettingsExpiryType) IsValid() bool {
	switch g {
	case NEVER_EXPIRE, EXPIRY_PERIOD:
		return true
	default:
		return false
	}
}

type ShopStaffFilterOptions struct {
	CommonQueryOptions
}

func ShopStaffPreSave(shopStaff *model.ShopStaff) {
	if shopStaff.ID == "" {
		shopStaff.ID = NewId()
	}
	if shopStaff.CreatedAt == 0 {
		shopStaff.CreatedAt = GetMillis()
	}
	ShopStaffCommonPre(shopStaff)
}

func ShopStaffCommonPre(shopStaff *model.ShopStaff) {
	if shopStaff.SalaryCurrency.IsValid() != nil {
		shopStaff.SalaryCurrency = DEFAULT_CURRENCY
	}
	if shopStaff.SalaryPeriod.IsValid() != nil {
		shopStaff.SalaryPeriod = model.StaffSalaryPeriodMonthly
	}
}

func ShopStaffIsValid(shopStaff model.ShopStaff) *AppError {
	if !IsValidId(shopStaff.ID) {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(shopStaff.StaffID) {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff.is_valid.staff_id.app_error", nil, "", http.StatusBadRequest)
	}
	if shopStaff.SalaryCurrency.IsValid() != nil {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff.is_valid.salary_currency.app_error", nil, "", http.StatusBadRequest)
	}
	if shopStaff.SalaryPeriod.IsValid() != nil {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff.is_valid.salary_period.app_error", nil, "", http.StatusBadRequest)
	}
	if shopStaff.CreatedAt <= 0 {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if shopStaff.Salary.LessThan(decimal.Zero) {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff.is_valid.salary.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}
