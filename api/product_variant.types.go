package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type ProductVariant struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Sku             *string         `json:"sku"`
	TrackInventory  bool            `json:"trackInventory"`
	Weight          *Weight         `json:"weight"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	Channel         *string         `json:"channel"`
	Margin          *int32          `json:"margin"`
	QuantityOrdered *int32          `json:"quantityOrdered"`

	// Translation     *ProductVariantTranslation `json:"translation"`
	// DigitalContent  *DigitalContent            `json:"digitalContent"`
	// Stocks            []*Stock                        `json:"stocks"`
	// QuantityAvailable int32                           `json:"quantityAvailable"`
	// Preorder          *PreorderData                   `json:"preorder"`
	// ChannelListings   []*ProductVariantChannelListing `json:"channelListings"`
	// Pricing           *VariantPricingInfo             `json:"pricing"`
	// Attributes        []*SelectedAttribute            `json:"attributes"`
	// Product           *Product                        `json:"product"`
	// Revenue           *TaxedMoney                     `json:"revenue"`
	// Media             []*ProductMedia                 `json:"media"`
}

func SystemProductVariantToGraphqlProductVariant(variant *model.ProductVariant) *ProductVariant {
	if variant == nil {
		return nil
	}

	res := &ProductVariant{
		ID:              variant.Id,
		Name:            variant.Name,
		Sku:             variant.Sku,
		TrackInventory:  *variant.TrackInventory,
		Channel:         model.NewString("unknown"), // ??
		Metadata:        MetadataToSlice(variant.Metadata),
		PrivateMetadata: MetadataToSlice(variant.PrivateMetadata),
		Margin:          model.NewInt32(0), // ??
		QuantityOrdered: model.NewInt32(0), // ??
	}
	if variant.Weight != nil {
		res.Weight = &Weight{WeightUnitsEnum(variant.WeightUnit), float64(*variant.Weight)}
	}

	return res
}

func productVariantByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*ProductVariant] {
	if len(ids) == 0 {
		return []*dataloader.Result[*ProductVariant]{}
	}

	var (
		productVariants model.ProductVariants
		appErr          *model.AppError
		res             []*dataloader.Result[*ProductVariant]
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productVariants, appErr = embedCtx.
		App.
		Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			Id: squirrel.Eq{store.ProductVariantTableName + ".Id": ids},
		})

	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, variant := range productVariants {
		res = append(res, &dataloader.Result[*ProductVariant]{Data: SystemProductVariantToGraphqlProductVariant(variant)})
	}
	return res

errorLabel:
	for range ids {
		res = append(res, &dataloader.Result[*ProductVariant]{Error: err})
	}
	return res
}

func graphqlProductVariantsByProductIDLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*ProductVariant] {
	var (
		productVariants model.ProductVariants
		appErr          *model.AppError
		res             []*dataloader.Result[[]*ProductVariant]

		// keys are product ids
		variantsMap = map[string][]*ProductVariant{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productVariants, appErr = embedCtx.
		App.
		Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			ProductID: squirrel.Eq{store.ProductVariantTableName + ".ProductID": productIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, variant := range productVariants {
		if variant != nil {
			variantsMap[variant.ProductID] = append(variantsMap[variant.ProductID], SystemProductVariantToGraphqlProductVariant(variant))
		}
	}

	for _, productID := range productIDs {
		res = append(res, &dataloader.Result[[]*ProductVariant]{Data: variantsMap[productID]})
	}
	return res

errorLabel:
	for range productIDs {
		res = append(res, &dataloader.Result[[]*ProductVariant]{Error: err})
	}
	return res
}