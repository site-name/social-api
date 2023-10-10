package discount

import (
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

func (a *ServiceDiscount) UpsertSale(transaction *gorm.DB, sale *model.Sale) (*model.Sale, *model.AppError) {
	sale, err := a.srv.Store.DiscountSale().Upsert(transaction, sale)
	if err != nil {
		return nil, model.NewAppError("UpsertSale", "app.discount.error_upsert_sale.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return sale, nil
}

func (a *ServiceDiscount) GetSaleDiscount(sale *model.Sale, saleChannelListing *model.SaleChannelListing) (types.DiscountCalculator, *model.AppError) {
	if saleChannelListing == nil {
		return nil, model.NewAppError("GetSaleDiscount", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "saleChannelListing"}, "", http.StatusBadRequest)
	}

	if sale.Type == model.DISCOUNT_VALUE_TYPE_FIXED {
		discountAmount := &goprices.Money{ // can use directly here since sale channel listings are validated before saving
			Amount:   *saleChannelListing.DiscountValue,
			Currency: saleChannelListing.Currency,
		}
		return a.Decorator(discountAmount), nil
	}
	return a.Decorator(saleChannelListing.DiscountValue), nil
}

// FilterSalesByOption should be used to filter active or expired sales
// refer: saleor/discount/models.SaleQueryset for details
func (a *ServiceDiscount) FilterSalesByOption(option *model.SaleFilterOption) (int64, []*model.Sale, *model.AppError) {
	totalCount, sales, err := a.srv.Store.DiscountSale().FilterSalesByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return 0, nil, model.NewAppError("ServiceDiscount.FilterSalesByOption", "app.discount.filter_sales_by_options.app_error", nil, err.Error(), statusCode)
	}

	return totalCount, sales, nil
}

// ActiveSales finds active sales by given date. If date is nil then set date to UTC now
//
//	(end_date == NULL || end_date >= date) && start_date <= date
func (a *ServiceDiscount) ActiveSales(date *time.Time) (model.Sales, *model.AppError) {
	if date == nil {
		date = util.NewTime(time.Now().UTC())
	}

	_, activeSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(&model.SaleFilterOption{
			Conditions: squirrel.And{
				squirrel.LtOrEq{model.SaleTableName + ".StartDate": *date},
				squirrel.Or{
					squirrel.Eq{model.SaleTableName + ".EndDate": nil},
					squirrel.GtOrEq{model.SaleTableName + ".EndDate": *date},
				},
			},
		})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ServiceDiscount.ActiveSales", "app.discount.active_sales_by_date.app_error", nil, err.Error(), statusCode)
	}

	return activeSalesByDate, nil
}

// ExpiredSales returns sales that are expired by date. If date is nil, default to UTC now
//
//	end_date <= date && start_date <= date
func (a *ServiceDiscount) ExpiredSales(date *time.Time) ([]*model.Sale, *model.AppError) {
	if date == nil {
		date = util.NewTime(time.Now().UTC())
	}

	_, expiredSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(&model.SaleFilterOption{
			Conditions: squirrel.Lt{
				model.SaleTableName + ".EndDate":   *date,
				model.SaleTableName + ".StartDate": *date,
			},
		})

	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ServiceDiscount.ExpiredSales", "app.discount.expired_sales_by_date.app_error", nil, err.Error(), statusCode)
	}

	return expiredSalesByDate, nil
}

func (s *ServiceDiscount) ToggleSaleRelations(transaction *gorm.DB, saleID string, productIDs, variantIDs, categoryIDs, collectionIDs []string, isDelete bool) *model.AppError {
	err := s.srv.Store.DiscountSale().ToggleSaleRelations(transaction, model.Sales{{Id: saleID}}, collectionIDs, productIDs, variantIDs, categoryIDs, isDelete)
	if err != nil {
		return model.NewAppError("ToggleSaleRelations", "app.discount.insert_sale_relations.app_error", nil, "failed to insert sale relations", http.StatusInternalServerError)
	}

	return nil
}

// SaleCollectionsByOptions returns a slice of sale-collection relations filtered using given options
func (s *ServiceDiscount) SaleCollectionsByOptions(options squirrel.Sqlizer) ([]*model.SaleCollection, *model.AppError) {
	var res []*model.SaleCollection
	err := s.srv.Store.GetReplica().Table(model.SaleCollectionTableName).Find(&res, store.BuildSqlizer(options)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleCollectionsByOptions", "app.discount.sale_collections_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

// SaleCategoriesByOption returns sale-category relations with an app error
func (s *ServiceDiscount) SaleCategoriesByOption(option squirrel.Sqlizer) ([]*model.SaleCategory, *model.AppError) {
	var res []*model.SaleCategory
	err := s.srv.Store.GetReplica().Table("SaleCategories").Find(&res, store.BuildSqlizer(option)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleCategoriesByOption", "app.discount.sale_categories_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

// SaleProductsByOptions returns a slice of sale-product relations filtered using given options
func (s *ServiceDiscount) SaleProductsByOptions(options squirrel.Sqlizer) ([]*model.SaleProduct, *model.AppError) {
	var res []*model.SaleProduct
	err := s.srv.Store.GetReplica().Table("sale_collections").Find(&res, store.BuildSqlizer(options)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleProductsByOptions", "app.discount.sale_product_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

// SaleProductVariantsByOptions returns a list of sale-product variant relations filtered using given options
func (s *ServiceDiscount) SaleProductVariantsByOptions(options squirrel.Sqlizer) ([]*model.SaleProductVariant, *model.AppError) {
	var res []*model.SaleProductVariant
	err := s.srv.Store.GetReplica().Table("sale_productvariants").Find(&res, store.BuildSqlizer(options)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleProductVariantsByOptions", "app.discount.error_finding_sale_product_variants_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

func (s *ServiceDiscount) SaleChannelListingsByOptions(options *model.SaleChannelListingFilterOption) ([]*model.SaleChannelListing, *model.AppError) {
	listings, err := s.srv.Store.DiscountSaleChannelListing().SaleChannelListingsWithOption(options)
	if err != nil {
		return nil, model.NewAppError("SaleChannelListingsByOptions", "app.discount.error_finding_sale_channel_listings_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return listings, nil
}
