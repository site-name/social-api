package product

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// getVariantPricesInChannelsDict
func (a *AppProduct) getVariantPricesInChannelsDict(product *product_and_discount.Product) (map[string][]*goprices.Money, *model.AppError) {
	variantChannelListings, appErr := a.
		ProductVariantChannelListingsByOption(&product_and_discount.ProductVariantChannelListingFilterOption{
			VariantProductID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: product.Id,
				},
			},
			PriceAmount: &model.NumberFilter{
				NumberOption: &model.NumberOption{
					NULL: model.NewBool(false),
				},
			},
		})
	if appErr != nil {
		return nil, appErr
	}

	pricesDict := map[string][]*goprices.Money{}
	for _, listing := range variantChannelListings {
		listing.PopulateNonDbFields() // must run this first
		pricesDict[listing.ChannelID] = append(pricesDict[listing.ChannelID], listing.Price)
	}

	return pricesDict, nil
}
