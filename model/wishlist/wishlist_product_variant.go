package wishlist

import "github.com/sitename/sitename/model"

// WishlistVariant represents relationships between wishlists and product variants
type WishlistProductVariant struct {
	Id               string `json:"id"`
	WishlistItemID   string `json:"wishlist_item_id"`
	ProductVariantID string `json:"product_variant_id"`
}

func (w *WishlistProductVariant) PreSave() {
	if w.Id == "" {
		w.Id = model.NewId()
	}
}

func (w *WishlistProductVariant) ToJson() string {
	return model.ModelToJson(w)
}

func (w *WishlistProductVariant) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.wishlist_product_variant.is_valid.%s.app_error",
		"wishlist_product_variant_id=",
		"WishlistProductVariant.IsValid",
	)
	if !model.IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(w.WishlistItemID) {
		return outer("wishlist_item_id", &w.Id)
	}
	if !model.IsValidId(w.ProductVariantID) {
		return outer("product_variant_id", &w.Id)
	}

	return nil
}
