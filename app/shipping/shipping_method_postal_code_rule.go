package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

func (s *ServiceShipping) ShippingMethodPostalCodeRulesByOptions(options *model.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, *model.AppError) {
	rules, err := s.srv.Store.ShippingMethodPostalCodeRule().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("ShippingMethodPostalCodeRulesByOptions", "aap.shipping.shipping_method_postal_code_rules_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return rules, nil
}

func (s *ServiceShipping) CreateShippingMethodPostalCodeRules(transaction store_iface.SqlxTxExecutor, rules model.ShippingMethodPostalCodeRules) (model.ShippingMethodPostalCodeRules, *model.AppError) {
	rules, err := s.srv.Store.ShippingMethodPostalCodeRule().Save(transaction, rules)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("CreateShippingMethodPostalCodeRules", "app.shipping.save_shipping_method_postal_code_rules.app_error", nil, err.Error(), statusCode)
	}

	return rules, nil
}
