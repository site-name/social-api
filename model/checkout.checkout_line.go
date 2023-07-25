package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// A single checkout line.
// Multiple lines in the same checkout can refer to the same product variant if
// their `data` field is different.
type CheckoutLine struct {
	Id         string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt   int64  `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	CheckoutID string `json:"checkout_id" gorm:"type:uuid;column:CheckoutID"`
	VariantID  string `json:"variant_id" gorm:"type:uuid;column:VariantID"`
	Quantity   int    `json:"quantity" gorm:"column:Quantity;check:Quantity >= 1"` // min 1
}

func (c *CheckoutLine) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *CheckoutLine) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *CheckoutLine) TableName() string             { return CheckoutLineTableName }

// CheckoutLineFilterOption is used to build squirrel sql queries
type CheckoutLineFilterOption struct {
	Conditions squirrel.Sqlizer
}

type CheckoutLines []*CheckoutLine

func (c CheckoutLines) VariantIDs() []string {
	return lo.Map(c, func(l *CheckoutLine, _ int) string { return l.VariantID })
}

func (c CheckoutLines) IDs() []string {
	return lo.Map(c, func(l *CheckoutLine, _ int) string { return l.Id })
}

func (c *CheckoutLine) Equal(other *CheckoutLine) bool {
	return c.VariantID == other.VariantID && c.Quantity == other.Quantity
}

func (c *CheckoutLine) NotEqual(other *CheckoutLine) bool {
	return !c.Equal(other)
}

func (c *CheckoutLine) IsValid() *AppError {
	if !IsValidId(c.CheckoutID) {
		return NewAppError("CheckoutLine.IsValid", "model.checkout_line.is_valid.checkout_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(c.VariantID) {
		return NewAppError("CheckoutLine.IsValid", "model.checkout_line.is_valid.variant_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (c *CheckoutLine) DeepCopy() *CheckoutLine {
	res := *c
	return &res
}
