package account

import (
	"github.com/sitename/sitename/model"
)

type UserAddress struct {
	Id        string `json:"id"`
	UserID    string `json:"user_id"`
	AddressID string `json:"address_id"`
}

func (ua *UserAddress) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.user_address.is_valid.%s.app_error",
		"user_address_id=",
		"model.UserAddress",
	)
	if !model.IsValidId(ua.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(ua.UserID) {
		return outer("user_id", &ua.Id)
	}
	if !model.IsValidId(ua.AddressID) {
		return outer("address_id", &ua.Id)
	}

	return nil
}

func (ua *UserAddress) PreSave() {
	if ua.Id == "" {
		ua.Id = model.NewId()
	}
}

func (ua *UserAddress) ToJSON() string {
	return model.ModelToJson(ua)
}
