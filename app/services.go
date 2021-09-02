package app

import (
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/modules/slog"
)

// all sub applications of platform
var (
	accountApp   func(*Server) (sub_app_iface.AccountService, error)
	giftcardApp  func(*Server) (sub_app_iface.GiftcardService, error)
	paymentApp   func(*Server) (sub_app_iface.PaymentService, error)
	checkoutApp  func(*Server) (sub_app_iface.CheckoutService, error)
	warehouseApp func(*Server) (sub_app_iface.WarehouseService, error)
	productApp   func(*Server) (sub_app_iface.ProductService, error)
	wishlistApp  func(*Server) (sub_app_iface.WishlistService, error)
	orderApp     func(*Server) (sub_app_iface.OrderService, error)
	webhookApp   func(*Server) (sub_app_iface.WebhookService, error)
	menuApp      func(*Server) (sub_app_iface.MenuService, error)
	pageApp      func(*Server) (sub_app_iface.PageService, error)
	seoApp       func(*Server) (sub_app_iface.SeoService, error)
	shopApp      func(*Server) (sub_app_iface.ShopService, error)
	shippingApp  func(*Server) (sub_app_iface.ShippingService, error)
	discountApp  func(*Server) (sub_app_iface.DiscountService, error)
	csvApp       func(*Server) (sub_app_iface.CsvService, error)
	attributeApp func(*Server) (sub_app_iface.AttributeService, error)
	channelApp   func(*Server) (sub_app_iface.ChannelService, error)
	invoiceApp   func(*Server) (sub_app_iface.InvoiceService, error)
	fileApp      func(*Server) (sub_app_iface.FileService, error)
	pluginApp    func(*Server) (sub_app_iface.PluginService, error)
)

// RegisterPluginApp
func RegisterPluginApp(f func(*Server) (sub_app_iface.PluginService, error)) {
	pluginApp = f
}

// RegisterFileApp
func RegisterFileApp(f func(*Server) (sub_app_iface.FileService, error)) {
	fileApp = f
}

// RegisterGiftcardApp
func RegisterGiftcardApp(f func(*Server) (sub_app_iface.GiftcardService, error)) {
	giftcardApp = f
}

// RegisterPaymentApp
func RegisterPaymentApp(f func(*Server) (sub_app_iface.PaymentService, error)) {
	paymentApp = f
}

func RegisterProductApp(f func(*Server) (sub_app_iface.ProductService, error)) {
	productApp = f
}

func RegisterWarehouseApp(f func(*Server) (sub_app_iface.WarehouseService, error)) {
	warehouseApp = f
}

func RegisterWishlistApp(f func(*Server) (sub_app_iface.WishlistService, error)) {
	wishlistApp = f
}

func RegisterCheckoutApp(f func(*Server) (sub_app_iface.CheckoutService, error)) {
	checkoutApp = f
}

func RegisterOrderApp(f func(*Server) (sub_app_iface.OrderService, error)) {
	orderApp = f
}

func RegisterWebhookApp(f func(*Server) (sub_app_iface.WebhookService, error)) {
	webhookApp = f
}

func RegisterMenuApp(f func(*Server) (sub_app_iface.MenuService, error)) {
	menuApp = f
}

func RegisterPageApp(f func(*Server) (sub_app_iface.PageService, error)) {
	pageApp = f
}

func RegisterSeoApp(f func(*Server) (sub_app_iface.SeoService, error)) {
	seoApp = f
}

func RegisterShopApp(f func(*Server) (sub_app_iface.ShopService, error)) {
	shopApp = f
}

func RegisterShippingApp(f func(*Server) (sub_app_iface.ShippingService, error)) {
	shippingApp = f
}

func RegisterDiscountApp(f func(*Server) (sub_app_iface.DiscountService, error)) {
	discountApp = f
}

func RegisterCsvApp(f func(*Server) (sub_app_iface.CsvService, error)) {
	csvApp = f
}

func RegisterAttributeApp(f func(*Server) (sub_app_iface.AttributeService, error)) {
	attributeApp = f
}

func RegisterChannelApp(f func(*Server) (sub_app_iface.ChannelService, error)) {
	channelApp = f
}

func RegisterAccountApp(f func(*Server) (sub_app_iface.AccountService, error)) {
	accountApp = f
}

func RegisterInvoiceApp(f func(*Server) (sub_app_iface.InvoiceService, error)) {
	invoiceApp = f
}

func criticalLog(app string) {
	slog.Critical("Failed to register. Please check again", slog.String("app", app))
}

// registerSubServices register all sub services to App.
func (s *Server) registerSubServices() error {
	slog.Info("Registering all sub services...")

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
