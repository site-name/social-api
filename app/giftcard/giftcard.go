package giftcard

import (
	"net/http"

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

// PromoCodeIsGiftCard checks whether there is giftcard with given code
func (a *AppGiftcard) PromoCodeIsGiftCard(code string) (bool, *model.AppError) {
	giftcards, err := a.Srv().Store.GiftCard().FilterByOption(&giftcard.GiftCardFilterOption{
		Code: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: code,
			},
		},
	})

	if err != nil {
		if _, ok := err.(*store.ErrNotFound); ok {
			return false, nil
		}
		return false, model.NewAppError("PromoCodeIsGiftCard", "app.giftcard.error_finding_giftcards_with_option", nil, err.Error(), http.StatusInternalServerError)
	}

	return len(giftcards) != 0, nil
}

// GiftcardsByOption finds a list of giftcards with given option
func (a *AppGiftcard) GiftcardsByOption(option *giftcard.GiftCardFilterOption) ([]*giftcard.GiftCard, *model.AppError) {
	giftcards, err := a.Srv().Store.GiftCard().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GiftCardsByOption", "app.giftcard.error_finding_giftcards_by_option.app_error", err)
	}

	return giftcards, nil
}

// UpdateGiftCard updates given giftcard. You must changed the giftcard's properties before giving it to me
func (a *AppGiftcard) UpdateGiftCard(giftcard *giftcard.GiftCard) (*giftcard.GiftCard, *model.AppError) {
	giftcard, err := a.Srv().Store.GiftCard().Upsert(giftcard)
	if err != nil {
		return nil, model.NewAppError("UpdateGiftCard", "app.giftcard.error_updating_giftcard.app_error", nil, err.Error(), http.StatusExpectationFailed)
	}

	return giftcard, nil
}
