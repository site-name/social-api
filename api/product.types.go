package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
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
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	return &embedCtx.CurrentChannelID, nil
}

func (p *Product) AvailableForPurchase(ctx context.Context) (*Date, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		return nil, nil
	}

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
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
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		return nil, nil
	}

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
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

// NOTE: Refer to ./schemas/product.graphqls for details on directives used.
func (p *Product) Collections(ctx context.Context) ([]*Collection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	collections, err := CollectionsByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	requesterIsShopStaff := embedCtx.AppContext.Session().
		GetUserRoles().
		InterSection([]string{model.ShopStaffRoleId, model.ShopAdminRoleId}).
		Len() > 0
	if requesterIsShopStaff {
		return systemRecordsToGraphql(collections, systemCollectionToGraphqlCollection), nil
	}

	keys := lo.Map(collections, func(c *model.Collection, _ int) string { return fmt.Sprintf("%s__%s", c.Id, embedCtx.CurrentChannelID) })
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
		listing := channelListingsDict[c.Id]
		return listing != nil && listing.IsVisible()
	})

	return systemRecordsToGraphql(visibleCollections, systemCollectionToGraphqlCollection), nil
}

// NOTE: Refer to ./schemas/product.graphqls for details on directives used.
func (p *Product) ChannelListings(ctx context.Context) ([]*ProductChannelListing, error) {
	channelListings, err := ProductChannelListingByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(channelListings, systemProductChannelListingToGraphqlProductChannelListing), nil
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
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	taxType, appErr := pluginMng.GetTaxCodeFromObjectMeta(p.p)
	if appErr != nil {
		return nil, appErr
	}

	return &TaxType{
		Description: &taxType.Descriptiton,
		TaxCode:     &taxType.Code,
	}, nil
}

func (p *Product) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*ProductPricingInfo, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		return nil, nil
	}

	discountInfos, err := DiscountsByDateTimeLoader.Load(ctx, time.Now())()
	if err != nil {
		return nil, err
	}

	channel, err := ChannelByIdLoader.Load(ctx, embedCtx.CurrentChannelID)()
	if err != nil {
		return nil, err
	}

	productChannelLiting, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	variants, err := ProductVariantsByProductIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	variantChannelListings, err := VariantsChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
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

	var countryCode model.CountryCode
	if args.Address != nil && args.Address.Country != nil {
		countryCode = *args.Address.Country
	}
	if countryCode == "" {
		countryCode = channel.DefaultCountry
	}

	localCurrency := util.GetCurrencyForCountry(countryCode.String())

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	availability, appErr := embedCtx.App.Srv().ProductService().GetProductAvailability(
		*p.p,
		productChannelLiting,
		variants,
		variantChannelListings,
		collections,
		discountInfos,
		*channel,
		pluginMng,
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

// NOTE: Refer to ./schemas/product.graphqls for details on directives used.
func (p *Product) IsAvailable(ctx context.Context, args struct{ Address *AddressInput }) (*bool, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		return nil, nil
	}

	var countryCode string
	if args.Address != nil && args.Address.Country != nil {
		countryCode = string(*args.Address.Country)
	}

	requesterIsStaffOfShop := embedCtx.AppContext.
		Session().
		GetUserRoles().
		InterSection([]string{model.ShopStaffRoleId, model.ShopAdminRoleId}).
		Len() > 0

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	if productChannelListing != nil && productChannelListing.IsAvailableForPurchase() {
		// check variant availability:
		var variants model.ProductVariants
		if requesterIsStaffOfShop && embedCtx.CurrentChannelID == "" {
			variants, err = ProductVariantsByProductIdLoader.Load(ctx, p.ID)()
		} else if requesterIsStaffOfShop && embedCtx.CurrentChannelID != "" {
			variants, err = ProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
		} else {
			variants, err = AvailableProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
		}
		if err != nil {
			return nil, err
		}

		keys := lo.Map(variants, func(v *model.ProductVariant, _ int) string {
			return fmt.Sprintf("%s__%s__%s", v.Id, countryCode, embedCtx.CurrentChannelID)
		})
		quantities, errs := AvailableQuantityByProductVariantIdCountryCodeAndChannelIDLoader.LoadMany(ctx, keys)()
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

	return systemRecordsToGraphql(medias, systemProductMediaToGraphqlProductMedia), nil
}

func (p *Product) Variants(ctx context.Context) ([]*ProductVariant, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("Product.variants", model.ShopStaffRoleId, model.ShopAdminRoleId)
	requesterIsShopStaff := embedCtx.Err == nil

	var variants model.ProductVariants
	var err error

	if requesterIsShopStaff && embedCtx.CurrentChannelID == "" {
		variants, err = ProductVariantsByProductIdLoader.Load(ctx, p.ID)()
	} else if requesterIsShopStaff && embedCtx.CurrentChannelID != "" {
		variants, err = ProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
	} else {
		variants, err = AvailableProductVariantsByProductIdAndChannel.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
	}
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(variants, SystemProductVariantToGraphqlProductVariant), nil
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
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	taxType, appErr := pluginMng.GetTaxCodeFromObjectMeta(p.p)
	if appErr != nil {
		return nil, appErr
	}
	return &TaxType{
		Description: &taxType.Descriptiton,
		TaxCode:     &taxType.Code,
	}, nil
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

// attributes ORDER BY Slug
func (p *ProductType) AvailableAttributes(ctx context.Context, args struct {
	GraphqlParams
	Filter *AttributeFilterInput
}) (*AttributeCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("ProductType.AvailableAttributes", model.ShopAdminRoleId, model.ShopStaffRoleId)

	filterOptions := &model.AttributeFilterOption{}
	if args.Filter != nil {
		var appErr *model.AppError
		filterOptions, appErr = args.Filter.parse("ProductType.AvailableAttributes")
		if appErr != nil {
			return nil, appErr
		}
	}

	// this is needed
	filterOptions.UserIsShopStaff = embedCtx.Err == nil

	attributes, err := embedCtx.App.Srv().Store.Attribute().GetProductTypeAttributes(p.ID, true, filterOptions)
	if err != nil {
		return nil, model.NewAppError("AvailableAttributes", "app.attribute.filter_product_type_attributes.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	keyFunc := func(a *model.Attribute) []any { return []any{model.AttributeTableName + ".Slug", a.Slug} }
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

	return systemRecordsToGraphql(attributes, SystemAttributeToGraphqlAttribute), nil
}

func (p *ProductType) VariantAttributes(ctx context.Context, args struct{ VariantSelection *VariantAttributeScope }) ([]*Attribute, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	attributes, err := VariantAttributesByProductTypeIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	if args.VariantSelection == nil || *args.VariantSelection == VariantAttributeScopeAll {
		return systemRecordsToGraphql(attributes, SystemAttributeToGraphqlAttribute), nil
	}

	variantSelectionAttributes := embedCtx.App.Srv().ProductService().GetVariantSelectionAttributes(attributes)
	if *args.VariantSelection == VariantAttributeScopeVariantSelection {
		return systemRecordsToGraphql(variantSelectionAttributes, SystemAttributeToGraphqlAttribute), nil
	}

	variantSelectionAttributesMap := lo.SliceToMap(variantSelectionAttributes, func(v *model.Attribute) (string, struct{}) { return v.Id, struct{}{} })
	attributes = lo.Filter(attributes, func(a *model.Attribute, _ int) bool {
		_, exist := variantSelectionAttributesMap[a.Id]
		return !exist
	})

	return systemRecordsToGraphql(attributes, SystemAttributeToGraphqlAttribute), nil
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
