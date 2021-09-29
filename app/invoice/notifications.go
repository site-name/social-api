package invoice

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/invoice"
)

func (s *ServiceInvoice) GetInvoicePayload(inVoice *invoice.Invoice) model.StringInterface {
	return model.StringInterface{
		"id":           inVoice.Id,
		"number":       inVoice.Number,
		"download_url": inVoice.ExternalUrl,
		"order_id":     inVoice.OrderID,
	}
}

// SendInvoice Send an invoice to user of related order with URL to download it
func (s *ServiceInvoice) SendInvoice(inVoice *invoice.Invoice, staffUser *account.User, _ interface{}, manager interface{}) *model.AppError {
	panic("not implemented")
}
