package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
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
	Id               UUID          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	ShippingMethodID UUID          `json:"shipping_method_id" gorm:"index:shippingmethodid_start_end_key;type:uuid;column:ShippingMethodID"`
	Start            string        `json:"start" gorm:"index:shippingmethodid_start_end_key;type:varchar(32);column:Start"`
	End              string        `json:"end" gorm:"index:shippingmethodid_start_end_key;type:varchar(32);column:End"`
	InclusionType    InclusionType `json:"inclusion_type" gorm:"type:varchar(32);column:InclusionType"`
}

func (c *ShippingMethodPostalCodeRule) BeforeCreate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}
func (c *ShippingMethodPostalCodeRule) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}
func (c *ShippingMethodPostalCodeRule) TableName() string {
	return ShippingMethodPostalCodeRuleTableName
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
	Conditions squirrel.Sqlizer
}

func (r *ShippingMethodPostalCodeRule) commonPre() {
	if !r.InclusionType.IsValid() {
		r.InclusionType = EXCLUDE
	}
}

func (r *ShippingMethodPostalCodeRule) DeepCopy() *ShippingMethodPostalCodeRule {
	if r == nil {
		return new(ShippingMethodPostalCodeRule)
	}

	res := *r
	return &res
}

func (s *ShippingMethodPostalCodeRule) IsValid() *AppError {
	if !IsValidId(s.ShippingMethodID) {
		return NewAppError("ShippingMethodPostalCodeRule.IsValid", "model.shipping_method_postal_code_rule.is_valid.shippingMethod_id.app_error", nil, "please provide valid shipping method id", http.StatusBadRequest)
	}
	if !s.InclusionType.IsValid() {
		return NewAppError("ShippingMethodPostalCodeRule.IsValid", "model.shipping_method_postal_code_rule.is_valid.inclusion_type.app_error", nil, "please provide valid inclusion type", http.StatusBadRequest)
	}

	return nil
}
