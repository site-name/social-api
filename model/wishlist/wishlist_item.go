package wishlist

import (
	"io"

	"github.com/sitename/sitename/model"
)

type WishlistItem struct {
	Id         string `json:"id"`
	WishlistID string `json:"wishlist_id"`
	ProductID  string `json:"product_id"`
	CreateAt   int64  `json:"create_at"`
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

func (w *WishlistItem) ToJson() string {
	return model.ModelToJson(w)
}

func WishlistItemFromJson(data io.Reader) *WishlistItem {
	var w WishlistItem
	model.ModelFromJson(&w, data)
	return &w
}
