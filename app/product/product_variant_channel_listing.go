package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// ProductVariantChannelListingsByOption returns a slice of product variant channel listings by given option
func (a *AppProduct) ProductVariantChannelListingsByOption(option *product_and_discount.ProductVariantChannelListingFilterOption) ([]*product_and_discount.ProductVariantChannelListing, *model.AppError) {
	listings, err := a.Srv().Store.ProductVariantChannelListing().FilterbyOption(option)
	var (
		statusCode   int
		errorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	}
	if len(listings) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ProductVariantChannelListingsByOption", "app.product_error_finding_product_variant_channel_listings_by_option.app_error", nil, errorMessage, statusCode)
	}

	return listings, nil
}
