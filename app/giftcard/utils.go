package giftcard

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/modules/util"
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
func (s *ServiceGiftcard) FulfillNonShippableGiftcards(orDer *order.Order, orderLines order.OrderLines, siteSettings *shop.Shop, user *account.User, _ interface{}, manager interfaces.PluginManagerInterface) ([]*giftcard.GiftCard, *model.AppError) {
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
		Id:                 squirrel.Eq{s.srv.Store.OrderLine().TableName("Id"): giftcardLines.IDs()},
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
func (s *ServiceGiftcard) GiftcardsCreate(orDer *order.Order, giftcardLines order.OrderLines, quantities map[string]int, settings *shop.Shop, requestorUser *account.User, _ interface{}, manager interfaces.PluginManagerInterface) ([]*giftcard.GiftCard, *model.AppError) {
	var (
		customerUser          *account.User = nil
		customerUserID        *string
		appErr                *model.AppError
		userEmail             = orDer.UserEmail
		giftcards             = []*giftcard.GiftCard{}
		nonShippableGiftcards = []*giftcard.GiftCard{}
		expiryDate            = s.CalculateExpiryDate(settings)
	)

	if orDer.UserID != nil {
		customerUser, appErr = s.srv.AccountService().UserById(context.Background(), *orDer.UserID)
		if appErr != nil {
			return nil, appErr
		}
	}
	if customerUser != nil {
		customerUserID = &customerUser.Id
	}

	// refetch order lines with prefetching options
	giftcardLines, appErr = s.srv.OrderService().OrderLinesByOption(&order.OrderLineFilterOption{
		Id: squirrel.Eq{s.srv.Store.OrderLine().TableName("Id"): giftcardLines.IDs()},
		PrefetchRelated: order.OrderLinePrefetchRelated{
			VariantProduct: true,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	for _, orderLine := range giftcardLines {
		var (
			priceAmount   = orderLine.UnitPriceGrossAmount
			lineGiftcards = []*giftcard.GiftCard{}
			productID     *string
		)
		if orderLine.VariantID != nil && orderLine.ProductVariant != nil {
			productID = &orderLine.ProductVariant.ProductID
		}

		for i := 0; i < quantities[orderLine.Id]; i++ {

			lineGiftcards = append(lineGiftcards, &giftcard.GiftCard{
				Code:                 s.srv.DiscountService().GeneratePromoCode(),
				InitialBalanceAmount: priceAmount,
				CurrentBalanceAmount: priceAmount,
				CreatedByID:          customerUserID,
				CreatedByEmail:       &userEmail,
				ProductID:            productID,
				ExpiryDate:           expiryDate,
			})
		}

		giftcards = append(giftcards, lineGiftcards...)
		if !orderLine.IsShippingRequired {
			nonShippableGiftcards = append(nonShippableGiftcards, lineGiftcards...)
		}
	}

	giftcards, appErr = s.UpsertGiftcards(nil, giftcards...)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = s.GiftcardsBoughtEvent(nil, giftcards, orDer.Id, requestorUser, nil)
	if appErr != nil {
		return nil, appErr
	}

	channelOfOrder, appErr := s.srv.ChannelService().ChannelByOption(&channel.ChannelFilterOption{
		Id: squirrel.Eq{s.srv.Store.Channel().TableName("Id"): orDer.ChannelID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// send to customer all non-shippable gift cards
	appErr = s.SendGiftcardsToCustomer(nonShippableGiftcards, userEmail, requestorUser, nil, customerUser, manager, channelOfOrder.Slug)
	if appErr != nil {
		return nil, appErr
	}

	return giftcards, nil
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

func (s *ServiceGiftcard) FulfillGiftcardLines(giftcardLines order.OrderLines, requestorUser *account.User, _ interface{}, orDer *order.Order, manager interfaces.PluginManagerInterface) (interface{}, *model.AppError) {
	panic("not implt")
}

// CalculateExpiryDate calculate expiry date based on giftcard settings.
func (s *ServiceGiftcard) CalculateExpiryDate(shopSettings *shop.Shop) *time.Time {
	var (
		today      = util.StartOfDay(time.Now())
		expiryDate *time.Time
	)

	if shopSettings.GiftcardExpiryType == shop.EXPIRY_PERIOD {
		if expiryPeriod := shopSettings.GiftcardExpiryPeriod; expiryPeriod != nil {
			switch shopSettings.GiftcardExpiryPeriodType {
			case model.DAY:
				expiryDate = model.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * time.Hour))
			case model.WEEK:
				expiryDate = model.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * 7 * time.Hour))
			case model.MONTH:
				expiryDate = model.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * 30 * time.Hour))
			case model.YEAR:
				expiryDate = model.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * 365 * time.Hour))
			}
		}
	}

	return expiryDate
}

func (s *ServiceGiftcard) SendGiftcardsToCustomer(giftcards []*giftcard.GiftCard, userEmail string, requestorUser *account.User, _ interface{}, customerUser *account.User, manager interfaces.PluginManagerInterface, channelSlug string) *model.AppError {
	panic("not implemented")

}

func (s *ServiceGiftcard) DeactivateOrderGiftcards(orderID string, user *account.User, _ interface{}) *model.AppError {
	// giftcardEvents, appErr := s.GiftcardEventsByOptions(&giftcard.GiftCardEventFilterOption{
	// 	Type:       squirrel.Eq{s.srv.Store.GiftcardEvent().TableName("Type"): giftcard.BOUGHT},
	// 	Parameters: squirrel.Eq{s.srv.Store.GiftcardEvent().TableName("Parameters -> 'order_id'"): orderID}, // WHERE GiftcardEvents.Parameters -> 'order_id' = <something>
	// })
	// if appErr != nil {
	// 	if appErr.StatusCode == http.StatusInternalServerError {
	// 		return appErr
	// 	}
	// }

	// s.GiftcardsByOption(nil, &giftcard.GiftCardFilterOption{

	// })
	panic("not implemented")
}

func (s *ServiceGiftcard) OrderHasGiftcardLines(orDer *order.Order) (bool, *model.AppError) {
	orderLines, appErr := s.srv.OrderService().OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID:    squirrel.Eq{s.srv.Store.OrderLine().TableName("OrderID"): orDer.Id},
		IsGiftcard: model.NewBool(true),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
	}

	return len(orderLines) > 0, nil
}
