package warehouse

import (
	"io"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
)

// max lengths for some warehouse's fields
const (
	WAREHOUSE_NAME_MAX_LENGTH = 250
	WAREHOUSE_SLUG_MAX_LENGTH = 255
)

type WareHouse struct {
	Id        string `json:"id"`
	Name      string `json:"name"` // unique
	Slug      string `json:"slug"` // unique
	AddressID string `json:"address_id"`
	Email     string `json:"email"`
	model.ModelMetadata
}

func (w *WareHouse) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.warehouse.is_valid.%s.app_error",
		"warehouse_id=",
		"WareHouse.IsValid",
	)
	if !model.IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(w.AddressID) {
		return outer("address_id", nil)
	}
	if utf8.RuneCountInString(w.Name) > WAREHOUSE_NAME_MAX_LENGTH {
		return outer("name", &w.Id)
	}
	if len(w.Slug) > WAREHOUSE_SLUG_MAX_LENGTH {
		return outer("slug", &w.Id)
	}
	if w.Email != "" && !model.IsValidEmail(w.Email) {
		return outer("email", &w.Id)
	}

	return nil
}

func (w *WareHouse) PreSave() {
	if w.Id == "" {
		w.Id = model.NewId()
	}
	w.Name = model.SanitizeUnicode(w.Name)
	w.Slug = slug.Make(w.Name)
	w.ModelMetadata.PreSave()
}

func (w *WareHouse) PreUpdate() {
	w.Name = model.SanitizeUnicode(w.Name)
	w.Slug = slug.Make(w.Name)
	w.ModelMetadata.PreUpdate()
}

func (w *WareHouse) ToJson() string {
	return model.ModelToJson(w)
}

func WareHouseFromJson(data io.Reader) *WareHouse {
	var w WareHouse
	model.ModelFromJson(&w, data)
	return &w
}
