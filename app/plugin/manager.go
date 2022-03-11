package plugin

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

var _ interfaces.PluginManagerInterface = (*PluginManager)(nil)

type PluginManager struct {
	Srv        *app.Server
	AllPlugins []interfaces.BasePluginInterface
	ShopID     string
}

// NewPluginManager returns a new plugin manager
func (s *ServicePlugin) NewPluginManager(shopID string) (interfaces.PluginManagerInterface, *model.AppError) {
	m := &PluginManager{
		Srv:    s.srv,
		ShopID: shopID,
	}

	// find all channels belong to given shop
	channels, appErr := m.Srv.ChannelService().ChannelsByOption(&channel.ChannelFilterOption{
		ShopID: squirrel.Eq{store.ChannelTableName + ".ShopID": shopID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// finds a list of plugin configs belong found channels
	pluginConfigsOfChannels, appErr := m.Srv.PluginService().FilterPluginConfigurations(&plugins.PluginConfigurationFilterOptions{
		ChannelID: squirrel.Eq{store.PluginConfigurationTableName + ".ChannelID": channels.IDs()},
	})
	if appErr != nil {
		return nil, appErr
	}

	// keys are plugin configurations's identifiers
	var configsMap = map[string]*plugins.PluginConfiguration{}
	for _, config := range pluginConfigsOfChannels {
		configsMap[config.Identifier] = config
	}

	for _, pluginInitObj := range pluginInitObjects {

		var (
			pluginConfig []model.StringInterface = pluginInitObj.Manifest.DefaultConfiguration
			active       bool                    = pluginInitObj.Manifest.DefaultActive
			channelID    string
		)
		if existingConfig, ok := configsMap[pluginInitObj.Manifest.PluginID]; ok {
			pluginConfig = existingConfig.Configuration
			active = existingConfig.Active
			channelID = existingConfig.ChannelID
		}

		plugin := pluginInitObj.NewPluginFunc(&NewPluginConfig{
			Manager:       m,
			Configuration: pluginConfig,
			Active:        active,
			ChannelID:     channelID,
		})

		m.AllPlugins = append(m.AllPlugins, plugin)
	}

	return m, nil
}

func (m *PluginManager) GetShopID() string {
	return m.ShopID
}

func (m *PluginManager) getPlugins(channelID string, active bool) []interfaces.BasePluginInterface {
	res := []interfaces.BasePluginInterface{}

	for _, plg := range m.AllPlugins {
		if plg != nil && active == plg.IsActive() && (channelID == "" || channelID == plg.ChannelId()) {
			res = append(res, plg)
		}
	}

	return res
}

func (m *PluginManager) ChangeUserAddress(address account.Address, addressType string, user *account.User) (*account.Address, *model.AppError) {
	var (
		appErr        *model.AppError
		previousValue account.Address = address
		address_      *account.Address
	)

	for _, plg := range m.getPlugins("", true) {
		address_, appErr = plg.ChangeUserAddress(address, addressType, user, previousValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				address_ = &previousValue
				continue
			}
			return nil, appErr
		}
		previousValue = *address_
	}

	return address_, nil
}

func (m *PluginManager) CalculateCheckoutTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateCheckoutTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutSubTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	lineTotalSum, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency)
	var err error

	for _, line := range lines.FilterNils() {
		taxedMoney, appErr := m.CalculateCheckoutLineTotal(checkoutInfo, lines, *line, address, discounts)
		if appErr != nil {
			return nil, appErr
		}

		lineTotalSum, err = lineTotalSum.Add(taxedMoney)
		if err != nil {
			return nil, model.NewAppError("CalculateCheckoutSubTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
	}

	quantizedTaxedMoney, _ := lineTotalSum.Quantize(nil, goprices.Up)
	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutShipping(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateCheckoutShipping", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutLineTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateCheckoutLineTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateOrderShipping(orDer order.Order) (*goprices.TaxedMoney, *model.AppError) {
	if orDer.ShippingMethodID == nil {
		zero, _ := util.ZeroTaxedMoney(orDer.Currency)
		return zero, nil
	}

	shippingMethodChannelListings, appErr := m.Srv.ShippingService().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: squirrel.Eq{m.Srv.Store.ShippingMethodChannelListing().TableName("ShippingMethodID"): orDer.ShippingMethodID},
		ChannelID:        squirrel.Eq{m.Srv.Store.ShippingMethodChannelListing().TableName("ChannelID"): orDer.ChannelID},
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
		Quantize(nil, goprices.Up)

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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateOrderShipping", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) GetCheckoutShippingTaxRate(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, shippingPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewDecimal(deci.Round(4)), nil
}

func (m *PluginManager) GetOrderShippingTaxRate(orDer order.Order, shippingPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewDecimal(deci.Round(4)), nil
}

func (m *PluginManager) CalculateOrderlineTotal(orDer order.Order, orderLine order.OrderLine, variant product_and_discount.ProductVariant, product product_and_discount.Product) (*goprices.TaxedMoney, *model.AppError) {
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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateOrderlineTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateCheckoutLineUnitPrice(totalLinePrice goprices.TaxedMoney, quantity int, checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, appErr := m.Srv.CheckoutService().BaseCheckoutLineUnitPrice(&totalLinePrice, quantity)
	if appErr != nil {
		return nil, appErr
	}

	var taxedMoney *goprices.TaxedMoney

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, appErr = plg.CalculateCheckoutLineUnitPrice(checkoutInfo, lines, checkoutLineInfo, address, discounts, *defaultValue)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotImplemented {
				taxedMoney = defaultValue
				continue
			}
			return nil, appErr
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateCheckoutLineUnitPrice", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) CalculateOrderLineUnit(orDer order.Order, orderLine order.OrderLine, variant product_and_discount.ProductVariant, product product_and_discount.Product) (*goprices.TaxedMoney, *model.AppError) {
	orderLine.PopulateNonDbFields() // this is needed
	defaultValue, err := orderLine.UnitPrice.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateOrderLineUnit", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CalculateOrderLineUnit", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) GetCheckoutLineTaxRate(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, checkoutLineInfo checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo, unitPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewDecimal(deci.RoundUp(4)), nil
}

func (m *PluginManager) GetOrderLineTaxRate(orDer order.Order, product product_and_discount.Product, variant product_and_discount.ProductVariant, address *account.Address, unitPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
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

	return model.NewDecimal(deci.RoundUp(4)), nil
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

func (m *PluginManager) ApplyTaxesToProduct(product product_and_discount.Product, price goprices.Money, country string, channelID string) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(nil, goprices.Up)

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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("ApplyTaxesToProduct", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	return quantizedTaxedMoney, nil
}

func (m *PluginManager) ApplyTaxesToShipping(price goprices.Money, shippingAddress account.Address, channelID string) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(nil, goprices.Up)

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

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("ApplyTaxesToShipping", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) PreprocessOrderCreation(checkoutInfo checkout.CheckoutInfo, discounts []*product_and_discount.DiscountInfo, lines checkout.CheckoutLineInfos) (interface{}, *model.AppError) {
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

func (m *PluginManager) CustomerCreated(customer account.User) (interface{}, *model.AppError) {
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

func (m *PluginManager) CustomerUpdated(customer account.User) (interface{}, *model.AppError) {
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

func (m *PluginManager) ProductCreated(product product_and_discount.Product) (interface{}, *model.AppError) {
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

func (m *PluginManager) ProductUpdated(product product_and_discount.Product) (interface{}, *model.AppError) {
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

func (m *PluginManager) ProductDeleted(product product_and_discount.Product, variants []int) (interface{}, *model.AppError) {
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

func (m *PluginManager) ProductVariantCreated(variant product_and_discount.ProductVariant) (interface{}, *model.AppError) {
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

func (m *PluginManager) ProductVariantUpdated(variant product_and_discount.ProductVariant) (interface{}, *model.AppError) {
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

func (m *PluginManager) ProductVariantDeleted(variant product_and_discount.ProductVariant) (interface{}, *model.AppError) {
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

func (m *PluginManager) ProductVariantOutOfStock(stock warehouse.Stock) *model.AppError {
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

func (m *PluginManager) ProductVariantBackInStock(stock warehouse.Stock) *model.AppError {
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

func (m *PluginManager) OrderCreated(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) OrderConfirmed(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) DraftOrderCreated(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) DraftOrderDeleted(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) DraftOrderUpdated(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) SaleCreated(sale product_and_discount.Sale, currentCatalogue product_and_discount.NodeCatalogueInfo) (interface{}, *model.AppError) {
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

func (m *PluginManager) SaleDeleted(sale product_and_discount.Sale, previousCatalogue product_and_discount.NodeCatalogueInfo) (interface{}, *model.AppError) {
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

func (m *PluginManager) SaleUpdated(sale product_and_discount.Sale, previousCatalogue, currentCatalogue product_and_discount.NodeCatalogueInfo) (interface{}, *model.AppError) {
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

func (m *PluginManager) InvoiceRequest(orDer order.Order, inVoice invoice.Invoice, number string) (interface{}, *model.AppError) {
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

func (m *PluginManager) InvoiceDelete(inVoice invoice.Invoice) (interface{}, *model.AppError) {
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

func (m *PluginManager) InvoiceSent(inVoice invoice.Invoice, email string) (interface{}, *model.AppError) {
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

func (m *PluginManager) OrderFullyPaid(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) OrderUpdated(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) OrderCancelled(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) OrderFulfilled(orDer order.Order) (interface{}, *model.AppError) {
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

func (m *PluginManager) FulfillmentCreated(fulfillment order.Fulfillment) (interface{}, *model.AppError) {
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

func (m *PluginManager) FulfillmentCanceled(fulfillment order.Fulfillment) (interface{}, *model.AppError) {
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

func (m *PluginManager) CheckoutCreated(checkOut checkout.Checkout) (interface{}, *model.AppError) {
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

func (m *PluginManager) CheckoutUpdated(checkOut checkout.Checkout) (interface{}, *model.AppError) {
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

func (m *PluginManager) PageCreated(paGe page.Page) (interface{}, *model.AppError) {
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

func (m *PluginManager) PageUpdated(paGe page.Page) (interface{}, *model.AppError) {
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

func (m *PluginManager) PageDeleted(paGe page.Page) (interface{}, *model.AppError) {
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
	for _, plg := range m.AllPlugins {
		if plg != nil && plg.CheckPluginId(pluginID) && (channelID == "" || plg.ChannelId() == channelID) {
			return plg
		}
	}

	return nil
}

func (m *PluginManager) InitializePayment(gateway string, paymentData model.StringInterface, channelID string) *payment.InitializedPaymentResponse {
	plg := m.getPlugin(gateway, channelID)
	if plg == nil {
		return nil
	}

	value, _ := plg.InitializePayment(paymentData, nil)
	return value
}

func (m *PluginManager) runPaymentMethod(gateway, methodName string, paymentInformation payment.PaymentData, channelID string) (*payment.GatewayResponse, error) {
	plg := m.getPlugin(gateway, channelID)

	if plg != nil {
		var (
			value  *payment.GatewayResponse
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

func (m *PluginManager) AuthorizePayment(gateway string, paymentInformation payment.PaymentData, channelID string) (*payment.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "authorize_payment", paymentInformation, channelID)
}

func (m *PluginManager) CapturePayment(gateway string, paymentInformation payment.PaymentData, channelID string) (*payment.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "capture_payment", paymentInformation, channelID)
}

func (m *PluginManager) RefundPayment(gateway string, paymentInformation payment.PaymentData, channelID string) (*payment.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "refund_payment", paymentInformation, channelID)
}

func (m *PluginManager) VoidPayment(gateway string, paymentInformation payment.PaymentData, channelID string) (*payment.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "void_payment", paymentInformation, channelID)
}

func (m *PluginManager) ConfirmPayment(gateway string, paymentInformation payment.PaymentData, channelID string) (*payment.GatewayResponse, error) {
	return m.runPaymentMethod(gateway, "confirm_payment", paymentInformation, channelID)
}

func (m *PluginManager) ProcessPayment(gateway string, paymentInformation payment.PaymentData, channelID string) (*payment.GatewayResponse, error) {
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

func (m *PluginManager) GetClientToken(gateway string, tokenConfig payment.TokenConfig, channelID string) (string, *model.AppError) {
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

func (m *PluginManager) ListPaymentSources(gateway, customerID, channelID string) ([]*payment.CustomerSource, error) {
	plg := m.getPlugin(gateway, channelID)
	if plg != nil {
		defaultValue := []*payment.CustomerSource{}
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

func (m *PluginManager) ListPaymentGateways(currency string, checkOut *checkout.Checkout, channelID string, activeOnly bool) []*payment.PaymentGateway {
	if checkOut != nil {
		channelID = checkOut.ChannelID
	}
	plugins := m.getPlugins(channelID, activeOnly)

	// if currency is given return only gateways which support given currency
	var gateways []*payment.PaymentGateway

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
		_, appErr := plg.ExternalObtainAccessTokens(nil, nil, plugins.ExternalAccessTokens{})
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
	case product_and_discount.Product,
		product_and_discount.ProductType,
		*product_and_discount.Product,
		*product_and_discount.ProductType:
	default:
		return nil, model.NewAppError("AssignTaxCodeToObjectMeta", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
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
	case product_and_discount.Product,
		product_and_discount.ProductType,
		*product_and_discount.Product,
		*product_and_discount.ProductType:
	default:
		return nil, model.NewAppError("GetTaxCodeFromObjectMeta", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
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
	case product_and_discount.Product,
		product_and_discount.ProductType,
		*product_and_discount.Product,
		*product_and_discount.ProductType:
	default:
		return nil, model.NewAppError("GetTaxRatePercentageValue", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
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

	return model.NewDecimal(deci.Round(0)), nil
}

func (m *PluginManager) SavePluginConfiguration(pluginID, channelID string, cleanedData model.StringInterface) (*plugins.PluginConfiguration, *model.AppError) {
	if !model.IsValidId(channelID) {
		return nil, model.NewAppError("SavePluginConfiguration", app.InvalidArgumentAppErrorID, nil, "", http.StatusBadRequest)
	}

	var pluginList []interfaces.BasePluginInterface
	if channelID != "" {
		pluginList = m.getPlugins(channelID, true)
	} else {
		pluginList = m.AllPlugins
	}

	for _, plg := range pluginList {
		manifest := plg.GetManifest()
		if manifest.PluginID == pluginID {

			// try get or create plugin configuration
			pluginConfig, appErr := m.Srv.PluginService().GetPluginConfiguration(&plugins.PluginConfigurationFilterOptions{
				Identifier: squirrel.Eq{m.Srv.Store.PluginConfiguration().TableName("Identifier"): pluginID},
				ChannelID:  squirrel.Eq{m.Srv.Store.PluginConfiguration().TableName("ChannelID"): channelID},
			})
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError {
					return nil, appErr
				}

				pluginConfig, appErr = m.Srv.PluginService().UpsertPluginConfiguration(&plugins.PluginConfiguration{
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

func (m *PluginManager) ExternalObtainAccessTokens(pluginID string, data model.StringInterface, req *http.Request) (*plugins.ExternalAccessTokens, *model.AppError) {
	var defaultValue plugins.ExternalAccessTokens
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

func (m *PluginManager) ExternalRefresh(pluginID string, data model.StringInterface, req *http.Request) (*plugins.ExternalAccessTokens, *model.AppError) {
	var defaultValue plugins.ExternalAccessTokens

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

func (m *PluginManager) AuthenticateUser(req *http.Request) (*account.User, *model.AppError) {
	var (
		defaultValue *account.User = nil
		value        *account.User
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

func (m *PluginManager) ExternalVerify(pluginID string, data model.StringInterface, req *http.Request) (*account.User, model.StringInterface, *model.AppError) {
	var (
		defaultData = model.StringInterface{}
		defaultUser *account.User
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
