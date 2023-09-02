package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type Wishlist struct {
	Id       UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Token    string `json:"token" gorm:"type:uuid;default:gen_random_uuid();column:Token;unique"` // uuid, unique, not editable
	UserID   *UUID  `json:"user_id" gorm:"type:uuid;column:UserID"`
	CreateAt int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
}

func (t *Wishlist) TableName() string             { return WishlistTableName }
func (t *Wishlist) BeforeCreate(_ *gorm.DB) error { return t.IsValid() }
func (t *Wishlist) BeforeUpdate(_ *gorm.DB) error {
	// prevent update
	t.Token = ""
	t.CreateAt = 0
	return t.IsValid()
}

type WishlistFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (w *Wishlist) IsValid() *AppError {
	if w.UserID != nil && !IsValidId(*w.UserID) {
		return NewAppError("Wishlist.IsValid", "model.wishlist.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}

	return nil
}
