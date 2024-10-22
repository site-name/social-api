package product

import (
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

func GetProductCostsData(
	variantChannelListings []*model.ProductVariantChannelListing,
	hasVariants bool,
	currency string,
) (*goprices.MoneyRange, []float64, *model_helper.AppError) {
	purchaseCostsRange, _ := util.ZeroMoneyRange(currency)
	margins := []float64{0.0, 0.0}

	if !hasVariants {
		return purchaseCostsRange, margins, nil
	}

	costsData := GetCostDataFromVariantChannelListing(variantChannelListings)
	if len(costsData.Costs()) > 0 {
		purchaseCostsRange, _ = goprices.NewMoneyRange(util.MinMaxMoneyInMoneySlice(costsData.Costs()))
	}
	if length := len(costsData.Margins()); length > 0 {
		margins = []float64{costsData.Margins()[0], costsData.Margins()[length-1]}
	}

	return purchaseCostsRange, margins, nil
}

func GetCostDataFromVariantChannelListing(variantChannelListings []*model.ProductVariantChannelListing) *model.CostsData {
	var (
		costs   []*goprices.Money
		margins []float64
	)
	for _, listing := range variantChannelListings {
		costsData := GetvariantCostsData(listing)
		costs = append(costs, costsData.Costs()...)
		margins = append(margins, costsData.Margins()...)
	}

	return model.NewCostsData(costs, margins)
}

func GetvariantCostsData(variantChannelListing *model.ProductVariantChannelListing) *model.CostsData {
	var (
		costs   []*goprices.Money
		margins []float64
	)
	costs = append(costs, GetCostPrice(variantChannelListing))
	if margin := GetMarginForVariantChannelListing(variantChannelListing); margin != nil {
		margins = append(margins, *margin)
	}

	return model.NewCostsData(costs, margins)
}

func GetCostPrice(variantChannelListing *model.ProductVariantChannelListing) *goprices.Money {
	variantChannelListing.PopulateNonDbFields()

	if variantChannelListing.CostPrice == nil {
		return &goprices.Money{
			Amount:   decimal.Zero,
			Currency: variantChannelListing.Currency,
		}
	}

	return variantChannelListing.CostPrice
}

func GetMarginForVariantChannelListing(variantChannelListing *model.ProductVariantChannelListing) *float64 {
	variantChannelListing.PopulateNonDbFields()

	if variantChannelListing.CostPrice == nil || variantChannelListing.Price == nil {
		return nil
	}

	margin, _ := variantChannelListing.Price.Sub(variantChannelListing.CostPrice)
	fl64Percent := margin.Amount.
		Div(*variantChannelListing.PriceAmount).
		Mul(decimal.NewFromInt(100)).
		Round(0).
		InexactFloat64()

	return &fl64Percent
}
