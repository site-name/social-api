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
	ID              string          `json:"id"`
	SeoTitle        *string         `json:"seoTitle"`
	SeoDescription  *string         `json:"seoDescription"`
	Name            string          `json:"name"`
	Description     JSONString      `json:"description"`
	Slug            string          `json:"slug"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`

	parentID *string

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
		ID:              c.Id,
		SeoTitle:        c.SeoTitle,
		SeoDescription:  c.SeoDescription,
		Name:            c.Name,
		Description:     JSONString(c.Description),
		Slug:            c.Slug,
		Metadata:        MetadataToSlice(c.Metadata),
		PrivateMetadata: MetadataToSlice(c.PrivateMetadata),

		parentID: c.ParentID,
	}
}

func (c *Category) BackgroundImage(ctx context.Context, args struct{ Size *int32 }) (*Image, error) {
	panic("not implemented")
}

func (c *Category) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*CategoryTranslation, error) {
	panic("not implemented")
}

func (c *Category) Parent(ctx context.Context) (*Category, error) {
	panic("not implemented")
}

func (c *Category) Ancestors(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CategoryCountableConnection, error) {
	panic("not implemented")
}

func (c *Category) Children(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CategoryCountableConnection, error) {
	panic("not implemented")
}

func (c *Category) Products(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductCountableConnection, error) {
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

func productsByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Product] {
	var (
		res             = make([]*dataloader.Result[[]*model.Product], len(voucherIDs))
		products        model.Products
		appErr          *model.AppError
		voucherProducts []*model.VoucherProduct

		voucherProductMap = map[string]string{}         // keys are product ids, values are voucher ids
		productMap        = map[string]model.Products{} // keys are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	products, appErr = embedCtx.App.Srv().ProductService().ProductsByOption(&model.ProductFilterOption{
		VoucherID: squirrel.Eq{store.VoucherProductTableName + ".VoucherID": voucherIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	voucherProducts, err = embedCtx.App.Srv().Store.VoucherProduct().FilterByOptions(&model.VoucherProductFilterOptions{
		VoucherID: squirrel.Eq{store.VoucherProductTableName + ".VoucherID": voucherIDs},
	})
	if err != nil {
		err = model.NewAppError("productsByVoucherIDLoader", "app.discount.voucher_products_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range voucherProducts {
		voucherProductMap[rel.ProductID] = rel.VoucherID
	}

	for _, prd := range products {
		voucherID, ok := voucherProductMap[prd.Id]
		if ok {
			productMap[voucherID] = append(productMap[voucherID], prd)
		}
	}

	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Data: productMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
	}
	return res
}

func productVariantsByVoucherIdLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res        = make([]*dataloader.Result[[]*model.ProductVariant], len(voucherIDs))
		variants   model.ProductVariants
		appErr     *model.AppError
		variantMap = map[string]model.ProductVariants{} // keys are voucher ids

		variantVouchers   []*model.VoucherProductVariant
		variantVoucherMap = map[string]string{} // keys are variant ids, values are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variants, appErr = embedCtx.App.Srv().ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			VoucherID: squirrel.Eq{store.VoucherProductVariantTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	variantVouchers, err = embedCtx.App.Srv().Store.VoucherProductVariant().FilterByOptions(&model.VoucherProductVariantFilterOption{
		VoucherID: squirrel.Eq{store.VoucherProductVariantTableName + ".VoucherID": voucherIDs},
	})
	if err != nil {
		err = model.NewAppError("productVariantsByVoucherIdLoader", "app.discount.voucher_variants_reations_by_options", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range variantVouchers {
		variantVoucherMap[rel.ProductVariantID] = rel.VoucherID
	}

	for _, variant := range variants {
		voucherID, ok := variantVoucherMap[variant.Id]
		if ok {
			variantMap[voucherID] = append(variantMap[voucherID], variant)
		}
	}

	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: variantMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
	}
	return res
}

func categoriesBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Category] {
	var (
		res             = make([]*dataloader.Result[[]*model.Category], len(saleIDs))
		categories      model.Categories
		appErr          *model.AppError
		saleCategories  []*model.SaleCategoryRelation
		categoryMap     = map[string]model.Categories{} // keys are sale ids
		saleCategoryMap = map[string]string{}           // keys are category ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	categories, appErr = embedCtx.App.Srv().
		ProductService().
		CategoriesByOption(&model.CategoryFilterOption{
			SaleID: squirrel.Eq{store.SaleCategoryRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleCategories, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleCategoriesByOption(&model.SaleCategoryRelationFilterOption{
			SaleID: squirrel.Eq{store.SaleCategoryRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleCategories {
		saleCategoryMap[rel.CategoryID] = rel.SaleID
	}

	for _, cate := range categories {
		saleID, ok := saleCategoryMap[cate.Id]
		if ok {
			categoryMap[saleID] = append(categoryMap[saleID], cate)
		}
	}

	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Data: categoryMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Error: err}
	}
	return res
}

func collectionsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Collection] {
	var (
		res               = make([]*dataloader.Result[[]*model.Collection], len(saleIDs))
		collections       model.Collections
		appErr            *model.AppError
		saleCollections   []*model.SaleCollectionRelation
		collectionMap     = map[string]model.Collections{} // keys are sale ids
		saleCollectionMap = map[string]string{}            // keys are collection ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collections, appErr = embedCtx.App.Srv().
		ProductService().
		CollectionsByOption(&model.CollectionFilterOption{
			SaleID: squirrel.Eq{store.SaleCollectionRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleCollections, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleCollectionsByOptions(&model.SaleCollectionRelationFilterOption{
			SaleID: squirrel.Eq{store.SaleCollectionRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleCollections {
		saleCollectionMap[rel.CollectionID] = rel.SaleID
	}

	for _, collection := range collections {
		saleID, ok := saleCollectionMap[collection.Id]
		if ok {
			collectionMap[saleID] = append(collectionMap[saleID], collection)
		}
	}

	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Data: collectionMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Error: err}
	}
	return res
}

func productsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Product] {
	var (
		res            = make([]*dataloader.Result[[]*model.Product], len(saleIDs))
		products       model.Products
		appErr         *model.AppError
		saleProducts   []*model.SaleProductRelation
		productMap     = map[string]model.Products{} // keys are sale ids
		saleProductMap = map[string]string{}         // keys are product ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	products, appErr = embedCtx.App.Srv().
		ProductService().
		ProductsByOption(&model.ProductFilterOption{
			SaleID: squirrel.Eq{store.SaleProductRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleProducts, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleProductsByOptions(&model.SaleProductRelationFilterOption{
			SaleID: squirrel.Eq{store.SaleProductRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleProducts {
		saleProductMap[rel.ProductID] = rel.SaleID
	}

	for _, product := range products {
		saleID, ok := saleProductMap[product.Id]
		if ok {
			productMap[saleID] = append(productMap[saleID], product)
		}
	}

	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Data: productMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
	}
	return res
}

func productVariantsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res                   = make([]*dataloader.Result[[]*model.ProductVariant], len(saleIDs))
		variants              model.ProductVariants
		appErr                *model.AppError
		saleVariants          []*model.SaleProductVariant
		productVariantMap     = map[string]model.ProductVariants{} // keys are sale ids
		saleProductVariantMap = map[string]string{}                // keys are product variant ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variants, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			SaleID: squirrel.Eq{store.SaleProductVariantTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleVariants, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleProductVariantsByOptions(&model.SaleProductVariantFilterOption{
			SaleID: squirrel.Eq{store.SaleProductVariantTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleVariants {
		saleProductVariantMap[rel.ProductVariantID] = rel.SaleID
	}

	for _, product := range variants {
		saleID, ok := saleProductVariantMap[product.Id]
		if ok {
			productVariantMap[saleID] = append(productVariantMap[saleID], product)
		}
	}

	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: productVariantMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
	}
	return res
}
