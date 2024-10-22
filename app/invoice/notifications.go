package invoice

import (
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
)

func GetInvoicePayload(inVoice model.Invoice) model_types.JSONString {
	return model_types.JSONString{
		"id":           inVoice.ID,
		"number":       inVoice.Number,
		"download_url": inVoice.ExternalURL,
		"order_id":     inVoice.OrderID,
	}
}

// SendInvoice Send an invoice to user of related order with URL to download it
func (s *ServiceInvoice) SendInvoice(inVoice model.Invoice, staffUser model.User, _ any, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	var (
		orDer  *model.Order
		appErr *model_helper.AppError
	)

	if !inVoice.OrderID.IsNil() {
		orDer, appErr = s.srv.Order.OrderById(*inVoice.OrderID.String)
		if appErr != nil {
			return appErr
		}
	}

	recipientEmail, appErr := s.srv.Order.CustomerEmail(orDer)
	if appErr != nil {
		return appErr
	}

	payload := map[string]any{
		"invoice":           GetInvoicePayload(inVoice),
		"recipient_email":   recipientEmail,
		"requester_app_id":  nil,
		"requester_user_id": staffUser.ID,
		"domain":            *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":         s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr = manager.Notify(model_helper.INVOICE_READY, payload, orDer.ChannelID, "")
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.InvoiceSent(inVoice, recipientEmail)
	if appErr != nil {
		return appErr
	}

	return nil
}
