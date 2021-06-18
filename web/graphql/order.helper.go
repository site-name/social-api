package graphql

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

type orderError struct {
	message string
	code    gqlmodel.OrderErrorCode
}

// getOrderCountryCode try find country code for order
//
// If both shipping/billing addresses are not provided, then return default contry code from config
func (r *orderResolver) getOrderCountryCode(order *gqlmodel.Order) (string, *model.AppError) {
	addressID := order.BillingAddressID
	requireShip, appErr := r.OrderShippingIsRequired(order.ID)
	if appErr != nil {
		return "", appErr
	}
	if requireShip {
		addressID = order.ShippingAddressID
	}
	if addressID == nil {
		return *r.Config().LocalizationSettings.DefaultCountryCode, nil
	}

	address, appErr := r.GetAddressById(*addressID)
	if appErr != nil {
		return "", appErr
	}

	return address.Country, nil
}

func validateBillingAddress(order *gqlmodel.Order) *orderError {
	if order.BillingAddressID == nil {
		return &orderError{
			message: "graphql.order.invalid_billing_address.app_error",
			code:    gqlmodel.OrderErrorCodeBillingAddressNotSet,
		}
	}
	return nil
}

// func (r *orderResolver) ValidateDraftOrder(order *gqlmodel.Order) error {

// }
