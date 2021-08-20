package product

import (
	"net/http"
	"time"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
)

// CollectionsByOption returns all collections that satisfy given option.
//
// NOTE: `ShopID` is required.
func (a *AppProduct) CollectionsByOption(option *product_and_discount.CollectionFilterOption) ([]*product_and_discount.Collection, *model.AppError) {
	// validate if shopID is provided
	if !model.IsValidId(option.ShopID) {
		return nil, model.NewAppError("CollectionsByOption", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "option.ShopID"}, "", http.StatusBadRequest)
	}

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
func (a *AppProduct) PublishedCollections(channelSlug string, shopID string) ([]*product_and_discount.Collection, *model.AppError) {
	today := util.StartOfDay(time.Now().UTC())

	return a.CollectionsByOption(&product_and_discount.CollectionFilterOption{
		ShopID: shopID,
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

// VisibleCollectionsToUser returns all collections that belong to given shop and can be viewed by given user
func (a *AppProduct) VisibleCollectionsToUser(userID string, shopID string, channelSlug string) ([]*product_and_discount.Collection, *model.AppError) {
	// check if shop and user has relationship (shop-staff)
	_, appErr := a.ShopApp().ShopStaffRelationByShopIDAndStaffID(shopID, userID)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // return immediately if error is caused by system
		}
		return a.PublishedCollections(channelSlug, shopID) // not found error, returns only published collections
	}

	if channelSlug != "" {
		return a.CollectionsByOption(&product_and_discount.CollectionFilterOption{
			ShopID: shopID,
			ChannelListingChannelSlug: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: channelSlug,
				},
			},
		})
	}

	return a.CollectionsByOption(&product_and_discount.CollectionFilterOption{
		ShopID:    shopID,
		SelectAll: true,
	})
}
