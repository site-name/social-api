package giftcard

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// UpsertOrderGiftcardRelations takes an order-giftcard relation instance then save it
func (a *ServiceGiftcard) UpsertOrderGiftcardRelations(transaction store_iface.SqlxTxExecutor, orderGiftCards ...*model.OrderGiftCard) ([]*model.OrderGiftCard, *model.AppError) {
	orderGiftCards, err := a.srv.Store.GiftCardOrder().BulkUpsert(transaction, orderGiftCards...)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		if _, ok := err.(*store.ErrInvalidInput); ok {
			return nil, model.NewAppError("UpsertOrderGiftcardRelations", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderGiftcard"}, err.Error(), http.StatusBadRequest)
		}
		var statusCode int = http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("UpsertOrderGiftcardRelations", "app.giftcard.error_upserting_order_giftcard_relations.app_error", nil, err.Error(), statusCode)
	}

	return orderGiftCards, nil
}
