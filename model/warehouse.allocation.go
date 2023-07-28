package model

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Allocation struct {
	Id                string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt          int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	OrderLineID       string `json:"order_line_id" gorm:"type:uuid;column:OrderLineID;index:orderlineid_stockid_key"` // NOT NULL
	StockID           string `json:"stock_id" gorm:"type:uuid;column:StockID;index:orderlineid_stockid_key"`          // NOT NULL
	QuantityAllocated int    `json:"quantity_allocated" gorm:"column:QuantityAllocated"`                              // default 0

	stockAvailableQuantity int        // this field is set when AllocationFilterOption's `AnnotateStockAvailableQuantity` is true
	stock                  *Stock     // this field is populated with related stock
	orderLine              *OrderLine //
}

func (s *Allocation) SetStock(stk *Stock)                 { s.stock = stk }
func (s *Allocation) GetStock() *Stock                    { return s.stock }
func (s *Allocation) SetOrderLine(line *OrderLine)        { s.orderLine = line }
func (s *Allocation) GetOrderLine() *OrderLine            { return s.orderLine }
func (s *Allocation) SetStockAvailableQuantity(value int) { s.stockAvailableQuantity = value }
func (s *Allocation) GetStockAvailableQuantity() int      { return s.stockAvailableQuantity }

func (c *Allocation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Allocation) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent updating
	return c.IsValid()
}
func (c *Allocation) TableName() string { return AllocationTableName }

// AllocationFilterOption is used to build sql queries to filtering warehouse allocations
type AllocationFilterOption struct {
	Conditions squirrel.Sqlizer

	OrderLineOrderID squirrel.Sqlizer // INNER JOIN OrderLines ON (...) WHERE OrderLines.OrderID ...

	// if true, `FOR UPDATE` will be placed in the end of sqlqueries.
	// NOTE: Only apply if `Transaction` is set
	LockForUpdate bool
	ForUpdateOf   string // this is placed after `FOR UPDATE`. E.g: "Warehouses" => `FOR UPDATE OF Warehouses`
	Transaction   *gorm.DB

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
	if !IsValidId(a.OrderLineID) {
		return NewAppError("Allocation.IsValid", "model.allocation.is_valid.orderline_id.app_error", nil, "please provide valid order line id", http.StatusBadRequest)
	}
	if !IsValidId(a.StockID) {
		return NewAppError("Allocation.IsValid", "model.allocation.is_valid.stock_id.app_error", nil, "please provide valid stock id", http.StatusBadRequest)
	}

	return nil
}

func (a *Allocation) commonPre() {
	if a.QuantityAllocated < 0 {
		a.QuantityAllocated = 0
	}
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
