/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package giftcard

import (
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
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

func (a *ServiceGiftcard) GetGiftCard(id string) (*model.GiftCard, *model.AppError) {
	giftcard, err := a.srv.Store.GiftCard().GetById(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetGiftCard", "app.giftcard.giftcard_missing.app_error", nil, err.Error(), statusCode)
	}

	return giftcard, nil
}

func (a *ServiceGiftcard) GiftcardsByCheckout(checkoutToken string) ([]*model.GiftCard, *model.AppError) {
	return a.GiftcardsByOption(nil, &model.GiftCardFilterOption{
		CheckoutToken: squirrel.Eq{store.GiftcardCheckoutTableName + ".CheckoutID": checkoutToken},
	})
}

// PromoCodeIsGiftCard checks whether there is giftcard with given code
func (a *ServiceGiftcard) PromoCodeIsGiftCard(code string) (bool, *model.AppError) {
	giftcards, appErr := a.GiftcardsByOption(nil, &model.GiftCardFilterOption{
		Code: squirrel.Eq{store.GiftcardTableName + ".Code": code},
	})
	if appErr != nil {
		return false, appErr
	}

	return len(giftcards) > 0, nil
}

// GiftcardsByOption finds a list of giftcards with given option
func (a *ServiceGiftcard) GiftcardsByOption(transaction store_iface.SqlxTxExecutor, option *model.GiftCardFilterOption) ([]*model.GiftCard, *model.AppError) {
	giftcards, err := a.srv.Store.GiftCard().FilterByOption(transaction, option)
	if err != nil {
		return nil, model.NewAppError("GiftcardsByOption", "app.giftcard.error_finding_giftcards_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return giftcards, nil
}

// UpsertGiftcards depends on given giftcard's Id to decide saves or updates it
func (a *ServiceGiftcard) UpsertGiftcards(transaction store_iface.SqlxTxExecutor, giftcards ...*model.GiftCard) ([]*model.GiftCard, *model.AppError) {
	giftcards, err := a.srv.Store.GiftCard().BulkUpsert(transaction, giftcards...)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError

		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusInternalServerError
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("UpsertGiftcards", "app.giftcard.error_upserting_giftcards.app_error", nil, err.Error(), statusCode)
	}

	return giftcards, nil
}

// ActiveGiftcards finds giftcards wich have `ExpiryDate` are either NULL OR >= given date
func (s *ServiceGiftcard) ActiveGiftcards(date time.Time) ([]*model.GiftCard, *model.AppError) {
	return s.GiftcardsByOption(nil, &model.GiftCardFilterOption{
		ExpiryDate: squirrel.Or{
			squirrel.Eq{store.GiftcardTableName + ".ExpiryDate": nil},
			squirrel.GtOrEq{store.GiftcardTableName + ".ExpiryDate": util.StartOfDay(date)},
		},
		IsActive: model.NewPrimitive(true),
	})
}
