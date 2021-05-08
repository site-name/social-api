package shipping

import (
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"golang.org/x/text/language"
)

// max lengths for some fields
const (
	SHIPPING_METHOD_NAME_MAX_LENGTH = 100
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
	Id                  string                          `json:"id"`
	Name                string                          `json:"name"`
	Type                string                          `json:"type"`
	ShippingZoneID      string                          `json:"shipping_zone_id"`
	MinimumOrderWeight  *float32                        `json:"minimum_order_weight"`
	MaximumOrderWeight  *float32                        `json:"maximum_order_weight"`
	WeightUnit          string                          `json:"weight_unit"`
	ExcludedProducts    []*product_and_discount.Product `json:"excluded_products" db:"-"`
	MaximumDeliveryDays *uint                           `json:"maximum_delivery_days"`
	MinimumDeliveryDays *uint                           `json:"minimum_delivery_days"`
	Description         *model.StringInterface          `json:"description"`
	model.ModelMetadata `db:"-"`
}

func (s *ShippingMethod) String() string {
	return s.Name
}

func (s *ShippingMethod) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shipping_method.is_valid.%s.app_error",
		"shipping_method_id=",
		"ShippingMethod.IsValid",
	)

	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.ShippingZoneID) {
		return outer("shipping_zone_id", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SHIPPING_METHOD_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	if ShippingMethodTypeString[strings.ToLower(s.Type)] == "" {
		return outer("type", &s.Id)
	}
	return nil
}

func (s *ShippingMethod) ToJson() string {
	return model.ModelToJson(s)
}

// --------------------

const SHIPPING_METHOD_TRANSLATION_NAME_MAX_LENGTH = 255

type ShippingMethodTranslation struct {
	Id               string                 `json:"id"`
	ShippingMethodID string                 `json:"shipping_method_id"`
	LanguageCode     string                 `json:"language_code"`
	Name             string                 `json:"name"`
	Description      *model.StringInterface `json:"description"`
}

func (s *ShippingMethodTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shipping_method_translation.is_valid.%s.app_error",
		"shipping_method_translation_id=",
		"ShippingMethodTranslation.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.ShippingMethodID) {
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
		s.Id = model.NewId()
	}
	s.Name = model.SanitizeUnicode(s.Name)
}

func (s *ShippingMethodTranslation) PreUpdate() {
	s.Name = model.SanitizeUnicode(s.Name)
}
