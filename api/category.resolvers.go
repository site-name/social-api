package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	SortBy *CategorySortingInput
	Level  *int
	GraphqlParams
}) (*CategoryCountableConnection, error) {
	var categories model.Categories // NOTE: don't modify inner categories

	if args.Level != nil {
		switch *args.Level {
		case 0:
			categories = model.FirstLevelCategories
		case 1:
			categories = model.SecondLevelCategories
		case 2:
			categories = model.ThirdLevelCategories
		case 3:
			categories = model.FourthLevelCategories
		case 4:
			categories = model.FifthhLevelCategories
		default:
			return nil, model.NewAppError("Resolver.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "level"}, "level can be in range [0, 4] only", http.StatusBadRequest)
		}
	}

	if categories.Len() == 0 {
		categories = model.AllCategories
	}

	if args.Filter != nil {
		if s := args.Filter.Search; s != nil && *s != "" {
			search := strings.ToLower(*s)

			categories = lo.Filter(categories, func(c *model.Category, _ int) bool {
				lowerName := strings.ToLower(c.Name)
				lowerNameEn := strings.ToLower(c.NameEn)

				return strings.Contains(lowerNameEn, search) ||
					strings.Contains(lowerName, search) ||
					strings.Contains(c.Slug, search)
			})
		}

		if ids := args.Filter.Ids; len(ids) > 0 {
			idMap := lo.SliceToMap(ids, func(item string) (string, struct{}) { return item, struct{}{} })
			categories = lo.Filter(categories, func(c *model.Category, _ int) bool {
				_, exist := idMap[c.Id]
				return exist
			})
		}
	}

	// default to sort by english name
	var keyFunc any = func(c *model.Category) string { return c.NameEn }
	if args.SortBy != nil && args.SortBy.Field.IsValid() {

		switch args.SortBy.Field {
		case CategorySortFieldSubcategoryCount:
			keyFunc = func(c *model.Category) int { return len(c.Children) }

		case CategorySortFieldProductCount:

		}
	}

	res, appErr := newGraphqlPaginator(categories, nil, systemCategoryToGraphqlCategory, args.GraphqlParams).parse("Resolver.Categories")
	if appErr != nil {
		return nil, appErr
	}
}

func (r *Resolver) Category(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*Category, error) {
	var id, slug string
	if args.Id != nil {
		id = *args.Id
	}
	if args.Slug != nil {
		slug = *args.Slug
	}
	if id == "" && slug == "" {
		return nil, model.NewAppError("Resolver.Category", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id/Slug"}, "please provide either id or slug", http.StatusBadRequest)
	}

	category := model.FirstLevelCategories.Search(func(c *model.Category) bool {
		return c.Id == id || c.Slug == slug
	})
	return systemCategoryToGraphqlCategory(category), nil
}
