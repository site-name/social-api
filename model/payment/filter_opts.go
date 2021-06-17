package payment

// PaymentFilterOpts contains options for filter payments
type PaymentFilterOpts struct {
	IsActive bool   `json:"is_active"`
	OrderID  string `json:"order_id"`
	PaymentTransactionFilterOpts
}

// PaymentTransactionFilterOpts contains options for filter payment's transactions
type PaymentTransactionFilterOpts struct {
	Kind           string `json:"kind"`
	ActionRequired bool   `json:"action_required"`
	IsSuccess      bool   `json:"is_success"`
}
