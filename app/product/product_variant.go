package product

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/measurement"
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

// ProductVariantGetPrice returns price
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

// ProductVariantByOrderLineID returns a product variant by given order line id
func (a *AppProduct) ProductVariantByOrderLineID(orderLineID string) (*product_and_discount.ProductVariant, *model.AppError) {
	productVariant, err := a.Srv().Store.ProductVariant().GetByOrderLineID(orderLineID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantByOrderLineID", "app.product.error_finding_product_variant_by_order_line_id.app_error", err)
	}

	return productVariant, nil
}

// ProductVariantsByOption returns a list of product variants satisfy given option
func (a *AppProduct) ProductVariantsByOption(option *product_and_discount.ProductVariantFilterOption) ([]*product_and_discount.ProductVariant, *model.AppError) {
	productVariants, err := a.Srv().Store.ProductVariant().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantsByOption", "app.product.error_finding_product_variants_by_option.app_error", err)
	}

	return productVariants, nil
}

// ProductVariantGetWeight returns weight of given product variant
func (a *AppProduct) ProductVariantGetWeight(productVariantID string) (*measurement.Weight, *model.AppError) {
	weight, err := a.Srv().Store.ProductVariant().GetWeight(productVariantID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantGetWeight", "app.product.error_getting_product_variant_weight.app_error", err)
	}

	return weight, nil
}
