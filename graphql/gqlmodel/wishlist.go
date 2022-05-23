package gqlmodel

import (
	"time"

	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/modules/util"
)

type Wishlist struct {
	ID       string    `json:"id"`
	Token    string    `json:"token"`
	CreateAt time.Time `json:"createAt"`
	ItemIDs  []string  `json:"items"` // []*WishlistItem
}

func (Wishlist) IsNode() {}

type WishlistItem struct {
	ID         string    `json:"id"`
	CreateAt   time.Time `json:"createAt"`
	ProductID  *string   `json:"product"`  // *Product
	VariantIDs []string  `json:"variants"` // []*ProductVariant
}

func (WishlistItem) IsNode() {}

// SystemWishlistToGraphqlWishlist converts normal wishlist to graphql wishlist
func SystemWishlistToGraphqlWishlist(w *wishlist.Wishlist) *Wishlist {
	if w == nil {
		return nil
	}
	return &Wishlist{
		ID:       w.Id,
		CreateAt: util.TimeFromMillis(w.CreateAt),
		Token:    w.Token,
	}
}

// SystemWishlistItemsToGraphqlWishlistItems converts a list of wishlist items to a graphql list of wishlist items
func SystemWishlistItemsToGraphqlWishlistItems(ws []*wishlist.WishlistItem) []*WishlistItem {
	res := make([]*WishlistItem, len(ws))
	for i := range ws {
		res[i] = SystemWishlistItemToGraphqlWishlistItem(ws[i])
	}

	return res
}

// SystemWishlistItemToGraphqlWishlistItem converts database wishlist item to graphql wislist item
func SystemWishlistItemToGraphqlWishlistItem(wi *wishlist.WishlistItem) *WishlistItem {
	if wi == nil {
		return nil
	}

	return &WishlistItem{
		ID:        wi.Id,
		ProductID: &wi.ProductID,
		CreateAt:  util.TimeFromMillis(wi.CreateAt),
	}
}
