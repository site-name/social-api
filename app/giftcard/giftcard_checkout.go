package giftcard

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

// DeleteGiftCardCheckout drops a giftcard-checkout relation
func (s *AppGiftcard) DeleteGiftCardCheckout(giftcardID string, checkoutToken string) *model.AppError {
	err := s.Srv().Store.GiftCardCheckout().Delete(giftcardID, checkoutToken)
	if err != nil {
		if _, ok := err.(*store.ErrNotFound); ok {
			return nil
		}
		return model.NewAppError("DeleteGiftcardCheckout", "app.giftcard.error_deleting_giftcard_checkout_relation.app_error", nil, err.Error(), http.StatusExpectationFailed)
	}

	return nil
}

// CreateGiftCardCheckout create a new giftcard-checkout relation and returns it
func (a *AppGiftcard) CreateGiftCardCheckout(giftcardID string, checkoutToken string) (*giftcard.GiftCardCheckout, *model.AppError) {
	giftCardCheckout, err := a.Srv().Store.GiftCardCheckout().Save(&giftcard.GiftCardCheckout{
		GiftcardID: giftcardID,
		CheckoutID: checkoutToken,
	})
	if err != nil {
		var statusCode = http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("CreateGiftCardCheckout", "app.giftcard.error_creating_new_giftcard_checkout_relation.app_error", nil, err.Error(), statusCode)
	}

	return giftCardCheckout, nil
}
