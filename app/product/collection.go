package product

import (
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// CollectionsByOption returns all collections that satisfy given option.
//
// NOTE: `ShopID` is required.
func (a *ServiceProduct) CollectionsByOption(option *model.CollectionFilterOption) (model.Collections, *model.AppError) {
	collections, err := a.srv.Store.Collection().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("CollectionsByOption", "app.product.error_finding_collections_by_option", nil, err.Error(), http.StatusInternalServerError)
	}

	return collections, nil
}

// CollectionsByVoucherID finds all collections that have relationships with given voucher
func (a *ServiceProduct) CollectionsByVoucherID(voucherID string) ([]*model.Collection, *model.AppError) {
	return a.CollectionsByOption(&model.CollectionFilterOption{
		VoucherID: squirrel.Eq{store.VoucherCollectionTableName + ".VoucherID": voucherID},
	})
}

// CollectionsByProductID finds and returns all collections related to given product
func (a *ServiceProduct) CollectionsByProductID(productID string) ([]*model.Collection, *model.AppError) {
	return a.CollectionsByOption(&model.CollectionFilterOption{
		ProductID: squirrel.Eq{store.CollectionProductRelationTableName + ".ProductID": productID},
	})
}

// PublishedCollections returns all published collections
func (a *ServiceProduct) PublishedCollections(channelSlug string, shopID string) ([]*model.Collection, *model.AppError) {
	today := util.StartOfDay(time.Now())

	return a.CollectionsByOption(&model.CollectionFilterOption{
		ShopID: shopID,
		ChannelListingPublicationDate: squirrel.Or{
			squirrel.LtOrEq{store.CollectionChannelListingTableName + ".PublicationDate": today},
			squirrel.Eq{store.CollectionChannelListingTableName + ".PublicationDate": nil},
		},
		ChannelListingIsPublished:     model.NewPrimitive(true),
		ChannelListingChannelSlug:     squirrel.Eq{store.ChannelTableName + ".Slug": channelSlug},
		ChannelListingChannelIsActive: model.NewPrimitive(true),
	})
}

// VisibleCollectionsToUser returns all collections that belong to given shop and can be viewed by given user
func (a *ServiceProduct) VisibleCollectionsToUser(userID, shopID, channelSlug string) ([]*model.Collection, *model.AppError) {
	// check if shop and user has relationship (shop-staff)
	_, appErr := a.srv.ShopService().ShopStaffByOptions(&model.ShopStaffRelationFilterOptions{
		ShopID:  squirrel.Eq{store.ShopStaffTableName + ".ShopID": shopID},
		StaffID: squirrel.Eq{store.ShopStaffTableName + ".StaffID": userID},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // return immediately if error is caused by system
		}
		return a.PublishedCollections(channelSlug, shopID) // not found error, returns only published collections
	}

	if channelSlug != "" {
		return a.CollectionsByOption(&model.CollectionFilterOption{
			ShopID:                    shopID,
			ChannelListingChannelSlug: squirrel.Eq{store.ChannelTableName + ".Slug": channelSlug},
		})
	}

	return a.CollectionsByOption(&model.CollectionFilterOption{
		ShopID:    shopID,
		SelectAll: true,
	})
}

func (a *ServiceProduct) CollectionChannelListingsByOptions(options *model.CollectionChannelListingFilterOptions) ([]*model.CollectionChannelListing, *model.AppError) {
	rels, err := a.srv.Store.CollectionChannelListing().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("CollectionChannelListingsByOptions", "app.product.collection_channel_listings_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return rels, nil
}
