package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func (s *ServiceShipping) ShippingMethodPostalCodeRulesByOptions(options *model.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, *model.AppError) {
	rules, err := s.srv.Store.ShippingMethodPostalCodeRule().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("ShippingMethodPostalCodeRulesByOptions", "aap.shipping.shipping_method_postal_code_rules_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return rules, nil
}
