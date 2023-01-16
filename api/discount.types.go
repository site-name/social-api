package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func discountsByDateTimeLoader(ctx context.Context, dateTimes []time.Time) []*dataloader.Result[[]*model.DiscountInfo] {
	var (
		res             = make([]*dataloader.Result[[]*model.DiscountInfo], len(dateTimes))
		appErr          *model.AppError
		salesMap        = map[time.Time]model.Sales{}
		saleIDS         []string
		collections     = map[string][]string{}
		channelListings = map[string]map[string]*model.SaleChannelListing{}
		products        = map[string][]string{}
		categories      = map[string][]string{}
		variants        = map[string][]string{}

		discountService sub_app_iface.DiscountService
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	discountService = embedCtx.App.Srv().DiscountService()

	for _, dateTime := range dateTimes {
		sales, appErr := discountService.ActiveSales(&dateTime)
		if appErr != nil {
			err = appErr
			goto errorLabel
		}

		for _, sale := range sales {
			saleIDS = append(saleIDS, sale.Id)
		}

		salesMap[dateTime] = sales
	}

	collections, appErr = discountService.FetchCollections(saleIDS)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	channelListings, appErr = discountService.FetchSaleChannelListings(saleIDS)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	products, appErr = discountService.FetchProducts(saleIDS)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	categories, appErr = discountService.FetchCategories(saleIDS)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	variants, appErr = discountService.FetchVariants(saleIDS)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for i, datetime := range dateTimes {
		items := make([]*model.DiscountInfo, len(salesMap[datetime]))

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

		res[i] = &dataloader.Result[[]*model.DiscountInfo]{Data: items}
	}

	return res

errorLabel:
	for i := range dateTimes {
		res[i] = &dataloader.Result[[]*model.DiscountInfo]{Error: err}
	}
	return res
}

// NOTE: saleIDChannelIDPairs are strings with format of saleID__channelID.
func saleChannelListingBySaleIdAndChanneSlugLoader(ctx context.Context, saleIDChannelIDPairs []string) []*dataloader.Result[*model.SaleChannelListing] {
	var (
		res        = make([]*dataloader.Result[*model.SaleChannelListing], len(saleIDChannelIDPairs))
		saleIDs    []string
		channelIDs []string
		listings   []*model.SaleChannelListing
		appErr     *model.AppError
		listingMap = map[string]*model.SaleChannelListing{} // keys are string format of saleID__channelID
	)

	for _, item := range saleIDChannelIDPairs {
		index := strings.Index(item, "__")
		if index < 0 {
			continue
		}

		saleIDs = append(saleIDs, item[:index])
		channelIDs = append(channelIDs, item[index+2:])
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	listings, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleChannelListingsByOptions(&model.SaleChannelListingFilterOption{
			SaleID:    squirrel.Eq{store.SaleChannelListingTableName + ".SaleID": saleIDs},
			ChannelID: squirrel.Eq{store.SaleChannelListingTableName + ".ChannelID": channelIDs},
			// SelectRelatedChannel: true,
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	listingMap = lo.SliceToMap(listings, func(l *model.SaleChannelListing) (string, *model.SaleChannelListing) {
		return l.SaleID + "__" + l.ChannelID, l
	})

	for idx, pair := range saleIDChannelIDPairs {
		res[idx] = &dataloader.Result[*model.SaleChannelListing]{Data: listingMap[pair]}
	}
	return res

errorLabel:
	for idx := range saleIDChannelIDPairs {
		res[idx] = &dataloader.Result[*model.SaleChannelListing]{Error: err}
	}
	return res
}

func saleChannelListingBySaleIdLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.SaleChannelListing] {
	var (
		res        = make([]*dataloader.Result[[]*model.SaleChannelListing], len(saleIDs))
		listings   []*model.SaleChannelListing
		appErr     *model.AppError
		listingMap = map[string][]*model.SaleChannelListing{} // keys are sale ids
	)
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	listings, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleChannelListingsByOptions(&model.SaleChannelListingFilterOption{
			SaleID: squirrel.Eq{store.SaleChannelListingTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range listings {
		listingMap[listing.SaleID] = append(listingMap[listing.SaleID], listing)
	}

	for idx, saleID := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.SaleChannelListing]{Data: listingMap[saleID]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.SaleChannelListing]{Error: err}
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
	reason         *string               `json:"reason"`
	Amount         *Money                `json:"amount"`
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
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()
	if embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
		return o.reason, nil
	}

	return nil, model.NewAppError("OrderDiscount.Reason", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}

func orderDiscountsByOrderIDLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.OrderDiscount] {
	var (
		res              = make([]*dataloader.Result[[]*model.OrderDiscount], len(orderIDs))
		orderDiscounts   model.OrderDiscounts
		appErr           *model.AppError
		orderDiscountMap = map[string]model.OrderDiscounts{} // keys are order ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	orderDiscounts, appErr = embedCtx.App.Srv().DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		OrderID: squirrel.Eq{store.OrderDiscountTableName + ".OrderID": orderIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
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

errorLabel:
	for idx := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.OrderDiscount]{Error: err}
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
		StartDate:                DateTime{util.TimeFromMillis(v.StartDate)},
		ApplyOncePerOrder:        v.ApplyOncePerOrder,
		ApplyOncePerCustomer:     v.ApplyOncePerCustomer,
		DiscountValueType:        DiscountValueTypeEnum(v.DiscountValueType),
		MinCheckoutItemsQuantity: model.NewPrimitive[int32](int32(v.MinCheckoutItemsQuantity)),
		Metadata:                 MetadataToSlice(v.Metadata),
		PrivateMetadata:          MetadataToSlice(v.PrivateMetadata),
	}

	countries := strings.Fields(v.Countries)
	if len(countries) > 0 {
		for _, code := range countries {
			res.Countries = append(res.Countries, &CountryDisplay{
				Code:    code,
				Country: model.Countries[code],
			})
		}
	}

	if v.EndDate != nil {
		res.EndDate = &DateTime{util.TimeFromMillis(*v.EndDate)}
	}
	if v.UsageLimit != nil {
		res.UsageLimit = model.NewPrimitive[int32](int32(*v.UsageLimit))
	}

	return res
}

func (v *Voucher) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*VoucherTranslation, error) {
	panic("not implemented")
}

// categories are order by names
func (v *Voucher) Categories(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CategoryCountableConnection, error) {
	categories, err := CategoriesByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	// unbase64:
	var (
		before *string
		after  *string
	)

	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Voucher.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}
		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Voucher.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}
		after = model.NewPrimitive(string(data))
	}

	g := graphqlPaginator[*model.Category, string]{
		data:    categories,
		keyFunc: func(c *model.Category) string { return c.Name },

		before: before,
		after:  after,
		first:  args.First,
		last:   args.Last,
	}

	data, hasPrev, hasNext, appErr := g.parse("Voucher.Categories")
	if appErr != nil {
		return nil, appErr
	}

	res := &CategoryCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(categories))),
		Edges: lo.Map(data, func(c *model.Category, _ int) *CategoryCountableEdge {
			return &CategoryCountableEdge{
				systemCategoryToGraphqlCategory(c),
				base64.StdEncoding.EncodeToString([]byte(c.Name)),
			}
		}),
	}
	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(res.Edges)-1].Cursor,
	}

	return res, nil
}

// collections order by slugs
func (v *Voucher) Collections(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CollectionCountableConnection, error) {
	collections, err := CollectionsByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	// unbase64
	var (
		before *string
		after  *string
	)

	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Voucher.Collections", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}
		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Voucher.Collections", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}
		after = model.NewPrimitive(string(data))
	}

	g := graphqlPaginator[*model.Collection, string]{
		data:    collections,
		keyFunc: func(c *model.Collection) string { return c.Slug },
		before:  before,
		after:   after,
		first:   args.First,
		last:    args.Last,
	}

	data, hasPrev, hasNext, appErr := g.parse("Voucher.Collections")
	if err != nil {
		return nil, appErr
	}

	res := &CollectionCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(collections))),
		Edges: lo.Map(data, func(c *model.Collection, _ int) *CollectionCountableEdge {
			return &CollectionCountableEdge{
				systemCollectionToGraphqlCollection(c),
				base64.StdEncoding.EncodeToString([]byte(c.Slug)),
			}
		}),
	}

	res.PageInfo = &PageInfo{hasNext, hasPrev, &res.Edges[0].Cursor, &res.Edges[len(res.Edges)-1].Cursor}

	return res, nil
}

func (v *Voucher) Products(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductCountableConnection, error) {
	products, err := ProductsByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	// unbase64
	var (
		before *string
		after  *string
	)

	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Voucher.Products", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}

		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Voucher.Products", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}

		after = model.NewPrimitive(string(data))
	}

	g := graphqlPaginator[*model.Product, string]{
		data:    products,
		keyFunc: func(p *model.Product) string { return p.Slug },
		before:  before,
		after:   after,
		first:   args.First,
		last:    args.Last,
	}

	data, hasPrev, hasNext, appErr := g.parse("voucher.Products")
	if appErr != nil {
		return nil, appErr
	}

	res := &ProductCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(products))),
		Edges: lo.Map(data, func(p *model.Product, _ int) *ProductCountableEdge {
			return &ProductCountableEdge{
				Node:   SystemProductToGraphqlProduct(p),
				Cursor: base64.StdEncoding.EncodeToString([]byte(p.Slug)),
			}
		}),
	}
	res.PageInfo = &PageInfo{
		hasNext,
		hasPrev,
		&res.Edges[0].Cursor,
		&res.Edges[len(res.Edges)-1].Cursor,
	}

	return res, nil
}

func (v *Voucher) Variants(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductVariantCountableConnection, error) {
	variants, err := ProductVariantsByVoucherIDLoader.Load(ctx, v.ID)()
	if err != nil {
		return nil, err
	}

	// unbase 64
	var before *string
	var after *string

	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Voucher.Variants", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}

		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Voucher.Variants", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}

		after = model.NewPrimitive(string(data))
	}

	p := graphqlPaginator[*model.ProductVariant, string]{
		data:    variants,
		keyFunc: func(pv *model.ProductVariant) string { return pv.Sku },
		before:  before,
		after:   after,
		first:   args.First,
		last:    args.Last,
	}

	data, hasPrev, hasNext, appErr := p.parse("Voucher.Variants")
	if appErr != nil {
		return nil, appErr
	}

	res := &ProductVariantCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(variants))),
		Edges: lo.Map(data, func(v *model.ProductVariant, _ int) *ProductVariantCountableEdge {
			return &ProductVariantCountableEdge{
				Node:   SystemProductVariantToGraphqlProductVariant(v),
				Cursor: base64.StdEncoding.EncodeToString([]byte(v.Sku)),
			}
		}),
	}

	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(res.Edges)-1].Cursor,
	}

	return res, nil
}

func (v *Voucher) DiscountValue(ctx context.Context) (*float64, error) {
	// VoucherChannelListingByVoucherIdLoader.Load(ctx, v.ID)()
	panic("not implemented")
}

func (v *Voucher) Currency(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (v *Voucher) MinSpent(ctx context.Context) (*Money, error) {
	panic("not implemented")
}

func (v *Voucher) ChannelListings(ctx context.Context) ([]*VoucherChannelListing, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()
	if embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageDiscounts) {
		listings, err := VoucherChannelListingByVoucherIdLoader.Load(ctx, v.ID)()
		if err != nil {
			return nil, err
		}

		return DataloaderResultMap(listings, systemVoucherChannelListingToGraphqlVoucherChannelListing), nil
	}

	return nil, model.NewAppError("Voucher.ChannelListings", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}

func voucherByIDLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Voucher] {
	var (
		res        = make([]*dataloader.Result[*model.Voucher], len(ids))
		voucherMap = map[string]*model.Voucher{}
		appErr     *model.AppError
		vouchers   []*model.Voucher
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	vouchers, appErr = embedCtx.App.Srv().DiscountService().VouchersByOption(&model.VoucherFilterOption{
		Id: squirrel.Eq{store.VoucherTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, v := range vouchers {
		voucherMap[v.Id] = v
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Voucher]{Data: voucherMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Voucher]{Error: err}
	}
	return res
}

// NOTE: idPairs contains strings with format of voucherID__channelID
func voucherChannelListingByVoucherIdAndChanneSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.VoucherChannelListing] {
	var (
		res                      = make([]*dataloader.Result[*model.VoucherChannelListing], len(idPairs))
		voucherChannelListings   []*model.VoucherChannelListing
		appErr                   *model.AppError
		voucherChannelListingMap = map[string]*model.VoucherChannelListing{} // keys are voucher channel listing ids

		voucherIDs []string
		channelIDs []string
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	for _, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index == -1 {
			continue
		}

		voucherIDs = append(voucherIDs, pair[:index])
		channelIDs = append(channelIDs, pair[index+2:])
	}

	voucherChannelListings, appErr = embedCtx.App.Srv().DiscountService().
		VoucherChannelListingsByOption(&model.VoucherChannelListingFilterOption{
			VoucherID: squirrel.Eq{store.VoucherChannelListingTableName + ".VoucherID": voucherIDs},
			ChannelID: squirrel.Eq{store.VoucherChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range voucherChannelListings {
		voucherChannelListingMap[rel.VoucherID+"__"+rel.ChannelID] = rel
	}

	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[*model.VoucherChannelListing]{Data: voucherChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[*model.VoucherChannelListing]{Error: err}
	}
	return res
}

func voucherChannelListingByVoucherIdLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.VoucherChannelListing] {
	var (
		res                      = make([]*dataloader.Result[[]*model.VoucherChannelListing], len(voucherIDs))
		voucherChannelListings   []*model.VoucherChannelListing
		appErr                   *model.AppError
		voucherChannelListingMap = map[string][]*model.VoucherChannelListing{} // keys are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	voucherChannelListings, appErr = embedCtx.App.Srv().DiscountService().
		VoucherChannelListingsByOption(&model.VoucherChannelListingFilterOption{
			VoucherID: squirrel.Eq{store.VoucherChannelListingTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range voucherChannelListings {
		voucherChannelListingMap[rel.VoucherID] = append(voucherChannelListingMap[rel.VoucherID], rel)
	}

	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.VoucherChannelListing]{Data: voucherChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.VoucherChannelListing]{Error: err}
	}
	return res
}

// ------------ voucher channel listing

type VoucherChannelListing struct {
	ID string `json:"id"`
	// Channel       *Channel `json:"channel"`
	DiscountValue float64 `json:"discountValue"`
	Currency      string  `json:"currency"`
	MinSpent      *Money  `json:"minSpent"`

	channelID string
}

func systemVoucherChannelListingToGraphqlVoucherChannelListing(l *model.VoucherChannelListing) *VoucherChannelListing {
	if l == nil {
		return nil
	}

	l.PopulateNonDbFields()

	res := &VoucherChannelListing{
		ID:        l.Id,
		Currency:  l.Currency,
		MinSpent:  SystemMoneyToGraphqlMoney(l.MinSpent),
		channelID: l.ChannelID,
	}

	flt, _ := l.DiscountValue.Float64()
	res.DiscountValue = flt
	return res
}

func (v *VoucherChannelListing) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, v.channelID)()
	if err != nil {
		return nil, err
	}

	return SystemChannelToGraphqlChannel(channel), nil
}
