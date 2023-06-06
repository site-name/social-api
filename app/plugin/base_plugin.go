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
	"github.com/sitename/sitename/modules/util"
)

const ErrorPluginbMethodNotImplemented = "app.plugin.method_not_implemented.app_error"

// PluginConfig contains configurations to initialize a new plugin
type PluginConfig struct {
	Active        bool
	ChannelID     string
	Configuration interfaces.PluginConfigurationType
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
	Configuration interfaces.PluginConfigurationType
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

func (b *BasePlugin) GetConfiguration() interfaces.PluginConfigurationType {
	return b.Configuration
}

func (b *BasePlugin) SetConfiguration(config interfaces.PluginConfigurationType) {
	b.Configuration = config
}

func (b *BasePlugin) SetActive(active bool) {
	b.Active = active
}

func (b *BasePlugin) String() string {
	return b.Manifest.PluginName
}

func (b *BasePlugin) ExternalObtainAccessTokens(data model.StringInterface, request *http.Request, previousValue model.ExternalAccessTokens) (*model.ExternalAccessTokens, *model.AppError) {
	return nil, model.NewAppError("ExternalObtainAccessTokens", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalRefresh(data model.StringInterface, request *http.Request, previousValue model.ExternalAccessTokens) (*model.ExternalAccessTokens, *model.AppError) {
	return nil, model.NewAppError("ExternalRefresh", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalLogout(data model.StringInterface, request *http.Request, previousValue model.StringInterface) *model.AppError {
	return model.NewAppError("ExternalLogout", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalVerify(data model.StringInterface, request *http.Request, previousValue interfaces.AType) (*model.User, model.StringInterface, *model.AppError) {
	return nil, nil, model.NewAppError("ExternalVerify", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AuthenticateUser(request *http.Request, previousValue interface{}) (*model.User, *model.AppError) {
	return nil, model.NewAppError("AuthenticateUser", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) Webhook(request *http.Request, path string, previousValue http.Response) (*http.Response, *model.AppError) {
	return nil, model.NewAppError("Webhook", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) Notify(event string, payload model.StringInterface, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("Notify", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ChangeUserAddress(address model.Address, addressType *model.AddressTypeEnum, user *model.User, previousValue model.Address) (*model.Address, *model.AppError) {
	return nil, model.NewAppError("ChangeUserAddress", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("CalculateCheckoutTotal", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutShipping(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("CalculateCheckoutShipping", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderShipping(orDer *model.Order, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("CalculateOrderShipping", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("CalculateCheckoutLineTotal", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineTotal(orDer *model.Order, orderLine *model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("CalculateOrderLineTotal", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineUnitPrice(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("CalculateCheckoutLineUnitPrice", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineUnit(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("CalculateOrderLineUnit", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetCheckoutLineTaxRate(checkoutInfo *model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("GetCheckoutLineTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetOrderLineTaxRate(orDer model.Order, product model.Product, variant model.ProductVariant, address *model.Address, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("GetOrderLineTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetCheckoutShippingTaxRate(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("GetCheckoutShippingTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetOrderShippingTaxRate(orDer model.Order, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("GetOrderShippingTaxRate", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxRateTypeChoices(previousValue []*model.TaxType) ([]*model.TaxType, *model.AppError) {
	return nil, model.NewAppError("GetTaxRateTypeChoices", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ShowTaxesOnStorefront(previousValue bool) (bool, *model.AppError) {
	return false, model.NewAppError("ShowTaxesOnStorefront", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToShipping(price goprices.Money, shippingAddress model.Address, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("ApplyTaxesToShipping", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToProduct(product model.Product, price goprices.Money, country model.CountryCode, previousVlaue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("ApplyTaxesToProduct", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PreprocessOrderCreation(checkoutInfo model.CheckoutInfo, discounts []*model.DiscountInfo, lines model.CheckoutLineInfos, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("PreprocessOrderCreation", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderCreated(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("OrderCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderConfirmed(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return b.OrderCreated(orDer, previousValue)
}

func (b *BasePlugin) SaleCreated(sale model.Sale, currentCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("SaleCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) SaleDeleted(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("SaleDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) SaleUpdated(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, currentCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("SaleUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceRequest(orDer model.Order, inVoice model.Invoice, number string, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("InvoiceRequest", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceDelete(inVoice model.Invoice, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("InvoiceDelete", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceSent(inVoice model.Invoice, email string, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("InvoiceSent", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AssignTaxCodeToObjectMeta(obj interface{}, taxCode string, previousValue model.TaxType) (*model.TaxType, *model.AppError) {
	return nil, model.NewAppError("AssignTaxCodeToObjectMeta", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxRatePercentageValue(obj interface{}, country string, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("GetTaxRatePercentageValue", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CustomerCreated(customer model.User, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("CustomerCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CustomerUpdated(customer model.User, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("CustomerUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductCreated(product model.Product, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("ProductCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductUpdated(product model.Product, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("ProductUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductDeleted(product model.Product, variants []int, previousVale interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("ProductDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantCreated(productVariant model.ProductVariant, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("ProductVariantCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantUpdated(variant model.ProductVariant, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("ProductVariantUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantOutOfStock(stock model.Stock, defaultValue interface{}) *model.AppError {
	return model.NewAppError("ProductVariantOutOfStock", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantBackInStock(stock model.Stock, defaultValue interface{}) *model.AppError {
	return model.NewAppError("ProductVariantBackInStock", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantDeleted(productVariant model.ProductVariant, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("ProductVariantDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderFullyPaid(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("OrderFullyPaid", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderUpdated(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("OrderUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderCancelled(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("OrderCancelled", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderFulfilled(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("OrderFulfilled", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderCreated(orDer model.Order, defaultValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("DraftOrderCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderUpdated(orDer model.Order, defaultValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("DraftOrderUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderDeleted(orDer model.Order, defaultValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("DraftOrderDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FulfillmentCreated(fulfillment model.Fulfillment, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("FulfillmentCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FulfillmentCanceled(fulfillment model.Fulfillment, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("FulfillmentCanceled", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckoutCreated(checkOut model.Checkout, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("CheckoutCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckoutUpdated(checkOut model.Checkout, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("CheckoutUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageUpdated(page_ model.Page, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("PageUpdated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageCreated(page_ model.Page, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("PageCreated", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageDeleted(page_ model.Page, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("PageDeleted", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FetchTaxesData(previousValue bool) (bool, *model.AppError) {
	return false, model.NewAppError("FetchTaxesData", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InitializePayment(paymentData model.StringInterface, previousValue interface{}) (*model.InitializedPaymentResponse, *model.AppError) {
	return nil, model.NewAppError("InitializePayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AuthorizePayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("AuthorizePayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CapturePayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("CapturePayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) VoidPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("VoidPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) RefundPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("RefundPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ConfirmPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("ConfirmPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProcessPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("ProcessPayment", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ListPaymentSources(customerID string, previousValue interface{}) ([]*model.CustomerSource, *model.AppError) {
	return nil, model.NewAppError("ListPaymentSources", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetClientToken(tokenConfig model.TokenConfig, previousValue interface{}) (string, *model.AppError) {
	return "", model.NewAppError("GetClientToken", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetPaymentConfig(previousValue interface{}) ([]model.StringInterface, *model.AppError) {
	return nil, model.NewAppError("GetPaymentConfig", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetSupportedCurrencies(previousValue interface{}) ([]string, *model.AppError) {
	return nil, model.NewAppError("GetSupportedCurrencies", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxCodeFromObjectMeta(obj interface{}, previousValue model.TaxType) (*model.TaxType, *model.AppError) {
	return nil, model.NewAppError("GetTaxCodeFromObjectMeta", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) TokenIsRequiredAsPaymentInput(previousValue bool) (bool, *model.AppError) {
	return previousValue, nil
}

func (b *BasePlugin) GetPaymentGateways(currency string, checkOut *model.Checkout, previousValue interface{}) ([]*model.PaymentGateway, *model.AppError) {
	paymentConfig, notImplt := b.GetPaymentConfig(previousValue)
	if notImplt != nil {
		paymentConfig = []model.StringInterface{}
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

func (b *BasePlugin) ExternalAuthenticationUrl(data model.StringInterface, request *http.Request, previousValue model.StringInterface) (model.StringInterface, *model.AppError) {
	return nil, model.NewAppError("ExternalAuthenticationUrl", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckPluginId(pluginID string) bool {
	return b.Manifest.PluginID == pluginID
}

func (b *BasePlugin) GetDefaultActive() (bool, *model.AppError) {
	return b.Manifest.DefaultActive, nil
}

func (b *BasePlugin) UpdateConfigurationStructure(config []model.StringInterface) (interfaces.PluginConfigurationType, *model.AppError) {
	var updatedConfiguration []model.StringInterface

	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model.StringInterface)
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

func (b *BasePlugin) GetPluginConfiguration(config interfaces.PluginConfigurationType) (interfaces.PluginConfigurationType, *model.AppError) {
	if config == nil {
		config = interfaces.PluginConfigurationType{}
	}

	config, _ = b.UpdateConfigurationStructure(config)

	var notImplt *model.AppError
	if len(config) > 0 {
		config, notImplt = b.AppendConfigStructure(config)
		if notImplt != nil {
			return nil, notImplt
		}
	}

	return config, nil
}

func (b *BasePlugin) AppendConfigStructure(config interfaces.PluginConfigurationType) (interfaces.PluginConfigurationType, *model.AppError) {
	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model.StringInterface)
	}

	fieldsWithoutStructure := []model.StringInterface{}

	for _, configurationField := range config {
		structureToAdd, ok := configStructure[configurationField.Get("name", "").(string)]
		if ok && structureToAdd != nil {
			for key, value := range structureToAdd {
				configurationField[key] = value
			}
		} else {
			fieldsWithoutStructure = append(fieldsWithoutStructure, configurationField)
		}
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

func (b *BasePlugin) UpdateConfigItems(configurationToUpdate []model.StringInterface, currentConfig []model.StringInterface) ([]model.StringInterface, *model.AppError) {
	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model.StringInterface)
	}

	for _, configItem := range currentConfig {
		for _, configItemToUpdate := range configurationToUpdate {
			configItemName := configItemToUpdate.Get("name", "")

			if configItem.Get("name", "") == configItemName {
				newValue, ok3 := configItemToUpdate["value"]
				configStructureValue, ok4 := configStructure[configItemName.(string)]

				if !ok4 || configStructureValue == nil {
					configStructureValue = make(model.StringInterface)
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

	configurationToUpdateDict := make(model.StringInterface)
	for _, item := range configurationToUpdate {
		configurationToUpdateDict[item["name"].(string)] = item["value"]
	}
	configurationToUpdateDictKeys := lo.Keys(configurationToUpdateDict)

	for _, item := range configurationToUpdateDictKeys {
		if !currentConfigKeys.Contains(item) {
			if val, ok := configStructure[item]; !ok || val == nil {
				continue
			}

			currentConfig = append(currentConfig, model.StringInterface{
				"name":  item,
				"value": configurationToUpdateDict[item],
			})
		}
	}

	return currentConfig, nil
}

func (b *BasePlugin) SavePluginConfiguration(pluginConfiguration *model.PluginConfiguration, cleanedData model.StringInterface) (*model.PluginConfiguration, *model.AppError) {
	currentConfig := pluginConfiguration.Configuration
	configurationToUpdate, ok := cleanedData["configuration"]

	if ok && configurationToUpdate != nil {
		pluginConfiguration.Configuration, _ = b.UpdateConfigItems(configurationToUpdate.([]model.StringInterface), currentConfig)
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

func (b *BasePlugin) ValidatePluginConfiguration(pluginConfiguration *model.PluginConfiguration) *model.AppError {
	return model.NewAppError("ValidatePluginConfiguration", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PreSavePluginConfiguration(pluginConfiguration *model.PluginConfiguration) *model.AppError {
	return model.NewAppError("PreSavePluginConfiguration", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}
