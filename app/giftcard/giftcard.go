/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package giftcard

import (
	"net/http"
	"time"

	"github.com/mattermost/gorp"
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
	app.RegisterGiftcardService(func(s *app.Server) (sub_app_iface.GiftcardService, error) {
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
	return a.GiftcardsByOption(nil, &giftcard.GiftCardFilterOption{
		CheckoutToken: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkoutToken,
			},
		},
	})
}

// PromoCodeIsGiftCard checks whether there is giftcard with given code
func (a *ServiceGiftcard) PromoCodeIsGiftCard(code string) (bool, *model.AppError) {
	giftcards, appErr := a.GiftcardsByOption(nil, &giftcard.GiftCardFilterOption{
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
func (a *ServiceGiftcard) GiftcardsByOption(transaction *gorp.Transaction, option *giftcard.GiftCardFilterOption) ([]*giftcard.GiftCard, *model.AppError) {
	giftcards, err := a.srv.Store.GiftCard().FilterByOption(transaction, option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		errMessage = err.Error()
		statusCode = http.StatusInternalServerError
	} else if len(giftcards) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("GiftcardsByOption", "app.giftcard.error_finding_giftcards_by_option.app_error", nil, errMessage, statusCode)
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

// ActiveGiftcards finds giftcards wich have `ExpiryDate` are either NULL OR >= given date
func (s *ServiceGiftcard) ActiveGiftcards(date *time.Time) ([]*giftcard.GiftCard, *model.AppError) {
	return s.GiftcardsByOption(nil, &giftcard.GiftCardFilterOption{
		ExpiryDate: &model.TimeFilter{
			Or: &model.TimeOption{
				GtE:               date,
				NULL:              model.NewBool(true),
				CompareStartOfDay: true,
			},
		},
		IsActive: model.NewBool(true),
	})
}
