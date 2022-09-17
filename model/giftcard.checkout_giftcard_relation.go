package model

// GiftCardCheckout represents relationship between giftcards-checkouts (m2m)
type GiftCardCheckout struct {
	Id         string `json:"id"`
	GiftcardID string `json:"giftcard_id"`
	CheckoutID string `json:"checkout_id"`
}

func (o *GiftCardCheckout) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"order_giftcard.is_valid.%s.app_error",
		"order_giftcard_id=",
		"GiftCardCheckout.IsValid",
	)
	if !IsValidId(o.Id) {
		return outer("id", nil)
	}
	if !IsValidId(o.CheckoutID) {
		return outer("checkout_id", &o.Id)
	}
	if !IsValidId(o.GiftcardID) {
		return outer("giftcard_id", &o.Id)
	}

	return nil
}

func (o *GiftCardCheckout) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}
}
