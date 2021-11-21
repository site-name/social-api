package gqlmodel

// ------------------ original implementation -----------------------

// type Warehouse struct {
// 	ID                    string                             `json:"id"`
// 	Name                  string                             `json:"name"`
// 	Slug                  string                             `json:"slug"`
// 	ShippingZones         *ShippingZoneCountableConnection   `json:"shippingZones"`
// 	Address               *Address                           `json:"address"`
// 	Email                 string                             `json:"email"`
// 	IsPrivate             bool                               `json:"isPrivate"`
// 	PrivateMetadata       []*MetadataItem                    `json:"privateMetadata"`
// 	Metadata              []*MetadataItem                    `json:"metadata"`
// 	ClickAndCollectOption WarehouseClickAndCollectOptionEnum `json:"clickAndCollectOption"`
// }

// func (Warehouse) IsDeliveryMethod()     {}
// func (Warehouse) IsNode()               {}
// func (Warehouse) IsObjectWithMetadata() {}

type Warehouse struct {
	ID              string                           `json:"id"`
	Name            string                           `json:"name"`
	Slug            string                           `json:"slug"`
	CompanyName     string                           `json:"companyName"`
	ShippingZones   *ShippingZoneCountableConnection `json:"shippingZones"`
	AddressID       *string                          `json:"address"`
	Email           string                           `json:"email"`
	PrivateMetadata []*MetadataItem                  `json:"privateMetadata"`
	Metadata        []*MetadataItem                  `json:"metadata"`
}

func (Warehouse) IsDeliveryMethod()     {}
func (Warehouse) IsNode()               {}
func (Warehouse) IsObjectWithMetadata() {}
