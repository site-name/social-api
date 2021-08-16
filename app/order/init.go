package order

import (
	"net/http"
	"sync"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/modules/util"
)

type RecalculateOrderPricesFunc func(*order.Order, map[string]interface{}) *model.AppError

type AppOrder struct {
	app.AppIface
	wg    sync.WaitGroup
	mutex sync.Mutex

	RecalculateOrderPrices RecalculateOrderPricesFunc // This attribute is initialized as this app is started
}

func init() {
	app.RegisterOrderApp(func(a app.AppIface) sub_app_iface.OrderApp {
		orderApp := &AppOrder{
			AppIface: a,
		}

		orderApp.RecalculateOrderPrices = orderApp.UpdateVoucherDiscount(orderApp.DecoratedFunc)

		return orderApp
	})
}

// UpdateVoucherDiscount Recalculate order discount amount based on order voucher
func (a *AppOrder) UpdateVoucherDiscount(fun RecalculateOrderPricesFunc) RecalculateOrderPricesFunc {
	return func(ord *order.Order, kwargs map[string]interface{}) *model.AppError {
		if kwargs == nil {
			kwargs = make(map[string]interface{})
		}

		var (
			discount         interface{}
			notApplicableErr *model.NotApplicable
			appErr           *model.AppError
		)

		if kwargs["update_voucher_discount"] != nil || true {
			discount, notApplicableErr, appErr = a.GetVoucherDiscountForOrder(ord)
			if appErr != nil {
				appErr.Where = "UpdateVoucherDiscount"
				return appErr
			}
			if notApplicableErr != nil {
				discount, _ = util.ZeroMoney(ord.Currency)
			}
		}

		// set discount here
		kwargs["discount"] = discount

		return fun(ord, kwargs)
	}
}

func (a *AppOrder) DecoratedFunc(ord *order.Order, kwargs map[string]interface{}) (appErr *model.AppError) {
	defer func() {
		if appErr != nil {
			appErr.Where = "RecalculateOrderPrices"
		}
	}()

	ord.PopulateNonDbFields() // NOTE: must call this func before doing money calculations

	// avoid using prefetched order lines
	orderLines, apErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if apErr != nil {
		appErr = apErr
		return
	}

	totalPrice := ord.ShippingPrice
	for _, orderLine := range orderLines {
		orderLine.PopulateNonDbFields() // NOTE: call this before performing money calculations

		addedPrice, err := totalPrice.Add(orderLine.TotalPrice)
		if err != nil {
			appErr = model.NewAppError("RecalculateOrderPrices", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
			return
		}
		totalPrice = addedPrice
	}

	unDiscountedTotal, _ := goprices.NewTaxedMoney(totalPrice.Net, totalPrice.Gross)

	voucherDiscount, _ := util.ZeroMoney(ord.Currency)
	if discountIface := kwargs["discount"]; discountIface != nil {
		if discountValue, ok := discountIface.(*goprices.Money); ok {
			voucherDiscount = discountValue
		}
	}

	// discount amount can't be greater than order total
	if totalPrice.Gross.Amount.LessThan(*voucherDiscount.Amount) {
		voucherDiscount = totalPrice.Gross
	}
	subResult, err := totalPrice.Sub(voucherDiscount)
	if err != nil {
		appErr = model.NewAppError("RecalculateOrderPrices", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	totalPrice = subResult

	ord.Total = totalPrice
	ord.UnDiscountedTotal = unDiscountedTotal

	if !voucherDiscount.Amount.Equal(decimal.Zero) {
		assignedOrderDiscount, apErr := a.GetVoucherDiscountAssignedToOrder(ord)
		if apErr != nil {
			appErr = apErr
			return
		}

		if assignedOrderDiscount != nil {
			assignedOrderDiscount.AmountValue = voucherDiscount.Amount
			assignedOrderDiscount.Value = voucherDiscount.Amount
			_, appErr = a.DiscountApp().UpsertOrderDiscount(assignedOrderDiscount)
			if appErr != nil {
				return
			}
		}
	}

	return
}
