/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package discount

import (
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
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

// Decorator returns a function to calculate discount
func (*ServiceDiscount) Decorator(preValue any) types.DiscountCalculator {
	if preValue == nil {
		return nil
	}

	return func(first any, fromGross *bool) (any, error) {
		// validate first
		switch first.(type) {
		case *goprices.Money,
			*goprices.MoneyRange,
			*goprices.TaxedMoney,
			*goprices.TaxedMoneyRange:
		default:
			return nil, model_helper.NewAppError("DiscountCalculator", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "first"}, "", http.StatusBadRequest)
		}

		if fromGross == nil { // fixed discount
			switch t := first.(type) {
			case *goprices.Money:
				return goprices.FixedDiscount[*goprices.Money](t, preValue.(*goprices.Money))
			case *goprices.MoneyRange:
				return goprices.FixedDiscount[*goprices.MoneyRange](t, preValue.(*goprices.Money))
			case *goprices.TaxedMoney:
				return goprices.FixedDiscount[*goprices.TaxedMoney](t, preValue.(*goprices.Money))
			case *goprices.TaxedMoneyRange:
				return goprices.FixedDiscount[*goprices.TaxedMoneyRange](t, preValue.(*goprices.Money))
			default:
				return nil, nil
			}
		}

		flt, _ := preValue.(*decimal.Decimal).Float64()
		switch t := first.(type) {
		case *goprices.Money:
			return goprices.PercentageDiscount[*goprices.Money](t, flt, *fromGross)
		case *goprices.MoneyRange:
			return goprices.PercentageDiscount[*goprices.MoneyRange](t, flt, *fromGross)
		case *goprices.TaxedMoney:
			return goprices.PercentageDiscount[*goprices.TaxedMoney](t, flt, *fromGross)
		case *goprices.TaxedMoneyRange:
			return goprices.PercentageDiscount[*goprices.TaxedMoneyRange](t, flt, *fromGross)
		default:
			return nil, nil
		}
	}
}
