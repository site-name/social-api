package plugin

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

var _ interfaces.PluginManagerInterface = (*PluginManager)(nil)

type PluginManager struct {
	Srv        *app.Server
	allPlugins []interfaces.BasePluginInterface
}

// NewPluginManager returns a new plugin manager
func (s *ServicePlugin) newPluginManager() (interfaces.PluginManagerInterface, *model.AppError) {
	manager := &PluginManager{
		Srv: s.srv,
	}

	channels, appErr := manager.Srv.
		ChannelService().
		ChannelsByOption(&model.ChannelFilterOption{})
	if appErr != nil {
		return nil, appErr
	}

	// finds a list of plugin configs belong found channels
	pluginConfigsOfChannels, appErr := manager.Srv.
		PluginService().
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

func (m *PluginManager) ChangeUserAddress(address model.Address, addressType *model.AddressTypeEnum, user *model.User) (*model.Address, *model.AppError) {
	var (
		appErr        *model.AppError
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

func (m *PluginManager) CalculateCheckoutTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
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
		return nil, model.NewAppError("CalculateCheckoutTotal", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutSubTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	lineTotalSum, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency)
	var err error

	for _, line := range lines.FilterNils() {
		taxedMoney, appErr := m.CalculateCheckoutLineTotal(checkoutInfo, lines, *line, address, discounts)
		if appErr != nil {
			return nil, appErr
		}

		lineTotalSum, err = lineTotalSum.Add(taxedMoney)
		if err != nil {
			return nil, model.NewAppError("CalculateCheckoutSubTotal", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
	}

	quantizedTaxedMoney, _ := lineTotalSum.Quantize(goprices.Up, -1)
	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutShipping(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
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
		return nil, model.NewAppError("CalculateCheckoutShipping", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutLineTotal(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
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
		return nil, model.NewAppError("CalculateCheckoutLineTotal", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateOrderShipping(orDer model.Order) (*goprices.TaxedMoney, *model.AppError) {
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
		return nil, model.NewAppError("CalculateOrderShipping", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) GetCheckoutShippingTaxRate(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, address *model.Address, discounts []*model.DiscountInfo, shippingPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewPrimitive(deci.Round(4)), nil
}

func (m *PluginManager) GetOrderShippingTaxRate(orDer model.Order, shippingPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewPrimitive(deci.Round(4)), nil
}

func (m *PluginManager) CalculateOrderlineTotal(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product) (*goprices.TaxedMoney, *model.AppError) {
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
		return nil, model.NewAppError("CalculateOrderlineTotal", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutLineUnitPrice(totalLinePrice goprices.TaxedMoney, quantity int, checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue := m.Srv.CheckoutService().BaseCheckoutLineUnitPrice(&totalLinePrice, quantity)

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
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
		return nil, model.NewAppError("CalculateCheckoutLineUnitPrice", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateOrderLineUnit(orDer model.Order, orderLine model.OrderLine, variant model.ProductVariant, product model.Product) (*goprices.TaxedMoney, *model.AppError) {
	orderLine.PopulateNonDbFields() // this is needed
	defaultValue, err := orderLine.UnitPrice.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model.NewAppError("CalculateOrderLineUnit", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	var (
		taxedMoney *goprices.TaxedMoney
		appErr     *model.AppError
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
		return nil, model.NewAppError("CalculateOrderLineUnit", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) GetCheckoutLineTaxRate(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, checkoutLineInfo model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo, unitPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewPrimitive(deci.RoundUp(4)), nil
}

func (m *PluginManager) GetOrderLineTaxRate(orDer model.Order, product model.Product, variant model.ProductVariant, address *model.Address, unitPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewPrimitive(deci.RoundUp(4)), nil
}

func (m *PluginManager) GetTaxRateTypeChoices() ([]*model.TaxType, *model.AppError) {
	defaultValue := []*model.TaxType{}

	var (
		taxTypes []*model.TaxType
		appErr   *model.AppError
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

func (m *PluginManager) ShowTaxesOnStoreFront() (bool, *model.AppError) {
	defaultValue := false

	var (
		showTax bool
		appErr  *model.AppError
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

func (m *PluginManager) ApplyTaxesToProduct(product model.Product, price goprices.Money, country model.CountryCode, channelID string) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(goprices.Up, -1)

	var (
		taxedMoney *goprices.TaxedMoney
		appErr     *model.AppError
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
		return nil, model.NewAppError("ApplyTaxesToProduct", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	return quantizedTaxedMoney, nil
}

func (m *PluginManager) ApplyTaxesToShipping(price goprices.Money, shippingAddress model.Address, channelID string) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(goprices.Up, -1)

	var (
		taxedMoney *goprices.TaxedMoney
		appErr     *model.AppError
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
		return nil, model.NewAppError("ApplyTaxesToShipping", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) PreprocessOrderCreation(checkoutInfo model.CheckoutInfo, discounts []*model.DiscountInfo, lines model.CheckoutLineInfos) (interface{}, *model.AppError) {
	var defaultValue interface{} = nil

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) CustomerCreated(customer model.User) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) CustomerUpdated(customer model.User) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) ProductCreated(product model.Product) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) ProductUpdated(product model.Product) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) ProductDeleted(product model.Product, variants []int) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) ProductVariantCreated(variant model.ProductVariant) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) ProductVariantUpdated(variant model.ProductVariant) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) ProductVariantDeleted(variant model.ProductVariant) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) ProductVariantOutOfStock(stock model.Stock) *model.AppError {
	var defaultValue interface{}

	var appErr *model.AppError

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

func (m *PluginManager) ProductVariantBackInStock(stock model.Stock) *model.AppError {
	var defaultValue interface{}

	var appErr *model.AppError

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

func (m *PluginManager) OrderCreated(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) OrderConfirmed(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) DraftOrderCreated(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) DraftOrderDeleted(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) DraftOrderUpdated(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) SaleCreated(sale model.Sale, currentCatalogue model.NodeCatalogueInfo) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) SaleDeleted(sale model.Sale, previousCatalogue model.NodeCatalogueInfo) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) SaleUpdated(sale model.Sale, previousCatalogue, currentCatalogue model.NodeCatalogueInfo) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) InvoiceRequest(orDer model.Order, inVoice model.Invoice, number string) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) InvoiceDelete(inVoice model.Invoice) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var channelID string
	if inVoice.OrderID != nil {
		orDer, appErr := m.Srv.OrderService().OrderById(*inVoice.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = orDer.ChannelID
	}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) InvoiceSent(inVoice model.Invoice, email string) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var channelID string
	if inVoice.OrderID != nil {
		orDer, appErr := m.Srv.OrderService().OrderById(*inVoice.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = orDer.ChannelID
	}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) OrderFullyPaid(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) OrderUpdated(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) OrderCancelled(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) OrderFulfilled(orDer model.Order) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) FulfillmentCreated(fulfillment model.Fulfillment) (interface{}, *model.AppError) {
	var defaultValue interface{}

	orDer, appErr := m.Srv.OrderService().OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	var value interface{}

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

func (m *PluginManager) FulfillmentCanceled(fulfillment model.Fulfillment) (interface{}, *model.AppError) {
	var defaultValue interface{}

	orDer, appErr := m.Srv.OrderService().OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	var value interface{}

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

func (m *PluginManager) CheckoutCreated(checkOut model.Checkout) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) CheckoutUpdated(checkOut model.Checkout) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) PageCreated(paGe model.Page) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) PageUpdated(paGe model.Page) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) PageDeleted(paGe model.Page) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var (
		value  interface{}
		appErr *model.AppError
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

func (m *PluginManager) InitializePayment(gateway string, paymentData model.StringInterface, channelID string) *model.InitializedPaymentResponse {
	plg := m.getPlugin(gateway, channelID)
	if plg == nil {
		return nil
	}

	value, _ := plg.InitializePayment(paymentData, nil)
	return value
}

func (m *PluginManager) runPaymentMethod(gateway, methodName string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error) {
	plg := m.getPlugin(gateway, channelID)

	if plg != nil {
		var (
			value  *model.GatewayResponse
			appErr *model.AppError
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

func (m *PluginManager) AuthorizePayment(gateway string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "authorize_payment", paymentInformation, channelID)
}

func (m *PluginManager) CapturePayment(gateway string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "capture_payment", paymentInformation, channelID)
}

func (m *PluginManager) RefundPayment(gateway string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "refund_payment", paymentInformation, channelID)
}

func (m *PluginManager) VoidPayment(gateway string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "void_payment", paymentInformation, channelID)
}

func (m *PluginManager) ConfirmPayment(gateway string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "confirm_payment", paymentInformation, channelID)
}

func (m *PluginManager) ProcessPayment(gateway string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "process_payment", paymentInformation, channelID)
}

func (m *PluginManager) TokenIsRequiredAsPaymentInput(gateway, channelID string) (bool, *model.AppError) {
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

func (m *PluginManager) GetClientToken(gateway string, tokenConfig model.TokenConfig, channelID string) (string, *model.AppError) {
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

func (m *PluginManager) ListPaymentSources(gateway, customerID, channelID string) ([]*model.CustomerSource, error) {
	plg := m.getPlugin(gateway, channelID)
	if plg != nil {
		defaultValue := []*model.CustomerSource{}
		return plg.ListPaymentSources(customerID, defaultValue)
	}

	return nil, fmt.Errorf("payment plugin %s is inaccessible", gateway)
}

func (m *PluginManager) TranslationCreated(translation interface{}) {
	panic("not implemented")
}

func (m *PluginManager) TranslationUpdated(translation interface{}) {
	panic("not implemented")
}

func (m *PluginManager) ListPaymentGateways(currency string, checkOut *model.Checkout, channelID string, activeOnly bool) []*model.PaymentGateway {
	if checkOut != nil {
		channelID = checkOut.ChannelID
	}
	plugins := m.getPlugins(channelID, activeOnly)

	// if currency is given return only gateways which support given currency
	var gateways []*model.PaymentGateway

	for _, plg := range plugins {
		value, appErr := plg.GetPaymentGateways(currency, checkOut, nil)
		if appErr != nil {
			continue
		}
		gateways = append(gateways, value...)
	}

	return gateways
}

func (m *PluginManager) ListExternalAuthentications(activeOnly bool) ([]model.StringInterface, *model.AppError) {
	filteredPlugins := m.getPlugins("", activeOnly)

	res := []model.StringInterface{}

	for _, plg := range filteredPlugins {
		_, appErr := plg.ExternalObtainAccessTokens(nil, nil, model.ExternalAccessTokens{})
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				continue
			}
			return nil, appErr
		}
		manifest := plg.GetManifest()
		res = append(res, model.StringInterface{
			"id":   manifest.PluginID,
			"name": manifest.PluginName,
		})
	}

	return res, nil
}

// AssignTaxCodeToObjectMeta requires obj must be Product or ProductType
func (m *PluginManager) AssignTaxCodeToObjectMeta(obj interface{}, taxCode string) (*model.TaxType, *model.AppError) {

	// validate obj
	switch obj.(type) {
	case model.Product,
		model.ProductType,
		*model.Product,
		*model.ProductType:
	default:
		return nil, model.NewAppError("AssignTaxCodeToObjectMeta", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	var (
		defaultValue = new(model.TaxType)
		value        *model.TaxType
		appErr       *model.AppError
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
func (m *PluginManager) GetTaxCodeFromObjectMeta(obj interface{}) (*model.TaxType, *model.AppError) {
	// validate obj
	switch obj.(type) {
	case model.Product,
		model.ProductType,
		*model.Product,
		*model.ProductType:
	default:
		return nil, model.NewAppError("GetTaxCodeFromObjectMeta", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	var (
		defaultValue = new(model.TaxType)
		value        *model.TaxType
		appErr       *model.AppError
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
func (m *PluginManager) GetTaxRatePercentageValue(obj interface{}, country string) (*decimal.Decimal, *model.AppError) {
	switch obj.(type) {
	case model.Product,
		model.ProductType,
		*model.Product,
		*model.ProductType:
	default:
		return nil, model.NewAppError("GetTaxRatePercentageValue", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	defaultValue := decimal.Zero.Round(0)

	var (
		deci   *decimal.Decimal
		appErr *model.AppError
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

	return model.NewPrimitive(deci.Round(0)), nil
}

func (m *PluginManager) SavePluginConfiguration(pluginID, channelID string, cleanedData model.StringInterface) (*model.PluginConfiguration, *model.AppError) {
	if !model.IsValidId(channelID) {
		return nil, model.NewAppError("SavePluginConfiguration", model.InvalidArgumentAppErrorID, nil, "", http.StatusBadRequest)
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
			pluginConfig, appErr := m.Srv.PluginService().GetPluginConfiguration(&model.PluginConfigurationFilterOptions{
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

func (m *PluginManager) FetchTaxesData() (bool, *model.AppError) {
	defaultValue := false

	var (
		value  bool
		appErr *model.AppError
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

func (m *PluginManager) WebhookEndpointWithoutChannel(req *http.Request, pluginID string) (*http.Response, *model.AppError) {
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

func (m *PluginManager) Webhook(req *http.Request, pluginID, channelID string) (*http.Response, *model.AppError) {
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

func (m *PluginManager) Notify(event string, payload model.StringInterface, channelID string, pluginID string) (interface{}, *model.AppError) {
	var defaultValue interface{}

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

func (m *PluginManager) ExternalObtainAccessTokens(pluginID string, data model.StringInterface, req *http.Request) (*model.ExternalAccessTokens, *model.AppError) {
	var defaultValue model.ExternalAccessTokens
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

func (m *PluginManager) ExternalAuthenticationUrl(pluginID string, data model.StringInterface, req *http.Request) (model.StringInterface, *model.AppError) {
	defaultValue := model.StringInterface{}

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

func (m *PluginManager) ExternalRefresh(pluginID string, data model.StringInterface, req *http.Request) (*model.ExternalAccessTokens, *model.AppError) {
	var defaultValue model.ExternalAccessTokens

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

func (m *PluginManager) AuthenticateUser(req *http.Request) (*model.User, *model.AppError) {
	var (
		defaultValue *model.User = nil
		value        *model.User
		appErr       *model.AppError
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

func (m *PluginManager) ExternalLogout(pluginID string, data model.StringInterface, req *http.Request) (model.StringInterface, *model.AppError) {
	defaultValue := model.StringInterface{}

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

func (m *PluginManager) ExternalVerify(pluginID string, data model.StringInterface, req *http.Request) (*model.User, model.StringInterface, *model.AppError) {
	var (
		defaultData = model.StringInterface{}
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
