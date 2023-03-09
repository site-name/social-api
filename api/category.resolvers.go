package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) CategoryCreate(ctx context.Context, args struct {
	Input  CategoryInput
	Parent *string
}) (*CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryDelete(ctx context.Context, args struct{ Id string }) (*CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryBulkDelete(ctx context.Context, args struct{ Ids []string }) (*CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryUpdate(ctx context.Context, args struct {
	Id    string
	Input CategoryInput
}) (*CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
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
	Level  *int32 // 0 <= level <= 4
	GraphqlParams
}) (*CategoryCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	var levelFilter, searchFilter, idFilter func(c *model.Category) bool

	if lv := args.Level; lv != nil && *lv >= model.CATEGORY_MIN_LEVEL &&
		*lv <= model.CATEGORY_MAX_LEVEL {
		levelFilter = func(c *model.Category) bool {
			return c.Level == uint8(*lv)
		}
	}

	if args.Filter != nil {
		if search := args.Filter.Search; search != nil && *search != "" {
			lowSearch := strings.ToLower(*search)

			searchFilter = func(c *model.Category) bool {
				lowerSlug := strings.ToLower(c.Slug)
				lowerName := strings.ToLower(c.Name)
				return strings.Contains(lowerName, lowSearch) || strings.Contains(lowerSlug, lowSearch)
			}
		}

		if len(args.Filter.Ids) > 0 {
			idMap := map[string]struct{}{}
			for _, id := range args.Filter.Ids {
				idMap[id] = struct{}{}
			}

			idFilter = func(c *model.Category) bool {
				_, ok := idMap[c.Id]
				return ok
			}
		}
	}

	filter := func(c *model.Category) bool {
		return (levelFilter != nil && levelFilter(c)) ||
			(searchFilter != nil && searchFilter(c)) ||
			(idFilter != nil && idFilter(c))
	}

	// find categories:
	categories := embedCtx.App.Srv().ProductService().FilterCategoriesFromCache(filter)

	// default to sort by english name
	var res *CountableConnection[*Category]
	var appErr *model.AppError

	switch args.SortBy.Field {
	case CategorySortFieldSubcategoryCount:
		keyFunc := func(c *model.Category) int { return c.NumOfChildren }
		res, appErr = newGraphqlPaginator(categories, keyFunc, systemCategoryToGraphqlCategory, args.GraphqlParams).parse("Resolver.Categories")

	case CategorySortFieldProductCount:
		keyFunc := func(c *model.Category) uint64 { return c.NumOfProducts }
		res, appErr = newGraphqlPaginator(categories, keyFunc, systemCategoryToGraphqlCategory, args.GraphqlParams).parse("Resolver.Categories")

	default:
		keyFunc := func(c *model.Category) string { return c.Name }
		if args.SortBy.Field != CategorySortFieldName {
			keyFunc = func(c *model.Category) string { return c.Slug }
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
	var id, slug string
	if args.Id != nil && IdsAreValidUUIDs(*args.Id) {
		id = *args.Id
	}
	if args.Slug != nil {
		slug = *args.Slug
	}
	if id == "" && slug == "" {
		return nil, model.NewAppError("Resolver.Category", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id/Slug"}, "please provide either id or slug", http.StatusBadRequest)
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	category, appErr := embedCtx.App.Srv().ProductService().CategoryByOption(&model.CategoryFilterOption{
		Extra: squirrel.Or{
			squirrel.Eq{store.CategoryTableName + ".Id": id},
			squirrel.Eq{store.CategoryTableName + ".Slug": slug},
		},
	})
	if appErr != nil {
		return nil, appErr
	}
	return systemCategoryToGraphqlCategory(category), nil
}
