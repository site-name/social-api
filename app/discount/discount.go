/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package discount

import (
	"errors"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceDiscount struct {
	srv *app.Server
}

func init() {
	app.RegisterDiscountService(func(s *app.Server) (sub_app_iface.DiscountService, error) {
		return &ServiceDiscount{
			srv: s,
		}, nil
	})
}

// Decorator returns a function to calculate discount
func Decorator(preValue interface{}) types.DiscountCalculator {
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
