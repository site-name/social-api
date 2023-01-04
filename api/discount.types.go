package api

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
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

func SystemOrderDiscountToGraphqlOrderDiscount(r *model.OrderDiscount) *OrderDiscount {
	if r == nil {
		return &OrderDiscount{}
	}

	if r.Value == nil {
		r.Value = &decimal.Zero
	}

	return &OrderDiscount{
		ID:             r.Id,
		Type:           OrderDiscountType(r.Type),
		ValueType:      DiscountValueTypeEnum(r.ValueType),
		Value:          PositiveDecimal(*r.Value),
		Name:           r.Name,
		TranslatedName: r.TranslatedName,
		Reason:         r.Reason,
		Amount:         SystemMoneyToGraphqlMoney(r.Amount),
	}
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
	// Categories               *CategoryCountableConnection       `json:"categories"`
	// Collections              *CollectionCountableConnection     `json:"collections"`
	// Products                 *ProductCountableConnection        `json:"products"`
	// Variants                 *ProductVariantCountableConnection `json:"variants"`
	Countries   []*CountryDisplay   `json:"countries"`
	Translation *VoucherTranslation `json:"translation"`
	// DiscountValue            *float64                           `json:"discountValue"`
	// Currency                 *string                            `json:"currency"`
	// MinSpent                 *Money                             `json:"minSpent"`
	// ChannelListings          []*VoucherChannelListing           `json:"channelListings"`
}

func systemVoucherToGraphqlVoucher(v *model.Voucher) *Voucher {
	if v == nil {
		return nil
	}
	panic("not implemented")
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
