package product

import (
	"github.com/sitename/sitename/model/shop"
)

func (a *AppProduct) GetDefaultDigitalContentSettings(aShop *shop.Shop) *shop.ShopDefaultDigitalContentSettings {
	return &shop.ShopDefaultDigitalContentSettings{
		AutomaticFulfillmentDigitalProducts: aShop.AutomaticFulfillmentDigitalProducts,
		DefaultDigitalMaxDownloads:          aShop.DefaultDigitalMaxDownloads,
		DefaultDigitalUrlValidDays:          aShop.DefaultDigitalUrlValidDays,
	}
}
