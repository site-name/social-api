package plugins

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/modules/util"
)

type PluginManager struct {
	srv        *app.Server
	AllPlugins []BasePluginInterface // keys are channel id
}

func NewPluginManager(srv *app.Server, ch *channel.Channel) (*PluginManager, *model.AppError) {
	m := &PluginManager{
		srv: srv,
	}

	// finds a list of plugin configs belong to a specific channel
	pluginConfigsOfChannel, appErr := m.srv.PluginService().FilterPluginConfigurations(&plugins.PluginConfigurationFilterOptions{
		ChannelID:              squirrel.Eq{m.srv.Store.PluginConfiguration().TableName("ChannelID"): ch.Id},
		PrefetchRelatedChannel: true, //
	})
	if appErr != nil {
		return nil, appErr
	}

	var configsMap = map[string]*plugins.PluginConfiguration{}
	for _, config := range pluginConfigsOfChannel {
		configsMap[config.Identifier] = config
	}

	for _, pluginInitObj := range pluginInitObjects {

		existingConfig := configsMap[pluginInitObj.PluginID]

		plugin := pluginInitObj.NewPluginFunc(NewPluginConfig{
			Srv:           srv,
			Channel:       ch,
			Manager:       m,
			Configuration: existingConfig.Configuration,
			Active:        existingConfig.Active,
		})

		m.AllPlugins = append(m.AllPlugins, plugin)
	}

	return m, nil
}

func (m *PluginManager) ActivePlugins() []BasePluginInterface {
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

	for _, plg := range m.ActivePlugins() {
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
	for _, plg := range m.ActivePlugins() {
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

	for _, plg := range m.ActivePlugins() {
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

	for _, plg := range m.ActivePlugins() {
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
	for _, plg := range m.ActivePlugins() {
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
