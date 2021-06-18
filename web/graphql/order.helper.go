package graphql

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

type orderError struct {
	message string
	code    gqlmodel.OrderErrorCode
}

func newOrderError(messageFormat string, code gqlmodel.OrderErrorCode) *orderError {
	return &orderError{
		message: messageFormat,
		code:    code,
	}
}

// getOrderCountryCode try find country code for order
//
// If both shipping/billing addresses are not provided, then return default contry code from config
func (r *orderResolver) getOrderCountryCode(order *gqlmodel.Order) (string, *model.AppError) {
	// remember to grant value to BillingAddressID, ShippingAddressID if available
	// in graphql Order model initializations
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

// validateBillingAddress validates if given order has billing address or not
func validateBillingAddress(order *gqlmodel.Order) *orderError {
	if order.BillingAddressID == nil {
		return &orderError{
			message: "graphql.order.billing_address_not_set.app_error",
			code:    gqlmodel.OrderErrorCodeBillingAddressNotSet,
		}
	}
	return nil
}

// validateShippingAddress validates if given order has shipping address or not
func validateShippingAddress(order *gqlmodel.Order) *orderError {
	if order.ShippingAddressID == nil {
		return &orderError{
			message: "graphql.order.shipping_address.not_set.app_error",
			code:    gqlmodel.OrderErrorCodeBillingAddressNotSet,
		}
	}
	return nil
}

// validateShippingMethod validates if given order has shipping method or not
func validateShippingMethod(order *gqlmodel.Order) *orderError {
	if order.ShippingMethodID == nil {
		return &orderError{
			message: "graphql.order.shipping_method_not_set.app_error",
			code:    gqlmodel.OrderErrorCodeShippingMethodRequired,
		}
	}
	return nil
}

func (r *orderResolver) validateOrderLines(order *gqlmodel.Order, orderCountryCode string) (*orderError, *model.AppError) {
	orderLines, appErr := r.GetAllOrderLinesByOrderId(order.ID)
	if appErr != nil {
		return nil, appErr
	}

	for _, line := range orderLines {
		if line.VariantID == nil {
			return newOrderError("graphql.order.order_line_no_variant.app_error", gqlmodel.OrderErrorCodeNotFound), nil
		}
		// get product variant
		variant, appErr := r.ProductVariantById(*line.VariantID)
		if appErr != nil {
			return nil, appErr
		}
		if variant.TrackInventory != nil && *variant.TrackInventory {

		}
	}
}

// ValidateDraftOrder validates
func (r *orderResolver) ValidateDraftOrder(order *gqlmodel.Order, orderCountryCode string) (*orderError, *model.AppError) {
	if err := validateBillingAddress(order); err != nil {
		return err, nil
	}

	requireShip, appErr := r.OrderShippingIsRequired(order.ID)
	if appErr != nil {
		return nil, appErr
	}

	// if order requires ship
	if requireShip {
		if err := validateShippingAddress(order); err != nil {
			return err, nil
		}
		if err := validateShippingMethod(order); err != nil {
			return err, nil
		}
	}
	// check total quantity
	if totalQuantity, appErr := r.OrderTotalQuantity(order.ID); appErr != nil {
		return nil, appErr
	} else if totalQuantity == 0 {
		return newOrderError("graphql.order.quantity_empty.app_error", gqlmodel.OrderErrorCodeRequired), nil
	}

	return nil
}
