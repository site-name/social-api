package giftcard

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type AppGiftcard struct {
	app.AppIface
}

func NewAppGiftcard(a app.AppIface) app.GiftCardInterface {
	return &AppGiftcard{AppIface: a}
}

func (agc *AppGiftcard) GetAllByUserId(userID string) ([]*giftcard.GiftCard, *model.AppError) {
	if giftcards, err := agc.Srv().Store.GiftCard().GetAllByUserId(userID); err != nil {
		var nfErr *store.ErrNotFound
		statusCode := http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetAllByUserId", "app.giftcard.giftcards_by_user.app_error", nil, err.Error(), statusCode)
	} else {
		return giftcards, nil
	}
}

func (agc *AppGiftcard) GetAll() ([]*giftcard.GiftCard, *model.AppError) {
	if giftcards, err := agc.Srv().Store.GiftCard().GetAllByUserId(""); err != nil {
		var nfErr *store.ErrNotFound
		statusCode := http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetAllByUserId", "app.giftcard.giftcards_by_user.app_error", nil, err.Error(), statusCode)
	} else {
		return giftcards, nil
	}
}
