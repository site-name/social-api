package wishlist

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// WishlistItemsByOption returns a slice of wishlist items filtered using given option
func (a *ServiceWishlist) WishlistItemsByOption(option model_helper.WishlistItemFilterOption) (model.WishlistItemSlice, *model_helper.AppError) {
	items, err := a.srv.Store.WishlistItem().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("WishlistItemsByOption", "app.model.error_finding_wishlist_items_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return items, nil
}

// WishlistItemByOption returns 1 wishlist item filtered using given option
func (a *ServiceWishlist) WishlistItemByOption(option model_helper.WishlistItemFilterOption) (*model.WishlistItem, *model_helper.AppError) {
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
func (a *ServiceWishlist) BulkUpsertWishlistItems(transaction boil.ContextTransactor, wishlistItems model.WishlistItemSlice) (model.WishlistItemSlice, *model_helper.AppError) {
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
func (a *ServiceWishlist) GetOrCreateWishlistItem(wishlistItem model.WishlistItem) (*model.WishlistItem, *model_helper.AppError) {
	conditions := make([]qm.QueryMod, 0, 3)

	if wishlistItem.ID != "" {
		conditions = append(conditions, model.WishlistItemWhere.ID.EQ(wishlistItem.ID))
	}
	if wishlistItem.WishlistID != "" {
		conditions = append(conditions, model.WishlistItemWhere.WishlistID.EQ(wishlistItem.WishlistID))
	}
	if wishlistItem.ProductID != "" {
		conditions = append(conditions, model.WishlistItemWhere.ProductID.EQ(wishlistItem.ProductID))
	}

	_, appErr := a.WishlistItemByOption(model_helper.WishlistItemFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(conditions...),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// this means wishlist item not found, we need to create a new one
		items, appErr := a.BulkUpsertWishlistItems(nil, model.WishlistItemSlice{&wishlistItem})
		if appErr != nil {
			return nil, appErr
		}
		wishlistItem = *items[0]
	}

	return &wishlistItem, appErr
}

// DeleteWishlistItemsByOption tell store to delete wishlist items that satisfy given option, then returns a number of items deleted
func (a *ServiceWishlist) DeleteWishlistItemsByOption(transaction boil.ContextTransactor, option *model.WishlistItemFilterOption) (int64, *model_helper.AppError) {
	numDeleted, err := a.srv.Store.WishlistItem().Delete(transaction, option)
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

	var wishlistItems model.WishlistItemSlice
	err := transaction.
		Preload("ProductVariants").
		Find(&wishlistItems, "WishlistID IN ?", []string{srcWishlist.ID, dstWishlist.ID}).Error
	if err != nil {
		return model_helper.NewAppError("MoveItemsBetweenWishlists", "app.wishlist.wishlist_items_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	destWishlistMap := map[string]*model.WishlistItem{} // keys are product ids

	for _, item := range wishlistItems {
		if item.WishlistID == dstWishlist.ID {
			destWishlistMap[item.ProductID] = item
		}
	}

	for _, item := range wishlistItems {
		if item.WishlistID == srcWishlist.ID {
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
