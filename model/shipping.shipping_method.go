package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/modules/measurement"
	"golang.org/x/text/language"
)

// max lengths for some fields
const (
	SHIPPING_METHOD_NAME_MAX_LENGTH = 100
	SHIPPING_METHOD_TYPE_MAX_LENGTH = 30
)

// shipping method valid types
const (
	PRICE_BASED  = "price"
	WEIGHT_BASED = "weight"
)

var ShippingMethodTypeString = map[string]string{
	PRICE_BASED:  "Price based shipping",
	WEIGHT_BASED: "Weight based shipping",
}

type ShippingMethod struct {
	Id                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Type                string                 `json:"type"`
	ShippingZoneID      string                 `json:"shipping_zone_id"`
	MinimumOrderWeight  float32                `json:"minimum_order_weight"` // default0 0
	MaximumOrderWeight  *float32               `json:"maximum_order_weight"`
	WeightUnit          measurement.WeightUnit `json:"weight_unit"`
	MinOrderWeight      *measurement.Weight    `json:"min_order_weight" db:"-"`
	MaxOrderWeight      *measurement.Weight    `json:"max_order_weight" db:"-"`
	MaximumDeliveryDays *int                   `json:"maximum_delivery_days"`
	MinimumDeliveryDays *int                   `json:"minimum_delivery_days"`
	Description         StringInterface        `json:"description"`
	ModelMetadata

	shippingZones                 ShippingZones                 `db:"-"` // this field is used for holding prefetched related instances
	shippingMethodPostalCodeRules ShippingMethodPostalCodeRules `db:"-"` // this field is used for holding prefetched related instances
}

func (s *ShippingMethod) GetShippingZones() ShippingZones {
	return s.shippingZones
}

func (s *ShippingMethod) SetShippingZones(zones ShippingZones) {
	s.shippingZones = zones
}

func (s *ShippingMethod) AppendShippingZone(zone *ShippingZone) {
	if zone != nil {
		s.shippingZones = append(s.shippingZones, zone)
	}
}

func (s *ShippingMethod) GetshippingMethodPostalCodeRules() ShippingMethodPostalCodeRules {
	return s.shippingMethodPostalCodeRules
}

func (s *ShippingMethod) SetshippingMethodPostalCodeRules(r ShippingMethodPostalCodeRules) {
	s.shippingMethodPostalCodeRules = r
}

func (s *ShippingMethod) AppendShippingMethodPostalCodeRule(rule *ShippingMethodPostalCodeRule) {
	if rule != nil {
		s.shippingMethodPostalCodeRules = append(s.shippingMethodPostalCodeRules, rule)
	}
}

// ShippingMethodFilterOption is used for filtering shipping methods
type ShippingMethodFilterOption struct {
	Id                         squirrel.Sqlizer
	Type                       squirrel.Sqlizer
	MinimumOrderWeight         squirrel.Sqlizer
	MaximumOrderWeight         squirrel.Sqlizer
	ShippingZoneChannelSlug    squirrel.Sqlizer
	ChannelListingsChannelSlug squirrel.Sqlizer
	ShippingZoneID             squirrel.Sqlizer
}

func (s *ShippingMethod) PopulateNonDbFields() {
	s.MinOrderWeight = &measurement.Weight{
		Amount: s.MinimumOrderWeight,
		Unit:   s.WeightUnit,
	}
	if s.MaximumOrderWeight != nil {
		s.MaxOrderWeight = &measurement.Weight{
			Amount: *s.MaximumOrderWeight,
			Unit:   s.WeightUnit,
		}
	}
}

func (s *ShippingMethod) String() string {
	if s.Type == PRICE_BASED {
		return fmt.Sprintf("ShippingMethod(type=%s)", s.Type)
	}

	return fmt.Sprintf("ShippingMethod(type=%s weight_range=(%s))", s.Type, s.getWeightTypeDisplay())
}

func (s *ShippingMethod) getWeightTypeDisplay() string {
	s.PopulateNonDbFields()

	if s.MinOrderWeight.Unit != measurement.STANDARD_WEIGHT_UNIT {
		minWeight, _ := s.MinOrderWeight.ConvertTo(measurement.STANDARD_WEIGHT_UNIT)
		s.MinOrderWeight = minWeight
	}
	if s.MaxOrderWeight == nil {
		return fmt.Sprintf("%s and up", s.MinOrderWeight.String())
	}

	if s.MaxOrderWeight != nil && s.MaxOrderWeight.Unit != measurement.STANDARD_WEIGHT_UNIT {
		maxWeight, _ := s.MinOrderWeight.ConvertTo(measurement.STANDARD_WEIGHT_UNIT)
		s.MaxOrderWeight = maxWeight
	}
	return fmt.Sprintf("%s to %s", s.MinOrderWeight.String(), s.MaxOrderWeight.String())
}

func (s *ShippingMethod) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.Name = SanitizeUnicode(s.Name)
}

func (s *ShippingMethod) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"shipping_method.is_valid.%s.app_error",
		"shipping_method_id=",
		"ShippingMethod.IsValid",
	)

	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.ShippingZoneID) {
		return outer("shipping_zone_id", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SHIPPING_METHOD_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	if ShippingMethodTypeString[strings.ToLower(s.Type)] == "" || len(s.Type) > SHIPPING_METHOD_TYPE_MAX_LENGTH {
		return outer("type", &s.Id)
	}
	return nil
}

func (s *ShippingMethod) DeepCopy() *ShippingMethod {
	if s == nil {
		return new(ShippingMethod)
	}

	res := *s

	if s.MaximumOrderWeight != nil {
		res.MaximumOrderWeight = NewPrimitive(*s.MaximumOrderWeight)
	}

	if s.MaximumDeliveryDays != nil {
		res.MaximumDeliveryDays = NewPrimitive(*s.MaximumDeliveryDays)
	}
	if s.MinimumDeliveryDays != nil {
		res.MinimumDeliveryDays = NewPrimitive(*s.MinimumDeliveryDays)
	}
	if s.MinOrderWeight != nil {
		res.MinOrderWeight = &measurement.Weight{Amount: s.MinOrderWeight.Amount, Unit: s.MinOrderWeight.Unit}
	}
	if s.MaxOrderWeight != nil {
		res.MaxOrderWeight = &measurement.Weight{Amount: s.MaxOrderWeight.Amount, Unit: s.MaxOrderWeight.Unit}
	}
	res.Description = s.Description.DeepCopy()
	res.ModelMetadata = s.ModelMetadata.DeepCopy()

	res.shippingZones = s.shippingZones.DeepCopy()
	res.shippingMethodPostalCodeRules = s.shippingMethodPostalCodeRules.DeepCopy()

	return &res
}

const SHIPPING_METHOD_TRANSLATION_NAME_MAX_LENGTH = 255

type ShippingMethodTranslation struct {
	Id               string           `json:"id"`
	ShippingMethodID string           `json:"shipping_method_id"`
	LanguageCode     string           `json:"language_code"`
	Name             string           `json:"name"`
	Description      *StringInterface `json:"description"`
}

func (s *ShippingMethodTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"shipping_method_translation.is_valid.%s.app_error",
		"shipping_method_translation_id=",
		"ShippingMethodTranslation.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.ShippingMethodID) {
		return outer("shipping_method_id", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SHIPPING_METHOD_TRANSLATION_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	if tag, err := language.Parse(s.LanguageCode); err != nil || !strings.EqualFold(tag.String(), s.LanguageCode) {
		return outer("language_code", &s.Id)
	}

	return nil
}

func (s *ShippingMethodTranslation) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.Name = SanitizeUnicode(s.Name)
}

func (s *ShippingMethodTranslation) PreUpdate() {
	s.Name = SanitizeUnicode(s.Name)
}
