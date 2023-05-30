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
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) CategoryCreate(ctx context.Context, args struct {
	Input  CategoryInput
	Parent *string
}) (*CategoryCreate, error) {
	// check user permissions
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoles("CategoryCreate", model.ShopStaffRoleId)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate parent
	if pr := args.Parent; pr != nil && !model.IsValidId(*pr) {
		return nil, model.NewAppError("CategoryCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "parent"}, fmt.Sprintf("%s is not a valid category id", *pr), http.StatusBadRequest)
	}
	if appErr := args.Input.Validate(); appErr != nil {
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
	// requester must have delete category permission to delete
	// embedCtx:= GetContextValue[*web.Context](ctx, WebCtx)
	// embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionDeleteCategory)
	// if embedCtx.Err != nil {
	// 	return nil, embedCtx.Err
	// }

	// if !model.IsValidId(args.Id) {
	// 	return nil, model.NewAppError("CategoryDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid category id", http.StatusBadRequest)
	// }

	// embedCtx.App.Srv().ProductService().DeleteCategories()
	panic("not implemented")
}

func (r *Resolver) CategoryBulkDelete(ctx context.Context, args struct{ Ids []string }) (*CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryUpdate(ctx context.Context, args struct {
	Id    string
	Input CategoryInput
}) (*CategoryUpdate, error) {
	// requester must be authenticated and has category_update permission to do this
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoles("CategoryUpdate", model.ShopStaffRoleId)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate given id
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("CategoryUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}
	if appErr := args.Input.Validate(); appErr != nil {
		return nil, appErr
	}

	categories, appErr := embedCtx.App.Srv().ProductService().CategoryByIds([]string{args.Id}, true)
	if appErr != nil {
		return nil, appErr
	}
	if categories.Len() == 0 {
		return nil, model.NewAppError("CategoryUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("category with id=%s not found", args.Id), http.StatusBadRequest)
	}

	category := categories[0]
	args.Input.PatchCategory(category)

	category, appErr = r.srv.ProductService().UpsertCategory(category)
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
	if args.Id == nil && args.Slug == nil {
		return nil, model.NewAppError("Category", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id and slug"}, "id or slug must be provided", http.StatusBadRequest)
	}
	if args.Id != nil && model.IsValidId(*args.Id) {
		return nil, model.NewAppError("Category", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprint("%s is invalid id", *args.Id), http.StatusBadRequest)
	}
	if args.Slug != nil && !slug.IsSlug(*args.Slug) {
		return nil, model.NewAppError("Category", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, fmt.Sprint("%s is invalid slug", *args.Slug), http.StatusBadRequest)
	}

	categories := r.srv.ProductService().FilterCategoriesFromCache(func(c *model.Category) bool {
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
