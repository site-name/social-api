package giftcard

import (
	"net/http"
	"time"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
)

// AddGiftcardCodeToCheckout adds giftcard data to checkout by code.
func (a *AppGiftcard) AddGiftcardCodeToCheckout(ckout *checkout.Checkout, promoCode string) *model.AppError {
	now := model.NewTime(time.Now()) // NOT: not sure use UTC or system time

	giftcards, appErr := a.GiftcardsByOption(&giftcard.GiftCardFilterOption{
		Code: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: promoCode,
			},
		},
		EndDate: &model.TimeFilter{
			Or: &model.TimeOption{
				NULL: model.NewBool(true),
				GtE:  now,
			},
		},
		StartDate: &model.TimeFilter{
			TimeOption: &model.TimeOption{
				LtE: now,
			},
		},
		IsActive: model.NewBool(true),
	})

	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound { // not found means promo code is invalid
			return model.NewAppError("AddGiftcardCodeToCheckout", app.InvalidPromoCodeAppErrorID, map[string]interface{}{"PromoCode": promoCode}, "", http.StatusBadRequest)
		}
		return appErr // if this is system error
	}

	_, appErr = a.CreateGiftCardCheckout(giftcards[0].Id, ckout.Token)
	return appErr
}

// RemoveGiftcardCodeFromCheckout drops a relation between giftcard and checkout
func (a *AppGiftcard) RemoveGiftcardCodeFromCheckout(ckout *checkout.Checkout, giftcardCode string) *model.AppError {
	giftcards, appErr := a.GiftcardsByOption(&giftcard.GiftCardFilterOption{
		Code: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: giftcardCode,
			},
		},
	})

	if appErr != nil {
		return appErr
	}

	if len(giftcards) > 0 {
		appErr := a.DeleteGiftCardCheckout(giftcards[0].Id, ckout.Token)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// ToggleGiftcardStatus set status of given giftcard to inactive/active
func (a *AppGiftcard) ToggleGiftcardStatus(giftCard *giftcard.GiftCard) *model.AppError {
	if *giftCard.IsActive {
		giftCard.IsActive = model.NewBool(false)
	} else {
		giftCard.IsActive = model.NewBool(true)
	}

	_, appErr := a.UpdateGiftCard(giftCard)
	if appErr != nil {
		appErr.Where = "ToggleGiftcardStatus" // this lets us know where does the error come from
		return appErr
	}

	return nil
}