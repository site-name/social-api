package wishlist

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// WishlistItemsByOption returns a slice of wishlist items filtered using given option
func (a *ServiceWishlist) WishlistItemsByOption(option *model.WishlistItemFilterOption) ([]*model.WishlistItem, *model.AppError) {
	items, err := a.srv.Store.WishlistItem().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("WishlistItemsByOption", "app.model.error_finding_wishlist_items_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return items, nil
}

// WishlistItemByOption returns 1 wishlist item filtered using given option
func (a *ServiceWishlist) WishlistItemByOption(option *model.WishlistItemFilterOption) (*model.WishlistItem, *model.AppError) {
	item, err := a.srv.Store.WishlistItem().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("WishlistItemByOption", "app.model.error_finding_wishlist_item_by_option.app_error", nil, err.Error(), statusCode)
	}

	return item, nil
}

// BulkUpsertWishlistItems updates or inserts given wishlist item into database then returns it
func (a *ServiceWishlist) BulkUpsertWishlistItems(transaction *gorm.DB, wishlistItems model.WishlistItems) (model.WishlistItems, *model.AppError) {
	wishlistItems, err := a.srv.Store.WishlistItem().BulkUpsert(transaction, wishlistItems)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		} else if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("BulkUpsertWishlistItems", "app.model.error_upserting_wishlist_item.app_error", nil, err.Error(), statusCode)
	}

	return wishlistItems, nil
}

// GetOrCreateWishlistItem insert or get wishlist items
func (a *ServiceWishlist) GetOrCreateWishlistItem(wishlistItem *model.WishlistItem) (*model.WishlistItem, *model.AppError) {
	conditions := squirrel.And{}

	if model.IsValidId(wishlistItem.Id) {
		conditions = append(conditions, squirrel.Eq{model.WishlistItemTableName + ".Id": wishlistItem.Id})
	}
	if model.IsValidId(wishlistItem.WishlistID) {
		conditions = append(conditions, squirrel.Eq{model.WishlistItemTableName + ".WishlistID": wishlistItem.WishlistID})
	}
	if model.IsValidId(wishlistItem.ProductID) {
		conditions = append(conditions, squirrel.Eq{model.WishlistItemTableName + ".ProductID": wishlistItem.ProductID})
	}

	item, appErr := a.WishlistItemByOption(&model.WishlistItemFilterOption{
		Conditions: conditions,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// this means wishlist item not found, we need to create a new one
		items, appErr := a.BulkUpsertWishlistItems(nil, model.WishlistItems{wishlistItem})
		if appErr != nil {
			return nil, appErr
		}
		item = items[0]
	}

	return item, appErr
}

// DeleteWishlistItemsByOption tell store to delete wishlist items that satisfy given option, then returns a number of items deleted
func (a *ServiceWishlist) DeleteWishlistItemsByOption(transaction *gorm.DB, option *model.WishlistItemFilterOption) (int64, *model.AppError) {
	numDeleted, err := a.srv.Store.WishlistItem().DeleteItemsByOption(transaction, option)
	if err != nil {
		return 0, model.NewAppError("DeleteWishlistItemsByOption", "app.wishlist.error_deleting_wishlist_items_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numDeleted, nil
}

// MoveItemsBetweenWishlists moves items from given srcWishlist to given dstWishlist
func (a *ServiceWishlist) MoveItemsBetweenWishlists(srcWishlist *model.Wishlist, dstWishlist *model.Wishlist) *model.AppError {
	transaction := a.srv.Store.GetMaster().Begin()
	defer transaction.Rollback()

	itemsFromBothWishlists, appErr := a.WishlistItemsByOption(&model.WishlistItemFilterOption{
		Conditions: squirrel.Eq{model.WishlistItemTableName + ".WishlistID": []string{srcWishlist.Id, dstWishlist.Id}},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		return nil
	}

	// categorize which item belongs to which wishlist
	var (
		itemsOfSourceWishlist              model.WishlistItems
		destinationWishlistMap             = map[string]*model.WishlistItem{}                 // destinationWishlistMap is a map with keys are product ids of destination wishlist's items
		productVariantsOfSourceWishlistMap = map[string][]*model.WishlistItemProductVariant{} // productVariantsOfSourceWishlistMap has keys are source wishlist items's ids
	)

	for _, item := range itemsFromBothWishlists {
		if item != nil && item.WishlistID == srcWishlist.Id {
			itemsOfSourceWishlist = append(itemsOfSourceWishlist, item)
			continue
		}
		destinationWishlistMap[item.ProductID] = item
	}

	// this function will not execute if not triggered
	populate_productVariantsOfSourceWishlistMap := func() *model.AppError {
		productVariantsOfSourceWishlist, appErr := a.srv.ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
			WishlistItemID: squirrel.Eq{model.WishlistItemProductVariantTableName + ".WishlistItemID": itemsOfSourceWishlist.IDs()},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return appErr
			}
			productVariantsOfSourceWishlist = []*model.ProductVariant{}
		}

		for _, item := range itemsOfSourceWishlist {
			productVariantsOfSourceWishlistMap[item.Id] = []*model.WishlistItemProductVariant{}

			for index, variant := range productVariantsOfSourceWishlist {
				if variant.ProductID == item.ProductID {

					productVariantsOfSourceWishlistMap[item.Id] = append(
						productVariantsOfSourceWishlistMap[item.Id],

						&model.WishlistItemProductVariant{
							WishlistItemID:   item.Id,
							ProductVariantID: variant.Id,
						},
					)

					// this filters out matched product variants, make later loops faster
					productVariantsOfSourceWishlist = append(productVariantsOfSourceWishlist[:index], productVariantsOfSourceWishlist[index+1:]...)
				}
			}
		}

		return nil
	}

	// Copying the items from the source to the destination wishlist.
	for index, srcItem := range itemsOfSourceWishlist {
		if dstItem, exist := destinationWishlistMap[srcItem.ProductID]; exist && dstItem != nil {
			// This wishlist srcItem's product already exist.
			// Adding and the variants, "add" already handles duplicates.

			if index == 0 {
				appErr = populate_productVariantsOfSourceWishlistMap()
				if appErr != nil {
					return appErr
				}
			}

			_, appErr = a.BulkUpsertWishlistItemProductVariantRelations(transaction, productVariantsOfSourceWishlistMap[srcItem.Id])
			if appErr != nil {
				return appErr
			}

			_, appErr = a.DeleteWishlistItemsByOption(transaction, &model.WishlistItemFilterOption{
				Conditions: squirrel.Eq{model.WishlistItemTableName + ".Id": srcItem.Id},
			})
			if appErr != nil {
				return appErr
			}
		} else {
			// This wishlist srcItem contains a new product.
			// It can be reassigned to the destination wishlist.
			srcItem.WishlistID = dstWishlist.Id
			_, appErr = a.BulkUpsertWishlistItems(transaction, model.WishlistItems{srcItem})
			if appErr != nil {
				return appErr
			}
		}
	}

	for _, item := range itemsOfSourceWishlist {
		item.WishlistID = dstWishlist.Id
	}

	if err := transaction.Commit().Error; err != nil {
		return model.NewAppError("MoveItemsBetweenWishlists", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
