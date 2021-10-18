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

// DatabaseWishlistToGraphqlWishlist converts normal wishlist to graphql wishlist
func DatabaseWishlistToGraphqlWishlist(w *wishlist.Wishlist) *Wishlist {

	return &Wishlist{
		ID:       w.Id,
		CreateAt: util.TimeFromMillis(w.CreateAt),
		Token:    w.Token,
	}
}

// DatabaseWishlistItemsToGraphqlWishlistItems converts a list of wishlist items to a graphql list of wishlist items
func DatabaseWishlistItemsToGraphqlWishlistItems(ws []*wishlist.WishlistItem) []*WishlistItem {
	res := make([]*WishlistItem, len(ws))
	for i := range ws {
		res[i] = DatabaseWishlistItemToGraphqlWishlistItem(ws[i])
	}

	return res
}

// DatabaseWishlistItemToGraphqlWishlistItem converts database wishlist item to graphql wislist item
func DatabaseWishlistItemToGraphqlWishlistItem(wi *wishlist.WishlistItem) *WishlistItem {

	return &WishlistItem{
		ID:        wi.Id,
		ProductID: &wi.ProductID,
		CreateAt:  util.TimeFromMillis(wi.CreateAt),
	}
}
