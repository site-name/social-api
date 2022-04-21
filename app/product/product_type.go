package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *ServiceProduct) ProductTypesByCheckoutToken(checkoutToken string) ([]*product_and_discount.ProductType, *model.AppError) {
	productTypes, err := a.srv.Store.ProductType().FilterProductTypesByCheckoutToken(checkoutToken)
	var (
		statusCode int
		errMsg     string
	)

	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(productTypes) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ProductTypesByCheckoutToken", "app.product.error_finding_product_types_by_checkout_token.app_error", nil, errMsg, statusCode)
	}

	return productTypes, nil
}

// ProductTypesByProductIDs returns all product types that belong to given products
func (a *ServiceProduct) ProductTypesByProductIDs(productIDs []string) ([]*product_and_discount.ProductType, *model.AppError) {
	productTypes, err := a.srv.Store.ProductType().ProductTypesByProductIDs(productIDs)

	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(productTypes) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ProductTypesByProductIDs", "app.product.error_finding_product_types_by_product_ids.app_error", nil, errMsg, statusCode)
	}

	return productTypes, nil
}

// ProductTypeByOption returns a product type with given option
func (s *ServiceProduct) ProductTypeByOption(options *product_and_discount.ProductTypeFilterOption) (*product_and_discount.ProductType, *model.AppError) {
	productType, err := s.srv.Store.ProductType().GetByOption(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("ProductTypeByOption", "app.product.error_finding_product_type_by_options.app_error", nil, err.Error(), statusCode)
	}

	return productType, nil
}

func (s *ServiceProduct) ProductTypesByOptions(options *product_and_discount.ProductTypeFilterOption) ([]*product_and_discount.ProductType, *model.AppError) {
	prdTypes, err := s.srv.Store.ProductType().FilterbyOption(options)
	var (
		statusCode int
		errMsg     string
	)

	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(prdTypes) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ProductTypesByOptions", "app.product.error_finding_product_types_by_options.app_error", nil, errMsg, statusCode)
	}

	return prdTypes, nil
}
