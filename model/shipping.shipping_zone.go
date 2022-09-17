package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/modules/util"
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
	ModelMetadata

	RelativeWarehouseIDs []string `json:"-" db:"-"`
}

// ShippingZoneFilterOption is used to build sql queries to finds shipping zones
type ShippingZoneFilterOption struct {
	Id           squirrel.Sqlizer
	DefaultValue *bool
	WarehouseID  squirrel.Sqlizer // INNER JOIN WarehouseShippingZones ON ... WHERE WarehouseShippingZones.WarehouseID

	SelectRelatedThroughData bool // if true, `RelativeWarehouseIDs` property get populated with related data
}

type ShippingZones []*ShippingZone

func (s ShippingZones) IDs() []string {
	var res []string
	for _, zone := range s {
		if zone != nil {
			res = append(res, zone.Id)
		}
	}

	return res
}

// RelativeWarehouseIDsFlat joins all `RelativeWarehouseIDs` fields of all shipping zones into single slice of strings
//
// E.g: [["a", "b"], ["c", "d"]] => ["a", "b", "c", "d"]
func (s ShippingZones) RelativeWarehouseIDsFlat(keepDuplicates bool) []string {
	var res []string

	for _, item := range s {
		res = append(res, item.RelativeWarehouseIDs...)
	}

	if keepDuplicates {
		return res
	}
	return util.Dedup(res)
}

func (s *ShippingZone) String() string {
	return s.Name
}

func (s *ShippingZone) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"shipping_zone.is_valid.%s.app_error",
		"shipping_zone_id=",
		"ShippingZone.IsValid",
	)

	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SHIPPING_ZONE_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	for _, country := range strings.Fields(s.Countries) {
		if _, ok := Countries[country]; !ok {
			return outer("country", &s.Id)
		}
	}

	return nil
}

func (s *ShippingZone) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
	s.Name = SanitizeUnicode(s.Name)
	if s.Default == nil {
		s.Default = NewBool(false)
	}
	s.Description = SanitizeUnicode(s.Description)
	s.Countries = strings.ToUpper(s.Countries)
}

func (s *ShippingZone) PreUpdate() {
	s.Name = SanitizeUnicode(s.Name)
	s.Description = SanitizeUnicode(s.Description)
	s.Countries = strings.ToUpper(s.Countries)
	if s.Default == nil {
		s.Default = NewBool(false)
	}
}

func (s *ShippingZone) DeepCopy() *ShippingZone {
	res := *s

	res.ModelMetadata = *s.ModelMetadata.DeepCopy()
	return &res
}
