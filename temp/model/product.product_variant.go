package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

// sort by sku
type ProductVariant struct {
	Id                      string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:Id"`
	Name                    string                 `json:"name" gorm:"type:varchar(255);column:Name"` // varchar(255)
	ProductID               string                 `json:"product_id" gorm:"type:uuid;index:productvariants_productid_index_key;column:ProductID"`
	Sku                     string                 `json:"sku" gorm:"type:varchar(255);column:Sku;unique"` // varchar(255)
	Weight                  *float32               `json:"weight" gorm:"column:Weight"`
	WeightUnit              measurement.WeightUnit `json:"weight_unit" gorm:"column:WeightUnit"`
	TrackInventory          *bool                  `json:"track_inventory" gorm:"column:TrackInventory"` // default *true
	IsPreOrder              bool                   `json:"is_preorder" gorm:"column:IsPreOrder"`
	PreorderEndDate         *time.Time             `json:"preorder_end_date" gorm:"type:bigint;column:PreorderEndDate"`
	PreOrderGlobalThreshold *int                   `json:"preorder_global_threshold" gorm:"type:smallint;column:PreOrderGlobalThreshold"`
	Sortable
	ModelMetadata

	Stocks            Stocks                      `json:"-" gorm:"foreignKey:ProductVariantID"`
	DigitalContent    *DigitalContent             `json:"-" gorm:"foreignKey:ProductVariantID"` // for storing value returned by prefetching
	Attributes        []*AssignedVariantAttribute `json:"-" gorm:"foreignKey:VariantID"`
	OrderLines        OrderLines                  `json:"-" gorm:"foreignKey:VariantID"`
	Sales             Sales                       `json:"-" gorm:"many2many:SaleProductVariants"`
	Vouchers          Vouchers                    `json:"-" gorm:"many2many:VoucherVariants"`
	ProductMedias     ProductMedias               `json:"-" gorm:"many2many:VariantMedias"`
	AttributesRelated []*AttributeVariant         `json:"-" gorm:"many2many:AssignedVariantAttributes"`
	WishlistItems     WishlistItems               `json:"-" gorm:"many2many:WishlistItemProductVariants"`
	Product           *Product                    `json:"-"`
}

// column names for product variant table
const (
	ProductVariantColumnId                      = "Id"
	ProductVariantColumnName                    = "Name"
	ProductVariantColumnProductID               = "ProductID"
	ProductVariantColumnSku                     = "Sku"
	ProductVariantColumnWeight                  = "Weight"
	ProductVariantColumnWeightUnit              = "WeightUnit"
	ProductVariantColumnTrackInventory          = "TrackInventory"
	ProductVariantColumnIsPreOrder              = "IsPreOrder"
	ProductVariantColumnPreorderEndDate         = "PreorderEndDate"
	ProductVariantColumnPreOrderGlobalThreshold = "PreOrderGlobalThreshold"
)

func (c *ProductVariant) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariant) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariant) TableName() string             { return ProductVariantTableName }

// ProductVariantFilterOption is used to build sql queries
type ProductVariantFilterOption struct {
	Conditions squirrel.Sqlizer

	WishlistItemID squirrel.Sqlizer // INNER JOIN WishlistItemProductVariants ON (...) WHERE WishlistItemProductVariants.WishlistItemID ...
	WishlistID     squirrel.Sqlizer // INNER JOIN WishlistItemProductVariants ON (...) INNER JOIN WishlistItems ON (...) WHERE WishlistItems.WishlistID ...

	RelatedProductVariantChannelListingConditions squirrel.Sqlizer // INNER JOIN ProductVariantChannelListing ON ... WHERE ProductVariantChannelListing ...
	ProductVariantChannelListingChannelSlug       squirrel.Sqlizer // INNER JOIN `ProductVariantChannelListing` ON ... INNER JOIN Channels ON ... WHERE Channels.Slug ...

	Distinct bool // if true, use SELECT DISTINCT

	// NOTE: the key must be:
	//  "VariantVouchers.voucher_id"
	VoucherID squirrel.Sqlizer // INNER JOIN VariantVouchers ON ... WHERE VariantVouchers.voucher_id ...
	// NOTE: the key must be:
	//  "VariantSales.sale_id"
	SaleID squirrel.Sqlizer // INNER JOIN VariantSales ON ... WHERE VariantSales.sale_id ...

	// can be:
	//  "DigitalContent", ...
	Preloads []string
}

func (p *ProductVariant) WeightString() string {
	if p == nil || p.Weight == nil {
		return ""
	}

	u := p.WeightUnit
	if measurement.WEIGHT_UNIT_STRINGS[u] == "" {
		u = measurement.G
	}
	return fmt.Sprintf("%f %s", *p.Weight, u)
}

type ProductVariants []*ProductVariant

func (p ProductVariants) DeepCopy() ProductVariants {
	return lo.Map(p, func(v *ProductVariant, _ int) *ProductVariant { return v.DeepCopy() })
}

// FilterNils returns new ProductVariants contains all non-nil items from current ProductVariants
func (p ProductVariants) FilterNils() ProductVariants {
	return lo.Filter(p, func(v *ProductVariant, _ int) bool { return v != nil })
}

func (p ProductVariants) IDs() []string {
	return lo.Map(p, func(v *ProductVariant, _ int) string { return v.Id })
}

// ProductIDs returns all product ids of current product variants
func (p ProductVariants) ProductIDs() []string {
	return lo.Map(p, func(v *ProductVariant, _ int) string { return v.ProductID })
}

func (p *ProductVariant) IsValid() *AppError {
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}
	if p.Weight != nil && *p.Weight <= 0 {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight.app_error", nil, "please provide valid weight", http.StatusBadRequest)
	}
	if _, ok := measurement.WEIGHT_UNIT_CONVERSION[p.WeightUnit]; !ok {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight_unit.app_error", nil, "please provide valid weight unit", http.StatusBadRequest)
	}

	return nil
}

// String returns exact product variant name or Sku depends on their truth value
func (p *ProductVariant) String() string {
	if p.Name != "" {
		return p.Name
	}

	return fmt.Sprintf("ID:%s", p.Id)
}

func (p *ProductVariant) IsPreorderActive() bool {
	return p.IsPreOrder && (p.PreorderEndDate == nil || p.PreorderEndDate.After(util.StartOfDay(time.Now())))
}

func (p *ProductVariant) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
	if p.TrackInventory == nil {
		p.TrackInventory = GetPointerOfValue(true)
	}
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.G
	}
	p.ModelMetadata.PopulateFields()

}

func (p *ProductVariant) DeepCopy() *ProductVariant {
	if p == nil {
		return nil
	}

	res := *p

	res.Weight = CopyPointer(p.Weight)
	res.TrackInventory = CopyPointer(p.TrackInventory)
	res.PreorderEndDate = CopyPointer(p.PreorderEndDate)
	res.PreOrderGlobalThreshold = CopyPointer(p.PreOrderGlobalThreshold)
	res.SortOrder = CopyPointer(p.SortOrder)

	res.ModelMetadata = p.ModelMetadata.DeepCopy()
	if p.Product != nil {
		res.Product = p.Product.DeepCopy()
	}
	if p.DigitalContent != nil {
		res.DigitalContent = p.DigitalContent.DeepCopy()
	}
	if p.Stocks != nil {
		res.Stocks = p.Stocks.DeepCopy()
	}

	return &res
}

type ProductVariantTranslation struct {
	Id               string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode     LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode;index:languagecode_productvariantid_key"`
	ProductVariantID string           `json:"product_variant_id" gorm:"type:uuid;column:ProductVariantID;index:languagecode_productvariantid_key"`
	Name             string           `json:"name" gorm:"type:varchar(255);column:Name"`
}

func (c *ProductVariantTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariantTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariantTranslation) TableName() string             { return ProductVariantTranslationTableName }

// ProductVariantTranslationFilterOption is used to build squirrel sql queries
type ProductVariantTranslationFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (p *ProductVariantTranslation) String() string {
	if p.Name != "" {
		return p.Name
	}

	return p.ProductVariantID
}

func (p *ProductVariantTranslation) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
}

func (p *ProductVariantTranslation) IsValid() *AppError {
	if !IsValidId(p.ProductVariantID) {
		return NewAppError("ProductVariantTranslation.IsValid", "model.product_variant_translation.is_valid.product_variant_id.app_error", nil, "please provide valid product variant id", http.StatusBadRequest)
	}
	if !p.LanguageCode.IsValid() {
		return NewAppError("ProductVariantTranslation.IsValid", "model.product_variant_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}
