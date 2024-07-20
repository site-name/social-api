package discount

import (
	"net/http"
	"time"

	"github.com/mattermost/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *ServiceDiscount) UpsertSale(transaction boil.ContextTransactor, sale model.Sale) (*model.Sale, *model_helper.AppError) {
	upsertSale, err := a.srv.Store.DiscountSale().Upsert(transaction, sale)
	if err != nil {
		return nil, model_helper.NewAppError("UpsertSale", "app.discount.error_upsert_sale.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return upsertSale, nil
}

func (a *ServiceDiscount) GetSaleDiscount(sale model.Sale, saleChannelListing model.SaleChannelListing) (types.DiscountCalculator, *model_helper.AppError) {
	if sale.Type == model.DiscountValueTypeFixed {
		discountAmount, _ := goprices.NewMoneyFromDecimal(saleChannelListing.DiscountValue, saleChannelListing.Currency.String())
		return a.Decorator(discountAmount), nil
	}
	return a.Decorator(saleChannelListing.DiscountValue), nil
}

// FilterSalesByOption should be used to filter active or expired sales
// refer: saleor/discount/models.SaleQueryset for details
func (a *ServiceDiscount) FilterSalesByOption(option model_helper.SaleFilterOption) (model.SaleSlice, *model_helper.AppError) {
	sales, err := a.srv.Store.DiscountSale().FilterSalesByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("ServiceDiscount.FilterSalesByOption", "app.discount.filter_sales_by_options.app_error", nil, err.Error(), statusCode)
	}

	return sales, nil
}

// ActiveSales finds active sales by given date. If date is nil then set date to UTC now
//
//	(end_date == NULL || end_date >= date) && start_date <= date
func (a *ServiceDiscount) ActiveSales(date *time.Time) (model.SaleSlice, *model_helper.AppError) {
	if date == nil {
		date = model_helper.GetPointerOfValue(time.Now().UTC())
	}

	activeSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(model_helper.SaleFilterOption{
			// Conditions: squirrel.And{
			// 	squirrel.LtOrEq{model.SaleTableName + ".StartDate": *date},
			// 	squirrel.Or{
			// 		squirrel.Eq{model.SaleTableName + ".EndDate": nil},
			// 		squirrel.GtOrEq{model.SaleTableName + ".EndDate": *date},
			// 	},
			// },

			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.SaleWhere.StartDate.LTE(*date),
			),
		})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("ServiceDiscount.ActiveSales", "app.discount.active_sales_by_date.app_error", nil, err.Error(), statusCode)
	}

	return activeSalesByDate, nil
}

// ExpiredSales returns sales that are expired by date. If date is nil, default to UTC now
//
//	end_date <= date && start_date <= date
func (a *ServiceDiscount) ExpiredSales(date *time.Time) ([]*model.Sale, *model_helper.AppError) {
	if date == nil {
		date = model_helper.GetPointerOfValue(time.Now().UTC())
	}

	expiredSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(model_helper.SaleFilterOption{
			// Conditions: squirrel.Lt{
			// 	model.SaleTableName + ".EndDate":   *date,
			// 	model.SaleTableName + ".StartDate": *date,
			// },
			CommonQueryOptions: model_helper.NewCommonQueryOptions(),
		})

	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("ServiceDiscount.ExpiredSales", "app.discount.expired_sales_by_date.app_error", nil, err.Error(), statusCode)
	}

	return expiredSalesByDate, nil
}

func (s *ServiceDiscount) ToggleSaleRelations(transaction boil.ContextTransactor, saleID string, productIDs, variantIDs, categoryIDs, collectionIDs []string, isDelete bool) *model_helper.AppError {
	err := s.srv.Store.DiscountSale().ToggleSaleRelations(transaction, model.SaleSlice{{ID: saleID}}, collectionIDs, productIDs, variantIDs, categoryIDs, isDelete)
	if err != nil {
		return model_helper.NewAppError("ToggleSaleRelations", "app.discount.insert_sale_relations.app_error", nil, "failed to insert sale relations", http.StatusInternalServerError)
	}

	return nil
}

// SaleCollectionsByOptions returns a slice of sale-collection relations filtered using given options
func (s *ServiceDiscount) SaleCollectionsByOptions(options squirrel.Sqlizer) ([]*model.SaleCollection, *model_helper.AppError) {
	args, err := store.BuildSqlizer(options, "SaleCollectionsByOptions")
	if err != nil {
		return nil, model_helper.NewAppError("SaleCollectionsByOptions", model_helper.InvalidArgumentAppErrorID, nil, err.Error(), http.StatusBadRequest)
	}

	var res []*model.SaleCollection
	err = s.srv.Store.GetReplica().Table(model.SaleCollectionTableName).Find(&res, args...).Error
	if err != nil {
		return nil, model_helper.NewAppError("SaleCollectionsByOptions", "app.discount.sale_collections_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

func (s *ServiceDiscount) SaleCategoriesByOption(option squirrel.Sqlizer) ([]*model.SaleCategory, *model_helper.AppError) {
	args, err := store.BuildSqlizer(option, "SaleCategoriesByOption")
	if err != nil {
		return nil, model_helper.NewAppError("SaleCategoriesByOption", model_helper.InvalidArgumentAppErrorID, nil, err.Error(), http.StatusBadRequest)
	}

	var res []*model.SaleCategory
	err = s.srv.Store.GetReplica().Table(model.SaleCategoryTableName).Find(&res, args...).Error
	if err != nil {
		return nil, model_helper.NewAppError("SaleCategoriesByOption", "app.discount.sale_categories_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

// SaleProductsByOptions returns a slice of sale-product relations filtered using given options
func (s *ServiceDiscount) SaleProductsByOptions(options squirrel.Sqlizer) ([]*model.SaleProduct, *model_helper.AppError) {
	args, err := store.BuildSqlizer(options, "SaleProductsByOptions")
	if err != nil {
		return nil, model_helper.NewAppError("SaleProductsByOptions", model_helper.InvalidArgumentAppErrorID, nil, err.Error(), http.StatusBadRequest)
	}
	var res []*model.SaleProduct
	err = s.srv.Store.GetReplica().Table(model.SaleProductTableName).Find(&res, args...).Error
	if err != nil {
		return nil, model_helper.NewAppError("SaleProductsByOptions", "app.discount.sale_product_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

// SaleProductVariantsByOptions returns a list of sale-product variant relations filtered using given options
func (s *ServiceDiscount) SaleProductVariantsByOptions(options squirrel.Sqlizer) ([]*model.SaleProductVariant, *model_helper.AppError) {
	args, err := store.BuildSqlizer(options, "SaleProductVariantsByOptions")
	if err != nil {
		return nil, model_helper.NewAppError("SaleProductVariantsByOptions", model_helper.InvalidArgumentAppErrorID, nil, err.Error(), http.StatusBadRequest)
	}
	var res []*model.SaleProductVariant
	err = s.srv.Store.GetReplica().Table(model.SaleProductVariantTableName).Find(&res, args...).Error
	if err != nil {
		return nil, model_helper.NewAppError("SaleProductVariantsByOptions", "app.discount.error_finding_sale_product_variants_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

func (s *ServiceDiscount) SaleChannelListingsByOptions(options *model.SaleChannelListingFilterOption) ([]*model.SaleChannelListing, *model_helper.AppError) {
	listings, err := s.srv.Store.DiscountSaleChannelListing().SaleChannelListingsWithOption(options)
	if err != nil {
		return nil, model_helper.NewAppError("SaleChannelListingsByOptions", "app.discount.error_finding_sale_channel_listings_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return listings, nil
}
