package model

import "github.com/Masterminds/squirrel"

// ShopStaffRelation represents a relation between a shop and an staff user
type ShopStaffRelation struct {
	Id       string `json:"id"`
	ShopID   string `json:"shop_id"`  // unique with staffID
	StaffID  string `json:"staff_id"` //
	CreateAt int64  `json:"create_at"`
	EndAt    *int64 `json:"end_at"`

	staff *User
	shop  *Shop
}

func (r *ShopStaffRelation) GetStaff() *User {
	return r.staff
}

func (r *ShopStaffRelation) SetStaff(u *User) {
	r.staff = u
}

func (r *ShopStaffRelation) GetShop() *Shop {
	return r.shop
}

func (r *ShopStaffRelation) SetShop(s *Shop) {
	r.shop = s
}

type ShopStaffRelationFilterOptions struct {
	ShopID   squirrel.Sqlizer
	StaffID  squirrel.Sqlizer
	CreateAt squirrel.Sqlizer
	EndAt    squirrel.Sqlizer

	SelectRelatedStaff bool
	SelectRelatedShop  bool
}

func (s *ShopStaffRelation) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
}

func (s *ShopStaffRelation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shop_staff_relation.is_valid.%s.app_error",
		"shop_staff_relation_id=",
		"ShopStaffRelation.IsValid",
	)

	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.ShopID) {
		return outer("shop_id", &s.Id)
	}
	if !IsValidId(s.StaffID) {
		return outer("staff_id", &s.Id)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if s.EndAt != nil && *s.EndAt == 0 {
		return outer("end_at", &s.Id)
	}

	return nil
}

func (s *ShopStaffRelation) DeepCopy() *ShopStaffRelation {
	res := *s

	if s.staff != nil {
		res.staff = s.staff.DeepCopy()
	}
	if s.shop != nil {
		res.shop = s.shop.DeepCopy()
	}

	return &res
}
