/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package giftcard

import (
	"net/http"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
	giftcard, err := a.srv.Store.GiftCard().GetById(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("GetGiftCard", "app.giftcard.giftcard_missing.app_error", nil, err.Error(), statusCode)
	}

	return giftcard, nil
}

func (a *ServiceGiftcard) GiftcardsByCheckout(checkoutToken string) (model.GiftcardSlice, *model_helper.AppError) {
	giftcards, appErr := a.GiftcardsByOption(model_helper.GiftcardFilterOption{
		CheckoutToken: model.GiftcardCheckoutWhere.CheckoutID.EQ(checkoutToken),
	})
	return giftcards, appErr
}

func (a *ServiceGiftcard) PromoCodeIsGiftCard(code string) (bool, *model_helper.AppError) {
	giftcards, appErr := a.GiftcardsByOption(model_helper.GiftcardFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.GiftcardWhere.Code.EQ(code),
		),
	})
	if appErr != nil {
		return false, appErr
	}

	return len(giftcards) > 0, nil
}

func (a *ServiceGiftcard) GiftcardsByOption(option model_helper.GiftcardFilterOption) (model.GiftcardSlice, *model_helper.AppError) {
	giftcards, err := a.srv.Store.GiftCard().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("GiftcardsByOption", "app.giftcard.error_finding_giftcards_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return giftcards, nil
}

func (a *ServiceGiftcard) UpsertGiftcards(transaction boil.ContextTransactor, giftcards model.GiftcardSlice) (model.GiftcardSlice, *model_helper.AppError) {
	giftcards, err := a.srv.Store.GiftCard().BulkUpsert(transaction, giftcards)
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

func (s *ServiceGiftcard) ActiveGiftcards(date time.Time) (model.GiftcardSlice, *model_helper.AppError) {
	giftcards, appErr := s.GiftcardsByOption(model_helper.GiftcardFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model_helper.Or{
				squirrel.Eq{model.GiftcardTableColumns.ExpiryDate: nil},
				squirrel.GtOrEq{model.GiftcardTableColumns.ExpiryDate: util.StartOfDay(date)},
			},
			model.GiftcardWhere.IsActive.EQ(model_types.NewNullBool(true)),
		),
	})
	return giftcards, appErr
}

func (s *ServiceGiftcard) DeleteGiftcards(transaction boil.ContextTransactor, ids []string) *model_helper.AppError {
	err := s.srv.Store.GiftCard().Delete(transaction, ids)
	if err != nil {
		return model_helper.NewAppError("DeleteGiftcards", "app.giftcard.error_deleting_giftcards.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// relations must be []*Order || []*Checkout
func (s *ServiceGiftcard) AddGiftcardRelations(transaction boil.ContextTransactor, giftcards model.GiftcardSlice, relations any) *model_helper.AppError {
	err := s.srv.Store.GiftCard().AddRelations(transaction, giftcards, relations)
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
func (s *ServiceGiftcard) RemoveGiftcardRelations(transaction boil.ContextTransactor, giftcards model.GiftcardSlice, relations any) *model_helper.AppError {
	err := s.srv.Store.GiftCard().RemoveRelations(transaction, giftcards, relations)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return model_helper.NewAppError("AddGiftcardRelations", "app.giftcard.add_relations.app_error", nil, err.Error(), statusCode)
	}

	return nil
}
