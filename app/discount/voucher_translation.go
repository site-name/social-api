package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

// VoucherTranslationsByOption returns a list of voucher translations filtered using given option
func (s *ServiceDiscount) VoucherTranslationsByOption(option *model.VoucherTranslationFilterOption) ([]*model.VoucherTranslation, *model_helper.AppError) {
	translations, err := s.srv.Store.VoucherTranslation().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("VoucherTranslationsByOption", "app.discount.error_finding_voucher_translations_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return translations, nil
}

// GetVoucherTranslationByOption returns a voucher translation by given options
func (s *ServiceDiscount) GetVoucherTranslationByOption(option *model.VoucherTranslationFilterOption) (*model.VoucherTranslation, *model_helper.AppError) {
	translation, err := s.srv.Store.VoucherTranslation().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("GetVoucherTranslationByOption", "app.discount.error_finding_voucher_translation_by_option.app_error", nil, err.Error(), statusCode)
	}

	return translation, nil
}
