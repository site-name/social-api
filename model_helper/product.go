package model_helper

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func ProductPreSave(p *model.Product) {
	productCommonPre(p)
	if p.ID == "" {
		p.ID = NewId()
	}
	p.CreatedAt = GetMillis()
	p.UpdatedAt = p.CreatedAt
	p.Slug = slug.Make(p.Name)
}

func productCommonPre(p *model.Product) {
	p.Name = SanitizeUnicode(p.Name)
	p.DescriptionPlainText = SanitizeUnicode(p.DescriptionPlainText)
	p.SeoTitle = SanitizeUnicode(p.SeoTitle)
	p.SeoDescription = SanitizeUnicode(p.SeoDescription)
	if p.ChargeTaxes.IsNil() {
		p.ChargeTaxes = model_types.NewNullBool(true)
	}
	if p.WeightUnit == "" {
		p.WeightUnit = measurement.G.String()
	}
}

func ProductPreUpdate(p *model.Product) {
	p.UpdatedAt = GetMillis()
	productCommonPre(p)
}

func ProductIsValid(p model.Product) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(p.CategoryID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.category_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.CreatedAt <= 0 {
		return NewAppError("Product.IsValid", "model.product.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if p.UpdatedAt <= 0 {
		return NewAppError("Product.IsValid", "model.product.is_valid.updated_at.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.DefaultVariantID.IsNil() && !IsValidId(*p.DefaultVariantID.String) {
		return NewAppError("Product.IsValid", "model.product.is_valid.default_variant_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Name == "" {
		return NewAppError("Product.IsValid", "model.product.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(p.Slug) {
		return NewAppError("Product.IsValid", "model.product.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}
	if measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(p.WeightUnit)] == "" {
		return NewAppError("Product.IsValid", "model.product.is_valid.weight_unit.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type ProductFilterOptions struct {
	CommonQueryOptions
}

func CategoryPreSave(c *model.Category) {
	if c.ID == "" {
		c.ID = NewId()
	}
	c.Slug = slug.Make(c.Name)
	CategoryCommonPre(c)
}

func CategoryCommonPre(c *model.Category) {
	c.Name = SanitizeUnicode(c.Name)
	c.SeoTitle = SanitizeUnicode(c.SeoTitle)
	c.SeoDescription = SanitizeUnicode(c.SeoDescription)
}

func CategoryIsValid(c model.Category) *AppError {
	if !IsValidId(c.ID) {
		return NewAppError("Category.IsValid", "model.category.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if c.Name == "" {
		return NewAppError("Category.IsValid", "model.category.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(c.Slug) {
		return NewAppError("Category.IsValid", "model.category.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}
	if !c.ParentID.IsNil() && !IsValidId(*c.ParentID.String) {
		return NewAppError("Category.IsValid", "model.category.is_valid.parent_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func ProductVariantPreSave(pv *model.ProductVariant) {
	if pv.ID == "" {
		pv.ID = NewId()
	}
	ProductVariantCommonPre(pv)
}

func ProductVariantCommonPre(pv *model.ProductVariant) {
	pv.Name = SanitizeUnicode(pv.Name)
	if pv.TrackInventory.IsNil() {
		pv.TrackInventory = model_types.NewNullBool(true)
	}
	if !pv.Weight.IsNil() && pv.WeightUnit == "" {
		pv.WeightUnit = measurement.G.String()
	}
}

func ProductVariantIsValid(p model.ProductVariant) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.product_id.app_error", nil, "", http.StatusBadRequest)
	}
	if p.Name == "" {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.Weight.IsNil() && *p.Weight.Float32 < 0 {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight.app_error", nil, "", http.StatusBadRequest)
	}
	if !p.Weight.IsNil() && measurement.WEIGHT_UNIT_CONVERSION[measurement.WeightUnit(p.WeightUnit)] == 0 {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight_unit.app_error", nil, "", http.StatusBadRequest)
	}
	if p.IsPreorder && (p.PreorderEndDate.IsNil() || *p.PreorderEndDate.Int64 < GetMillis()) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.preorder_end_date.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func CollectionPreSave(c *model.Collection) {
	if c.ID == "" {
		c.ID = NewId()
	}
	CollectionCommonPre(c)
	c.Slug = slug.Make(c.Name)
}

func CollectionCommonPre(c *model.Collection) {
	c.Name = SanitizeUnicode(c.Name)
	c.SeoTitle = SanitizeUnicode(c.SeoTitle)
	c.SeoDescription = SanitizeUnicode(c.SeoDescription)
}

func CollectionIsValid(c model.Collection) *AppError {
	if !IsValidId(c.ID) {
		return NewAppError("Collection.IsValid", "model.collection.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if c.Name == "" {
		return NewAppError("Collection.IsValid", "model.collection.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(c.Slug) {
		return NewAppError("Collection.IsValid", "model.collection.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type CategoryFilterOption struct {
	CommonQueryOptions
	SaleID    qm.QueryMod // INNER JOIN SaleCategories ON ... WHERE SaleCategories.SaleID ...
	ProductID qm.QueryMod // INNER JOIN Products ON ... WHERE Products.ID ...
	VoucherID qm.QueryMod // INNER JOIN VoucherCategories ON ... WHERE VoucherCategories.VoucherID ...
}
