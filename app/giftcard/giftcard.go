/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package giftcard

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type ServiceGiftcard struct {
	srv *app.Server
}

func init() {
	app.RegisterGiftcardApp(func(s *app.Server) (sub_app_iface.GiftcardService, error) {
		return &ServiceGiftcard{
			srv: s,
		}, nil
	})
}

func (a *ServiceGiftcard) GetGiftCard(id string) (*giftcard.GiftCard, *model.AppError) {
	gc, err := a.srv.Store.GiftCard().GetById(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetGiftCard", "app.giftcard.giftcard_missing.app_error", err)
	}

	return gc, nil
}

func (a *ServiceGiftcard) GiftcardsByCheckout(checkoutToken string) ([]*giftcard.GiftCard, *model.AppError) {
	// validate given checkout token is valid uuid
	if !model.IsValidId(checkoutToken) {
		return nil, model.NewAppError("GiftcardsByCheckout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkoutToken"}, "", http.StatusBadRequest)
	}

	giftcardsOfCheckout, err := a.srv.Store.GiftCard().GetAllByCheckout(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GiftcardsByCheckout", "app.giftcard.giftcards_by_checkout_missing.app_error", err)
	}
	return giftcardsOfCheckout, nil
}

// PromoCodeIsGiftCard checks whether there is giftcard with given code
func (a *ServiceGiftcard) PromoCodeIsGiftCard(code string) (bool, *model.AppError) {
	giftcards, appErr := a.GiftcardsByOption(&giftcard.GiftCardFilterOption{
		Code: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: code,
			},
		},
	})

	if appErr != nil {
		return false, appErr
	}

	return len(giftcards) != 0, nil
}

// GiftcardsByOption finds a list of giftcards with given option
func (a *ServiceGiftcard) GiftcardsByOption(option *giftcard.GiftCardFilterOption) ([]*giftcard.GiftCard, *model.AppError) {
	giftcards, err := a.srv.Store.GiftCard().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GiftCardsByOption", "app.giftcard.error_finding_giftcards_by_option.app_error", err)
	}

	return giftcards, nil
}

// UpsertGiftcard depends on given giftcard's Id to decide saves or updates it
func (a *ServiceGiftcard) UpsertGiftcard(giftcard *giftcard.GiftCard) (*giftcard.GiftCard, *model.AppError) {
	giftcard, err := a.srv.Store.GiftCard().Upsert(giftcard)
	if err != nil {
		return nil, model.NewAppError("UpsertGiftcard", "app.giftcard.error_updating_giftcard.app_error", nil, err.Error(), http.StatusExpectationFailed)
	}

	return giftcard, nil
}
