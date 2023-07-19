package model

import (
	"fmt"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"gorm.io/gorm"
)

// max lengths for some fields of product variant
const (
	PRODUCT_VARIANT_NAME_MAX_LENGTH = 255
	PRODUCT_VARIANT_SKU_MAX_LENGTH  = 255
)

// sort by sku
type ProductVariant struct {
	Id                      string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:Id"`
	Name                    string                 `json:"name" gorm:"type:varchar(255);column:Name"` // varchar(255)
	ProductID               string                 `json:"product_id" gorm:"type:uuid;index:productvariants_productid_index_key;column:ProductID"`
	Sku                     string                 `json:"sku" gorm:"type:varchar(255);column:Sku"` // varchar(255)
	Weight                  *float32               `json:"weight" gorm:"column:Weight"`
	WeightUnit              measurement.WeightUnit `json:"weight_unit" gorm:"column:WeightUnit"`
	TrackInventory          *bool                  `json:"track_inventory" gorm:"column:TrackInventory"` // default *true
	IsPreOrder              bool                   `json:"is_preorder" gorm:"column:IsPreOrder"`
	PreorderEndDate         *int64                 `json:"preorder_end_date" column:"type:bigint;column:PreorderEndDate"`
	PreOrderGlobalThreshold *int                   `json:"preorder_global_threshold" gorm:"type:smallint;column:PreOrderGlobalThreshold"`
	Sortable
	ModelMetadata

	digitalContent         *DigitalContent               `gorm:"-"` // for storing value returned by prefetching
	product                *Product                      `gorm:"-"`
	stocks                 Stocks                        `gorm:"-"`
	variantChannelListings ProductVariantChannelListings `gorm:"-"`

	Sales             Sales                       `json:"-" gorm:"many2many:SaleProductVariants"`
	Vouchers          Vouchers                    `json:"-" gorm:"many2many:voucherproductvariants"`
	ProductMedias     ProductMedias               `json:"-" gorm:"many2many:VariantMedias"`
	Attributes        []*AssignedVariantAttribute `json:"-" gorm:"foreignKey:VariantID"`
	AttributesRelated []*AttributeVariant         `json:"-" gorm:"many2many:AssignedVariantAttributes"`
}

func (p *ProductVariant) BeforeCreate(_ *gorm.DB) error {
	p.commonPre()
	return nil
}

// ProductVariantFilterOption is used to build sql queries
type ProductVariantFilterOption struct {
	Id        squirrel.Sqlizer
	Name      squirrel.Sqlizer
	ProductID squirrel.Sqlizer

	WishlistItemID squirrel.Sqlizer // INNER JOIN WishlistItemProductVariants ON (...) WHERE WishlistItemProductVariants.WishlistItemID ...
	WishlistID     squirrel.Sqlizer // INNER JOIN WishlistItemProductVariants ON (...) INNER JOIN WishlistItems ON (...) WHERE WishlistItems.WishlistID ...

	ProductVariantChannelListingPriceAmount squirrel.Sqlizer // INNER JOIN `ProductVariantChannelListing` ON ... WHERE ProductVariantChannelListing.PriceAmount ...
	ProductVariantChannelListingChannelSlug squirrel.Sqlizer // INNER JOIN `ProductVariantChannelListing` ON ... INNER JOIN Channels ON ... WHERE Channels.Slug ...
	ProductVariantChannelListingChannelID   squirrel.Sqlizer // INNER JOIN ProductVariantChannelListing ON ... WHERE ProductVariantChannelListing.ChannelID ...

	Distinct bool // if true, use SELECT DISTINCT

	VoucherID squirrel.Sqlizer // INNER JOIN VariantVouchers ON ... WHERE VariantVouchers.VoucherID ...
	SaleID    squirrel.Sqlizer // INNER JOIN VariantSales ON ... WHERE VariantSales.SaleID ...

	SelectRelatedDigitalContent bool // if true, JOIN Digital content table and attach related values to returning values(s)
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

func (s *ProductVariant) SetStocks(stk Stocks) {
	s.stocks = stk
}
func (p *ProductVariant) GetStocks() Stocks {
	return p.stocks
}

func (s *ProductVariant) SetProduct(prd *Product) {
	s.product = prd
}
func (p *ProductVariant) GetProduct() *Product {
	return p.product
}

func (s *ProductVariant) SetDigitalContent(d *DigitalContent) {
	s.digitalContent = d
}
func (p *ProductVariant) GetDigitalContent() *DigitalContent {
	return p.digitalContent
}

func (s *ProductVariant) SetVariantChannelListings(d ProductVariantChannelListings) {
	s.variantChannelListings = d
}
func (s *ProductVariant) AppendVariantChannelListing(d *ProductVariantChannelListing) {
	s.variantChannelListings = append(s.variantChannelListings, d)
}
func (p *ProductVariant) GetVariantChannelListings() ProductVariantChannelListings {
	return p.variantChannelListings
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
	outer := CreateAppErrorForModel(
		"model.product_variant.is_valid.%s.app_error",
		"product_variant_id=",
		"ProductVariant.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if len(p.Sku) > PRODUCT_VARIANT_SKU_MAX_LENGTH {
		return outer("sku", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_VARIANT_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
	}
	if p.Weight != nil && *p.Weight <= 0 {
		return outer("weight", &p.Id)
	}
	if p.WeightUnit != "" {
		if _, ok := measurement.WEIGHT_UNIT_CONVERSION[p.WeightUnit]; !ok {
			return outer("weight_unit", &p.Id)
		}
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
	return p.IsPreOrder && (p.PreorderEndDate == nil || (p.PreorderEndDate != nil && GetMillis() <= *p.PreorderEndDate))

}

func (p *ProductVariant) ToJSON() string {
	return ModelToJson(p)
}

func (p *ProductVariant) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.commonPre()
	p.ModelMetadata.PopulateFields()
}

func (p *ProductVariant) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
	if p.TrackInventory == nil {
		p.TrackInventory = NewPrimitive(true)
	}
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

func (p *ProductVariant) PreUpdate() {
	p.commonPre()
	p.ModelMetadata.PopulateFields()
}

func (p *ProductVariant) DeepCopy() *ProductVariant {
	if p == nil {
		return nil
	}

	res := *p

	if p.Weight != nil {
		res.Weight = NewPrimitive(*p.Weight)
	}
	if p.TrackInventory != nil {
		res.TrackInventory = NewPrimitive(*p.TrackInventory)
	}
	if p.PreorderEndDate != nil {
		res.PreorderEndDate = NewPrimitive(*p.PreorderEndDate)
	}
	if p.PreOrderGlobalThreshold != nil {
		res.PreOrderGlobalThreshold = NewPrimitive(*p.PreOrderGlobalThreshold)
	}
	if p.SortOrder != nil {
		res.SortOrder = NewPrimitive(*p.SortOrder)
	}

	res.ModelMetadata = p.ModelMetadata.DeepCopy()

	if p.product != nil {
		res.product = p.product.DeepCopy()
	}
	if p.digitalContent != nil {
		res.digitalContent = p.digitalContent.DeepCopy()
	}
	if p.stocks != nil {
		res.stocks = p.stocks.DeepCopy()
	}

	return &res
}

type ProductVariantTranslation struct {
	Id               string           `json:"id"`
	LanguageCode     LanguageCodeEnum `json:"language_code"`
	ProductVariantID string           `json:"product_variant_id"`
	Name             string           `json:"name"`
}

// ProductVariantTranslationFilterOption is used to build squirrel sql queries
type ProductVariantTranslationFilterOption struct {
	Id               squirrel.Sqlizer
	LanguageCode     squirrel.Sqlizer
	ProductVariantID squirrel.Sqlizer
	Name             squirrel.Sqlizer
}

func (p *ProductVariantTranslation) String() string {
	if p.Name != "" {
		return p.Name
	}

	return p.ProductVariantID
}

func (p *ProductVariantTranslation) PreSave() {
	if !IsValidId(p.Id) {
		p.Id = NewId()
	}
	p.commonPre()
}

func (p *ProductVariantTranslation) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
}

func (p *ProductVariantTranslation) PreUpdate() {
	p.commonPre()
}

func (p *ProductVariantTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.product_variant_translation.is_valid.%s.app_error",
		"product_variant_translation_id=",
		"ProductVariantTranslation.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.ProductVariantID) {
		return outer("product_variant_id", &p.Id)
	}
	if !p.LanguageCode.IsValid() {
		return outer("language_code", &p.Id)
	}

	return nil
}

func (p *ProductVariantTranslation) ToJSON() string {
	return ModelToJson(p)
}
