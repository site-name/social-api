package wishlist

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
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

func (w *WishlistItem) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.wishlist_item.is_valid.%s.app_error",
		"wishlist_item_id=",
		"WishlistItem.IsValid",
	)
	if !model.IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(w.WishlistID) {
		return outer("wishlist_id", &w.Id)
	}
	if !model.IsValidId(w.ProductID) {
		return outer("product_id", &w.Id)
	}

	return nil
}

func (w *WishlistItem) PreSave() {
	if w.Id == "" {
		w.Id = model.NewId()
	}
	w.CreateAt = model.GetMillis()
}

func (w *WishlistItem) ToJSON() string {
	return model.ModelToJson(w)
}
