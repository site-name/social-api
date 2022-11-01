package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func systemProductVariantToGraphqlProductVariant(variant *model.ProductVariant) *ProductVariant {
	if variant == nil {
		return nil
	}

	res := &ProductVariant{}
	panic("not implemented")
	return res
}

func graphqlProductVariantsByIDsLoader(ctx context.Context, ids []string) []*dataloader.Result[*ProductVariant] {
	var (
		productVariants model.ProductVariants
		appErr          *model.AppError
		res             []*dataloader.Result[*ProductVariant]
	)

	if len(ids) == 0 {
		return res
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productVariants, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			Id: squirrel.Eq{store.ProductVariantTableName + ".Id": ids},
		})

	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, variant := range productVariants {
		res = append(res, &dataloader.Result[*ProductVariant]{Data: systemProductVariantToGraphqlProductVariant(variant)})
	}
	return res

errorLabel:
	for range ids {
		res = append(res, &dataloader.Result[*ProductVariant]{Error: err})
	}
	return res
}
