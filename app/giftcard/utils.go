package giftcard

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

// AddGiftcardCodeToCheckout adds giftcard data to checkout by code. Raise InvalidPromoCode if gift card cannot be applied.
func (a *ServiceGiftcard) AddGiftcardCodeToCheckout(checkout *model.Checkout, email, promoCode, currency string) (*model.InvalidPromoCode, *model_helper.AppError) {
	now := time.Now()

	_, giftcards, appErr := a.GiftcardsByOption(&model.GiftCardFilterOption{
		Conditions: squirrel.And{
			squirrel.Expr(model.GiftcardTableName+".Code = ?", promoCode),
			squirrel.Expr(model.GiftcardTableName+".Currency = ?", strings.ToUpper(currency)),
			squirrel.Expr(model.GiftcardTableName+".StartDate <= ?", now),
			squirrel.Expr(model.GiftcardTableName + ".IsActive"),
			squirrel.Or{
				squirrel.Expr(model.GiftcardTableName+".ExpiryDate >= ?", now),
				squirrel.Expr(model.GiftcardTableName + ".ExpiryDate IS NULL"),
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(giftcards) == 0 {
		return &model.InvalidPromoCode{}, nil
	}

	giftcard := giftcards[0]

	// giftcard can be used only by one user
	if giftcard.UsedByEmail != nil && *giftcard.UsedByEmail != email {
		return &model.InvalidPromoCode{}, nil
	}

	return nil, a.AddGiftcardRelations(nil, model.Giftcards{giftcard}, []*model.Checkout{checkout})
}

// RemoveGiftcardCodeFromCheckout drops a relation between giftcard and checkout
func (a *ServiceGiftcard) RemoveGiftcardCodeFromCheckout(checkout *model.Checkout, giftcardCode string) *model_helper.AppError {
	_, giftcards, appErr := a.GiftcardsByOption(&model.GiftCardFilterOption{
		Conditions: squirrel.Expr(model.GiftcardTableName+".Code = ?", giftcardCode),
	})
	if appErr != nil {
		return appErr
	}
	if len(giftcards) == 0 {
		return nil
	}

	return a.RemoveGiftcardRelations(nil, giftcards, []*model.Checkout{checkout})
}

// ToggleGiftcardStatus set status of given giftcard to inactive/active
func (a *ServiceGiftcard) ToggleGiftcardStatus(giftCard *model.GiftCard) *model_helper.AppError {
	if *giftCard.IsActive {
		giftCard.IsActive = model.GetPointerOfValue(false)
	} else {
		giftCard.IsActive = model.GetPointerOfValue(true)
	}

	_, appErr := a.UpsertGiftcards(nil, giftCard)
	if appErr != nil {
		return appErr
	}

	return nil
}

// FulfillNonShippableGiftcards
func (s *ServiceGiftcard) FulfillNonShippableGiftcards(order *model.Order, orderLines model.OrderLines, siteSettings model.ShopSettings, user *model.User, _ interface{}, manager interfaces.PluginManagerInterface) ([]*model.GiftCard, *model.InsufficientStock, *model_helper.AppError) {
	if user != nil && !model_helper.IsValidId(user.Id) {
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

	_, inSufErr, appErr := s.FulfillGiftcardLines(giftcardLines, user, nil, order, manager)
	if inSufErr != nil || appErr != nil {
		return nil, inSufErr, appErr
	}

	var orderLineIDQuantityMap = map[string]int{} // orderLineIDQuantityMap has keys are order line ids
	for _, line := range giftcardLines {
		orderLineIDQuantityMap[line.Id] = line.Quantity
	}

	res, appErr := s.GiftcardsCreate(nil, order, giftcardLines, orderLineIDQuantityMap, siteSettings, user, nil, manager)
	return res, nil, appErr
}

func (s *ServiceGiftcard) GetNonShippableGiftcardLines(lines model.OrderLines) (model.OrderLines, *model_helper.AppError) {
	giftcardLines := GetGiftcardLines(lines)
	nonShippableLines, appErr := s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{
			model.OrderLineTableName + ".Id":                 giftcardLines.IDs(),
			model.OrderLineTableName + ".IsShippingRequired": true,
		},
	})

	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	return nonShippableLines, nil
}

// GiftcardsCreate creates purchased gift cards
func (s *ServiceGiftcard) GiftcardsCreate(tx *gorm.DB, order *model.Order, giftcardLines model.OrderLines, quantities map[string]int, settings model.ShopSettings, requestorUser *model.User, _ interface{}, manager interfaces.PluginManagerInterface) ([]*model.GiftCard, *model_helper.AppError) {
	var (
		customerUser          *model.User = nil
		customerUserID        *string
		appErr                *model_helper.AppError
		userEmail             = order.UserEmail
		giftcards             = []*model.GiftCard{}
		nonShippableGiftcards = []*model.GiftCard{}
		expiryDate            = s.CalculateExpiryDate(settings)
	)

	if order.UserID != nil {
		customerUser, appErr = s.srv.AccountService().UserById(context.Background(), *order.UserID)
		if appErr != nil {
			return nil, appErr
		}
	}
	if customerUser != nil {
		customerUserID = &customerUser.Id
	}

	// refetch order lines with prefetching options
	giftcardLines, appErr = s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": giftcardLines.IDs()},
		Preload:    []string{"ProductVariant.Product"},
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
		if orderLine.VariantID != nil && orderLine.ProductVariant != nil {
			productID = &orderLine.ProductVariant.ProductID
		}

		for i := 0; i < quantities[orderLine.Id]; i++ {
			lineGiftcards = append(lineGiftcards, &model.GiftCard{
				Code:                 model.NewPromoCode(),
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

	giftcards, appErr = s.UpsertGiftcards(tx, giftcards...)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = s.GiftcardsBoughtEvent(tx, giftcards, order.Id, requestorUser, nil)
	if appErr != nil {
		return nil, appErr
	}

	channelOfOrder, appErr := s.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": order.ChannelID},
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

func (s *ServiceGiftcard) FulfillGiftcardLines(giftcardLines model.OrderLines, requestorUser *model.User, _ interface{}, order *model.Order, manager interfaces.PluginManagerInterface) ([]*model.Fulfillment, *model.InsufficientStock, *model_helper.AppError) {
	if len(giftcardLines) == 0 {
		return nil, nil, model_helper.NewAppError("FulfillGiftcardLines", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "giftcardLines"}, "", http.StatusBadRequest)
	}

	// check if we need to prefetch related values for given order lines:
	if giftcardLines[0].Allocations.Len() == 0 ||
		giftcardLines[0].ProductVariant == nil {

		var appErr *model_helper.AppError
		giftcardLines, appErr = s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
			Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": giftcardLines.IDs()},
			Preload: []string{ // TODO: check if this works
				"Allocations.Stock",
				"ProductVariant.Stocks",
			},
		})
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	var linesForWarehouses = map[string][]*model.QuantityOrderLine{}

	for _, orderLine := range giftcardLines {
		if orderLine.Allocations.Len() > 0 {
			for _, allocation := range orderLine.Allocations {

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
				Conditions: squirrel.Eq{model.StockTableName + ".Id": orderLine.ProductVariant.Stocks.IDs()},
				ChannelID:  order.ChannelID,
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
								Variant: *orderLine.ProductVariant,
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
func (s *ServiceGiftcard) CalculateExpiryDate(shopSettings model.ShopSettings) *time.Time {
	var (
		today      = util.StartOfDay(time.Now().UTC())
		expiryDate *time.Time
	)

	if *shopSettings.GiftcardExpiryType == model.EXPIRY_PERIOD {
		if expiryPeriod := shopSettings.GiftcardExpiryPeriod; expiryPeriod != nil {
			switch *shopSettings.GiftcardExpiryPeriodType {
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

func (s *ServiceGiftcard) SendGiftcardsToCustomer(giftcards []*model.GiftCard, userEmail string, requestorUser *model.User, _ interface{}, customerUser *model.User, manager interfaces.PluginManagerInterface, channelSlug string) *model_helper.AppError {
	for _, gc := range giftcards {
		appErr := s.SendGiftcardNotification(requestorUser, nil, customerUser, userEmail, *gc, manager, channelSlug, false)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (s *ServiceGiftcard) DeactivateOrderGiftcards(tx *gorm.DB, orderID string, user *model.User, _ interface{}) *model_helper.AppError {
	giftcardIDs, err := s.srv.Store.GiftCard().DeactivateOrderGiftcards(tx, orderID)
	if err != nil {
		return model_helper.NewAppError("DeactivateOrderGiftcards", "app.giftcard.error_updating_giftcards.app_error", nil, err.Error(), http.StatusInternalServerError)
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

	_, appErr := s.BulkUpsertGiftcardEvents(tx, events...)
	return appErr
}

func (s *ServiceGiftcard) OrderHasGiftcardLines(order *model.Order) (bool, *model_helper.AppError) {
	orderLines, appErr := s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{
			model.OrderLineTableName + ".OrderID":    order.Id,
			model.OrderLineTableName + ".IsGiftcard": true,
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
	}

	return len(orderLines) > 0, nil
}
