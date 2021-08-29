package wishlist

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type AppWishlist struct {
	app.AppIface
}

func init() {
	app.RegisterWishlistApp(func(a app.AppIface) sub_app_iface.WishlistApp {
		return &AppWishlist{a}
	})
}

// UpsertWishlist inserts a new wishlist instance into database with given userID
func (a *AppWishlist) UpsertWishlist(wishList *wishlist.Wishlist) (*wishlist.Wishlist, *model.AppError) {
	newWl, err := a.Srv().Store.Wishlist().Upsert(wishList)
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
func (a *AppWishlist) WishlistByOption(option *wishlist.WishlistFilterOption) (*wishlist.Wishlist, *model.AppError) {
	wl, err := a.Srv().Store.Wishlist().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WishlistByOption", "app.wishlist.error_finding_wishlist.app_error", err)
	}

	return wl, nil
}

// SetUser assigns given user to given wishlist
func (a *AppWishlist) SetUserForWishlist(wishList *wishlist.Wishlist, userID string) *model.AppError {
	// validate given user is valid
	if !model.IsValidId(userID) || wishList.UserID == &userID {
		return model.NewAppError("SetUserForWishlist", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "userID"}, "", http.StatusBadRequest)
	}
	wishList.UserID = &userID

	_, appErr := a.UpsertWishlist(wishList)
	return appErr
}

// GetAllVariants returns all product variants in child wishlist items of given wishlist
func (a *AppWishlist) GetAllVariants(wishlistID string) ([]*product_and_discount.ProductVariant, *model.AppError) {
	productVariants, appErr := a.ProductApp().ProductVariantsByOption(&product_and_discount.ProductVariantFilterOption{
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
func (a *AppWishlist) AddProduct(wishlistID string, productID string) (*wishlist.WishlistItem, *model.AppError) {
	item, appErr := a.GetOrCreateWishlistItem(&wishlist.WishlistItem{
		WishlistID: wishlistID,
		ProductID:  productID,
	})

	return item, appErr
}

// RemoveProduct removes a wishlist item of given wishlist that have ProductID property is given productID
func (a *AppWishlist) RemoveProduct(wishlistID string, productID string) *model.AppError {
	_, appErr := a.DeleteWishlistItemsByOption(nil, &wishlist.WishlistItemFilterOption{
		WishlistID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: wishlistID,
			},
		},
		ProductID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: productID,
			},
		},
	})

	return appErr
}

// AddProductVariant add given product variant into given wishlist
func (a *AppWishlist) AddProductVariant(wishlistID string, productVariant *product_and_discount.ProductVariant) (*wishlist.WishlistItem, *model.AppError) {
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
func (a *AppWishlist) RemoveProductVariant(wishlistID string, productVariant *product_and_discount.ProductVariant) *model.AppError {
	wishlistItem, appErr := a.WishlistItemByOption(&wishlist.WishlistItemFilterOption{
		WishlistID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: wishlistID,
			},
		},
		ProductID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: productVariant.ProductID,
			},
		},
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
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: wishlistItem.Id,
				},
			},
		})
	}

	return appErr
}
