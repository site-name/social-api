package shipping

import (
	"strings"

	"github.com/sitename/sitename/model"
)

// max lengths for some fields
const (
	SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH = 32
)

// standard choices for inclusion_type
const (
	INCLUDE = "include"
	EXCLUDE = "exclude"
)

var PostalCodeRuleInclusionTypeString = map[string]string{
	INCLUDE: "Shipping method should include postal code rule",
	EXCLUDE: "Shipping method should exclude postal code rule",
}

type ShippingMethodPostalCodeRule struct {
	Id               string `json:"id"`
	ShippingMethodID string `json:"shipping_method_id"`
	Start            string `json:"start"`
	End              string `json:"end"`
	InclusionType    string `json:"inclusion_type"`
}

func (s *ShippingMethodPostalCodeRule) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shipping_method_postal_code.is_valid.%s.app_error",
		"shipping_method_postal_code_id=",
		"ShippingMethodPostalCodeRule.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.ShippingMethodID) {
		return outer("shipping_method_is", &s.Id)
	}
	if len(s.Start) > SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH {
		return outer("start", &s.Id)
	}
	if len(s.End) > SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH {
		return outer("end", &s.Id)
	}
	if PostalCodeRuleInclusionTypeString[strings.ToLower(s.InclusionType)] == "" {
		return outer("inclusion_type", &s.Id)
	}

	return nil
}

func (s *ShippingMethodPostalCodeRule) ToJson() string {
	return model.ModelToJson(s)
}
