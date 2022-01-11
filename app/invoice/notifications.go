package invoice

import (
	"github.com/sitename/sitename/app/plugin"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/model/order"
)

func GetInvoicePayload(inVoice invoice.Invoice) model.StringInterface {
	return model.StringInterface{
		"id":           inVoice.Id,
		"number":       inVoice.Number,
		"download_url": inVoice.ExternalUrl,
		"order_id":     inVoice.OrderID,
	}
}

// SendInvoice Send an invoice to user of related order with URL to download it
func (s *ServiceInvoice) SendInvoice(inVoice invoice.Invoice, staffUser *account.User, _ interface{}, manager interfaces.PluginManagerInterface) *model.AppError {

	var (
		orDer  *order.Order
		appErr *model.AppError
	)

	if inVoice.OrderID != nil {
		orDer, appErr = s.srv.OrderService().OrderById(*inVoice.OrderID)
		if appErr != nil {
			return appErr
		}
	}

	recipientEmail, appErr := s.srv.OrderService().CustomerEmail(orDer)
	if appErr != nil {
		return appErr
	}

	shop, appErr := s.srv.ShopService().ShopById(manager.GetShopID())
	if appErr != nil {
		return appErr
	}

	payload := map[string]interface{}{
		"invoice":          GetInvoicePayload(inVoice),
		"recipient_email":  recipientEmail,
		"requester_app_id": nil,
		"domain":           *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":        shop.Name,
	}

	_, appErr = manager.Notify(plugin.INVOICE_READY, payload, orDer.ChannelID, "")
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.InvoiceSent(inVoice, recipientEmail)
	if appErr != nil {
		return appErr
	}

	return nil
}
