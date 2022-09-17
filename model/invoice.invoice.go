package model

// max lengths for Invoice
const (
	INVOICE_NUMBER_MAX_LENGTH       = 255
	INVOICE_EXTERNAL_URL_MAX_LENGTH = 2048
)

type Invoice struct {
	Id          string  `json:"id"`
	OrderID     *string `json:"order_id"`
	Number      string  `json:"number"`
	CreateAt    int64   `json:"create_at"`
	ExternalUrl string  `json:"external_url"`
	ModelMetadata
}

func (i *Invoice) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"invoice.is_valid.%s.app_error",
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
}
