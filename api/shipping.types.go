package api

import (
	"context"
	"net/http"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type ShippingMethod struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	Description         JSONString              `json:"description"`
	MinimumOrderWeight  *Weight                 `json:"minimumOrderWeight"`
	MaximumOrderWeight  *Weight                 `json:"maximumOrderWeight"`
	MaximumDeliveryDays *int32                  `json:"maximumDeliveryDays"`
	MinimumDeliveryDays *int32                  `json:"minimumDeliveryDays"`
	PrivateMetadata     []*MetadataItem         `json:"privateMetadata"`
	Metadata            []*MetadataItem         `json:"metadata"`
	Type                *ShippingMethodTypeEnum `json:"type"`

	s *model.ShippingMethod

	// Translation         *ShippingMethodTranslation `json:"translation"`
	// ChannelListings     []*ShippingMethodChannelListing `json:"channelListings"`
	// Price               *Money                          `json:"price"`
	// MaximumOrderPrice   *Money                          `json:"maximumOrderPrice"`
	// MinimumOrderPrice   *Money                          `json:"minimumOrderPrice"`
	// PostalCodeRules     []*ShippingMethodPostalCodeRule `json:"postalCodeRules"`
	// ExcludedProducts    *ProductCountableConnection     `json:"excludedProducts"`
}

func SystemShippingMethodToGraphqlShippingMethod(m *model.ShippingMethod) *ShippingMethod {
	if m == nil {
		return nil
	}

	res := &ShippingMethod{
		ID:              m.Id,
		Name:            m.Name,
		Description:     JSONString(m.Description),
		PrivateMetadata: MetadataToSlice(m.PrivateMetadata),
		Metadata:        MetadataToSlice(m.Metadata),
		Type:            model.NewPrimitive(ShippingMethodTypeEnum(m.Type)),
		s:               m,
		MinimumOrderWeight: &Weight{
			Unit:  WeightUnitsEnum(m.WeightUnit),
			Value: float64(m.MinimumOrderWeight),
		},
	}

	if m.MaximumOrderWeight != nil {
		res.MaximumOrderWeight = &Weight{
			Unit:  WeightUnitsEnum(m.WeightUnit),
			Value: float64(*m.MaximumOrderWeight),
		}
	}

	if m.MaximumDeliveryDays != nil {
		res.MaximumDeliveryDays = model.NewPrimitive(int32(*m.MaximumDeliveryDays))
	}
	if m.MinimumDeliveryDays != nil {
		res.MinimumDeliveryDays = model.NewPrimitive(int32(*m.MinimumDeliveryDays))
	}

	return res
}

func (s *ShippingMethod) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*ShippingMethodTranslation, error) {
	panic("not implemented")
}

func (s *ShippingMethod) ChannelListings(ctx context.Context) ([]*ShippingMethodChannelListing, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	if embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageShipping) {
		listings, err := ShippingMethodChannelListingByShippingMethodIdLoader.Load(ctx, s.ID)()
		if err != nil {
			return nil, err
		}

		return DataloaderResultMap(listings, systemShippingMethodChannelListingToGraphqlShippingMethodChannelListing), nil
	}

	return nil, model.NewAppError("ShippingMethod.ChannelListings", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}

func (s *ShippingMethod) Price(ctx context.Context) (*Money, error) {
	embedChannel, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	if embedChannel == "" {
		return nil, nil
	}

	listing, err := ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader.Load(ctx, s.ID+"__"+embedChannel)()
	if err != nil {
		return nil, err
	}

	if p := listing.PriceAmount; p != nil {
		return &Money{
			Amount:   p.InexactFloat64(),
			Currency: listing.Currency,
		}, nil
	}

	return nil, nil
}

func (s *ShippingMethod) MaximumOrderPrice(ctx context.Context) (*Money, error) {
	embedChannel, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	if embedChannel == "" {
		return nil, nil
	}

	listing, err := ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader.Load(ctx, s.ID+"__"+embedChannel)()
	if err != nil {
		return nil, err
	}

	if p := listing.MaximumOrderPriceAmount; p != nil {
		return &Money{
			Amount:   p.InexactFloat64(),
			Currency: listing.Currency,
		}, nil
	}
	return nil, nil
}

func (s *ShippingMethod) MinimumOrderPrice(ctx context.Context) (*Money, error) {
	embedChannel, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	if embedChannel == "" {
		return nil, nil
	}

	listing, err := ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader.Load(ctx, s.ID+"__"+embedChannel)()
	if err != nil {
		return nil, err
	}

	if p := listing.MinimumOrderPriceAmount; p != nil {
		return &Money{
			Amount:   p.InexactFloat64(),
			Currency: listing.Currency,
		}, nil
	}
	return nil, nil
}

func (s *ShippingMethod) PostalCodeRules(ctx context.Context) ([]*ShippingMethodPostalCodeRule, error) {
	postalCodeRules, err := PostalCodeRulesByShippingMethodIdLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	return lo.Map(postalCodeRules, func(r *model.ShippingMethodPostalCodeRule, _ int) *ShippingMethodPostalCodeRule {
		inclusionType := PostalCodeRuleInclusionTypeEnum(r.InclusionType)

		return &ShippingMethodPostalCodeRule{
			Start:         &r.Start,
			End:           &r.End,
			ID:            r.Id,
			InclusionType: &inclusionType,
		}
	}), nil
}

// NOTE: products are ordered by their slugs
func (s *ShippingMethod) ExcludedProducts(ctx context.Context, args GraphqlParams) (*ProductCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageShipping) {
		return nil, model.NewAppError("ShippingMethod.ExcludedProducts", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

	products, err := ExcludedProductByShippingMethodIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(p *model.Product) string { return p.Slug }
	res, appErr := newGraphqlPaginator(products, keyFunc, SystemProductToGraphqlProduct, args).parse("ShippingMethod.ExcludedProducts")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductCountableConnection)(unsafe.Pointer(res)), nil
}

// ---------------- shipping zone -------------------------

type ShippingZone struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Default         bool              `json:"default"`
	PrivateMetadata []*MetadataItem   `json:"privateMetadata"`
	Metadata        []*MetadataItem   `json:"metadata"`
	Countries       []*CountryDisplay `json:"countries"`
	Description     *string           `json:"description"`

	// PriceRange      *MoneyRange       `json:"priceRange"`
	// ShippingMethods []*ShippingMethod `json:"shippingMethods"`
	// Warehouses      []*Warehouse      `json:"warehouses"`
	// Channels        []*Channel        `json:"channels"`
}

func SystemShippingZoneToGraphqlShippingZone(s *model.ShippingZone) *ShippingZone {
	if s == nil {
		return nil
	}

	res := &ShippingZone{
		ID:              s.Id,
		Name:            s.Name,
		Default:         *s.Default,
		PrivateMetadata: MetadataToSlice(s.PrivateMetadata),
		Metadata:        MetadataToSlice(s.Metadata),
		Description:     &s.Description,
	}

	if s.Countries != "" {
		splitCountries := strings.FieldsFunc(strings.TrimSpace(s.Countries), func(r rune) bool { return r == ' ' || r == ',' })

		for _, code := range splitCountries {
			res.Countries = append(res.Countries, &CountryDisplay{
				Code:    code,
				Country: model.Countries[model.CountryCode(code)],
			})
		}
	}

	return res
}

func (s *ShippingZone) PriceRange(ctx context.Context) (*MoneyRange, error) {
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

	listings, appErr := embedCtx.App.Srv().ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			ChannelID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ChannelID": channelID},
		})
	if appErr != nil {
		return nil, appErr
	}
	if len(listings) == 0 {
		return nil, nil
	}

	var prices = lo.Map(listings, func(l *model.ShippingMethodChannelListing, _ int) *goprices.Money { return l.GetTotal() })
	min, max := util.MinMaxMoneyInMoneySlice(prices)

	return SystemMoneyRangeToGraphqlMoneyRange(&goprices.MoneyRange{
		Start:    min,
		Stop:     max,
		Currency: min.Currency,
	}), nil
}

func (s *ShippingZone) ShippingMethods(ctx context.Context) ([]*ShippingMethod, error) {
	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}

	var shippingMethods []*model.ShippingMethod
	if model.IsValidId(channelID) {
		shippingMethods, err = ShippingMethodsByShippingZoneIdAndChannelSlugLoader.Load(ctx, s.ID+"__"+channelID)()
	} else {
		shippingMethods, err = ShippingMethodsByShippingZoneIdLoader.Load(ctx, s.ID)()
	}

	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(shippingMethods, SystemShippingMethodToGraphqlShippingMethod), nil
}

func (s *ShippingZone) Warehouses(ctx context.Context) ([]*Warehouse, error) {
	warehouses, err := WarehousesByShippingZoneIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(warehouses, SystemWarehouseToGraphqlWarehouse), nil
}

func (s *ShippingZone) Channels(ctx context.Context) ([]*Channel, error) {
	channels, err := ChannelsByShippingZoneIdLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(channels, SystemChannelToGraphqlChannel), nil
}

// ------------------
type ShippingMethodChannelListing struct {
	ID                string `json:"id"`
	MinimumOrderPrice *Money `json:"minimumOrderPrice"`
	MaximumOrderPrice *Money `json:"maximumOrderPrice"`
	Price             *Money `json:"price"`
	s                 *model.ShippingMethodChannelListing

	// Channel           *Channel `json:"channel"`
}

func systemShippingMethodChannelListingToGraphqlShippingMethodChannelListing(s *model.ShippingMethodChannelListing) *ShippingMethodChannelListing {
	if s == nil {
		return nil
	}

	s.PopulateNonDbFields()

	res := &ShippingMethodChannelListing{
		ID: s.Id,
		s:  s,
	}
	if p := s.MinimumOrderPrice; p != nil {
		res.MinimumOrderPrice = SystemMoneyToGraphqlMoney(p)
	}
	if p := s.MaximumOrderPrice; p != nil {
		res.MaximumOrderPrice = SystemMoneyToGraphqlMoney(p)
	}
	if p := s.Price; p != nil {
		res.Price = SystemMoneyToGraphqlMoney(p)
	}
	return res
}

func (s *ShippingMethodChannelListing) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, s.s.ChannelID)()
	if err != nil {
		return nil, err
	}
	return SystemChannelToGraphqlChannel(channel), nil
}
