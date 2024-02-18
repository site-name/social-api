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
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type ServiceGiftcard struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Giftcard = &ServiceGiftcard{s}
		return nil
	})
}

func (a *ServiceGiftcard) GetGiftCard(id string) (*model.Giftcard, *model_helper.AppError) {
	giftcard, err := a.srv.Store.Giftcard().GetById(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("GetGiftCard", "app.giftcard.giftcard_missing.app_error", nil, err.Error(), statusCode)
	}

	return giftcard, nil
}

func (a *ServiceGiftcard) GiftcardsByCheckout(checkoutToken string) ([]*model.Giftcard, *model_helper.AppError) {
	_, giftcards, appErr := a.GiftcardsByOption(model_helper.GiftcardFilterOption{
		CheckoutToken: squirrel.Eq{model.GiftcardCheckoutTableName + ".CheckoutID": checkoutToken},
	})
	return giftcards, appErr
}

// PromoCodeIsGiftCard checks whether there is giftcard with given code
func (a *ServiceGiftcard) PromoCodeIsGiftCard(code string) (bool, *model_helper.AppError) {
	_, giftcards, appErr := a.GiftcardsByOption(&model.GiftCardFilterOption{
		Conditions: squirrel.Eq{model.GiftcardTableName + ".Code": code},
	})
	if appErr != nil {
		return false, appErr
	}

	return len(giftcards) > 0, nil
}

// GiftcardsByOption finds a list of giftcards with given option
func (a *ServiceGiftcard) GiftcardsByOption(option model_helper.GiftcardFilterOption) (int64, []*model.Giftcard, *model_helper.AppError) {
	totalCount, giftcards, err := a.srv.Store.Giftcard().FilterByOption(option)
	if err != nil {
		return 0, nil, model_helper.NewAppError("GiftcardsByOption", "app.giftcard.error_finding_giftcards_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return totalCount, giftcards, nil
}

// UpsertGiftcards depends on given giftcard's Id to decide saves or updates it
func (a *ServiceGiftcard) UpsertGiftcards(transaction *gorm.DB, giftcards ...*model.Giftcard) ([]*model.Giftcard, *model_helper.AppError) {
	giftcards, err := a.srv.Store.Giftcard().BulkUpsert(transaction, giftcards...)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError

		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusInternalServerError
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("UpsertGiftcards", "app.giftcard.error_upserting_giftcards.app_error", nil, err.Error(), statusCode)
	}

	return giftcards, nil
}

// ActiveGiftcards finds giftcards wich have `ExpiryDate` are either NULL OR >= given date
func (s *ServiceGiftcard) ActiveGiftcards(date time.Time) ([]*model.Giftcard, *model_helper.AppError) {
	_, giftcards, appErr := s.GiftcardsByOption(&model.GiftCardFilterOption{
		Conditions: squirrel.And{
			squirrel.Or{
				squirrel.Eq{model.GiftcardTableName + ".ExpiryDate": nil},
				squirrel.GtOrEq{model.GiftcardTableName + ".ExpiryDate": util.StartOfDay(date)},
			},
			squirrel.Eq{model.GiftcardTableName + ".IsActive": true},
		},
	})
	return giftcards, appErr
}

func (s *ServiceGiftcard) DeleteGiftcards(transaction *gorm.DB, ids []string) *model_helper.AppError {
	err := s.srv.Store.Giftcard().DeleteGiftcards(transaction, ids)
	if err != nil {
		return model_helper.NewAppError("DeleteGiftcards", "app.giftcard.error_deleting_giftcards.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// relations must be []*Order || []*Checkout
func (s *ServiceGiftcard) AddGiftcardRelations(transaction *gorm.DB, giftcards model.Giftcards, relations any) *model_helper.AppError {
	err := s.srv.Store.Giftcard().AddRelations(transaction, giftcards, relations)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return model_helper.NewAppError("AddGiftcardRelations", "app.giftcard.add_relations.app_error", nil, err.Error(), statusCode)
	}

	return nil
}

// relations must be []*Order || []*Checkout
func (s *ServiceGiftcard) RemoveGiftcardRelations(transaction *gorm.DB, giftcards model.Giftcards, relations any) *model_helper.AppError {
	err := s.srv.Store.Giftcard().RemoveRelations(transaction, giftcards, relations)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return model_helper.NewAppError("AddGiftcardRelations", "app.giftcard.add_relations.app_error", nil, err.Error(), statusCode)
	}

	return nil
}
