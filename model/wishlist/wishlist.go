package wishlist

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
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

func (w *Wishlist) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.wishlist.is_valid.%s.app_error",
		"wishlist_id=",
		"Wishlist.IsValid",
	)
	if !model.IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(w.Token) {
		return outer("token", &w.Id)
	}
	if w.UserID != nil && !model.IsValidId(*w.UserID) {
		return outer("user_id", &w.Id)
	}
	if w.CreateAt == 0 {
		return outer("create_at", &w.Id)
	}

	return nil
}

func (w *Wishlist) ToJSON() string {
	return model.ModelToJson(w)
}

func (w *Wishlist) PreSave() {
	if !model.IsValidId(w.Id) {
		w.Id = model.NewId()
	}
	if !model.IsValidId(w.Token) {
		w.Token = model.NewId()
	}
	w.CreateAt = model.GetMillis()
}

func (w *Wishlist) PreUpdate() {
	if !model.IsValidId(w.Token) {
		w.Token = model.NewId()
	}
}
