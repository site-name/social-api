package model

import "github.com/Masterminds/squirrel"

// OrderGiftCard is a relationship model between Order & GiftCard (m2m)
type OrderGiftCard struct {
	Id         string `json:"id"`
	GiftCardID string `json:"giftcard_id"` // unique together with orderID
	OrderID    string `json:"order_id"`    // unique together with GiftCardID
}

type OrderGiftCardFilterOptions struct {
	GiftCardID squirrel.Sqlizer
	OrderID    squirrel.Sqlizer
}

func (o *OrderGiftCard) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"order_giftcard.is_valid.%s.app_error",
		"order_giftcard_id=",
		"OrderGiftCard.IsValid",
	)
	if !IsValidId(o.Id) {
		return outer("id", nil)
	}
	if !IsValidId(o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if !IsValidId(o.GiftCardID) {
		return outer("giftcard_id", &o.Id)
	}

	return nil
}

func (o *OrderGiftCard) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}
}
