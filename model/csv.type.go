package model

import "time"

type AttributeFilter struct {
	Slug        string
	Values      []string
	ValuesRange *struct {
		Gte *int32
		Lte *int32
	}
	DateTime *struct {
		Gte *time.Time
		Lte *time.Time
	}
	Date *struct {
		Gte *time.Time
		Lte *time.Time
	}
	Boolean *bool
}

type ProductFilterInput struct {
	IsPublished       *bool
	Collections       []string
	Categories        []string
	HasCategory       *bool
	Attributes        []*AttributeFilter
	StockAvailability *StockAvailability // can either be Instock or outOfStock
	Stocks            *struct {
		WarehouseIds []string
		Quantity     *struct {
			Gte *int32
			Lte *int32
		}
	}
	Search   *string
	Metadata []*struct {
		Key   string
		Value string
	}
	Price *struct {
		Gte *float64
		Lte *float64
	}
	MinimalPrice *struct {
		Gte *float64
		Lte *float64
	}
	ProductTypes          []string
	GiftCard              *bool
	Ids                   []string
	HasPreorderedVariants *bool
	Channel               *string
}

// valid values for ProductOrder.Field
const (
	ProductOrderFieldName            ProductOrderField = "NAME"
	ProductOrderFieldRank            ProductOrderField = "RANK"
	ProductOrderFieldPrice           ProductOrderField = "PRICE"
	ProductOrderFieldMinimalPrice    ProductOrderField = "MINIMAL_PRICE"
	ProductOrderFieldDate            ProductOrderField = "DATE"
	ProductOrderFieldType            ProductOrderField = "TYPE"
	ProductOrderFieldPublished       ProductOrderField = "PUBLISHED"
	ProductOrderFieldPublicationDate ProductOrderField = "PUBLICATION_DATE"
	ProductOrderFieldCollection      ProductOrderField = "COLLECTION"
	ProductOrderFieldRating          ProductOrderField = "RATING"
)

type ProductOrderField string

func (p ProductOrderField) IsValid() bool {
	switch p {
	case ProductOrderFieldName, ProductOrderFieldRank,
		ProductOrderFieldPrice, ProductOrderFieldMinimalPrice,
		ProductOrderFieldDate, ProductOrderFieldType,
		ProductOrderFieldPublished, ProductOrderFieldPublicationDate,
		ProductOrderFieldCollection, ProductOrderFieldRating:
		return true
	default:
		return false
	}
}

type ProductOrder struct {
	Field       *ProductOrderField
	Direction   OrderDirection
	AttributeID *string
}

type ExportProductsFilterOptions struct {
	Scope      string // "all" or "ids" or "filter"
	Filter     *ProductFilterInput
	Ids        []string
	SortBy     *ProductOrder
	ExportInfo *struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	}
	FileType string // xlsx or csv
}
