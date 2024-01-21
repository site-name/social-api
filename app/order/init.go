package order

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/order/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

func init() {
	app.RegisterService(func(s *app.Server) error {
		sv := &ServiceOrder{srv: s}
		sv.RecalculateOrderPrices = sv.UpdateVoucherDiscount(sv.decoratedFunc)
		s.Order = sv

		return nil
	})
}

type ServiceOrder struct {
	srv *app.Server

	RecalculateOrderPrices types.RecalculateOrderPricesFunc // This attribute is initialized as this app is started
}

// UpdateVoucherDiscount Recalculate order discount amount based on order voucher
func (a *ServiceOrder) UpdateVoucherDiscount(fun types.RecalculateOrderPricesFunc) types.RecalculateOrderPricesFunc {
	return func(transaction *gorm.DB, order *model.Order, kwargs map[string]interface{}) *model_helper.AppError {
		if kwargs == nil {
			kwargs = make(map[string]interface{})
		}

		var (
			discount          interface{}
			notApplicableErr  *model.NotApplicable
			appErr            *model_helper.AppError
			calculateDiscount bool
		)

		if item := kwargs["update_voucher_discount"]; item == nil {
			calculateDiscount = true
		} else {
			if boolItem, ok := item.(bool); ok {
				calculateDiscount = boolItem
			}
		}

		if calculateDiscount {
			discount, notApplicableErr, appErr = a.GetVoucherDiscountForOrder(order)
			if appErr != nil {
				return appErr
			}
			if notApplicableErr != nil {
				discount, _ = util.ZeroMoney(order.Currency)
			}
		}

		// set discount
		kwargs["discount"] = discount

		return fun(transaction, order, kwargs)
	}
}

func (a *ServiceOrder) decoratedFunc(transaction *gorm.DB, order *model.Order, kwargs map[string]interface{}) *model_helper.AppError {
	order.PopulateNonDbFields() // NOTE: must call this func before doing money calculations

	// avoid using prefetched order lines
	orderLines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": order.Id},
	})
	if appErr != nil {
		return appErr
	}

	totalPrice := order.ShippingPrice
	for _, orderLine := range orderLines {
		orderLine.PopulateNonDbFields() // NOTE: call this before performing money calculations

		addedPrice, err := totalPrice.Add(orderLine.TotalPrice)
		if err != nil {
			return model_helper.NewAppError("RecalculateOrderPrices", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
		totalPrice = addedPrice
		// reassign value here since `addedPrice` can be nil if error occurs.
		// so `totalPrice` becomes wrong
	}

	unDiscountedTotal, _ := goprices.NewTaxedMoney(totalPrice.Net, totalPrice.Gross) // ignore error here

	voucherDiscount, _ := util.ZeroMoney(order.Currency) // ignore error since order's Currency is validated before being insert into db
	if discountIface := kwargs["discount"]; discountIface != nil {
		if discountValue, ok := discountIface.(*goprices.Money); ok {
			voucherDiscount = discountValue
		}
	}

	// discount amount can't be greater than order total
	if totalPrice.Gross.Amount.LessThan(voucherDiscount.Amount) {
		voucherDiscount = totalPrice.Gross
	}
	subResult, err := totalPrice.Sub(voucherDiscount)
	if err != nil {
		return model_helper.NewAppError("RecalculateOrderPrices", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	totalPrice = subResult

	order.Total = totalPrice
	order.UnDiscountedTotal = unDiscountedTotal

	if !voucherDiscount.Amount.Equal(decimal.Zero) { // != 0.0
		assignedOrderDiscount, apErr := a.GetVoucherDiscountAssignedToOrder(order)
		if apErr != nil {
			return apErr
		}

		if assignedOrderDiscount != nil {
			assignedOrderDiscount.AmountValue = &voucherDiscount.Amount
			assignedOrderDiscount.Value = &voucherDiscount.Amount
			_, appErr = a.srv.DiscountService().UpsertOrderDiscount(transaction, assignedOrderDiscount)
			if appErr != nil {
				return appErr
			}
		}
	}

	return nil
}
