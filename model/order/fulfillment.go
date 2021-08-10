package order

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
)

var (
	trackingNumberRegex = regexp.MustCompile(`^[-\w]+://`)
)

// max lengths for some fulfillment's fields
const (
	FULFILLMENT_STATUS_MAX_LENGTH          = 32
	FULFILLMENT_TRACKING_NUMBER_MAX_LENGTH = 255
)

// fulfillment statuses
const (
	FULFILLMENT_FULFILLED             = "fulfilled"             // group of products in an order marked as fulfilled
	FULFILLMENT_REFUNDED              = "refunded"              // group of refunded products
	FULFILLMENT_RETURNED              = "returned"              // group of returned products
	FULFILLMENT_REFUNDED_AND_RETURNED = "refunded_and_returned" // group of returned and replaced products
	FULFILLMENT_REPLACED              = "replaced"              // group of replaced products
	FULFILLMENT_CANCELED              = "canceled"              // fulfilled group of products in an order marked as canceled
)

var FulfillmentStrings = map[string]string{
	FULFILLMENT_FULFILLED:             "Fulfilled",
	FULFILLMENT_REFUNDED:              "Refunded",
	FULFILLMENT_RETURNED:              "Returned",
	FULFILLMENT_REPLACED:              "Replaced",
	FULFILLMENT_REFUNDED_AND_RETURNED: "Refunded and returned",
	FULFILLMENT_CANCELED:              "Canceled",
}

type Fulfillment struct {
	Id                   string           `json:"id"`
	FulfillmentOrder     uint             `json:"fulfillment_order"`
	OrderID              string           `json:"order_id"` // not null nor editable
	Status               string           `json:"status"`
	TrackingNumber       string           `json:"tracking_numdber"`
	CreateAt             int64            `json:"create_at"`
	ShippingRefundAmount *decimal.Decimal `json:"shipping_refund_amount"`
	TotalRefundAmount    *decimal.Decimal `json:"total_refund_amount"`
	model.ModelMetadata
}

// FulfillmentFilterOption is used to build squirrel sql queries
type FulfillmentFilterOption struct {
	Id      *model.StringFilter
	OrderID *model.StringFilter
	Status  *model.StringFilter
}

func (f *Fulfillment) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.fulfillment.is_valid.%s.app_error",
		"fulfillment_id=",
		"Fulfillment.IsValid",
	)
	if !model.IsValidId(f.Id) {
		return outer("id", nil)
	}
	if f.CreateAt == 0 {
		return outer("create_at", &f.Id)
	}
	if len(f.Status) > FULFILLMENT_STATUS_MAX_LENGTH || FulfillmentStrings[strings.ToLower(f.Status)] == "" {
		return outer("status", &f.Id)
	}
	if len(f.TrackingNumber) > FULFILLMENT_TRACKING_NUMBER_MAX_LENGTH {
		return outer("tracking_number", &f.Id)
	}

	return nil
}

func (f *Fulfillment) ToJson() string {
	return model.ModelToJson(f)
}

func (f *Fulfillment) PreSave() {
	if f.Id == "" {
		f.Id = model.NewId()
	}
	f.CreateAt = model.GetMillis()
	if f.Status == "" {
		f.Status = FULFILLED
	}
}

func (f *Fulfillment) ComposedId() (string, error) {
	if !model.IsValidId(f.Id) {
		return "", errors.New("please save me first")
	}
	return fmt.Sprintf("%s-%d", f.OrderID, f.FulfillmentOrder), nil
}

// CanEdit checks if current Fulfillment's Status is "canceled"
func (f *Fulfillment) CanEdit() bool {
	return f.Status != CANCELED
}

func (f *Fulfillment) IstrackingNumber() bool {
	return trackingNumberRegex.MatchString(f.TrackingNumber)
}
