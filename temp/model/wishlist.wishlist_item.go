package model

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type WishlistItem struct {
	Id         string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	WishlistID string `json:"wishlist_id" gorm:"type:uuid;column:WishlistID;index:wishlistid_productid_key"`
	ProductID  string `json:"product_id" gorm:"type:uuid;column:ProductID;index:wishlistid_productid_key"`
	CreateAt   int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`

	ProductVariants ProductVariants `json:"-" gorm:"many2many:WishlistItemProductVariants"`
}

// column names for table wishlistItem
const (
	WishlistItemColumnId         = "Id"
	WishlistItemColumnWishlistID = "WishlistID"
	WishlistItemColumnProductID  = "ProductID"
	WishlistItemColumnCreateAt   = "CreateAt"
)

func (t *WishlistItem) TableName() string             { return WishlistItemTableName }
func (t *WishlistItem) BeforeCreate(_ *gorm.DB) error { return t.IsValid() }
func (t *WishlistItem) BeforeUpdate(_ *gorm.DB) error {
	// prevent update
	t.CreateAt = 0
	return t.IsValid()
}

// WishlistItemFilterOption is used to build squirrel filtering options
type WishlistItemFilterOption struct {
	Conditions squirrel.Sqlizer
}

type WishlistItems []*WishlistItem

func (w WishlistItems) IDs() []string {
	return lo.Map(w, func(item *WishlistItem, _ int) string { return item.Id })
}

func (w WishlistItems) ProductIDs() []string {
	return lo.Map(w, func(item *WishlistItem, _ int) string { return item.ProductID })
}

func (w *WishlistItem) IsValid() *AppError {
	if !IsValidId(w.WishlistID) {
		return NewAppError("WishlistItem.IsValid", "model.wishlist_item.is_valid.wishlist_id.app_error", nil, "please proivde valid wishlist id", http.StatusBadRequest)
	}
	if !IsValidId(w.ProductID) {
		return NewAppError("WishlistItem.IsValid", "model.wishlist_item.is_valid.product_id.app_error", nil, "please proivde valid product id", http.StatusBadRequest)
	}
	return nil
}
