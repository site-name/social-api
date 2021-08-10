package discount

import (
	"net/http"
	"sync"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
)

type AppDiscount struct {
	app.AppIface
	wg    sync.WaitGroup
	mutex sync.Mutex
}

func init() {
	app.RegisterDiscountApp(func(a app.AppIface) sub_app_iface.DiscountApp {
		return &AppDiscount{
			AppIface: a,
		}
	})
}

// DiscountCalculator number of `args` must be 1 or 2
//
//  if len(args) == 1 {
//		args[0].(type) == (*Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange)
//  }
//  if len(args) == 2 {
//		args[0].(type) == (*Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange) && args[1].(type) == bool
//  }
type DiscountCalculator func(args ...interface{}) (interface{}, error)

func Decorator(preValue interface{}) DiscountCalculator {
	return func(args ...interface{}) (interface{}, error) {
		// validating number of args
		if l := len(args); l < 1 || l > 2 {
			return nil, model.NewAppError("app.Discount.decorator", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "args"}, "you must provide either 1 or 2 arguments", http.StatusBadRequest)
		}

		if len(args) == 1 { // fixed discount
			return goprices.FixedDiscount(args[0], preValue.(*goprices.Money))
		}
		return goprices.PercentageDiscount(args[0], preValue, args[1].(bool))
	}
}
