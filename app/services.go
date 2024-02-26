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

// registerSubServices register all sub services to App.
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

func (s *Server) PluginService() sub_app_iface.PluginService {
	return s.Plugin
}

func (s *Server) OrderService() sub_app_iface.OrderService {
	return s.Order
}

func (s *Server) CsvService() sub_app_iface.CsvService {
	return s.Csv
}

func (s *Server) ProductService() sub_app_iface.ProductService {
	return s.Product
}

func (s *Server) PaymentService() sub_app_iface.PaymentService {
	return s.Payment
}

func (s *Server) GiftcardService() sub_app_iface.GiftcardService {
	return s.Giftcard
}

func (s *Server) SeoService() sub_app_iface.SeoService {
	return s.Seo
}

func (s *Server) ShippingService() sub_app_iface.ShippingService {
	return s.Shipping
}

func (s *Server) WishlistService() sub_app_iface.WishlistService {
	return s.Wishlist
}

func (s *Server) PageService() sub_app_iface.PageService {
	return s.Page
}

func (s *Server) MenuService() sub_app_iface.MenuService {
	return s.Menu
}

func (s *Server) AttributeService() sub_app_iface.AttributeService {
	return s.Attribute
}

func (s *Server) WarehouseService() sub_app_iface.WarehouseService {
	return s.Warehouse
}

func (s *Server) CheckoutService() sub_app_iface.CheckoutService {
	return s.Checkout
}

func (s *Server) WebhookService() sub_app_iface.WebhookService {
	return s.Webhook
}

func (s *Server) ChannelService() sub_app_iface.ChannelService {
	return s.Channel
}

func (s *Server) AccountService() sub_app_iface.AccountService {
	return s.Account
}

func (s *Server) InvoiceService() sub_app_iface.InvoiceService {
	return s.Invoice
}

func (s *Server) FileService() sub_app_iface.FileService {
	return s.File
}

func (s *Server) DiscountService() sub_app_iface.DiscountService {
	return s.Discount
}
