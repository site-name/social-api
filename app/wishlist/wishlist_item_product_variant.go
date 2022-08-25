package wishlist

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// BulkUpsertWishlistItemProductVariantRelations
func (a *ServiceWishlist) BulkUpsertWishlistItemProductVariantRelations(transaction store_iface.SqlxTxExecutor, relations []*wishlist.WishlistItemProductVariant) ([]*wishlist.WishlistItemProductVariant, *model.AppError) {
	relations, err := a.srv.Store.WishlistItemProductVariant().BulkUpsert(transaction, relations)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("BulkUpsertWishlistItemProductVariantRelations", "app.wishlist_error_upserting_wishlist_item_product_variant_relations.app_error", nil, err.Error(), statusCode)
	}

	return relations, nil
}

// AddWishlistItemProductVariantRelation adds given wishlist item-product variant relation into database and returns it
func (a *ServiceWishlist) AddWishlistItemProductVariantRelation(relation *wishlist.WishlistItemProductVariant) (*wishlist.WishlistItemProductVariant, *model.AppError) {
	relation, err := a.srv.Store.WishlistItemProductVariant().Save(relation)
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
func (a *ServiceWishlist) DeleteWishlistItemProductVariantRelation(relation *wishlist.WishlistItemProductVariant) (int64, *model.AppError) {
	numberOfRelationsLeft, err := a.srv.Store.WishlistItemProductVariant().DeleteRelation(relation)
	if err != nil {
		return 0, model.NewAppError("DeleteWishlistItemProductVariantRelation", "app.wishlist.error_deleting_wishlist_item_product_variant_relation.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numberOfRelationsLeft, nil
}
