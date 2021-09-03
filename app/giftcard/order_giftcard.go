package giftcard

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

// CreateOrderGiftcardRelation takes an order-giftcard relation instance then save it
func (a *ServiceGiftcard) CreateOrderGiftcardRelation(orderGiftCard *giftcard.OrderGiftCard) (*giftcard.OrderGiftCard, *model.AppError) {
	orderGiftCard, err := a.srv.Store.GiftCardOrder().Save(orderGiftCard)
	if err != nil {
		if _, ok := err.(*store.ErrInvalidInput); ok {
			return nil, model.NewAppError("CreateOrderGiftcardRelation", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderGiftcard"}, err.Error(), http.StatusBadRequest)
		}
		return nil, model.NewAppError("CreateOrderGiftcardRelation", "app.giftcard.error_creating_order_giftcard_relation.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orderGiftCard, nil
}
