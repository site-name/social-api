package model

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"gorm.io/gorm"
)

var (
	trackingNumberRegex = regexp.MustCompile(`^[-\w]+://`)
)

// max lengths for some fulfillment's fields
// const (
// FULFILLMENT_DECIMAL_FIELD_MAX_DIGITS     = 12
// FULFILLMENT_DECIMAL_FIELD_DECIMAL_PLACES = 3
// )

type FulfillmentStatus string

// fulfillment statuses
const (
	FULFILLMENT_FULFILLED             FulfillmentStatus = "fulfilled"             // group of products in an order marked as fulfilled
	FULFILLMENT_REFUNDED              FulfillmentStatus = "refunded"              // group of refunded products
	FULFILLMENT_RETURNED              FulfillmentStatus = "returned"              // group of returned products
	FULFILLMENT_REFUNDED_AND_RETURNED FulfillmentStatus = "refunded_and_returned" // group of returned and replaced products
	FULFILLMENT_REPLACED              FulfillmentStatus = "replaced"              // group of replaced products
	FULFILLMENT_CANCELED              FulfillmentStatus = "canceled"              // fulfilled group of products in an order marked as canceled
	FULFILLMENT_WAITING_FOR_APPROVAL  FulfillmentStatus = "waiting_for_approval"  // group of products waiting for approval
)

var FulfillmentStrings = map[FulfillmentStatus]string{
	FULFILLMENT_FULFILLED:             "Fulfilled",
	FULFILLMENT_REFUNDED:              "Refunded",
	FULFILLMENT_RETURNED:              "Returned",
	FULFILLMENT_REPLACED:              "Replaced",
	FULFILLMENT_REFUNDED_AND_RETURNED: "Refunded and returned",
	FULFILLMENT_CANCELED:              "Canceled",
	FULFILLMENT_WAITING_FOR_APPROVAL:  "Waiting for approval",
}

func (f FulfillmentStatus) IsValid() bool {
	return FulfillmentStrings[f] != ""
}

type Fulfillment struct {
	Id                   string            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	FulfillmentOrder     int               `json:"fulfillment_order" gorm:"column:FulfillmentOrder;check:FulfillmentOrder >= 0"` // >= 0; not editable
	OrderID              string            `json:"order_id" gorm:"type:uuid;column:OrderID"`                                     // not null nor editable
	Status               FulfillmentStatus `json:"status" gorm:"type:varchar(32);column:Status"`                                 // default "fulfilled"
	TrackingNumber       string            `json:"tracking_numdber" gorm:"type:varchar(255);column:TrackingNumber"`
	CreateAt             int64             `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	ShippingRefundAmount *decimal.Decimal  `json:"shipping_refund_amount" gorm:"column:ShippingRefundAmount;type:decimal(12,3)"` // max digits 12, decimal places 3
	TotalRefundAmount    *decimal.Decimal  `json:"total_refund_amount" gorm:"column:TotalRefundAmount;type:decimal(12,3)"`       // max digits 12, decimal places 3
	ModelMetadata

	FulfillmentLines FulfillmentLines `json:"-" gorm:"foreignKey:FulfillmentID;constraint:OnDelete:CASCADE"`

	order *Order // this field get populated in queries that require select related data
}

func (c *Fulfillment) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Fulfillment) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Fulfillment) TableName() string             { return FulfillmentTableName }
func (f *Fulfillment) GetOrder() *Order              { return f.order }
func (f *Fulfillment) SetOrder(o *Order)             { f.order = o }

// FulfillmentFilterOption is used to build squirrel sql queries
// NOTE: `FulfillmentLineID` and `HaveNoFulfillmentLines` fields are evaluated exclusively. if ... else if
type FulfillmentFilterOption struct {
	Conditions squirrel.Sqlizer

	// INNER JOIN FulfillmentLines ON (...) WHERE FulfillmentLines.Id ...
	FulfillmentLineID squirrel.Sqlizer
	// LEFT JOIN FulfillmentLines ON ... WHERE FulfillmentLines.FulfillmentID = NULL
	HaveNoFulfillmentLines bool

	SelectRelatedOrder bool // if true, tells store to select related order also
	SelectForUpdate    bool // if true, add `FOR UPDATE`to the end of sql queries. NOTE: Only applied when Transaction field is set
	Transaction        *gorm.DB
}

type Fulfillments []*Fulfillment

func (f Fulfillments) IDs() []string {
	return lo.Map(f, func(item *Fulfillment, _ int) string { return item.Id })
}

func (f *Fulfillment) IsValid() *AppError {
	if !f.Status.IsValid() {
		return NewAppError("Fulfillment.IsValid", "model.fulfillment.is_valid.status.app_error", nil, "please provide valid status", http.StatusBadRequest)
	}
	return nil
}

func (f *Fulfillment) commonPre() {
	if f.Status == "" {
		f.Status = FULFILLMENT_FULFILLED
	}
}

func (f *Fulfillment) ComposedId() string {
	return fmt.Sprintf("%s-%d", f.OrderID, f.FulfillmentOrder)
}

// CanEdit checks if current Fulfillment's Status is "canceled"
func (f *Fulfillment) CanEdit() bool {
	return f.Status != FULFILLMENT_FULFILLED
}

func (f *Fulfillment) IstrackingNumber() bool {
	return trackingNumberRegex.MatchString(f.TrackingNumber)
}

func (f *Fulfillment) DeepCopy() *Fulfillment {
	res := *f
	if f.ShippingRefundAmount != nil {
		res.ShippingRefundAmount = NewPrimitive(*f.ShippingRefundAmount)
	}
	if f.TotalRefundAmount != nil {
		res.TotalRefundAmount = NewPrimitive(*f.TotalRefundAmount)
	}
	res.ModelMetadata = f.ModelMetadata.DeepCopy()
	if f.order != nil {
		res.order = f.order.DeepCopy()
	}
	return &res
}
