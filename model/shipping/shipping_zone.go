package shipping

import (
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
)

// max length for some fields
const (
	SHIPPING_ZONE_NAME_MAX_LENGTH = 100
)

type ShippingZone struct {
	Id                  string             `json:"id"`
	Name                string             `json:"name"`
	Contries            string             `json:"countries"` // multiple allowed
	Default             *bool              `json:"default"`
	Description         string             `json:"description"`
	Channels            []*channel.Channel `json:"channels"`
	model.ModelMetadata `db:"-"`
}

func (s *ShippingZone) String() string {
	return s.Name
}

func (s *ShippingZone) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shipping_zone.is_valid.%s.app_error",
		"shipping_zone_id=",
		"ShippingZone.IsValid",
	)

	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(s.Name) > SHIPPING_ZONE_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	for _, country := range strings.Fields(s.Contries) {
		if model.Countries[strings.ToUpper(country)] == "" {
			return outer("country", &s.Id)
		}
	}

	return nil
}

func (s *ShippingZone) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	s.Name = model.SanitizeUnicode(s.Name)
	if s.Default == nil {
		s.Default = model.NewBool(false)
	}
}

func (s *ShippingZone) PreUpdate() {
	s.Name = model.SanitizeUnicode(s.Name)
}

func (s *ShippingZone) ToJson() string {
	return model.ModelToJson(s)
}