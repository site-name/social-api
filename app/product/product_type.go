package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

func (a *ServiceProduct) ProductTypesByCheckoutToken(checkoutToken string) ([]*model.ProductType, *model.AppError) {
	productTypes, err := a.srv.Store.ProductType().FilterProductTypesByCheckoutToken(checkoutToken)
	if err != nil {
		return nil, model.NewAppError("ProductTypesByCheckoutToken", "app.product.error_finding_product_types_by_checkout_token.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return productTypes, nil
}

// ProductTypesByProductIDs returns all product types that belong to given products
func (a *ServiceProduct) ProductTypesByProductIDs(productIDs []string) ([]*model.ProductType, *model.AppError) {
	productTypes, err := a.srv.Store.ProductType().ProductTypesByProductIDs(productIDs)
	if err != nil {
		return nil, model.NewAppError("ProductTypesByProductIDs", "app.product.error_finding_product_types_by_product_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return productTypes, nil
}

// ProductTypeByOption returns a product type with given option
func (s *ServiceProduct) ProductTypeByOption(options *model.ProductTypeFilterOption) (*model.ProductType, *model.AppError) {
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

func (s *ServiceProduct) ProductTypesByOptions(options *model.ProductTypeFilterOption) (int64, []*model.ProductType, *model.AppError) {
	count, productTypes, err := s.srv.Store.ProductType().FilterbyOption(options)
	if err != nil {
		return 0, nil, model.NewAppError("ProductTypesByOptions", "app.product.error_finding_product_types_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return count, productTypes, nil
}

// func (s *ServiceProduct) CountProductTypesByOptions(options *model.ProductTypeFilterOption) (int64, *model.AppError) {
// 	if options == nil {
// 		options = &model.ProductTypeFilterOption{}
// 	}
// 	count, err := s.srv.Store.ProductType().Count(options)
// 	if err != nil {
// 		return 0, model.NewAppError("CountProductTypesByOptions", "app.product.error_counting_product_types_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
// 	}

// 	return count, nil
// }

func (s *ServiceProduct) UpsertProductType(tx *gorm.DB, productType *model.ProductType) (*model.ProductType, *model.AppError) {
	productType, err := s.srv.Store.ProductType().Save(tx, productType)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertProductType", "app.product.upsert_product_type.app_error", nil, err.Error(), statusCode)
	}

	return productType, nil
}

func (s *ServiceProduct) DeleteProductTypes(tx *gorm.DB, ids []string) (int64, *model.AppError) {
	numDeleted, err := s.srv.Store.ProductType().Delete(tx, ids)
	if err != nil {
		return 0, model.NewAppError("DeleteProductTypes", "app.product.delete_product_types.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return numDeleted, nil
}

func (s *ServiceProduct) ToggleProductTypeAttributeRelations(tx *gorm.DB, productTypeID string, variantAttributes, productAttributes model.Attributes, isDelete bool) *model.AppError {
	err := s.srv.Store.ProductType().ToggleProductTypeRelations(tx, productTypeID, productAttributes, variantAttributes, isDelete)
	if err != nil {
		return model.NewAppError("ToggleProductTypeAttributeRelations", "app.product.toggle_product_type_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
