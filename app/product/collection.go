package product

import (
	"net/http"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
)

// CollectionsByOption returns all collections that satisfy given option
func (a *AppProduct) CollectionsByOption(option *product_and_discount.CollectionFilterOption) ([]*product_and_discount.Collection, *model.AppError) {
	collections, err := a.Srv().Store.Collection().FilterByOption(option)
	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(collections) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("CollectionsByOption", "app.product.error_finding_collections_by_option", nil, errMsg, statusCode)
	}

	return collections, nil
}

// CollectionsByVoucherID finds all collections that have relationships with given voucher
func (a *AppProduct) CollectionsByVoucherID(voucherID string) ([]*product_and_discount.Collection, *model.AppError) {
	return a.CollectionsByOption(&product_and_discount.CollectionFilterOption{
		VoucherIDs: []string{voucherID},
	})
}

// CollectionsByProductID finds and returns all collections related to given product
func (a *AppProduct) CollectionsByProductID(productID string) ([]*product_and_discount.Collection, *model.AppError) {
	return a.CollectionsByOption(&product_and_discount.CollectionFilterOption{
		ProductIDs: []string{productID},
	})
}

// PublishedCollections returns all published collections
func (a *AppProduct) PublishedCollections(channelSlug string) ([]*product_and_discount.Collection, *model.AppError) {
	today := util.StartOfDay(time.Now().UTC())

	return a.CollectionsByOption(&product_and_discount.CollectionFilterOption{
		ChannelListingPublicationDate: &model.TimeFilter{
			Or: &model.TimeOption{
				LtE:               &today,
				CompareStartOfDay: true, // since `CollectionChannelListing's PublicationDate field has type date
				NULL:              model.NewBool(true),
			},
		},
		ChannelListingIsPublished: model.NewBool(true),
		ChannelListingChannelSlug: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: channelSlug,
			},
		},
		ChannelListingChannelIsActive: model.NewBool(true),
	})
}
