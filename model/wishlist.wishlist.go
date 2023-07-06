package model

import (
	"github.com/Masterminds/squirrel"
)

type Wishlist struct {
	Id       string  `json:"id"`
	Token    string  `json:"token"` // uuid, unique, not editable
	UserID   *string `json:"user_id"`
	CreateAt int64   `json:"create_at"`
}

// WishlistFilterOption is used to build squirrel sql queries
type WishlistFilterOption struct {
	Id     squirrel.Sqlizer
	Token  squirrel.Sqlizer
	UserID squirrel.Sqlizer
}

func (w *Wishlist) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.wishlist.is_valid.%s.app_error",
		"wishlist_id=",
		"Wishlist.IsValid",
	)
	if !IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !IsValidId(w.Token) {
		return outer("token", &w.Id)
	}
	if w.UserID != nil && !IsValidId(*w.UserID) {
		return outer("user_id", &w.Id)
	}
	if w.CreateAt == 0 {
		return outer("create_at", &w.Id)
	}

	return nil
}

func (w *Wishlist) PreSave() {
	if !IsValidId(w.Id) {
		w.Id = NewId()
	}
	if !IsValidId(w.Token) {
		w.Token = NewId()
	}
	w.CreateAt = GetMillis()
}

func (w *Wishlist) PreUpdate() {
	if !IsValidId(w.Token) {
		w.Token = NewId()
	}
}
