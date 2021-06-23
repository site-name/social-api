package imports

import (
	_ "github.com/sitename/sitename/modules/jobs/active_users"
	// _ "github.com/sitename/sitename/modules/migrations"
	// register all sub applications
	_ "github.com/sitename/sitename/app/account"
	_ "github.com/sitename/sitename/app/attribute"
	_ "github.com/sitename/sitename/app/channel"
	_ "github.com/sitename/sitename/app/checkout"
	_ "github.com/sitename/sitename/app/csv"
	_ "github.com/sitename/sitename/app/discount"
	_ "github.com/sitename/sitename/app/giftcard"
	_ "github.com/sitename/sitename/app/invoice"
	_ "github.com/sitename/sitename/app/menu"
	_ "github.com/sitename/sitename/app/order"
	_ "github.com/sitename/sitename/app/page"
	_ "github.com/sitename/sitename/app/payment"
	_ "github.com/sitename/sitename/app/product"
	_ "github.com/sitename/sitename/app/seo"
	_ "github.com/sitename/sitename/app/shipping"
	_ "github.com/sitename/sitename/app/site"
	_ "github.com/sitename/sitename/app/warehouse"
	_ "github.com/sitename/sitename/app/webhook"
	_ "github.com/sitename/sitename/app/wishlist"
)
