package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
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
		"ProductVariants.IsPreOrder",
		"ProductVariants.PreorderEndDate",
		"ProductVariants.PreOrderGlobalThreshold",
		"ProductVariants.SortOrder",
		"ProductVariants.Metadata",
		"ProductVariants.PrivateMetadata",
	}
}

func (ps *SqlProductVariantStore) ScanFields(variant product_and_discount.ProductVariant) []interface{} {
	return []interface{}{
		&variant.Id,
		&variant.Name,
		&variant.ProductID,
		&variant.Sku,
		&variant.Weight,
		&variant.WeightUnit,
		&variant.TrackInventory,
		&variant.IsPreOrder,
		&variant.PreorderEndDate,
		&variant.PreOrderGlobalThreshold,
		&variant.SortOrder,
		&variant.Metadata,
		&variant.PrivateMetadata,
	}
}

func (ps *SqlProductVariantStore) Save(transaction *gorp.Transaction, variant *product_and_discount.ProductVariant) (*product_and_discount.ProductVariant, error) {
	var upsertor store.Upsertor = ps.GetMaster()
	if transaction != nil {
		upsertor = transaction
	}

	variant.PreSave()
	if err := variant.IsValid(); err != nil {
		return nil, err
	}

	if err := upsertor.Insert(variant); err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Sku", "idx_productvariants_sku_unique", "productvariants_sku_key"}) {
			return nil, store.NewErrInvalidInput(store.ProductVariantTableName, "Sku", variant.Sku)
		}
		return nil, errors.Wrapf(err, "failed to save product variant with id=%s", variant.Id)
	}

	return variant, nil
}

// Update updates given product variant and returns it
func (ps *SqlProductVariantStore) Update(transaction *gorp.Transaction, variant *product_and_discount.ProductVariant) (*product_and_discount.ProductVariant, error) {
	variant.PreUpdate()
	if err := variant.IsValid(); err != nil {
		return nil, err
	}

	var selectUpsertor store.SelectUpsertor = ps.GetMaster()
	if transaction != nil {
		selectUpsertor = transaction
	}

	err := selectUpsertor.SelectOne(&product_and_discount.ProductVariant{}, "SELECT * FROM "+store.ProductVariantTableName+" WHERE Id = :ID", map[string]interface{}{"ID": variant.Id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantTableName, variant.Id)
		}
		return nil, errors.Wrapf(err, "failed to check if a product variant with id=%s does exist", variant.Id)
	}

	numUpdated, err := selectUpsertor.Update(variant)
	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Sku", "idx_productvariants_sku_unique", "productvariants_sku_key"}) {
			return nil, store.NewErrInvalidInput(store.ProductVariantTableName, "Sku", variant.Sku)
		}
		return nil, errors.Wrapf(err, "failed to update product variant with id=%s", variant.Id)
	}
	if numUpdated != 1 {
		return nil, errors.Errorf("%d product variant(s) were/was updated instead of 1", numUpdated)
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

	query, args, _ := vs.GetQueryBuilder().
		Select(vs.ModelFields()...).
		From(store.ProductVariantTableName).
		InnerJoin(store.OrderLineTableName + " ON (ProductVariants.Id = Orderlines.VariantID)").
		Where(squirrel.Eq{"Orderlines.Id": orderLineID}).
		ToSql()

	err := vs.GetReplica().SelectOne(&res, query, args...)
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

	selectFields := vs.ModelFields()
	if option.SelectRelatedDigitalContent {
		selectFields = append(selectFields, vs.DigitalContent().ModelFields()...)
	}

	query := vs.GetQueryBuilder().
		Select(selectFields...).
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

	var joined_WishlistProductVariantTableName_table bool

	if option.WishlistItemID != nil {
		query = query.
			InnerJoin(store.WishlistProductVariantTableName + " ON (WishlistItemProductVariants.ProductVariantID = ProductVariants.Id)").
			Where(option.WishlistItemID.ToSquirrel("WishlistItemProductVariants.WishlistItemID"))

		joined_WishlistProductVariantTableName_table = true // indicate joined `WishlistProductVariantTableName`
	}

	if option.WishlistID != nil {
		if !joined_WishlistProductVariantTableName_table {
			query = query.InnerJoin(store.WishlistProductVariantTableName + " ON (WishlistItemProductVariants.ProductVariantID = ProductVariants.Id)")
		}
		query = query.
			InnerJoin(store.WishlistItemTableName + " ON (WishlistItemProductVariants.WishlistItemID = WishlistItems.Id)").
			Where(option.WishlistID.ToSquirrel("WishlistItems.WishlistID"))
	}

	if option.SelectRelatedDigitalContent {
		query = query.InnerJoin(store.ProductDigitalContentTableName + " ON (ProductVariants.Id = DigitalContents.ProductVariantID)")
	}

	rows, err := query.RunWith(vs.GetReplica()).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variants by options")
	}

	var (
		res            []*product_and_discount.ProductVariant
		variant        product_and_discount.ProductVariant
		digitalContent product_and_discount.DigitalContent

		scanFields = vs.ScanFields(variant)
	)
	if option.SelectRelatedDigitalContent {
		scanFields = append(scanFields, vs.DigitalContent().ScanFields(digitalContent)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product")
		}

		if option.SelectRelatedDigitalContent {
			variant.DigitalContent = digitalContent.DeepCopy()
		}
		res = append(res, variant.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return res, nil
}
