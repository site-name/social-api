package model

// ShopStaffRelation represents a relation between a shop and an staff user
type ShopStaffRelation struct {
	Id       string `json:"id"`
	ShopID   string `json:"shop_id"`  // unique with staffID
	StaffID  string `json:"staff_id"` //
	CreateAt int64  `json:"create_at"`
	EndAt    *int64 `json:"end_at"`
}

func (s *ShopStaffRelation) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
}

func (s *ShopStaffRelation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"shop_staff_relation.is_valid.%s.app_error",
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