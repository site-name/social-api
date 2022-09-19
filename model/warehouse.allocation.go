package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

type Allocation struct {
	Id                string `json:"id"`
	CreateAt          int64  `json:"create_at"`
	OrderLineID       string `json:"order_ldine_id"`     // NOT NULL
	StockID           string `json:"stock_id"`           // NOT NULL
	QuantityAllocated int    `json:"quantity_allocated"` // default 0

	StockAvailableQuantity int        `json:"-" db:"-"` // this field is set when AllocationFilterOption's `AnnotateStockAvailableQuantity` is true
	Stock                  *Stock     `json:"-" db:"-"` // this field is populated with related stock
	OrderLine              *OrderLine `json:"-" db:"-"`
}

// AllocationFilterOption is used to build sql queries to filtering warehouse allocations
type AllocationFilterOption struct {
	Id                squirrel.Sqlizer
	OrderLineID       squirrel.Sqlizer
	OrderLineOrderID  squirrel.Sqlizer // INNER JOIN OrderLines ON (...) WHERE OrderLines.OrderID = ...
	StockID           squirrel.Sqlizer
	QuantityAllocated squirrel.Sqlizer

	LockForUpdate bool   // if true, `FOR UPDATE` will be placed in the end of sqlqueries
	ForUpdateOf   string // this is placed after `FOR UPDATE`. E.g: "Warehouses" => `FOR UPDATE OF Warehouses`

	SelectedRelatedStock   bool
	SelectRelatedOrderLine bool

	AnnotateStockAvailableQuantity bool
}

type Allocations []*Allocation

func (a Allocations) IDs() []string {
	res := []string{}
	for _, item := range a {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (a Allocations) StockIDs() []string {
	res := []string{}
	for _, item := range a {
		if item != nil {
			res = append(res, item.StockID)
		}
	}

	return res
}

func (a Allocations) Len() int {
	return len(a)
}

func (a *Allocation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.allocation.is_valid.%s.app_error",
		"allocation_id=",
		"Allocation.isValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if a.CreateAt == 0 {
		return outer("create_at", &a.Id)
	}
	if !IsValidId(a.OrderLineID) {
		return outer("order_line_id", &a.Id)
	}
	if !IsValidId(a.StockID) {
		return outer("stock_id", &a.Id)
	}

	return nil
}

func (a *Allocation) ToJSON() string {
	return ModelToJson(a)
}

func (a *Allocation) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	a.CreateAt = GetMillis()
	a.commonPre()
}

func (a *Allocation) commonPre() {
	if a.QuantityAllocated < 0 {
		a.QuantityAllocated = 0
	}
}

func (a *Allocation) PreUpdate() {
	a.commonPre()
}

type AllocationError struct {
	OrderLines OrderLines
	builder    strings.Builder
}

func (a *AllocationError) Error() string {
	a.builder.Reset()

	a.builder.WriteString("Unable to deallocate stock for lines ")

	var suffix string = ", "
	for i, line := range a.OrderLines {
		if i == len(a.OrderLines)-1 {
			suffix = ""
		}
		a.builder.WriteString(line.String() + suffix)
	}

	return a.builder.String()
}

func (a *Allocation) DeepCopy() *Allocation {
	res := *a

	if a.Stock != nil {
		res.Stock = a.Stock.DeepCopy()
	}
	if a.OrderLine != nil {
		res.OrderLine = a.OrderLine.DeepCopy()
	}
	return &res
}
