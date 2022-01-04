package plugins

import (
	"fmt"
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
)

type ConfigurationTypeField string

// PluginMethodNotImplemented is used to indicate if a method is implemented or not
type PluginMethodNotImplemented struct {
	MethodName string
}

func (p *PluginMethodNotImplemented) Error() string {
	return fmt.Sprintf("Method: %s is not implemented on this plugin", p.MethodName)
}

const (
	STRING           ConfigurationTypeField = "String"
	MULTILINE        ConfigurationTypeField = "Multiline"
	BOOLEAN          ConfigurationTypeField = "Boolean"
	SECRET           ConfigurationTypeField = "Secret"
	SECRET_MULTILINE ConfigurationTypeField = "SecretMultiline"
	PASSWORD         ConfigurationTypeField = "Password"
	OUTPUT           ConfigurationTypeField = "OUTPUT"
)

type PluginConfigurationType []model.StringInterface

type ExternalAccessToken struct {
	Token        *string
	RefreshToken *string
	CsrfToken    *string
	User         *account.User
}

// PluginManifest
type PluginManifest struct {
	PluginName              string
	PluginID                string
	Description             string
	ConfigStructure         map[string]model.StringInterface
	ConfigurationPerChannel bool
	DefaultConfiguration    []model.StringInterface
	DefaultActive           bool
	MetaCodeKey             string
	MetaDescriptionKey      string
}

type BasePluginInterface interface {
	fmt.Stringer
	// Check if given plugin_id matches with the PLUGIN_ID of this plugin
	CheckPluginId(pluginID string) bool
	// Handle authentication request responsible for obtaining access tokens.
	// Overwrite this method if the plugin handles authentication flow.
	ExternalAuthenticationUrl(data model.StringInterface, request *http.Request, previousValue interface{}) (model.StringInterface, *PluginMethodNotImplemented)
	// Handle authentication request responsible for obtaining access tokens.
	// Overwrite this method if the plugin handles authentication flow.
	ExternalObtainAccessTokens(data model.StringInterface, request *http.Request, previousValue interface{}) (*ExternalAccessToken, *PluginMethodNotImplemented)
	// Handle authentication refresh request.
	// Overwrite this method if the plugin handles authentication flow and supports
	// refreshing the access.
	ExternalRefresh(data model.StringInterface, request *http.Request, previousValue interface{}) (*ExternalAccessToken, *PluginMethodNotImplemented)
	// Handle logout request.
	// Overwrite this method if the plugin handles logout flow.
	ExternalLogout(data model.StringInterface, request *http.Request, previousValue interface{}) *PluginMethodNotImplemented
	// Verify the provided authentication data.
	// Overwrite this method if the plugin should validate the authentication data.
	ExternalVerify(data model.StringInterface, request *http.Request, previousValue interface{}) (*account.User, model.StringInterface, *PluginMethodNotImplemented)
	// Authenticate user which should be assigned to the request.
	// Overwrite this method if the plugin handles authentication flow.
	AuthenticateUser(request *http.Request, previousValue interface{}) (*account.User, *PluginMethodNotImplemented)
	// Handle received http request.
	// Overwrite this method if the plugin expects the incoming requests.
	Webhook(request *http.Request, path string, previousValue interface{}) (http.Response, *PluginMethodNotImplemented)
	// Handle notification request.
	// Overwrite this method if the plugin is responsible for sending notifications.
	Notify(event interface{}, payload model.StringInterface, previousValue interface{}) *PluginMethodNotImplemented
	//
	ChangeUserAddress(address *account.Address, addressType string, user *account.User, previousValue *account.Address) (*account.Address, *PluginMethodNotImplemented)
	// Calculate the total for checkout.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of a checkout total. Return TaxedMoney.
	CalculateCheckoutTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Calculate the shipping costs for checkout.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of shipping costs. Return TaxedMoney.
	CalculateCheckoutShipping(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Calculate the shipping costs for the order.
	// Update shipping costs in the order in case of changes in shipping address or
	// changes in draft order. Return TaxedMoney.
	CalculateOrderShipping(orDer *order.Order, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Calculate checkout line total.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of a checkout line total. Return TaxedMoney.
	CalculateCheckoutLineTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Calculate order line total.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of a order line total. Return TaxedMoney.
	CalculateOrderLineTotal(orDer *order.Order, orderLine *order.OrderLine, variant product_and_discount.ProductVariant, product product_and_discount.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Calculate checkout line unit price
	CalculateCheckoutLineUnitPrice(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Calculate order line unit price.
	// Update order line unit price in the order in case of changes in draft order.
	// Return TaxedMoney.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of an order line unit price.
	CalculateOrderLineUnit(orDer order.Order, orderLine order.OrderLine, variant product_and_discount.ProductVariant, product product_and_discount.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	//
	GetCheckoutLineTaxRate(checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented)
	//
	GetOrderLineTaxRate(orDer order.Order, product product_and_discount.Product, variant product_and_discount.ProductVariant, address *account.Address, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented)
	//
	GetCheckoutShippingTaxRate(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented)
	//
	GetOrderShippingTaxRate(orDer order.Order, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented)
	// Return list of all tax categories.
	// The returned list will be used to provide staff users with the possibility to
	// assign tax categories to a product. It can be used by tax plugins to properly
	// calculate taxes for products.
	// Overwrite this method in case your plugin provides a list of tax categories.
	GetTaxRateTypeChoices(previousValue []*model.TaxType) ([]*model.TaxType, *PluginMethodNotImplemented)
	// Define if storefront should add info about taxes to the price.
	// It is used only by the old storefront. The returned value determines if
	// storefront should append info to the price about "including/excluding X% VAT"
	ShowTaxesOnStorefront(previousValue bool) (bool, *PluginMethodNotImplemented)
	// Apply taxes to the shipping costs based on the shipping address.
	// Overwrite this method if you want to show available shipping methods with taxes.
	ApplyTaxesToShipping(price goprices.Money, shippingAddress account.Address, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Apply taxes to the product price based on the customer country.
	// Overwrite this method if you want to show products with taxes.
	ApplyTaxesToProduct(product product_and_discount.Product, price goprices.Money, country string, previousVlaue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented)
	// Trigger directly before order creation.
	// Overwrite this method if you need to trigger specific logic before an order is created.
	PreprocessOrderCreation(checkoutInfo checkout.CheckoutInfo, discounts []*product_and_discount.DiscountInfo, lines checkout.CheckoutLineInfos, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when order is created.
	// Overwrite this method if you need to trigger specific logic after an order is created.
	OrderCreated(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	//
	DraftOrderCreated(orDer order.Order, defaultValue interface{}) (interface{}, *PluginMethodNotImplemented)
	//
	DraftOrderUpdated(orDer order.Order, defaultValue interface{}) (interface{}, *PluginMethodNotImplemented)
	//
	DraftOrderDeleted(orDer order.Order, defaultValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when order is confirmed by staff.
	// Overwrite this method if you need to trigger specific logic after an order is
	// confirmed.
	OrderConfirmed(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when sale is created.
	// Overwrite this method if you need to trigger specific logic after sale is created.
	SaleCreated(sale product_and_discount.Sale, currentCatalogue product_and_discount.NodeCatalogueInfo, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when sale is deleted.
	// Overwrite this method if you need to trigger specific logic after sale is deleted.
	SaleDeleted(sale product_and_discount.Sale, previousCatalogue product_and_discount.NodeCatalogueInfo, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when sale is updated.
	// Overwrite this method if you need to trigger specific logic after sale is updated.
	SaleUpdated(sale product_and_discount.Sale, previousCatalogue product_and_discount.NodeCatalogueInfo, currentCatalogue product_and_discount.NodeCatalogueInfo, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when invoice creation starts.
	// Overwrite to create invoice with proper data, call invoice.update_invoice.
	InvoiceRequest(orDer order.Order, inVoice invoice.Invoice, number string, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger before invoice is deleted.
	// Perform any extra logic before the invoice gets deleted.
	// Note there is no need to run invoice.delete() as it will happen in mutation.
	InvoiceDelete(inVoice invoice.Invoice, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger after invoice is sent.
	InvoiceSent(inVoice invoice.Invoice, email string, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Return tax code from object meta.
	//
	// NOTE: obj can be 'Product' or 'ProductType'
	AssignTaxCodeToObjectMeta(obj interface{}, taxCode string, previousValue model.TaxType) (*model.TaxType, *PluginMethodNotImplemented)
	// Return tax code from object meta
	//
	// NOTE: obj must be either Product or ProductType
	GetTaxCodeFromObjectMeta(obj interface{}, previousValue model.TaxType) (*model.TaxType, *PluginMethodNotImplemented)
	// Return tax rate percentage value for a given tax rate type in a country.
	// It is used only by the old storefront.
	GetTaxRatePercentageValue(obj interface{}, country interface{}, previousValue interface{}) *PluginMethodNotImplemented
	// Trigger when user is created.
	// Overwrite this method if you need to trigger specific logic after a user is created.
	CustomerCreated(customer account.User, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when user is updated.
	// Overwrite this method if you need to trigger specific logic after a user is
	// updated.
	CustomerUpdated(customer account.User, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when product is created.
	// Overwrite this method if you need to trigger specific logic after a product is created.
	ProductCreated(product product_and_discount.Product, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when product is updated.
	// Overwrite this method if you need to trigger specific logic after a product is updated.
	ProductUpdated(product product_and_discount.Product, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when product is deleted.
	// Overwrite this method if you need to trigger specific logic after a product is deleted.
	ProductDeleted(product product_and_discount.Product, variants []int, previousVale interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when product variant is created.
	// Overwrite this method if you need to trigger specific logic after a product
	// variant is created.
	ProductVariantCreated(productVariant product_and_discount.ProductVariant, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when product variant is updated.
	// Overwrite this method if you need to trigger specific logic after a product
	// variant is updated.
	ProductVariantUpdated(variant product_and_discount.ProductVariant, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when product variant is deleted.
	// Overwrite this method if you need to trigger specific logic after a product
	// variant is deleted.
	ProductVariantDeleted(productVariant product_and_discount.ProductVariant, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// ProductVariantOutOfStock triggered when a product variant is out of stock
	ProductVariantOutOfStock(stock warehouse.Stock, defaultValue interface{}) *PluginMethodNotImplemented
	// ProductVariantBackInStock is triggered when a product is available again in stock
	ProductVariantBackInStock(stock warehouse.Stock, defaultValue interface{}) *PluginMethodNotImplemented
	// Trigger when order is fully paid.
	// Overwrite this method if you need to trigger specific logic when an order is
	// fully paid.
	OrderFullyPaid(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when order is updated.
	// Overwrite this method if you need to trigger specific logic when an order is changed.
	OrderUpdated(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when order is cancelled.
	// Overwrite this method if you need to trigger specific logic when an order is
	// canceled.
	OrderCancelled(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when order is fulfilled.
	// Overwrite this method if you need to trigger specific logic when an order is fulfilled.
	OrderFulfilled(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when fulfillemnt is created.
	// Overwrite this method if you need to trigger specific logic when a fulfillment is created.
	FulfillmentCreated(fulfillment order.Fulfillment, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when fulfillemnt is cancelled.
	// Overwrite this method if you need to trigger specific logic when a fulfillment is cancelled.
	FulfillmentCanceled(fulfillment order.Fulfillment, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when checkout is created.
	// Overwrite this method if you need to trigger specific logic when a checkout is created.
	CheckoutCreated(checkOut checkout.Checkout, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when checkout is updated.
	// Overwrite this method if you need to trigger specific logic when a checkout is updated.
	CheckoutUpdated(checkOut checkout.Checkout, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when page is updated.
	// Overwrite this method if you need to trigger specific logic when a page is updated.
	PageUpdated(page_ page.Page, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when page is created.
	// Overwrite this method if you need to trigger specific logic when a page is created.
	PageCreated(page_ page.Page, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Trigger when page is deleted.
	// Overwrite this method if you need to trigger specific logic when a page is deleted.
	PageDeleted(page_ page.Page, previousValue interface{}) (interface{}, *PluginMethodNotImplemented)
	// Triggered when ShopFetchTaxRates mutation is called.
	FetchTaxesData(previousValue interface{}) (bool, *PluginMethodNotImplemented)
	//
	InitializePayment(paymentData model.StringInterface, previousValue interface{}) (*payment.InitializedPaymentResponse, *PluginMethodNotImplemented)
	//
	AuthorizePayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented)
	//
	CapturePayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented)
	//
	VoidPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented)
	//
	RefundPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented)
	//
	ConfirmPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented)
	//
	ProcessPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented)
	//
	ListPaymentSources(customerID string, previousValue interface{}) ([]*payment.CustomerSource, *PluginMethodNotImplemented)
	//
	GetClientToken(tokenConfig payment.TokenConfig, previousValue interface{}) (string, *PluginMethodNotImplemented)
	//
	GetPaymentConfig(previousValue interface{}) ([]model.StringInterface, *PluginMethodNotImplemented)
	//
	GetSupportedCurrencies(previousValue interface{}) ([]string, *PluginMethodNotImplemented)
	//
	TokenIsRequiredAsPaymentInput(previousValue bool) (bool, *PluginMethodNotImplemented)
	//
	GetPaymentGateways(currency string, checkOut *checkout.Checkout, previousValue interface{}) ([]*payment.PaymentGateway, *PluginMethodNotImplemented)
	//
	UpdateConfigItems(configurationToUpdate []model.StringInterface, currentConfig []model.StringInterface) ([]model.StringInterface, *PluginMethodNotImplemented)
	// Validate if provided configuration is correct.
	// Raise django.core.exceptions.ValidationError otherwise.
	ValidatePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration) (*model.AppError, *PluginMethodNotImplemented)
	// Trigger before plugin configuration will be saved.
	// Overwrite this method if you need to trigger specific logic before saving a
	// plugin configuration.
	PreSavePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration) (*model.AppError, *PluginMethodNotImplemented)
	//
	SavePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration, cleanedData model.StringInterface) (*plugins.PluginConfiguration, *model.AppError, *PluginMethodNotImplemented)
	// Append configuration structure to config from the database.
	// Database stores "key: value" pairs, the definition of fields should be declared
	// inside of the plugin. Based on this, the plugin will generate a structure of
	// configuration with current values and provide access to it via API.
	AppendConfigStructure(configuration PluginConfigurationType) (PluginConfigurationType, *PluginMethodNotImplemented)
	//
	UpdateConfigurationStructure(configuration []model.StringInterface) (PluginConfigurationType, *PluginMethodNotImplemented)
	//
	GetDefaultActive() (bool, *PluginMethodNotImplemented)
	//
	GetPluginConfiguration(configuration PluginConfigurationType) (PluginConfigurationType, *PluginMethodNotImplemented)
	//
	IsActive() bool
	ChannelId() string
	GetManifest() *PluginManifest
}
