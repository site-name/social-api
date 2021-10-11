package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *ServiceProduct) ProductTypesByCheckoutToken(checkoutToken string) ([]*product_and_discount.ProductType, *model.AppError) {
	productTypes, err := a.srv.Store.ProductType().FilterProductTypesByCheckoutToken(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductTypesByCheckoutToken", "app.product.product_types_by_checkout_missing.app_error", err)
	}

	return productTypes, nil
}

// ProductTypesByProductIDs returns all product types that belong to given products
func (a *ServiceProduct) ProductTypesByProductIDs(productIDs []string) ([]*product_and_discount.ProductType, *model.AppError) {
	productTypes, err := a.srv.Store.ProductType().ProductTypesByProductIDs(productIDs)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductTypesByProductIDs", "app.product.product_types_by_product_ids.app_error", err)
	}

	return productTypes, nil
}

// ProductTypeByOption returns a product type with given option
func (s *ServiceProduct) ProductTypeByOption(options *product_and_discount.ProductTypeFilterOption) (*product_and_discount.ProductType, *model.AppError) {
	productType, err := s.srv.Store.ProductType().GetByOption(options)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductTypeByOption", "app.product_error_finding_product_type_by_option.app_error", err)
	}

	return productType, nil
}
