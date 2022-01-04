package plugins

import (
	"net/http"

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
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
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
		ChannelID: squirrel.Eq{
			m.srv.Store.PluginConfiguration().TableName("ChannelID"): channels.IDs(),
		},
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

func (m *PluginManager) filterPlugins(channelID string, active bool) []BasePluginInterface {
	res := []BasePluginInterface{}

	for _, plg := range m.AllPlugins {

		if plg.IsActive() {
			res = append(res, plg)
		}
	}

	return res
}

func (m *PluginManager) ChangeUserAddress(address *account.Address, addressType string, user *account.User) *account.Address {
	var (
		notImplt      *PluginMethodNotImplemented
		previousValue account.Address = *address
	)

	for _, plg := range m.filterPlugins() {
		address, notImplt = plg.ChangeUserAddress(address, addressType, user, &previousValue)
		if notImplt != nil {
			address = &previousValue
			continue
		}
		previousValue = *address
	}

	return address
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
	for _, plg := range m.filterPlugins() {
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

	for _, plg := range m.filterPlugins() {
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

	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
		showTax, notImplt = plg.ShowTaxesOnStorefront(defaultValue)
		if notImplt != nil {
			showTax = defaultValue
			continue
		}
		defaultValue = showTax
	}

	return showTax
}

func (m *PluginManager) ApplyTaxesToProduct(product product_and_discount.Product, price goprices.Money, country string) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(nil, goprices.Up)

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.filterPlugins() {
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

func (m *PluginManager) ApplyTaxesToShipping(price goprices.Money, shippingAddress account.Address) (*goprices.TaxedMoney, *model.AppError) {
	defaultValue, _ := (&goprices.TaxedMoney{
		Net:      &price,
		Gross:    &price,
		Currency: price.Currency,
	}).Quantize(nil, goprices.Up)

	var (
		taxedMoney *goprices.TaxedMoney
		notImplt   *PluginMethodNotImplemented
	)
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
		value, notImplt = plg.InvoiceRequest(orDer, inVoice, number, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) InvoiceDelete(inVoice invoice.Invoice) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.filterPlugins() {
		value, notImplt = plg.InvoiceDelete(inVoice, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) InvoiceSent(inVoice invoice.Invoice, email string) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.filterPlugins() {
		value, notImplt = plg.InvoiceSent(inVoice, email, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}

func (m *PluginManager) OrderFullyPaid(orDer order.Order) interface{} {
	var defaultValue interface{}

	var (
		value    interface{}
		notImplt *PluginMethodNotImplemented
	)
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
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
	for _, plg := range m.filterPlugins() {
		value, notImplt = plg.OrderFulfilled(orDer, defaultValue)
		if notImplt != nil {
			value = defaultValue
			continue
		}
		defaultValue = value
	}

	return value
}
