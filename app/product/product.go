package product

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type ServiceProduct struct {
	srv *app.Server
}

func init() {
	app.RegisterProductService(func(s *app.Server) (sub_app_iface.ProductService, error) {
		return &ServiceProduct{
			srv: s,
		}, nil
	})
}

// ProductById returns 1 product by given id
func (a *ServiceProduct) ProductById(productID string) (*model.Product, *model.AppError) {
	return a.ProductByOption(&model.ProductFilterOption{
		Id: squirrel.Eq{store.ProductTableName + ".Id": productID},
	})
}

// ProductsByOption returns a list of products that satisfy given option
func (a *ServiceProduct) ProductsByOption(option *model.ProductFilterOption) ([]*model.Product, *model.AppError) {
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
		VoucherID: squirrel.Eq{store.VoucherProductTableName + ".VoucherID": voucherID},
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
		ProductID: squirrel.Eq{store.ProductMediaTableName + ".ProductID": productID},
		Type:      squirrel.Eq{store.ProductMediaTableName + ".Type": model.IMAGE},
	})
	if appErr != nil {
		return nil, appErr
	}

	return productMedias[0], nil
}

func (a *ServiceProduct) GetVisibleProductsToUser(userSession *model.Session, channel_IdOrSlug string) (model.Products, *model.AppError) {
	userHasOneOfProductPermissions := a.srv.AccountService().SessionHasPermissionToAny(userSession, model.ProductPermissions...)
	productQuery := a.srv.Store.Product().VisibleToUserProducts(channel_IdOrSlug, userHasOneOfProductPermissions)
	products, err := a.srv.Store.Product().FilterByQuery(productQuery)
	if err != nil {
		return nil, model.NewAppError("GetVisibleProductsToUser", "app.product.get_visible_products_for_user.app_error", nil, err.Error(), http.StatusInternalServerError)
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
