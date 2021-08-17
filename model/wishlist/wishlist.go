package wishlist

import (
	"github.com/sitename/sitename/model"
)

type Wishlist struct {
	Id       string  `json:"id"`
	Token    string  `json:"token"` // uuid, unique, not editable
	UserID   *string `json:"user_id"`
	CreateAt int64   `json:"create_at"`
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

func (w *Wishlist) ToJson() string {
	return model.ModelToJson(w)
}

func (w *Wishlist) PreSave() {
	if w.Id == "" {
		w.Id = model.NewId()
	}
	if w.Token == "" {
		w.Token = model.NewId()
	}
	w.CreateAt = model.GetMillis()
}
