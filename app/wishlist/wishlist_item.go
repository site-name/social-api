package wishlist

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

// WishlistItemsByOption returns a slice of wishlist items filtered using given option
func (a *AppWishlist) WishlistItemsByOption(option *wishlist.WishlistItemFilterOption) ([]*wishlist.WishlistItem, *model.AppError) {
	items, err := a.Srv().Store.WishlistItem().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WishlistItemsByOption", "app.wishlist.error_finding_wishlist_items_by_option.app_error", err)
	}
	return items, nil
}

// WishlistItemByOption returns 1 wishlist item filtered using given option
func (a *AppWishlist) WishlistItemByOption(option *wishlist.WishlistItemFilterOption) (*wishlist.WishlistItem, *model.AppError) {
	item, err := a.Srv().Store.WishlistItem().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WishlistItemByOption", "app.wishlist.error_finding_wishlist_item_by_option.app_error", err)
	}

	return item, nil
}

// UpsertWishlistItem updates or inserts given wishlist item into database then returns it
func (a *AppWishlist) UpsertWishlistItem(wishlistItem *wishlist.WishlistItem) (*wishlist.WishlistItem, *model.AppError) {
	wishlistItem, err := a.Srv().Store.WishlistItem().Upsert(wishlistItem)
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

		return nil, model.NewAppError("UpsertWishlistItem", "app.wishlist.error_upserting_wishlist_item.app_error", nil, err.Error(), statusCode)
	}

	return wishlistItem, nil
}

// GetOrCreateWishlistItem insert or get wishlist items
func (a *AppWishlist) GetOrCreateWishlistItem(wishlistItem *wishlist.WishlistItem) (*wishlist.WishlistItem, *model.AppError) {
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
		item, appErr = a.UpsertWishlistItem(wishlistItem)
	}

	return item, appErr
}

// DeleteWishlistItemsByOption tell store to delete wishlist items that satisfy given option, then returns a number of items deleted
func (a *AppWishlist) DeleteWishlistItemsByOption(option *wishlist.WishlistItemFilterOption) (int64, *model.AppError) {
	numDeleted, err := a.Srv().Store.WishlistItem().DeleteItemsByOption(option)
	if err != nil {
		return 0, model.NewAppError("DeleteWishlistItemsByOption", "app.wishlist.error_deleting_wishlist_items_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numDeleted, nil
}

// MoveItemsBetweenWishlists moves items from given srcWishlist to given dstWishlist
func (a *AppWishlist) MoveItemsBetweenWishlists(srcWishlist *wishlist.Wishlist, dstWishlist *wishlist.Wishlist) *model.AppError {
	transaction, err := a.Srv().Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("MoveItemsBetweenWishlists", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.Srv().Store.FinalizeTransaction(transaction)

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

	// categorize which items belong to which list
	var (
		itemsOfDestWishlist   []*wishlist.WishlistItem
		itemsOfSourceWishlist []*wishlist.WishlistItem
	)

	for _, item := range itemsFromBothWishlists {
		if item != nil && item.WishlistID == srcWishlist.Id {
			itemsOfSourceWishlist = append(itemsOfSourceWishlist, item)
			continue
		}
		itemsOfDestWishlist = append(itemsOfDestWishlist, item)
	}

	dstWishlistMap := map[string]*wishlist.WishlistItem{}
	for _, item := range itemsOfDestWishlist {
		dstWishlistMap[item.ProductID] = item
	}

	// Copying the items from the source to the destination wishlist.
	for _, srcItem := range itemsOfSourceWishlist {
		if anItem, exist := dstWishlistMap[srcItem.ProductID]; exist && anItem != nil {
			// This wishlist srcItem's product already exist.
			// Adding and the variants, "add" already handles duplicates.

			_, appErr = a.DeleteWishlistItemsByOption(&wishlist.WishlistItemFilterOption{
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
			a.UpsertWishlistItem(srcItem)
		}
	}

	// a.Upsert

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("MoveItemsBetweenWishlists", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
