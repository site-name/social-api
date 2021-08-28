package product

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantStore struct {
	store.Store
}

func NewSqlProductVariantStore(s store.Store) store.ProductVariantStore {
	pvs := &SqlProductVariantStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariant{}, store.ProductVariantTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Sku").SetMaxSize(product_and_discount.PRODUCT_VARIANT_SKU_MAX_LENGTH).SetUnique(true)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_VARIANT_NAME_MAX_LENGTH)
	}
	return pvs
}

func (ps *SqlProductVariantStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_variants_sku", store.ProductVariantTableName, "Sku")
	ps.CreateForeignKeyIfNotExists(store.ProductVariantTableName, "ProductID", store.ProductTableName, "Id", true)
}

func (ps *SqlProductVariantStore) ModelFields() []string {
	return []string{
		"ProductVariants.Id",
		"ProductVariants.Name",
		"ProductVariants.ProductID",
		"ProductVariants.Sku",
		"ProductVariants.Weight",
		"ProductVariants.WeightUnit",
		"ProductVariants.TrackInventory",
		"ProductVariants.SortOrder",
		"ProductVariants.Metadata",
		"ProductVariants.PrivateMetadata",
	}
}

func (ps *SqlProductVariantStore) Save(variant *product_and_discount.ProductVariant) (*product_and_discount.ProductVariant, error) {
	variant.PreSave()
	if err := variant.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(variant); err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Sku", "idx_productvariants_sku_unique", "productvariants_sku_key"}) {
			return nil, store.NewErrInvalidInput(store.ProductVariantTableName, "Sku", variant.Sku)
		}
		return nil, errors.Wrapf(err, "failed to save product variant with id=%s", variant.Id)
	}

	return variant, nil
}

func (ps *SqlProductVariantStore) Get(id string) (*product_and_discount.ProductVariant, error) {
	var variant product_and_discount.ProductVariant
	err := ps.GetReplica().SelectOne(
		&variant,
		"SELECT * FROM "+store.ProductVariantTableName+" WHERE Id = :id",
		map[string]interface{}{"id": id},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find product variant with id=%s", id)
	}

	return &variant, nil
}

// GetWeight returns either given variant's weight or the accompany product's weight or product type of accompany product's weight
func (ps *SqlProductVariantStore) GetWeight(productVariantID string) (*measurement.Weight, error) {
	rowScanner := ps.GetQueryBuilder().
		Select(
			"ProductVariants.Weight",
			"ProductVariants.WeightUnit",
			"Products.Weight",
			"Products.WeightUnit",
			"ProductTypes.Weight",
			"ProductTypes.WeightUnit",
		).
		From(store.ProductVariantTableName).
		InnerJoin(store.ProductTableName + " ON Products.Id = ProductVariants.ProductID").
		InnerJoin(store.ProductTypeTableName + " ON ProductTypes.Id = Products.ProductTypeID").
		Where(squirrel.Eq{"ProductVariants.Id": productVariantID}).
		RunWith(ps.GetReplica()).
		QueryRow()

	var (
		variantWeight     measurement.Weight
		productWeight     measurement.Weight
		productTypeWeight measurement.Weight
	)
	err := rowScanner.Scan(
		&variantWeight.Amount,
		&variantWeight.Unit,
		&productWeight.Amount,
		&productWeight.Unit,
		&productTypeWeight.Amount,
		&productTypeWeight.Unit,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantTableName, productVariantID)
		}
		return nil, errors.Wrapf(err, "failed to scan result for productVariantId=%s", productVariantID)
	}

	if variantWeight.Amount != nil {
		return &variantWeight, nil
	}
	if productWeight.Amount != nil {
		return &productWeight, nil
	}
	if productTypeWeight.Amount != nil {
		return &productTypeWeight, nil
	}

	return nil, errors.Errorf("weight for product variant with id=%s is not set", productVariantID)
}

// GetByOrderLineID finds and returns a product variant by given orderLineID
func (vs *SqlProductVariantStore) GetByOrderLineID(orderLineID string) (*product_and_discount.ProductVariant, error) {
	var res product_and_discount.ProductVariant
	err := vs.GetReplica().SelectOne(
		&res,
		`SELECT `+strings.Join(vs.ModelFields(), ", ")+`
		FROM `+store.ProductVariantTableName+`
		INNER JOIN `+store.OrderLineTableName+` ON (
			ProductVariants.Id = Orderlines.VariantID
		)
		WHERE Orderlines.Id = :OrderLineID`,
		map[string]interface{}{
			"OrderLineID": orderLineID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantTableName, "orderLineID="+orderLineID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant with order line id=%s", orderLineID)
	}

	return &res, nil
}

// FilterByOption finds and returns product variants based on given option
func (vs *SqlProductVariantStore) FilterByOption(option *product_and_discount.ProductVariantFilterOption) ([]*product_and_discount.ProductVariant, error) {
	query := vs.GetQueryBuilder().
		Select(vs.ModelFields()...).
		From(store.ProductVariantTableName).
		OrderBy(store.TableOrderingMap[store.ProductVariantTableName])

	// parse option
	if option.Distinct {
		query = query.Distinct()
	}
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("ProductVariants.Id"))
	}
	if option.Name != nil {
		query = query.Where(option.Name.ToSquirrel("ProductVariants.Name"))
	}

	var joinedProductVariantChannelListingTable bool
	if option.ProductVariantChannelListingPriceAmount != nil {
		query = query.
			InnerJoin(store.ProductVariantChannelListingTableName + " ON (ProductVariantChannelListings.VariantID = ProductVariants.Id)").
			Where(option.ProductVariantChannelListingPriceAmount.ToSquirrel("ProductVariantChannelListings.PriceAmount"))
		joinedProductVariantChannelListingTable = true // indicate that already joined
	}
	if option.ProductVariantChannelListingChannelSlug != nil {
		if !joinedProductVariantChannelListingTable { // check if joined or not
			query = query.InnerJoin(store.ProductVariantChannelListingTableName + " ON (ProductVariantChannelListings.VariantID = ProductVariants.Id)")
		}
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = ProductVariantChannelListings.ChannelID)").
			Where(option.ProductVariantChannelListingChannelSlug.ToSquirrel("Channels.Slug"))
	}

	if option.WishlistID != nil {
		query = query.
			InnerJoin(store.WishlistProductVariantTableName + " ON (WishlistItemProductVariants.ProductVariantID = ProductVariants.Id)").
			InnerJoin(store.WishlistItemTableName + " ON (WishlistItemProductVariants.WishlistItemID = WishlistItems.Id)").
			Where(option.WishlistID.ToSquirrel("WishlistItems.WishlistID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.ProductVariant
	_, err = vs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variants by given option")
	}

	return res, nil
}
