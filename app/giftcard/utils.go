package giftcard

import (
	"net/http"
	"strings"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
)

// AddGiftcardCodeToCheckout adds giftcard data to checkout by code. Raise InvalidPromoCode if gift card cannot be applied.
func (a *ServiceGiftcard) AddGiftcardCodeToCheckout(ckout *checkout.Checkout, email, promoCode, currency string) (*giftcard.InvalidPromoCode, *model.AppError) {
	now := model.NewTime(time.Now())

	giftcards, appErr := a.GiftcardsByOption(nil, &giftcard.GiftCardFilterOption{
		Code: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: promoCode,
			},
		},
		Currency: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: strings.ToUpper(currency),
			},
		},
		ExpiryDate: &model.TimeFilter{
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
			return &giftcard.InvalidPromoCode{}, nil
		}
		return nil, appErr // if this is system error
	}

	// giftcard can be used only by one user
	if giftcards[0].UsedByEmail != nil || *giftcards[0].UsedByEmail != email {
		return &giftcard.InvalidPromoCode{}, nil
	}

	_, appErr = a.CreateGiftCardCheckout(giftcards[0].Id, ckout.Token)
	return nil, appErr
}

// RemoveGiftcardCodeFromCheckout drops a relation between giftcard and checkout
func (a *ServiceGiftcard) RemoveGiftcardCodeFromCheckout(ckout *checkout.Checkout, giftcardCode string) *model.AppError {
	giftcards, appErr := a.GiftcardsByOption(nil, &giftcard.GiftCardFilterOption{
		Code: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: giftcardCode,
			},
		},
	})

	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		giftcards = []*giftcard.GiftCard{}
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
func (a *ServiceGiftcard) ToggleGiftcardStatus(giftCard *giftcard.GiftCard) *model.AppError {
	if *giftCard.IsActive {
		giftCard.IsActive = model.NewBool(false)
	} else {
		giftCard.IsActive = model.NewBool(true)
	}

	_, appErr := a.UpsertGiftcard(giftCard)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (s *ServiceGiftcard) FulfillNonShippableGiftcards(orDer *order.Order, orderLines order.OrderLines) {

}

func (s *ServiceGiftcard) GetNonShippableGiftcardLines(lineIDs []string) {

}
