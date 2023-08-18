package api

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
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
	variant, err := ProductVariantByIdLoader.Load(ctx, line.variantID)()
	if err != nil {
		return nil, err
	}

	return SystemProductVariantToGraphqlProductVariant(variant), nil
}

func (line *CheckoutLine) TotalPrice(ctx context.Context) (*TaxedMoney, error) {
	checkout, err := CheckoutByTokenLoader.Load(ctx, line.checkoutID)()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	discounts, err := DiscountsByDateTimeLoader.Load(ctx, now)()
	if err != nil {
		return nil, err
	}

	checkoutInfo, err := CheckoutInfoByCheckoutTokenLoader.Load(ctx, checkout.Token)()
	if err != nil {
		return nil, err
	}

	checkoutLineInfos, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, checkout.Token)()
	if err != nil {
		return nil, err
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()

	for _, lineInfo := range checkoutLineInfos {
		if lineInfo.Line.Id == line.ID {
			address := checkoutInfo.ShippingAddress
			if address == nil {
				address = checkoutInfo.BillingAddress
			}

			taxedMoney, appErr := pluginManager.CalculateCheckoutLineTotal(
				*checkoutInfo,
				checkoutLineInfos,
				*lineInfo,
				address,
				discounts,
			)
			if appErr != nil {
				return nil, appErr
			}

			return SystemTaxedMoneyToGraphqlTaxedMoney(taxedMoney), nil
		}
	}

	return nil, nil
}

func (line *CheckoutLine) RequiresShipping(ctx context.Context) (*bool, error) {
	productType, err := ProductTypeByVariantIdLoader.Load(ctx, line.variantID)()
	if err != nil {
		return nil, err
	}

	return productType.IsShippingRequired, nil
}

func checkoutLinesByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[[]*model.CheckoutLine] {
	res := make([]*dataloader.Result[[]*model.CheckoutLine], len(tokens))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	checkoutLines, appErr := embedCtx.App.Srv().
		CheckoutService().
		CheckoutLinesByOption(&model.CheckoutLineFilterOption{
			Conditions: squirrel.Eq{model.CheckoutLineTableName + ".CheckoutID": tokens},
		})
	if appErr != nil {
		for idx := range tokens {
			res[idx] = &dataloader.Result[[]*model.CheckoutLine]{Error: appErr}
		}
		return res
	}

	// checkoutLinesMap has keys are checkout tokens.
	// values are checkout lines belong to the checkout parent
	var checkoutLinesMap = map[string][]*model.CheckoutLine{}
	for _, line := range checkoutLines {
		if line != nil {
			checkoutLinesMap[line.CheckoutID] = append(checkoutLinesMap[line.CheckoutID], line)
		}
	}

	for idx, token := range tokens {
		res[idx] = &dataloader.Result[[]*model.CheckoutLine]{Data: checkoutLinesMap[token]}
	}
	return res
}

func checkoutLineByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.CheckoutLine] {
	res := make([]*dataloader.Result[*model.CheckoutLine], len(ids))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkoutLines, appErr := embedCtx.App.Srv().
		CheckoutService().
		CheckoutLinesByOption(&model.CheckoutLineFilterOption{
			Conditions: squirrel.Eq{model.CheckoutLineTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.CheckoutLine]{Error: appErr}
		}
		return res
	}

	var checkoutLineMap = map[string]*model.CheckoutLine{}
	for _, line := range checkoutLines {
		checkoutLineMap[line.Id] = line
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.CheckoutLine]{Data: checkoutLineMap[id]}
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

	checkouts, errs := CheckoutByTokenLoader.LoadMany(ctx, tokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	checkoutLines, errs = CheckoutLinesByCheckoutTokenLoader.LoadMany(ctx, tokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for _, lines := range checkoutLines {
		for _, line := range lines {
			variantIDS = append(variantIDS, line.VariantID)
		}
	}

	variantIDS = lo.Uniq(variantIDS)

	if len(variantIDS) == 0 {
		return res
	}

	channelIDS = lo.Map(checkouts, func(c *model.Checkout, _ int) string { return c.ChannelID })

	variants, errs = ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	products, errs = ProductByVariantIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	productTypes, errs = ProductTypeByVariantIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	collections, errs = CollectionsByVariantIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for i := 0; i < min(len(channelIDS), len(checkoutLines)); i++ {
		for _, line := range checkoutLines[i] {
			variantIDChannelIDPairs = append(variantIDChannelIDPairs, line.VariantID+"__"+channelIDS[i])
		}
	}

	channelListings, errs = VariantChannelListingByVariantIdAndChannelLoader.LoadMany(ctx, variantIDChannelIDPairs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for i := 0; i < min(len(variantIDS), len(variants)); i++ {
		variantsMap[variantIDS[i]] = variants[i]
	}

	for i := 0; i < min(len(variantIDS), len(products)); i++ {
		productsMap[variantIDS[i]] = products[i]
	}

	for i := 0; i < min(len(variantIDS), len(productTypes)); i++ {
		productTypesMap[variantIDS[i]] = productTypes[i]
	}

	for i := 0; i < min(len(variantIDS), len(collections)); i++ {
		collectionsMap[variantIDS[i]] = collections[i]
	}

	for i := 0; i < min(len(variantIDChannelIDPairs), len(channelListings)); i++ {
		channelListingsMap[variantIDChannelIDPairs[i]] = channelListings[i]
	}

	for i := 0; i < min(len(checkouts), len(checkoutLines)); i++ {
		for _, line := range checkoutLines[i] {
			linesInfoMap[checkouts[i].Token] = append(linesInfoMap[checkouts[i].Token], &model.CheckoutLineInfo{
				Line:           *line,
				Variant:        *variantsMap[line.VariantID],
				ChannelListing: *channelListingsMap[line.VariantID+"__"+checkouts[i].ChannelID],
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
