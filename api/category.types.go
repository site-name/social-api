package api

import (
	"context"
	"unsafe"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

type Category struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	SeoTitle        *string         `json:"seoTitle"`
	SeoDescription  *string         `json:"seoDescription"`
	Level           int32           `json:"level"`
	Slug            string          `json:"slug"`
	Description     JSONString      `json:"description"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`

	c *model.Category

	// BackgroundImage *Image          `json:"backgroundImage"`
	// Level           int32           `json:"level"`
	// Parent          *Category       `json:"parent"`
	// Ancestors       *CategoryCountableConnection `json:"ancestors"`
	// Products        *ProductCountableConnection  `json:"products"`
	// Children        *CategoryCountableConnection `json:"children"`
	// Translation     *CategoryTranslation         `json:"translation"`
}

func systemCategoryToGraphqlCategory(c *model.Category) *Category {
	if c == nil {
		return nil
	}

	return &Category{
		Id:              c.Id,
		Name:            c.Name,
		Slug:            c.Slug,
		SeoTitle:        &c.SeoTitle,
		SeoDescription:  &c.SeoDescription,
		Level:           int32(c.Level),
		Description:     JSONString(c.Description),
		Metadata:        MetadataToSlice(c.Metadata),
		PrivateMetadata: MetadataToSlice(c.PrivateMetadata),
		c:               c,
	}
}

func (c *Category) BackgroundImage(ctx context.Context, args struct{ Size *int32 }) (*Image, error) {
	panic("not implemented")
}

func (c *Category) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*CategoryTranslation, error) {
	panic("not implemented")
}

func (c *Category) Parent(ctx context.Context) (*Category, error) {
	if c.c.ParentID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	categories, appErr := embedCtx.App.Srv().ProductService().CategoryByIds([]string{*c.c.ParentID}, true)
	if appErr != nil {
		return nil, appErr
	}

	return systemCategoryToGraphqlCategory(categories[0]), nil
}

func (c *Category) Ancestors(ctx context.Context, args GraphqlParams) (*CategoryCountableConnection, error) {
	panic("not implemented")
}

func (c *Category) Children(ctx context.Context, args GraphqlParams) (*CategoryCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	filter := func(c *model.Category) bool { return c.ParentID == &c.Id }
	children := embedCtx.App.Srv().ProductService().FilterCategoriesFromCache(filter)

	keyFunc := func(c *model.Category) string { return c.Slug }
	res, appErr := newGraphqlPaginator(children, keyFunc, systemCategoryToGraphqlCategory, args).parse("Category.Children")
	if appErr != nil {
		return nil, appErr
	}

	return (*CategoryCountableConnection)(unsafe.Pointer(res)), nil
}

func (c *Category) Products(ctx context.Context, args struct {
	Channel *string
	GraphqlParams
}) (*ProductCountableConnection, error) {
	panic("not implemented")
}
