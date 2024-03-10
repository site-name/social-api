package model_helper

import (
	"fmt"
	"net/http"
	"time"

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

type ProductFilterOption struct {
	CommonQueryOptions
	HasNoProductVariants bool
	ProductVariantID     qm.QueryMod
	VoucherID            qm.QueryMod
	SaleID               qm.QueryMod
	CollectionID         qm.QueryMod
	Preloads             []string
}

type ProductCountByCategoryID struct {
	CategoryID   string `json:"category_id"`
	ProductCount uint64 `json:"product_count"`
}

func ProductPreUpdate(p *model.Product) {
	p.UpdatedAt = GetMillis()
	productCommonPre(p)
}

func ProductIsValid(p model.Product) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !IsValidId(p.CategoryID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.category_id.app_error", nil, "please provide valid category id", http.StatusBadRequest)
	}
	if p.CreatedAt <= 0 {
		return NewAppError("Product.IsValid", "model.product.is_valid.created_at.app_error", nil, "please specify create time", http.StatusBadRequest)
	}
	if p.UpdatedAt <= 0 {
		return NewAppError("Product.IsValid", "model.product.is_valid.updated_at.app_error", nil, "please specify update time", http.StatusBadRequest)
	}
	if !p.DefaultVariantID.IsNil() && !IsValidId(*p.DefaultVariantID.String) {
		return NewAppError("Product.IsValid", "model.product.is_valid.default_variant_id.app_error", nil, "please provide valid default variant id", http.StatusBadRequest)
	}
	if p.Name == "" {
		return NewAppError("Product.IsValid", "model.product.is_valid.name.app_error", nil, "please rovide valid name", http.StatusBadRequest)
	}
	if !slug.IsSlug(p.Slug) {
		return NewAppError("Product.IsValid", "model.product.is_valid.slug.app_error", nil, "please provide valid slug", http.StatusBadRequest)
	}
	if measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(p.WeightUnit)] == "" {
		return NewAppError("Product.IsValid", "model.product.is_valid.weight_unit.app_error", nil, "please provide valid weignt unit", http.StatusBadRequest)
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

// NOTE: client code doesn't need to pass in select mods, store code handles that
type CollectionFilterOptions struct {
	CommonQueryOptions
	ProductID                                             qm.QueryMod // INNER JOIN product_cpllections ON ... WHERE product_collections.product_id ...
	VoucherID                                             qm.QueryMod // INNER JOIN voucher_collections ON ... WHERE voucher_collections.voucher_id ...
	SaleID                                                qm.QueryMod // INNER JOIN sale_collections ON ... WHERE sale_collections.sale_id ...
	RelatedCollectionChannelListingConds                  qm.QueryMod // INNER JOIN collection_channel_listings ON ... WHERE collection_channel_listings...
	RelatedCollectionChannelListingChannelConds           qm.QueryMod // INNER JOIN collection_channel_listings ON ... INNER JOIN chanel ON ... WHERE channels...
	AnnotateProductCount                                  bool
	AnnotateIsPublished                                   bool
	AnnotatePublicationDate                               bool
	ChannelSlugForIsPublishedAndPublicationDateAnnotation string
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
	ProductVariantID         qm.QueryMod // INNER JOIN products ON ... INNER JOIN product_variants ON ... WHERE product_variants.id ...
	RelatedChannelConditions qm.QueryMod // INNER JOIN channels ON ... WHERE channels ...
	Preloads                 []string
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

var CustomCollectionColumns = struct {
	ProductCount    string
	IsPublished     string
	PublicationDate string
}{
	ProductCount:    "product_count",
	IsPublished:     "is_published",
	PublicationDate: "publication_date",
}

type CustomCollection struct {
	model.Collection `boil:",bind"`
	ProductCount     *int       `boil:"product_count" json:"product_count"`
	IsPublished      *bool      `boil:"is_published" json:"is_published"`
	PublicationDate  *time.Time `boil:"publication_date" json:"publication_date"`
}

type CustomCollectionSlice []*CustomCollection

type ProductVariantFilterOptions struct {
	CommonQueryOptions
	WishlistItemID qm.QueryMod // INNER JOIN WishlistItemProductVariants ON (...) WHERE WishlistItemProductVariants.WishlistItemID ...
	WishlistID     qm.QueryMod // INNER JOIN WishlistItemProductVariants ON (...) INNER JOIN WishlistItems ON (...) WHERE WishlistItems.WishlistID ...

	VoucherID qm.QueryMod // INNER JOIN voucher_product_variants ON ... WHERE voucher_product_variants.voucher_id ...
	SaleID    qm.QueryMod // INNER JOIN sale_product_variants ON ... WHERE sale_product_variants.sale_id ...

	RelatedProductVariantChannelListingConds qm.QueryMod // INNER JOIN product_variant_channel_listings ON ... WHERE product_variant_channel_listings ...
	ProductVariantChannelListingChannelSlug  qm.QueryMod // INNER JOIN `product_variant_channel_listings` ON ... INNER JOIN channels ON ... WHERE channels.slug ...

	Preloads []string
}

func ProductVariantChannelListingPreSave(p *model.ProductVariantChannelListing) {
	if p.ID == "" {
		p.ID = NewId()
	}
	if p.CreatedAt == 0 {
		p.CreatedAt = GetMillis()
	}
	ProductVariantChannelListingCommonPre(p)
}

func ProductVariantChannelListingCommonPre(p *model.ProductVariantChannelListing) {
	p.Annotations = nil
}

func ProductVariantChannelListingIsValid(p model.ProductVariantChannelListing) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !IsValidId(p.VariantID) {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.variant_id.app_error", nil, "please provide valid variant id", http.StatusBadRequest)
	}
	if !IsValidId(p.ChannelID) {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}
	if p.CreatedAt <= 0 {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.created_at.app_error", nil, "please specify creation time", http.StatusBadRequest)
	}
	if p.Currency.Valid && p.Currency.Val.IsValid() != nil {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}

	return nil
}

type ProductVariantChannelListingFilterOption struct {
	CommonQueryOptions

	VariantProductID qm.QueryMod // INNER JOIN product_variants ON ... WHERE product_variants.product_id ...
	Preloads         []string

	AnnotatePreorderQuantityAllocated bool
	AnnotateAvailablePreorderQuantity bool
}

var ProductVariantChannelListingAnnotationKeys = struct {
	AvailablePreorderQuantity string
	PreorderQuantityAllocated string
}{
	AvailablePreorderQuantity: "available_preorder_quantity",
	PreorderQuantityAllocated: "preorder_quantity_allocated",
}

func PreorderAllocationPreSave(p *model.PreorderAllocation) {
	if p.ID == "" {
		p.ID = NewId()
	}
}

func PreorderAllocationIsValid(p model.PreorderAllocation) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("PreorderAllocation.IsValid", "model.preorder_allocation.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !IsValidId(p.OrderLineID) {
		return NewAppError("PreorderAllocation.IsValid", "model.preorder_allocation.is_valid.order_line_id.app_error", nil, "please provide valid order line id", http.StatusBadRequest)
	}
	if !IsValidId(p.ProductVariantChannelListingID) {
		return NewAppError("PreorderAllocation.IsValid", "model.preorder_allocation.is_valid.product_variant_channel_listing_id.app_error", nil, "please provide valid order line id", http.StatusBadRequest)
	}
	if p.Quantity < 0 {
		return NewAppError("PreorderAllocation.IsValid", "model.preorder_allocation.is_valid.quantity.app_error", nil, "please provide valid quantity", http.StatusBadRequest)
	}

	return nil
}

func ProductVariantString(p model.ProductVariant) string {
	if p.Name != "" {
		return p.Name
	}

	return fmt.Sprintf("ID:%s", p.ID)
}

type ProductTranslationFilterOption struct {
	CommonQueryOptions
}

type ProductVariantTranslationFilterOption struct {
	CommonQueryOptions
}
