package api

type OrderEvent struct {
	ID                    string                       `json:"id"`
	Date                  *DateTime                    `json:"date"`
	Type                  *OrderEventsEnum             `json:"type"`
	User                  *User                        `json:"user"`
	Message               *string                      `json:"message"`
	Email                 *string                      `json:"email"`
	EmailType             *OrderEventsEmailsEnum       `json:"emailType"`
	Amount                *float64                     `json:"amount"`
	PaymentID             *string                      `json:"paymentId"`
	PaymentGateway        *string                      `json:"paymentGateway"`
	Quantity              *int32                       `json:"quantity"`
	ComposedID            *string                      `json:"composedId"`
	OrderNumber           *string                      `json:"orderNumber"`
	InvoiceNumber         *string                      `json:"invoiceNumber"`
	OversoldItems         []string                     `json:"oversoldItems"`
	Lines                 []*OrderEventOrderLineObject `json:"lines"`
	FulfilledItems        []*FulfillmentLine           `json:"fulfilledItems"`
	Warehouse             *Warehouse                   `json:"warehouse"`
	TransactionReference  *string                      `json:"transactionReference"`
	ShippingCostsIncluded *bool                        `json:"shippingCostsIncluded"`
	RelatedOrder          *Order                       `json:"relatedOrder"`
	Discount              *OrderEventDiscountObject    `json:"discount"`
}

func (o *OrderEvent) A() {

}
