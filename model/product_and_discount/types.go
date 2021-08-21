package product_and_discount

import goprices "github.com/site-name/go-prices"

type ProductAvailability struct {
	OnSale                  bool
	PriceRange              *goprices.TaxedMoneyRange // can be nil
	PriceRangeUnDiscounted  *goprices.TaxedMoneyRange // can be nil
	Discount                *goprices.TaxedMoney      // can be nil
	PriceRangeLocalCurrency *goprices.TaxedMoneyRange // can be nil
	DiscountLocalCurrency   *goprices.TaxedMoneyRange // can be nil
}

type VariantAvailability struct {
	OnSale                bool
	Price                 goprices.TaxedMoney
	PriceUnDiscounted     goprices.TaxedMoney
	Discount              *goprices.TaxedMoney // can be nil
	PriceLocalCurrency    *goprices.TaxedMoney // can be nil
	DiscountLocalCurrency *goprices.TaxedMoney // can be nil
}
