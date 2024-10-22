package plugin

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

var _ interfaces.PluginManagerInterface = (*PluginManager)(nil)

type PluginManager struct {
	Srv        *app.Server
	allPlugins []interfaces.BasePluginInterface
}

// NewPluginManager returns a new plugin manager
func (s *ServicePlugin) newPluginManager() (interfaces.PluginManagerInterface, *model_helper.AppError) {
	manager := &PluginManager{
		Srv: s.srv,
	}

	channels, appErr := manager.Srv.
		Channel.
		ChannelsByOption(model_helper.ChannelFilterOptions{})
	if appErr != nil {
		return nil, appErr
	}

	// finds a list of plugin configs belong found channels
	pluginConfigsOfChannels, appErr := manager.Srv.
		Plugin.
		FilterPluginConfigurations(&model.PluginConfigurationFilterOptions{
			Conditions: squirrel.Eq{model.PluginConfigurationTableName + ".ChannelID": channels.IDs()},
		})
	if appErr != nil {
		return nil, appErr
	}

	// keys are plugin configurations's identifiers
	var configsMap = lo.SliceToMap(pluginConfigsOfChannels, func(p *model.PluginConfiguration) (string, *model.PluginConfiguration) { return p.Identifier, p })

	for _, pluginInitObj := range pluginInitObjects {
		var (
			pluginConfig = pluginInitObj.Manifest.DefaultConfiguration
			active       = pluginInitObj.Manifest.DefaultActive
			channelID    string
		)
		if existingConfig, ok := configsMap[pluginInitObj.Manifest.PluginID]; ok {
			pluginConfig = existingConfig.Configuration
			active = existingConfig.Active
			channelID = existingConfig.ChannelID
		}

		plugin := pluginInitObj.NewPluginFunc(&PluginConfig{
			Manager:       manager,
			Configuration: pluginConfig,
			Active:        active,
			ChannelID:     channelID,
			Manifest:      pluginInitObj.Manifest,
		})

		manager.allPlugins = append(manager.allPlugins, plugin)
	}

	return manager, nil
}

func (m *PluginManager) getPlugins(channelID string, active bool) []interfaces.BasePluginInterface {
	res := []interfaces.BasePluginInterface{}

	for _, plg := range m.allPlugins {
		if plg != nil && active == plg.IsActive() && (channelID == "" || channelID == plg.ChannelId()) {
			res = append(res, plg)
		}
	}

	return res
}

func (m *PluginManager) ChangeUserAddress(address model.Address, addressType model_helper.AddressTypeEnum, user *model.User) (*model.Address, *model_helper.AppError) {
	var (
		appErr        *model_helper.AppError
		previousValue model.Address = address
		anAddress     *model.Address
	)

	for _, plg := range m.getPlugins("", true) {
		anAddress, appErr = plg.ChangeUserAddress(address, addressType, user, previousValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				anAddress = &previousValue
				continue
			}
			return nil, appErr
		}
		previousValue = *anAddress
	}

	return anAddress, nil
}

func (m *PluginManager) CalculateCheckoutTotal(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	subTotal, appErr := m.CalculateCheckoutSubTotal(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	shippingPrice, appErr := m.CalculateCheckoutShipping(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	defaultValue, appErr := m.Srv.CheckoutService().BaseCheckoutTotal(subTotal, shippingPrice, checkoutInfo.Checkout.Discount, checkoutInfo.Checkout.Currency)
	if appErr != nil {
		return nil, appErr
	}

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, appErr = plg.CalculateCheckoutTotal(checkoutInfo, lines, address, discounts, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateCheckoutTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutSubTotal(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	lineTotalSum, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency)
	var err error

	for _, line := range lines.FilterNils() {
		taxedMoney, appErr := m.CalculateCheckoutLineTotal(checkoutInfo, lines, *line, address, discounts)
		if appErr != nil {
			return nil, appErr
		}

		lineTotalSum, err = lineTotalSum.Add(taxedMoney)
		if err != nil {
			return nil, model_helper.NewAppError("CalculateCheckoutSubTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
	}

	quantizedTaxedMoney, _ := lineTotalSum.Quantize(goprices.Up, -1)
	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutShipping(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseCheckoutShippingPrice(&checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, appErr = plg.CalculateCheckoutShipping(checkoutInfo, lines, address, discounts, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateCheckoutShipping", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutLineTotal(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, checkoutLineInfo model_helper.CheckoutLineInfo, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseCheckoutLineTotal(&checkoutLineInfo, &checkoutInfo.Channel, discounts)
	if appErr != nil {
		return nil, appErr
	}

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, appErr = plg.CalculateCheckoutLineTotal(checkoutInfo, lines, checkoutLineInfo, address, discounts, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateCheckoutLineTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateOrderShipping(orDer model.Order) (*goprices.TaxedMoney, *model_helper.AppError) {
	if orDer.ShippingMethodID == nil {
		zero, _ := util.ZeroTaxedMoney(orDer.Currency)
		return zero, nil
	}

	shippingMethodChannelListings, appErr := m.Srv.ShippingService().ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
		Conditions: squirrel.Eq{
			model.ShippingMethodChannelListingTableName + ".ShippingMethodID": orDer.ShippingMethodID,
			model.ShippingMethodChannelListingTableName + ".ChannelID":        orDer.ChannelID,
		},
	})
	if appErr != nil {
		return nil, appErr
	}
	firstItem := shippingMethodChannelListings[0]
	firstItem.PopulateNonDbFields() // this call is mandatory.

	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      firstItem.Price,
		Gross:    firstItem.Price,
		Currency: firstItem.Currency,
	}).
		Quantize(goprices.Up, -1)

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		taxedMoney, appErr = plg.CalculateOrderShipping(&orDer, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateOrderShipping", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) GetCheckoutShippingTaxRate(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo, shippingPrice goprices.TaxedMoney) (*decimal.Decimal, *model_helper.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseTaxRate(&shippingPrice)
	if appErr != nil {
		return nil, appErr
	}

	var deci *decimal.Decimal

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		deci, appErr = plg.GetCheckoutShippingTaxRate(checkoutInfo, lines, address, discounts, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				deci = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = deci
	}

	return model_helper.GetPointerOfValue(deci.Round(4)), nil
}

func (m *PluginManager) GetOrderShippingTaxRate(orDer model.Order, shippingPrice goprices.TaxedMoney) (*decimal.Decimal, *model_helper.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseTaxRate(&shippingPrice)
	if appErr != nil {
		return nil, appErr
	}

	var deci *decimal.Decimal

	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		deci, appErr = plg.GetOrderShippingTaxRate(orDer, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				deci = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = deci
	}

	return model_helper.GetPointerOfValue(deci.Round(4)), nil
}

func (m *PluginManager) CalculateOrderlineTotal(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product) (*goprices.TaxedMoney, *model_helper.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseOrderLineTotal(&orderLine)
	if appErr != nil {
		return nil, appErr
	}

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		taxedMoney, appErr = plg.CalculateOrderLineTotal(&orDer, &orderLine, variant, product, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateOrderlineTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutLineUnitPrice(totalLinePrice goprices.TaxedMoney, quantity int, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, checkoutLineInfo model_helper.CheckoutLineInfo, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	defaultValue := m.Srv.Checkout.BaseCheckoutLineUnitPrice(&totalLinePrice, quantity)

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(checkoutInfo.Channel.ID, true) {
		taxedMoney, appErr := plg.CalculateCheckoutLineUnitPrice(checkoutInfo, lines, checkoutLineInfo, address, discounts, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateCheckoutLineUnitPrice", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateOrderLineUnit(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product) (*goprices.TaxedMoney, *model_helper.AppError) {
	orderLine.PopulateNonDbFields() // this is needed
	defaultValue, err := orderLine.UnitPrice.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateOrderLineUnit", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	var (
		taxedMoney *goprices.TaxedMoney
		appErr     *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		taxedMoney, appErr = plg.CalculateOrderLineUnit(orDer, orderLine, variant, product, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateOrderLineUnit", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) GetCheckoutLineTaxRate(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, checkoutLineInfo model_helper.CheckoutLineInfo, address *model.Address, discounts []*model_helper.DiscountInfo, unitPrice goprices.TaxedMoney) (*decimal.Decimal, *model_helper.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseTaxRate(&unitPrice)
	if appErr != nil {
		return nil, appErr
	}

	var deci *decimal.Decimal

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		deci, appErr = plg.GetCheckoutLineTaxRate(&checkoutInfo, lines, checkoutLineInfo, address, discounts, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				deci = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = deci
	}

	return model_helper.GetPointerOfValue(deci.RoundUp(4)), nil
}

func (m *PluginManager) GetOrderLineTaxRate(orDer model.Order, product model.Product, variant model.ProductVariant, address *model.Address, unitPrice goprices.TaxedMoney) (*decimal.Decimal, *model_helper.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseTaxRate(&unitPrice)
	if appErr != nil {
		return nil, appErr
	}

	var deci *decimal.Decimal

	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		deci, appErr = plg.GetOrderLineTaxRate(orDer, product, variant, address, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				deci = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = deci
	}

	return model_helper.GetPointerOfValue(deci.RoundUp(4)), nil
}

func (m *PluginManager) GetTaxRateTypeChoices() ([]*model_helper.TaxType, *model_helper.AppError) {
	defaultValue := []*model_helper.TaxType{}

	var (
		taxTypes []*model_helper.TaxType
		appErr   *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		taxTypes, appErr = plg.GetTaxRateTypeChoices(defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxTypes = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxTypes
	}

	return taxTypes, nil
}

func (m *PluginManager) ShowTaxesOnStoreFront() (bool, *model_helper.AppError) {
	defaultValue := false

	var (
		showTax bool
		appErr  *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		showTax, appErr = plg.ShowTaxesOnStorefront(defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				showTax = defaultValue
				continue
			}
			return false, appErr
		}
		defaultValue = showTax
	}

	return showTax, nil
}

func (m *PluginManager) ApplyTaxesToProduct(product model.Product, price goprices.Money, country model.CountryCode, channelID string) (*goprices.TaxedMoney, *model_helper.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(goprices.Up, -1)

	var (
		taxedMoney *goprices.TaxedMoney
		appErr     *model_helper.AppError
	)
	for _, plg := range m.getPlugins(channelID, true) {
		taxedMoney, appErr = plg.ApplyTaxesToProduct(product, price, country, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("ApplyTaxesToProduct", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	return quantizedTaxedMoney, nil
}

func (m *PluginManager) ApplyTaxesToShipping(price goprices.Money, shippingAddress model.Address, channelID string) (*goprices.TaxedMoney, *model_helper.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(goprices.Up, -1)

	var (
		taxedMoney *goprices.TaxedMoney
		appErr     *model_helper.AppError
	)
	for _, plg := range m.getPlugins(channelID, true) {
		taxedMoney, appErr = plg.ApplyTaxesToShipping(price, shippingAddress, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("ApplyTaxesToShipping", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) PreprocessOrderCreation(checkoutInfo model_helper.CheckoutInfo, discounts []*model_helper.DiscountInfo, lines model_helper.CheckoutLineInfos) (any, *model_helper.AppError) {
	var defaultValue any = nil

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		value, appErr = plg.PreprocessOrderCreation(checkoutInfo, discounts, lines, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) CustomerCreated(customer model.User) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.CustomerCreated(customer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, appErr
}

func (m *PluginManager) CustomerUpdated(customer model.User) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.CustomerUpdated(customer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) ProductCreated(product model.Product) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.ProductCreated(product, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) ProductUpdated(product model.Product) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.ProductUpdated(product, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) ProductDeleted(product model.Product, variants []int) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.ProductDeleted(product, variants, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, appErr
}

func (m *PluginManager) ProductVariantCreated(variant model.ProductVariant) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.ProductVariantCreated(variant, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) ProductVariantUpdated(variant model.ProductVariant) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.ProductVariantUpdated(variant, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) ProductVariantDeleted(variant model.ProductVariant) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.ProductVariantDeleted(variant, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) ProductVariantOutOfStock(stock model.Stock) *model_helper.AppError {
	var defaultValue any

	var appErr *model_helper.AppError

	for _, plg := range m.getPlugins("", true) {
		appErr = plg.ProductVariantOutOfStock(stock, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				continue
			}
			return appErr
		}
	}

	return nil
}

func (m *PluginManager) ProductVariantBackInStock(stock model.Stock) *model_helper.AppError {
	var defaultValue any

	var appErr *model_helper.AppError

	for _, plg := range m.getPlugins("", true) {
		appErr = plg.ProductVariantBackInStock(stock, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				continue
			}
			return appErr
		}
	}

	return nil
}

func (m *PluginManager) OrderCreated(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.OrderCreated(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) OrderConfirmed(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.OrderConfirmed(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) DraftOrderCreated(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.DraftOrderCreated(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, appErr
}

func (m *PluginManager) DraftOrderDeleted(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.DraftOrderDeleted(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) DraftOrderUpdated(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.DraftOrderUpdated(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) SaleCreated(sale model.Sale, currentCatalogue model_helper.NodeCatalogueInfo) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.SaleCreated(sale, currentCatalogue, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) SaleDeleted(sale model.Sale, previousCatalogue model_helper.NodeCatalogueInfo) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.SaleDeleted(sale, previousCatalogue, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) SaleUpdated(sale model.Sale, previousCatalogue, currentCatalogue model_helper.NodeCatalogueInfo) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.SaleUpdated(sale, previousCatalogue, currentCatalogue, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) InvoiceRequest(orDer model.Order, inVoice model.Invoice, number string) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.InvoiceRequest(orDer, inVoice, number, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) InvoiceDelete(inVoice model.Invoice) (any, *model_helper.AppError) {
	var defaultValue any

	var channelID string
	if inVoice.OrderID != nil {
		orDer, appErr := m.Srv.OrderService().OrderById(*inVoice.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = orDer.ChannelID
	}

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(channelID, true) {
		value, appErr = plg.InvoiceDelete(inVoice, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) InvoiceSent(inVoice model.Invoice, email string) (any, *model_helper.AppError) {
	var defaultValue any

	var channelID string
	if inVoice.OrderID != nil {
		orDer, appErr := m.Srv.OrderService().OrderById(*inVoice.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = orDer.ChannelID
	}

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(channelID, true) {
		value, appErr = plg.InvoiceSent(inVoice, email, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) OrderFullyPaid(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.OrderFullyPaid(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) OrderUpdated(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.OrderUpdated(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) OrderCancelled(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.OrderCancelled(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) OrderFulfilled(orDer model.Order) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.OrderFulfilled(orDer, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) FulfillmentCreated(fulfillment model.Fulfillment) (any, *model_helper.AppError) {
	var defaultValue any

	orDer, appErr := m.Srv.OrderService().OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	var value any

	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.FulfillmentCreated(fulfillment, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) FulfillmentCanceled(fulfillment model.Fulfillment) (any, *model_helper.AppError) {
	var defaultValue any

	orDer, appErr := m.Srv.OrderService().OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	var value any

	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, appErr = plg.FulfillmentCanceled(fulfillment, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) CheckoutCreated(checkOut model.Checkout) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(checkOut.ChannelID, true) {
		value, appErr = plg.CheckoutCreated(checkOut, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) CheckoutUpdated(checkOut model.Checkout) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins(checkOut.ChannelID, true) {
		value, appErr = plg.CheckoutUpdated(checkOut, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) PageCreated(paGe model.Page) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.PageCreated(paGe, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) PageUpdated(paGe model.Page) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.PageUpdated(paGe, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) PageDeleted(paGe model.Page) (any, *model_helper.AppError) {
	var defaultValue any

	var (
		value  any
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.PageDeleted(paGe, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) getPlugin(pluginID string, channelID string) interfaces.BasePluginInterface {
	for _, plg := range m.allPlugins {
		if plg != nil && plg.CheckPluginId(pluginID) && (channelID == "" || plg.ChannelId() == channelID) {
			return plg
		}
	}

	return nil
}

func (m *PluginManager) InitializePayment(gateway string, paymentData model_types.JSONString, channelID string) *model_helper.InitializedPaymentResponse {
	plg := m.getPlugin(gateway, channelID)
	if plg == nil {
		return nil
	}

	value, _ := plg.InitializePayment(paymentData, nil)
	return value
}

func (m *PluginManager) runPaymentMethod(gateway, methodName string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error) {
	plg := m.getPlugin(gateway, channelID)

	if plg != nil {
		var (
			value  *model_helper.GatewayResponse
			appErr *model_helper.AppError
		)

		switch methodName {
		case "authorize_payment":
			value, appErr = plg.AuthorizePayment(paymentInformation, nil)
		case "capture_payment":
			value, appErr = plg.CapturePayment(paymentInformation, nil)
		case "refund_payment":
			value, appErr = plg.RefundPayment(paymentInformation, nil)
		case "void_payment":
			value, appErr = plg.VoidPayment(paymentInformation, nil)
		case "confirm_payment":
			value, appErr = plg.ConfirmPayment(paymentInformation, nil)
		case "process_payment":
			value, appErr = plg.ProcessPayment(paymentInformation, nil)

		default:
			return nil, fmt.Errorf("no method found")
		}

		if appErr != nil {
			return nil, appErr
		}

		return value, nil
	}

	return nil, fmt.Errorf("payment plugin %s for %s payment method is in-accessible", gateway, methodName)
}

func (m *PluginManager) AuthorizePayment(gateway string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "authorize_payment", paymentInformation, channelID)
}

func (m *PluginManager) CapturePayment(gateway string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "capture_payment", paymentInformation, channelID)
}

func (m *PluginManager) RefundPayment(gateway string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "refund_payment", paymentInformation, channelID)
}

func (m *PluginManager) VoidPayment(gateway string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "void_payment", paymentInformation, channelID)
}

func (m *PluginManager) ConfirmPayment(gateway string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "confirm_payment", paymentInformation, channelID)
}

func (m *PluginManager) ProcessPayment(gateway string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "process_payment", paymentInformation, channelID)
}

func (m *PluginManager) TokenIsRequiredAsPaymentInput(gateway, channelID string) (bool, *model_helper.AppError) {
	plg := m.getPlugin(gateway, channelID)
	defaultValue := true

	if plg != nil {
		value, appErr := plg.TokenIsRequiredAsPaymentInput(defaultValue)
		if appErr != nil {
			return false, appErr
		}
		return value, nil
	}

	return defaultValue, nil
}

func (m *PluginManager) GetClientToken(gateway string, tokenConfig model_helper.TokenConfig, channelID string) (string, *model_helper.AppError) {
	plg := m.getPlugin(gateway, channelID)
	if plg != nil {
		value, appErr := plg.GetClientToken(tokenConfig, nil)
		if appErr != nil {
			return "", appErr
		}
		return value, nil
	}

	return "", nil
}

func (m *PluginManager) ListPaymentSources(gateway, customerID, channelID string) ([]*model_helper.CustomerSource, error) {
	plg := m.getPlugin(gateway, channelID)
	if plg != nil {
		defaultValue := []*model_helper.CustomerSource{}
		return plg.ListPaymentSources(customerID, defaultValue)
	}

	return nil, fmt.Errorf("payment plugin %s is inaccessible", gateway)
}

func (m *PluginManager) TranslationCreated(translation any) {
	panic("not implemented")
}

func (m *PluginManager) TranslationUpdated(translation any) {
	panic("not implemented")
}

func (m *PluginManager) ListPaymentGateways(currency string, checkOut *model.Checkout, channelID string, activeOnly bool) []*model_helper.PaymentGateway {
	if checkOut != nil {
		channelID = checkOut.ChannelID
	}
	plugins := m.getPlugins(channelID, activeOnly)

	// if currency is given return only gateways which support given currency
	var gateways []*model_helper.PaymentGateway

	for _, plg := range plugins {
		value, appErr := plg.GetPaymentGateways(currency, checkOut, nil)
		if appErr != nil {
			continue
		}
		gateways = append(gateways, value...)
	}

	return gateways
}

func (m *PluginManager) ListExternalAuthentications(activeOnly bool) ([]model_types.JSONString, *model_helper.AppError) {
	filteredPlugins := m.getPlugins("", activeOnly)

	res := []model_types.JSONString{}

	for _, plg := range filteredPlugins {
		_, appErr := plg.ExternalObtainAccessTokens(nil, nil, model_helper.ExternalAccessTokens{})
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				continue
			}
			return nil, appErr
		}
		manifest := plg.GetManifest()
		res = append(res, model_types.JSONString{
			"id":   manifest.PluginID,
			"name": manifest.PluginName,
		})
	}

	return res, nil
}

// AssignTaxCodeToObjectMeta requires obj must be Product or ProductType
func (m *PluginManager) AssignTaxCodeToObjectMeta(obj any, taxCode string) (*model_helper.TaxType, *model_helper.AppError) {
	// validate obj
	switch obj.(type) {
	case model.Product,
		model.ProductType,
		*model.Product,
		*model.ProductType:
	default:
		return nil, model_helper.NewAppError("AssignTaxCodeToObjectMeta", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	var (
		defaultValue = new(model_helper.TaxType)
		value        *model_helper.TaxType
		appErr       *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.AssignTaxCodeToObjectMeta(obj, taxCode, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

// GetTaxCodeFromObjectMeta
//
// NOTE: obj must be either Product or ProductType
func (m *PluginManager) GetTaxCodeFromObjectMeta(obj any) (*model_helper.TaxType, *model_helper.AppError) {
	// validate obj
	switch obj.(type) {
	case model.Product,
		model.ProductType,
		*model.Product,
		*model.ProductType:
	default:
		return nil, model_helper.NewAppError("GetTaxCodeFromObjectMeta", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	var (
		defaultValue = new(model_helper.TaxType)
		value        *model_helper.TaxType
		appErr       *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.GetTaxCodeFromObjectMeta(obj, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

// GetTaxRatePercentageValue
//
// obj must be either Product or ProductType
func (m *PluginManager) GetTaxRatePercentageValue(obj any, country string) (*decimal.Decimal, *model_helper.AppError) {
	switch obj.(type) {
	case model.Product,
		model.ProductType,
		*model.Product,
		*model.ProductType:
	default:
		return nil, model_helper.NewAppError("GetTaxRatePercentageValue", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	defaultValue := decimal.Zero.Round(0)

	var (
		deci   *decimal.Decimal
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		deci, appErr = plg.GetTaxRatePercentageValue(obj, country, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				deci = &defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = *deci
	}

	return model_helper.GetPointerOfValue(deci.Round(0)), nil
}

func (m *PluginManager) SavePluginConfiguration(pluginID, channelID string, cleanedData model_types.JSONString) (*model.PluginConfiguration, *model_helper.AppError) {
	if !model_helper.IsValidId(channelID) {
		return nil, model_helper.NewAppError("SavePluginConfiguration", model_helper.InvalidArgumentAppErrorID, nil, "", http.StatusBadRequest)
	}

	var pluginList []interfaces.BasePluginInterface
	if channelID != "" {
		pluginList = m.getPlugins(channelID, true)
	} else {
		pluginList = m.allPlugins
	}

	for _, plg := range pluginList {
		manifest := plg.GetManifest()
		if manifest.PluginID == pluginID {

			// try get or create plugin configuration
			pluginConfig, appErr := m.Srv.Plugin.GetPluginConfiguration(&model.PluginConfigurationFilterOptions{
				Conditions: squirrel.Eq{
					model.PluginConfigurationTableName + ".Identifier": pluginID,
					model.PluginConfigurationTableName + ".ChannelID":  channelID,
				},
			})
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError {
					return nil, appErr
				}

				pluginConfig, appErr = m.Srv.PluginService().UpsertPluginConfiguration(&model.PluginConfiguration{
					Identifier:    pluginID,
					ChannelID:     channelID,
					Configuration: plg.GetConfiguration(),
				})
				if appErr != nil {
					return nil, appErr
				}
			}

			pluginConfig, appErr = plg.SavePluginConfiguration(pluginConfig, cleanedData)
			if appErr != nil {
				if appErr.StatusCode == http.StatusNotImplemented {
					m.Srv.Log.Warn("Method not implemented", slog.String("method", appErr.Where), slog.Err(appErr))
				}
				return nil, appErr
			}

			pluginConfig.Name = manifest.PluginName
			pluginConfig.Description = manifest.Description
			plg.SetActive(pluginConfig.Active)
			plg.SetConfiguration(pluginConfig.Configuration)

			return pluginConfig, nil
		}
	}

	return nil, nil
}

func (m *PluginManager) FetchTaxesData() (bool, *model_helper.AppError) {
	defaultValue := false

	var (
		value  bool
		appErr *model_helper.AppError
	)
	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.FetchTaxesData(defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return false, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) WebhookEndpointWithoutChannel(req *http.Request, pluginID string) (*http.Response, *model_helper.AppError) {
	splitPath := strings.SplitN(req.URL.Path, pluginID, 1)

	var path string
	if len(splitPath) == 2 {
		path = splitPath[1]
	}

	defaultValue := http.Response{
		StatusCode: http.StatusNotFound,
	}
	plg := m.getPlugin(pluginID, "")
	if plg == nil {
		return &defaultValue, nil
	}

	value, appErr := plg.Webhook(req, path, defaultValue)
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotImplemented {
			return &defaultValue, nil
		}
		return nil, appErr
	}

	return value, nil
}

func (m *PluginManager) Webhook(req *http.Request, pluginID, channelID string) (*http.Response, *model_helper.AppError) {
	splitPath := strings.SplitN(req.URL.Path, pluginID, 1)

	var path string
	if len(splitPath) == 2 {
		path = splitPath[1]
	}

	defaultValue := &http.Response{
		StatusCode: http.StatusNotFound,
	}

	plg := m.getPlugin(pluginID, channelID)
	if plg == nil {
		return defaultValue, nil
	}

	if !plg.IsActive() {
		return defaultValue, nil
	}

	if manifest := plg.GetManifest(); manifest.ConfigurationPerChannel && channelID == "" {
		return &http.Response{
			Body: io.NopCloser(strings.NewReader("Incorrect endpoint. Use /plugins/channel/<channel_id>/" + manifest.PluginID)),
		}, nil
	}

	res, appErr := plg.Webhook(req, path, *defaultValue)
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotImplemented {
			return defaultValue, nil
		}
		return nil, appErr
	}
	return res, nil
}

func (m *PluginManager) Notify(event string, payload model_types.JSONString, channelID string, pluginID string) (any, *model_helper.AppError) {
	var defaultValue any

	if pluginID != "" {
		plg := m.getPlugin(pluginID, channelID)
		value, appErr := plg.Notify(event, payload, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				return defaultValue, nil
			}
			return nil, appErr
		}
		return value, nil
	}

	for _, plg := range m.getPlugins(channelID, true) {
		value, appErr := plg.Notify(event, payload, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return defaultValue, nil
}

func (m *PluginManager) ExternalObtainAccessTokens(pluginID string, data model_types.JSONString, req *http.Request) (*model_helper.ExternalAccessTokens, *model_helper.AppError) {
	var defaultValue model_helper.ExternalAccessTokens
	plg := m.getPlugin(pluginID, "")

	if plg != nil {
		res, appErr := plg.ExternalObtainAccessTokens(data, req, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				return &defaultValue, nil
			}
			return nil, appErr
		}
		return res, nil
	}

	return &defaultValue, nil
}

func (m *PluginManager) ExternalAuthenticationUrl(pluginID string, data model_types.JSONString, req *http.Request) (model_types.JSONString, *model_helper.AppError) {
	defaultValue := model_types.JSONString{}

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		res, appErr := plg.ExternalAuthenticationUrl(data, req, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				return defaultValue, nil
			}
			return nil, appErr
		}
		return res, nil
	}

	return defaultValue, nil
}

func (m *PluginManager) ExternalRefresh(pluginID string, data model_types.JSONString, req *http.Request) (*model_helper.ExternalAccessTokens, *model_helper.AppError) {
	var defaultValue model_helper.ExternalAccessTokens

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		res, appErr := plg.ExternalRefresh(data, req, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				return &defaultValue, nil
			}
			return nil, appErr
		}
		return res, nil
	}

	return &defaultValue, nil
}

func (m *PluginManager) AuthenticateUser(req *http.Request) (*model.User, *model_helper.AppError) {
	var (
		defaultValue *model.User = nil
		value        *model.User
		appErr       *model_helper.AppError
	)

	for _, plg := range m.getPlugins("", true) {
		value, appErr = plg.AuthenticateUser(req, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				value = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) ExternalLogout(pluginID string, data model_types.JSONString, req *http.Request) (model_types.JSONString, *model_helper.AppError) {
	defaultValue := model_types.JSONString{}

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		appErr := plg.ExternalLogout(data, req, defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				return defaultValue, nil
			}
			return nil, appErr
		}
	}

	return defaultValue, nil
}

func (m *PluginManager) ExternalVerify(pluginID string, data model_types.JSONString, req *http.Request) (*model.User, model_types.JSONString, *model_helper.AppError) {
	var (
		defaultData = model_types.JSONString{}
		defaultUser *model.User
	)

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		user, data, appErr := plg.ExternalVerify(data, req, interfaces.AType{
			User: defaultUser,
			Data: defaultData,
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				return defaultUser, defaultData, nil
			}
			return nil, nil, appErr
		}
		return user, data, nil
	}

	return defaultUser, defaultData, nil
}
