package giftcard

import (
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

// CommonCreateGiftcardEvent is common method for creating giftcard events
func (s *ServiceGiftcard) CommonCreateGiftcardEvent(giftcardID, userID string, parameters model.StringMap, Type string) (*giftcard.GiftCardEvent, *model.AppError) {
	panic("not implemented")
}

// BulkUpsertGiftcardEvents tells store to upsert given giftcard events into database then returns them
func (s *ServiceGiftcard) BulkUpsertGiftcardEvents(transaction *gorp.Transaction, events []*giftcard.GiftCardEvent) ([]*giftcard.GiftCardEvent, *model.AppError) {
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
func (s *ServiceGiftcard) GiftcardsUsedInOrderEvent(transaction *gorp.Transaction, balanceData giftcard.BalanceData, orderID string, user *account.User, _ interface{}) ([]*giftcard.GiftCardEvent, *model.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	var events []*giftcard.GiftCardEvent
	for _, item := range balanceData {
		events = append(events, &giftcard.GiftCardEvent{
			GiftcardID: item.Giftcard.Id,
			UserID:     userID,
			Type:       giftcard.USED_IN_ORDER,
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

	return s.BulkUpsertGiftcardEvents(transaction, events)
}
