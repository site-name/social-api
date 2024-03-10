package model_helper

import (
	"sort"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

type ProductAvailability struct {
	OnSale                  bool
	PriceRange              *goprices.TaxedMoneyRange // can be nil
	PriceRangeUnDiscounted  *goprices.TaxedMoneyRange // can be nil
	Discount                *goprices.TaxedMoney      // can be nil
	PriceRangeLocalCurrency *goprices.TaxedMoneyRange // can be nil
	DiscountLocalCurrency   *goprices.TaxedMoney      // can be nil
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
	Sale            any // either *Sale || *Voucher
	ChannelListings map[string]*model.SaleChannelListing
	ProductIDs      util.AnyArray[string]
	CategoryIDs     util.AnyArray[string]
	CollectionIDs   util.AnyArray[string]
	VariantsIDs     util.AnyArray[string]
}

func (d *DiscountInfo) IsSaleValid() bool {
	switch d.Sale.(type) {
	case *model.Sale, *model.Voucher, model.Sale, model.Voucher:
		return true

	default:
		return false
	}
}

type CostsData struct {
	costs   []goprices.Money
	margins []float64
}

func NewCostsData(costs []goprices.Money, margins []float64) *CostsData {
	sort.SliceStable(costs, func(i, j int) bool {
		return costs[i].LessThan(costs[j])
	})
	sort.Float64s(margins)

	return &CostsData{
		costs:   costs,
		margins: margins,
	}
}

func (c *CostsData) Costs() []goprices.Money {
	return c.costs
}

func (c *CostsData) Margins() []float64 {
	return c.margins
}

type NodeCatalogueInfo map[string][]string
