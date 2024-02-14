package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

type WishlistFilterOption struct {
	CommonQueryOptions
}

type WishlistItemFilterOption struct {
	CommonQueryOptions
}

func WishlistIsValid(w model.Wishlist) *AppError {
	if !IsValidId(w.UserID) {
		return NewAppError("Wishlist.IsValid", "model.wishlist.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(w.ID) {
		return NewAppError("Wishlist.IsValid", "model.wishlist.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if w.CreatedAt <= 0 {
		return NewAppError("Wishlist.IsValid", "model.wishlist.is_valid.created_at.app_error", nil, "please provide valid create_at", http.StatusBadRequest)
	}

	return nil
}

func WishlistItemIsValid(wi model.WishlistItem) *AppError {
	if !IsValidId(wi.WishlistID) {
		return NewAppError("WishlistItem.IsValid", "model.wishlist_item.is_valid.wishlist_id.app_error", nil, "please provide valid wishlist id", http.StatusBadRequest)
	}
	if !IsValidId(wi.ProductID) {
		return NewAppError("WishlistItem.IsValid", "model.wishlist_item.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}
	if wi.CreatedAt <= 0 {
		return NewAppError("WishlistItem.IsValid", "model.wishlist_item.is_valid.created_at.app_error", nil, "please provide valid create_at", http.StatusBadRequest)
	}
	if !IsValidId(wi.ID) {
		return NewAppError("WishlistItem.IsValid", "model.wishlist_item.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}

	return nil
}

func WishlistPreSave(w *model.Wishlist) {
	w.ID = NewId()
	w.CreatedAt = GetMillis()
}

func WishlistItemPreSave(wi *model.WishlistItem) {
	wi.ID = NewId()
	wi.CreatedAt = GetMillis()
}
