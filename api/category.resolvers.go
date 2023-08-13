package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/category.graphqls for details on directive used.
func (r *Resolver) CategoryCreate(ctx context.Context, args struct {
	Input  CategoryInput
	Parent *string
}) (*CategoryCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate parent
	if pr := args.Parent; pr != nil && !model.IsValidId(*pr) {
		return nil, model.NewAppError("CategoryCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "parent"}, fmt.Sprintf("%s is not a valid category id", *pr), http.StatusBadRequest)
	}
	if appErr := args.Input.Validate("CategoryCreate"); appErr != nil {
		return nil, appErr
	}

	// construct category instance
	category := new(model.Category)
	args.Input.PatchCategory(category)

	// save category
	category, appErr := embedCtx.App.Srv().ProductService().UpsertCategory(category)
	if appErr != nil {
		return nil, appErr
	}

	// TODO: check if we need create bg image thumbnail

	return &CategoryCreate{
		Category: systemCategoryToGraphqlCategory(category),
	}, nil
}

func (r *Resolver) CategoryDelete(ctx context.Context, args struct{ Id string }) (*CategoryDelete, error) {
	panic("not implemented")
}

func (r *Resolver) CategoryBulkDelete(ctx context.Context, args struct{ Ids []string }) (*CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/category.graphqls for details on directive used.
func (r *Resolver) CategoryUpdate(ctx context.Context, args struct {
	Id    string
	Input CategoryInput
}) (*CategoryUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate given id
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("CategoryUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}
	if appErr := args.Input.Validate("CategoryCreate"); appErr != nil {
		return nil, appErr
	}

	categories, appErr := embedCtx.App.Srv().ProductService().CategoryByIds([]string{args.Id}, true)
	if appErr != nil {
		return nil, appErr
	}
	if categories.Len() == 0 {
		return nil, model.NewAppError("CategoryUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("category with id=%s not found", args.Id), http.StatusBadRequest)
	}

	category := categories[0]
	args.Input.PatchCategory(category)

	category, appErr = embedCtx.App.Srv().ProductService().UpsertCategory(category)
	if appErr != nil {
		return nil, appErr
	}

	return &CategoryUpdate{
		Category: systemCategoryToGraphqlCategory(category),
	}, nil
}

func (r *Resolver) CategoryTranslate(ctx context.Context, args struct {
	Id           string
	Input        TranslationInput
	LanguageCode LanguageCodeEnum
}) (*CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Categories(ctx context.Context, args struct {
	Filter *CategoryFilterInput
	SortBy CategorySortingInput
	Level  *int32 // 0 <= level
	GraphqlParams
}) (*CategoryCountableConnection, error) {
	var levelFilter, searchFilter, idFilter, metadataFilter func(c *model.Category) bool

	// parse filter
	if args.Filter != nil {
		// parse search
		if search := args.Filter.Search; search != nil && *search != "" {
			lowerSearch := strings.ToLower(*search)

			searchFilter = func(c *model.Category) bool {
				lowerSlug := strings.ToLower(c.Slug)
				lowerName := strings.ToLower(c.Name)
				return strings.Contains(lowerName, lowerSearch) || strings.Contains(lowerSlug, lowerSearch)
			}
		}

		// parse ids
		if ids := args.Filter.Ids; len(ids) > 0 {
			if !lo.EveryBy(ids, model.IsValidId) {
				return nil, model.NewAppError("Categories", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Filter.Ids"}, "please provide valid uuids", http.StatusBadRequest)
			}
			idMap := map[string]bool{}
			for _, id := range ids {
				idMap[id] = true
			}
			idFilter = func(c *model.Category) bool {
				return idMap[c.Id]
			}
		}

		// parse meta
		if metas := args.Filter.Metadata; len(metas) > 0 {
			metadataFilter = func(c *model.Category) bool {
				for _, meta := range metas {
					if meta.Key != "" {
						if meta.Value != "" {
							value := c.Metadata[meta.Key]
							if value == meta.Value {
								return true
							}
							continue
						}

						if _, ok := c.Metadata[meta.Key]; ok {
							return true
						}
					}
				}
				return false
			}
		}
	}

	// parse level
	if lv := args.Level; lv != nil {
		if *lv < model.CATEGORY_MIN_LEVEL {
			return nil, model.NewAppError("Categories", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Level"}, fmt.Sprintf("Level must be >= %d", model.CATEGORY_MIN_LEVEL), http.StatusBadRequest)
		}
		levelFilter = func(c *model.Category) bool {
			return c.Level == uint8(*lv)
		}
	}

	noNeedFilter := levelFilter == nil &&
		searchFilter == nil &&
		idFilter == nil &&
		metadataFilter == nil

	filter := func(c *model.Category) bool {
		return noNeedFilter ||
			(levelFilter != nil && levelFilter(c)) ||
			(searchFilter != nil && searchFilter(c)) ||
			(idFilter != nil && idFilter(c)) ||
			(metadataFilter != nil && metadataFilter(c))
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// find categories:
	categories := embedCtx.App.Srv().ProductService().FilterCategoriesFromCache(filter)

	// default to sort by english name
	var res *CountableConnection[*Category]
	var appErr *model.AppError

	switch args.SortBy.Field {
	case CategorySortFieldSubcategoryCount:
		keyFunc := func(c *model.Category) []any {
			return []any{model.CategoryTableName + ".NumOfChildren", c.NumOfChildren}
		}
		res, appErr = newGraphqlPaginator(categories, keyFunc, systemCategoryToGraphqlCategory, args.GraphqlParams).parse("Resolver.Categories")

	case CategorySortFieldProductCount:
		keyFunc := func(c *model.Category) []any {
			return []any{model.CategoryTableName + ".NumOfProducts", c.NumOfProducts}
		}
		res, appErr = newGraphqlPaginator(categories, keyFunc, systemCategoryToGraphqlCategory, args.GraphqlParams).parse("Resolver.Categories")

	default:
		keyFunc := func(c *model.Category) []any { return []any{model.CategoryTableName + ".Name", c.Name} }
		if args.SortBy.Field != CategorySortFieldName {
			keyFunc = func(c *model.Category) []any { return []any{model.CategoryTableName + ".Slug", c.Slug} }
		}
		res, appErr = newGraphqlPaginator(categories, keyFunc, systemCategoryToGraphqlCategory, args.GraphqlParams).parse("Resolver.Categories")
	}
	if appErr != nil {
		return nil, appErr
	}
	return (*CategoryCountableConnection)(unsafe.Pointer(res)), nil
}

func (r *Resolver) Category(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*Category, error) {
	if args.Id == nil && args.Slug == nil {
		return nil, model.NewAppError("Category", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id and slug"}, "id or slug must be provided", http.StatusBadRequest)
	}
	if args.Id != nil && model.IsValidId(*args.Id) {
		return nil, model.NewAppError("Category", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", *args.Id), http.StatusBadRequest)
	}
	if args.Slug != nil && !slug.IsSlug(*args.Slug) {
		return nil, model.NewAppError("Category", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, fmt.Sprintf("%s is invalid slug", *args.Slug), http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	categories := embedCtx.App.Srv().ProductService().FilterCategoriesFromCache(func(c *model.Category) bool {
		if args.Id != nil {
			return c.Id == *args.Id
		}
		return c.Slug == *args.Slug
	})
	if categories.Len() == 0 {
		return nil, nil
	}

	return systemCategoryToGraphqlCategory(categories[0]), nil
}
