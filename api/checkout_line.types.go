package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type CheckoutLine struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`

	variantID  string
	checkoutID string

	// RequiresShipping *bool           `json:"requiresShipping"`
	// TotalPrice       *TaxedMoney     `json:"totalPrice"`
	// Variant          *ProductVariant `json:"variant"`
}

func SystemCheckoutLineToGraphqlCheckoutLine(line *model.CheckoutLine) *CheckoutLine {
	if line == nil {
		return nil
	}

	res := &CheckoutLine{
		ID:         line.Id,
		Quantity:   int32(line.Quantity),
		variantID:  line.VariantID,
		checkoutID: line.CheckoutID,
	}
	return res
}

func (line *CheckoutLine) Variant(ctx context.Context) (*ProductVariant, error) {
	variant, err := dataloaders.ProductVariantByIdLoader.Load(ctx, line.variantID)()
	if err != nil {
		return nil, err
	}

	return SystemProductVariantToGraphqlProductVariant(variant), nil
}

func (line *CheckoutLine) TotalPrice(ctx context.Context) (*TaxedMoney, error) {
	// checkout, err := dataloaders.CheckoutByTokenLoader.Load(ctx, line.checkoutID)()
	// if err != nil {
	// 	return nil, err
	// }
	panic("not implemented")
}

func (line *CheckoutLine) RequiresShipping(ctx context.Context) (*bool, error) {
	productType, err := dataloaders.ProductTypeByVariantIdLoader.Load(ctx, line.variantID)()
	if err != nil {
		return nil, err
	}

	return productType.IsShippingRequired, nil
}

func checkoutLinesByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[[]*model.CheckoutLine] {
	var (
		res           = make([]*dataloader.Result[[]*model.CheckoutLine], len(tokens))
		appErr        *model.AppError
		checkoutLines model.CheckoutLines

		// checkoutLinesMap has keys are checkout tokens.
		// values are checkout lines belong to the checkout parent
		checkoutLinesMap = map[string][]*model.CheckoutLine{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	checkoutLines, appErr = embedCtx.App.Srv().
		CheckoutService().
		CheckoutLinesByOption(&model.CheckoutLineFilterOption{
			CheckoutID: squirrel.Eq{store.CheckoutLineTableName + ".CheckoutID": tokens},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, line := range checkoutLines {
		if line != nil {
			checkoutLinesMap[line.CheckoutID] = append(checkoutLinesMap[line.CheckoutID], line)
		}
	}

	for idx, token := range tokens {
		res[idx] = &dataloader.Result[[]*model.CheckoutLine]{Data: checkoutLinesMap[token]}
	}
	return res

errorLabel:
	for idx := range tokens {
		res[idx] = &dataloader.Result[[]*model.CheckoutLine]{Error: err}
	}
	return res
}

func checkoutLineByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*CheckoutLine] {
	var (
		res             = make([]*dataloader.Result[*CheckoutLine], len(ids))
		checkoutLines   model.CheckoutLines
		appErr          *model.AppError
		checkoutLineMap = map[string]*CheckoutLine{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	checkoutLines, appErr = embedCtx.App.Srv().
		CheckoutService().
		CheckoutLinesByOption(&model.CheckoutLineFilterOption{
			Id: squirrel.Eq{store.CheckoutLineTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, line := range checkoutLines {
		checkoutLineMap[line.Id] = SystemCheckoutLineToGraphqlCheckoutLine(line)
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*CheckoutLine]{Data: checkoutLineMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*CheckoutLine]{Error: err}
	}
	return res
}

func checkoutLinesInfoByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[[]*model.CheckoutLineInfo] {
	var (
		res = make([]*dataloader.Result[[]*model.CheckoutLineInfo], len(tokens))

		variantIDS      []string
		channelIDS      []string
		checkoutLines   [][]*model.CheckoutLine
		variants        []*model.ProductVariant
		products        []*model.Product
		productTypes    []*model.ProductType
		collections     [][]*model.Collection
		channelListings []*model.ProductVariantChannelListing

		variantIDChannelIDPairs []string // slice of variantID__channelID values

		variantsMap        = map[string]*model.ProductVariant{}               // keys are product variant ids
		productsMap        = map[string]*model.Product{}                      // keys are product variant ids
		productTypesMap    = map[string]*model.ProductType{}                  // keys are product variant ids
		collectionsMap     = map[string][]*model.Collection{}                 // keys are product variant ids
		channelListingsMap = map[string]*model.ProductVariantChannelListing{} // keys are variantID__channelID format

		linesInfoMap = map[string][]*model.CheckoutLineInfo{} // keys are checkout tokens
	)

	checkouts, errs := dataloaders.CheckoutByTokenLoader.LoadMany(ctx, tokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	checkoutLines, errs = dataloaders.CheckoutLinesByCheckoutTokenLoader.LoadMany(ctx, tokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for _, lines := range checkoutLines {
		for _, line := range lines {
			variantIDS = append(variantIDS, line.VariantID)
		}
	}

	if len(variantIDS) == 0 {
		return res
	}

	channelIDS = lo.Map(checkouts, func(c *Checkout, _ int) string { return c.channelID })

	variants, errs = dataloaders.ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	products, errs = dataloaders.ProductByVariantIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	productTypes, errs = dataloaders.ProductTypeByVariantIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	collections, errs = dataloaders.CollectionsByVariantIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for i := 0; i < util.Min(len(channelIDS), len(checkoutLines)); i++ {
		for _, line := range checkoutLines[i] {
			variantIDChannelIDPairs = append(variantIDChannelIDPairs, line.VariantID+"__"+channelIDS[i])
		}
	}

	channelListings, errs = dataloaders.VariantChannelListingByVariantIdAndChannelIdLoader.LoadMany(ctx, variantIDChannelIDPairs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for i := 0; i < util.Min(len(variantIDS), len(variants)); i++ {
		variantsMap[variantIDS[i]] = variants[i]
	}

	for i := 0; i < util.Min(len(variantIDS), len(products)); i++ {
		productsMap[variantIDS[i]] = products[i]
	}

	for i := 0; i < util.Min(len(variantIDS), len(productTypes)); i++ {
		productTypesMap[variantIDS[i]] = productTypes[i]
	}

	for i := 0; i < util.Min(len(variantIDS), len(collections)); i++ {
		collectionsMap[variantIDS[i]] = collections[i]
	}

	for i := 0; i < util.Min(len(variantIDChannelIDPairs), len(channelListings)); i++ {
		channelListingsMap[variantIDChannelIDPairs[i]] = channelListings[i]
	}

	for i := 0; i < util.Min(len(checkouts), len(checkoutLines)); i++ {
		for _, line := range checkoutLines[i] {
			linesInfoMap[checkouts[i].Token] = append(linesInfoMap[checkouts[i].Token], &model.CheckoutLineInfo{
				Line:           *line,
				Variant:        *variantsMap[line.VariantID],
				ChannelListing: *channelListingsMap[line.VariantID+"__"+checkouts[i].channelID],
				Product:        *productsMap[line.VariantID],
				ProductType:    *productTypesMap[line.VariantID],
				Collections:    collectionsMap[line.VariantID],
			})
		}
	}

	for idx, token := range tokens {
		res[idx] = &dataloader.Result[[]*model.CheckoutLineInfo]{Data: linesInfoMap[token]}
	}
	return res

errorLabel:
	for idx := range tokens {
		res[idx] = &dataloader.Result[[]*model.CheckoutLineInfo]{Error: errs[0]}
	}
	return res
}
