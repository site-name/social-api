package api

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

func discountsByDateTimeLoader(ctx context.Context, dateTimes []time.Time) []*dataloader.Result[[]*model_helper.DiscountInfo] {
	var (
		res             = make([]*dataloader.Result[[]*model_helper.DiscountInfo], len(dateTimes))
		salesMap        = map[time.Time]model.Sales{}
		saleIDS         []string
		channelListings = map[string]map[string]*model.SaleChannelListing{}
		products        = map[string][]string{}
		categories      = map[string][]string{}
		variants        = map[string][]string{}
		collections     map[string][]string
		appErr          *model_helper.AppError
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	discountService := embedCtx.App.Srv().DiscountService()

	for _, dateTime := range dateTimes {
		sales, appErr := discountService.ActiveSales(&dateTime)
		if appErr != nil {
			goto errorLabel
		}

		for _, sale := range sales {
			saleIDS = append(saleIDS, sale.Id)
		}

		salesMap[dateTime] = sales
	}

	collections, appErr = discountService.FetchCollections(saleIDS)
	if appErr != nil {
		goto errorLabel
	}

	channelListings, appErr = discountService.FetchSaleChannelListings(saleIDS)
	if appErr != nil {
		goto errorLabel
	}

	products, appErr = discountService.FetchProducts(saleIDS)
	if appErr != nil {
		goto errorLabel
	}

	categories, appErr = discountService.FetchCategories(saleIDS)
	if appErr != nil {
		goto errorLabel
	}

	variants, appErr = discountService.FetchVariants(saleIDS)
	if appErr != nil {
		goto errorLabel
	}

	for i, datetime := range dateTimes {
		items := make([]*model_helper.DiscountInfo, len(salesMap[datetime]))

		for idx, sale := range salesMap[datetime] {
			items[idx] = &model.DiscountInfo{
				Sale:            sale,
				ChannelListings: channelListings[sale.Id],
				CategoryIDs:     categories[sale.Id],
				CollectionIDs:   collections[sale.Id],
				ProductIDs:      products[sale.Id],
				VariantsIDs:     variants[sale.Id],
			}
		}

		res[i] = &dataloader.Result[[]*model_helper.DiscountInfo]{Data: items}
	}

	return res

errorLabel:
	for i := range dateTimes {
		res[i] = &dataloader.Result[[]*model_helper.DiscountInfo]{Error: appErr}
	}
	return res
}

// NOTE: saleIDChannelIDPairs are strings with format of saleID__channelID.
func saleChannelListingBySaleIdAndChanneSlugLoader(ctx context.Context, saleIDChannelIDPairs []string) []*dataloader.Result[*model.SaleChannelListing] {
	var (
		res        = make([]*dataloader.Result[*model.SaleChannelListing], len(saleIDChannelIDPairs))
		saleIDs    []string
		channelIDs []string
	)

	for _, item := range saleIDChannelIDPairs {
		index := strings.Index(item, "__")
		if index >= 0 {
			saleIDs = append(saleIDs, item[:index])
			channelIDs = append(channelIDs, item[index+2:])
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	listings, appErr := embedCtx.App.Srv().
		DiscountService().
		SaleChannelListingsByOptions(&model.SaleChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.SaleChannelListingTableName + ".SaleID":    saleIDs,
				model.SaleChannelListingTableName + ".ChannelID": channelIDs,
			},
		})
	if appErr != nil {
		for idx := range saleIDChannelIDPairs {
			res[idx] = &dataloader.Result[*model.SaleChannelListing]{Error: appErr}
		}
		return res
	}

	listingMap := lo.SliceToMap(listings, func(l *model.SaleChannelListing) (string, *model.SaleChannelListing) {
		return l.SaleID + "__" + l.ChannelID, l
	})

	for idx, pair := range saleIDChannelIDPairs {
		res[idx] = &dataloader.Result[*model.SaleChannelListing]{Data: listingMap[pair]}
	}
	return res
}

func saleChannelListingBySaleIdLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.SaleChannelListing] {
	var (
		res        = make([]*dataloader.Result[[]*model.SaleChannelListing], len(saleIDs))
		listingMap = map[string][]*model.SaleChannelListing{} // keys are sale ids
	)
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	listings, appErr := embedCtx.App.Srv().
		DiscountService().
		SaleChannelListingsByOptions(&model.SaleChannelListingFilterOption{
			Conditions: squirrel.Eq{model.SaleChannelListingTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		for idx := range saleIDs {
			res[idx] = &dataloader.Result[[]*model.SaleChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range listings {
		listingMap[listing.SaleID] = append(listingMap[listing.SaleID], listing)
	}

	for idx, saleID := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.SaleChannelListing]{Data: listingMap[saleID]}
	}
	return res
}

// ------------------------- order discount --------------------

type OrderDiscount struct {
	ID             string                `json:"id"`
	Type           OrderDiscountType     `json:"type"`
	ValueType      DiscountValueTypeEnum `json:"valueType"`
	Value          PositiveDecimal       `json:"value"`
	Name           *string               `json:"name"`
	TranslatedName *string               `json:"translatedName"`
	Amount         *Money                `json:"amount"`
	reason         *string
}

func SystemOrderDiscountToGraphqlOrderDiscount(r *model.OrderDiscount) *OrderDiscount {
	if r == nil {
		return nil
	}

	return &OrderDiscount{
		ID:             r.Id,
		Type:           OrderDiscountType(r.Type),
		ValueType:      DiscountValueTypeEnum(r.ValueType),
		Value:          PositiveDecimal(*r.Value),
		Name:           r.Name,
		TranslatedName: r.TranslatedName,
		reason:         r.Reason,
		Amount:         SystemMoneyToGraphqlMoney(r.Amount),
	}
}

func (o *OrderDiscount) Reason(ctx context.Context) (*string, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionReadOrderDiscount, model.PermissionCreateOrderDiscount, model.PermissionUpdateOrderDiscount)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	return o.reason, nil
}

func orderDiscountsByOrderIDLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.OrderDiscount] {
	var (
		res              = make([]*dataloader.Result[[]*model.OrderDiscount], len(orderIDs))
		orderDiscountMap = map[string]model.OrderDiscounts{} // keys are order ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	orderDiscounts, appErr := embedCtx.App.Srv().DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		Conditions: squirrel.Eq{model.OrderDiscountTableName + ".OrderID": orderIDs},
	})
	if appErr != nil {
		for idx := range orderIDs {
			res[idx] = &dataloader.Result[[]*model.OrderDiscount]{Error: appErr}
		}
		return res
	}

	for _, rel := range orderDiscounts {
		if rel.OrderID == nil {
			continue
		}
		orderDiscountMap[*rel.OrderID] = append(orderDiscountMap[*rel.OrderID], rel)
	}

	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.OrderDiscount]{Data: orderDiscountMap[id]}
	}
	return res
}

// ------------------------- voucher --------------------
type Voucher struct {
	ID                       string                `json:"id"`
	Name                     *string               `json:"name"`
	Type                     VoucherTypeEnum       `json:"type"`
	Code                     string                `json:"code"`
	UsageLimit               *int32                `json:"usageLimit"`
	Used                     int32                 `json:"used"`
	StartDate                DateTime              `json:"startDate"`
	EndDate                  *DateTime             `json:"endDate"`
	ApplyOncePerOrder        bool                  `json:"applyOncePerOrder"`
	ApplyOncePerCustomer     bool                  `json:"applyOncePerCustomer"`
	DiscountValueType        DiscountValueTypeEnum `json:"discountValueType"`
	MinCheckoutItemsQuantity *int32                `json:"minCheckoutItemsQuantity"`
	PrivateMetadata          []*MetadataItem       `json:"privateMetadata"`
	Metadata                 []*MetadataItem       `json:"metadata"`
	Countries                []*CountryDisplay     `json:"countries"`

	// Translation              *VoucherTranslation   `json:"translation"`
	// Categories               *CategoryCountableConnection       `json:"categories"`
	// Collections              *CollectionCountableConnection     `json:"collections"`
	// Products                 *ProductCountableConnection        `json:"products"`
	// Variants                 *ProductVariantCountableConnection `json:"variants"`
	// DiscountValue            *float64                           `json:"discountValue"`
	// Currency                 *string                            `json:"currency"`
	// MinSpent                 *Money                             `json:"minSpent"`
	// ChannelListings          []*VoucherChannelListing           `json:"channelListings"`
}

func systemVoucherToGraphqlVoucher(v *model.Voucher) *Voucher {
	if v == nil {
		return nil
	}

	res := &Voucher{
		ID:                       v.Id,
		Name:                     v.Name,
		Type:                     VoucherTypeEnum(v.Type),
		Code:                     v.Code,
		Used:                     int32(v.Used),
		StartDate:                DateTime{v.StartDate},
		ApplyOncePerOrder:        v.ApplyOncePerOrder,
		ApplyOncePerCustomer:     v.ApplyOncePerCustomer,
		DiscountValueType:        DiscountValueTypeEnum(v.DiscountValueType),
		MinCheckoutItemsQuantity: model_helper.GetPointerOfValue(int32(v.MinCheckoutItemsQuantity)),
		Metadata:                 MetadataToSlice(v.Metadata),
		PrivateMetadata:          MetadataToSlice(v.PrivateMetadata),
	}

	countries := strings.Fields(v.Countries)
	if len(countries) > 0 {
		for _, code := range countries {
			res.Countries = append(res.Countries, &CountryDisplay{
				Code:    code,
				Country: model.Countries[model.CountryCode(code)],
			})
		}
	}

	if v.EndDate != nil {
		res.EndDate = &DateTime{*v.EndDate}
	}
	if v.UsageLimit != nil {
		res.UsageLimit = model_helper.GetPointerOfValue(int32(*v.UsageLimit))
	}

	return res
}

func (v *Voucher) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*VoucherTranslation, error) {
	panic("not implemented")
}

// categories are order by names
func (v *Voucher) Categories(ctx context.Context, args GraphqlParams) (*CategoryCountableConnection, error) {
	categories, err := CategoriesByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(c *model.Category) []any { return []any{model.CategoryTableName + ".Name", c.Name} }
	res, appErr := newGraphqlPaginator(categories, keyFunc, systemCategoryToGraphqlCategory, args).parse("Voucher.Categories")
	if appErr != nil {
		return nil, appErr
	}

	return (*CategoryCountableConnection)(unsafe.Pointer(res)), nil
}

// collections order by slugs
func (v *Voucher) Collections(ctx context.Context, args GraphqlParams) (*CollectionCountableConnection, error) {
	collections, err := CollectionsByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(c *model.Collection) []any { return []any{model.CollectionTableName + ".Slug", c.Slug} }
	res, appErr := newGraphqlPaginator(collections, keyFunc, systemCollectionToGraphqlCollection, args).parse("Voucher.Collections")
	if err != nil {
		return nil, appErr
	}

	return (*CollectionCountableConnection)(unsafe.Pointer(res)), nil
}

func (v *Voucher) Products(ctx context.Context, args GraphqlParams) (*ProductCountableConnection, error) {
	products, err := ProductsByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(p *model.Product) []any { return []any{model.ProductTableName + ".Slug", p.Slug} }
	res, appErr := newGraphqlPaginator(products, keyFunc, SystemProductToGraphqlProduct, args).parse("voucher.Products")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductCountableConnection)(unsafe.Pointer(res)), nil
}

func (v *Voucher) Variants(ctx context.Context, args GraphqlParams) (*ProductVariantCountableConnection, error) {
	variants, err := ProductVariantsByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(pv *model.ProductVariant) []any { return []any{model.ProductVariantTableName + ".Sku", pv.Sku} }
	res, appErr := newGraphqlPaginator(variants, keyFunc, SystemProductVariantToGraphqlProductVariant, args).parse("Voucher.Variants")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductVariantCountableConnection)(unsafe.Pointer(res)), nil
}

func (v *Voucher) DiscountValue(ctx context.Context) (*float64, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		embedCtx.SetInvalidUrlParam("channel_id")
		return nil, embedCtx.Err
	}

	voucherChannelListing, err := VoucherChannelListingByVoucherIdAndChanneSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", v.ID, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	res := voucherChannelListing.DiscountValue.InexactFloat64()
	return &res, nil
}

func (v *Voucher) Currency(ctx context.Context) (*string, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		embedCtx.SetInvalidUrlParam("channel_id")
		return nil, embedCtx.Err
	}

	voucherChannelListing, err := VoucherChannelListingByVoucherIdAndChanneSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", v.ID, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	return &voucherChannelListing.Currency, nil
}

func (v *Voucher) MinSpent(ctx context.Context) (*Money, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		embedCtx.SetInvalidUrlParam("channel_id")
		return nil, embedCtx.Err
	}

	voucherChannelListing, err := VoucherChannelListingByVoucherIdAndChanneSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", v.ID, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	return SystemMoneyToGraphqlMoney(voucherChannelListing.MinSpent), nil
}

func (v *Voucher) ChannelListings(ctx context.Context) ([]*VoucherChannelListing, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionReadVoucher, model.PermissionUpdateVoucher, model.PermissionDeleteVoucher)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	listings, err := VoucherChannelListingByVoucherIdLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(listings, systemVoucherChannelListingToGraphqlVoucherChannelListing), nil
}

func voucherByIDLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Voucher] {
	var (
		res        = make([]*dataloader.Result[*model.Voucher], len(ids))
		voucherMap = map[string]*model.Voucher{}
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, vouchers, appErr := embedCtx.App.Srv().DiscountService().VouchersByOption(&model.VoucherFilterOption{
		Conditions: squirrel.Eq{model.VoucherTableName + ".Id": ids},
	})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Voucher]{Error: appErr}
		}
		return res
	}

	for _, v := range vouchers {
		voucherMap[v.Id] = v
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Voucher]{Data: voucherMap[id]}
	}
	return res
}

// NOTE: idPairs contains strings with format of voucherID__channelID
func voucherChannelListingByVoucherIdAndChanneSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.VoucherChannelListing] {
	var (
		res                      = make([]*dataloader.Result[*model.VoucherChannelListing], len(idPairs))
		voucherChannelListingMap = map[string]*model.VoucherChannelListing{} // keys are voucher channel listing ids

		voucherIDs []string
		channelIDs []string
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	for _, pair := range idPairs {
		if index := strings.Index(pair, "__"); index >= 0 {
			voucherIDs = append(voucherIDs, pair[:index])
			channelIDs = append(channelIDs, pair[index+2:])
		}
	}

	voucherChannelListings, appErr := embedCtx.App.Srv().DiscountService().
		VoucherChannelListingsByOption(&model.VoucherChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.VoucherChannelListingTableName + ".VoucherID": voucherIDs,
				model.VoucherChannelListingTableName + ".ChannelID": channelIDs,
			},
		})
	if appErr != nil {
		for idx := range idPairs {
			res[idx] = &dataloader.Result[*model.VoucherChannelListing]{Error: appErr}
		}
		return res
	}

	for _, rel := range voucherChannelListings {
		voucherChannelListingMap[rel.VoucherID+"__"+rel.ChannelID] = rel
	}

	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[*model.VoucherChannelListing]{Data: voucherChannelListingMap[id]}
	}
	return res
}

func voucherChannelListingByVoucherIdLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.VoucherChannelListing] {
	var (
		res                      = make([]*dataloader.Result[[]*model.VoucherChannelListing], len(voucherIDs))
		voucherChannelListingMap = map[string][]*model.VoucherChannelListing{} // keys are voucher ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	voucherChannelListings, appErr := embedCtx.App.Srv().DiscountService().
		VoucherChannelListingsByOption(&model.VoucherChannelListingFilterOption{
			Conditions: squirrel.Eq{model.VoucherChannelListingTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		for idx := range voucherIDs {
			res[idx] = &dataloader.Result[[]*model.VoucherChannelListing]{Error: appErr}
		}
		return res
	}

	for _, rel := range voucherChannelListings {
		voucherChannelListingMap[rel.VoucherID] = append(voucherChannelListingMap[rel.VoucherID], rel)
	}

	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.VoucherChannelListing]{Data: voucherChannelListingMap[id]}
	}
	return res
}

// ------------ voucher channel listing ---------------

type VoucherChannelListing struct {
	ID            string  `json:"id"`
	DiscountValue float64 `json:"discountValue"`
	Currency      string  `json:"currency"`
	MinSpent      *Money  `json:"minSpent"`

	// Channel       *Channel `json:"channel"`
	vcl *model.VoucherChannelListing
}

func systemVoucherChannelListingToGraphqlVoucherChannelListing(vcl *model.VoucherChannelListing) *VoucherChannelListing {
	if vcl == nil {
		return nil
	}

	vcl.PopulateNonDbFields()

	return &VoucherChannelListing{
		ID:            vcl.Id,
		Currency:      vcl.Currency,
		MinSpent:      SystemMoneyToGraphqlMoney(vcl.MinSpent),
		DiscountValue: vcl.DiscountValue.InexactFloat64(),
		vcl:           vcl,
	}
}

func (v *VoucherChannelListing) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, v.vcl.ChannelID)()
	if err != nil {
		return nil, err
	}

	return SystemChannelToGraphqlChannel(channel), nil
}

// ---------------------- sale channel listing

type SaleChannelListing struct {
	ID            string  `json:"id"`
	DiscountValue float64 `json:"discountValue"`
	Currency      string  `json:"currency"`

	// Channel       *Channel `json:"channel"`
	scl *model.SaleChannelListing
}

func systemSaleChannelListingToGraphqlSaleChannelListing(scl *model.SaleChannelListing) *SaleChannelListing {
	if scl == nil {
		return nil
	}

	return &SaleChannelListing{
		ID:            scl.Id,
		DiscountValue: scl.DiscountValue.InexactFloat64(),
		Currency:      scl.Currency,
		scl:           scl,
	}
}

func (s *SaleChannelListing) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, s.scl.ChannelID)()
	if err != nil {
		return nil, err
	}

	return SystemChannelToGraphqlChannel(channel), nil
}
