package product_and_discount

import (
	"sort"

	goprices "github.com/site-name/go-prices"
)

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

// DiscountInfo contains information of a discount
type DiscountInfo struct {
	Sale            interface{} // either *Sale || *Voucher
	ChannelListings map[string]*SaleChannelListing
	ProductIDs      []string
	CategoryIDs     []string
	CollectionIDs   []string
	VariantsIDs     []string
}

type CostsData struct {
	costs   []*goprices.Money
	margins []float64
}

func NewCostsData(costs []*goprices.Money, margins []float64) *CostsData {
	// sorting:
	sort.Slice(costs, func(i, j int) bool {
		less, err := costs[i].LessThan(costs[j])
		return less && err == nil
	})
	sort.Float64s(margins)

	c := &CostsData{
		costs:   costs,
		margins: margins,
	}
	return c
}

func (c *CostsData) Costs() []*goprices.Money {
	return c.costs
}

func (c *CostsData) Margins() []float64 {
	return c.margins
}

type NodeCatalogueInfo map[string][]string
