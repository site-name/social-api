package checkout

import (
	"io"

	"github.com/sitename/sitename/model"
)

const (
	CHECKOUT_LINE_MIN_QUANTITY = 1
)

type CheckoutLine struct {
	Id         string `json:"id"`
	CheckoutID string `json:"checkout_id"`
	VariantID  string `json:"variant_id"`
	Quantity   uint   `json:"quantity"`
}

func (c *CheckoutLine) ToJson() string {
	return model.ModelToJson(c)
}

func CheckoutLineFromJson(data io.Reader) *CheckoutLine {
	var checkoutLine CheckoutLine
	model.ModelFromJson(&checkoutLine, data)
	return &checkoutLine
}

func (c *CheckoutLine) Equal(other *CheckoutLine) bool {
	return c.VariantID == other.VariantID && c.Quantity == other.Quantity
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
	if c.CheckoutID == "" {
		return outer("checkout_id", &c.Id)
	}
	if c.VariantID == "" {
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
}
