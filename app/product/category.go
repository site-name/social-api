package product

import (
	"net/http"
	"sort"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (s *ServiceProduct) FilterCategoriesFromCache(filter func(c *model.Category) bool) model.CategorySlice {
	var res model.CategorySlice

	s.categoryMap.Range(func(id, value any) bool {
		if category := value.(*model.Category); filter(category) {
			res = append(res, category)
		}
		return true
	})

	return res
}

func (s *ServiceProduct) CategoryByIds(ids []string, allowFromCache bool) (model.CategorySlice, *model_helper.AppError) {
	if allowFromCache {
		var res model.CategorySlice
		notFoundCategoryIdMap := lo.SliceToMap(ids, func(id string) (string, struct{}) { return id, struct{}{} })

		s.categoryMap.Range(func(id, value any) bool {
			_, ok := notFoundCategoryIdMap[id.(string)]
			if ok {
				res = append(res, value.(*model.Category))
				delete(notFoundCategoryIdMap, id.(string))
			}
			return true
		})

		if len(notFoundCategoryIdMap) > 0 {
			categories, appErr := s.CategoriesByOption(&model.CategoryFilterOption{
				Conditions: squirrel.Eq{model.CategoryTableName + ".id": lo.Keys(notFoundCategoryIdMap)},
			})
			if appErr != nil {
				return nil, appErr
			}
			res = append(res, categories...)
		}

		return res, nil
	}

	return s.CategoriesByOption(&model.CategoryFilterOption{
		Conditions: squirrel.Eq{model.CategoryTableName + ".id": ids},
	})
}

// CategoriesByOption returns all categories that satisfy given option
func (a *ServiceProduct) CategoriesByOption(option *model.CategoryFilterOption) (model.CategorySlice, *model_helper.AppError) {
	categories, err := a.srv.Store.Category().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("CategoriesByOption", "app.product.error_finding_categories_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return categories, nil
}

// CategoryByOption returns 1 category that satisfies given option
func (a *ServiceProduct) CategoryByOption(option *model.CategoryFilterOption) (*model.Category, *model_helper.AppError) {
	category, err := a.srv.Store.Category().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("CategoryByOption", "app.product.error_finding_category_by_option.app_error", nil, err.Error(), statusCode)
	}

	return category, nil
}

// DoAnalyticCategories finds all categories in system.
// Counts number of products of each category.
// Sets NumberOfProducts, NumberOfChildren, Children attributes of each category.
// Stores classified categories in cache
func (s *ServiceProduct) DoAnalyticCategories() *model_helper.AppError {
	slog.Info("Analyzing categories")

	var allCategories model.CategorySlice
	const limit = 500
	var lastCategorySlug string

	for {
		filterOpts := &model.CategoryFilterOption{
			Limit:   limit,
			OrderBy: model.CategoryTableName + ".Slug ASC",
		}
		if lastCategorySlug != "" {
			filterOpts.Conditions = squirrel.Gt{model.CategoryTableName + ".Slug": lastCategorySlug}
		}
		categories, appErr := s.CategoriesByOption(filterOpts)
		if appErr != nil {
			return appErr
		}

		lastCategorySlug = categories[categories.Len()-1].Slug
		allCategories = append(allCategories, categories...)

		if len(categories) < limit {
			break
		}
	}

	countObjs, err := s.srv.Store.Product().CountByCategoryIDs(allCategories.IDs(false))
	if err != nil {
		return model_helper.NewAppError("DoAnalyticCategories", "app.product.counting_products_by_category_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	categoryMap := lo.SliceToMap(allCategories, func(c *model.Category) (string, *model.Category) { return c.Id, c })
	for _, count := range countObjs {
		cate, ok := categoryMap[count.CategoryID]
		if ok && cate != nil {
			cate.NumOfProducts = count.ProductCount
		}
	}
	allCategories = s.ClassifyCategories(allCategories)
	for _, cate := range allCategories {
		s.categoryMap.Store(cate.Id, cate)
	}

	return nil
}

// ClassifyCategories takes a slice of single categories.
// Returns a slice of category families
// NOTE: you can call this function
func (s *ServiceProduct) ClassifyCategories(categories model.CategorySlice) model.CategorySlice {
	if len(categories) <= 1 {
		return categories
	}

	var res model.CategorySlice
	sort.SliceStable(categories, func(i, j int) bool {
		return categories[i].Level > categories[j].Level
	})
	var categoryMap = lo.SliceToMap(categories, func(c *model.Category) (string, *model.Category) { return c.Id, c })

	for _, cate := range categories {
		if cate == nil {
			continue
		}
		res = append(res, cate)

		if cate.ParentID != nil {
			parent, ok := categoryMap[*cate.ParentID]
			if ok && parent != nil {
				// parent.Children = append(parent.Children, cate)
				parent.NumOfChildren++
				parent.NumOfProducts += cate.NumOfProducts
			}
		}
	}

	return res
}

// UpsertCategory first checks if given category need a Level number.
// Performs upsert given category into database.
// asynchronously does category anayltic to update category cache.
func (s *ServiceProduct) UpsertCategory(cate *model.Category) (*model.Category, *model_helper.AppError) {
	if !model_helper.IsValidId(cate.Id) && cate.ParentID != nil { // meaning saving category
		parentCate, ok := s.categoryMap.Load(*cate.ParentID)
		if ok && parentCate != nil {
			cate.Level = parentCate.(*model.Category).Level + 1
		}
	}

	cate, err := s.srv.Store.Category().Upsert(cate)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("UpsertCategory", "app.product.upsert_category.app_error", nil, err.Error(), statusCode)
	}

	s.srv.Go(func() {
		appErr := s.DoAnalyticCategories()
		if appErr != nil {
			slog.Error("failed to do category analyse", slog.Err(appErr))
		}
	})

	return cate, nil
}
