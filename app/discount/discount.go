/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package discount

import (
	"fmt"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
)

type ServiceDiscount struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Discount = &ServiceDiscount{s}
		return nil
	})
}

// Decorator returns a function to calculate discount.
// `preValue` must has type of goprices.Money || decimal.Decimal
func (*ServiceDiscount) Decorator(preValue any) types.DiscountCalculator {
	if preValue == nil {
		return nil
	}

	return func(first any, fromGross *bool) (any, error) {
		if fromGross == nil { // fixed discount
			switch t := first.(type) {
			case goprices.Money:
				return goprices.FixedDiscount[*goprices.Money](&t, preValue.(goprices.Money))
			case goprices.MoneyRange:
				return goprices.FixedDiscount[*goprices.MoneyRange](&t, preValue.(goprices.Money))
			case goprices.TaxedMoney:
				return goprices.FixedDiscount[*goprices.TaxedMoney](&t, preValue.(goprices.Money))
			case goprices.TaxedMoneyRange:
				return goprices.FixedDiscount[*goprices.TaxedMoneyRange](&t, preValue.(goprices.Money))
			default:
				return nil, fmt.Errorf("invalid first value provided with type: %T", first)
			}
		}

		f64 := preValue.(decimal.Decimal).InexactFloat64()
		switch t := first.(type) {
		case goprices.Money:
			return goprices.PercentageDiscount[*goprices.Money](&t, f64, *fromGross)
		case goprices.MoneyRange:
			return goprices.PercentageDiscount[*goprices.MoneyRange](&t, f64, *fromGross)
		case goprices.TaxedMoney:
			return goprices.PercentageDiscount[*goprices.TaxedMoney](&t, f64, *fromGross)
		case goprices.TaxedMoneyRange:
			return goprices.PercentageDiscount[*goprices.TaxedMoneyRange](&t, f64, *fromGross)
		default:
			return nil, fmt.Errorf("invalid first value provided with type: %T", first)
		}
	}
}
