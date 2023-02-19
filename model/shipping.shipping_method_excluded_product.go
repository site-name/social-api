package model

import "github.com/Masterminds/squirrel"

// ShippingMethodExcludedProduct is relation model for shipping methods and products
type ShippingMethodExcludedProduct struct {
	Id               string `json:"id"`
	ShippingMethodID string `json:"shipping_method_id"`
	ProductID        string `json:"product_id"`

	product *Product `db:"-"`
}

type ShippingMethodExcludedProductFilterOptions struct {
	ShippingMethodID squirrel.Sqlizer
	ProductID        squirrel.Sqlizer

	SelectRelatedProduct bool
}

func (s *ShippingMethodExcludedProduct) GetProduct() *Product {
	return s.product
}

func (s *ShippingMethodExcludedProduct) SetProduct(p *Product) {
	s.product = p
}

func (s *ShippingMethodExcludedProduct) DeepCopy() *ShippingMethodExcludedProduct {
	res := *s
	if s.product != nil {
		res.product = s.product.DeepCopy()
	}
	return &res
}

func (s *ShippingMethodExcludedProduct) PreSave() {
	if !IsValidId(s.Id) {
		s.Id = NewId()
	}
}

func (s *ShippingMethodExcludedProduct) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"shipping_method_excluded_product.is_valid.%s.app_error",
		"shipping_method_excluded_product_id=",
		"ShippingMethodExcludedProduct.IsValid",
	)

	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.ShippingMethodID) {
		return outer("shipping_method_id", &s.Id)
	}
	if !IsValidId(s.Id) {
		return outer("product_id", &s.Id)
	}

	return nil
}
