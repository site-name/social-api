package plugins

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
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
)

type PluginManager struct {
	srv        *app.Server
	AllPlugins []BasePluginInterface
}

// NewPluginManager returns a new plugin manager
func NewPluginManager(srv *app.Server, shopID string) (*PluginManager, *model.AppError) {
	m := &PluginManager{
		srv: srv,
	}

	// find all channels belong to given shop
	channels, appErr := m.srv.ChannelService().ChannelsByOption(&channel.ChannelFilterOption{
		ShopID: squirrel.Eq{m.srv.Store.Channel().TableName("ShopID"): shopID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// finds a list of plugin configs belong found channels
	pluginConfigsOfChannels, appErr := m.srv.PluginService().FilterPluginConfigurations(&plugins.PluginConfigurationFilterOptions{
		ChannelID: squirrel.Eq{m.srv.Store.PluginConfiguration().TableName("ChannelID"): channels.IDs()},
		// PrefetchRelatedChannel: true, //
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

func (m *PluginManager) getPlugins(channelID string, active bool) []BasePluginInterface {
	res := []BasePluginInterface{}

	for _, plg := range m.AllPlugins {
		if ((active && plg.IsActive()) || (!active && !plg.IsActive())) &&
			(channelID == "" || channelID == plg.ChannelId()) {
			res = append(res, plg)
		}
	}

	return res
}

func (m *PluginManager) ChangeUserAddress(address account.Address, addressType string, user *account.User) *account.Address {
	var (
		notImplt      *PluginMethodNotImplemented
		previousValue account.Address = address
		address_      *account.Address
	)

	for _, plg := range m.getPlugins("", true) {
		address_, notImplt = plg.ChangeUserAddress(address, addressType, user, previousValue)
		if notImplt != nil {
			address_ = &previousValue
			continue
		}
		previousValue = *address_
	}

	return address_
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

	defaultValue, appErr := m.srv.CheckoutService().BaseCheckoutTotal(subTotal, shippingPrice, checkoutInfo.Checkout.Discount, checkoutInfo.Checkout.Currency)
	if appErr != nil {
		return nil, appErr
	}

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, notImplt = plg.CalculateCheckoutTotal(checkoutInfo, lines, address, discounts, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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
	defaultValue, appErr := m.srv.CheckoutService().BaseCheckoutShippingPrice(&checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, notImplt = plg.CalculateCheckoutShipping(checkoutInfo, lines, address, discounts, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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
	defaultValue, appErr := m.srv.CheckoutService().BaseCheckoutLineTotal(&checkoutLineInfo, &checkoutInfo.Channel, discounts)
	if appErr != nil {
		return nil, appErr
	}

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)

	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, notImplt = plg.CalculateCheckoutLineTotal(checkoutInfo, lines, checkoutLineInfo, address, discounts, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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

	shippingMethodChannelListings, appErr := m.srv.ShippingService().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: squirrel.Eq{m.srv.Store.ShippingMethodChannelListing().TableName("ShippingMethodID"): orDer.ShippingMethodID},
		ChannelID:        squirrel.Eq{m.srv.Store.ShippingMethodChannelListing().TableName("ChannelID"): orDer.ChannelID},
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

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		taxedMoney, notImplt = plg.CalculateOrderShipping(&orDer, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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
	defaultValue, appErr := m.srv.CheckoutService().BaseTaxRate(&shippingPrice)
	if appErr != nil {
		return nil, appErr
	}

	var (
		deci     *decimal.Decimal
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		deci, notImplt = plg.GetCheckoutShippingTaxRate(checkoutInfo, lines, address, discounts, *defaultValue)
		if notImplt != nil {
			deci = defaultValue
			continue
		}
		defaultValue = deci
	}

	return model.NewDecimal(deci.Round(4)), nil
}

func (m *PluginManager) GetOrderShippingTaxRate(orDer order.Order, shippingPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
	defaultValue, appErr := m.srv.CheckoutService().BaseTaxRate(&shippingPrice)
	if appErr != nil {
		return nil, appErr
	}

	var (
		deci     *decimal.Decimal
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		deci, notImplt = plg.GetOrderShippingTaxRate(orDer, *defaultValue)
		if notImplt != nil {
			deci = defaultValue
			continue
		}
		defaultValue = deci
	}

	return model.NewDecimal(deci.Round(4)), nil
}

func (m *PluginManager) CalculateOrderlineTotal(orDer order.Order, orderLine order.OrderLine, variant product_and_discount.ProductVariant, product product_and_discount.Product) (interface{}, *model.AppError) {
	defaultValue, appErr := m.srv.CheckoutService().BaseOrderLineTotal(&orderLine)
	if appErr != nil {
		return nil, appErr
	}

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		taxedMoney, notImplt = plg.CalculateOrderLineTotal(&orDer, &orderLine, variant, product, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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
	defaultValue, appErr := m.srv.CheckoutService().BaseCheckoutLineUnitPrice(&totalLinePrice, quantity)
	if appErr != nil {
		return nil, appErr
	}

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		taxedMoney, notImplt = plg.CalculateCheckoutLineUnitPrice(checkoutInfo, lines, checkoutLineInfo, address, discounts, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		taxedMoney, notImplt = plg.CalculateOrderLineUnit(orDer, orderLine, variant, product, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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
	defaultValue, appErr := m.srv.CheckoutService().BaseTaxRate(&unitPrice)
	if appErr != nil {
		return nil, appErr
	}

	var (
		deci     *decimal.Decimal
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		deci, notImplt = plg.GetCheckoutLineTaxRate(&checkoutInfo, lines, checkoutLineInfo, address, discounts, *defaultValue)
		if notImplt != nil {
			deci = defaultValue
			continue
		}
		defaultValue = deci
	}

	return model.NewDecimal(deci.RoundUp(4)), nil
}

func (m *PluginManager) GetOrderLineTaxRate(orDer order.Order, product product_and_discount.Product, variant product_and_discount.ProductVariant, address *account.Address, unitPrice goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
	defaultValue, appErr := m.srv.CheckoutService().BaseTaxRate(&unitPrice)
	if appErr != nil {
		return nil, appErr
	}

	var (
		deci     *decimal.Decimal
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		deci, notImplt = plg.GetOrderLineTaxRate(orDer, product, variant, address, *defaultValue)
		if notImplt != nil {
			deci = defaultValue
			continue
		}
		defaultValue = deci
	}

	return model.NewDecimal(deci.RoundUp(4)), nil
}

func (m *PluginManager) GetTaxRateTypeChoices() []*model.TaxType {
	defaultValue := []*model.TaxType{}

	var (
		taxTypes []*model.TaxType
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		taxTypes, notImplt = plg.GetTaxRateTypeChoices(defaultValue)
		if notImplt != nil {
			taxTypes = defaultValue
			continue
		}
		defaultValue = taxTypes
	}

	return taxTypes
}

func (m *PluginManager) ShowTaxesOnStoreFront() bool {
	defaultValue := false

	var (
		showTax  bool
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		showTax, notImplt = plg.ShowTaxesOnStorefront(defaultValue)
		if notImplt != nil {
			showTax = defaultValue
			continue
		}
		defaultValue = showTax
	}

	return showTax
}

func (m *PluginManager) ApplyTaxesToProduct(product product_and_discount.Product, price goprices.Money, country string, channelID string) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(nil, goprices.Up)

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(channelID, true) {
		taxedMoney, notImplt = plg.ApplyTaxesToProduct(product, price, country, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
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
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(channelID, true) {
		taxedMoney, notImplt = plg.ApplyTaxesToShipping(price, shippingAddress, *defaultValue)
		if notImplt != nil {
			taxedMoney = defaultValue
			continue
		}
		defaultValue = taxedMoney
	}

	quantizedTaxedMoney, err := taxedMoney.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("ApplyTaxesToShipping", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return quantizedTaxedMoney, nil
}

func (m *PluginManager) PreprocessOrderCreation(checkoutInfo checkout.CheckoutInfo, discounts []*product_and_discount.DiscountInfo, lines checkout.CheckoutLineInfos) interface{} {
	var defaultValue interface{} = nil

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(checkoutInfo.Channel.Id, true) {
		value, notImplt = plg.PreprocessOrderCreation(checkoutInfo, discounts, lines, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) CustomerCreated(customer account.User) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.CustomerCreated(customer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) CustomerUpdated(customer account.User) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.CustomerUpdated(customer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ProductCreated(product product_and_discount.Product) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.ProductCreated(product, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ProductUpdated(product product_and_discount.Product) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.ProductUpdated(product, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ProductDeleted(product product_and_discount.Product, variants []int) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.ProductDeleted(product, variants, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ProductVariantCreated(variant product_and_discount.ProductVariant) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.ProductVariantCreated(variant, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ProductVariantUpdated(variant product_and_discount.ProductVariant) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.ProductVariantUpdated(variant, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ProductVariantDeleted(variant product_and_discount.ProductVariant) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.ProductVariantDeleted(variant, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ProductVariantOutOfStock(stock warehouse.Stock) {
	var defaultValue interface{}

	var (
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		notImplt = plg.ProductVariantOutOfStock(stock, defaultValue)
		if notImplt != nil {
			continue
		}
	}
}

func (m *PluginManager) ProductVariantBackInStock(stock warehouse.Stock) {
	var defaultValue interface{}

	var (
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		notImplt = plg.ProductVariantBackInStock(stock, defaultValue)
		if notImplt != nil {
			continue
		}
	}
}

func (m *PluginManager) OrderCreated(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.OrderCreated(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) OrderConfirmed(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.OrderConfirmed(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) DraftOrderCreated(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.DraftOrderCreated(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) DraftOrderDeleted(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.DraftOrderDeleted(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) DraftOrderUpdated(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.DraftOrderUpdated(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) SaleCreated(sale product_and_discount.Sale, currentCatalogue product_and_discount.NodeCatalogueInfo) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.SaleCreated(sale, currentCatalogue, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) SaleDeleted(sale product_and_discount.Sale, previousCatalogue product_and_discount.NodeCatalogueInfo) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.SaleDeleted(sale, previousCatalogue, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) SaleUpdated(sale product_and_discount.Sale, previousCatalogue, currentCatalogue product_and_discount.NodeCatalogueInfo) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.SaleUpdated(sale, previousCatalogue, currentCatalogue, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) InvoiceRequest(orDer order.Order, inVoice invoice.Invoice, number string) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.InvoiceRequest(orDer, inVoice, number, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) InvoiceDelete(inVoice invoice.Invoice) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var channelID string
	if inVoice.OrderID != nil {
		orDer, appErr := m.srv.OrderService().OrderById(*inVoice.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = orDer.ChannelID
	}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(channelID, true) {
		value, notImplt = plg.InvoiceDelete(inVoice, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) InvoiceSent(inVoice invoice.Invoice, email string) (interface{}, *model.AppError) {
	var defaultValue interface{}

	var channelID string
	if inVoice.OrderID != nil {
		orDer, appErr := m.srv.OrderService().OrderById(*inVoice.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		channelID = orDer.ChannelID
	}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(channelID, true) {
		value, notImplt = plg.InvoiceSent(inVoice, email, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) OrderFullyPaid(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.OrderFullyPaid(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) OrderUpdated(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.OrderUpdated(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) OrderCancelled(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.OrderCancelled(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) OrderFulfilled(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.OrderFulfilled(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) FulfillmentCreated(fulfillment order.Fulfillment) (interface{}, *model.AppError) {
	var defaultValue interface{}

	orDer, appErr := m.srv.OrderService().OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.FulfillmentCreated(fulfillment, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) FulfillmentCanceled(fulfillment order.Fulfillment) (interface{}, *model.AppError) {
	var defaultValue interface{}

	orDer, appErr := m.srv.OrderService().OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(orDer.ChannelID, true) {
		value, notImplt = plg.FulfillmentCanceled(fulfillment, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) CheckoutCreated(checkOut checkout.Checkout) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(checkOut.ChannelID, true) {
		value, notImplt = plg.CheckoutCreated(checkOut, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) CheckoutUpdated(checkOut checkout.Checkout) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins(checkOut.ChannelID, true) {
		value, notImplt = plg.CheckoutUpdated(checkOut, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) PageCreated(paGe page.Page) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.PageCreated(paGe, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) PageUpdated(paGe page.Page) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.PageUpdated(paGe, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) PageDeleted(paGe page.Page) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.PageDeleted(paGe, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) getPlugin(pluginID string, channelID string) BasePluginInterface {
	for _, plg := range m.AllPlugins {
		if plg.CheckPluginId(pluginID) && (channelID == "" || plg.ChannelId() == channelID) {
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
			value    *payment.GatewayResponse
			notImplt *PluginMethodNotImplemented
		)

		switch methodName {
		case "authorize_payment":
			value, notImplt = plg.AuthorizePayment(paymentInformation, nil)
		case "capture_payment":
			value, notImplt = plg.CapturePayment(paymentInformation, nil)
		case "refund_payment":
			value, notImplt = plg.RefundPayment(paymentInformation, nil)
		case "void_payment":
			value, notImplt = plg.VoidPayment(paymentInformation, nil)
		case "confirm_payment":
			value, notImplt = plg.ConfirmPayment(paymentInformation, nil)
		case "process_payment":
			value, notImplt = plg.ProcessPayment(paymentInformation, nil)

		default:
			return nil, fmt.Errorf("no method found")
		}

		if notImplt == nil && value != nil {
			return value, nil
		}
	}

	return nil, fmt.Errorf("Payment plugin %s for %s payment method is in-accessible", gateway, methodName)
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

func (m *PluginManager) TokenIsRequiredAsPaymentInput(gateway, channelID string) bool {
	plg := m.getPlugin(gateway, channelID)
	defaultValue := true

	if plg != nil {
		value, _ := plg.TokenIsRequiredAsPaymentInput(defaultValue)
		// ignore not implement since it is default to nil
		return value
	}

	return defaultValue
}

func (m *PluginManager) GetClientToken(gateway string, tokenConfig payment.TokenConfig, channelID string) string {
	plg := m.getPlugin(gateway, channelID)
	if plg != nil {
		value, _ := plg.GetClientToken(tokenConfig, nil)
		return value
	}

	return ""
}

func (m *PluginManager) ListPaymentSources(gateway, customerID, channelID string) ([]*payment.CustomerSource, error) {
	plg := m.getPlugin(gateway, channelID)
	if plg != nil {
		defaultValue := []*payment.CustomerSource{}
		return plg.ListPaymentSources(customerID, defaultValue)
	}

	return nil, fmt.Errorf("Payment plugin %s is inaccessible", gateway)
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
		value, notImplt := plg.GetPaymentGateways(currency, checkOut, nil)
		if notImplt != nil { // this indicates the plugin does not implement its own method
			continue
		}
		gateways = append(gateways, value...)
	}

	return gateways
}

func (m *PluginManager) ListExternalAuthentications(activeOnly bool) []model.StringInterface {
	plugins := m.getPlugins("", activeOnly)

	res := []model.StringInterface{}

	for _, plg := range plugins {
		_, notImplt := plg.ExternalObtainAccessTokens(nil, nil, ExternalAccessTokens{})
		if notImplt == nil {
			manifest := plg.GetManifest()
			res = append(res, model.StringInterface{
				"id":   manifest.PluginID,
				"name": manifest.PluginName,
			})
		}
	}

	return res
}

// AssignTaxCodeToObjectMeta requires obj must be Product or ProductType
func (m *PluginManager) AssignTaxCodeToObjectMeta(obj interface{}, taxCode string) (*model.TaxType, *model.AppError) {

	// validate obj
	switch obj.(type) {
	case product_and_discount.Product,
		product_and_discount.ProductType,
		*product_and_discount.Product,
		**product_and_discount.ProductType:
	default:
		return nil, model.NewAppError("AssignTaxCodeToObjectMeta", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	var (
		defaultValue = new(model.TaxType)
		value        *model.TaxType
		notImplt     *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.AssignTaxCodeToObjectMeta(obj, taxCode, *defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
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
		**product_and_discount.ProductType:
	default:
		return nil, model.NewAppError("GetTaxCodeFromObjectMeta", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "obj"}, "obj must be either Product or ProductType", http.StatusBadRequest)
	}

	var (
		defaultValue = new(model.TaxType)
		value        *model.TaxType
		notImplt     *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.GetTaxCodeFromObjectMeta(obj, *defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value, nil
}

func (m *PluginManager) GetTaxRatePercentageValue(obj interface{}, country string) {
	panic("not implemented")
}

func (m *PluginManager) SavePluginConfiguration(pluginID, channelID string, cleanedData model.StringInterface) (*plugins.PluginConfiguration, *model.AppError) {
	if !model.IsValidId(channelID) {
		return nil, model.NewAppError("SavePluginConfiguration", app.InvalidArgumentAppErrorID, nil, "", http.StatusBadRequest)
	}

	var pluginList []BasePluginInterface
	if channelID != "" {
		pluginList = m.getPlugins(channelID, true)
	} else {
		pluginList = m.AllPlugins
	}

	for _, plg := range pluginList {
		manifest := plg.GetManifest()
		if manifest.PluginID == pluginID {

			// try get or create plugin configuration
			pluginConfig, appErr := m.srv.PluginService().GetPluginConfiguration(&plugins.PluginConfigurationFilterOptions{
				Identifier: squirrel.Eq{m.srv.Store.PluginConfiguration().TableName("Identifier"): pluginID},
				ChannelID:  squirrel.Eq{m.srv.Store.PluginConfiguration().TableName("ChannelID"): channelID},
			})
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError {
					return nil, appErr
				}

				pluginConfig, appErr = m.srv.PluginService().UpsertPluginConfiguration(&plugins.PluginConfiguration{
					Identifier:    pluginID,
					ChannelID:     channelID,
					Configuration: plg.GetConfiguration(),
				})
				if appErr != nil {
					return nil, appErr
				}
			}

			pluginConfig, appErr, notImplt := plg.SavePluginConfiguration(pluginConfig, cleanedData)
			if notImplt != nil {
				m.srv.Log.Warn("Method not implemented", slog.Err(notImplt))
			}
			if appErr != nil {
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

func (m *PluginManager) FetchTaxesData() bool {
	defaultValue := false

	var (
		value    bool
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.FetchTaxesData(defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return defaultValue
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

	value, notImplt := plg.Webhook(req, path, defaultValue)
	if notImplt != nil {
		return nil, model.NewAppError("WebhookEndpointWithoutChannel", "modules.plugins.method_not_implemented.app_error", nil, notImplt.Error(), http.StatusNotImplemented)
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

	res, notImplt := plg.Webhook(req, path, *defaultValue)
	if notImplt != nil {
		return nil, model.NewAppError("Webhook", "modules.plugins.method_not_implemented.app_error", nil, notImplt.Error(), http.StatusNotImplemented)
	}
	return res, nil
}

func (m *PluginManager) Notify(event string, payload model.StringInterface, channelID string, pluginID string) interface{} {
	var defaultValue interface{}

	if pluginID != "" {
		plg := m.getPlugin(pluginID, channelID)
		value, notImplt := plg.Notify(event, payload, defaultValue)
		if notImplt != nil {
			value = defaultValue
		}
		return value
	}

	for _, plg := range m.getPlugins(channelID, true) {
		value, notImplt := plg.Notify(event, payload, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return defaultValue
}

func (m *PluginManager) ExternalObtainAccessTokens(pluginID string, data model.StringInterface, req *http.Request) (*ExternalAccessTokens, *model.AppError) {
	var defaultValue ExternalAccessTokens
	plg := m.getPlugin(pluginID, "")

	res, notImplt := plg.ExternalObtainAccessTokens(data, req, defaultValue)
	if notImplt != nil {
		return nil, model.NewAppError("ExternalObtainAccessTokens", "modules.plugins.method_not_implemented.app_error", nil, notImplt.Error(), http.StatusNotImplemented)
	}

	return res, nil
}

func (m *PluginManager) ExternalAuthenticationUrl(pluginID string, data model.StringInterface, req *http.Request) model.StringInterface {
	defaultValue := model.StringInterface{}

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		res, notImplt := plg.ExternalAuthenticationUrl(data, req, defaultValue)
		if notImplt != nil {
			return defaultValue
		}
		return res
	}

	return defaultValue
}

func (m *PluginManager) ExternalRefresh(pluginID string, data model.StringInterface, req *http.Request) ExternalAccessTokens {
	defaultValue := ExternalAccessTokens{}

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		res, notImplt := plg.ExternalRefresh(data, req, defaultValue)
		if notImplt != nil {
			return defaultValue
		}
		return *res
	}

	return defaultValue
}

func (m *PluginManager) AuthenticateUser(req *http.Request) *account.User {
	var (
		defaultValue *account.User = nil
		value        *account.User
		notImplt     *PluginMethodNotImplemented
	)

	for _, plg := range m.getPlugins("", true) {
		value, notImplt = plg.AuthenticateUser(req, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) ExternalLogout(pluginID string, data model.StringInterface, req *http.Request) model.StringInterface {
	defaultValue := model.StringInterface{}

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		notImplt := plg.ExternalLogout(data, req, defaultValue)
		if notImplt != nil {
			return defaultValue
		}
	}

	return defaultValue
}

func (m *PluginManager) ExternalVerify(pluginID string, data model.StringInterface, req *http.Request) (*account.User, model.StringInterface) {
	var (
		defaultData = model.StringInterface{}
		defaultUser *account.User
	)

	plg := m.getPlugin(pluginID, "")
	if plg != nil {
		user, data, notImplt := plg.ExternalVerify(data, req, AType{
			User: defaultUser,
			Data: defaultData,
		})
		if notImplt != nil {
			return defaultUser, defaultData
		}
		return user, data
	}

	return defaultUser, defaultData
}
