package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
)

type Allocation struct {
	Id                string `json:"id"`
	CreateAt          int64  `json:"create_at"`
	OrderLineID       string `json:"order_ldine_id"`     // NOT NULL
	StockID           string `json:"stock_id"`           // NOT NULL
	QuantityAllocated int    `json:"quantity_allocated"` // default 0

	stockAvailableQuantity int        // this field is set when AllocationFilterOption's `AnnotateStockAvailableQuantity` is true
	stock                  *Stock     // this field is populated with related stock
	orderLine              *OrderLine //
}

func (s *Allocation) SetStock(stk *Stock) {
	s.stock = stk
}

func (s *Allocation) GetStock() *Stock {
	return s.stock
}

func (s *Allocation) SetOrderLine(line *OrderLine) {
	s.orderLine = line
}

func (s *Allocation) GetOrderLine() *OrderLine {
	return s.orderLine
}

func (s *Allocation) SetStockAvailableQuantity(value int) {
	s.stockAvailableQuantity = value
}

func (s *Allocation) GetStockAvailableQuantity() int {
	return s.stockAvailableQuantity
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
	return lo.Map(a, func(al *Allocation, _ int) string { return al.Id })
}

func (a Allocations) StockIDs() []string {
	return lo.Map(a, func(al *Allocation, _ int) string { return al.StockID })
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

	if a.stock != nil {
		res.stock = a.stock.DeepCopy()
	}
	if a.orderLine != nil {
		res.orderLine = a.orderLine.DeepCopy()
	}
	return &res
}
