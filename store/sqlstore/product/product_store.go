package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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

// Save inserts given product into database then returns it
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

// FilterByOption finds and returns all products that satisfy given option
func (ps *SqlProductStore) FilterByOption(option *product_and_discount.ProductFilterOption) ([]*product_and_discount.Product, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductTableName).
		OrderBy(store.TableOrderingMap[store.ProductTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Products.Id"))
	}
	if option.ProductVariantID != nil {
		query = query.
			LeftJoin(store.ProductVariantTableName + " ON (Products.Id = ProductVariants.ProductID)").
			Where(option.ProductVariantID.ToSquirrel("ProductVariants.Id"))
	}
	if option.VoucherID != nil {
		query = query.Where(option.VoucherID.ToSquirrel("")) // ne need to provide key value here
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.Product
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given option")
	}

	return res, nil
}

// GetByOption finds and returns 1 product that satisfies given option
func (ps *SqlProductStore) GetByOption(option *product_and_discount.ProductFilterOption) (*product_and_discount.Product, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductTableName).
		OrderBy(store.TableOrderingMap[store.ProductTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Products.Id"))
	}
	if option.ProductVariantID != nil {
		query = query.
			LeftJoin(store.ProductVariantTableName + " ON (Products.Id = ProductVariants.ProductID)").
			Where(option.ProductVariantID.ToSquirrel("ProductVariants.Id"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res product_and_discount.Product
	err = ps.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find product by given option")
	}

	return &res, nil
}

// ProductsByVoucherID finds all products that have relationships with given voucher
func (ps *SqlProductStore) ProductsByVoucherID(voucherID string) ([]*product_and_discount.Product, error) {
	return ps.FilterByOption(&product_and_discount.ProductFilterOption{
		VoucherID: &model.StringFilter{
			StringOption: &model.StringOption{
				ExtraExpr: []squirrel.Sqlizer{
					squirrel.Expr("Products.Id IN (SELECT ProductID FROM ? WHERE VoucherID = ?)", store.VoucherProductTableName, voucherID),
				},
			},
		},
	})
}

// FilterPublishedProducts finds and returns products that belong to given channel slug and are published
func (ps *SqlProductStore) FilterPublishedProducts(channelSlug string) ([]*product_and_discount.Product, error) {

	channelBySlugAndExist := ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where("Channels.IsActive AND Channels.Slug = ? AND Channels.Id = ProductChannelListings.ChannelID", channelSlug).
		Limit(1)
}
