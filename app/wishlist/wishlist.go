package wishlist

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type ServiceWishlist struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Wishlist = &ServiceWishlist{s}
		return nil
	})
}

// UpsertWishlist inserts a new wishlist instance into database with given userID
func (a *ServiceWishlist) UpsertWishlist(wishList *model.Wishlist) (*model.Wishlist, *model.AppError) {
	newWl, err := a.srv.Store.Wishlist().Upsert(wishList)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertWishlist", "app.wishlist.error_upserting_wishlist.app_error", nil, err.Error(), statusCode)
	}

	return newWl, nil
}

// WishlistByOption returns 1 wishlist filtered by given option
func (a *ServiceWishlist) WishlistByOption(option *model.WishlistFilterOption) (*model.Wishlist, *model.AppError) {
	wishlist, err := a.srv.Store.Wishlist().GetByOption(option)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			status = http.StatusNotFound
		}
		return nil, model.NewAppError("WishlistByOption", "app.wishlist.error_finding_wishlist.app_error", nil, err.Error(), status)
	}

	return wishlist, nil
}

// SetUser assigns given user to given wishlist
func (a *ServiceWishlist) SetUserForWishlist(wishList *model.Wishlist, userID string) (*model.Wishlist, *model.AppError) {
	wishList.UserID = &userID

	return a.UpsertWishlist(wishList)
}

// GetAllVariants returns all product variants in child wishlist items of given wishlist
func (a *ServiceWishlist) GetAllVariants(wishlistID string) ([]*model.ProductVariant, *model.AppError) {
	productVariants, appErr := a.srv.ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		WishlistID: squirrel.Eq{model.WishlistItemTableName + ".WishlistID": wishlistID},
		Distinct:   true,
	})
	if appErr != nil {
		return nil, appErr
	}

	return productVariants, nil
}

// AddProduct add or create a wishlist item that belongs to given wishlist and contains given product
func (a *ServiceWishlist) AddProduct(wishlistID string, productID string) (*model.WishlistItem, *model.AppError) {
	item, appErr := a.GetOrCreateWishlistItem(&model.WishlistItem{
		WishlistID: wishlistID,
		ProductID:  productID,
	})

	return item, appErr
}

// RemoveProduct removes a wishlist item of given wishlist that have ProductID property is given productID
func (a *ServiceWishlist) RemoveProduct(wishlistID string, productID string) *model.AppError {
	_, appErr := a.DeleteWishlistItemsByOption(nil, &model.WishlistItemFilterOption{
		Conditions: squirrel.Eq{
			model.WishlistItemTableName + ".WishlistID": wishlistID,
			model.WishlistItemTableName + ".ProductID":  productID,
		},
	})

	return appErr
}

// AddProductVariant add given product variant into given wishlist
func (a *ServiceWishlist) AddProductVariant(wishlistID string, productVariant *model.ProductVariant) (*model.WishlistItem, *model.AppError) {
	item, appErr := a.AddProduct(wishlistID, productVariant.ProductID)
	if appErr != nil {
		return nil, appErr
	}

	err := a.srv.Store.GetMaster().Model(item).Association("ProductVariants").Append(productVariant)
	if err != nil {
		return nil, model.NewAppError("AddProductVariant", "app.wishlist.add_variant_to_wishlist_item.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return item, nil
}

// RemoveProductVariant remove a wishlist item from given wishlist
func (a *ServiceWishlist) RemoveProductVariant(wishlistID string, productVariant *model.ProductVariant) *model.AppError {
	wishlistItem, appErr := a.WishlistItemByOption(&model.WishlistItemFilterOption{
		Conditions: squirrel.Eq{
			model.WishlistItemTableName + ".WishlistID": wishlistID,
			model.WishlistItemTableName + ".ProductID":  productVariant.ProductID,
		},
	})
	if appErr != nil {
		return appErr
	}

	err := a.srv.Store.GetMaster().Model(wishlistItem).Association("ProductVariants").Delete(productVariant)
	if err != nil {
		return model.NewAppError("RemoveProductVariant", "app.wishlist.remove_product_variant_from_wishlist_item.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	numOfProductVariantsInWishlistItem := a.srv.Store.GetMaster().Model(wishlistItem).Association("ProductVariants").Count()

	if numOfProductVariantsInWishlistItem == 0 {
		_, appErr = a.DeleteWishlistItemsByOption(nil, &model.WishlistItemFilterOption{
			Conditions: squirrel.Eq{model.WishlistItemTableName + ".Id": wishlistItem.Id},
		})
	}

	return appErr
}
