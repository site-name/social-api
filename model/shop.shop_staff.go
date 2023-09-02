package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	"gorm.io/gorm"
)

type StaffSalaryPeriod string

const (
	StaffSalaryPeriodHourly  StaffSalaryPeriod = "hourly"
	StaffSalaryPeriodDaily   StaffSalaryPeriod = "daily"
	StaffSalaryPeriodMonthly StaffSalaryPeriod = "monthly"
)

func (s StaffSalaryPeriod) IsValid() bool {
	switch s {
	case StaffSalaryPeriodHourly, StaffSalaryPeriodDaily, StaffSalaryPeriodMonthly:
		return true
	default:
		return false
	}
}

// ShopStaff represents a relation between a shop and an staff user
type ShopStaff struct {
	Id             UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	StaffID        UUID              `json:"staff_id" gorm:"type:uuid;column:StaffID;unique"`
	CreateAt       int64             `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	EndAt          *int64            `json:"end_at" gorm:"type:bigint;column:EndAt"`
	SalaryPeriod   StaffSalaryPeriod `json:"salary_period" gorm:"type:varchar(20);column:SalaryPeriod"`
	Salary         decimal.Decimal   `json:"salary" gorm:"default:0;column:Salary"` // default 0
	SalaryCurrency string            `json:"salary_currency" gorm:"type:varchar(5);column:SalaryCurrency"`

	staff *User
}

func (r *ShopStaff) GetStaff() *User               { return r.staff }
func (r *ShopStaff) SetStaff(u *User)              { r.staff = u }
func (c *ShopStaff) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *ShopStaff) BeforeUpdate(_ *gorm.DB) error {
	c.CreateAt = 0 // prevent update
	return c.IsValid()
}
func (c *ShopStaff) TableName() string { return ShopStaffTableName }

type ShopStaffFilterOptions struct {
	Conditions squirrel.Sqlizer

	SelectRelatedStaff bool
}

func (s *ShopStaff) IsValid() *AppError {
	if !IsValidId(s.StaffID) {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff_relation.is_valid.staff_id.app_error", nil, "please provide valid shop staff id", http.StatusBadRequest)
	}
	if !s.SalaryPeriod.IsValid() {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff_relation.is_valid.salary_period.app_error", nil, "please provide valid salary period", http.StatusBadRequest)
	}
	if s.EndAt != nil && *s.EndAt == 0 {
		return NewAppError("ShopStaff.IsValid", "model.shop_staff_relation.is_valid.end_at.app_error", nil, "please provide valid end time", http.StatusBadRequest)
	}

	return nil
}

func (s *ShopStaff) DeepCopy() *ShopStaff {
	res := *s

	res.EndAt = CopyPointer(s.EndAt)
	if s.staff != nil {
		res.staff = s.staff.DeepCopy()
	}

	return &res
}
