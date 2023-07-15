package discount

import (
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

func (a *ServiceDiscount) GetSaleDiscount(sale *model.Sale, saleChannelListing *model.SaleChannelListing) (types.DiscountCalculator, *model.AppError) {
	if saleChannelListing == nil {
		return nil, model.NewAppError("GetSaleDiscount", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "saleChannelListing"}, "", http.StatusBadRequest)
	}

	if sale.Type == model.FIXED {
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
func (a *ServiceDiscount) FilterSalesByOption(option *model.SaleFilterOption) ([]*model.Sale, *model.AppError) {
	sales, err := a.srv.Store.DiscountSale().FilterSalesByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ServiceDiscount.FilterSalesByOption", "app.discount.filter_sales_by_options.app_error", nil, err.Error(), statusCode)
	}

	return sales, nil
}

// ActiveSales finds active sales by given date. If date is nil then set date to UTC now
//
//	(end_date == NULL || end_date >= date) && start_date <= date
func (a *ServiceDiscount) ActiveSales(date *time.Time) (model.Sales, *model.AppError) {
	if date == nil {
		date = util.NewTime(time.Now().UTC())
	}

	activeSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(&model.SaleFilterOption{
			EndDate: squirrel.Or{
				squirrel.Eq{model.SaleTableName + ".EndDate": nil},
				squirrel.GtOrEq{model.SaleTableName + ".EndDate": *date},
			},
			StartDate: squirrel.LtOrEq{model.SaleTableName + ".StartDate": *date},
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

	expiredSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(&model.SaleFilterOption{
			EndDate:   squirrel.Lt{model.SaleTableName + ".EndDate": *date},
			StartDate: squirrel.Lt{model.SaleTableName + ".StartDate": *date},
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

func (s *ServiceDiscount) AddSaleRelations(transaction *gorm.DB, saleID string, productIDs, variantIDs, categoryIDs, collectionIDs []string) *model.AppError {
	if len(productIDs) > 0 {
		products := lo.Map(productIDs, func(id string, _ int) *model.Product { return &model.Product{Id: id} })
		err := s.srv.Store.DiscountSale().AddSaleRelations(transaction, model.Sales{{Id: saleID}}, products)
		if err != nil {
			return model.NewAppError("UpsertSaleRelations", "app.discount.insert_sale_relations.app_error", nil, "failed to insert sale-product relations", http.StatusInternalServerError)
		}
	}

	if len(variantIDs) > 0 {
		variants := lo.Map(variantIDs, func(id string, _ int) *model.ProductVariant { return &model.ProductVariant{Id: id} })
		err := s.srv.Store.DiscountSale().AddSaleRelations(transaction, model.Sales{{Id: saleID}}, variants)
		if err != nil {
			return model.NewAppError("UpsertSaleRelations", "app.discount.insert_sale_relations.app_error", nil, "failed to insert sale-variant relations", http.StatusInternalServerError)
		}
	}

	if len(collectionIDs) > 0 {
		collections := lo.Map(collectionIDs, func(id string, _ int) *model.Collection { return &model.Collection{Id: id} })
		err := s.srv.Store.DiscountSale().AddSaleRelations(transaction, model.Sales{{Id: saleID}}, collections)
		if err != nil {
			return model.NewAppError("UpsertSaleRelations", "app.discount.insert_sale_relations.app_error", nil, "failed to insert sale-collection relations", http.StatusInternalServerError)
		}
	}

	if len(categoryIDs) > 0 {
		categories := lo.Map(categoryIDs, func(id string, _ int) *model.Category { return &model.Category{Id: id} })
		err := s.srv.Store.DiscountSale().AddSaleRelations(transaction, model.Sales{{Id: saleID}}, categories)
		if err != nil {
			return model.NewAppError("UpsertSaleRelations", "app.discount.insert_sale_relations.app_error", nil, "failed to insert sale-collection relations", http.StatusInternalServerError)
		}
	}

	return nil
}
