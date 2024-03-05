package model_helper

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func ProductPreSave(p *model.Product) {
	productCommonPre(p)
	if p.ID == "" {
		p.ID = NewId()
	}
	if p.CreatedAt == 0 {
		p.CreatedAt = GetMillis()
	}
	p.UpdatedAt = p.CreatedAt
	p.Slug = slug.Make(p.Name)
}

func productCommonPre(p *model.Product) {
	p.Name = SanitizeUnicode(p.Name)
	p.DescriptionPlainText = SanitizeUnicode(p.DescriptionPlainText)
	p.SeoTitle = SanitizeUnicode(p.SeoTitle)
	p.SeoDescription = SanitizeUnicode(p.SeoDescription)
	if p.ChargeTaxes.IsNil() {
		p.ChargeTaxes = model_types.NewNullBool(true)
	}
	if p.WeightUnit == "" {
		p.WeightUnit = measurement.G.String()
	}
}

func ProductPreUpdate(p *model.Product) {
	p.UpdatedAt = GetMillis()
	productCommonPre(p)
}

func ProductIsValid(p model.Product) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(p.CategoryID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.category_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.CreatedAt <= 0 {
		return NewAppError("Product.IsValid", "model.product.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if p.UpdatedAt <= 0 {
		return NewAppError("Product.IsValid", "model.product.is_valid.updated_at.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.DefaultVariantID.IsNil() && !IsValidId(*p.DefaultVariantID.String) {
		return NewAppError("Product.IsValid", "model.product.is_valid.default_variant_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Name == "" {
		return NewAppError("Product.IsValid", "model.product.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(p.Slug) {
		return NewAppError("Product.IsValid", "model.product.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}
	if measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(p.WeightUnit)] == "" {
		return NewAppError("Product.IsValid", "model.product.is_valid.weight_unit.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type ProductFilterOptions struct {
	CommonQueryOptions
}

func CategoryPreSave(c *model.Category) {
	if c.ID == "" {
		c.ID = NewId()
	}
	c.Slug = slug.Make(c.Name)
	CategoryCommonPre(c)
}

func CategoryCommonPre(c *model.Category) {
	c.Name = SanitizeUnicode(c.Name)
	c.SeoTitle = SanitizeUnicode(c.SeoTitle)
	c.SeoDescription = SanitizeUnicode(c.SeoDescription)
}

func CategoryIsValid(c model.Category) *AppError {
	if !IsValidId(c.ID) {
		return NewAppError("Category.IsValid", "model.category.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if c.Name == "" {
		return NewAppError("Category.IsValid", "model.category.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(c.Slug) {
		return NewAppError("Category.IsValid", "model.category.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}
	if !c.ParentID.IsNil() && !IsValidId(*c.ParentID.String) {
		return NewAppError("Category.IsValid", "model.category.is_valid.parent_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func ProductVariantPreSave(pv *model.ProductVariant) {
	if pv.ID == "" {
		pv.ID = NewId()
	}
	ProductVariantCommonPre(pv)
}

func ProductVariantCommonPre(pv *model.ProductVariant) {
	pv.Name = SanitizeUnicode(pv.Name)
	if pv.TrackInventory.IsNil() {
		pv.TrackInventory = model_types.NewNullBool(true)
	}
	if !pv.Weight.IsNil() && pv.WeightUnit == "" {
		pv.WeightUnit = measurement.G.String()
	}
}

func ProductVariantIsValid(p model.ProductVariant) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.product_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Name == "" {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.Weight.IsNil() && *p.Weight.Float32 < 0 {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.Weight.IsNil() && measurement.WEIGHT_UNIT_CONVERSION[measurement.WeightUnit(p.WeightUnit)] == 0 {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight_unit.app_error", nil, "", http.StatusBadRequest)
	}
	if p.IsPreorder && (p.PreorderEndDate.IsNil() || *p.PreorderEndDate.Int64 < GetMillis()) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.preorder_end_date.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func CollectionPreSave(c *model.Collection) {
	if c.ID == "" {
		c.ID = NewId()
	}
	CollectionCommonPre(c)
	c.Slug = slug.Make(c.Name)
}

func CollectionCommonPre(c *model.Collection) {
	c.Name = SanitizeUnicode(c.Name)
	c.SeoTitle = SanitizeUnicode(c.SeoTitle)
	c.SeoDescription = SanitizeUnicode(c.SeoDescription)
}

func CollectionIsValid(c model.Collection) *AppError {
	if !IsValidId(c.ID) {
		return NewAppError("Collection.IsValid", "model.collection.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if c.Name == "" {
		return NewAppError("Collection.IsValid", "model.collection.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(c.Slug) {
		return NewAppError("Collection.IsValid", "model.collection.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type CategoryFilterOption struct {
	CommonQueryOptions
	SaleID    qm.QueryMod // INNER JOIN SaleCategories ON ... WHERE SaleCategories.SaleID ...
	ProductID qm.QueryMod // INNER JOIN Products ON ... WHERE Products.ID ...
	VoucherID qm.QueryMod // INNER JOIN VoucherCategories ON ... WHERE VoucherCategories.VoucherID ...
}

type CollectionFilterOptions struct {
	CommonQueryOptions
	ProductID                                   qm.QueryMod // INNER JOIN product_cpllections ON ... WHERE product_collections.product_id ...
	VoucherID                                   qm.QueryMod // INNER JOIN voucher_collections ON ... WHERE voucher_collections.voucher_id ...
	SaleID                                      qm.QueryMod // INNER JOIN sale_collections ON ... WHERE sale_collections.sale_id ...
	RelatedCollectionChannelListingConds        qm.QueryMod // INNER JOIN collection_channel_listings ON ... WHERE collection_channel_listings...
	RelatedCollectionChannelListingChannelConds qm.QueryMod // INNER JOIN collection_channel_listings ON ... INNER JOIN chanel ON ... WHERE channels...
	AnnotateProductCount                        bool
}

func CollectionChannelListingPreSave(c *model.CollectionChannelListing) {
	if c.ID == "" {
		c.ID = NewId()
	}
	if c.CreatedAt == 0 {
		c.CreatedAt = GetMillis()
	}
}

func CollectionChannelListingIsValid(c model.CollectionChannelListing) *AppError {
	if !IsValidId(c.ID) {
		return NewAppError("CollectionChannelListing.IsValid", "model.collection_channel_listing.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(c.CollectionID) {
		return NewAppError("CollectionChannelListing.IsValid", "model.collection_channel_listing.is_valid.collection_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !c.ChannelID.IsNil() && !IsValidId(*c.ChannelID.String) {
		return NewAppError("CollectionChannelListing.IsValid", "model.collection_channel_listing.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}
	if c.CreatedAt <= 0 {
		return NewAppError("CollectionChannelListing.IsValid", "model.collection_channel_listing.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type CollectionChannelListingFilterOptions struct {
	CommonQueryOptions
}

type DigitalContentFilterOption struct {
	CommonQueryOptions
}

func DigitalContentPreSave(d *model.DigitalContent) {
	if d.ID == "" {
		d.ID = NewId()
	}
}

func DigitalContentIsValid(d model.DigitalContent) *AppError {
	if !IsValidId(d.ID) {
		return NewAppError("DigitalContent.IsValid", "model.digital_content.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if d.ContentType.IsValid() != nil {
		return NewAppError("DigitalContent.IsValid", "model.digital_content.is_valid.content_type.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type DigitalContentUrlFilterOptions struct {
	CommonQueryOptions
}

func DigitalContentUrlPreSave(d *model.DigitalContentURL) {
	if d.ID == "" {
		d.ID = NewId()
	}
	if d.CreatedAt == 0 {
		d.CreatedAt = GetMillis()
	}
}

func DigitalContentUrlIsValid(d model.DigitalContentURL) *AppError {
	if !IsValidId(d.ID) {
		return NewAppError("DigitalContentURL.IsValid", "model.digital_content_url.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(d.Token) {
		return NewAppError("DigitalContentURL.IsValid", "model.digital_content_url.is_valid.token.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(d.ContentID) {
		return NewAppError("DigitalContentURL.IsValid", "model.digital_content_url.is_valid.content_id.app_error", nil, "", http.StatusBadRequest)
	}
	if d.CreatedAt <= 0 {
		return NewAppError("DigitalContentURL.IsValid", "model.digital_content_url.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if d.DownloadNum < 0 {
		return NewAppError("DigitalContentURL.IsValid", "model.digital_content_url.is_valid.download_num.app_error", nil, "", http.StatusBadRequest)
	}
	if !d.LineID.IsNil() && !IsValidId(*d.LineID.String) {
		return NewAppError("DigitalContentURL.IsValid", "model.digital_content_url.is_valid.line_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func ProductChannelListingPreSave(p *model.ProductChannelListing) {
	if p.ID == "" {
		p.ID = NewId()
	}
	if p.CreatedAt == 0 {
		p.CreatedAt = GetMillis()
	}
	ProductChannelListingCommonPre(p)
}

func ProductChannelListingCommonPre(p *model.ProductChannelListing) {
	if p.Currency.IsValid() != nil {
		p.Currency = DEFAULT_CURRENCY
	}
}

func ProductChannelListingIsValid(p model.ProductChannelListing) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.product_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(p.ChannelID) {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.CreatedAt <= 0 {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Currency.IsValid() != nil {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type ProductChannelListingFilterOption struct {
	CommonQueryOptions
}

func ProductMediaPreSave(p *model.ProductMedium) {
	if p.ID == "" {
		p.ID = NewId()
	}
	if p.CreatedAt == 0 {
		p.CreatedAt = GetMillis()
	}
	ProductMediaCommonPre(p)
}

func ProductMediaCommonPre(p *model.ProductMedium) {
	p.Alt = SanitizeUnicode(p.Alt)
	if p.Type.IsValid() != nil {
		p.Type = model.ProductMediaTypeIMAGE
	}
}

func ProductMediaIsValid(p model.ProductMedium) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("ProductMedia.IsValid", "model.product_media.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductMedia.IsValid", "model.product_media.is_valid.product_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.CreatedAt <= 0 {
		return NewAppError("ProductMedia.IsValid", "model.product_media.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Type.IsValid() != nil {
		return NewAppError("ProductMedia.IsValid", "model.product_media.is_valid.type.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type ProductMediaFilterOption struct {
	CommonQueryOptions
	Preloads  []string
	VariantID qm.QueryMod // INNER JOIN VariantMedias ON VariantMedias.MediaID = ProductMedias.Id Where VariantMedias.VariantID ...
}
