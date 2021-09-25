package checkout

import (
	"errors"

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
	Variant        product_and_discount.ProductVariant
	ChannelListing product_and_discount.ProductVariantChannelListing
	Product        product_and_discount.Product
	ProductType    product_and_discount.ProductType
	Collections    []*product_and_discount.Collection
}

type CheckoutLineInfos []*CheckoutLineInfo

func (c CheckoutLineInfos) CheckoutLines() CheckoutLines {
	var res CheckoutLines
	for _, item := range c {
		if item != nil {
			res = append(res, &item.Line)
		}
	}

	return res
}

func (c CheckoutLineInfos) Products() product_and_discount.Products {
	res := product_and_discount.Products{}
	for _, item := range c {
		res = append(res, &item.Product)
	}

	return res
}

func (c CheckoutLineInfos) FilterNils() CheckoutLineInfos {
	res := CheckoutLineInfos{}
	for _, item := range c {
		if item != nil {
			res = append(res, item)
		}
	}

	return res
}

// CheckoutInfo contains information of a checkout
type CheckoutInfo struct {
	Checkout                      Checkout
	User                          *account.User
	Channel                       channel.Channel
	BillingAddress                *account.Address
	ShippingAddress               *account.Address
	DeliveryMethodInfo            DeliveryMethodBaseInterface
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

// GetCustomerEmail returns either current checkout info's user's email or checkout's email
func (c *CheckoutInfo) GetCustomerEmail() string {
	if c.User != nil {
		return c.User.Email
	}
	return c.Checkout.Email
}

type DeliveryMethodBaseInterface interface {
	WarehousePK() string
	IsLocalCollectionPoint() bool
	DeliveryMethodName() model.StringMap
	GetWarehouseFilterLookup() model.StringInterface // the returning map contains nothing or {"warehouse_id": <an UUID>}
	IsValidDeliveryMethod() bool
	IsMethodInValidMethods(checkoutInfo *CheckoutInfo) bool
	UpdateChannelListings(checkoutInfo *CheckoutInfo) error

	GetDeliveryMethod() interface{} // GetDeliveryMethod returns an interface{}, can be either *ShippingMethod or *Warehouse.
	GetShippingAddress() *account.Address
	GetOrderKey() string

	String() string
	Self() interface{}
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

func (d *DeliveryMethodBase) Self() interface{} {
	return d
}

func (d *DeliveryMethodBase) String() string {
	return "DeliveryMethodBase"
}

func (d *DeliveryMethodBase) WarehousePK() string {
	return ""
}

func (d *DeliveryMethodBase) IsLocalCollectionPoint() bool {
	return false
}

func (d *DeliveryMethodBase) DeliveryMethodName() model.StringMap {
	return model.StringMap{"shipping_method_name": ""}
}

func (d *DeliveryMethodBase) GetWarehouseFilterLookup() model.StringInterface {
	return model.StringInterface{}
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
func (d *DeliveryMethodBase) GetShippingAddress() *account.Address {
	return d.ShippingAddress
}
func (d *DeliveryMethodBase) GetOrderKey() string {
	return d.OrderKey
}

// ShippingMethodInfo should not be modified after initializing
type ShippingMethodInfo struct {
	DeliveryMethodBase
	DeliveryMethod  shipping.ShippingMethod
	ShippingAddress *account.Address // can be nil
	OrderKey        string           // default to "shipping_method"
}

func (d *ShippingMethodInfo) Self() interface{} {
	return d
}

func (d *ShippingMethodInfo) String() string {
	return "ShippingMethodInfo"
}

func (s *ShippingMethodInfo) DeliveryMethodName() model.StringMap {
	return model.StringMap{"shipping_method_name": s.DeliveryMethod.String()}
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

var ErrorNotUsable = errors.New("this method is not usable, please use a method with same name and is a instance method of service checkout")

func (d *ShippingMethodInfo) UpdateChannelListings(_ *CheckoutInfo) error {
	return ErrorNotUsable
}

// CollectionPointInfo should not be modified after initializing
type CollectionPointInfo struct {
	DeliveryMethodBase
	DeliveryMethod  warehouse.WareHouse
	ShippingAddress *account.Address
	OrderKey        string // default to "collection_point"
}

func (d *CollectionPointInfo) Self() interface{} {
	return d
}

func (d *CollectionPointInfo) String() string {
	return "CollectionPointInfo"
}

func (c *CollectionPointInfo) WarehousePK() string {
	return c.DeliveryMethod.Id
}

func (c *CollectionPointInfo) IsLocalCollectionPoint() bool {
	return c.DeliveryMethod.ClickAndCollectOption == warehouse.LOCAL_STOCK
}

func (c *CollectionPointInfo) DeliveryMethodName() model.StringMap {
	return model.StringMap{"collection_point_name": c.DeliveryMethod.String()}
}

func (c *CollectionPointInfo) GetWarehouseFilterLookup() model.StringInterface {
	if c.IsLocalCollectionPoint() {
		return model.StringInterface{"warehouse_id": c.DeliveryMethod.Id}
	}
	return make(model.StringInterface)
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
