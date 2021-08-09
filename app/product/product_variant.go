package product

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// ProductVariantById finds product variant by given id
func (a *AppProduct) ProductVariantById(id string) (*product_and_discount.ProductVariant, *model.AppError) {
	variant, err := a.Srv().Store.ProductVariant().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantbyId", "app.product.product_variant_missing.app_error", err)
	}

	return variant, nil
}

func (a *AppProduct) ProductVariantGetPrice(
	product *product_and_discount.Product,
	collections []*product_and_discount.Collection,
	channel *channel.Channel,
	channelListing *product_and_discount.ProductVariantChannelListing,
	discounts []*product_and_discount.DiscountInfo, // optional
) (*goprices.Money, *model.AppError) {
	return a.DiscountApp().CalculateDiscountedPrice(product, channelListing.Price, collections, discounts, channel)
}

// ProductVariantIsDigital finds product type that related to given product variant and check if that product type is digital and does not require shipping
func (a *AppProduct) ProductVariantIsDigital(productVariantID string) (bool, *model.AppError) {
	productType, err := a.Srv().Store.ProductType().ProductTypeByProductVariantID(productVariantID)
	if err != nil {
		return false, store.AppErrorFromDatabaseLookupError("ProductVariantIsDigital", "app.product.product_type_by_product_variant_id.app_error", err)
	}

	return *productType.IsDigital && !*productType.IsShippingRequired, nil
}
