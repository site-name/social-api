package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// VoucherTranslationsByOption returns a list of voucher translations filtered using given option
func (s *ServiceDiscount) VoucherTranslationsByOption(option *product_and_discount.VoucherTranslationFilterOption) ([]*product_and_discount.VoucherTranslation, *model.AppError) {
	translations, err := s.srv.Store.VoucherTranslation().FilterByOption(option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(translations) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("VoucherTranslationsByOption", "app.discount.error_finding_voucher_translations_by_option.app_error", nil, errMessage, statusCode)
	}

	return translations, nil
}

// GetVoucherTranslationByOption returns a voucher translation by given options
func (s *ServiceDiscount) GetVoucherTranslationByOption(option *product_and_discount.VoucherTranslationFilterOption) (*product_and_discount.VoucherTranslation, *model.AppError) {
	translation, err := s.srv.Store.VoucherTranslation().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetVoucherTranslationByOption", "app.discount.error_finding_voucher_translation_by_option.app_error", err)
	}

	return translation, nil
}
