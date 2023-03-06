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
	StockAvailability *string
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
	ProductOrderFieldName            = "NAME"
	ProductOrderFieldRank            = "RANK"
	ProductOrderFieldPrice           = "PRICE"
	ProductOrderFieldMinimalPrice    = "MINIMAL_PRICE"
	ProductOrderFieldDate            = "DATE"
	ProductOrderFieldType            = "TYPE"
	ProductOrderFieldPublished       = "PUBLISHED"
	ProductOrderFieldPublicationDate = "PUBLICATION_DATE"
	ProductOrderFieldCollection      = "COLLECTION"
	ProductOrderFieldRating          = "RATING"
)

// TODO: add support for field AttributeID
type ProductOrder struct {
	Field     *string
	Direction OrderDirection

	// AttributeID *string // TODO: add support filtering by this field
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
