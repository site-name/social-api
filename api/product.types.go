package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func SystemProductToGraphqlProduct(prd *model.Product) *Product {
	if prd == nil {
		return nil
	}

	panic("not implemented")
}

func graphqlProductsByIDsLoader(ctx context.Context, ids []string) []*dataloader.Result[*Product] {
	var (
		res      []*dataloader.Result[*Product]
		appErr   *model.AppError
		products model.Products
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	products, appErr = embedCtx.App.
		Srv().
		ProductService().
		ProductsByOption(&model.ProductFilterOption{
			Id: squirrel.Eq{store.ProductTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, prd := range products {
		res = append(res, &dataloader.Result[*Product]{Data: SystemProductToGraphqlProduct(prd)})
	}
	return res

errorLabel:
	for range ids {
		res = append(res, &dataloader.Result[*Product]{Error: err})
	}
	return res
}
