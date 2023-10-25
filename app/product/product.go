package product

import (
	"net/http"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type ServiceProduct struct {
	srv         *app.Server
	categoryMap sync.Map // this is cache for all categories of this system, look up by category ids
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		service := &ServiceProduct{srv: s}

		appErr := service.DoAnalyticCategories()
		if appErr != nil {
			return appErr
		}
		s.Product = service

		return nil
	})
}

func (s *ServiceProduct) UpsertProduct(tx *gorm.DB, product *model.Product) (*model.Product, *model.AppError) {
	product, err := s.srv.Store.Product().Save(tx, product)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertProduct", "app.product.upsert_product.app_error", nil, err.Error(), statusCode)
	}

	return product, nil
}

// ProductById returns 1 product by given id
func (a *ServiceProduct) ProductById(productID string) (*model.Product, *model.AppError) {
	return a.ProductByOption(&model.ProductFilterOption{
		Conditions: squirrel.Expr(model.ProductTableName+".Id = ?", productID),
	})
}

// ProductsByOption returns a list of products that satisfy given option
func (a *ServiceProduct) ProductsByOption(option *model.ProductFilterOption) (model.Products, *model.AppError) {
	products, err := a.srv.Store.Product().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ProductsByOption", "app.product.error_finding_products_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return products, nil
}

// ProductByOption returns 1 product that satisfy given option
func (a *ServiceProduct) ProductByOption(option *model.ProductFilterOption) (*model.Product, *model.AppError) {
	product, err := a.srv.Store.Product().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ProductByOption", "app.error_finding_product_by_option.app_error", nil, err.Error(), statusCode)
	}

	return product, nil
}

// ProductsByVoucherID finds all products that have relationships with given voucher
func (a *ServiceProduct) ProductsByVoucherID(voucherID string) ([]*model.Product, *model.AppError) {
	products, appErr := a.ProductsByOption(&model.ProductFilterOption{
		VoucherID: squirrel.Expr(model.VoucherProductTableName+".voucher_id = ?", voucherID),
	})
	if appErr != nil {
		return nil, appErr
	}

	return products, nil
}

// ProductsRequireShipping checks if at least 1 product require shipping, then return true, false otherwise
func (a *ServiceProduct) ProductsRequireShipping(productIDs []string) (bool, *model.AppError) {
	productTypes, appErr := a.ProductTypesByProductIDs(productIDs)
	if appErr != nil { // this error caused by system
		return false, appErr
	}

	return lo.SomeBy(productTypes, func(p *model.ProductType) bool {
		return p != nil &&
			p.IsShippingRequired != nil &&
			*p.IsShippingRequired
	}), nil
}

// ProductGetFirstImage returns first media of given product
func (a *ServiceProduct) ProductGetFirstImage(productID string) (*model.ProductMedia, *model.AppError) {
	productMedias, appErr := a.ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Expr(model.ProductMediaTableName+".ProductID = ? AND ProductMedias.Type = ?", productID, model.IMAGE),
	})
	if appErr != nil {
		return nil, appErr
	}

	return productMedias[0], nil
}

func (a *ServiceProduct) GetVisibleToUserProducts(channelIdOrSlug string, userCanSeeAllProducts bool) (model.Products, *model.AppError) {
	productQuery := a.srv.Store.Product().VisibleToUserProductsQuery(channelIdOrSlug, userCanSeeAllProducts)
	products, err := a.srv.Store.Product().FilterByQuery(productQuery)
	if err != nil {
		return nil, model.NewAppError("GetVisibleToUserProducts", "app.product.get_visible_products_for_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return products, nil
}

func (s *ServiceProduct) FilterProductsAdvanced(options *model.ExportProductsFilterOptions, channelIdOrSlug string) (model.Products, *model.AppError) {
	productsQuery := s.srv.Store.Product().AdvancedFilterQueryBuilder(options)

	if channelIdOrSlug != "" {
		productsQuery = productsQuery.
			Where(`EXISTS (
				SELECT (1) AS "a"
				FROM ProductChannelListings PC
				WHERE (
					EXISTS (
						SELECT (1) AS "a"
						FROM Channels C
						WHERE (
							(C.Slug = ? OR C.Id = ?)
							AND C.Id = PC.ChannelId
						)
						LIMIT 1
					)
				)
				LIMIT 1
			)`, channelIdOrSlug, channelIdOrSlug)
	}

	products, err := s.srv.Store.Product().FilterByQuery(productsQuery)
	if err != nil {
		return nil, model.NewAppError("FilterProductsAdvanced", "app.product.filter_advanced_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return products, nil
}

func (s *ServiceProduct) SetDefaultProductVariantForProduct(productID, variantID string) (*model.Product, *model.AppError) {
	// validate if given variant belongs to given product
	variant, appErr := s.ProductVariantById(variantID)
	if appErr != nil {
		return nil, appErr
	}

	if variant.ProductID != productID {
		return nil, model.NewAppError("SetDefaultProductVariantForProduct", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "VariantID"}, "given product does not have given variant", http.StatusBadRequest)
	}

	// begin tx
	tx := s.srv.Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("SetDefaultProductVariantForProduct", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(tx)

	product, appErr := s.UpsertProduct(tx, &model.Product{
		Id:               productID,
		DefaultVariantID: &variantID,
	})
	if appErr != nil {
		return nil, appErr
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("SetDefaultProductVariantForProduct", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	pluginMng := s.srv.PluginService().GetPluginManager()
	_, appErr = pluginMng.ProductUpdated(*product)
	if appErr != nil {
		return nil, appErr
	}

	return product, nil
}
