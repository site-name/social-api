package model

import (
	"errors"

	"github.com/samber/lo"
)

// CheckoutLineInfo contains information of a checkout line
type CheckoutLineInfo struct {
	Line           CheckoutLine
	Variant        ProductVariant
	ChannelListing ProductVariantChannelListing
	Product        Product
	ProductType    ProductType
	Collections    []*Collection
}

// CheckoutLineInfos is a slice contains checkout line info(s)
type CheckoutLineInfos []*CheckoutLineInfo

// CheckoutLines returns a list of checkout lines
func (cs CheckoutLineInfos) CheckoutLines() CheckoutLines {
	return lo.Map(cs, func(c *CheckoutLineInfo, _ int) *CheckoutLine { return &c.Line })
}

// Products returns a list of products from current checkout line infos
func (cs CheckoutLineInfos) Products() Products {
	return lo.Map(cs, func(c *CheckoutLineInfo, _ int) *Product { return &c.Product })
}

func (cs CheckoutLineInfos) ProductVariants() ProductVariants {
	return lo.Map(cs, func(c *CheckoutLineInfo, _ int) *ProductVariant { return &c.Variant })
}

// FilterNils returns a list of non-nil checkout line info(s)
func (cs CheckoutLineInfos) FilterNils() CheckoutLineInfos {
	return lo.Filter(cs, func(c *CheckoutLineInfo, _ int) bool { return c != nil })
}

// CheckoutInfo contains information of a checkout
type CheckoutInfo struct {
	Checkout                      Checkout
	User                          *User
	Channel                       Channel
	BillingAddress                *Address
	ShippingAddress               *Address
	DeliveryMethodInfo            DeliveryMethodBaseInterface
	ValidShippingMethods          []*ShippingMethod
	ValidPickupPoints             []*WareHouse
	ShippingMethodChannelListings *ShippingMethodChannelListing
}

// ValidDeliveryMethods returns a slice of interfaces.
//
// NOTE: These interfaces can be *Warehouse or *ShippingMethod
func (c *CheckoutInfo) ValidDeliveryMethods() []any {
	var res []any

	for _, item := range c.ValidShippingMethods {
		if item != nil {
			res = append(res, item)
		}
	}

	for _, item := range c.ValidPickupPoints {
		if item != nil {
			res = append(res, item)
		}
	}

	return res
}

// GetCountry
func (c *CheckoutInfo) GetCountry() CountryCode {
	addr := c.ShippingAddress
	if addr == nil {
		addr = c.BillingAddress
	}

	if addr == nil || Countries[addr.Country] == "" {
		return c.Checkout.Country
	}

	return addr.Country
}

// GetCustomerEmail returns either current checkout info's user's email or checkout's email
func (c *CheckoutInfo) GetCustomerEmail() string {
	if c.User != nil {
		return c.User.Email
	}
	return c.Checkout.Email
}

// DeliveryMethodBaseInterface
type DeliveryMethodBaseInterface interface {
	WarehousePK() string
	// DeliveryMethodOrderField() StringInterface
	IsLocalCollectionPoint() bool
	DeliveryMethodName() StringMap
	GetWarehouseFilterLookup() StringInterface // the returning map contains nothing or {"warehouse_id": <an UUID>}
	IsValidDeliveryMethod() bool
	IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool
	UpdateChannelListings(checkoutInfo *CheckoutInfo) error

	GetDeliveryMethod() interface{} // GetDeliveryMethod returns an interface{}, can be either *ShippingMethod or *Warehouse.
	GetShippingAddress() *Address
	GetOrderKey() string

	String() string
	Self() interface{} // Self returns the current object which implements DeliveryMethodBaseInterface
}

// check if some struct types satisfy DeliveryMethodBaseInterface
var (
	_ DeliveryMethodBaseInterface = (*DeliveryMethodBase)(nil)
	_ DeliveryMethodBaseInterface = (*ShippingMethodInfo)(nil)
	_ DeliveryMethodBaseInterface = (*CollectionPointInfo)(nil)
)

// DeliveryMethodBase should not be modified after initialized
type DeliveryMethodBase struct {
	DeliveryMethod  interface{} // either *ShippingMethodData or *Warehouse. Can be nil
	ShippingAddress *Address    // can be nil
}

func (d *DeliveryMethodBase) Self() interface{} {
	return d
}

func (d *DeliveryMethodBase) String() string {
	return "DeliveryMethodBase"
}

func (d *DeliveryMethodBase) WarehousePK() string {
	return ""
}

// func (d *DeliveryMethodBase) DeliveryMethodOrderField() StringInterface {
// 	return StringInterface{
// 		"shipping_method": d.DeliveryMethod,
// 	}
// }

func (d *DeliveryMethodBase) IsLocalCollectionPoint() bool {
	return false
}

func (d *DeliveryMethodBase) DeliveryMethodName() StringMap {
	return StringMap{"shipping_method_name": ""}
}

func (d *DeliveryMethodBase) GetWarehouseFilterLookup() StringInterface {
	return StringInterface{}
}

func (d *DeliveryMethodBase) IsValidDeliveryMethod() bool {
	return false
}

func (d *DeliveryMethodBase) IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool {
	return false
}

func (d *DeliveryMethodBase) UpdateChannelListings(checkoutInfo *CheckoutInfo) error {
	checkoutInfo.ShippingMethodChannelListings = nil
	return nil
}

func (d *DeliveryMethodBase) GetDeliveryMethod() interface{} {
	return d.DeliveryMethod
}

func (d *DeliveryMethodBase) GetShippingAddress() *Address {
	return d.ShippingAddress
}

func (d *DeliveryMethodBase) GetOrderKey() string {
	return "shipping_method"
}

// ShippingMethodInfo should not be modified after initializing
type ShippingMethodInfo struct {
	DeliveryMethodBase
	DeliveryMethod  ShippingMethod
	ShippingAddress *Address // can be nil
}

func (d *ShippingMethodInfo) Self() interface{} {
	return d
}

func (d *ShippingMethodInfo) String() string {
	return "ShippingMethodInfo"
}

func (s *ShippingMethodInfo) DeliveryMethodName() StringMap {
	return StringMap{"shipping_method_name": s.DeliveryMethod.Name}
}

// func (s *ShippingMethodInfo) DeliveryMethodOrderField() StringInterface {
// 	if !s.DeliveryMethod.IsExternal() {
// 		return StringInterface{
// 			"shipping_method_id": s.DeliveryMethod.Id,
// 		}
// 	}

// 	return StringInterface{}
// }

func (s *ShippingMethodInfo) IsValidDeliveryMethod() bool {
	return s.ShippingAddress != nil
}

func (s *ShippingMethodInfo) IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool {
	validDeliveryMethods := checkoutInfo.ValidDeliveryMethods()
	if len(validDeliveryMethods) == 0 {
		return false
	}

	for _, item := range validDeliveryMethods {
		if shippingMethod, ok := item.(*ShippingMethod); ok && shippingMethod.Id == s.DeliveryMethod.Id {
			return true
		}
	}

	return false
}

// ErrorNotUsable
var ErrorNotUsable = errors.New("this method is not usable, please use a method with same name and is a instance method of service checkout")

func (d *ShippingMethodInfo) UpdateChannelListings(_ *CheckoutInfo) error {
	return ErrorNotUsable
}

// CollectionPointInfo should not be modified after initializing
type CollectionPointInfo struct {
	DeliveryMethodBase
	DeliveryMethod  WareHouse
	ShippingAddress *Address
}

func (d *CollectionPointInfo) Self() interface{} {
	return d
}

func (d *CollectionPointInfo) String() string {
	return "CollectionPointInfo"
}

// func (d *CollectionPointInfo) DeliveryMethodOrderField() StringInterface {
// 	return StringInterface{
// 		"collection_point": d.DeliveryMethod,
// 	}
// }

func (c *CollectionPointInfo) WarehousePK() string {
	return c.DeliveryMethod.Id
}

func (c *CollectionPointInfo) IsLocalCollectionPoint() bool {
	return c.DeliveryMethod.ClickAndCollectOption == LOCAL_STOCK
}

func (c *CollectionPointInfo) DeliveryMethodName() StringMap {
	return StringMap{"collection_point_name": c.DeliveryMethod.String()}
}

func (c *CollectionPointInfo) GetWarehouseFilterLookup() StringInterface {
	if c.IsLocalCollectionPoint() {
		return StringInterface{"warehouse_id": c.DeliveryMethod.Id}
	}
	return make(StringInterface)
}

func (c *CollectionPointInfo) IsValidDeliveryMethod() bool {
	return c.ShippingAddress != nil && c.DeliveryMethod.AddressID != nil && c.ShippingAddress.Id == *c.DeliveryMethod.AddressID
}

func (c *CollectionPointInfo) IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool {
	if len(checkoutInfo.ValidDeliveryMethods()) == 0 {
		return false
	}

	for _, method := range checkoutInfo.ValidDeliveryMethods() {
		if wareHouse, ok := method.(*WareHouse); ok && wareHouse.Id == c.DeliveryMethod.Id {
			return true
		}
	}

	return false
}

func (c *CollectionPointInfo) GetOrderKey() string {
	return "collection_point"
}
