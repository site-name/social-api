package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductTypeStore struct {
	store.Store
}

func NewSqlProductTypeStore(s store.Store) store.ProductTypeStore {
	pts := &SqlProductTypeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductType{}, store.ProductTypeTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_TYPE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(product_and_discount.PRODUCT_TYPE_SLUG_MAX_LENGTH)
		table.ColMap("Kind").SetMaxSize(product_and_discount.PRODUCT_TYPE_KIND_MAX_LENGTH)
	}
	return pts
}

func (ps *SqlProductTypeStore) ModelFields() []string {
	return []string{
		"ProductTypes.Id",
		"ProductTypes.Name",
		"ProductTypes.Slug",
		"ProductTypes.Kind",
		"ProductTypes.HasVariants",
		"ProductTypes.IsShippingRequired",
		"ProductTypes.IsDigital",
		"ProductTypes.Weight",
		"ProductTypes.WeightUnit",
		"ProductTypes.Metadata",
		"ProductTypes.PrivateMetadata",
	}
}

func (ps *SqlProductTypeStore) ScanFields(productType product_and_discount.ProductType) []interface{} {
	return []interface{}{
		&productType.Id,
		&productType.Name,
		&productType.Slug,
		&productType.Kind,
		&productType.HasVariants,
		&productType.IsShippingRequired,
		&productType.IsDigital,
		&productType.Weight,
		&productType.WeightUnit,
		&productType.Metadata,
		&productType.PrivateMetadata,
	}
}

func (ps *SqlProductTypeStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_types_name", store.ProductTypeTableName, "Name")
	ps.CreateIndexIfNotExists("idx_product_types_name_lower_textpattern", store.ProductTypeTableName, "lower(Name) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_product_types_slug", store.ProductTypeTableName, "Slug")
}

func (ps *SqlProductTypeStore) Save(productType *product_and_discount.ProductType) (*product_and_discount.ProductType, error) {
	productType.PreSave()
	if err := productType.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(productType); err != nil {
		return nil, errors.Wrapf(err, "failed to save product type withh id=%s", productType.Id)
	}

	return productType, nil
}

func (ps *SqlProductTypeStore) Get(id string) (*product_and_discount.ProductType, error) {
	var res product_and_discount.ProductType
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductTypeTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTypeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find product type with id=%s", id)
	}

	return &res, nil
}

func (ps *SqlProductTypeStore) FilterProductTypesByCheckoutID(checkoutToken string) ([]*product_and_discount.ProductType, error) {
	/*
					checkout
					|      |
		...<--|		   |--> checkoutLine <-- productVariant <-- product <-- productType
																							|												     |
													 ...checkoutLine <--|              ...product <--|
	*/

	rows, err := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductTypeTableName).
		InnerJoin(store.ProductTableName + " ON (ProductTypes.Id = Products.ProductTypeID)").
		InnerJoin(store.ProductVariantTableName + " ON (ProductVariants.ProductID = Products.Id)").
		InnerJoin(store.CheckoutLineTableName + " ON (CheckoutLines.VariantID = ProductVariants.Id)").
		InnerJoin(store.CheckoutTableName + " ON (Checkouts.Token = CheckoutLines.CheckoutID)").
		Where(squirrel.Eq{"Checkouts.Token": checkoutToken}).
		RunWith(ps.GetReplica()).
		Query()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find product types belong to given checkout with id=%s", checkoutToken)
	}

	var (
		productTypes []*product_and_discount.ProductType
		productType  product_and_discount.ProductType
		scanFields   = ps.ScanFields(productType)
	)

	for rows.Next() {
		err := rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse a row")
		}
		productTypes = append(productTypes, &productType)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "error closing rows")
	}
	if rows.Err() != nil {
		return nil, errors.Wrapf(rows.Err(), "error occured during rows iteration")
	}

	return productTypes, nil
}

func (pts *SqlProductTypeStore) ProductTypesByProductIDs(productIDs []string) ([]*product_and_discount.ProductType, error) {
	var productTypes []*product_and_discount.ProductType
	_, err := pts.GetReplica().Select(
		&productTypes,
		`SELECT * FROM `+store.ProductTypeTableName+` AS PT 
		INNER JOIN `+store.ProductTableName+` AS P ON (
			PT.Id = P.ProductTypeID
		) WHERE (
			P.Id IN :IDs
		)`,
		map[string]interface{}{
			"IDs": productIDs,
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find product types with given product ids")
	}

	return productTypes, nil
}

// ProductTypeByProductVariantID finds and returns 1 product type that is related to given product variant
func (pts *SqlProductTypeStore) ProductTypeByProductVariantID(variantID string) (*product_and_discount.ProductType, error) {
	query := pts.GetQueryBuilder().
		Select(pts.ModelFields()...).
		From(store.ProductTypeTableName).
		OrderBy(store.TableOrderingMap[store.ProductTypeTableName]).
		InnerJoin(store.ProductTableName + " ON (Products.ProductTypeID = ProductTypes.Id)").
		InnerJoin(store.ProductVariantTableName + " ON (Products.Id = ProductVariants.ProductID)").
		Where(squirrel.Eq{"ProductVariants.Id": variantID})

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ProductTypeByProductVariantID_ToSql")
	}

	var res product_and_discount.ProductType
	err = pts.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTypeTableName, "variantID="+variantID)
		}
		return nil, errors.Wrapf(err, "failed to find product type with product variant id=%s", variantID)
	}

	return &res, nil
}
