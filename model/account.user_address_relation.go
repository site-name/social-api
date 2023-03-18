package model

import (
	"github.com/Masterminds/squirrel"
)

type UserAddress struct {
	Id        string `json:"id"`
	UserID    string `json:"user_id"`
	AddressID string `json:"address_id"`
}

type UserAddressFilterOptions struct {
	Id        squirrel.Sqlizer
	UserID    squirrel.Sqlizer
	AddressID squirrel.Sqlizer
}

func (ua *UserAddress) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.user_address.is_valid.%s.app_error",
		"user_address_id=",
		"UserAddress",
	)
	if !IsValidId(ua.Id) {
		return outer("id", nil)
	}
	if !IsValidId(ua.UserID) {
		return outer("user_id", &ua.Id)
	}
	if !IsValidId(ua.AddressID) {
		return outer("address_id", &ua.Id)
	}

	return nil
}

func (ua *UserAddress) PreSave() {
	if ua.Id == "" {
		ua.Id = NewId()
	}
}

func (ua *UserAddress) ToJSON() string {
	return ModelToJson(ua)
}
