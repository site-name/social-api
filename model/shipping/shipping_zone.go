package shipping

import (
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
)

// max length for some fields
const (
	SHIPPING_ZONE_NAME_MAX_LENGTH = 100
)

type ShippingZone struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Countries   string `json:"countries"` // multiple allowed, a, b, b, ...
	Default     *bool  `json:"default"`   // default false
	Description string `json:"description"`
	CreateAt    int64  `json:"create_at"`
	model.ModelMetadata
}

// ShippingZoneFilterOption is used to build sql queries to finds shipping zones
type ShippingZoneFilterOption struct {
	Id           *model.StringFilter // filter on Id field
	DefaultValue *bool               // filter on Default field
	WarehouseIDs []string
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
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SHIPPING_ZONE_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	for _, country := range strings.Fields(s.Countries) {
		if _, ok := model.Countries[country]; !ok {
			return outer("country", &s.Id)
		}
	}

	return nil
}

func (s *ShippingZone) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	s.CreateAt = model.GetMillis()
	s.Name = model.SanitizeUnicode(s.Name)
	if s.Default == nil {
		s.Default = model.NewBool(false)
	}
	s.Description = model.SanitizeUnicode(s.Description)
	s.Countries = strings.ToUpper(s.Countries)
}

func (s *ShippingZone) PreUpdate() {
	s.Name = model.SanitizeUnicode(s.Name)
	s.Description = model.SanitizeUnicode(s.Description)
	s.Countries = strings.ToUpper(s.Countries)
	if s.Default == nil {
		s.Default = model.NewBool(false)
	}
}
