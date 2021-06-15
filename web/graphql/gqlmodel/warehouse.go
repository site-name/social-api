package gqlmodel

type Warehouse struct {
	ID              string                           `json:"id"`
	Name            string                           `json:"name"`
	Slug            string                           `json:"slug"`
	CompanyName     string                           `json:"companyName"`
	ShippingZones   *ShippingZoneCountableConnection `json:"shippingZones"`
	AddressID       *string                          `json:"address"` // *Address
	Email           string                           `json:"email"`
	PrivateMetadata []*MetadataItem                  `json:"privateMetadata"`
	Metadata        []*MetadataItem                  `json:"metadata"`
}

func (Warehouse) IsNode()               {}
func (Warehouse) IsObjectWithMetadata() {}
