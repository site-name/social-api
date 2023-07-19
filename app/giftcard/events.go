package giftcard

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// GiftcardEventsByOptions returns a list of giftcard events filtered using given options
func (s *ServiceGiftcard) GiftcardEventsByOptions(options *model.GiftCardEventFilterOption) ([]*model.GiftCardEvent, *model.AppError) {
	events, err := s.srv.Store.GiftcardEvent().FilterByOptions(options)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(events) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("GiftcardEventsByOptions", "app.giftcard.error_finding_giftcard_events_by_options.app_error", nil, errMessage, statusCode)
	}

	return events, nil
}

// BulkUpsertGiftcardEvents tells store to upsert given giftcard events into database then returns them
func (s *ServiceGiftcard) BulkUpsertGiftcardEvents(transaction *gorm.DB, events ...*model.GiftCardEvent) ([]*model.GiftCardEvent, *model.AppError) {
	events, err := s.srv.Store.GiftcardEvent().BulkUpsert(transaction, events...)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("BulkUpsertGiftcardEvents", "app.giftcard.error_upserting_giftcard_events.app_error", nil, err.Error(), statusCode)
	}

	return events, nil
}

// GiftcardsUsedInOrderEvent bulk creates giftcard events
func (s *ServiceGiftcard) GiftcardsUsedInOrderEvent(transaction *gorm.DB, balanceData model.BalanceData, orderID string, user *model.User, _ interface{}) ([]*model.GiftCardEvent, *model.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	var events []*model.GiftCardEvent
	for _, item := range balanceData {
		events = append(events, &model.GiftCardEvent{
			GiftcardID: item.Giftcard.Id,
			UserID:     userID,
			Type:       model.GIFT_CARD_EVENT_TYPE_USED_IN_ORDER,
			Parameters: model.StringInterface{
				"order_id": orderID,
				"balance": model.StringInterface{
					"currency":             item.Giftcard.Currency,
					"current_balance":      item.Giftcard.CurrentBalanceAmount,
					"old_currency_balance": item.PreviousBalance,
				},
			},
		})
	}

	return s.BulkUpsertGiftcardEvents(transaction, events...)
}

func (s *ServiceGiftcard) GiftcardsBoughtEvent(transaction *gorm.DB, giftcards []*model.GiftCard, orderID string, user *model.User, _ interface{}) ([]*model.GiftCardEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	events := []*model.GiftCardEvent{}
	for _, giftCard := range giftcards {
		events = append(events, &model.GiftCardEvent{
			GiftcardID: giftCard.Id,
			UserID:     userID,
			Type:       model.GIFT_CARD_EVENT_TYPE_BOUGHT,
			Parameters: model.StringInterface{
				"order_id":    orderID,
				"expiry_date": giftCard.ExpiryDate,
			},
		})
	}

	return s.BulkUpsertGiftcardEvents(transaction, events...)
}
