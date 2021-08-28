package wishlist

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

// AddWishlistItemProductVariantRelation adds given wishlist item-product variant relation into database and returns it
func (a *AppWishlist) AddWishlistItemProductVariantRelation(relation *wishlist.WishlistItemProductVariant) (*wishlist.WishlistItemProductVariant, *model.AppError) {
	relation, err := a.Srv().Store.WishlistItemProductVariant().Save(relation)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("AddWishlistItemProductVariantRelation", "app.wishlist.error_adding_wishlist_item_product_variant_relation.app_error", nil, err.Error(), statusCode)
	}

	return relation, nil
}

// DeleteWishlistItemProductVariantRelation deletes a wishlist item-product variant relation and returns a number of remaining relations in database
func (a *AppWishlist) DeleteWishlistItemProductVariantRelation(relation *wishlist.WishlistItemProductVariant) (int64, *model.AppError) {
	numberOfRelationsLeft, err := a.Srv().Store.WishlistItemProductVariant().DeleteRelation(relation)
	if err != nil {
		return 0, model.NewAppError("DeleteWishlistItemProductVariantRelation", "app.wishlist.error_deleting_wishlist_item_product_variant_relation.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numberOfRelationsLeft, nil
}
