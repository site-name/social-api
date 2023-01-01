package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"golang.org/x/text/language"
)

// max lengths for some fields of product variant
const (
	PRODUCT_VARIANT_NAME_MAX_LENGTH = 255
	PRODUCT_VARIANT_SKU_MAX_LENGTH  = 255
)

type ProductVariant struct {
	Id                      string                 `json:"id"`
	Name                    string                 `json:"name"`
	ProductID               string                 `json:"product_id"`
	Sku                     *string                `json:"sku"`
	Weight                  *float32               `json:"weight"`
	WeightUnit              measurement.WeightUnit `json:"weight_unit"`
	TrackInventory          *bool                  `json:"track_inventory"` // default *true
	IsPreOrder              bool                   `json:"is_preorder"`
	PreorderEndDate         *int64                 `json:"preorder_end_date"`
	PreOrderGlobalThreshold *int                   `json:"preorder_global_threshold"`
	Sortable
	ModelMetadata

	DigitalContent *DigitalContent `json:"-" db:"-"` // for storing value returned by prefetching
	Product        *Product        `json:"-" db:"-"`
	stocks         Stocks          `json:"-" db:"-"`
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

	Distinct bool // if true, use SELECT DISTINCT

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
	if p.Sku != nil && len(*p.Sku) > PRODUCT_VARIANT_SKU_MAX_LENGTH {
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
	if p.Sku != nil {
		return *p.Sku
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
		p.TrackInventory = NewBool(true)
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

	if p.Sku != nil {
		res.Sku = NewString(*p.Sku)
	}
	if p.Weight != nil {
		res.Weight = NewFloat32(*p.Weight)
	}
	if p.TrackInventory != nil {
		res.TrackInventory = NewBool(*p.TrackInventory)
	}
	if p.PreorderEndDate != nil {
		res.PreorderEndDate = NewInt64(*p.PreorderEndDate)
	}
	if p.PreOrderGlobalThreshold != nil {
		res.PreOrderGlobalThreshold = NewInt(*p.PreOrderGlobalThreshold)
	}
	if p.SortOrder != nil {
		res.SortOrder = NewInt(*p.SortOrder)
	}

	res.ModelMetadata = p.ModelMetadata.DeepCopy()

	if p.Product != nil {
		res.Product = p.Product.DeepCopy()
	}
	if p.DigitalContent != nil {
		res.DigitalContent = p.DigitalContent.DeepCopy()
	}
	if p.stocks != nil {
		res.stocks = p.stocks.DeepCopy()
	}

	return &res
}

type ProductVariantTranslation struct {
	Id               string `json:"id"`
	LanguageCode     string `json:"language_code"`
	ProductVariantID string `json:"product_variant_id"`
	Name             string `json:"name"`
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
	tag, err := language.Parse(p.LanguageCode)
	if err != nil || !strings.EqualFold(tag.String(), p.LanguageCode) || Languages[strings.ToLower(p.LanguageCode)] == "" {
		return outer("language_code", &p.Id)
	}

	return nil
}

func (p *ProductVariantTranslation) ToJSON() string {
	return ModelToJson(p)
}
