package api

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
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
		Id:              c.ID,
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
	if c.c.ParentID.IsNil() {
		return nil, nil
	}

	category, err := CategoryByIdLoader.Load(ctx, *c.c.ParentID.String)()
	if err != nil {
		return nil, err
	}

	return systemCategoryToGraphqlCategory(category), nil
}

func (c *Category) Ancestors(ctx context.Context, args GraphqlParams) (*CategoryCountableConnection, error) {
	panic("not implemented")
}

func (c *Category) Children(ctx context.Context, args GraphqlParams) (*CategoryCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	filter := func(c *model.Category) bool { return model_types.PrimitiveIsNotNilAndEqual(c.ParentID.String, c.ID) }
	children := embedCtx.App.Srv().ProductService().FilterCategoriesFromCache(filter)

	keyFunc := func(c *model.Category) []any { return []any{model.CategoryTableName + ".Slug", c.Slug} }
	res, appErr := newGraphqlPaginator(children, keyFunc, systemCategoryToGraphqlCategory, args).parse("Category.Children")
	if appErr != nil {
		return nil, appErr
	}

	return (*CategoryCountableConnection)(unsafe.Pointer(res)), nil
}

func (c *Category) Products(ctx context.Context, args struct {
	Channel *string // NOTE: Channel can be channel id or slug
	GraphqlParams
}) (*ProductCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("Category.Products", model_helper.ShopStaffRoleId, model_helper.ShopAdminRoleId)
	userCanSeeAllProducts := embedCtx.Err == nil

	// validate user input params
	var channelIdOrSlug string
	if args.Channel != nil {
		if !slug.IsSlug(*args.Channel) && !model_helper.IsValidId(*args.Channel) {
			return nil, model_helper.NewAppError("Category.Products", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "channel"}, fmt.Sprintf("%s is not a channel slug nor id", *args.Channel), http.StatusBadRequest)
		}
		channelIdOrSlug = *args.Channel
	}

	if appErr := args.GraphqlParams.validate("Category.Products"); appErr != nil {
		return nil, appErr
	}

	products, appErr := embedCtx.App.Srv().
		ProductService().
		GetVisibleToUserProducts(channelIdOrSlug, userCanSeeAllProducts)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := func(p *model.Product) []any { return []any{model.ProductTableName + ".Slug", p.Slug} }
	res, appErr := newGraphqlPaginator(products, keyFunc, SystemProductToGraphqlProduct, args.GraphqlParams).parse("Category.Products")
	if appErr != nil {
		return nil, appErr
	}
	return (*ProductCountableConnection)(unsafe.Pointer(res)), nil
}
