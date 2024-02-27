package app

import (
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/modules/slog"
)

var serverOpts []func(s *Server) error

func RegisterService(f func(s *Server) error) {
	if f == nil {
		panic("f cannot be nil")
	}
	serverOpts = append(serverOpts, f)
}

func (s *Server) registerSubServices() error {
	slog.Info("Registering all sub services...")

	for _, opt := range serverOpts {
		err := opt(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) PluginService() sub_app_iface.PluginService {
	return a.srv.Plugin
}

func (a *App) OrderService() sub_app_iface.OrderService {
	return a.srv.Order
}

func (a *App) CsvService() sub_app_iface.CsvService {
	return a.srv.Csv
}

func (a *App) ProductService() sub_app_iface.ProductService {
	return a.srv.Product
}

func (a *App) PaymentService() sub_app_iface.PaymentService {
	return a.srv.Payment
}

func (a *App) GiftcardService() sub_app_iface.GiftcardService {
	return a.srv.Giftcard
}

func (a *App) SeoService() sub_app_iface.SeoService {
	return a.srv.Seo
}

func (a *App) ShippingService() sub_app_iface.ShippingService {
	return a.srv.Shipping
}

func (a *App) WishlistService() sub_app_iface.WishlistService {
	return a.srv.Wishlist
}

func (a *App) PageService() sub_app_iface.PageService {
	return a.srv.Page
}

func (a *App) MenuService() sub_app_iface.MenuService {
	return a.srv.Menu
}

func (a *App) AttributeService() sub_app_iface.AttributeService {
	return a.srv.Attribute
}

func (a *App) WarehouseService() sub_app_iface.WarehouseService {
	return a.srv.Warehouse
}

func (a *App) CheckoutService() sub_app_iface.CheckoutService {
	return a.srv.Checkout
}

func (a *App) WebhookService() sub_app_iface.WebhookService {
	return a.srv.Webhook
}

func (a *App) ChannelService() sub_app_iface.ChannelService {
	return a.srv.Channel
}

func (a *App) AccountService() sub_app_iface.AccountService {
	return a.srv.Account
}

func (a *App) InvoiceService() sub_app_iface.InvoiceService {
	return a.srv.Invoice
}

func (a *App) FileService() sub_app_iface.FileService {
	return a.srv.File
}

func (a *App) DiscountService() sub_app_iface.DiscountService {
	return a.srv.Discount
}
