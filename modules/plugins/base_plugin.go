package plugins

import (
	"net/http"
	"reflect"
	"strings"

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
	"github.com/sitename/sitename/modules/util"
)

var (
	_ BasePluginInterface = (*BasePlugin)(nil)
)

type NewPluginConfig struct {
	Active        bool
	ChannelID     string
	Configuration PluginConfigurationType
	Manager       *PluginManager
}

type BasePlugin struct {
	Manifest *PluginManifest

	Active        bool
	ChannelID     string
	Configuration PluginConfigurationType
	Manager       *PluginManager
}

func NewBasePlugin(cfg *NewPluginConfig) *BasePlugin {
	manifest := &PluginManifest{
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

func (b *BasePlugin) String() string {
	return b.Manifest.PluginName
}

func (b *BasePlugin) ExternalObtainAccessTokens(data model.StringInterface, request *http.Request, previousValue interface{}) (*ExternalAccessToken, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ExternalRefresh(data model.StringInterface, request *http.Request, previousValue interface{}) (*ExternalAccessToken, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ExternalLogout(data model.StringInterface, request *http.Request, previousValue interface{}) *PluginMethodNotImplemented {
	return new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ExternalVerify(data model.StringInterface, request *http.Request, previousValue interface{}) (*account.User, model.StringInterface, *PluginMethodNotImplemented) {
	return nil, nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) AuthenticateUser(request *http.Request, previousValue interface{}) (*account.User, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) Webhook(request *http.Request, path string, previousValue interface{}) (http.Response, *PluginMethodNotImplemented) {
	return http.Response{}, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) Notify(event interface{}, payload model.StringInterface, previousValue interface{}) *PluginMethodNotImplemented {
	return new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ChangeUserAddress(address *account.Address, addressType string, user *account.User, previousValue *account.Address) (*account.Address, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutShipping(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CalculateOrderShipping(orDer *order.Order, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineTotal(orDer *order.Order, orderLine *order.OrderLine, variant product_and_discount.ProductVariant, product product_and_discount.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CalculateCheckoutLineUnitPrice(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CalculateOrderLineUnit(orDer order.Order, orderLine order.OrderLine, variant product_and_discount.ProductVariant, product product_and_discount.Product, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetCheckoutLineTaxRate(checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetOrderLineTaxRate(orDer order.Order, product product_and_discount.Product, variant product_and_discount.ProductVariant, address *account.Address, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetCheckoutShippingTaxRate(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetOrderShippingTaxRate(orDer order.Order, previousValue decimal.Decimal) (*decimal.Decimal, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetTaxRateTypeChoices(previousValue []*model.TaxType) ([]*model.TaxType, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ShowTaxesOnStorefront(previousValue bool) (bool, *PluginMethodNotImplemented) {
	return false, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToShipping(price goprices.Money, shippingAddress account.Address, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ApplyTaxesToProduct(product product_and_discount.Product, price goprices.Money, country string, previousVlaue goprices.TaxedMoney) (*goprices.TaxedMoney, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) PreprocessOrderCreation(checkoutInfo checkout.CheckoutInfo, discounts []*product_and_discount.DiscountInfo, lines checkout.CheckoutLineInfos, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) OrderCreated(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) OrderConfirmed(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return b.OrderCreated(orDer, previousValue)
}

func (b *BasePlugin) SaleCreated(sale product_and_discount.Sale, currentCatalogue product_and_discount.NodeCatalogueInfo, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) SaleDeleted(sale product_and_discount.Sale, previousCatalogue product_and_discount.NodeCatalogueInfo, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) SaleUpdated(sale product_and_discount.Sale, previousCatalogue product_and_discount.NodeCatalogueInfo, currentCatalogue product_and_discount.NodeCatalogueInfo, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) InvoiceRequest(orDer order.Order, inVoice invoice.Invoice, number string, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) InvoiceDelete(inVoice invoice.Invoice, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) InvoiceSent(inVoice invoice.Invoice, email string, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) AssignTaxCodeToObjectMeta(obj interface{}, previousValue model.TaxType) (*model.TaxType, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetTaxRatePercentageValue(obj interface{}, country interface{}, previousValue interface{}) *PluginMethodNotImplemented {
	return new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CustomerCreated(customer account.User, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CustomerUpdated(customer account.User, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductCreated(product product_and_discount.Product, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductUpdated(product product_and_discount.Product, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductDeleted(product product_and_discount.Product, variants []int, previousVale interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductVariantCreated(productVariant product_and_discount.ProductVariant, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductVariantUpdated(variant product_and_discount.ProductVariant, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductVariantOutOfStock(stock warehouse.Stock, defaultValue interface{}) *PluginMethodNotImplemented {
	return new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductVariantBackInStock(stock warehouse.Stock, defaultValue interface{}) *PluginMethodNotImplemented {
	return new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProductVariantDeleted(productVariant product_and_discount.ProductVariant, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) OrderFullyPaid(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) OrderUpdated(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) OrderCancelled(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) OrderFulfilled(orDer order.Order, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) DraftOrderCreated(orDer order.Order, defaultValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) DraftOrderUpdated(orDer order.Order, defaultValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) DraftOrderDeleted(orDer order.Order, defaultValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) FulfillmentCreated(fulfillment order.Fulfillment, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) FulfillmentCanceled(fulfillment order.Fulfillment, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CheckoutCreated(checkOut checkout.Checkout, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CheckoutUpdated(checkOut checkout.Checkout, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) PageUpdated(page_ page.Page, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) PageCreated(page_ page.Page, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) PageDeleted(page_ page.Page, previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) FetchTaxesData(previousValue interface{}) (bool, *PluginMethodNotImplemented) {
	return false, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) InitializePayment(paymentData model.StringInterface, previousValue interface{}) (*payment.InitializedPaymentResponse, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) AuthorizePayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CapturePayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) VoidPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) RefundPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ConfirmPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ProcessPayment(paymentInformation payment.PaymentData, previousValue interface{}) (*payment.GatewayResponse, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) ListPaymentSources(customerID string, previousValue interface{}) ([]*payment.CustomerSource, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetClientToken(tokenConfig interface{}, previousValue interface{}) (string, *PluginMethodNotImplemented) {
	return "", new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetPaymentConfig(previousValue interface{}) ([]model.StringInterface, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) GetSupportedCurrencies(previousValue interface{}) ([]string, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) TokenIsRequiredAsPaymentInput(previousValue interface{}) (interface{}, *PluginMethodNotImplemented) {
	return previousValue, nil
}

func (b *BasePlugin) GetPaymentGateways(currency string, checkOut *checkout.Checkout, previousValue interface{}) ([]*payment.PaymentGateway, *PluginMethodNotImplemented) {
	paymentConfig, notImplt := b.GetPaymentConfig(previousValue)
	if notImplt != nil {
		paymentConfig = []model.StringInterface{}
	}

	currencies, notImplt := b.GetSupportedCurrencies(previousValue)
	if notImplt != nil {
		currencies = []string{}
	}

	if currency != "" && !util.StringInSlice(currency, currencies) {
		return []*payment.PaymentGateway{}, nil
	}

	return []*payment.PaymentGateway{
		{
			Id:         b.Manifest.PluginID,
			Name:       b.Manifest.PluginName,
			Config:     paymentConfig,
			Currencies: currencies,
		},
	}, nil
}

func (b *BasePlugin) ExternalAuthenticationUrl(data model.StringInterface, request *http.Request, previousValue interface{}) (model.StringInterface, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CheckPluginId(pluginID string) (bool, *PluginMethodNotImplemented) {
	return b.Manifest.PluginID == pluginID, nil
}

func (b *BasePlugin) GetDefaultActive() (bool, *PluginMethodNotImplemented) {
	return b.Manifest.DefaultActive, nil
}

func (b *BasePlugin) UpdateConfigurationStructure(config []model.StringInterface) (PluginConfigurationType, *PluginMethodNotImplemented) {
	var updatedConfiguration []model.StringInterface

	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model.StringInterface)
	}

	desiredConfigKeys := []string{}
	for key := range configStructure {
		desiredConfigKeys = append(desiredConfigKeys, key)
	}
	desiredConfigKeys = util.RemoveDuplicatesFromStringArray(desiredConfigKeys)

	for _, configField := range config {
		if name, ok := configField["name"]; ok && !util.StringInSlice(name.(string), desiredConfigKeys) {
			continue
		}

		updatedConfiguration = append(updatedConfiguration, model.CopyStringInterface(configField))
	}

	configuredKeys := []string{}
	for _, cfg := range updatedConfiguration {
		configuredKeys = append(configuredKeys, cfg["name"].(string)) // name should exist
	}
	configuredKeys = util.RemoveDuplicatesFromStringArray(configuredKeys)

	missingKeys := []string{}
	for _, value := range desiredConfigKeys {
		if !util.StringInSlice(value, configuredKeys) {
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
		if util.StringInSlice(item["name"].(string), missingKeys) {
			updatedValues = append(updatedValues, model.CopyStringInterface(item))
		}
	}

	if len(updatedValues) > 0 {
		updatedConfiguration = append(updatedConfiguration, updatedValues...)
	}

	return updatedConfiguration, nil
}

func (b *BasePlugin) GetPluginConfiguration(config PluginConfigurationType) (PluginConfigurationType, *PluginMethodNotImplemented) {
	if config == nil {
		config = PluginConfigurationType{}
	}

	config, _ = b.UpdateConfigurationStructure(config)

	var notImplt *PluginMethodNotImplemented
	if len(config) > 0 {
		config, notImplt = b.AppendConfigStructure(config)
		if notImplt != nil {
			return nil, notImplt
		}
	}

	return config, nil
}

func (b *BasePlugin) AppendConfigStructure(config PluginConfigurationType) (PluginConfigurationType, *PluginMethodNotImplemented) {
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

func (b *BasePlugin) UpdateConfigItems(configurationToUpdate []model.StringInterface, currentConfig []model.StringInterface) ([]model.StringInterface, *PluginMethodNotImplemented) {
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
					itemType.(ConfigurationTypeField) == BOOLEAN &&
					newValueIsNotNullNorBoolean {
					newValue = strings.ToLower(newValue.(string)) == "true"
				}

				if val, ok := itemType.(ConfigurationTypeField); ok && val == OUTPUT {
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
	currentConfigKeys = util.RemoveDuplicatesFromStringArray(currentConfigKeys)

	configurationToUpdateDict := make(model.StringInterface)
	configurationToUpdateDictKeys := []string{}

	for _, item := range configurationToUpdate {
		configurationToUpdateDict[item["name"].(string)] = item["value"]
		configurationToUpdateDictKeys = append(configurationToUpdateDictKeys, item["name"].(string))
	}
	configurationToUpdateDictKeys = util.RemoveDuplicatesFromStringArray(configurationToUpdateDictKeys)

	for _, item := range configurationToUpdateDictKeys {
		if !util.StringInSlice(item, currentConfigKeys) {
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

func (b *BasePlugin) SavePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration, cleanedData model.StringInterface) (*plugins.PluginConfiguration, *model.AppError, *PluginMethodNotImplemented) {
	currentConfig := pluginConfiguration.Configuration
	configurationToUpdate, ok := cleanedData["configuration"]

	if ok && configurationToUpdate != nil {
		pluginConfiguration.Configuration, _ = b.UpdateConfigItems(configurationToUpdate.([]model.StringInterface), currentConfig)
	}

	if active, ok := cleanedData["active"]; ok && active != nil {
		pluginConfiguration.Active = active.(bool)
	}

	appErr, notImplt := b.ValidatePluginConfiguration(pluginConfiguration)
	if notImplt != nil {
		return nil, nil, notImplt
	}
	if appErr != nil {
		return nil, appErr, nil
	}
	appErr, notImplt = b.PreSavePluginConfiguration(pluginConfiguration)
	if notImplt != nil {
		return nil, nil, notImplt
	}
	if appErr != nil {
		return nil, appErr, nil
	}

	pluginConfiguration, appErr = b.Manager.srv.PluginService().UpsertPluginConfiguration(pluginConfiguration)
	if appErr != nil {
		return nil, appErr, nil
	}

	if len(pluginConfiguration.Configuration) > 0 {
		pluginConfiguration.Configuration, notImplt = b.AppendConfigStructure(pluginConfiguration.Configuration)
		if notImplt != nil {
			return nil, nil, notImplt
		}
	}

	return pluginConfiguration, nil, nil
}

func (b *BasePlugin) ValidatePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration) (*model.AppError, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) PreSavePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration) (*model.AppError, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}
