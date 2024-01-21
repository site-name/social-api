package interfaces

import (
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type ConfigurationTypeField string

const (
	STRING           ConfigurationTypeField = "String"
	MULTILINE        ConfigurationTypeField = "Multiline"
	BOOLEAN          ConfigurationTypeField = "Boolean"
	SECRET           ConfigurationTypeField = "Secret"
	SECRET_MULTILINE ConfigurationTypeField = "SecretMultiline"
	PASSWORD         ConfigurationTypeField = "Password"
	OUTPUT           ConfigurationTypeField = "OUTPUT"
)

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

type AType struct {
	User *model.User
	Data model.StringInterface
}

type BasePluginInterface interface {
	String() string
	// Check if given plugin_id matches with the PLUGIN_ID of this plugin
	CheckPluginId(pluginID string) bool
	// Handle authentication request responsible for obtaining access tokens.
	// Overwrite this method if the plugin handles authentication flow.
	ExternalAuthenticationUrl(data model.StringInterface, request *http.Request, previousValue model.StringInterface) (model.StringInterface, *model_helper.AppError)
	// Handle authentication request responsible for obtaining access tokens.
	// Overwrite this method if the plugin handles authentication flow.
	ExternalObtainAccessTokens(data model.StringInterface, request *http.Request, previousValue model.ExternalAccessTokens) (*model.ExternalAccessTokens, *model_helper.AppError)
	// Handle authentication refresh request.
	// Overwrite this method if the plugin handles authentication flow and supports
	// refreshing the access.
	ExternalRefresh(data model.StringInterface, request *http.Request, previousValue model.ExternalAccessTokens) (*model.ExternalAccessTokens, *model_helper.AppError)
	// Handle logout request.
	// Overwrite this method if the plugin handles logout flow.
	ExternalLogout(data model.StringInterface, request *http.Request, previousValue model.StringInterface) *model_helper.AppError
	// Verify the provided authentication data.
	// Overwrite this method if the plugin should validate the authentication data.
	ExternalVerify(data model.StringInterface, request *http.Request, previousValue AType) (*model.User, model.StringInterface, *model_helper.AppError)
	// Authenticate user which should be assigned to the request.
	// Overwrite this method if the plugin handles authentication flow.
	AuthenticateUser(request *http.Request, previousValue interface{}) (*model.User, *model_helper.AppError)
	// Handle received http request.
	// Overwrite this method if the plugin expects the incoming requests.
	Webhook(request *http.Request, path string, previousValue http.Response) (*http.Response, *model_helper.AppError)
	// Handle notification request.
	// Overwrite this method if the plugin is responsible for sending notifications.
	Notify(event string, payload model.StringInterface, previousValue interface{}) (interface{}, *model_helper.AppError)
	//
	ChangeUserAddress(address model.Address, addressType *model.AddressTypeEnum, user *model.User, previousValue model.Address) (*model.Address, *model_helper.AppError)
	// Calculate the total for checkout.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of a checkout total. Return TaxedMoney.
	CalculateCheckoutTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Calculate the shipping costs for model.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of shipping costs. Return TaxedMoney.
	CalculateCheckoutShipping(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Calculate the shipping costs for the order.
	// Update shipping costs in the order in case of changes in shipping address or
	// changes in draft order. Return TaxedMoney.
	CalculateOrderShipping(orDer *model.Order, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Calculate checkout line total.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of a checkout line total. Return TaxedMoney.
	CalculateCheckoutLineTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Calculate order line total.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of a order line total. Return TaxedMoney.
	CalculateOrderLineTotal(orDer *model.Order, orderLine *model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Calculate checkout line unit price
	CalculateCheckoutLineUnitPrice(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Calculate order line unit price.
	// Update order line unit price in the order in case of changes in draft order.
	// Return TaxedMoney.
	// Overwrite this method if you need to apply specific logic for the calculation
	// of an order line unit price.
	CalculateOrderLineUnit(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	//
	GetCheckoutLineTaxRate(checkoutInfo *model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError)
	//
	GetOrderLineTaxRate(orDer model.Order, product model.Product, variant model.ProductVariant, address *model.Address, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError)
	//
	GetCheckoutShippingTaxRate(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError)
	//
	GetOrderShippingTaxRate(orDer model.Order, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError)
	// Return list of all tax categories.
	// The returned list will be used to provide staff users with the possibility to
	// assign tax categories to a product. It can be used by tax plugins to properly
	// calculate taxes for products.
	// Overwrite this method in case your plugin provides a list of tax categories.
	GetTaxRateTypeChoices(previousValue []*model.TaxType) ([]*model.TaxType, *model_helper.AppError)
	// Define if storefront should add info about taxes to the price.
	// It is used only by the old storefront. The returned value determines if
	// storefront should append info to the price about "including/excluding X% VAT"
	ShowTaxesOnStorefront(previousValue bool) (bool, *model_helper.AppError)
	// Apply taxes to the shipping costs based on the shipping address.
	// Overwrite this method if you want to show available shipping methods with taxes.
	ApplyTaxesToShipping(price goprices.Money, shippingAddress model.Address, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Apply taxes to the product price based on the customer country.
	// Overwrite this method if you want to show products with taxes.
	ApplyTaxesToProduct(product model.Product, price goprices.Money, country model.CountryCode, previousVlaue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError)
	// Trigger directly before order creation.
	// Overwrite this method if you need to trigger specific logic before an order is created.
	PreprocessOrderCreation(checkoutInfo model.CheckoutInfo, discounts []*model.DiscountInfo, lines model.CheckoutLineInfos, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when order is created.
	// Overwrite this method if you need to trigger specific logic after an order is created.
	OrderCreated(orDer model.Order, previousValue interface{}) (interface{}, *model_helper.AppError)
	//
	DraftOrderCreated(orDer model.Order, defaultValue interface{}) (interface{}, *model_helper.AppError)
	//
	DraftOrderUpdated(orDer model.Order, defaultValue interface{}) (interface{}, *model_helper.AppError)
	//
	DraftOrderDeleted(orDer model.Order, defaultValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when order is confirmed by staff.
	// Overwrite this method if you need to trigger specific logic after an order is
	// confirmed.
	OrderConfirmed(orDer model.Order, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when sale is created.
	// Overwrite this method if you need to trigger specific logic after sale is created.
	SaleCreated(sale model.Sale, currentCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when sale is deleted.
	// Overwrite this method if you need to trigger specific logic after sale is deleted.
	SaleDeleted(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when sale is updated.
	// Overwrite this method if you need to trigger specific logic after sale is updated.
	SaleUpdated(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, currentCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when invoice creation starts.
	// Overwrite to create invoice with proper data, call invoice.update_invoice.
	InvoiceRequest(orDer model.Order, inVoice model.Invoice, number string, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger before invoice is deleted.
	// Perform any extra logic before the invoice gets deleted.
	// Note there is no need to run invoice.delete() as it will happen in mutation.
	InvoiceDelete(inVoice model.Invoice, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger after invoice is sent.
	InvoiceSent(inVoice model.Invoice, email string, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Return tax code from object meta.
	//
	// NOTE: obj can be 'Product' or 'ProductType'
	AssignTaxCodeToObjectMeta(obj interface{}, taxCode string, previousValue model.TaxType) (*model.TaxType, *model_helper.AppError)
	// Return tax code from object meta
	//
	// NOTE: obj must be either Product or ProductType
	GetTaxCodeFromObjectMeta(obj interface{}, previousValue model.TaxType) (*model.TaxType, *model_helper.AppError)
	// Return tax rate percentage value for a given tax rate type in a country.
	// It is used only by the old storefront.
	GetTaxRatePercentageValue(obj interface{}, country string, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError)
	// Trigger when user is created.
	// Overwrite this method if you need to trigger specific logic after a user is created.
	CustomerCreated(customer model.User, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when user is updated.
	// Overwrite this method if you need to trigger specific logic after a user is
	// updated.
	CustomerUpdated(customer model.User, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when product is created.
	// Overwrite this method if you need to trigger specific logic after a product is created.
	ProductCreated(product model.Product, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when product is updated.
	// Overwrite this method if you need to trigger specific logic after a product is updated.
	ProductUpdated(product model.Product, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when product is deleted.
	// Overwrite this method if you need to trigger specific logic after a product is deleted.
	ProductDeleted(product model.Product, variants []int, previousVale interface{}) (interface{}, *model_helper.AppError)
	// Trigger when product variant is created.
	// Overwrite this method if you need to trigger specific logic after a product
	// variant is created.
	ProductVariantCreated(productVariant model.ProductVariant, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when product variant is updated.
	// Overwrite this method if you need to trigger specific logic after a product
	// variant is updated.
	ProductVariantUpdated(variant model.ProductVariant, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when product variant is deleted.
	// Overwrite this method if you need to trigger specific logic after a product
	// variant is deleted.
	ProductVariantDeleted(productVariant model.ProductVariant, previousValue interface{}) (interface{}, *model_helper.AppError)
	// ProductVariantOutOfStock triggered when a product variant is out of stock
	ProductVariantOutOfStock(stock model.Stock, defaultValue interface{}) *model_helper.AppError
	// ProductVariantBackInStock is triggered when a product is available again in stock
	ProductVariantBackInStock(stock model.Stock, defaultValue interface{}) *model_helper.AppError
	// Trigger when order is fully paid.
	// Overwrite this method if you need to trigger specific logic when an order is
	// fully paid.
	OrderFullyPaid(orDer model.Order, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when order is updated.
	// Overwrite this method if you need to trigger specific logic when an order is changed.
	OrderUpdated(orDer model.Order, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when order is cancelled.
	// Overwrite this method if you need to trigger specific logic when an order is
	// canceled.
	OrderCancelled(orDer model.Order, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when order is fulfilled.
	// Overwrite this method if you need to trigger specific logic when an order is fulfilled.
	OrderFulfilled(orDer model.Order, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when fulfillemnt is created.
	// Overwrite this method if you need to trigger specific logic when a fulfillment is created.
	FulfillmentCreated(fulfillment model.Fulfillment, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when fulfillemnt is cancelled.
	// Overwrite this method if you need to trigger specific logic when a fulfillment is cancelled.
	FulfillmentCanceled(fulfillment model.Fulfillment, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when checkout is created.
	// Overwrite this method if you need to trigger specific logic when a checkout is created.
	CheckoutCreated(checkOut model.Checkout, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when checkout is updated.
	// Overwrite this method if you need to trigger specific logic when a checkout is updated.
	CheckoutUpdated(checkOut model.Checkout, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when page is updated.
	// Overwrite this method if you need to trigger specific logic when a page is updated.
	PageUpdated(page_ model.Page, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when page is created.
	// Overwrite this method if you need to trigger specific logic when a page is created.
	PageCreated(page_ model.Page, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Trigger when page is deleted.
	// Overwrite this method if you need to trigger specific logic when a page is deleted.
	PageDeleted(page_ model.Page, previousValue interface{}) (interface{}, *model_helper.AppError)
	// Triggered when ShopFetchTaxRates mutation is called.
	FetchTaxesData(previousValue bool) (bool, *model_helper.AppError)
	//
	InitializePayment(paymentData model.StringInterface, previousValue interface{}) (*model.InitializedPaymentResponse, *model_helper.AppError)
	//
	AuthorizePayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model_helper.AppError)
	//
	CapturePayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model_helper.AppError)
	//
	VoidPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model_helper.AppError)
	//
	RefundPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model_helper.AppError)
	//
	ConfirmPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model_helper.AppError)
	//
	ProcessPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model_helper.AppError)
	//
	ListPaymentSources(customerID string, previousValue interface{}) ([]*model.CustomerSource, *model_helper.AppError)
	//
	GetClientToken(tokenConfig model.TokenConfig, previousValue interface{}) (string, *model_helper.AppError)
	//
	GetPaymentConfig(previousValue interface{}) ([]model.StringInterface, *model_helper.AppError)
	//
	GetSupportedCurrencies(previousValue interface{}) ([]string, *model_helper.AppError)
	//
	TokenIsRequiredAsPaymentInput(previousValue bool) (bool, *model_helper.AppError)
	//
	GetPaymentGateways(currency string, checkOut *model.Checkout, previousValue interface{}) ([]*model.PaymentGateway, *model_helper.AppError)
	//
	UpdateConfigItems(configurationToUpdate []model.StringInterface, currentConfig []model.StringInterface) ([]model.StringInterface, *model_helper.AppError)
	// Validate if provided configuration is correct.
	// Raise django.core.exceptions.ValidationError otherwise.
	ValidatePluginConfiguration(pluginConfiguration *model.PluginConfiguration) *model_helper.AppError
	// Trigger before plugin configuration will be saved.
	// Overwrite this method if you need to trigger specific logic before saving a
	// plugin configuration.
	PreSavePluginConfiguration(pluginConfiguration *model.PluginConfiguration) *model_helper.AppError
	//
	SavePluginConfiguration(pluginConfiguration *model.PluginConfiguration, cleanedData model.StringInterface) (*model.PluginConfiguration, *model_helper.AppError)
	// Append configuration structure to config from the database.
	// Database stores "key: value" pairs, the definition of fields should be declared
	// inside of the plugin. Based on this, the plugin will generate a structure of
	// configuration with current values and provide access to it via API.
	AppendConfigStructure(configuration model.StringInterfaces) (model.StringInterfaces, *model_helper.AppError)
	//
	UpdateConfigurationStructure(configuration []model.StringInterface) (model.StringInterfaces, *model_helper.AppError)
	//
	GetDefaultActive() (bool, *model_helper.AppError)
	//
	GetPluginConfiguration(configuration model.StringInterfaces) (model.StringInterfaces, *model_helper.AppError)
	//
	IsActive() bool
	ChannelId() string
	GetManifest() *PluginManifest
	GetConfiguration() model.StringInterfaces
	SetActive(active bool)
	SetConfiguration(config model.StringInterfaces)
}
