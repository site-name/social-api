package discount

import (
	"errors"
	"sync"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppDiscount struct {
	app.AppIface
	wg    sync.WaitGroup // this is for some methods that need concurrent executions
	mutex sync.Mutex     // this is for prevent data racing in methods that have concurrent executions
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
//  // pass 1 argument if you want to calculate fixed discount
//  if len(args) == 1 {
//		args[0].(type) == (*Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange)
//  }
//
//  // pass 2 arguments if you want to calculate percentage discount
//  if len(args) == 2 {
//		args[0].(type) == (*Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange) && args[1].(type) == bool
//  }
type DiscountCalculator func(args ...interface{}) (interface{}, error)

// Decorator returns a function to calculate discount
func Decorator(preValue interface{}) DiscountCalculator {
	return func(args ...interface{}) (interface{}, error) {
		// validating number of args
		if l := len(args); l < 1 || l > 2 {
			return nil, errors.New("at most 2 arguments only")
		}

		if len(args) == 1 { // fixed discount
			return goprices.FixedDiscount(args[0], preValue.(*goprices.Money))
		}
		return goprices.PercentageDiscount(args[0], preValue, args[1].(bool))
	}
}
