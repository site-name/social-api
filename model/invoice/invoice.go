package invoice

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/file"
)

// max lengths for Invoice
const (
	INVOICE_NUMBER_MAX_LENGTH       = 255
	INVOICE_EXTERNAL_URL_MAX_LENGTH = 2048
)

// TODO: considering add field Job to this model
type Invoice struct {
	Id          string        `json:"id"`
	OrderID     string        `json:"order_id"`
	Number      string        `json:"number"`
	CreateAt    int64         `json:"create_at"`
	ExternalUrl string        `json:"external_url"`
	InvoiceFile file.FileInfo `json:"invoice_file"`
	model.ModelMetadata
}

func (i *Invoice) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.invoice.is_valid.%s.app_error",
		"invoice_id=",
		"Invoice.IsValid",
	)
	if !model.IsValidId(i.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(i.OrderID) {
		return outer("order_id", &i.Id)
	}
	if len(i.Number) > INVOICE_NUMBER_MAX_LENGTH {
		return outer("number", &i.Id)
	}
	if len(i.ExternalUrl) > INVOICE_EXTERNAL_URL_MAX_LENGTH {
		return outer("external_url", &i.Id)
	}
	if i.CreateAt == 0 {
		return outer("create_at", &i.Id)
	}

	return nil
}

func (i *Invoice) PreSave() {
	if i.Id == "" {
		i.Id = model.NewId()
	}
	i.CreateAt = model.GetMillis()
}
