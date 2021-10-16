package giftcard

import (
	"net/http"
	"strings"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/shop"
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

	_, appErr := a.UpsertGiftcards(nil, giftCard)
	if appErr != nil {
		return appErr
	}

	return nil
}

// FulfillNonShippableGiftcards
func (s *ServiceGiftcard) FulfillNonShippableGiftcards(orDer *order.Order, orderLines order.OrderLines, siteSettings *shop.Shop, user *account.User, _ interface{}, manager interface{}) ([]*giftcard.GiftCard, *model.AppError) {
	if user != nil && !model.IsValidId(user.Id) {
		user = nil
	}

	giftcardLines, appErr := s.GetNonShippableGiftcardLines(orderLines)
	if appErr != nil {
		// this error caused by server
		return nil, appErr
	}

	if len(giftcardLines) == 0 {
		return nil, nil
	}

	_, appErr = s.FulfillGiftcardLines(giftcardLines, user, nil, orDer, manager)
	if appErr != nil {
		return nil, appErr
	}

	var orderLineIDQuantityMap = map[string]int{} // orderLineIDQuantityMap has keys are order line ids
	for _, line := range giftcardLines {
		orderLineIDQuantityMap[line.Id] = line.Quantity
	}

	return s.GiftcardsCreate(orDer, giftcardLines, orderLineIDQuantityMap, siteSettings, user, nil, manager)
}

func (s *ServiceGiftcard) GetNonShippableGiftcardLines(lines order.OrderLines) (order.OrderLines, *model.AppError) {
	giftcardLines := GetGiftcardLines(lines)
	nonShippableLines, appErr := s.srv.OrderService().OrderLinesByOption(&order.OrderLineFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: giftcardLines.IDs(),
			},
		},
		IsShippingRequired: model.NewBool(true),
	})

	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	return nonShippableLines, nil
}

// GiftcardsCreate creates purchased gift cards
func (s *ServiceGiftcard) GiftcardsCreate(orDer *order.Order, giftcardLines order.OrderLines, quantities map[string]int, settings *shop.Shop, requestorUser *account.User, _ interface{}, manager interface{}) ([]*giftcard.GiftCard, *model.AppError) {
	var (
		customerUserID        = orDer.UserID
		userEmail             = orDer.UserEmail
		giftcards             = []*giftcard.GiftCard{}
		nonShippableGiftcards = []*giftcard.GiftCard{}
	)
}

func GetGiftcardLines(lines order.OrderLines) order.OrderLines {
	res := order.OrderLines{}
	for _, line := range lines {
		if line != nil && line.IsGiftcard {
			res = append(res, line)
		}
	}

	return res
}

func (s *ServiceGiftcard) FulfillGiftcardLines(giftcardLines order.OrderLines, requestorUser *account.User, _ interface{}, orDer *order.Order, manager interface{}) (interface{}, *model.AppError) {
	panic("not implt")
}

// CalculateExpiryDate calculate expiry date based on giftcard settings.
func (s *ServiceGiftcard) CalculateExpiryDate(shopSettings *shop.Shop) {
	panic("not implt")
}
