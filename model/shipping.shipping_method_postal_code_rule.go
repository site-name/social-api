package model

import (
	"github.com/Masterminds/squirrel"
)

// max lengths for some fields
const (
	SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH = 32
)

type InclusionType string

func (i InclusionType) IsValid() bool {
	return PostalCodeRuleInclusionTypeString[i] != ""
}

// standard choices for inclusion_type
const (
	INCLUDE InclusionType = "include"
	EXCLUDE InclusionType = "exclude"
)

var PostalCodeRuleInclusionTypeString = map[InclusionType]string{
	INCLUDE: "Shipping method should include postal code rule",
	EXCLUDE: "Shipping method should exclude postal code rule",
}

type ShippingMethodPostalCodeRule struct {
	Id               string        `json:"id"`
	ShippingMethodID string        `json:"shipping_method_id"`
	Start            string        `json:"start"`
	End              string        `json:"end"`
	InclusionType    InclusionType `json:"inclusion_type"`
}

type ShippingMethodPostalCodeRules []*ShippingMethodPostalCodeRule

func (rs ShippingMethodPostalCodeRules) DeepCopy() ShippingMethodPostalCodeRules {
	res := make(ShippingMethodPostalCodeRules, len(rs))
	for idx, rule := range rs {
		res[idx] = rule.DeepCopy()
	}
	return res
}

type ShippingMethodPostalCodeRuleFilterOptions struct {
	Id               squirrel.Sqlizer
	ShippingMethodID squirrel.Sqlizer
}

func (r *ShippingMethodPostalCodeRule) DeepCopy() *ShippingMethodPostalCodeRule {
	if r == nil {
		return new(ShippingMethodPostalCodeRule)
	}

	res := *r
	return &res
}

func (s *ShippingMethodPostalCodeRule) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shipping_method_postal_code.is_valid.%s.app_error",
		"shipping_method_postal_code_id=",
		"ShippingMethodPostalCodeRule.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.ShippingMethodID) {
		return outer("shipping_method_is", &s.Id)
	}
	if len(s.Start) > SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH {
		return outer("start", &s.Id)
	}
	if len(s.End) > SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH {
		return outer("end", &s.Id)
	}
	if !s.InclusionType.IsValid() {
		return outer("inclusion_type", &s.Id)
	}

	return nil
}

func (s *ShippingMethodPostalCodeRule) ToJSON() string {
	return ModelToJson(s)
}
