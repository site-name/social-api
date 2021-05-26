package types

import (
	"github.com/sitename/sitename/web/model"
)

type ProductFilter struct {
	IsPublished       *bool                          `json:"is_published"`
	Collections       *[]string                      `json:"collections"`
	Categories        *[]string                      `json:"categories"`
	HasCategory       *bool                          `json:"has_category"`
	Price             *model.PriceRangeInput         `json:"price"`
	MinimalPrice      *model.PriceRangeInput         `json:"minimal_price"`
	Attributes        []*model.AttributeInput        `json:"attributes"`
	StockAvailability *model.StockAvailability       `json:"stock_availability"`
	ProductTypes      *[]string                      `json:"product_types"`
	Stocks            *model.ProductStockFilterInput `json:"stocks"`
	Search            *string                        `json:"search"` // search string
	Ids               *[]string                      `json:"ids"`    // ids of products
}

type ProductVariantFilter struct {
	Search *string   `json:"search"`
	Sku    *[]string `json:"sku"`
}

type CollectionFilter struct {
	Published *model.CollectionPublished `json:"published"`
	Search    *string                    `json:"search"`
	Ids       *[]string                  `json:"ids"`
}

type CategoryFilter struct {
	Search *string   `json:"search"`
	Ids    *[]string `json:"ids"`
}

type ProductTypeFilter struct {
	Search       *string                        `json:"search"`
	Configurable *model.ProductTypeConfigurable `json:"configurable"`
	ProductType  *model.ProductTypeEnum         `json:"product_type"`
	Ids          *[]string                      `json:"ids"`
}
