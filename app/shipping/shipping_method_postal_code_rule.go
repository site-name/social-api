package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (s *ServiceShipping) ShippingMethodPostalCodeRulesByOptions(options model_helper.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, *model_helper.AppError) {
	rules, err := s.srv.Store.ShippingMethodPostalCodeRule().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("ShippingMethodPostalCodeRulesByOptions", "aap.shipping.shipping_method_postal_code_rules_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return rules, nil
}

func (s *ServiceShipping) CreateShippingMethodPostalCodeRules(transaction boil.ContextTransactor, rules model.ShippingMethodPostalCodeRules) (model.ShippingMethodPostalCodeRules, *model_helper.AppError) {
	rules, err := s.srv.Store.ShippingMethodPostalCodeRule().Save(transaction, rules)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("CreateShippingMethodPostalCodeRules", "app.shipping.save_shipping_method_postal_code_rules.app_error", nil, err.Error(), statusCode)
	}

	return rules, nil
}
