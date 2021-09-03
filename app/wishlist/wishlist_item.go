package wishlist

import (
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

// WishlistItemsByOption returns a slice of wishlist items filtered using given option
func (a *ServiceWishlist) WishlistItemsByOption(option *wishlist.WishlistItemFilterOption) ([]*wishlist.WishlistItem, *model.AppError) {
	items, err := a.srv.Store.WishlistItem().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WishlistItemsByOption", "app.wishlist.error_finding_wishlist_items_by_option.app_error", err)
	}
	return items, nil
}

// WishlistItemByOption returns 1 wishlist item filtered using given option
func (a *ServiceWishlist) WishlistItemByOption(option *wishlist.WishlistItemFilterOption) (*wishlist.WishlistItem, *model.AppError) {
	item, err := a.srv.Store.WishlistItem().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WishlistItemByOption", "app.wishlist.error_finding_wishlist_item_by_option.app_error", err)
	}

	return item, nil
}

// BulkUpsertWishlistItems updates or inserts given wishlist item into database then returns it
func (a *ServiceWishlist) BulkUpsertWishlistItems(transaction *gorp.Transaction, wishlistItems wishlist.WishlistItems) (wishlist.WishlistItems, *model.AppError) {
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

		return nil, model.NewAppError("BulkUpsertWishlistItems", "app.wishlist.error_upserting_wishlist_item.app_error", nil, err.Error(), statusCode)
	}

	return wishlistItems, nil
}

// GetOrCreateWishlistItem insert or get wishlist items
func (a *ServiceWishlist) GetOrCreateWishlistItem(wishlistItem *wishlist.WishlistItem) (*wishlist.WishlistItem, *model.AppError) {
	option := &wishlist.WishlistItemFilterOption{}

	if model.IsValidId(wishlistItem.Id) {
		option.Id = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: wishlistItem.Id,
			},
		}
	}
	if model.IsValidId(wishlistItem.WishlistID) {
		option.WishlistID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: wishlistItem.WishlistID,
			},
		}
	}
	if model.IsValidId(wishlistItem.ProductID) {
		option.ProductID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: wishlistItem.ProductID,
			},
		}
	}

	item, appErr := a.WishlistItemByOption(option)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// this means wishlist item not found, we need to create a new one
		items, appErr := a.BulkUpsertWishlistItems(nil, wishlist.WishlistItems{wishlistItem})
		if appErr != nil {
			return nil, appErr
		}
		item = items[0]
	}

	return item, appErr
}

// DeleteWishlistItemsByOption tell store to delete wishlist items that satisfy given option, then returns a number of items deleted
func (a *ServiceWishlist) DeleteWishlistItemsByOption(transaction *gorp.Transaction, option *wishlist.WishlistItemFilterOption) (int64, *model.AppError) {
	numDeleted, err := a.srv.Store.WishlistItem().DeleteItemsByOption(transaction, option)
	if err != nil {
		return 0, model.NewAppError("DeleteWishlistItemsByOption", "app.wishlist.error_deleting_wishlist_items_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numDeleted, nil
}

// MoveItemsBetweenWishlists moves items from given srcWishlist to given dstWishlist
func (a *ServiceWishlist) MoveItemsBetweenWishlists(srcWishlist *wishlist.Wishlist, dstWishlist *wishlist.Wishlist) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("MoveItemsBetweenWishlists", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	itemsFromBothWishlists, appErr := a.WishlistItemsByOption(&wishlist.WishlistItemFilterOption{
		WishlistID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: []string{srcWishlist.Id, dstWishlist.Id},
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		return nil
	}

	// categorize which item belongs to which wishlist
	var (
		itemsOfSourceWishlist              wishlist.WishlistItems
		destinationWishlistMap             = map[string]*wishlist.WishlistItem{}                 // destinationWishlistMap is a map with keys are product ids of destination wishlist's items
		productVariantsOfSourceWishlistMap = map[string][]*wishlist.WishlistItemProductVariant{} // productVariantsOfSourceWishlistMap has keys are source wishlist items's ids
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
		productVariantsOfSourceWishlist, appErr := a.srv.ProductService().ProductVariantsByOption(&product_and_discount.ProductVariantFilterOption{
			WishlistItemID: &model.StringFilter{
				StringOption: &model.StringOption{
					In: itemsOfSourceWishlist.IDs(),
				},
			},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return appErr
			}
			productVariantsOfSourceWishlist = []*product_and_discount.ProductVariant{}
		}

		for _, item := range itemsOfSourceWishlist {
			productVariantsOfSourceWishlistMap[item.Id] = []*wishlist.WishlistItemProductVariant{}

			for index, variant := range productVariantsOfSourceWishlist {
				if variant.ProductID == item.ProductID {

					productVariantsOfSourceWishlistMap[item.Id] = append(
						productVariantsOfSourceWishlistMap[item.Id],

						&wishlist.WishlistItemProductVariant{
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

			_, appErr = a.DeleteWishlistItemsByOption(transaction, &wishlist.WishlistItemFilterOption{
				Id: &model.StringFilter{
					StringOption: &model.StringOption{
						Eq: srcItem.Id,
					},
				},
			})
			if appErr != nil {
				return appErr
			}
		} else {
			// This wishlist srcItem contains a new product.
			// It can be reassigned to the destination wishlist.
			srcItem.WishlistID = dstWishlist.Id
			_, appErr = a.BulkUpsertWishlistItems(transaction, wishlist.WishlistItems{srcItem})
			if appErr != nil {
				return appErr
			}
		}
	}

	for _, item := range itemsOfSourceWishlist {
		item.WishlistID = dstWishlist.Id
	}

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("MoveItemsBetweenWishlists", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
