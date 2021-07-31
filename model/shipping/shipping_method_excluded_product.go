package shipping

import "github.com/sitename/sitename/model"

// ShippingMethodExcludedProduct is relation model for shipping methods and products
type ShippingMethodExcludedProduct struct {
	Id               string `json:"id"`
	ShippingMethodID string `json:"shipping_method_id"`
	ProductID        string `json:"product_id"`
}

func (s *ShippingMethodExcludedProduct) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
}

func (s *ShippingMethodExcludedProduct) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shipping_method_excluded_product.is_valid.%s.app_error",
		"shipping_method_excluded_product_id=",
		"ShippingMethodExcludedProduct.IsValid",
	)

	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.ShippingMethodID) {
		return outer("shipping_method_id", &s.Id)
	}
	if !model.IsValidId(s.Id) {
		return outer("product_id", &s.Id)
	}

	return nil
}
