package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shop"
)

func (a *AppProduct) GetDefaultDigitalContentSettings(aShop *shop.Shop) *shop.ShopDefaultDigitalContentSettings {
	return &shop.ShopDefaultDigitalContentSettings{
		AutomaticFulfillmentDigitalProducts: aShop.AutomaticFulfillmentDigitalProducts,
		DefaultDigitalMaxDownloads:          aShop.DefaultDigitalMaxDownloads,
		DefaultDigitalUrlValidDays:          aShop.DefaultDigitalUrlValidDays,
	}
}

// DigitalContentUrlIsValid Check if digital url is still valid for customer.
//
// It takes default settings or digital product's settings
// to check if url is still valid.
func (a *AppProduct) DigitalContentUrlIsValid(contentURL *product_and_discount.DigitalContentUrl) (bool, *model.AppError) {
	panic("not implemented")
}
