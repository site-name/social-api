package product_and_discount

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
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
	model.Sortable
	model.ModelMetadata

	DigitalContent *DigitalContent                 `json:"-" db:"-"` // for storing value returned by prefetching
	Product        *Product                        `json:"-" db:"-"`
	ChannelListing []*ProductVariantChannelListing `json:"-" db:"-"`
}

// ProductVariantFilterOption is used to build sql queries
type ProductVariantFilterOption struct {
	Id   *model.StringFilter
	Name *model.StringFilter

	WishlistItemID *model.StringFilter // INNER JOIN WishlistItemProductVariants ON (...) WHERE WishlistItemProductVariants.WishlistItemID ...
	WishlistID     *model.StringFilter // INNER JOIN WishlistItemProductVariants ON (...) INNER JOIN WishlistItems ON (...) WHERE WishlistItems.WishlistID ...

	ProductVariantChannelListingPriceAmount *model.NumberFilter // LEFT JOIN `ProductVariantChannelListing`
	ProductVariantChannelListingChannelSlug *model.StringFilter // LEFT JOIN `ProductVariantChannelListing`

	Distinct bool // if true, use SELECT DISTINCT

	SelectRelatedDigitalContent bool // if true, JOIN Digital content table and attach related values to returning values(s)
}

type ProductVariants []*ProductVariant

// FilterNils returns new ProductVariants contains all non-nil items from current ProductVariants
func (p ProductVariants) FilterNils() ProductVariants {
	var res ProductVariants
	for _, item := range p {
		if item != nil {
			res = append(res, item)
		}
	}

	return res
}

func (p ProductVariants) IDs() []string {
	res := []string{}
	for _, item := range p {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

// ProductIDs returns all product ids of current product variants
func (p ProductVariants) ProductIDs() []string {
	res := make([]string, len(p))
	for i := range p {
		res[i] = p[i].ProductID
	}

	return res
}

func (p *ProductVariant) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_variant.is_valid.%s.app_error",
		"product_variant_id=",
		"ProductVariant.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductID) {
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
	return p.IsPreOrder && (p.PreorderEndDate == nil || (p.PreorderEndDate != nil && model.GetMillis() <= *p.PreorderEndDate))

}

func (p *ProductVariant) ToJSON() string {
	return model.ModelToJson(p)
}

func (p *ProductVariant) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.commonPre()
	p.ModelMetadata.PopulateFields()
}

func (p *ProductVariant) commonPre() {
	p.Name = model.SanitizeUnicode(p.Name)
	if p.TrackInventory == nil {
		p.TrackInventory = model.NewBool(true)
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
	res := *p
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
	Id               *model.StringFilter
	LanguageCode     *model.StringFilter
	ProductVariantID *model.StringFilter
	Name             *model.StringFilter
}

func (p *ProductVariantTranslation) String() string {
	if p.Name != "" {
		return p.Name
	}

	return p.ProductVariantID
}

func (p *ProductVariantTranslation) PreSave() {
	if !model.IsValidId(p.Id) {
		p.Id = model.NewId()
	}
	p.commonPre()
}

func (p *ProductVariantTranslation) commonPre() {
	p.Name = model.SanitizeUnicode(p.Name)
}

func (p *ProductVariantTranslation) PreUpdate() {
	p.commonPre()
}

func (p *ProductVariantTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_variant_translation.is_valid.%s.app_error",
		"product_variant_translation_id=",
		"ProductVariantTranslation.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductVariantID) {
		return outer("product_variant_id", &p.Id)
	}
	tag, err := language.Parse(p.LanguageCode)
	if err != nil || !strings.EqualFold(tag.String(), p.LanguageCode) || model.Languages[strings.ToLower(p.LanguageCode)] == "" {
		return outer("language_code", &p.Id)
	}

	return nil
}

func (p *ProductVariantTranslation) ToJSON() string {
	return model.ModelToJson(p)
}
