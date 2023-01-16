package api

import (
	"context"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Category struct {
	ID              string                       `json:"id"`
	SeoTitle        *string                      `json:"seoTitle"`
	SeoDescription  *string                      `json:"seoDescription"`
	Name            string                       `json:"name"`
	Description     JSONString                   `json:"description"`
	Slug            string                       `json:"slug"`
	Parent          *Category                    `json:"parent"`
	Level           int32                        `json:"level"`
	PrivateMetadata []*MetadataItem              `json:"privateMetadata"`
	Metadata        []*MetadataItem              `json:"metadata"`
	Ancestors       *CategoryCountableConnection `json:"ancestors"`
	Products        *ProductCountableConnection  `json:"products"`
	Children        *CategoryCountableConnection `json:"children"`
	BackgroundImage *Image                       `json:"backgroundImage"`
	Translation     *CategoryTranslation         `json:"translation"`
}

func systemCategoryToGraphqlCategory(c *model.Category) *Category {
	if c == nil {
		return nil
	}

	panic("not implemented")
}

func categoryByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Category] {
	var (
		res         = make([]*dataloader.Result[*model.Category], len(ids))
		categories  model.Categories
		appErr      *model.AppError
		categoryMap = map[string]*model.Category{} // keys are category ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	categories, appErr = embedCtx.App.Srv().ProductService().CategoriesByOption(&model.CategoryFilterOption{
		Id: squirrel.Eq{store.CategoryTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	categoryMap = lo.SliceToMap(categories, func(c *model.Category) (string, *model.Category) { return c.Id, c })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Category]{Data: categoryMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Category]{Error: err}
	}
	return res
}

func categoriesByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Category] {
	var (
		res                = make([]*dataloader.Result[[]*model.Category], len(voucherIDs))
		categories         model.Categories
		appErr             *model.AppError
		voucherCategories  []*model.VoucherCategory
		voucherCategoryMap = map[string]string{}           // values are voucher ids, keys are category ids
		categoryMap        = map[string]model.Categories{} // keys are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	categories, appErr = embedCtx.
		App.
		Srv().
		ProductService().
		CategoriesByOption(&model.CategoryFilterOption{
			VoucherID: squirrel.Eq{store.VoucherCategoryTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	voucherCategories, err = embedCtx.
		App.
		Srv().
		Store.
		VoucherCategory().
		FilterByOptions(&model.VoucherCategoryFilterOption{
			VoucherID: squirrel.Eq{store.VoucherCategoryTableName + ".VoucherID": voucherIDs},
		})
	if err != nil {
		err = model.NewAppError("categoriesByVoucherIDLoader", "app.discount.voucher_categories_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range voucherCategories {
		voucherCategoryMap[rel.CategoryID] = rel.VoucherID
	}

	for _, cate := range categories {
		voucherID, ok := voucherCategoryMap[cate.Id]
		if ok {
			categoryMap[voucherID] = append(categoryMap[voucherID], cate)
		}
	}

	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Data: categoryMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Error: err}
	}
	return res
}

func collectionsByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Collection] {
	var (
		collections          model.Collections
		appErr               *model.AppError
		res                  = make([]*dataloader.Result[[]*model.Collection], len(voucherIDs))
		voucherCollections   []*model.VoucherCollection
		collectionMap        = map[string]model.Collections{} // keys are voucher ids
		voucherCollectionMap = map[string]string{}            // keys are collection ids, values are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collections, appErr = embedCtx.App.Srv().
		ProductService().
		CollectionsByOption(&model.CollectionFilterOption{
			VoucherID: squirrel.Eq{store.VoucherCollectionTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	voucherCollections, err = embedCtx.App.Srv().Store.VoucherCollection().FilterByOptions(&model.VoucherCollectionFilterOptions{
		VoucherID: squirrel.Eq{store.VoucherCollectionTableName + ".VoucherID": voucherIDs},
	})
	if err != nil {
		err = model.NewAppError("collectionsByVoucherIDLoader", "app.discount.voucher_collections_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range voucherCollections {
		voucherCollectionMap[rel.CollectionID] = rel.VoucherID
	}

	for _, col := range collections {
		voucherID, ok := voucherCollectionMap[col.Id]
		if ok {
			collectionMap[voucherID] = append(collectionMap[voucherID], col)
		}
	}

	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Data: collectionMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Error: err}
	}
	return res
}
