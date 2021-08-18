package payment

import "github.com/sitename/sitename/model"

// PaymentTransactionFilterOpts contains options for filter payment's transactions
type PaymentTransactionFilterOpts struct {
	Id             *model.StringFilter
	PaymentID      *model.StringFilter
	Kind           *model.StringFilter
	ActionRequired *bool
	IsSuccess      *bool
}
