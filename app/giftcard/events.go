package giftcard

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gorm.io/gorm"
)

// GiftcardEventsByOptions returns a list of giftcard events filtered using given options
func (s *ServiceGiftcard) GiftcardEventsByOptions(options model_helper.GiftCardEventFilterOption) (model.GiftcardEventSlice, *model_helper.AppError) {
	events, err := s.srv.Store.GiftcardEvent().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("GiftcardEventsByOptions", "app.giftcard.error_finding_giftcard_events_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}

// BulkUpsertGiftcardEvents tells store to upsert given giftcard events into database then returns them
func (s *ServiceGiftcard) BulkUpsertGiftcardEvents(transaction boil.ContextTransactor, events model.GiftcardEventSlice) (model.GiftcardEventSlice, *model_helper.AppError) {
	events, err := s.srv.Store.GiftcardEvent().Upsert(transaction, events)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("BulkUpsertGiftcardEvents", "app.giftcard.error_upserting_giftcard_events.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}

// GiftcardsUsedInOrderEvent bulk creates giftcard events
func (s *ServiceGiftcard) GiftcardsUsedInOrderEvent(transaction *gorm.DB, balanceData model.BalanceData, orderID string, user *model.User, _ any) (model.GiftcardEventSlice, *model_helper.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	var events model.GiftcardEventSlice
	for _, item := range balanceData {
		events = append(events, &model.GiftCardEvent{
			GiftcardID: item.Giftcard.Id,
			UserID:     userID,
			Type:       model.GIFT_CARD_EVENT_TYPE_USED_IN_ORDER,
			Parameters: model_types.JSONString{
				"order_id": orderID,
				"balance": model_types.JSONString{
					"currency":             item.Giftcard.Currency,
					"current_balance":      item.Giftcard.CurrentBalanceAmount,
					"old_currency_balance": item.PreviousBalance,
				},
			},
		})
	}

	return s.BulkUpsertGiftcardEvents(transaction, events...)
}

func (s *ServiceGiftcard) GiftcardsBoughtEvent(transaction *gorm.DB, giftcards []*model.GiftCard, orderID string, user *model.User, _ any) (model.GiftcardEventSlice, *model_helper.AppError) {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	events := model.GiftcardEventSlice{}
	for _, giftCard := range giftcards {
		events = append(events, &model.GiftCardEvent{
			GiftcardID: giftCard.Id,
			UserID:     userID,
			Type:       model.GIFT_CARD_EVENT_TYPE_BOUGHT,
			Parameters: model_types.JSONString{
				"order_id":    orderID,
				"expiry_date": giftCard.ExpiryDate,
			},
		})
	}

	return s.BulkUpsertGiftcardEvents(transaction, events...)
}
