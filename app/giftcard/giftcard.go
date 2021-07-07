package giftcard

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type AppGiftcard struct {
	app.AppIface
}

func init() {
	app.RegisterGiftcardApp(func(a app.AppIface) sub_app_iface.GiftcardApp {
		return &AppGiftcard{a}
	})
}

func (a *AppGiftcard) GetGiftCard(id string) (*giftcard.GiftCard, *model.AppError) {
	gc, err := a.Srv().Store.GiftCard().GetById(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetGiftCard", "app.giftcard.giftcard_missing.app_error", err)
	}

	return gc, nil
}

func (a *AppGiftcard) GiftcardsByCheckout(checkoutID string) ([]*giftcard.GiftCard, *model.AppError) {
	gcs, err := a.Srv().Store.GiftCard().GetAllByCheckout(checkoutID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GiftcardsByCheckout", "app.giftcard.giftcards_by_checkout_missing.app_error", err)
	}
	return gcs, nil
}

func (a *AppGiftcard) GiftcardsByOrder(orderID string) ([]*giftcard.GiftCard, *model.AppError) {
	gcs, err := a.Srv().Store.GiftCard().GetAllByOrder(orderID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GiftcardsByOrder", "app.giftcard.giftcards_by_order_missing.app_error", err)
	}
	return gcs, nil
}
