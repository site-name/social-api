package model

import (
	"github.com/Masterminds/squirrel"
)

type WishlistItem struct {
	Id         string `json:"id"`
	WishlistID string `json:"wishlist_id"`
	ProductID  string `json:"product_id"`
	CreateAt   int64  `json:"create_at"`
}

// WishlistItemFilterOption is used to build squirrel filtering options
type WishlistItemFilterOption struct {
	Id         squirrel.Sqlizer
	WishlistID squirrel.Sqlizer
	ProductID  squirrel.Sqlizer
}

type WishlistItems []*WishlistItem

func (w WishlistItems) IDs() []string {
	var res []string
	for _, item := range w {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (w WishlistItems) ProductIDs() []string {
	var res []string
	for _, item := range w {
		if item != nil {
			res = append(res, item.ProductID)
		}
	}

	return res
}

func (w *WishlistItem) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.wishlist_item.is_valid.%s.app_error",
		"wishlist_item_id=",
		"WishlistItem.IsValid",
	)
	if !IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !IsValidId(w.WishlistID) {
		return outer("wishlist_id", &w.Id)
	}
	if !IsValidId(w.ProductID) {
		return outer("product_id", &w.Id)
	}

	return nil
}

func (w *WishlistItem) PreSave() {
	if w.Id == "" {
		w.Id = NewId()
	}
	w.CreateAt = GetMillis()
}

func (w *WishlistItem) ToJSON() string {
	return ModelToJson(w)
}
