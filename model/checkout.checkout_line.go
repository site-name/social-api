package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
)

// max lengths for checkout line table's fields
const (
	CHECKOUT_LINE_MIN_QUANTITY = 1
)

// A single checkout line.
// Multiple lines in the same checkout can refer to the same product variant if
// their `data` field is different.
type CheckoutLine struct {
	Id         string `json:"id"`
	CreateAt   int64  `json:"create_at"`
	CheckoutID string `json:"checkout_id"`
	VariantID  string `json:"variant_id"`
	Quantity   int    `json:"quantity"`

	productVariant *ProductVariant `db:"-"`
}

func (c *CheckoutLine) GetProductVariant() *ProductVariant {
	return c.productVariant
}

func (c *CheckoutLine) SetProductVariant(p *ProductVariant) {
	c.productVariant = p
}

// CheckoutLineFilterOption is used to build squirrel sql queries
type CheckoutLineFilterOption struct {
	Id         squirrel.Sqlizer
	CheckoutID squirrel.Sqlizer
	VariantID  squirrel.Sqlizer
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
	outer := CreateAppErrorForModel(
		"checkout_line.is_valid.%s.app_error",
		"checkout_id=",
		"CheckoutLine.IsValid",
	)
	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.CreateAt == 0 {
		return outer("create_at", &c.Id)
	}
	if !IsValidId(c.CheckoutID) {
		return outer("checkout_id", &c.Id)
	}
	if !IsValidId(c.VariantID) {
		return outer("variant_id", &c.Id)
	}
	if c.Quantity < CHECKOUT_LINE_MIN_QUANTITY {
		return outer("quantity", &c.Id)
	}

	return nil
}

func (c *CheckoutLine) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	if c.CreateAt == 0 {
		c.CreateAt = GetMillis()
	}
}

func (c *CheckoutLine) DeepCopy() *CheckoutLine {
	res := *c
	if c.productVariant != nil {
		res.productVariant = c.productVariant.DeepCopy()
	}
	return &res
}
