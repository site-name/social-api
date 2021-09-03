package checkout

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
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
	CreatAt    int64  `json:"create_at"`
	CheckoutID string `json:"checkout_id"`
	VariantID  string `json:"variant_id"`
	Quantity   int    `json:"quantity"`

	ProductVariant *product_and_discount.ProductVariant `json:"-" db:"-"`
}

type CheckoutLines []*CheckoutLine

func (c CheckoutLines) VariantIDs() []string {
	res := []string{}
	for _, item := range c {
		if item != nil {
			res = append(res, item.VariantID)
		}
	}

	return res
}

func (c CheckoutLines) IDs() []string {
	res := []string{}
	for _, item := range c {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (c *CheckoutLine) ToJson() string {
	return model.ModelToJson(c)
}

func (c *CheckoutLine) Equal(other *CheckoutLine) bool {
	return c.VariantID == other.VariantID && c.Quantity == other.Quantity
}

func (c *CheckoutLine) NotEqual(other *CheckoutLine) bool {
	return !c.Equal(other)
}

func (c *CheckoutLine) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.checkout_line.is_valid.%s.app_error",
		"checkout_id=",
		"CheckoutLine.IsValid",
	)
	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.CreatAt == 0 {
		return outer("create_at", &c.Id)
	}
	if !model.IsValidId(c.CheckoutID) {
		return outer("checkout_id", &c.Id)
	}
	if !model.IsValidId(c.VariantID) {
		return outer("variant_id", &c.Id)
	}
	if c.Quantity < CHECKOUT_LINE_MIN_QUANTITY {
		return outer("quantity", &c.Id)
	}

	return nil
}

func (c *CheckoutLine) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	if c.CreatAt == 0 {
		c.CreatAt = model.GetMillis()
	}
}
