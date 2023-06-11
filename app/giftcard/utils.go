package giftcard

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// AddGiftcardCodeToCheckout adds giftcard data to checkout by code. Raise InvalidPromoCode if gift card cannot be applied.
func (a *ServiceGiftcard) AddGiftcardCodeToCheckout(ckout *model.Checkout, email, promoCode, currency string) (*model.InvalidPromoCode, *model.AppError) {
	now := time.Now().UTC()

	giftcards, appErr := a.GiftcardsByOption(&model.GiftCardFilterOption{
		Code:     squirrel.Eq{store.GiftcardTableName + ".Code": promoCode},
		Currency: squirrel.Eq{store.GiftcardTableName + ".Currency": strings.ToUpper(currency)},
		ExpiryDate: squirrel.Or{
			squirrel.GtOrEq{store.GiftcardTableName + ".ExpiryDate": now},
			squirrel.Eq{store.GiftcardTableName + ".ExpiryDate": nil},
		},
		StartDate: squirrel.LtOrEq{store.GiftcardTableName + ".StartDate": now},
		IsActive:  squirrel.Eq{store.GiftcardTableName + ".IsActive": true},
	})

	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound { // not found means promo code is invalid
			return &model.InvalidPromoCode{}, nil
		}
		return nil, appErr // if this is system error
	}

	// giftcard can be used only by one user
	if giftcards[0].UsedByEmail != nil || *giftcards[0].UsedByEmail != email {
		return &model.InvalidPromoCode{}, nil
	}

	_, appErr = a.CreateGiftCardCheckout(giftcards[0].Id, ckout.Token)
	return nil, appErr
}

// RemoveGiftcardCodeFromCheckout drops a relation between giftcard and checkout
func (a *ServiceGiftcard) RemoveGiftcardCodeFromCheckout(ckout *model.Checkout, giftcardCode string) *model.AppError {
	giftcards, appErr := a.GiftcardsByOption(&model.GiftCardFilterOption{
		Code: squirrel.Eq{store.GiftcardTableName + ".Code": giftcardCode},
	})

	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		giftcards = []*model.GiftCard{}
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
func (a *ServiceGiftcard) ToggleGiftcardStatus(giftCard *model.GiftCard) *model.AppError {
	if *giftCard.IsActive {
		giftCard.IsActive = model.NewPrimitive(false)
	} else {
		giftCard.IsActive = model.NewPrimitive(true)
	}

	_, appErr := a.UpsertGiftcards(nil, giftCard)
	if appErr != nil {
		return appErr
	}

	return nil
}

// FulfillNonShippableGiftcards
func (s *ServiceGiftcard) FulfillNonShippableGiftcards(orDer *model.Order, orderLines model.OrderLines, siteSettings *model.Shop, user *model.User, _ interface{}, manager interfaces.PluginManagerInterface) ([]*model.GiftCard, *model.InsufficientStock, *model.AppError) {
	if user != nil && !model.IsValidId(user.Id) {
		user = nil
	}

	giftcardLines, appErr := s.GetNonShippableGiftcardLines(orderLines)
	if appErr != nil {
		// this error caused by server
		return nil, nil, appErr
	}

	if len(giftcardLines) == 0 {
		return nil, nil, nil
	}

	_, inSufErr, appErr := s.FulfillGiftcardLines(giftcardLines, user, nil, orDer, manager)
	if inSufErr != nil || appErr != nil {
		return nil, inSufErr, appErr
	}

	var orderLineIDQuantityMap = map[string]int{} // orderLineIDQuantityMap has keys are order line ids
	for _, line := range giftcardLines {
		orderLineIDQuantityMap[line.Id] = line.Quantity
	}

	res, appErr := s.GiftcardsCreate(orDer, giftcardLines, orderLineIDQuantityMap, siteSettings, user, nil, manager)
	return res, nil, appErr
}

func (s *ServiceGiftcard) GetNonShippableGiftcardLines(lines model.OrderLines) (model.OrderLines, *model.AppError) {
	giftcardLines := GetGiftcardLines(lines)
	nonShippableLines, appErr := s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Id:                 squirrel.Eq{store.OrderLineTableName + ".Id": giftcardLines.IDs()},
		IsShippingRequired: model.NewPrimitive(true),
	})

	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	return nonShippableLines, nil
}

// GiftcardsCreate creates purchased gift cards
func (s *ServiceGiftcard) GiftcardsCreate(orDer *model.Order, giftcardLines model.OrderLines, quantities map[string]int, settings *model.Shop, requestorUser *model.User, _ interface{}, manager interfaces.PluginManagerInterface) ([]*model.GiftCard, *model.AppError) {
	var (
		customerUser          *model.User = nil
		customerUserID        *string
		appErr                *model.AppError
		userEmail             = orDer.UserEmail
		giftcards             = []*model.GiftCard{}
		nonShippableGiftcards = []*model.GiftCard{}
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
	giftcardLines, appErr = s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Id: squirrel.Eq{store.OrderLineTableName + ".Id": giftcardLines.IDs()},
		PrefetchRelated: model.OrderLinePrefetchRelated{
			VariantProduct: true,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	for _, orderLine := range giftcardLines {
		var (
			priceAmount   = orderLine.UnitPriceGrossAmount
			lineGiftcards = []*model.GiftCard{}
			productID     *string
		)
		if orderLine.VariantID != nil && orderLine.GetProductVariant() != nil {
			productID = &orderLine.GetProductVariant().ProductID
		}

		for i := 0; i < quantities[orderLine.Id]; i++ {

			lineGiftcards = append(lineGiftcards, &model.GiftCard{
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

	channelOfOrder, appErr := s.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": orDer.ChannelID},
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

func GetGiftcardLines(lines model.OrderLines) model.OrderLines {
	res := model.OrderLines{}
	for _, line := range lines {
		if line != nil && line.IsGiftcard {
			res = append(res, line)
		}
	}

	return res
}

func (s *ServiceGiftcard) FulfillGiftcardLines(giftcardLines model.OrderLines, requestorUser *model.User, _ interface{}, order *model.Order, manager interfaces.PluginManagerInterface) ([]*model.Fulfillment, *model.InsufficientStock, *model.AppError) {
	if len(giftcardLines) == 0 {
		return nil, nil, model.NewAppError("FulfillGiftcardLines", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "giftcardLines"}, "", http.StatusBadRequest)
	}

	// check if we need to prefetch related values for given order lines:
	if giftcardLines[0].GetAllocations().Len() == 0 ||
		giftcardLines[0].GetProductVariant() == nil {

		var appErr *model.AppError
		giftcardLines, appErr = s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
			Id: squirrel.Eq{store.OrderLineTableName + ".Id": giftcardLines.IDs()},
			PrefetchRelated: model.OrderLinePrefetchRelated{
				AllocationsStock: true,
				VariantStocks:    true,
			},
		})
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	var linesForWarehouses = map[string][]*model.QuantityOrderLine{}

	for _, orderLine := range giftcardLines {
		if orderLine.GetAllocations().Len() > 0 {
			for _, allocation := range orderLine.GetAllocations() {

				if allocation.QuantityAllocated > 0 {

					linesForWarehouses[allocation.Stock.WarehouseID] = append(
						linesForWarehouses[allocation.Stock.WarehouseID],
						&model.QuantityOrderLine{
							OrderLine: orderLine,
							Quantity:  allocation.QuantityAllocated,
						},
					)
				}
			}
		} else {

			stocks, appErr := s.srv.WarehouseService().FilterStocksForChannel(&model.StockFilterForChannelOption{
				Id:        squirrel.Eq{store.StockTableName + ".Id": orderLine.GetProductVariant().GetStocks().IDs()},
				ChannelID: order.ChannelID,
			})
			if appErr != nil {
				if appErr.StatusCode != http.StatusNotFound {
					return nil, nil, appErr
				}

				return nil,
					&model.InsufficientStock{
						Code: model.GIFT_CARD_NOT_APPLICABLE,
						Items: []*model.InsufficientStockData{
							{
								Variant: *orderLine.GetProductVariant(),
							},
						},
					},
					nil
			}

			linesForWarehouses[stocks[0].WarehouseID] = append(
				linesForWarehouses[stocks[0].WarehouseID],
				&model.QuantityOrderLine{
					OrderLine: orderLine,
					Quantity:  orderLine.Quantity,
				},
			)
		}
	}

	return s.srv.OrderService().CreateFulfillments(requestorUser, nil, order, linesForWarehouses, manager, true, true, false)
}

// CalculateExpiryDate calculate expiry date based on giftcard settings.
func (s *ServiceGiftcard) CalculateExpiryDate(shopSettings *model.Shop) *time.Time {
	var (
		today      = util.StartOfDay(time.Now().UTC())
		expiryDate *time.Time
	)

	if shopSettings.GiftcardExpiryType == model.EXPIRY_PERIOD {
		if expiryPeriod := shopSettings.GiftcardExpiryPeriod; expiryPeriod != nil {
			switch shopSettings.GiftcardExpiryPeriodType {
			case model.DAY:
				expiryDate = util.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * time.Hour))
			case model.WEEK:
				expiryDate = util.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * 7 * time.Hour))
			case model.MONTH:
				expiryDate = util.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * 30 * time.Hour))
			case model.YEAR:
				expiryDate = util.NewTime(today.Add(time.Duration(*expiryPeriod) * 24 * 365 * time.Hour))
			}
		}
	}

	return expiryDate
}

func (s *ServiceGiftcard) SendGiftcardsToCustomer(giftcards []*model.GiftCard, userEmail string, requestorUser *model.User, _ interface{}, customerUser *model.User, manager interfaces.PluginManagerInterface, channelSlug string) *model.AppError {
	for _, gc := range giftcards {
		appErr := s.SendGiftcardNotification(requestorUser, nil, customerUser, userEmail, *gc, manager, channelSlug, false)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (s *ServiceGiftcard) DeactivateOrderGiftcards(orderID string, user *model.User, _ interface{}) *model.AppError {
	giftcardIDs, err := s.srv.Store.GiftCard().DeactivateOrderGiftcards(orderID)
	if err != nil {
		return model.NewAppError("DeactivateOrderGiftcards", "app.giftcard.error_updating_giftcards.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	var userID *string
	if user != nil {
		userID = &user.Id
	}

	var events []*model.GiftCardEvent
	for _, id := range giftcardIDs {
		events = append(events, &model.GiftCardEvent{
			UserID:     userID,
			GiftcardID: id,
			Type:       model.GIFT_CARD_EVENT_TYPE_DEACTIVATED,
		})
	}

	_, appErr := s.BulkUpsertGiftcardEvents(nil, events...)
	return appErr
}

func (s *ServiceGiftcard) OrderHasGiftcardLines(orDer *model.Order) (bool, *model.AppError) {
	orderLines, appErr := s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		OrderID:    squirrel.Eq{store.OrderLineTableName + ".OrderID": orDer.Id},
		IsGiftcard: model.NewPrimitive(true),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
	}

	return len(orderLines) > 0, nil
}
