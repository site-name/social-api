package order

import (
	"io"
	"regexp"
	"strings"

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

type Fulfillment struct {
	Id               string `json:"id"`
	FulfillmentOrder uint   `json:"fulfillment_order"`
	OrderID          string `json:"order_id"`
	Status           string `json:"status"`
	TrackingNumber   string `json:"tracking_numdber"`
	CreateAt         int64  `json:"create_at"`
	model.ModelMetadata
}

func (f *Fulfillment) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.fulfillment.is_valid.%s.app_error",
		"fulfillment_id=",
		"Fulfillment.IsValid",
	)
	if !model.IsValidId(f.Id) {
		return outer("Id", nil)
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

func FulfillmentFromJson(data io.Reader) *Fulfillment {
	var f Fulfillment
	model.ModelFromJson(&f, data)
	return &f
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

func (f *Fulfillment) CanEdit() bool {
	return f.Status != CANCELED
}

func (f *Fulfillment) IstrackingNumber() bool {
	return trackingNumberRegex.MatchString(f.TrackingNumber)
}
