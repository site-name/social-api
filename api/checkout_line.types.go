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
	return dataloaders.ProductVariantByIdLoader.Load(ctx, line.variantID)()
}

func (line *CheckoutLine) TotalPrice(ctx context.Context) (*TaxedMoney, error) {
	// checkout, err := dataloaders.CheckoutByTokenLoader.Load(ctx, line.checkoutID)()
	// if err != nil {
	// 	return nil, err
	// }
	panic("not implemented")
}

func (line *CheckoutLine) RequiresShipping(ctx context.Context) (*bool, error) {
	panic("not implemented")
}

func checkoutLinesByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[[]*CheckoutLine] {
	var (
		res           []*dataloader.Result[[]*CheckoutLine]
		appErr        *model.AppError
		checkoutLines model.CheckoutLines

		// checkoutLinesMap has keys are checkout tokens.
		// values are checkout lines belong to the checkout parent
		checkoutLinesMap = map[string][]*CheckoutLine{}
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
			checkoutLinesMap[line.CheckoutID] = append(
				checkoutLinesMap[line.CheckoutID],
				SystemCheckoutLineToGraphqlCheckoutLine(line))
		}
	}

	for _, token := range tokens {
		res = append(res, &dataloader.Result[[]*CheckoutLine]{Data: checkoutLinesMap[token]})
	}
	return res

errorLabel:
	for range tokens {
		res = append(res, &dataloader.Result[[]*CheckoutLine]{Error: err})
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

type CheckoutLineInfo struct {
	Line           CheckoutLine
	Variant        ProductVariant
	ChannelListing ProductVariantChannelListing
	Product        Product
	ProductType    ProductType
	Collections    []*Collection
}

func CheckoutLinesInfoByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[[]*CheckoutLineInfo] {
	var (
		res []*dataloader.Result[[]*CheckoutLineInfo]

		variantIDS      []string
		channelIDS      []string
		checkoutLines   [][]*CheckoutLine
		variants        []*ProductVariant
		products        []*Product
		productTypes    []*ProductType
		collections     [][]*Collection
		channelListings []*ProductVariantChannelListing

		variantIDChannelIDPairs []string // slice of variantID__channelID pairs

		variantsMap        = map[string]*ProductVariant{}               // keys are product variant ids
		productsMap        = map[string]*Product{}                      // keys are product variant ids
		productTypesMap    = map[string]*ProductType{}                  // keys are product variant ids
		collectionsMap     = map[string][]*Collection{}                 // keys are product variant ids
		channelListingsMap = map[string]*ProductVariantChannelListing{} // keys are variantID__channelID format

		linesInfoMap = map[string][]*CheckoutLineInfo{} // keys are checkout tokens
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
			if line != nil {
				variantIDS = append(variantIDS, line.variantID)
			}
		}
	}

	if len(variantIDS) == 0 {
		return make([]*dataloader.Result[[]*CheckoutLineInfo], len(tokens))
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
		lines := checkoutLines[i]
		for _, line := range lines {
			variantIDChannelIDPairs = append(variantIDChannelIDPairs, line.variantID+"__"+channelIDS[i])
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
		var (
			checkout = checkouts[i]
			lines    = checkoutLines[i]
		)

		for _, line := range lines {
			linesInfoMap[checkout.Token] = append(linesInfoMap[checkout.Token], &CheckoutLineInfo{
				Line:           *line,
				Variant:        *variantsMap[line.variantID],
				ChannelListing: *channelListingsMap[line.variantID+"__"+checkout.channelID],
				Product:        *productsMap[line.variantID],
				ProductType:    *productTypesMap[line.variantID],
				Collections:    collectionsMap[line.variantID],
			})
		}
	}

	for _, token := range tokens {
		res = append(res, &dataloader.Result[[]*CheckoutLineInfo]{Data: linesInfoMap[token]})
	}
	return res

errorLabel:
	for range tokens {
		res = append(res, &dataloader.Result[[]*CheckoutLineInfo]{Error: errs[0]})
	}
	return res
}
