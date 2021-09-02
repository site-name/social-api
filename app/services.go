package app

import (
	"github.com/sitename/sitename/app/account"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/modules/slog"
)

// all sub applications of platform
var (
	accountApp   func(AppIface) sub_app_iface.AccountService
	giftcardApp  func(AppIface) sub_app_iface.GiftcardService
	paymentApp   func(AppIface) sub_app_iface.PaymentService
	checkoutApp  func(AppIface) sub_app_iface.CheckoutService
	warehouseApp func(AppIface) sub_app_iface.WarehouseService
	productApp   func(AppIface) sub_app_iface.ProductService
	wishlistApp  func(AppIface) sub_app_iface.WishlistService
	orderApp     func(AppIface) sub_app_iface.OrderService
	webhookApp   func(AppIface) sub_app_iface.WebhookService
	menuApp      func(AppIface) sub_app_iface.MenuService
	pageApp      func(AppIface) sub_app_iface.PageService
	seoApp       func(AppIface) sub_app_iface.SeoService
	shopApp      func(AppIface) sub_app_iface.ShopService
	shippingApp  func(AppIface) sub_app_iface.ShippingService
	discountApp  func(AppIface) sub_app_iface.DiscountService
	csvApp       func(AppIface) sub_app_iface.CsvService
	attributeApp func(AppIface) sub_app_iface.AttributeService
	channelApp   func(AppIface) sub_app_iface.ChannelService
	invoiceApp   func(AppIface) sub_app_iface.InvoiceService
	fileApp      func(AppIface) sub_app_iface.FileService
	pluginApp    func(AppIface) sub_app_iface.PluginService
)

func criticalLog(app string) {
	slog.Critical("Failed to register. Please check again", slog.String("app", app))
}

// registerSubServices register all sub services to App.
func (server *Server) registerSubServices() error {
	slog.Info("Registering all sub services...")
	account.NewServiceAccount(&account.ServiceAccountConfig{
		CacheProvider: server.,
	})
}

// PluginService returns order sub app
func (s *Server) PluginService() sub_app_iface.PluginService {
	return s.plugin
}

// Order returns order sub app
func (s *Server) OrderService() sub_app_iface.OrderService {
	return s.order
}

// Csv returns csv sub app
func (s *Server) CsvService() sub_app_iface.CsvService {
	return s.csv
}

// Product returns product sub app
func (s *Server) ProductService() sub_app_iface.ProductService {
	return s.product
}

// Payment returns payment sub app
func (s *Server) PaymentService() sub_app_iface.PaymentService {
	return s.payment
}

// Giftcard returns giftcard sub app
func (s *Server) GiftcardService() sub_app_iface.GiftcardService {
	return s.giftcard
}

// ShopService returns shop sub app
func (s *Server) ShopService() sub_app_iface.ShopService {
	return s.shop
}

// Seo returns order seo app
func (s *Server) SeoService() sub_app_iface.SeoService {
	return s.seo
}

// Shipping returns shipping sub app
func (s *Server) ShippingService() sub_app_iface.ShippingService {
	return s.shipping
}

// Wishlist returns wishlist sub app
func (s *Server) WishlistService() sub_app_iface.WishlistService {
	return s.wishlist
}

// Page returns page sub app
func (s *Server) PageService() sub_app_iface.PageService {
	return s.page
}

// Menu returns menu sub app
func (s *Server) MenuService() sub_app_iface.MenuService {
	return s.menu
}

// Attribute returns attribute sub app
func (s *Server) AttributeService() sub_app_iface.AttributeService {
	return s.attribute
}

// Warehouse returns warehouse sub app
func (s *Server) WarehouseService() sub_app_iface.WarehouseService {
	return s.warehouse
}

// Checkout returns checkout sub app
func (s *Server) CheckoutService() sub_app_iface.CheckoutService {
	return s.checkout
}

// Webhook returns webhook sub app
func (s *Server) WebhookService() sub_app_iface.WebhookService {
	return s.webhook
}

// Channel returns channel sub app
func (s *Server) ChannelService() sub_app_iface.ChannelService {
	return s.channel
}

// Account returns account sub app
func (s *Server) AccountService() sub_app_iface.AccountService {
	return s.account
}

// Invoice returns invoice sub app
func (s *Server) InvoiceService() sub_app_iface.InvoiceService {
	return s.invoice
}

// FileService returns file sub app
func (s *Server) FileService() sub_app_iface.FileService {
	return s.file
}

// DiscountService returns discount sub app
func (s *Server) DiscountService() sub_app_iface.DiscountService {
	return s.discount
}
