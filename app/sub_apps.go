package app

import (
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/modules/slog"
)

// all sub applications of platform
var (
	accountApp   func(AppIface) sub_app_iface.AccountApp
	giftcardApp  func(AppIface) sub_app_iface.GiftcardApp
	paymentApp   func(AppIface) sub_app_iface.PaymentApp
	checkoutApp  func(AppIface) sub_app_iface.CheckoutApp
	warehouseApp func(AppIface) sub_app_iface.WarehouseApp
	productApp   func(AppIface) sub_app_iface.ProductApp
	wishlistApp  func(AppIface) sub_app_iface.WishlistApp
	orderApp     func(AppIface) sub_app_iface.OrderApp
	webhookApp   func(AppIface) sub_app_iface.WebhookApp
	menuApp      func(AppIface) sub_app_iface.MenuApp
	pageApp      func(AppIface) sub_app_iface.PageApp
	seoApp       func(AppIface) sub_app_iface.SeoApp
	siteApp      func(AppIface) sub_app_iface.SiteApp
	shippingApp  func(AppIface) sub_app_iface.ShippingApp
	discountApp  func(AppIface) sub_app_iface.DiscountApp
	csvApp       func(AppIface) sub_app_iface.CsvApp
	attributeApp func(AppIface) sub_app_iface.AttributeApp
	channelApp   func(AppIface) sub_app_iface.ChannelApp
	invoiceApp   func(AppIface) sub_app_iface.InvoiceApp
	fileApp      func(AppIface) sub_app_iface.FileApp
)

// RegisterFileApp
func RegisterFileApp(f func(AppIface) sub_app_iface.FileApp) {
	fileApp = f
}

// RegisterGiftcardApp
func RegisterGiftcardApp(f func(AppIface) sub_app_iface.GiftcardApp) {
	giftcardApp = f
}

// RegisterPaymentApp
func RegisterPaymentApp(f func(AppIface) sub_app_iface.PaymentApp) {
	paymentApp = f
}

func RegisterProductApp(f func(AppIface) sub_app_iface.ProductApp) {
	productApp = f
}

func RegisterWarehouseApp(f func(AppIface) sub_app_iface.WarehouseApp) {
	warehouseApp = f
}

func RegisterWishlistApp(f func(AppIface) sub_app_iface.WishlistApp) {
	wishlistApp = f
}

func RegisterCheckoutApp(f func(AppIface) sub_app_iface.CheckoutApp) {
	checkoutApp = f
}

func RegisterOrderApp(f func(AppIface) sub_app_iface.OrderApp) {
	orderApp = f
}

func RegisterWebhookApp(f func(AppIface) sub_app_iface.WebhookApp) {
	webhookApp = f
}

func RegisterMenuApp(f func(AppIface) sub_app_iface.MenuApp) {
	menuApp = f
}

func RegisterPageApp(f func(AppIface) sub_app_iface.PageApp) {
	pageApp = f
}

func RegisterSeoApp(f func(AppIface) sub_app_iface.SeoApp) {
	seoApp = f
}

func RegisterSiteApp(f func(AppIface) sub_app_iface.SiteApp) {
	siteApp = f
}

func RegisterShippingApp(f func(AppIface) sub_app_iface.ShippingApp) {
	shippingApp = f
}

func RegisterDiscountApp(f func(AppIface) sub_app_iface.DiscountApp) {
	discountApp = f
}

func RegisterCsvApp(f func(AppIface) sub_app_iface.CsvApp) {
	csvApp = f
}

func RegisterAttributeApp(f func(AppIface) sub_app_iface.AttributeApp) {
	attributeApp = f
}

func RegisterChannelApp(f func(AppIface) sub_app_iface.ChannelApp) {
	channelApp = f
}

func RegisterAccountApp(f func(AppIface) sub_app_iface.AccountApp) {
	accountApp = f
}

func RegisterInvoiceApp(f func(AppIface) sub_app_iface.InvoiceApp) {
	invoiceApp = f
}

func criticalLog(app string) {
	slog.Critical("Failed to register. Please check again", slog.String("app", app))
}

// registerAllSubApps register all sub app to App.
func registerAllSubApps() []AppOption {
	slog.Info("Registering all sub applications...")

	return []AppOption{
		func(a *App) {
			if productApp == nil {
				criticalLog("product")
				return
			}
			a.product = productApp(a)
		},
		func(a *App) {
			if seoApp == nil {
				criticalLog("seo")
				return
			}
			a.seo = seoApp(a)
		},
		func(a *App) {
			if siteApp == nil {
				criticalLog("site")
				return
			}
			a.site = siteApp(a)
		},
		func(a *App) {
			if paymentApp == nil {
				criticalLog("payment")
				return
			}
			a.payment = paymentApp(a)
		},
		func(a *App) {
			if shippingApp == nil {
				criticalLog("shipping")
				return
			}
			a.shipping = shippingApp(a)
		},
		func(a *App) {
			if menuApp == nil {
				criticalLog("menu")
				return
			}
			a.menu = menuApp(a)
		},
		func(a *App) {
			if orderApp == nil {
				criticalLog("order")
				return
			}
			a.order = orderApp(a)
		},
		func(a *App) {
			if webhookApp == nil {
				criticalLog("webhook")
				return
			}
			a.webhook = webhookApp(a)
		},
		func(a *App) {
			if warehouseApp == nil {
				criticalLog("warehouse")
				return
			}
			a.warehouse = warehouseApp(a)
		},
		func(a *App) {
			if checkoutApp == nil {
				criticalLog("checkout")
				return
			}
			a.checkout = checkoutApp(a)
		},
		func(a *App) {
			if discountApp == nil {
				criticalLog("discount")
				return
			}
			a.discount = discountApp(a)
		},
		func(a *App) {
			if wishlistApp == nil {
				criticalLog("wishlist")
				return
			}
			a.wishlist = wishlistApp(a)
		},
		func(a *App) {
			if giftcardApp == nil {
				criticalLog("giftcard")
				return
			}
			a.giftcard = giftcardApp(a)
		},
		func(a *App) {
			if pageApp == nil {
				criticalLog("page")
				return
			}
			a.page = pageApp(a)
		},
		func(a *App) {
			if csvApp == nil {
				criticalLog("csv")
				return
			}
			a.csv = csvApp(a)
		},
		func(a *App) {
			if attributeApp == nil {
				criticalLog("attribute")
				return
			}
			a.attribute = attributeApp(a)
		},
		func(a *App) {
			if channelApp == nil {
				criticalLog("channel")
				return
			}
			a.channel = channelApp(a)
		},
		func(a *App) {
			if accountApp == nil {
				criticalLog("account")
				return
			}
			a.account = accountApp(a)
		},
		func(a *App) {
			if invoiceApp == nil {
				criticalLog("invoice")
				return
			}
			a.invoice = invoiceApp(a)
		},
		func(a *App) {
			if fileApp == nil {
				criticalLog("file")
				return
			}
			a.file = fileApp(a)
		},
	}
}

// Order returns order sub app
func (a *App) OrderApp() sub_app_iface.OrderApp {
	return a.order
}

// Csv returns csv sub app
func (a *App) CsvApp() sub_app_iface.CsvApp {
	return a.csv
}

// Product returns product sub app
func (a *App) ProductApp() sub_app_iface.ProductApp {
	return a.product
}

// Payment returns payment sub app
func (a *App) PaymentApp() sub_app_iface.PaymentApp {
	return a.payment
}

// Giftcard returns giftcard sub app
func (a *App) GiftcardApp() sub_app_iface.GiftcardApp {
	return a.giftcard
}

// Site returns site sub app
func (a *App) SiteApp() sub_app_iface.SiteApp {
	return a.site
}

// Seo returns order seo app
func (a *App) SeoApp() sub_app_iface.SeoApp {
	return a.seo
}

// Shipping returns shipping sub app
func (a *App) ShippingApp() sub_app_iface.ShippingApp {
	return a.shipping
}

// Wishlist returns wishlist sub app
func (a *App) WishlistApp() sub_app_iface.WishlistApp {
	return a.wishlist
}

// Page returns page sub app
func (a *App) PageApp() sub_app_iface.PageApp {
	return a.page
}

// Menu returns menu sub app
func (a *App) MenuApp() sub_app_iface.MenuApp {
	return a.menu
}

// Attribute returns attribute sub app
func (a *App) AttributeApp() sub_app_iface.AttributeApp {
	return a.attribute
}

// Warehouse returns warehouse sub app
func (a *App) WarehouseApp() sub_app_iface.WarehouseApp {
	return a.warehouse
}

// Checkout returns checkout sub app
func (a *App) CheckoutApp() sub_app_iface.CheckoutApp {
	return a.checkout
}

// Webhook returns webhook sub app
func (a *App) WebhookApp() sub_app_iface.WebhookApp {
	return a.webhook
}

// Channel returns channel sub app
func (a *App) ChannelApp() sub_app_iface.ChannelApp {
	return a.channel
}

// Account returns account sub app
func (a *App) AccountApp() sub_app_iface.AccountApp {
	return a.account
}

// Invoice returns invoice sub app
func (a *App) InvoiceApp() sub_app_iface.InvoiceApp {
	return a.invoice
}

// FileApp returns file sub app
func (a *App) FileApp() sub_app_iface.FileApp {
	return a.file
}
