package api

import "context"

type Fulfillment struct {
	ID               string            `json:"id"`
	FulfillmentOrder int32             `json:"fulfillmentOrder"`
	Status           FulfillmentStatus `json:"status"`
	TrackingNumber   string            `json:"trackingNumber"`
	Created          DateTime          `json:"created"`
	PrivateMetadata  []*MetadataItem   `json:"privateMetadata"`
	Metadata         []*MetadataItem   `json:"metadata"`

	// Lines            []*FulfillmentLine `json:"lines"`
	// StatusDisplay    *string            `json:"statusDisplay"`
	// Warehouse        *Warehouse         `json:"warehouse"`
}

func (f *Fulfillment) Lines(ctx context.Context) ([]*FulfillmentLine, error) {
	panic("not implemented")
}

func (f *Fulfillment) StatusDisplay(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (f *Fulfillment) Warehouse(ctx context.Context) (*Warehouse, error) {
	panic("not implemented")
}

// ------------

type FulfillmentLine struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`
	// OrderLine *OrderLine `json:"orderLine"`
}

func (f *FulfillmentLine) OrderLine(ctx context.Context) (*OrderLine, error) {
	panic("not implemented")
}
