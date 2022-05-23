package wishlist

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type ServiceWishlist struct {
	srv *app.Server
}

func init() {
	app.RegisterWishlistService(func(s *app.Server) (sub_app_iface.WishlistService, error) {
		return &ServiceWishlist{
			srv: s,
		}, nil
	})
}

// UpsertWishlist inserts a new wishlist instance into database with given userID
func (a *ServiceWishlist) UpsertWishlist(wishList *wishlist.Wishlist) (*wishlist.Wishlist, *model.AppError) {
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
func (a *ServiceWishlist) WishlistByOption(option *wishlist.WishlistFilterOption) (*wishlist.Wishlist, *model.AppError) {
	wl, err := a.srv.Store.Wishlist().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WishlistByOption", "app.wishlist.error_finding_wishlist.app_error", err)
	}

	return wl, nil
}

// SetUser assigns given user to given wishlist
func (a *ServiceWishlist) SetUserForWishlist(wishList *wishlist.Wishlist, userID string) *model.AppError {
	// validate given user is valid
	if !model.IsValidId(userID) || wishList.UserID == &userID {
		return model.NewAppError("SetUserForWishlist", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "userID"}, "", http.StatusBadRequest)
	}
	wishList.UserID = &userID

	_, appErr := a.UpsertWishlist(wishList)
	return appErr
}

// GetAllVariants returns all product variants in child wishlist items of given wishlist
func (a *ServiceWishlist) GetAllVariants(wishlistID string) ([]*product_and_discount.ProductVariant, *model.AppError) {
	productVariants, appErr := a.srv.ProductService().ProductVariantsByOption(&product_and_discount.ProductVariantFilterOption{
		WishlistID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: wishlistID,
			},
		},
		Distinct: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	return productVariants, nil
}

// AddProduct add or create a wishlist item that belongs to given wishlist and contains given product
func (a *ServiceWishlist) AddProduct(wishlistID string, productID string) (*wishlist.WishlistItem, *model.AppError) {
	item, appErr := a.GetOrCreateWishlistItem(&wishlist.WishlistItem{
		WishlistID: wishlistID,
		ProductID:  productID,
	})

	return item, appErr
}

// RemoveProduct removes a wishlist item of given wishlist that have ProductID property is given productID
func (a *ServiceWishlist) RemoveProduct(wishlistID string, productID string) *model.AppError {
	_, appErr := a.DeleteWishlistItemsByOption(nil, &wishlist.WishlistItemFilterOption{
		WishlistID: squirrel.Eq{store.WishlistItemTableName + ".WishlistID": wishlistID},
		ProductID:  squirrel.Eq{store.WishlistItemTableName + ".ProductID": productID},
	})

	return appErr
}

// AddProductVariant add given product variant into given wishlist
func (a *ServiceWishlist) AddProductVariant(wishlistID string, productVariant *product_and_discount.ProductVariant) (*wishlist.WishlistItem, *model.AppError) {
	item, appErr := a.AddProduct(wishlistID, productVariant.ProductID)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = a.AddWishlistItemProductVariantRelation(&wishlist.WishlistItemProductVariant{
		WishlistItemID:   item.Id,
		ProductVariantID: productVariant.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	return item, nil
}

// RemoveProductVariant remove a wishlist item from given wishlist
func (a *ServiceWishlist) RemoveProductVariant(wishlistID string, productVariant *product_and_discount.ProductVariant) *model.AppError {
	wishlistItem, appErr := a.WishlistItemByOption(&wishlist.WishlistItemFilterOption{
		WishlistID: squirrel.Eq{store.WishlistItemTableName + ".WishlistID": wishlistID},
		ProductID:  squirrel.Eq{store.WishlistItemTableName + ".ProductID": productVariant.ProductID},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		return nil
	}

	numOfRelationsLeft, appErr := a.DeleteWishlistItemProductVariantRelation(&wishlist.WishlistItemProductVariant{
		ProductVariantID: productVariant.Id,
		WishlistItemID:   wishlistItem.Id,
	})
	if appErr != nil {
		return appErr
	}

	if numOfRelationsLeft == 0 {
		_, appErr = a.DeleteWishlistItemsByOption(nil, &wishlist.WishlistItemFilterOption{
			Id: squirrel.Eq{store.WishlistItemTableName + ".Id": wishlistItem.Id},
		})
	}

	return appErr
}
