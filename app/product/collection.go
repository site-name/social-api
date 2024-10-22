package product

import (
	"net/http"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

// CollectionsByOption returns all collections that satisfy given option.
//
// NOTE: `ShopID` is required.
func (a *ServiceProduct) CollectionsByOption(option *model.CollectionFilterOption) (int64, model.CollectionSlice, *model_helper.AppError) {
	totalCount, collections, err := a.srv.Store.Collection().FilterByOption(option)
	if err != nil {
		return 0, nil, model_helper.NewAppError("CollectionsByOption", "app.product.error_finding_collections_by_option", nil, err.Error(), http.StatusInternalServerError)
	}

	return totalCount, collections, nil
}

// CollectionsByVoucherID finds all collections that have relationships with given voucher
func (a *ServiceProduct) CollectionsByVoucherID(voucherID string) ([]*model.Collection, *model_helper.AppError) {
	_, collections, appErr := a.CollectionsByOption(&model.CollectionFilterOption{
		VoucherID: squirrel.Eq{model.VoucherCollectionTableName + ".VoucherID": voucherID},
	})
	return collections, appErr
}

// CollectionsByProductID finds and returns all collections related to given product
func (a *ServiceProduct) CollectionsByProductID(productID string) ([]*model.Collection, *model_helper.AppError) {
	_, collections, appEr := a.CollectionsByOption(&model.CollectionFilterOption{
		ProductID: squirrel.Eq{model.CollectionProductRelationTableName + ".ProductID": productID},
	})
	return collections, appEr
}

// PublishedCollections returns all published collections
func (a *ServiceProduct) PublishedCollections(channelSlug string) ([]*model.Collection, *model_helper.AppError) {
	today := util.StartOfDay(time.Now())

	publishedCollectionFilterOpts := &model.CollectionFilterOption{
		ChannelListingPublicationDate: squirrel.Expr(model.CollectionChannelListingTableName+".PublicationDate <= ? OR CollectionChannelListings.PublicationDate IS NULL", today),
		ChannelListingIsPublished:     squirrel.Expr(model.CollectionChannelListingTableName + ".IsPublished"),
		ChannelListingChannelIsActive: squirrel.Expr(model.ChannelTableName + ".IsActive"),
	}
	if channelSlug != "" {
		publishedCollectionFilterOpts.ChannelListingChannelSlug = squirrel.Expr(model.ChannelTableName+".Slug = ?", channelSlug)
	}

	_, collections, appErr := a.CollectionsByOption(publishedCollectionFilterOpts)
	return collections, appErr
}

func (a *ServiceProduct) VisibleCollectionsToUser(channelSlug string, userIsShopStaff bool) ([]*model.Collection, *model_helper.AppError) {
	if userIsShopStaff {
		collectionFilterOpts := &model.CollectionFilterOption{}
		if channelSlug != "" {
			collectionFilterOpts.ChannelListingChannelSlug = squirrel.Expr(model.ChannelTableName+".Slug = ?", channelSlug)
		}
		_, collections, appErr := a.CollectionsByOption(collectionFilterOpts)
		return collections, appErr
	}

	return a.PublishedCollections(channelSlug)
}

func (a *ServiceProduct) CollectionChannelListingsByOptions(options *model.CollectionChannelListingFilterOptions) ([]*model.CollectionChannelListing, *model_helper.AppError) {
	rels, err := a.srv.Store.CollectionChannelListing().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("CollectionChannelListingsByOptions", "app.product.collection_channel_listings_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return rels, nil
}
