package product

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
)

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

func (a *AppProduct) getProductPriceRange(discounted interface{}, unDiscounted interface{}, localCurrency string) (
	*struct {
		goprices.TaxedMoneyRange
		goprices.TaxedMoney
	},
	error,
) {
	// validate provided currency is valid
	if _, err := goprices.GetCurrencyPrecision(localCurrency); err != nil {
		// return nil,
	}
	panic("not implemented")
}

// GetVariantPrice
func (a *AppProduct) GetVariantPrice(
	variant *product_and_discount.ProductVariant,
	variantChannelListing *product_and_discount.ProductVariantChannelListing,
	product *product_and_discount.Product,
	collections []*product_and_discount.Collection,
	discounts []*product_and_discount.DiscountInfo,
	chanNel *channel.Channel,

) (*goprices.Money, *model.AppError) {

	variantChannelListing.PopulateNonDbFields() // must call this initially

	return a.DiscountApp().CalculateDiscountedPrice(
		product,
		variantChannelListing.Price,
		collections,
		discounts,
		chanNel,
	)
}

func (a *AppProduct) GetProductPriceRange(
	product *product_and_discount.Product,
	variants []*product_and_discount.ProductVariant,
	variantsChannelListing []*product_and_discount.ProductVariantChannelListing,
	collections []*product_and_discount.Collection,
	discounts []*product_and_discount.DiscountInfo,
	chanNel *channel.Channel,

) ([]*goprices.MoneyRange, *model.AppError) {

	// filter nil values (if exist) from variantsChannelListing
	for i, item := range variantsChannelListing {
		if item == nil {
			variantsChannelListing = append(variantsChannelListing[:i], variantsChannelListing[i+1:]...)
		}
	}
	if len(variants) > 0 {
		variantChannelListingsMap := model.MakeStringMapForModelSlice(
			variantsChannelListing,
			func(i interface{}) string {
				return i.(*product_and_discount.ProductVariantChannelListing).VariantID
			},
			nil,
		)

		prices := []*goprices.Money{}
		for _, variant := range variants {
			variantChannelListing := variantChannelListingsMap[variant.Id]
			if variantChannelListing != nil {
				price, appErr := a.GetVariantPrice(
					variant,
					variantChannelListing.(*product_and_discount.ProductVariantChannelListing),
					product,
					collections,
					discounts,
					chanNel,
				)
				if appErr != nil {
					return nil, appErr
				}

				prices = append(prices, price)
			}
		}

		if len(prices) > 0 {
			panic("not implemented")
		}
	}

	return nil, nil
}
