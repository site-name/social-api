package plugin

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
)

const ErrorPluginbMethodNotImplemented = "app.plugin.method_not_implemented.app_error"

// PluginConfig contains configurations to initialize a new plugin
type PluginConfig struct {
	Active        bool
	ChannelID     string
	Configuration model.StringInterfaces
	Manager       *PluginManager
	Manifest      *interfaces.PluginManifest
}

// type check
var _ interfaces.BasePluginInterface = (*BasePlugin)(nil)

// every newly added plugins must inherit from this one
type BasePlugin struct {
	Manifest      *interfaces.PluginManifest
	Active        bool
	ChannelID     string
	Configuration model.StringInterfaces
	Manager       *PluginManager
}

func NewBasePlugin(cfg *PluginConfig) *BasePlugin {
	return &BasePlugin{
		Active:        cfg.Active,
		ChannelID:     cfg.ChannelID,
		Configuration: cfg.Configuration,
		Manager:       cfg.Manager,
		Manifest:      cfg.Manifest,
	}
}

func (b *BasePlugin) IsActive() bool {
	return b.Active
}

func (b *BasePlugin) ChannelId() string {
	return b.ChannelID
}

func (b *BasePlugin) GetManifest() *interfaces.PluginManifest {
	return b.Manifest
}

func (b *BasePlugin) GetConfiguration() model.StringInterfaces {
	return b.Configuration
}

func (b *BasePlugin) SetConfiguration(config model.StringInterfaces) {
	b.Configuration = config
}

func (b *BasePlugin) SetActive(active bool) {
	b.Active = active
}

func (b *BasePlugin) String() string {
	return b.Manifest.PluginName
}

func (b *BasePlugin) ExternalObtainAccessTokens(data model_types.JSONString, request *http.Request, previousValue model.ExternalAccessTokens) (*model.ExternalAccessTokens, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ExternalObtainAccessTokens", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalRefresh(data model_types.JSONString, request *http.Request, previousValue model.ExternalAccessTokens) (*model.ExternalAccessTokens, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ExternalRefresh", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalLogout(data model_types.JSONString, request *http.Request, previousValue model_types.JSONString) *model_helper.AppError {
	return model_helper.NewAppError("ExternalLogout", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalVerify(data model_types.JSONString, request *http.Request, previousValue interfaces.AType) (*model.User, model_types.JSONString, *model_helper.AppError) {
	return nil, nil, model_helper.NewAppError("ExternalVerify", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AuthenticateUser(request *http.Request, previousValue any) (*model.User, *model_helper.AppError) {
	return nil, model_helper.NewAppError("AuthenticateUser", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) Webhook(request *http.Request, path string, previousValue http.Response) (*http.Response, *model_helper.AppError) {
	return nil, model_helper.NewAppError("Webhook", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) Notify(event string, payload model_types.JSONString, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("Notify", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ChangeUserAddress(address model.Address, addressType *model.AddressTypeEnum, user *model.User, previousValue model.Address) (*model.Address, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ChangeUserAddress", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutTotal(checkoutInfo model_helper.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CalculateCheckoutTotal", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutShipping(checkoutInfo model_helper.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CalculateCheckoutShipping", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderShipping(orDer *model.Order, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CalculateOrderShipping", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineTotal(checkoutInfo model_helper.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model_helper.CheckoutLineInfo, address *model.Address, discounts []*model_helper.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CalculateCheckoutLineTotal", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineTotal(orDer *model.Order, orderLine *model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CalculateOrderLineTotal", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineUnitPrice(checkoutInfo model_helper.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model_helper.CheckoutLineInfo, address *model.Address, discounts []*model_helper.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CalculateCheckoutLineUnitPrice", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineUnit(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CalculateOrderLineUnit", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetCheckoutLineTaxRate(checkoutInfo *model_helper.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model_helper.CheckoutLineInfo, address *model.Address, discounts []*model_helper.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetCheckoutLineTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetOrderLineTaxRate(orDer model.Order, product model.Product, variant model.ProductVariant, address *model.Address, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetOrderLineTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetCheckoutShippingTaxRate(checkoutInfo model_helper.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetCheckoutShippingTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetOrderShippingTaxRate(orDer model.Order, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetOrderShippingTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxRateTypeChoices(previousValue []*model.TaxType) ([]*model.TaxType, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetTaxRateTypeChoices", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ShowTaxesOnStorefront(previousValue bool) (bool, *model_helper.AppError) {
	return false, model_helper.NewAppError("ShowTaxesOnStorefront", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToShipping(price goprices.Money, shippingAddress model.Address, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ApplyTaxesToShipping", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToProduct(product model.Product, price goprices.Money, country model.CountryCode, previousVlaue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ApplyTaxesToProduct", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PreprocessOrderCreation(checkoutInfo model_helper.CheckoutInfo, discounts []*model_helper.DiscountInfo, lines model.CheckoutLineInfos, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("PreprocessOrderCreation", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderCreated(orDer model.Order, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("OrderCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderConfirmed(orDer model.Order, previousValue any) (any, *model_helper.AppError) {
	return b.OrderCreated(orDer, previousValue)
}

func (b *BasePlugin) SaleCreated(sale model.Sale, currentCatalogue model.NodeCatalogueInfo, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("SaleCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) SaleDeleted(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("SaleDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) SaleUpdated(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, currentCatalogue model.NodeCatalogueInfo, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("SaleUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceRequest(orDer model.Order, inVoice model.Invoice, number string, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("InvoiceRequest", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceDelete(inVoice model.Invoice, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("InvoiceDelete", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceSent(inVoice model.Invoice, email string, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("InvoiceSent", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AssignTaxCodeToObjectMeta(obj any, taxCode string, previousValue model.TaxType) (*model.TaxType, *model_helper.AppError) {
	return nil, model_helper.NewAppError("AssignTaxCodeToObjectMeta", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxRatePercentageValue(obj any, country string, previousValue decimal.Decimal) (*decimal.Decimal, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetTaxRatePercentageValue", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CustomerCreated(customer model.User, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CustomerCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CustomerUpdated(customer model.User, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CustomerUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductCreated(product model.Product, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ProductCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductUpdated(product model.Product, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ProductUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductDeleted(product model.Product, variants []int, previousVale any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ProductDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantCreated(productVariant model.ProductVariant, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ProductVariantCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantUpdated(variant model.ProductVariant, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ProductVariantUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantOutOfStock(stock model.Stock, defaultValue any) *model_helper.AppError {
	return model_helper.NewAppError("ProductVariantOutOfStock", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantBackInStock(stock model.Stock, defaultValue any) *model_helper.AppError {
	return model_helper.NewAppError("ProductVariantBackInStock", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantDeleted(productVariant model.ProductVariant, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ProductVariantDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderFullyPaid(orDer model.Order, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("OrderFullyPaid", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderUpdated(orDer model.Order, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("OrderUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderCancelled(orDer model.Order, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("OrderCancelled", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderFulfilled(orDer model.Order, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("OrderFulfilled", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderCreated(orDer model.Order, defaultValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("DraftOrderCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderUpdated(orDer model.Order, defaultValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("DraftOrderUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderDeleted(orDer model.Order, defaultValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("DraftOrderDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FulfillmentCreated(fulfillment model.Fulfillment, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("FulfillmentCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FulfillmentCanceled(fulfillment model.Fulfillment, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("FulfillmentCanceled", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckoutCreated(checkOut model.Checkout, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CheckoutCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckoutUpdated(checkOut model.Checkout, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CheckoutUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageUpdated(page_ model.Page, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("PageUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageCreated(page_ model.Page, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("PageCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageDeleted(page_ model.Page, previousValue any) (any, *model_helper.AppError) {
	return nil, model_helper.NewAppError("PageDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FetchTaxesData(previousValue bool) (bool, *model_helper.AppError) {
	return false, model_helper.NewAppError("FetchTaxesData", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InitializePayment(paymentData model_types.JSONString, previousValue any) (*model.InitializedPaymentResponse, *model_helper.AppError) {
	return nil, model_helper.NewAppError("InitializePayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AuthorizePayment(paymentInformation model.PaymentData, previousValue any) (*model.GatewayResponse, *model_helper.AppError) {
	return nil, model_helper.NewAppError("AuthorizePayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CapturePayment(paymentInformation model.PaymentData, previousValue any) (*model.GatewayResponse, *model_helper.AppError) {
	return nil, model_helper.NewAppError("CapturePayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) VoidPayment(paymentInformation model.PaymentData, previousValue any) (*model.GatewayResponse, *model_helper.AppError) {
	return nil, model_helper.NewAppError("VoidPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) RefundPayment(paymentInformation model.PaymentData, previousValue any) (*model.GatewayResponse, *model_helper.AppError) {
	return nil, model_helper.NewAppError("RefundPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ConfirmPayment(paymentInformation model.PaymentData, previousValue any) (*model.GatewayResponse, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ConfirmPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProcessPayment(paymentInformation model.PaymentData, previousValue any) (*model.GatewayResponse, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ProcessPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ListPaymentSources(customerID string, previousValue any) ([]*model.CustomerSource, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ListPaymentSources", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetClientToken(tokenConfig model.TokenConfig, previousValue any) (string, *model_helper.AppError) {
	return "", model_helper.NewAppError("GetClientToken", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetPaymentConfig(previousValue any) ([]model_types.JSONString, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetPaymentConfig", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetSupportedCurrencies(previousValue any) ([]string, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetSupportedCurrencies", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxCodeFromObjectMeta(obj any, previousValue model.TaxType) (*model.TaxType, *model_helper.AppError) {
	return nil, model_helper.NewAppError("GetTaxCodeFromObjectMeta", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) TokenIsRequiredAsPaymentInput(previousValue bool) (bool, *model_helper.AppError) {
	return previousValue, nil
}

func (b *BasePlugin) GetPaymentGateways(currency string, checkOut *model.Checkout, previousValue any) ([]*model.PaymentGateway, *model_helper.AppError) {
	paymentConfig, notImplt := b.GetPaymentConfig(previousValue)
	if notImplt != nil {
		paymentConfig = []model_types.JSONString{}
	}

	var currencies util.AnyArray[string]
	currencies, notImplt = b.GetSupportedCurrencies(previousValue)
	if notImplt != nil {
		currencies = []string{}
	}

	if currency != "" && !currencies.Contains(currency) {
		return []*model.PaymentGateway{}, nil
	}

	return []*model.PaymentGateway{
		{
			Id:         b.Manifest.PluginID,
			Name:       b.Manifest.PluginName,
			Config:     paymentConfig,
			Currencies: currencies,
		},
	}, nil
}

func (b *BasePlugin) ExternalAuthenticationUrl(data model_types.JSONString, request *http.Request, previousValue model_types.JSONString) (model_types.JSONString, *model_helper.AppError) {
	return nil, model_helper.NewAppError("ExternalAuthenticationUrl", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckPluginId(pluginID string) bool {
	return b.Manifest.PluginID == pluginID
}

func (b *BasePlugin) GetDefaultActive() (bool, *model_helper.AppError) {
	return b.Manifest.DefaultActive, nil
}

func (b *BasePlugin) UpdateConfigurationStructure(config []model_types.JSONString) (model.StringInterfaces, *model_helper.AppError) {
	var updatedConfiguration []model_types.JSONString

	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model_types.JSONString)
	}

	desiredConfigKeysMap := map[string]struct{}{}
	for key := range configStructure {
		if _, exist := desiredConfigKeysMap[key]; !exist {
			desiredConfigKeysMap[key] = struct{}{}
		}
	}

	for _, configField := range config {
		if name, ok := configField["name"]; ok {
			if _, exist := desiredConfigKeysMap[name.(string)]; !exist {
				continue
			}
		}

		updatedConfiguration = append(updatedConfiguration, configField.DeepCopy())
	}

	configuredKeysMap := map[string]struct{}{}
	for _, cfg := range updatedConfiguration {
		strName := cfg["name"].(string) // name should exist

		if _, exist := configuredKeysMap[strName]; !exist {
			configuredKeysMap[strName] = struct{}{}
		}
	}

	missingKeysMap := map[string]struct{}{} // items reside in desiredConfigKeys but not in configuredKeys
	for key := range desiredConfigKeysMap {
		if _, exist := configuredKeysMap[key]; !exist {
			missingKeysMap[key] = struct{}{}
		}
	}

	if len(missingKeysMap) == 0 {
		return updatedConfiguration, nil
	}

	if len(b.Manifest.DefaultConfiguration) == 0 {
		return updatedConfiguration, nil
	}

	for _, item := range b.Manifest.DefaultConfiguration {
		if _, exist := missingKeysMap[item["name"].(string)]; exist {
			updatedConfiguration = append(updatedConfiguration, item.DeepCopy())
		}
	}

	return updatedConfiguration, nil
}

func (b *BasePlugin) GetPluginConfiguration(config model.StringInterfaces) (model.StringInterfaces, *model_helper.AppError) {
	if config == nil {
		config = model.StringInterfaces{}
	}

	config, _ = b.UpdateConfigurationStructure(config)

	var notImplt *model_helper.AppError
	if len(config) > 0 {
		config, notImplt = b.AppendConfigStructure(config)
		if notImplt != nil {
			return nil, notImplt
		}
	}

	return config, nil
}

// Append configuration structure to config from the database.
//
// Database stores "key: value" pairs, the definition of fields should be declared
// inside of the plugin. Based on this, the plugin will generate a structure of
// configuration with current values and provide access to it via API.
func (b *BasePlugin) AppendConfigStructure(config model.StringInterfaces) (model.StringInterfaces, *model_helper.AppError) {
	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model_types.JSONString)
	}

	fieldsWithoutStructure := []model_types.JSONString{}

	for _, configurationField := range config {
		structureToAdd, ok := configStructure[configurationField.Get("name", "").(string)]
		if ok && structureToAdd != nil {
			for key, value := range structureToAdd {
				configurationField[key] = value
			}
			continue
		}

		fieldsWithoutStructure = append(fieldsWithoutStructure, configurationField)
	}

	if len(fieldsWithoutStructure) > 0 {
		for _, field := range fieldsWithoutStructure {
			for idx, item := range config {
				if reflect.DeepEqual(field, item) {
					config = append(config[:idx], config[idx+1:]...)
				}
			}
		}
	}

	return config, nil
}

func (b *BasePlugin) UpdateConfigItems(configurationToUpdate []model_types.JSONString, currentConfig []model_types.JSONString) ([]model_types.JSONString, *model_helper.AppError) {
	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model_types.JSONString)
	}

	for _, configItem := range currentConfig {
		for _, configItemToUpdate := range configurationToUpdate {
			configItemName := configItemToUpdate.Get("name", "")

			if configItem.Get("name", "") == configItemName {
				newValue, ok3 := configItemToUpdate["value"]
				configStructureValue, ok4 := configStructure[configItemName.(string)]

				if !ok4 || configStructureValue == nil {
					configStructureValue = make(model_types.JSONString)
				}
				itemType, ok5 := configStructureValue["type"]

				newValueIsNotNullNorBoolean := ok3 && newValue != nil
				if newValueIsNotNullNorBoolean {
					_, newValueIsBoolean := newValue.(bool)
					newValueIsNotNullNorBoolean = newValueIsNotNullNorBoolean && !newValueIsBoolean
				}

				if ok5 &&
					itemType != nil &&
					itemType.(interfaces.ConfigurationTypeField) == interfaces.BOOLEAN &&
					newValueIsNotNullNorBoolean {
					newValue = strings.ToLower(newValue.(string)) == "true"
				}

				if val, ok := itemType.(interfaces.ConfigurationTypeField); ok && val == interfaces.OUTPUT {
					// OUTPUT field is read only. No need to update it
					continue
				}

				configItem["value"] = newValue
			}
		}
	}

	// Get new keys that don't exist in currentConfig and extend it:
	currentConfigKeys := util.AnyArray[string]{}
	for _, cField := range currentConfig {
		currentConfigKeys = append(currentConfigKeys, cField["name"].(string))
	}

	configurationToUpdateDict := make(model_types.JSONString)
	for _, item := range configurationToUpdate {
		configurationToUpdateDict[item["name"].(string)] = item["value"]
	}
	configurationToUpdateDictKeys := lo.Keys(configurationToUpdateDict)

	for _, item := range configurationToUpdateDictKeys {
		if !currentConfigKeys.Contains(item) {
			if val, ok := configStructure[item]; !ok || val == nil {
				continue
			}

			currentConfig = append(currentConfig, model_types.JSONString{
				"name":  item,
				"value": configurationToUpdateDict[item],
			})
		}
	}

	return currentConfig, nil
}

func (b *BasePlugin) SavePluginConfiguration(pluginConfiguration *model.PluginConfiguration, cleanedData model_types.JSONString) (*model.PluginConfiguration, *model_helper.AppError) {
	currentConfig := pluginConfiguration.Configuration
	configurationToUpdate, ok := cleanedData["configuration"]

	if ok && configurationToUpdate != nil {
		pluginConfiguration.Configuration, _ = b.UpdateConfigItems(configurationToUpdate.([]model_types.JSONString), currentConfig)
	}

	if active, ok := cleanedData["active"]; ok && active != nil {
		pluginConfiguration.Active = active.(bool)
	}

	appErr := b.ValidatePluginConfiguration(pluginConfiguration)
	if appErr != nil {
		return nil, appErr
	}
	appErr = b.PreSavePluginConfiguration(pluginConfiguration)
	if appErr != nil {
		return nil, appErr
	}

	pluginConfiguration, appErr = b.Manager.Srv.PluginService().UpsertPluginConfiguration(pluginConfiguration)
	if appErr != nil {
		return nil, appErr
	}

	if len(pluginConfiguration.Configuration) > 0 {
		pluginConfiguration.Configuration, appErr = b.AppendConfigStructure(pluginConfiguration.Configuration)
		if appErr != nil {
			return nil, appErr
		}
	}

	return pluginConfiguration, nil
}

func (b *BasePlugin) ValidatePluginConfiguration(pluginConfiguration *model.PluginConfiguration) *model_helper.AppError {
	return model_helper.NewAppError("ValidatePluginConfiguration", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PreSavePluginConfiguration(pluginConfiguration *model.PluginConfiguration) *model_helper.AppError {
	return model_helper.NewAppError("PreSavePluginConfiguration", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}
