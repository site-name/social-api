package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
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
	Id             string            `json:"id"`
	StaffID        string            `json:"staff_id"` //
	CreateAt       int64             `json:"create_at"`
	EndAt          *int64            `json:"end_at"`
	SalaryPeriod   StaffSalaryPeriod `json:"salary_period"`
	Salary         decimal.Decimal   `json:"salary"` // default 0
	SalaryCurrency string            `json:"salary_currency"`

	staff *User
}

func (r *ShopStaff) GetStaff() *User {
	return r.staff
}

func (r *ShopStaff) SetStaff(u *User) {
	r.staff = u
}

type ShopStaffFilterOptions struct {
	Conditions squirrel.Sqlizer

	SelectRelatedStaff bool
}

func (s *ShopStaff) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
}

func (s *ShopStaff) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shop_staff_relation.is_valid.%s.app_error",
		"shop_staff_relation_id=",
		"ShopStaff.IsValid",
	)

	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.StaffID) {
		return outer("staff_id", &s.Id)
	}
	if !s.SalaryPeriod.IsValid() {
		return outer("salary_period", &s.Id)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if s.EndAt != nil && *s.EndAt == 0 {
		return outer("end_at", &s.Id)
	}

	return nil
}

func (s *ShopStaff) DeepCopy() *ShopStaff {
	res := *s

	if s.staff != nil {
		res.staff = s.staff.DeepCopy()
	}

	return &res
}
