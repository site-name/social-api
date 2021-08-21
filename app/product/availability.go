package product

import goprices "github.com/site-name/go-prices"

// getTotalDiscountFromRange Calculate the discount amount between two TaxedMoneyRange.
//
// Subtract two prices and return their total discount, if any.
// Otherwise, it returns None.
func getTotalDiscountFromRange(undiscounted *goprices.TaxedMoneyRange, discounted *goprices.TaxedMoneyRange) (*goprices.TaxedMoney, error) {
	return getTotalDiscount(undiscounted.Start, discounted.Start)
}

// getTotalDiscount  Calculate the discount amount between two TaxedMoney.
//
// Subtract two prices and return their total discount, if any.
// Otherwise, it returns None.
func getTotalDiscount(unDiscounted *goprices.TaxedMoney, discounted *goprices.TaxedMoney) (*goprices.TaxedMoney, error) {
	less, err := discounted.LessThan(unDiscounted)
	if err != nil {
		return nil, err
	}
	if less {
		return unDiscounted.Sub(discounted)
	}

	return nil, nil
}

func getProductPriceRange(discounted interface{}, unDiscounted interface{}, localCurrency string) (
	*struct {
		goprices.TaxedMoneyRange
		goprices.TaxedMoney
	},
	error,
) {
	if _, err := goprices.GetCurrencyPrecision(localCurrency); err != nil {

	}
}
