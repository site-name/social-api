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

// PluginService returns plugin sub app
func (s *Server) PluginService() sub_app_iface.PluginService {
	return s.Plugin
}

// Order returns order sub app
func (s *Server) OrderService() sub_app_iface.OrderService {
	return s.Order
}

// Csv returns csv sub app
func (s *Server) CsvService() sub_app_iface.CsvService {
	return s.Csv
}

// Product returns product sub app
func (s *Server) ProductService() sub_app_iface.ProductService {
	return s.Product
}

// Payment returns payment sub app
func (s *Server) PaymentService() sub_app_iface.PaymentService {
	return s.Payment
}

// Giftcard returns giftcard sub app
func (s *Server) GiftcardService() sub_app_iface.GiftcardService {
	return s.Giftcard
}

// Seo returns order seo app
func (s *Server) SeoService() sub_app_iface.SeoService {
	return s.Seo
}

// Shipping returns shipping sub app
func (s *Server) ShippingService() sub_app_iface.ShippingService {
	return s.Shipping
}

// Wishlist returns wishlist sub app
func (s *Server) WishlistService() sub_app_iface.WishlistService {
	return s.Wishlist
}

// Page returns page sub app
func (s *Server) PageService() sub_app_iface.PageService {
	return s.Page
}

// Menu returns menu sub app
func (s *Server) MenuService() sub_app_iface.MenuService {
	return s.Menu
}

// Attribute returns attribute sub app
func (s *Server) AttributeService() sub_app_iface.AttributeService {
	return s.Attribute
}

// Warehouse returns warehouse sub app
func (s *Server) WarehouseService() sub_app_iface.WarehouseService {
	return s.Warehouse
}

// Checkout returns checkout sub app
func (s *Server) CheckoutService() sub_app_iface.CheckoutService {
	return s.Checkout
}

// Webhook returns webhook sub app
func (s *Server) WebhookService() sub_app_iface.WebhookService {
	return s.Webhook
}

// Channel returns channel sub app
func (s *Server) ChannelService() sub_app_iface.ChannelService {
	return s.Channel
}

// Account returns account sub app
func (s *Server) AccountService() sub_app_iface.AccountService {
	return s.Account
}

// Invoice returns invoice sub app
func (s *Server) InvoiceService() sub_app_iface.InvoiceService {
	return s.Invoice
}

// FileService returns file sub app
func (s *Server) FileService() sub_app_iface.FileService {
	return s.File
}

// DiscountService returns discount sub app
func (s *Server) DiscountService() sub_app_iface.DiscountService {
	return s.Discount
}
