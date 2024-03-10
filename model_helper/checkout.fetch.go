package model_helper

import (
	"errors"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
)

// CheckoutLineInfo contains information of a checkout line
type CheckoutLineInfo struct {
	Line           model.CheckoutLine
	Variant        model.ProductVariant
	ChannelListing model.ProductVariantChannelListing
	Product        model.Product
	Collections    model.CollectionSlice
	// ProductType    model.ProductType
}

// CheckoutLineInfos is a slice contains checkout line info(s)
type CheckoutLineInfos []*CheckoutLineInfo

// CheckoutLines returns a list of checkout lines
func (cs CheckoutLineInfos) CheckoutLines() model.CheckoutLineSlice {
	return lo.Map(cs, func(c *CheckoutLineInfo, _ int) *model.CheckoutLine { return &c.Line })
}

// Products returns a list of products from current checkout line infos
func (cs CheckoutLineInfos) Products() model.ProductSlice {
	return lo.Map(cs, func(c *CheckoutLineInfo, _ int) *model.Product { return &c.Product })
}

func (cs CheckoutLineInfos) ProductVariants() model.ProductVariantSlice {
	return lo.Map(cs, func(c *CheckoutLineInfo, _ int) *model.ProductVariant { return &c.Variant })
}

// FilterNils returns a list of non-nil checkout line info(s)
func (cs CheckoutLineInfos) FilterNils() CheckoutLineInfos {
	return lo.Filter(cs, func(c *CheckoutLineInfo, _ int) bool { return c != nil })
}

// CheckoutInfo contains information of a checkout
type CheckoutInfo struct {
	Checkout                      model.Checkout
	User                          *model.User
	Channel                       model.Channel
	BillingAddress                *model.Address
	ShippingAddress               *model.Address
	DeliveryMethodInfo            DeliveryMethodBaseInterface
	ValidShippingMethods          model.ShippingMethodSlice
	ValidPickupPoints             model.WarehouseSlice
	ShippingMethodChannelListings *model.ShippingMethodChannelListing
}

// ValidDeliveryMethods returns a slice of interfaces.
//
// NOTE: These interfaces can be *Warehouse or *ShippingMethod
func (c CheckoutInfo) ValidDeliveryMethods() []any {
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
func (c CheckoutInfo) GetCountry() model.CountryCode {
	addr := c.ShippingAddress
	if addr == nil {
		addr = c.BillingAddress
	}

	if addr == nil || addr.Country.IsValid() != nil {
		return c.Checkout.Country
	}

	return addr.Country
}

// GetCustomerEmail returns either current checkout info's user's email or checkout's email
func (c CheckoutInfo) GetCustomerEmail() string {
	if c.User != nil {
		return c.User.Email
	}
	return c.Checkout.Email
}

// DeliveryMethodBaseInterface
type DeliveryMethodBaseInterface interface {
	WarehousePK() string
	IsLocalCollectionPoint() bool
	DeliveryMethodName() map[string]string
	GetWarehouseFilterLookup() map[string]any // the returning map contains nothing or {"warehouse_id": <an UUID>}
	IsValidDeliveryMethod() bool
	IsMethodInValidMethods(checkoutInfo CheckoutInfo) bool
	UpdateChannelListings(checkoutInfo *CheckoutInfo) error

	GetDeliveryMethod() any // GetDeliveryMethod returns an any, can be either *ShippingMethod or *Warehouse.
	GetShippingAddress() *model.Address
	GetOrderKey() string

	String() string
	Self() any // Self returns the current object which implements DeliveryMethodBaseInterface
	// DeliveryMethodOrderField() map[string]any
}

// check if some struct types satisfy DeliveryMethodBaseInterface
var (
	_ DeliveryMethodBaseInterface = (*DeliveryMethodBase)(nil)
	_ DeliveryMethodBaseInterface = (*ShippingMethodInfo)(nil)
	_ DeliveryMethodBaseInterface = (*CollectionPointInfo)(nil)
)

// DeliveryMethodBase should not be modified after initialized
type DeliveryMethodBase struct {
	DeliveryMethod  any            // either *ShippingMethodData or *Warehouse. Can be nil
	ShippingAddress *model.Address // can be nil
}

func (d DeliveryMethodBase) Self() any {
	return d
}

func (d DeliveryMethodBase) String() string {
	return "DeliveryMethodBase"
}

func (d DeliveryMethodBase) WarehousePK() string {
	return ""
}

// func (d *DeliveryMethodBase) DeliveryMethodOrderField() map[string]any {
// 	return map[string]any{
// 		"shipping_method": d.DeliveryMethod,
// 	}
// }

func (d DeliveryMethodBase) IsLocalCollectionPoint() bool {
	return false
}

func (d DeliveryMethodBase) DeliveryMethodName() map[string]string {
	return map[string]string{"shipping_method_name": ""}
}

func (d DeliveryMethodBase) GetWarehouseFilterLookup() map[string]any {
	return map[string]any{}
}

func (d DeliveryMethodBase) IsValidDeliveryMethod() bool {
	return false
}

func (d DeliveryMethodBase) IsMethodInValidMethods(checkoutInfo CheckoutInfo) bool {
	return false
}

func (d DeliveryMethodBase) UpdateChannelListings(checkoutInfo *CheckoutInfo) error {
	checkoutInfo.ShippingMethodChannelListings = nil
	return nil
}

func (d DeliveryMethodBase) GetDeliveryMethod() any {
	return d.DeliveryMethod
}

func (d DeliveryMethodBase) GetShippingAddress() *model.Address {
	return d.ShippingAddress
}

func (d DeliveryMethodBase) GetOrderKey() string {
	return "shipping_method"
}

// ShippingMethodInfo should not be modified after initializing
type ShippingMethodInfo struct {
	DeliveryMethodBase
	DeliveryMethod  model.ShippingMethod
	ShippingAddress *model.Address // can be nil
}

func (s ShippingMethodInfo) Self() any {
	return s
}

func (_ ShippingMethodInfo) String() string {
	return "ShippingMethodInfo"
}

func (s ShippingMethodInfo) DeliveryMethodName() map[string]string {
	return map[string]string{"shipping_method_name": s.DeliveryMethod.Name}
}

// func (s *ShippingMethodInfo) DeliveryMethodOrderField() map[string]any {
// 	if !s.DeliveryMethod.IsExternal() {
// 		return map[string]any{
// 			"shipping_method_id": s.DeliveryMethod.Id,
// 		}
// 	}

// 	return map[string]any{}
// }

func (s ShippingMethodInfo) IsValidDeliveryMethod() bool {
	return s.ShippingAddress != nil
}

func (s ShippingMethodInfo) IsMethodInValidMethods(checkoutInfo CheckoutInfo) bool {
	validDeliveryMethods := checkoutInfo.ValidDeliveryMethods()
	if len(validDeliveryMethods) == 0 {
		return false
	}

	for _, item := range validDeliveryMethods {
		if shippingMethod, ok := item.(*model.ShippingMethod); ok && shippingMethod.ID == s.DeliveryMethod.ID {
			return true
		}
	}

	return false
}

// ErrorNotUsable
var ErrorNotUsable = errors.New("this method is not usable, please use a method with same name and is a instance method of service checkout")

func (_ ShippingMethodInfo) UpdateChannelListings(_ *CheckoutInfo) error {
	return ErrorNotUsable
}

// CollectionPointInfo should not be modified after initializing
type CollectionPointInfo struct {
	DeliveryMethodBase
	DeliveryMethod  model.Warehouse
	ShippingAddress *model.Address
}

func (d CollectionPointInfo) Self() any {
	return d
}

func (_ CollectionPointInfo) String() string {
	return "CollectionPointInfo"
}

// func (d *CollectionPointInfo) DeliveryMethodOrderField() map[string]any {
// 	return map[string]any{
// 		"collection_point": d.DeliveryMethod,
// 	}
// }

func (c CollectionPointInfo) WarehousePK() string {
	return c.DeliveryMethod.ID
}

func (c CollectionPointInfo) IsLocalCollectionPoint() bool {
	return c.DeliveryMethod.ClickAndCollectOption == model.WarehouseClickAndCollectOptionLocal
}

func (c CollectionPointInfo) DeliveryMethodName() map[string]string {
	return map[string]string{"collection_point_name": c.DeliveryMethod.Name}
}

func (c CollectionPointInfo) GetWarehouseFilterLookup() map[string]any {
	if c.IsLocalCollectionPoint() {
		return map[string]any{"warehouse_id": c.DeliveryMethod.ID}
	}
	return make(map[string]any)
}

func (c CollectionPointInfo) IsValidDeliveryMethod() bool {
	return c.ShippingAddress != nil && model_types.PrimitiveIsNotNilAndEqual(c.DeliveryMethod.AddressID.String, c.ShippingAddress.ID)
}

func (c CollectionPointInfo) IsMethodInValidMethods(checkoutInfo CheckoutInfo) bool {
	if len(checkoutInfo.ValidDeliveryMethods()) == 0 {
		return false
	}

	for _, method := range checkoutInfo.ValidDeliveryMethods() {
		if wareHouse, ok := method.(*model.Warehouse); ok && wareHouse.ID == c.DeliveryMethod.ID {
			return true
		}
	}

	return false
}

func (_ CollectionPointInfo) GetOrderKey() string {
	return "collection_point"
}
