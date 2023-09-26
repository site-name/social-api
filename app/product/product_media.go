package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"gorm.io/gorm"
)

// ProductMediasByOption returns a list of product medias that satisfy given option
func (a *ServiceProduct) ProductMediasByOption(option *model.ProductMediaFilterOption) ([]*model.ProductMedia, *model.AppError) {
	productMedias, err := a.srv.Store.ProductMedia().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ProductMediasByOption", "app.product.error_finding_product_medias_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return productMedias, nil
}

func (s *ServiceProduct) DeleteProductMedias(tx *gorm.DB, ids []string) (int64, *model.AppError) {
	numDeleted, err := s.srv.Store.ProductMedia().Delete(tx, ids)
	if err != nil {
		return 0, model.NewAppError("DeleteProductMedias", "app.product.delete_medias.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return numDeleted, nil
}

func (s *ServiceProduct) UpsertProductMedias(tx *gorm.DB, medias model.ProductMedias) (model.ProductMedias, *model.AppError) {
	medias, err := s.srv.Store.ProductMedia().Upsert(tx, medias)
	if err != nil {
		return nil, model.NewAppError("UpsertProductMedias", "app.product.upsert_product_media.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return medias, nil
}
