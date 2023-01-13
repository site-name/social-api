package api

import (
	"context"

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
