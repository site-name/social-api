package plugin

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

var (
	_ interfaces.BasePluginInterface = (*BasePlugin)(nil)
)

const ErrorPluginbMethodNotImplemented = "app.plugin.method_not_implemented.app_error"

type NewPluginConfig struct {
	Active        bool
	ChannelID     string
	Configuration interfaces.PluginConfigurationType
	Manager       *PluginManager
}

type BasePlugin struct {
	Manifest *interfaces.PluginManifest

	Active        bool
	ChannelID     string
	Configuration interfaces.PluginConfigurationType
	Manager       *PluginManager
}

func NewBasePlugin(cfg *NewPluginConfig) *BasePlugin {
	manifest := &interfaces.PluginManifest{
		ConfigStructure:         make(map[string]model.StringInterface),
		ConfigurationPerChannel: true,
		DefaultConfiguration:    []model.StringInterface{},
	}

	return &BasePlugin{
		Manifest:      manifest,
		Active:        cfg.Active,
		ChannelID:     cfg.ChannelID,
		Configuration: cfg.Configuration,
		Manager:       cfg.Manager,
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
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalRefresh(data model.StringInterface, request *http.Request, previousValue model.ExternalAccessTokens) (*model.ExternalAccessTokens, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalLogout(data model.StringInterface, request *http.Request, previousValue model.StringInterface) *model.AppError {
	return model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ExternalVerify(data model.StringInterface, request *http.Request, previousValue interfaces.AType) (*model.User, model.StringInterface, *model.AppError) {
	return nil, nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AuthenticateUser(request *http.Request, previousValue interface{}) (*model.User, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) Webhook(request *http.Request, path string, previousValue http.Response) (*http.Response, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) Notify(event string, payload model.StringInterface, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ChangeUserAddress(address model.Address, addressType string, user *model.User, previousValue model.Address) (*model.Address, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutShipping(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderShipping(orDer *model.Order, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineTotal(orDer *model.Order, orderLine *model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineUnitPrice(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineUnit(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetCheckoutLineTaxRate(checkoutInfo *model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetOrderLineTaxRate(orDer model.Order, product model.Product, variant model.ProductVariant, address *model.Address, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetCheckoutShippingTaxRate(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetOrderShippingTaxRate(orDer model.Order, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxRateTypeChoices(previousValue []*model.TaxType) ([]*model.TaxType, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ShowTaxesOnStorefront(previousValue bool) (bool, *model.AppError) {
	return false, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToShipping(price goprices.Money, shippingAddress model.Address, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToProduct(product model.Product, price goprices.Money, country string, previousVlaue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PreprocessOrderCreation(checkoutInfo model.CheckoutInfo, discounts []*model.DiscountInfo, lines model.CheckoutLineInfos, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderCreated(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderConfirmed(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return b.OrderCreated(orDer, previousValue)
}

func (b *BasePlugin) SaleCreated(sale model.Sale, currentCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) SaleDeleted(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) SaleUpdated(sale model.Sale, previousCatalogue model.NodeCatalogueInfo, currentCatalogue model.NodeCatalogueInfo, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceRequest(orDer model.Order, inVoice model.Invoice, number string, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceDelete(inVoice model.Invoice, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InvoiceSent(inVoice model.Invoice, email string, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AssignTaxCodeToObjectMeta(obj interface{}, taxCode string, previousValue model.TaxType) (*model.TaxType, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxRatePercentageValue(obj interface{}, country string, previousValue decimal.Decimal) (*decimal.Decimal, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CustomerCreated(customer model.User, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CustomerUpdated(customer model.User, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductCreated(product model.Product, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductUpdated(product model.Product, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductDeleted(product model.Product, variants []int, previousVale interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantCreated(productVariant model.ProductVariant, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantUpdated(variant model.ProductVariant, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantOutOfStock(stock model.Stock, defaultValue interface{}) *model.AppError {
	return model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantBackInStock(stock model.Stock, defaultValue interface{}) *model.AppError {
	return model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProductVariantDeleted(productVariant model.ProductVariant, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderFullyPaid(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderUpdated(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderCancelled(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) OrderFulfilled(orDer model.Order, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderCreated(orDer model.Order, defaultValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderUpdated(orDer model.Order, defaultValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) DraftOrderDeleted(orDer model.Order, defaultValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FulfillmentCreated(fulfillment model.Fulfillment, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FulfillmentCanceled(fulfillment model.Fulfillment, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckoutCreated(checkOut model.Checkout, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CheckoutUpdated(checkOut model.Checkout, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageUpdated(page_ model.Page, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageCreated(page_ model.Page, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PageDeleted(page_ model.Page, previousValue interface{}) (interface{}, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) FetchTaxesData(previousValue bool) (bool, *model.AppError) {
	return false, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) InitializePayment(paymentData model.StringInterface, previousValue interface{}) (*model.InitializedPaymentResponse, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) AuthorizePayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) CapturePayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) VoidPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) RefundPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ConfirmPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ProcessPayment(paymentInformation model.PaymentData, previousValue interface{}) (*model.GatewayResponse, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) ListPaymentSources(customerID string, previousValue interface{}) ([]*model.CustomerSource, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetClientToken(tokenConfig model.TokenConfig, previousValue interface{}) (string, *model.AppError) {
	return "", model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetPaymentConfig(previousValue interface{}) ([]model.StringInterface, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetSupportedCurrencies(previousValue interface{}) ([]string, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) GetTaxCodeFromObjectMeta(obj interface{}, previousValue model.TaxType) (*model.TaxType, *model.AppError) {
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) TokenIsRequiredAsPaymentInput(previousValue bool) (bool, *model.AppError) {
	return previousValue, nil
}

func (b *BasePlugin) GetPaymentGateways(currency string, checkOut *model.Checkout, previousValue interface{}) ([]*model.PaymentGateway, *model.AppError) {
	paymentConfig, notImplt := b.GetPaymentConfig(previousValue)
	if notImplt != nil {
		paymentConfig = []model.StringInterface{}
	}

	currencies, notImplt := b.GetSupportedCurrencies(previousValue)
	if notImplt != nil {
		currencies = []string{}
	}

	if currency != "" && !util.ItemInSlice(currency, currencies) {
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
	return nil, model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
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

	desiredConfigKeys := []string{}
	for key := range configStructure {
		desiredConfigKeys = append(desiredConfigKeys, key)
	}
	desiredConfigKeys = util.Dedup(desiredConfigKeys)

	for _, configField := range config {
		if name, ok := configField["name"]; ok && !util.ItemInSlice(name.(string), desiredConfigKeys) {
			continue
		}

		updatedConfiguration = append(updatedConfiguration, configField.DeepCopy())
	}

	configuredKeys := []string{}
	for _, cfg := range updatedConfiguration {
		configuredKeys = append(configuredKeys, cfg["name"].(string)) // name should exist
	}
	configuredKeys = util.Dedup(configuredKeys)

	missingKeys := []string{}
	for _, value := range desiredConfigKeys {
		if !util.ItemInSlice(value, configuredKeys) {
			missingKeys = append(missingKeys, value)
		}
	}

	if len(missingKeys) == 0 {
		return updatedConfiguration, nil
	}

	if len(b.Manifest.DefaultConfiguration) == 0 {
		return updatedConfiguration, nil
	}

	updatedValues := []model.StringInterface{}
	for _, item := range b.Manifest.DefaultConfiguration {
		if util.ItemInSlice(item["name"].(string), missingKeys) {
			updatedValues = append(updatedValues, item.DeepCopy())
		}
	}

	if len(updatedValues) > 0 {
		updatedConfiguration = append(updatedConfiguration, updatedValues...)
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
		structureToAdd, ok := configStructure[configurationField["name"].(string)]
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

			configItemToUpdateName, ok1 := configItemToUpdate["name"]
			configItemName, ok2 := configItem["name"]

			if ok1 && ok2 && configItemToUpdateName == configItemName {

				newValue, ok3 := configItemToUpdate["value"]

				newValueIsNotNullNorBoolean := ok3 && newValue != nil
				if newValueIsNotNullNorBoolean {
					_, newValueIsBoolean := newValue.(bool)
					newValueIsNotNullNorBoolean = newValueIsNotNullNorBoolean && !newValueIsBoolean
				}

				configStructureValue, ok4 := configStructure[configItemToUpdateName.(string)]

				if !ok4 || configStructureValue == nil {
					configStructureValue = make(model.StringInterface)
				}

				itemType, ok5 := configStructureValue["type"]

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
	currentConfigKeys := []string{}
	for _, cField := range currentConfig {
		currentConfigKeys = append(currentConfigKeys, cField["name"].(string))
	}
	currentConfigKeys = util.Dedup(currentConfigKeys)

	configurationToUpdateDict := make(model.StringInterface)
	configurationToUpdateDictKeys := []string{}

	for _, item := range configurationToUpdate {
		configurationToUpdateDict[item["name"].(string)] = item["value"]
		configurationToUpdateDictKeys = append(configurationToUpdateDictKeys, item["name"].(string))
	}
	configurationToUpdateDictKeys = util.Dedup(configurationToUpdateDictKeys)

	for _, item := range configurationToUpdateDictKeys {
		if !util.ItemInSlice(item, currentConfigKeys) {
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
	return model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}

func (b *BasePlugin) PreSavePluginConfiguration(pluginConfiguration *model.PluginConfiguration) *model.AppError {
	return model.NewAppError("", ErrorPluginbMethodNotImplemented, nil, "", http.StatusNotImplemented)
}
