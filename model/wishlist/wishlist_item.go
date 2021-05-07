package wishlist

import (
	"github.com/sitename/sitename/model/product_and_discount"
)

type WishlishItem struct {
	Id         string                                 `json:"id"`
	WishListID string                                 `json:"wishlist_id"`
	ProductID  string                                 `json:"product_id"`
	Variants   []*product_and_discount.ProductVariant `json:"variants" db:"-"`
	CreateAt   int64                                  `json:"create_at"`
	Wishlish   *Wishlish                              `json:"wishlist" db:"-"`
	Product    *product_and_discount.Product          `json:"product" db:"-"`
}
