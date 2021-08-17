package warehouse

import (
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
)

type Allocation struct {
	Id                string `json:"id"`
	CreateAt          int64  `json:"create_at"`
	OrderLineID       string `json:"order_ldine_id"`     // NOT NULL
	StockID           string `json:"stock_id"`           // NOT NULL
	QuantityAllocated int    `json:"quantity_allocated"` // default 0
}

// AllocationFilterOption is used to build sql queries to filtering warehouse allocations
type AllocationFilterOption struct {
	Id          *model.StringFilter
	OrderLineID *model.StringFilter
	StockID     *model.StringFilter
}

type AllocationError struct {
	OrderLineDatas []*order.OrderLine
	builder        strings.Builder
}

func (a *AllocationError) OrderLineIDs() string {
	a.builder.Reset()

	var suffix string = ", "
	for i, line := range a.OrderLineDatas {
		if i == len(a.OrderLineDatas)-1 {
			suffix = ""
		}
		a.builder.WriteString(line.Id + suffix)
	}

	return a.builder.String()
}

func (a *AllocationError) Error() string {
	a.builder.Reset()

	var suffix string = ", "
	for i, line := range a.OrderLineDatas {
		if i == len(a.OrderLineDatas)-1 {
			suffix = ""
		}
		a.builder.WriteString(line.String() + suffix)
	}

	return a.builder.String()
}

func (a *Allocation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.allocation.is_valid.%s.app_error",
		"allocation_id=",
		"Allocation.isValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if a.CreateAt == 0 {
		return outer("create_at", &a.Id)
	}
	if !model.IsValidId(a.OrderLineID) {
		return outer("order_line_id", &a.Id)
	}
	if !model.IsValidId(a.StockID) {
		return outer("stock_id", &a.Id)
	}

	return nil
}

func (a *Allocation) ToJson() string {
	return model.ModelToJson(a)
}

func (a *Allocation) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
	a.CreateAt = model.GetMillis()
	if a.QuantityAllocated < 0 {
		a.QuantityAllocated = 0
	}
}
