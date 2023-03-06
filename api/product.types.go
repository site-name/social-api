package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Product struct {
	ID              string          `json:"id"`
	SeoTitle        *string         `json:"seoTitle"`
	SeoDescription  *string         `json:"seoDescription"`
	Name            string          `json:"name"`
	Description     JSONString      `json:"description"`
	Slug            string          `json:"slug"`
	UpdatedAt       *DateTime       `json:"updatedAt"`
	ChargeTaxes     bool            `json:"chargeTaxes"`
	Weight          *Weight         `json:"weight"`
	Rating          *float64        `json:"rating"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	p               *model.Product

	// Channel         *string         `json:"channel"`
	// AvailableForPurchase   *Date           `json:"availableForPurchase"`
	// IsAvailableForPurchase *bool           `json:"isAvailableForPurchase"`
	// Attributes             []*SelectedAttribute `json:"attributes"`
	// Pricing                *ProductPricingInfo  `json:"pricing"`
	// TaxType                *TaxType             `json:"taxType"`
	// IsAvailable            *bool                `json:"isAvailable"`
	// Thumbnail              *Image               `json:"thumbnail"`
	// ChannelListings        []*ProductChannelListing `json:"channelListings"`
	// MediaByID              *ProductMedia            `json:"mediaById"`
	// Variants               []*ProductVariant        `json:"variants"`
	// Media                  []*ProductMedia          `json:"media"`
	// Collections            []*Collection            `json:"collections"`
	// Translation            *ProductTranslation      `json:"translation"`
	// DefaultVariant         *ProductVariant          `json:"defaultVariant"`
	// ProductType            *ProductType             `json:"productType"`
	// Category               *Category                `json:"category"`
}

func SystemProductToGraphqlProduct(p *model.Product) *Product {
	if p == nil {
		return nil
	}

	res := &Product{
		ID:              p.Id,
		SeoTitle:        &p.SeoTitle,
		SeoDescription:  &p.SeoDescription,
		Name:            p.Name,
		Slug:            p.Slug,
		ChargeTaxes:     *p.ChargeTaxes,
		Description:     JSONString(p.Description),
		Metadata:        MetadataToSlice(p.Metadata),
		PrivateMetadata: MetadataToSlice(p.PrivateMetadata),
		p:               p,
	}
	if p.Rating != nil {
		res.Rating = model.NewPrimitive(float64(*p.Rating))
	}
	if p.Weight != nil {
		res.Weight = &Weight{
			Unit:  WeightUnitsEnum(p.WeightUnit),
			Value: float64(*p.Weight),
		}
	}

	return res
}

func (p *Product) Channel(ctx context.Context) (*string, error) {
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}

	return &channelID, nil
}

func (p *Product) AvailableForPurchase(ctx context.Context) (*Date, error) {
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	if channelID == "" {
		return nil, nil
	}

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
	if err != nil {
		return nil, err
	}
	if productChannelListing == nil {
		return nil, nil
	}
	if productChannelListing.AvailableForPurchase == nil {
		return nil, nil
	}

	return &Date{DateTime{*productChannelListing.AvailableForPurchase}}, nil
}

func (p *Product) IsAvailableForPurchase(ctx context.Context) (*bool, error) {
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	if channelID == "" {
		return nil, nil
	}

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
	if err != nil {
		return nil, err
	}
	if productChannelListing == nil {
		return nil, nil
	}

	available := productChannelListing.IsAvailableForPurchase()
	return &available, nil
}

func (p *Product) ProductType(ctx context.Context) (*ProductType, error) {
	productType, err := ProductTypeByIdLoader.Load(ctx, p.p.ProductTypeID)()
	if err != nil {
		return nil, err
	}

	return SystemProductTypeToGraphqlProductType(productType), nil
}

func (p *Product) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*ProductTranslation, error) {
	panic("not implemented")
}

func (p *Product) Collections(ctx context.Context) ([]*Collection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}

	hasProductPermission := embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.ProductPermissions...)

	collections, err := CollectionsByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	if hasProductPermission {
		return DataloaderResultMap(collections, systemCollectionToGraphqlCollection), nil
	}

	keys := lo.Map(collections, func(c *model.Collection, _ int) string { return fmt.Sprintf("%s__%s", c.Id, channelID) })
	collectionChannelListings, errs := CollectionChannelListingByCollectionIdAndChannelSlugLoader.LoadMany(ctx, keys)()
	if len(errs) > 0 && errs[0] != nil {
		return nil, errs[0]
	}

	channelListingsDict := map[string]*model.CollectionChannelListing{} // keys are collection ids
	for _, listing := range collectionChannelListings {
		if listing != nil {
			channelListingsDict[listing.CollectionID] = listing
		}
	}

	visibleCollections := lo.Filter(collections, func(c *model.Collection, _ int) bool {
		listing, ok := channelListingsDict[c.Id]
		return ok && listing != nil && listing.IsVisible()
	})

	return DataloaderResultMap(visibleCollections, systemCollectionToGraphqlCollection), nil
}

func (p *Product) ChannelListings(ctx context.Context) ([]*ProductChannelListing, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		return nil, model.NewAppError("Product.ChannelListings", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	channelListings, err := ProductChannelListingByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(channelListings, systemProductChannelListingToGraphqlProductChannelListing), nil
}

func (p *Product) Thumbnail(ctx context.Context, args struct{ Size *int32 }) (*Image, error) {
	panic("not implemented")
}

func (p *Product) DefaultVariant(ctx context.Context) (*ProductVariant, error) {
	if p.p.DefaultVariantID == nil {
		return nil, nil
	}

	variant, err := ProductVariantByIdLoader.Load(ctx, *p.p.DefaultVariantID)()
	if err != nil {
		return nil, err
	}

	return SystemProductVariantToGraphqlProductVariant(variant), nil
}

func (p *Product) Category(ctx context.Context) (*Category, error) {
	if p.p.CategoryID == nil {
		return nil, nil
	}

	category, err := CategoryByIdLoader.Load(ctx, *p.p.CategoryID)()
	if err != nil {
		return nil, err
	}

	return systemCategoryToGraphqlCategory(category), nil
}

func (p *Product) TaxType(ctx context.Context) (*TaxType, error) {
	panic("not implemented")
}

func (p *Product) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*ProductPricingInfo, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	if channelID == "" {
		return nil, nil
	}

	discountInfos, err := DiscountsByDateTimeLoader.Load(ctx, time.Now())()
	if err != nil {
		return nil, err
	}

	channel, err := ChannelByIdLoader.Load(ctx, channelID)()
	if err != nil {
		return nil, err
	}

	productChannelLiting, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
	if err != nil {
		return nil, err
	}

	variants, err := ProductVariantsByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	variantChannelListings, err := VariantsChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
	if err != nil {
		return nil, err
	}
	if len(variantChannelListings) == 0 {
		return nil, nil
	}

	collections, err := CollectionsByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	var countryCode string
	if args.Address != nil && args.Address.Country != nil {
		countryCode = string(*args.Address.Country)
	}
	if countryCode == "" {
		countryCode = channel.DefaultCountry
	}

	localCurrency := util.GetCurrencyForCountry(countryCode)

	panic("not implemented") // TODO: complete plugin manager before removeing this panic

	availability, appErr := embedCtx.App.Srv().ProductService().GetProductAvailability(
		*p.p,
		productChannelLiting,
		variants,
		variantChannelListings,
		collections,
		discountInfos,
		*channel,
		nil,
		countryCode,
		localCurrency,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductPricingInfo{
		OnSale:                  &availability.OnSale,
		Discount:                SystemTaxedMoneyToGraphqlTaxedMoney(availability.Discount),
		DiscountLocalCurrency:   SystemTaxedMoneyToGraphqlTaxedMoney(availability.DiscountLocalCurrency),
		PriceRange:              SystemTaxedMoneyRangeToGraphqlTaxedMoneyRange(availability.PriceRange),
		PriceRangeUndiscounted:  SystemTaxedMoneyRangeToGraphqlTaxedMoneyRange(availability.PriceRangeUnDiscounted),
		PriceRangeLocalCurrency: SystemTaxedMoneyRangeToGraphqlTaxedMoneyRange(availability.PriceRangeLocalCurrency),
	}, nil
}

func (p *Product) IsAvailable(ctx context.Context, args struct{ Address *AddressInput }) (*bool, error) {
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	var countryCode string
	if args.Address != nil && args.Address.Country != nil {
		countryCode = string(*args.Address.Country)
	}

	hasRequiredPermission := embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.ProductPermissions...)

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
	if err != nil {
		return nil, err
	}

	if productChannelListing != nil && productChannelListing.IsAvailableForPurchase() {
		// check variant availability:
		var variants model.ProductVariants
		if hasRequiredPermission && channelID == "" {
			variants, err = ProductVariantsByProductIdLoader.Load(ctx, p.ID)()
		} else if hasRequiredPermission && channelID != "" {
			variants, err = ProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
		} else {
			variants, err = AvailableProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
		}
		if err != nil {
			return nil, err
		}

		keys := lo.Map(variants, func(v *model.ProductVariant, _ int) string {
			return fmt.Sprintf("%s__%s__%s", v.Id, countryCode, channelID)
		})
		quantities, errs := AvailableQuantityByProductVariantIdCountryCodeAndChannelSlugLoader.LoadMany(ctx, keys)()
		if len(errs) > 0 && errs[0] != nil {
			return nil, errs[0]
		}

		res := lo.SomeBy(quantities, func(v int) bool { return v > 0 })
		return &res, nil
	}

	return model.NewPrimitive(false), nil
}

func (p *Product) Attributes(ctx context.Context) ([]*SelectedAttribute, error) {
	return SelectedAttributesByProductIdLoader.Load(ctx, p.ID)()
}

func (p *Product) MediaByID(ctx context.Context, args struct{ Id string }) (*ProductMedia, error) {
	media, err := ProductMediaByIdLoader.Load(ctx, args.Id)()
	if err != nil {
		return nil, err
	}

	return systemProductMediaToGraphqlProductMedia(media), nil
}

func (p *Product) Media(ctx context.Context) ([]*ProductMedia, error) {
	medias, err := MediaByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(medias, systemProductMediaToGraphqlProductMedia), nil
}

func (p *Product) Variants(ctx context.Context) ([]*ProductVariant, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}

	hasProductPermissions := embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.ProductPermissions...)

	var variants model.ProductVariants

	if hasProductPermissions && channelID == "" {
		variants, err = ProductVariantsByProductIdLoader.Load(ctx, p.ID)()
	} else if hasProductPermissions && channelID != "" {
		variants, err = ProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
	} else {
		variants, err = AvailableProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channelID))()
	}
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(variants, SystemProductVariantToGraphqlProductVariant), nil
}

type ProductType struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Slug               string              `json:"slug"`
	HasVariants        bool                `json:"hasVariants"`
	IsShippingRequired bool                `json:"isShippingRequired"`
	IsDigital          bool                `json:"isDigital"`
	PrivateMetadata    []*MetadataItem     `json:"privateMetadata"`
	Metadata           []*MetadataItem     `json:"metadata"`
	Kind               ProductTypeKindEnum `json:"kind"`
	p                  *model.ProductType

	// Weight              *Weight                       `json:"weight"`
	// TaxType             *TaxType                      `json:"taxType"`
	// VariantAttributes   []*Attribute                  `json:"variantAttributes"`
	// ProductAttributes   []*Attribute                  `json:"productAttributes"`
	// AvailableAttributes *AttributeCountableConnection `json:"availableAttributes"`
}

func SystemProductTypeToGraphqlProductType(t *model.ProductType) *ProductType {
	if t == nil {
		return nil
	}

	res := &ProductType{
		ID:              t.Id,
		Name:            t.Name,
		Slug:            t.Slug,
		Metadata:        MetadataToSlice(t.Metadata),
		PrivateMetadata: MetadataToSlice(t.PrivateMetadata),
		Kind:            ProductTypeKindEnum(t.Kind),
		p:               t,
	}

	if t.HasVariants != nil {
		res.HasVariants = *t.HasVariants
	}
	if t.IsShippingRequired != nil {
		res.IsShippingRequired = *t.IsShippingRequired
	}
	if t.IsDigital != nil {
		res.IsDigital = *t.IsDigital
	}
	return res
}

func (p *ProductType) TaxType(ctx context.Context) (*TaxType, error) {
	panic("not implemented")
}

func (p *ProductType) Weight(ctx context.Context) (*Weight, error) {
	if p.p.Weight == nil {
		return nil, nil
	}

	return &Weight{
		Value: float64(*p.p.Weight),
		Unit:  WeightUnitsEnum(p.p.WeightUnit),
	}, nil
}

// ORDER BY Slug
func (p *ProductType) AvailableAttributes(ctx context.Context, args struct {
	GraphqlParams
	Filter *AttributeFilterInput
}) (*AttributeCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		return nil, model.NewAppError("ProductType.AvailableAttributes", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	userHasOneOfProductPermission := embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.ProductPermissions...)

	filterOptions := &model.AttributeFilterOption{}
	if args.Filter != nil {
		filterOptions = args.Filter.ToAttributeFilterOption()
	}
	// NOTE: this permission check is necessary
	filterOptions.UserHasOneOfProductPermissions = &userHasOneOfProductPermission

	attributes, err := embedCtx.App.Srv().Store.Attribute().GetProductTypeAttributes(p.ID, true, filterOptions)
	if err != nil {
		return nil, model.NewAppError("GetProductTypeAttributes", "app.attribute.unassigned_product_type_attributes.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	keyFunc := func(a *model.Attribute) string { return a.Slug }
	res, appErr := newGraphqlPaginator(attributes, keyFunc, SystemAttributeToGraphqlAttribute, args.GraphqlParams).parse("ProductType.AvailableAttributes")
	if appErr != nil {
		return nil, appErr
	}

	return (*AttributeCountableConnection)(unsafe.Pointer(res)), nil
}

func (p *ProductType) ProductAttributes(ctx context.Context) ([]*Attribute, error) {
	attributes, err := ProductAttributesByProductTypeIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(attributes, SystemAttributeToGraphqlAttribute), nil
}

func (p *ProductType) VariantAttributes(ctx context.Context, args struct{ VariantSelection *VariantAttributeScope }) ([]*Attribute, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	attributes, err := VariantAttributesByProductTypeIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	if args.VariantSelection == nil || *args.VariantSelection == VariantAttributeScopeAll {
		return DataloaderResultMap(attributes, SystemAttributeToGraphqlAttribute), nil
	}

	variantSelectionAttributes := embedCtx.App.Srv().ProductService().GetVariantSelectionAttributes(attributes)
	if *args.VariantSelection == VariantAttributeScopeVariantSelection {
		return DataloaderResultMap(variantSelectionAttributes, SystemAttributeToGraphqlAttribute), nil
	}

	variantSelectionAttributesMap := lo.SliceToMap(variantSelectionAttributes, func(v *model.Attribute) (string, struct{}) { return v.Id, struct{}{} })
	attributes = lo.Filter(attributes, func(a *model.Attribute, _ int) bool {
		_, exist := variantSelectionAttributesMap[a.Id]
		return !exist
	})

	return DataloaderResultMap(attributes, SystemAttributeToGraphqlAttribute), nil
}

// -------------------- collection -----------------

type Collection struct {
	ID              string          `json:"id"`
	SeoTitle        *string         `json:"seoTitle"`
	SeoDescription  *string         `json:"seoDescription"`
	Name            string          `json:"name"`
	Description     JSONString      `json:"description"`
	Slug            string          `json:"slug"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`

	// Channel         *string         `json:"channel"`
	// Products        *ProductCountableConnection `json:"products"`
	// BackgroundImage *Image                      `json:"backgroundImage"`
	// Translation     *CollectionTranslation      `json:"translation"`
	// ChannelListings []*CollectionChannelListing `json:"channelListings"`
}

func systemCollectionToGraphqlCollection(c *model.Collection) *Collection {
	if c == nil {
		return nil
	}

	return &Collection{
		ID:              c.Id,
		SeoTitle:        &c.SeoTitle,
		SeoDescription:  &c.SeoDescription,
		Name:            c.Name,
		Slug:            c.Slug,
		Description:     JSONString(c.Description),
		Metadata:        MetadataToSlice(c.Metadata),
		PrivateMetadata: MetadataToSlice(c.PrivateMetadata),
	}
}

func (c *Collection) Channel(ctx context.Context) (*string, error) {
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, nil
	}

	return &channelID, nil
}

// TODO: add support filtering
func (c *Collection) Products(ctx context.Context, args struct {
	// Filter *ProductFilterInput
	// SortBy *ProductOrder
	GraphqlParams
}) (*ProductCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil && errors.Is(err, ErrorUnExpectedType) {
		return nil, err
	}

	products, appErr := embedCtx.App.Srv().ProductService().GetVisibleProductsToUser(embedCtx.AppContext.Session(), channelID)
	if appErr != nil {
		return nil, appErr
	}

	// filter to get products that belong to current collection:
	collectionProductRelations, appErr := embedCtx.App.Srv().ProductService().CollectionProductRelationsByOptions(&model.CollectionProductFilterOptions{
		CollectionID: squirrel.Eq{store.CollectionProductRelationTableName + ".CollectionID": c.ID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// keys are product ids
	validProductIdMap := lo.SliceToMap(collectionProductRelations, func(rel *model.CollectionProduct) (string, struct{}) { return rel.ProductID, struct{}{} })
	products = lo.Filter(products, func(p *model.Product, _ int) bool {
		_, exist := validProductIdMap[p.Id]
		return exist
	})

	keyFunc := func(p *model.Product) string { return p.Slug }
	res, appErr := newGraphqlPaginator(products, keyFunc, SystemProductToGraphqlProduct, args.GraphqlParams).parse("Collection.Products")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductCountableConnection)(unsafe.Pointer(res)), nil
}

func (c *Collection) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*CollectionTranslation, error) {
	panic("not implemented")
}

func (c *Collection) BackgroundImage(ctx context.Context, args struct{ Size *int32 }) (*Image, error) {
	panic("not implemented")
}

func (c *Collection) ChannelListings(ctx context.Context) ([]*CollectionChannelListing, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		return nil, model.NewAppError("Collection.ChannelListings", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	listings, err := CollectionChannelListingByCollectionIdLoader.Load(ctx, c.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(listings, systemCollectionChannelListingToGraphqlCollectionChannelListing), nil
}

// ------- ProductMedia

type ProductMedia struct {
	ID         string           `json:"id"`
	SortOrder  *int32           `json:"sortOrder"`
	Alt        string           `json:"alt"`
	Type       ProductMediaType `json:"type"`
	OembedData JSONString       `json:"oembedData"`
	// URL        string           `json:"url"`
}

func systemProductMediaToGraphqlProductMedia(p *model.ProductMedia) *ProductMedia {
	if p == nil {
		return nil
	}

	res := &ProductMedia{
		ID:   p.Id,
		Alt:  p.Alt,
		Type: ProductMediaType(p.Type),
	}
	if p.OembedData != nil {
		res.OembedData = JSONString(p.OembedData)
	}
	if s := p.SortOrder; s != nil {
		res.SortOrder = model.NewPrimitive(int32(*s))
	}

	return res
}

func (p *ProductMedia) URL(ctx context.Context, args struct{ Size *int32 }) (string, error) {
	panic("not implemented")
}
