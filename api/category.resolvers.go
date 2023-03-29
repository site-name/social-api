package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
)

func (r *Resolver) CategoryCreate(ctx context.Context, args struct {
	Input  CategoryInput
	Parent *string
}) (*CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryDelete(ctx context.Context, args struct{ Id string }) (*CategoryDelete, error) {
	// if !model.IsValidId(args.Id) {
	// 	return nil, model.NewAppError("CategoryDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid category id", http.StatusBadRequest)
	// }

	// // check permission to delete category:
	// embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	// if err != nil {
	// 	return nil, err
	// }
	// if !r.srv.AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionDeleteCategory) {
	// 	return nil, model.NewAppError("CategoryDelete", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	// }

	// categories, appErr := r.srv.ProductService().CategoryByIds([]string{args.Id}, true)
	// if appErr != nil {
	// 	return nil, appErr
	// }
	// if categories.Len() == 0 {
	// 	return nil, model.NewAppError("CategoryDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid category id", http.StatusBadRequest)
	// }

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

// TODO: Add support filter by metadata
func (r *Resolver) Categories(ctx context.Context, args struct {
	Filter *CategoryFilterInput
	SortBy CategorySortingInput
	Level  *int32 // 0 <= level <= 4
	GraphqlParams
}) (*CategoryCountableConnection, error) {
	var levelFilter, searchFilter, idFilter, metadataFilter func(c *model.Category) bool

	// parse filter
	if args.Filter != nil {
		// parse search
		if search := args.Filter.Search; search != nil && *search != "" {
			lowSearch := strings.ToLower(*search)

			searchFilter = func(c *model.Category) bool {
				lowerSlug := strings.ToLower(c.Slug)
				lowerName := strings.ToLower(c.Name)
				return strings.Contains(lowerName, lowSearch) || strings.Contains(lowerSlug, lowSearch)
			}
		}

		// parse ids
		if ids := args.Filter.Ids; len(ids) > 0 {
			if !lo.EveryBy(ids, model.IsValidId) {
				return nil, model.NewAppError("Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Filter.Ids"}, "please provide valid uuids", http.StatusBadRequest)
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
		if *lv < model.CATEGORY_MIN_LEVEL || *lv > model.CATEGORY_MAX_LEVEL {
			return nil, model.NewAppError("Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Level"}, fmt.Sprintf("Level must be >= %d and <= %d", model.CATEGORY_MIN_LEVEL, model.CATEGORY_MAX_LEVEL), http.StatusBadRequest)
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

	// find categories:
	categories := r.srv.ProductService().FilterCategoriesFromCache(filter)

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
	if args.Id != nil && model.IsValidId(*args.Id) {
		id = *args.Id
	}
	if args.Slug != nil && *args.Slug != "" {
		slug = *args.Slug
	}
	if id == "" && slug == "" {
		return nil, model.NewAppError("Resolver.Category", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id/Slug"}, "please provide either id or slug", http.StatusBadRequest)
	}

	categories := r.srv.ProductService().FilterCategoriesFromCache(func(c *model.Category) bool {
		if id != "" {
			return c.Id == id
		}
		return c.Slug == slug
	})
	if categories.Len() == 0 {
		return nil, nil
	}

	return systemCategoryToGraphqlCategory(categories[0]), nil
}
