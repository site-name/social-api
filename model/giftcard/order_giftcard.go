package giftcard

import "github.com/sitename/sitename/model"

// OrderGiftCard is a relationship model between Order & GiftCard (m2m)
type OrderGiftCard struct {
	Id         string `json:"id"`
	GiftCardID string `json:"giftcard_id"`
	OrderID    string `json:"order_id"`
}

func (o *OrderGiftCard) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.order_giftcard.is_valid.%s.app_error",
		"order_giftcard_id=",
		"OrderGiftCard.IsValid",
	)
	if !model.IsValidId(o.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if !model.IsValidId(o.GiftCardID) {
		return outer("giftcard_id", &o.Id)
	}

	return nil
}

func (o *OrderGiftCard) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
}
