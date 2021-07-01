package wishlist

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
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

func (a *AppWishlist) CreateWishlist(userID string) (*wishlist.Wishlist, *model.AppError) {
	wl := &wishlist.Wishlist{
		UserID: &userID,
	}
	newWl, err := a.Srv().Store.Wishlist().Save(wl)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			// invalid properties
			return nil, appErr
		} else if invlErr, ok := err.(*store.ErrInvalidInput); ok {
			// user id duplicate error
			return nil, model.NewAppError("CreateWishlist", "app.wishlist.wishlist_duplicate.app_error", nil, invlErr.Error(), http.StatusBadRequest)
		} else {
			// system saving error
			return nil, model.NewAppError("CreateWishlist", "app.wishlist.wislist_saving_error.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return newWl, nil
}

func (a *AppWishlist) WishlistByUserID(userID string) (*wishlist.Wishlist, *model.AppError) {
	wl, err := a.Srv().Store.Wishlist().GetByUserID(userID)
	if err != nil {
		if _, ok := err.(*store.ErrNotFound); ok {
			wl, appErr := a.CreateWishlist(userID)
			if appErr != nil {
				return nil, appErr
			}
			return wl, nil
		}
		return nil, model.NewAppError("", "app.wishlist.wishlist_finding_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return wl, nil
}
