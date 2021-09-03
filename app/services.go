package app

import (
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/modules/slog"
)

// all sub applications of platform
var (
	accountService   func(*Server) (sub_app_iface.AccountService, error)
	giftcardService  func(*Server) (sub_app_iface.GiftcardService, error)
	paymentService   func(*Server) (sub_app_iface.PaymentService, error)
	checkoutService  func(*Server) (sub_app_iface.CheckoutService, error)
	warehouseService func(*Server) (sub_app_iface.WarehouseService, error)
	productService   func(*Server) (sub_app_iface.ProductService, error)
	wishlistService  func(*Server) (sub_app_iface.WishlistService, error)
	orderService     func(*Server) (sub_app_iface.OrderService, error)
	webhookService   func(*Server) (sub_app_iface.WebhookService, error)
	menuService      func(*Server) (sub_app_iface.MenuService, error)
	pageService      func(*Server) (sub_app_iface.PageService, error)
	seoService       func(*Server) (sub_app_iface.SeoService, error)
	shopService      func(*Server) (sub_app_iface.ShopService, error)
	shippingService  func(*Server) (sub_app_iface.ShippingService, error)
	discountService  func(*Server) (sub_app_iface.DiscountService, error)
	csvService       func(*Server) (sub_app_iface.CsvService, error)
	attributeService func(*Server) (sub_app_iface.AttributeService, error)
	channelService   func(*Server) (sub_app_iface.ChannelService, error)
	invoiceService   func(*Server) (sub_app_iface.InvoiceService, error)
	fileService      func(*Server) (sub_app_iface.FileService, error)
	pluginService    func(*Server) (sub_app_iface.PluginService, error)
)

// RegisterPluginService
func RegisterPluginService(f func(*Server) (sub_app_iface.PluginService, error)) {
	pluginService = f
}

// RegisterFileService
func RegisterFileService(f func(*Server) (sub_app_iface.FileService, error)) {
	fileService = f
}

// RegisterGiftcardService
func RegisterGiftcardService(f func(*Server) (sub_app_iface.GiftcardService, error)) {
	giftcardService = f
}

// RegisterPaymentService
func RegisterPaymentService(f func(*Server) (sub_app_iface.PaymentService, error)) {
	paymentService = f
}

func RegisterProductService(f func(*Server) (sub_app_iface.ProductService, error)) {
	productService = f
}

func RegisterWarehouseService(f func(*Server) (sub_app_iface.WarehouseService, error)) {
	warehouseService = f
}

func RegisterWishlistService(f func(*Server) (sub_app_iface.WishlistService, error)) {
	wishlistService = f
}

func RegisterCheckoutService(f func(*Server) (sub_app_iface.CheckoutService, error)) {
	checkoutService = f
}

func RegisterOrderService(f func(*Server) (sub_app_iface.OrderService, error)) {
	orderService = f
}

func RegisterWebhookService(f func(*Server) (sub_app_iface.WebhookService, error)) {
	webhookService = f
}

func RegisterMenuService(f func(*Server) (sub_app_iface.MenuService, error)) {
	menuService = f
}

func RegisterPageService(f func(*Server) (sub_app_iface.PageService, error)) {
	pageService = f
}

func RegisterSeoService(f func(*Server) (sub_app_iface.SeoService, error)) {
	seoService = f
}

func RegisterShopService(f func(*Server) (sub_app_iface.ShopService, error)) {
	shopService = f
}

func RegisterShippingService(f func(*Server) (sub_app_iface.ShippingService, error)) {
	shippingService = f
}

func RegisterDiscountService(f func(*Server) (sub_app_iface.DiscountService, error)) {
	discountService = f
}

func RegisterCsvService(f func(*Server) (sub_app_iface.CsvService, error)) {
	csvService = f
}

func RegisterAttributeService(f func(*Server) (sub_app_iface.AttributeService, error)) {
	attributeService = f
}

func RegisterChannelService(f func(*Server) (sub_app_iface.ChannelService, error)) {
	channelService = f
}

func RegisterAccountService(f func(*Server) (sub_app_iface.AccountService, error)) {
	accountService = f
}

func RegisterInvoiceService(f func(*Server) (sub_app_iface.InvoiceService, error)) {
	invoiceService = f
}

func criticalLog(service string) {
	slog.Critical("Failed to register a service", slog.String("service", service))
}

// registerSubServices register all sub services to App.
func (s *Server) registerSubServices() error {
	slog.Info("Registering all sub services...")

	var err error
	s.account, err = accountService(s)
	if err != nil {
		return err
	}
	s.giftcard, err = giftcardService(s)
	if err != nil {
		return err
	}
	s.payment, err = paymentService(s)
	if err != nil {
		return err
	}
	s.checkout, err = checkoutService(s)
	if err != nil {
		return err
	}
	s.warehouse, err = warehouseService(s)
	if err != nil {
		return err
	}
	s.product, err = productService(s)
	if err != nil {
		return err
	}
	s.wishlist, err = wishlistService(s)
	if err != nil {
		return err
	}
	s.order, err = orderService(s)
	if err != nil {
		return err
	}
	s.webhook, err = webhookService(s)
	if err != nil {
		return err
	}
	s.menu, err = menuService(s)
	if err != nil {
		return err
	}
	s.page, err = pageService(s)
	if err != nil {
		return err
	}
	s.seo, err = seoService(s)
	if err != nil {
		return err
	}
	s.shop, err = shopService(s)
	if err != nil {
		return err
	}
	s.shipping, err = shippingService(s)
	if err != nil {
		return err
	}
	s.csv, err = csvService(s)
	if err != nil {
		return err
	}
	s.attribute, err = attributeService(s)
	if err != nil {
		return err
	}
	s.channel, err = channelService(s)
	if err != nil {
		return err
	}
	s.invoice, err = invoiceService(s)
	if err != nil {
		return err
	}
	s.file, err = fileService(s)
	if err != nil {
		return err
	}
	s.plugin, err = pluginService(s)
	if err != nil {
		return err
	}
	return nil
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
