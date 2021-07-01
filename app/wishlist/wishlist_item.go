package wishlist

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

func (a *AppWishlist) WishlistItemsByWishlistID(wishlistID string) ([]*wishlist.WishlistItem, *model.AppError) {
	items, err := a.Srv().Store.WishlistItem().WishlistItemsByWishlistId(wishlistID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("", "app.wishlist.items_by_wishlist.app_error", err)
	}
	return items, nil
}
