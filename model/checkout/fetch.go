package checkout

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
)

// CheckoutLineInfo contains information of a checkout line
type CheckoutLineInfo struct {
	Line           CheckoutLine
	Variant        *product_and_discount.ProductVariant
	ChannelListing *product_and_discount.ProductVariantChannelListing
	Product        product_and_discount.Product
	ProductType    product_and_discount.ProductType
	Collections    []*product_and_discount.Collection
}

// CheckoutInfo contains information of a checkout
type CheckoutInfo struct {
	Checkout                      Checkout
	User                          *account.User
	Channel                       channel.Channel
	BillingAddress                *account.Address
	ShippingAddress               *account.Address
	DeliveryMethodInfo            DeliveryMethodBase
	ShippingMethod                *shipping.ShippingMethod
	ValidShippingMethods          []*shipping.ShippingMethod
	ValidPickupPoints             []*warehouse.WareHouse
	ShippingMethodChannelListings *shipping.ShippingMethodChannelListing
}

// ValidDeliveryMethods returns a slice of interfaces.
// These interfaces can be *Warehouse or *ShippingMethod
func (c *CheckoutInfo) ValidDeliveryMethods() []interface{} {
	res := []interface{}{}
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

func (c *CheckoutInfo) GetCountry() string {
	addr := c.ShippingAddress
	if addr == nil {
		addr = c.BillingAddress
	}

	if addr == nil || addr.Country == "" {
		return c.Checkout.Country
	}

	return addr.Country
}

func (c *CheckoutInfo) GetCustomerEmail() string {
	if c.User != nil {
		return c.User.Email
	}
	return c.Checkout.Email
}

type DeliveryMethodBaseInterface interface {
	DeliveryMethodName() map[string]*string
	IsValidDeliveryMethod() bool
	IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool
}

// checking if some struct types satisfy DeliveryMethodBaseInterface
var (
	_ DeliveryMethodBaseInterface = &DeliveryMethodBase{}
	_ DeliveryMethodBaseInterface = &ShippingMethodInfo{}
	_ DeliveryMethodBaseInterface = &CollectionPointInfo{}
)

// DeliveryMethodBase should not be modified after initialized
type DeliveryMethodBase struct {
	DeliveryMethod  interface{}      // either *ShippingMethod or *Warehouse. Can be nil
	ShippingAddress *account.Address // can be nil
	OrderKey        string           // default to "shipping_method"
}

func (d *DeliveryMethodBase) WarehousePK() string

func (d *DeliveryMethodBase) IsLocalCollectionPoint() bool {
	return false
}

func (d *DeliveryMethodBase) DeliveryMethodName() map[string]*string {
	return map[string]*string{"shipping_method_name": nil}
}

func (d *DeliveryMethodBase) GetWarehouseFilterLookup() map[string]interface{} {
	return map[string]interface{}{}
}

func (d *DeliveryMethodBase) IsValidDeliveryMethod() bool {
	return false
}

func (d *DeliveryMethodBase) IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool {
	return false
}

func (d *DeliveryMethodBase) UpdateChannelListings(checkoutInfo *CheckoutInfo) {
	checkoutInfo.ShippingMethodChannelListings = nil
}

// ShippingMethodInfo should not be modified after initializing
type ShippingMethodInfo struct {
	// DeliveryMethodBase
	DeliveryMethod  shipping.ShippingMethod
	ShippingAddress *account.Address // can be nil
	OrderKey        string           // default to "shipping_method"
}

func (s *ShippingMethodInfo) DeliveryMethodName() map[string]*string {
	return map[string]*string{
		"shipping_method_name": model.NewString(s.DeliveryMethod.String()),
	}
}

func (s *ShippingMethodInfo) IsValidDeliveryMethod() bool {
	return s.ShippingAddress != nil
}

func (s *ShippingMethodInfo) IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool {
	validDeliveryMethods := checkoutInfo.ValidDeliveryMethods()
	if len(validDeliveryMethods) == 0 {
		return false
	}

	for _, item := range validDeliveryMethods {
		if shippingMethod, ok := item.(*shipping.ShippingMethod); ok && shippingMethod.Id == s.DeliveryMethod.Id {
			return true
		}
	}

	return false
}

// CollectionPointInfo should not be modified after initializing
type CollectionPointInfo struct {
	DeliveryMethod  warehouse.WareHouse
	ShippingAddress *account.Address
	OrderKey        string // default to "collection_point"
}

func (c *CollectionPointInfo) WarehousePK() string {
	return c.DeliveryMethod.Id
}

func (c *CollectionPointInfo) IsLocalCollectionPoint() bool {
	return c.DeliveryMethod.ClickAndCollectOption == warehouse.LOCAL_STOCK
}

func (c *CollectionPointInfo) DeliveryMethodName() map[string]*string {
	return map[string]*string{"collection_point_name": model.NewString(c.DeliveryMethod.String())}
}

func (c *CollectionPointInfo) GetWarehouseFilterLookup() map[string]interface{} {
	if c.IsLocalCollectionPoint() {
		return map[string]interface{}{"warehouse_id": c.DeliveryMethod.Id}
	}
	return make(map[string]interface{})
}

func (c *CollectionPointInfo) IsValidDeliveryMethod() bool {
	return c.ShippingAddress != nil && c.DeliveryMethod.AddressID != nil && c.ShippingAddress.Id == *c.DeliveryMethod.AddressID
}

func (c *CollectionPointInfo) IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool {
	if len(checkoutInfo.ValidDeliveryMethods()) == 0 {
		return false
	}

	for _, method := range checkoutInfo.ValidDeliveryMethods() {
		if wareHouse, ok := method.(*warehouse.WareHouse); ok && wareHouse.Id == c.DeliveryMethod.Id {
			return true
		}
	}

	return false
}
