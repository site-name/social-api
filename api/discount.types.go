package api

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

type DiscountInfo struct {
	Sale            any // either *Sale || *Voucher
	ChannelListings map[string]*SaleChannelListing
	ProductIDs      []string
	CategoryIDs     []string
	CollectionIDs   []string
	VariantsIDs     []string
}

func discountsByDateTimeLoader(ctx context.Context, dateTimes []*time.Time) []*dataloader.Result[[]*DiscountInfo] {
	var (
		res             = make([]*dataloader.Result[[]*DiscountInfo], len(dateTimes))
		appErr          *model.AppError
		salesMap        = map[*time.Time]model.Sales{}
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
		sales, appErr := discountService.ActiveSales(dateTime)
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

	// for i, datetime := range dateTimes {
	// 	items := make([]*DiscountInfo, len(salesMap[datetime]))

	// 	for idx, sale := range salesMap[datetime] {
	// 		items[idx] = &DiscountInfo{
	// 			Sale: sale,
	// 			ChannelListings: channelListings[sale.Id],
	// 		}
	// 	}
	// }
	panic("not implemented")

errorLabel:
	for i := range dateTimes {
		res[i] = &dataloader.Result[[]*DiscountInfo]{Error: err}
	}
	return res
}

func saleChannelListingBySaleIdAndChanneSlugLoader(ctx context.Context, saleIDChannelIDPairs []string) []*dataloader.Result[*model.DiscountInfo] {
	panic("not implemented")
}
