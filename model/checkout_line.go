package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sitename/sitename/modules/json"
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
	b, _ := json.JSON.Marshal(c)
	return string(b)
}

func CheckoutLineFromJson(data io.Reader) *CheckoutLine {
	var checkoutLine CheckoutLine
	err := json.JSON.NewDecoder(data).Decode(&checkoutLine)
	if err != nil {
		return nil
	}
	return &checkoutLine
}

func (c *CheckoutLine) checkoutLineAppErr(field string) *AppError {
	var details string
	if strings.ToLower(field) != "id" {
		details += "checkout_id=" + c.Id
	}
	id := fmt.Sprintf("model.checkout_line.is_valid.%s.app_error", field)

	return NewAppError("CheckoutLine.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *CheckoutLine) Equal(other *CheckoutLine) bool {
	return c.VariantID == other.VariantID && c.Quantity == other.Quantity
}

func (c *CheckoutLine) IsValid() *AppError {
	if c.Id == "" {
		return c.checkoutLineAppErr("id")
	}
	if c.CheckoutID == "" {
		return c.checkoutLineAppErr("checkout_id")
	}
	if c.VariantID == "" {
		return c.checkoutLineAppErr("variant_id")
	}
	if c.Quantity < CHECKOUT_LINE_MIN_QUANTITY {
		return c.checkoutLineAppErr("quantity")
	}

	return nil
}

func (c *CheckoutLine) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
}
