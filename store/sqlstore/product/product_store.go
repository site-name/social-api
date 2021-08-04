package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductStore struct {
	store.Store
}

func NewSqlProductStore(s store.Store) store.ProductStore {
	ps := &SqlProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Product{}, store.ProductTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("DefaultVariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.PRODUCT_SLUG_MAX_LENGTH).SetUnique(true)

		s.CommonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlProductStore) ModelFields() []string {
	return []string{
		"Products.Id",
		"Products.ProductTypeID",
		"Products.Name",
		"Products.Slug",
		"Products.Description",
		"Products.DescriptionPlainText",
		"Products.CategoryID",
		"Products.CreateAt",
		"Products.UpdateAt",
		"Products.ChargeTaxes",
		"Products.Weight",
		"Products.WeightUnit",
		"Products.DefaultVariantID",
		"Products.Rating",
		"Products.Metadata",
		"Products.PrivateMetadata",
		"Products.SeoTitle",
		"Products.SeoDescription",
	}
}

func (ps *SqlProductStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_products_name", store.ProductTableName, "Name")
	ps.CreateIndexIfNotExists("idx_products_slug", store.ProductTableName, "Slug")
	ps.CreateIndexIfNotExists("idx_products_name_lower_textpattern", store.ProductTableName, "lower(Name) text_pattern_ops")

	ps.CommonMetaDataIndex(store.ProductTableName)
}

func (ps *SqlProductStore) Save(prd *product_and_discount.Product) (*product_and_discount.Product, error) {
	prd.PreSave()
	if err := prd.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(prd); err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Name", "products_name_key", "idx_products_name_unique"}) {
			return nil, store.NewErrInvalidInput("Product", "name", prd.Name)
		}
		if ps.IsUniqueConstraintError(err, []string{"Slug", "products_slug_key", "idx_products_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Product", "slug", prd.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save Product with productId=%s", prd.Id)
	}

	return prd, nil
}

func (ps *SqlProductStore) Get(id string) (*product_and_discount.Product, error) {
	var res product_and_discount.Product
	err := ps.GetMaster().SelectOne(&res, "SELECT * FROM "+store.ProductTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Product", id)
		}
		return nil, errors.Wrapf(err, "failed to get Product with productId=%s", id)
	}

	return &res, nil
}

func (ps *SqlProductStore) GetProductsByIds(ids []string) ([]*product_and_discount.Product, error) {
	sqlQuery, args, err := ps.GetQueryBuilder().
		Select("*").
		From(store.ProductTableName).
		Where(squirrel.Eq{"Id": ids}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_products_by_ids")
	}

	products := []*product_and_discount.Product{}
	if _, err := ps.GetMaster().Select(&products, sqlQuery, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Products")
	}

	return products, nil
}

// func (ps *SqlProductStore) FilterProducts(filter *model.ProductFilterInput) ([]*product_and_discount.Product, error) {
// 	// table names
// 	const (
// 		TABLE_products               = "Products"
// 		TABLE_collectionProducts     = "CollectionProducts"
// 		TABLE_productChannelListings = "ProductChannelListings"
// 		TABLE_channels               = "Channels"
// 		TABLE_productVariants        = "ProductVariants"
// 	)

// 	const (
// 		INNER_JOIN = " INNER JOIN "
// 	)

// 	query := "SELECT * FROM Products"

// 	if len(filter.Collections) > 0 {
// 		query += " INNER JOIN CollectionProducts ON (Products.Id = CollectionProducts.ProductID)"
// 	}

// 	// TODO: fixme
// 	panic("not implemented")
// }
