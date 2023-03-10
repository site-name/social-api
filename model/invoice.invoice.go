package model

import "github.com/Masterminds/squirrel"

// max lengths for Invoice
const (
	INVOICE_NUMBER_MAX_LENGTH       = 255
	INVOICE_EXTERNAL_URL_MAX_LENGTH = 2048
	INVOICE_STATUS_MAX_LENGTH       = 50
	INVOICE_MESSAGE_MAX_LENGTH      = 255
)

type Invoice struct {
	Id          string  `json:"id"`
	OrderID     *string `json:"order_id"`
	Number      string  `json:"number"`
	CreateAt    int64   `json:"create_at"`
	ExternalUrl string  `json:"external_url"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	UpdateAt    int64   `json:"update_at"`
	ModelMetadata
}

type InvoiceFilterOptions struct {
	Id      squirrel.Sqlizer
	OrderID squirrel.Sqlizer
	Limit   uint64
}

func (i *Invoice) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.invoice.is_valid.%s.app_error",
		"invoice_id=",
		"Invoice.IsValid",
	)
	if !IsValidId(i.Id) {
		return outer("id", nil)
	}
	if i.OrderID != nil && !IsValidId(*i.OrderID) {
		return outer("order_id", &i.Id)
	}
	if len(i.Number) > INVOICE_NUMBER_MAX_LENGTH {
		return outer("number", &i.Id)
	}
	if len(i.ExternalUrl) > INVOICE_EXTERNAL_URL_MAX_LENGTH {
		return outer("external_url", &i.Id)
	}
	if len(i.Status) > INVOICE_STATUS_MAX_LENGTH || !ALL_JOB_STATUSES.Contains(i.Status) {
		return outer("status", &i.Id)
	}
	if len(i.Message) > INVOICE_MESSAGE_MAX_LENGTH {
		return outer("message", &i.Id)
	}
	if i.CreateAt == 0 {
		return outer("create_at", &i.Id)
	}

	return nil
}

func (i *Invoice) PreSave() {
	if i.Id == "" {
		i.Id = NewId()
	}
	i.CreateAt = GetMillis()
	i.UpdateAt = i.CreateAt
}

func (i *Invoice) PreUpdate() {
	i.UpdateAt = GetMillis()
}
