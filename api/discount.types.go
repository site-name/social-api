package api

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
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

// saleIDChannelIDPairs are strings with format of saleID__channelID.
func saleChannelListingBySaleIdAndChanneSlugLoader(ctx context.Context, saleIDChannelIDPairs []string) []*dataloader.Result[*model.DiscountInfo] {
	panic("not implemented")
}

func saleChannelListingBySaleIdLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.SaleChannelListing] {
	// 	var (
	// 		res = make([]*dataloader.Result[[]*model.SaleChannelListing], len(saleIDs))
	// 	)
	// 	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	// 	if err != nil {
	// 		goto errorLabel
	// 	}

	// errorLabel:
	//
	//	for idx := range saleIDs {
	//		res[idx] = &dataloader.Result[[]*model.SaleChannelListing]{Error: err}
	//	}
	//	return res
	panic("not implemented")
}
