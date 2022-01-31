package product

import (
	"net/http"

	"github.com/mattermost/gorp"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
)

// ProductVariantById finds product variant by given id
func (a *ServiceProduct) ProductVariantById(id string) (*product_and_discount.ProductVariant, *model.AppError) {
	variant, err := a.srv.Store.ProductVariant().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantbyId", "app.product.product_variant_missing.app_error", err)
	}

	return variant, nil
}

// ProductVariantGetPrice returns price
func (a *ServiceProduct) ProductVariantGetPrice(
	productVariant *product_and_discount.ProductVariant,
	product product_and_discount.Product,
	collections []*product_and_discount.Collection,
	channel channel.Channel,
	channelListing *product_and_discount.ProductVariantChannelListing,
	discounts []*product_and_discount.DiscountInfo, // optional
) (*goprices.Money, *model.AppError) {
	return a.srv.DiscountService().CalculateDiscountedPrice(product, channelListing.Price, collections, discounts, channel, productVariant.Id)
}

// ProductVariantIsDigital finds product type that related to given product variant and check if that product type is digital and does not require shipping
func (a *ServiceProduct) ProductVariantIsDigital(productVariantID string) (bool, *model.AppError) {
	productType, err := a.srv.Store.ProductType().ProductTypeByProductVariantID(productVariantID)
	if err != nil {
		return false, store.AppErrorFromDatabaseLookupError("ProductVariantIsDigital", "app.product.product_type_by_product_variant_id.app_error", err)
	}

	return *productType.IsDigital && !*productType.IsShippingRequired, nil
}

// ProductVariantByOrderLineID returns a product variant by given order line id
func (a *ServiceProduct) ProductVariantByOrderLineID(orderLineID string) (*product_and_discount.ProductVariant, *model.AppError) {
	productVariant, err := a.srv.Store.ProductVariant().GetByOrderLineID(orderLineID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantByOrderLineID", "app.product.error_finding_product_variant_by_order_line_id.app_error", err)
	}

	return productVariant, nil
}

// ProductVariantsByOption returns a list of product variants satisfy given option
func (a *ServiceProduct) ProductVariantsByOption(option *product_and_discount.ProductVariantFilterOption) ([]*product_and_discount.ProductVariant, *model.AppError) {
	productVariants, err := a.srv.Store.ProductVariant().FilterByOption(option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(productVariants) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ProductVariantsByOption", "app.product.error_finding_product_variants_by_options.app_error", nil, errMessage, statusCode)
	}

	return productVariants, nil
}

// ProductVariantGetWeight returns weight of given product variant
func (a *ServiceProduct) ProductVariantGetWeight(productVariantID string) (*measurement.Weight, *model.AppError) {
	weight, err := a.srv.Store.ProductVariant().GetWeight(productVariantID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantGetWeight", "app.product.error_getting_product_variant_weight.app_error", err)
	}

	return weight, nil
}

// DisplayProduct return display text for given product variant
//
// `translated` default to false
func (a *ServiceProduct) DisplayProduct(productVariant *product_and_discount.ProductVariant, translated bool) (stringm *model.AppError) {
	panic("not implt")
}

// ProductVariantsAvailableInChannel returns product variants based on given channel slug
func (a *ServiceProduct) ProductVariantsAvailableInChannel(channelSlug string) ([]*product_and_discount.ProductVariant, *model.AppError) {
	productVariants, appErr := a.ProductVariantsByOption(&product_and_discount.ProductVariantFilterOption{
		ProductVariantChannelListingPriceAmount: &model.NumberFilter{
			NumberOption: &model.NumberOption{
				NULL: model.NewBool(false),
			},
		},
		ProductVariantChannelListingChannelSlug: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: channelSlug,
			},
		},
	})

	if appErr != nil {
		return nil, appErr
	}

	return productVariants, nil
}

// UpsertProductVariant tells store to upsert given product variant and returns it
func (s *ServiceProduct) UpsertProductVariant(transaction *gorp.Transaction, variant *product_and_discount.ProductVariant) (*product_and_discount.ProductVariant, *model.AppError) {
	var (
		upsertedVariant *product_and_discount.ProductVariant
		err             error
	)
	if !model.IsValidId(variant.Id) {
		upsertedVariant, err = s.srv.Store.ProductVariant().Save(transaction, variant)
	} else {
		upsertedVariant, err = s.srv.Store.ProductVariant().Update(transaction, variant)
	}
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		var (
			statusCode = http.StatusInternalServerError
		)
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("UpsertProductVariant", "app.product.error_upserting_product_variant.app_error", nil, err.Error(), statusCode)
	}

	return upsertedVariant, nil
}
