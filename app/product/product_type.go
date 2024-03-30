package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// ProductTypeByOption returns a product type with given option
func (s *ServiceProduct) ProductTypeByOption(options *model.ProductTypeFilterOption) (*model.ProductType, *model_helper.AppError) {
	productType, err := s.srv.Store.ProductType().GetByOption(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("ProductTypeByOption", "app.product.error_finding_product_type_by_options.app_error", nil, err.Error(), statusCode)
	}

	return productType, nil
}

func (s *ServiceProduct) ProductTypesByOptions(options *model.ProductTypeFilterOption) (int64, []*model.ProductType, *model_helper.AppError) {
	count, productTypes, err := s.srv.Store.ProductType().FilterbyOption(options)
	if err != nil {
		return 0, nil, model_helper.NewAppError("ProductTypesByOptions", "app.product.error_finding_product_types_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return count, productTypes, nil
}

// func (s *ServiceProduct) CountProductTypesByOptions(options *model.ProductTypeFilterOption) (int64, *model_helper.AppError) {
// 	if options == nil {
// 		options = &model.ProductTypeFilterOption{}
// 	}
// 	count, err := s.srv.Store.ProductType().Count(options)
// 	if err != nil {
// 		return 0, model_helper.NewAppError("CountProductTypesByOptions", "app.product.error_counting_product_types_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
// 	}

// 	return count, nil
// }

func (s *ServiceProduct) UpsertProductType(tx *gorm.DB, productType *model.ProductType) (*model.ProductType, *model_helper.AppError) {
	productType, err := s.srv.Store.ProductType().Save(tx, productType)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("UpsertProductType", "app.product.upsert_product_type.app_error", nil, err.Error(), statusCode)
	}

	return productType, nil
}

func (s *ServiceProduct) DeleteProductTypes(tx *gorm.DB, ids []string) (int64, *model_helper.AppError) {
	numDeleted, err := s.srv.Store.ProductType().Delete(tx, ids)
	if err != nil {
		return 0, model_helper.NewAppError("DeleteProductTypes", "app.product.delete_product_types.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return numDeleted, nil
}

func (s *ServiceProduct) ToggleProductTypeAttributeRelations(tx *gorm.DB, productTypeID string, variantAttributes, productAttributes model.Attributes, isDelete bool) *model_helper.AppError {
	err := s.srv.Store.ProductType().ToggleProductTypeRelations(tx, productTypeID, productAttributes, variantAttributes, isDelete)
	if err != nil {
		return model_helper.NewAppError("ToggleProductTypeAttributeRelations", "app.product.toggle_product_type_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
