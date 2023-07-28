package model

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/measurement"
	"gorm.io/gorm"
)

type ShippingMethodType string

func (s ShippingMethodType) IsValid() bool {
	return ShippingMethodTypeString[s] != ""
}

// shipping method valid types
const (
	PRICE_BASED  = "price"
	WEIGHT_BASED = "weight"
)

var ShippingMethodTypeString = map[ShippingMethodType]string{
	PRICE_BASED:  "Price based shipping",
	WEIGHT_BASED: "Weight based shipping",
}

type ShippingMethod struct {
	Id                  string                 `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name                string                 `json:"name" gorm:"type:varchar(100);column:Name"`
	Type                ShippingMethodType     `json:"type" gorm:"type:varchar(30);column:Type"`
	ShippingZoneID      string                 `json:"shipping_zone_id" gorm:"type:uuid;column:ShippingZoneID"`
	MinimumOrderWeight  *float32               `json:"minimum_order_weight" gorm:"column:MinimumOrderWeight;default:0"` // default 0
	MaximumOrderWeight  *float32               `json:"maximum_order_weight" gorm:"column:MaximumOrderWeight"`
	WeightUnit          measurement.WeightUnit `json:"weight_unit" gorm:"column:WeightUnit"`
	MaximumDeliveryDays *int                   `json:"maximum_delivery_days" gorm:"column:MaximumDeliveryDays"`
	MinimumDeliveryDays *int                   `json:"minimum_delivery_days" gorm:"column:MinimumDeliveryDays"`
	Description         StringInterface        `json:"description" gorm:"type:jsonb;column:Description"`
	ModelMetadata

	ExcludedProducts Products `json:"-" gorm:"many2many:ShippingMethodExcludedProducts"`

	MinOrderWeight *measurement.Weight `json:"min_order_weight" gorm:"-"`
	MaxOrderWeight *measurement.Weight `json:"max_order_weight" gorm:"-"`

	shippingZone                  *ShippingZone                 `gorm:"-"` // this field is used for holding prefetched related instances
	shippingMethodPostalCodeRules ShippingMethodPostalCodeRules `gorm:"-"` // this field is used for holding prefetched related instances
	price                         *goprices.Money               `gorm:"-"` // this field is populated in some graphql resolvers
}

func (c *ShippingMethod) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ShippingMethod) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ShippingMethod) TableName() string             { return ShippingMethodTableName }

// ShippingMethodFilterOption is used for filtering shipping methods
type ShippingMethodFilterOption struct {
	Conditions squirrel.Sqlizer

	// INNER JOIN ShippingZones ON ...
	//
	// INNER JOIN ShippingZoneChannels ON ...
	//
	// INNER JOIN Channels ON ...
	//
	// WHERE Channels.Slug ...
	ShippingZoneChannelSlug squirrel.Sqlizer
	// INNER JOIN ShippingMethodChannelListings ON ...
	//
	// INNER JOIN Channels ON ...
	//
	// WHERE Channels.Slug ...
	ChannelListingsChannelSlug squirrel.Sqlizer
	// INNER JOIN ShippingZones ON ...
	//
	// WHERE ShippingZones.Countries ...
	ShippingZoneCountries     squirrel.Sqlizer
	SelectRelatedShippingZone bool
}

func (s *ShippingMethod) GetPrice() *goprices.Money          { return s.price }
func (s *ShippingMethod) SetPrice(p *goprices.Money)         { s.price = p }
func (s *ShippingMethod) GetShippingZone() *ShippingZone     { return s.shippingZone }
func (s *ShippingMethod) SetShippingZone(zone *ShippingZone) { s.shippingZone = zone }

type ShippingMethods []*ShippingMethod

func (ss ShippingMethods) IDs() []string {
	return lo.Map(ss, func(item *ShippingMethod, _ int) string { return item.Id })
}

func (s *ShippingMethod) GetshippingMethodPostalCodeRules() ShippingMethodPostalCodeRules {
	return s.shippingMethodPostalCodeRules
}

func (s *ShippingMethod) SetshippingMethodPostalCodeRules(r ShippingMethodPostalCodeRules) {
	s.shippingMethodPostalCodeRules = r
}

func (s *ShippingMethod) AppendShippingMethodPostalCodeRule(rule *ShippingMethodPostalCodeRule) {
	s.shippingMethodPostalCodeRules = append(s.shippingMethodPostalCodeRules, rule)
}

func (s *ShippingMethod) PopulateNonDbFields() {
	if s.MinimumOrderWeight != nil {
		s.MinOrderWeight = &measurement.Weight{
			Amount: *s.MinimumOrderWeight,
			Unit:   s.WeightUnit,
		}
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

func (s *ShippingMethod) commonPre() {
	s.Name = SanitizeUnicode(s.Name)
	if s.MinimumOrderWeight == nil {
		s.MinimumOrderWeight = NewPrimitive[float32](0)
	}
}

func (s *ShippingMethod) IsValid() *AppError {
	if !IsValidId(s.ShippingZoneID) {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.shipping_zone_id.app_error", nil, "please provide valid shipping zone id", http.StatusBadRequest)
	}
	if !s.Type.IsValid() {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	return nil
}

type ShippingMethodTranslation struct {
	Id               string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	ShippingMethodID string           `json:"shipping_method_id" gorm:"type:uuid;column:ShippingMethodID"`
	LanguageCode     LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode"`
	Name             string           `json:"name" gorm:"type:varchar(255);column:Name"`
	Description      StringInterface  `json:"description" gorm:"type:jsonb;column:Description"`
}

func (c *ShippingMethodTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ShippingMethodTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ShippingMethodTranslation) TableName() string             { return ShippingMethodTranslationTableName }

func (s *ShippingMethodTranslation) IsValid() *AppError {
	if !IsValidId(s.ShippingMethodID) {
		return NewAppError("ShippingMethodTranslation.IsValid", "model.shipping_method_translation.is_valid.shipping_method_id.app_error", nil, "please provide valid shipping method id", http.StatusBadRequest)
	}
	if !s.LanguageCode.IsValid() {
		return NewAppError("ShippingMethodTranslation.IsValid", "model.shipping_method_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func (s *ShippingMethodTranslation) commonPre() {
	s.Name = SanitizeUnicode(s.Name)
}
