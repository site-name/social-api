package giftcard

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (s *ServiceGiftcard) GiftcardEventsByOptions(options model_helper.GiftCardEventFilterOption) (model.GiftcardEventSlice, *model_helper.AppError) {
	events, err := s.srv.Store.GiftcardEvent().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("GiftcardEventsByOptions", "app.giftcard.error_finding_giftcard_events_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}

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

func (s *ServiceGiftcard) GiftcardsUsedInOrderEvent(transaction boil.ContextTransactor, balanceData model_helper.BalanceData, orderID string, user *model.User, _ any) (model.GiftcardEventSlice, *model_helper.AppError) {
	var userID *string
	if user != nil {
		userID = &user.ID
	}

	var events model.GiftcardEventSlice
	for _, item := range balanceData {
		events = append(events, &model.GiftcardEvent{
			GiftcardID: item.Giftcard.ID,
			UserID:     model_types.NullString{String: userID},
			Type:       model.GiftcardEventTypeUsedInOrder,
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

	return s.BulkUpsertGiftcardEvents(transaction, events)
}

func (s *ServiceGiftcard) GiftcardsBoughtEvent(transaction boil.ContextTransactor, giftcards model.GiftcardSlice, orderID string, user *model.User, _ any) (model.GiftcardEventSlice, *model_helper.AppError) {
	var userID *string
	if user != nil {
		userID = &user.ID
	}

	events := model.GiftcardEventSlice{}
	for _, giftCard := range giftcards {
		events = append(events, &model.GiftcardEvent{
			GiftcardID: giftCard.ID,
			UserID:     model_types.NullString{String: userID},
			Type:       model.GiftcardEventTypeBought,
			Parameters: model_types.JSONString{
				"order_id":    orderID,
				"expiry_date": giftCard.ExpiryDate,
			},
		})
	}

	return s.BulkUpsertGiftcardEvents(transaction, events)
}
