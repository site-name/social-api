package wishlist

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppWishlist struct {
	app.AppIface
}

func init() {
	app.RegisterWishlistApp(func(a app.AppIface) sub_app_iface.WishlistApp {
		return &AppWishlist{a}
	})
}
