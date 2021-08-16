package model

import "github.com/site-name/decimal"

// Exception raised when a discount is not applicable to a checkout.
//
// The error is raised if the order value is below the minimum required
// price or the order quantity is below the minimum quantity of items.
// Minimum price will be available as the `min_spent` attribute.
// Minimum quantity will be available as the `min_checkout_items_quantity` attribute.
type NotApplicable struct {
	Where                    string
	Message                  string
	MinSpent                 *decimal.Decimal
	MinCheckoutItemsQuantity int
}

func (a *NotApplicable) Error() string {
	return a.Message
}

// NewNotApplicable
func NewNotApplicable(where, message string, minSpent *decimal.Decimal, minCheckoutItemsQuantity int) *NotApplicable {
	return &NotApplicable{
		Where:                    where,
		Message:                  message,
		MinSpent:                 minSpent,
		MinCheckoutItemsQuantity: minCheckoutItemsQuantity,
	}
}
