package order

import (
	"io"

	"github.com/sitename/sitename/model"
)

type FulfillmentLine struct {
	Id            string  `json:"id"`
	OrderLineID   string  `json:"order_line_id"`
	FulfillmentID string  `json:"fulfillment_id"`
	Quantity      uint    `json:"quantity"`
	StockID       *string `json:"stock_id"`
}

// FulfillmentLineFilterOption is used to build sql queries
type FulfillmentLineFilterOption struct {
	Id            *model.StringFilter
	OrderLineID   *model.StringFilter
	FulfillmentID *model.StringFilter
}

type FulfillmentLines []*FulfillmentLine

func (f FulfillmentLines) IDs() []string {
	res := make([]string, len(f))
	for i := range f {
		res[i] = f[i].Id
	}

	return res
}

func (f FulfillmentLines) OrderLineIDs() []string {
	res := make([]string, len(f))
	for i := range f {
		res[i] = f[i].OrderLineID
	}

	return res
}

func (f *FulfillmentLine) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.fulfillment_line.is_valid.%s.app_error",
		"fulfillment_line_id=",
		"FulfillmentLine.IsValid",
	)
	if !model.IsValidId(f.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(f.OrderLineID) {
		return outer("order_id", &f.Id)
	}
	if !model.IsValidId(f.FulfillmentID) {
		return outer("fulfillment_id", &f.Id)
	}
	if f.StockID != nil && !model.IsValidId(*f.StockID) {
		return outer("stock_id", &f.Id)
	}

	return nil
}

func (f *FulfillmentLine) ToJson() string {
	return model.ModelToJson(f)
}

func FulfillmentLineFromJson(data io.Reader) *FulfillmentLine {
	var f FulfillmentLine
	model.ModelFromJson(&f, data)
	return &f
}

func (f *FulfillmentLine) PreSave() {
	if f.Id == "" {
		f.Id = model.NewId()
	}
}
