package model

import (
	"fmt"
	"regexp"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
)

var (
	trackingNumberRegex = regexp.MustCompile(`^[-\w]+://`)
)

// max lengths for some fulfillment's fields
const (
	FULFILLMENT_STATUS_MAX_LENGTH          = 32
	FULFILLMENT_TRACKING_NUMBER_MAX_LENGTH = 255
)

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

type Fulfillment struct {
	Id                   string            `json:"id"`
	FulfillmentOrder     int               `json:"fulfillment_order"`
	OrderID              string            `json:"order_id"` // not null nor editable
	Status               FulfillmentStatus `json:"status"`
	TrackingNumber       string            `json:"tracking_numdber"`
	CreateAt             int64             `json:"create_at"`
	ShippingRefundAmount *decimal.Decimal  `json:"shipping_refund_amount"`
	TotalRefundAmount    *decimal.Decimal  `json:"total_refund_amount"`
	ModelMetadata

	Order *Order `json:"-" db:"-"` // this field get populated in queries that require select related data
}

// FulfillmentFilterOption is used to build squirrel sql queries
type FulfillmentFilterOption struct {
	Id      squirrel.Sqlizer
	OrderID squirrel.Sqlizer
	Status  squirrel.Sqlizer

	SelectRelatedOrder bool // if true, tells store to select related order also

	FulfillmentLineID squirrel.Sqlizer // LEFT/INNER JOIN FulfillmentLines ON (...) WHERE FulfillmentLines.Id ...

	SelectForUpdate bool // if true, add `FOR UPDATE`to the end of sql queries
}

type Fulfillments []*Fulfillment

func (f Fulfillments) IDs() []string {
	res := []string{}
	for _, item := range f {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (f *Fulfillment) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"fulfillment.is_valid.%s.app_error",
		"fulfillment_id=",
		"Fulfillment.IsValid",
	)
	if !IsValidId(f.Id) {
		return outer("id", nil)
	}
	if f.CreateAt == 0 {
		return outer("create_at", &f.Id)
	}
	if len(f.Status) > FULFILLMENT_STATUS_MAX_LENGTH || FulfillmentStrings[f.Status] == "" {
		return outer("status", &f.Id)
	}
	if len(f.TrackingNumber) > FULFILLMENT_TRACKING_NUMBER_MAX_LENGTH {
		return outer("tracking_number", &f.Id)
	}

	return nil
}

func (f *Fulfillment) ToJSON() string {
	return ModelToJson(f)
}

func (f *Fulfillment) PreSave() {
	if f.Id == "" {
		f.Id = NewId()
	}
	f.CreateAt = GetMillis()
	f.commonPre()
}

func (f *Fulfillment) commonPre() {
	if f.Status == "" {
		f.Status = FULFILLMENT_FULFILLED
	}
}

func (f *Fulfillment) PreUpdate() {
	f.commonPre()
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
		res.ShippingRefundAmount = NewDecimal(*f.ShippingRefundAmount)
	}
	if f.TotalRefundAmount != nil {
		res.TotalRefundAmount = NewDecimal(*f.TotalRefundAmount)
	}
	res.ModelMetadata = f.ModelMetadata.DeepCopy()
	if f.Order != nil {
		res.Order = f.Order.DeepCopy()
	}
	return &res
}
