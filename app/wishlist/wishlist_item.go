package wishlist

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// WishlistItemsByOption returns a slice of wishlist items filtered using given option
func (a *ServiceWishlist) WishlistItemsByOption(option *model.WishlistItemFilterOption) ([]*model.WishlistItem, *model_helper.AppError) {
	items, err := a.srv.Store.WishlistItem().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("WishlistItemsByOption", "app.model.error_finding_wishlist_items_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return items, nil
}

// WishlistItemByOption returns 1 wishlist item filtered using given option
func (a *ServiceWishlist) WishlistItemByOption(option *model.WishlistItemFilterOption) (*model.WishlistItem, *model_helper.AppError) {
	item, err := a.srv.Store.WishlistItem().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("WishlistItemByOption", "app.model.error_finding_wishlist_item_by_option.app_error", nil, err.Error(), statusCode)
	}

	return item, nil
}

// BulkUpsertWishlistItems updates or inserts given wishlist item into database then returns it
func (a *ServiceWishlist) BulkUpsertWishlistItems(transaction boil.ContextTransactor, wishlistItems model.WishlistItems) (model.WishlistItems, *model_helper.AppError) {
	wishlistItems, err := a.srv.Store.WishlistItem().BulkUpsert(transaction, wishlistItems)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		} else if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("BulkUpsertWishlistItems", "app.model.error_upserting_wishlist_item.app_error", nil, err.Error(), statusCode)
	}

	return wishlistItems, nil
}

// GetOrCreateWishlistItem insert or get wishlist items
func (a *ServiceWishlist) GetOrCreateWishlistItem(wishlistItem *model.WishlistItem) (*model.WishlistItem, *model_helper.AppError) {
	conditions := squirrel.Eq{}

	if wishlistItem.Id != "" {
		conditions[model.WishlistItemTableName+".Id"] = wishlistItem.Id
	}
	if wishlistItem.WishlistID != "" {
		conditions[model.WishlistItemTableName+".WishlistID"] = wishlistItem.WishlistID
	}
	if wishlistItem.ProductID != "" {
		conditions[model.WishlistItemTableName+".ProductID"] = wishlistItem.ProductID
	}

	wishistItem, appErr := a.WishlistItemByOption(&model.WishlistItemFilterOption{
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
		wishistItem = items[0]
	}

	return wishistItem, appErr
}

// DeleteWishlistItemsByOption tell store to delete wishlist items that satisfy given option, then returns a number of items deleted
func (a *ServiceWishlist) DeleteWishlistItemsByOption(transaction boil.ContextTransactor, option *model.WishlistItemFilterOption) (int64, *model_helper.AppError) {
	numDeleted, err := a.srv.Store.WishlistItem().DeleteItemsByOption(transaction, option)
	if err != nil {
		return 0, model_helper.NewAppError("DeleteWishlistItemsByOption", "app.wishlist.error_deleting_wishlist_items_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numDeleted, nil
}

// MoveItemsBetweenWishlists moves items from given srcWishlist to given dstWishlist
func (a *ServiceWishlist) MoveItemsBetweenWishlists(srcWishlist *model.Wishlist, dstWishlist *model.Wishlist) *model_helper.AppError {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return model_helper.NewAppError("MoveItemsBetweenWishlists", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var wishlistItems model.WishlistItems
	err := transaction.
		Preload("ProductVariants").
		Find(&wishlistItems, "WishlistID IN ?", []string{srcWishlist.Id, dstWishlist.Id}).Error
	if err != nil {
		return model_helper.NewAppError("MoveItemsBetweenWishlists", "app.wishlist.wishlist_items_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	destWishlistMap := map[string]*model.WishlistItem{} // keys are product ids

	for _, item := range wishlistItems {
		if item.WishlistID == dstWishlist.Id {
			destWishlistMap[item.ProductID] = item
		}
	}

	for _, item := range wishlistItems {
		if item.WishlistID == srcWishlist.Id {
			if destItem, ok := destWishlistMap[item.ProductID]; ok {
				err := transaction.Model(destItem).Association("ProductVariants").Append(item.ProductVariants)
				if err != nil {
					return model_helper.NewAppError("MoveItemsBetweenWishlists", "app.wishlist.add_variants_to_wishlist_item.app_error", nil, err.Error(), http.StatusInternalServerError)
				}
				_, appErr := a.DeleteWishlistItemsByOption(transaction, &model.WishlistItemFilterOption{
					Conditions: squirrel.Eq{model.WishlistItemTableName + ".Id": item.Id},
				})
				if appErr != nil {
					return appErr
				}
			} else {
				item.WishlistID = dstWishlist.Id
				_, appErr := a.BulkUpsertWishlistItems(transaction, model.WishlistItems{item})
				if appErr != nil {
					return appErr
				}
			}
		}
	}

	if err := transaction.Commit().Error; err != nil {
		return model_helper.NewAppError("MoveItemsBetweenWishlists", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
