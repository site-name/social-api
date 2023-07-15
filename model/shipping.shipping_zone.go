package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
)

// max length for some fields
const (
	SHIPPING_ZONE_NAME_MAX_LENGTH = 100
)

// order by CreateAt
type ShippingZone struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Countries   string `json:"countries"` // multiple allowed. E.g: VN USA CN
	Default     *bool  `json:"default"`   // default false
	Description string `json:"description"`
	CreateAt    int64  `json:"create_at"`
	ModelMetadata

	RelativeWarehouseIDs []string `json:"-" db:"-"`
}

// ShippingZoneFilterOption is used to build sql queries to finds shipping zones
type ShippingZoneFilterOption struct {
	Conditions squirrel.Sqlizer

	WarehouseID squirrel.Sqlizer // INNER JOIN WarehouseShippingZones ON ... WHERE WarehouseShippingZones.WarehouseID
	ChannelID   squirrel.Sqlizer // inner join shippingZoneChannel on ... WHERE shippingZoneChannel.ChannelID ...

	SelectRelatedWarehouseIDs bool // if true, `RelativeWarehouseIDs` property get populated with related data
}

type ShippingZones []*ShippingZone

func (s ShippingZones) IDs() []string {
	return lo.Map(s, func(i *ShippingZone, _ int) string { return i.Id })
}

func (s ShippingZones) DeepCopy() ShippingZones {
	res := make(ShippingZones, len(s))
	for idx, shippingZone := range s {
		res[idx] = shippingZone.DeepCopy()
	}
	return res
}

func (s *ShippingZone) String() string {
	return s.Name
}

func (s *ShippingZone) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shipping_zone.is_valid.%s.app_error",
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
		if !CountryCode(country).IsValid() {
			return outer("country", &s.Id)
		}
	}

	return nil
}

func (s *ShippingZone) PreSave() {
	if IsValidId(s.Id) {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
	s.commonPre()
}

func (s *ShippingZone) commonPre() {
	s.Name = SanitizeUnicode(s.Name)
	if s.Default == nil {
		s.Default = NewPrimitive(false)
	}
	s.Description = SanitizeUnicode(s.Description)
	s.Countries = strings.ToUpper(s.Countries)
}

func (s *ShippingZone) PreUpdate() {
	s.commonPre()
}

func (s *ShippingZone) DeepCopy() *ShippingZone {
	res := *s

	if s.Default != nil {
		res.Default = NewPrimitive(*s.Default)
	}
	res.ModelMetadata = s.ModelMetadata.DeepCopy()
	if len(s.RelativeWarehouseIDs) > 0 {
		copy(res.RelativeWarehouseIDs, s.RelativeWarehouseIDs)
	}
	return &res
}
