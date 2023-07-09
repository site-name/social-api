package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlProductVariantStore struct {
	store.Store
}

func NewSqlProductVariantStore(s store.Store) store.ProductVariantStore {
	return &SqlProductVariantStore{s}
}

func (ps *SqlProductVariantStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Name",
		"ProductID",
		"Sku",
		"Weight",
		"WeightUnit",
		"TrackInventory",
		"IsPreOrder",
		"PreorderEndDate",
		"PreOrderGlobalThreshold",
		"SortOrder",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ps *SqlProductVariantStore) ScanFields(variant *model.ProductVariant) []interface{} {
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

func (ps *SqlProductVariantStore) Save(transaction store_iface.SqlxTxExecutor, variant *model.ProductVariant) (*model.ProductVariant, error) {
	var executor store_iface.SqlxExecutor = ps.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	variant.PreSave()
	if err := variant.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ProductVariantTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
	if _, err := executor.NamedExec(query, variant); err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Sku", "idx_productvariants_sku_unique", "productvariants_sku_key"}) {
			return nil, store.NewErrInvalidInput(store.ProductVariantTableName, "Sku", variant.Sku)
		}
		return nil, errors.Wrapf(err, "failed to save product variant with id=%s", variant.Id)
	}

	return variant, nil
}

// Update updates given product variant and returns it
func (ps *SqlProductVariantStore) Update(transaction store_iface.SqlxTxExecutor, variant *model.ProductVariant) (*model.ProductVariant, error) {
	variant.PreUpdate()
	if err := variant.IsValid(); err != nil {
		return nil, err
	}

	var executor store_iface.SqlxExecutor = ps.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	query := "UPDATE " + store.ProductVariantTableName + " SET " + ps.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	result, err := executor.NamedExec(query, variant)
	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Sku", "idx_productvariants_sku_unique", "productvariants_sku_key"}) {
			return nil, store.NewErrInvalidInput(store.ProductVariantTableName, "Sku", variant.Sku)
		}
		return nil, errors.Wrapf(err, "failed to update product variant with id=%s", variant.Id)
	}
	if numUpdated, _ := result.RowsAffected(); numUpdated > 1 {
		return nil, errors.Errorf("%d product variant(s) were/was updated instead of 1", numUpdated)
	}

	return variant, nil
}

func (ps *SqlProductVariantStore) Get(id string) (*model.ProductVariant, error) {
	var variant model.ProductVariant
	err := ps.GetReplicaX().Get(
		&variant,
		"SELECT * FROM "+store.ProductVariantTableName+" WHERE Id = ?",
		id,
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
	queryString, args, err := ps.GetQueryBuilder().
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
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "GetWeight_ToSql")
	}

	var (
		variantWeightAmount *float32
		variantWeightUnit   measurement.WeightUnit

		productWeightAmount *float32
		productWeightUnit   measurement.WeightUnit

		productTypeWeightAmount *float32
		productTypeWeightUnit   measurement.WeightUnit
	)
	err = ps.
		GetReplicaX().
		QueryRowX(queryString, args...).
		Scan(
			&variantWeightAmount,
			&variantWeightUnit,
			&productWeightAmount,
			&productWeightUnit,
			&productTypeWeightAmount,
			&productTypeWeightUnit,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantTableName, productVariantID)
		}
		return nil, errors.Wrapf(err, "failed to scan result for productVariantId=%s", productVariantID)
	}

	if variantWeightAmount != nil && variantWeightUnit != "" {
		return &measurement.Weight{Amount: *variantWeightAmount, Unit: variantWeightUnit}, nil
	}
	if productWeightAmount != nil && productWeightUnit != "" {
		return &measurement.Weight{Amount: *productTypeWeightAmount, Unit: productTypeWeightUnit}, nil
	}
	if productTypeWeightAmount != nil && productTypeWeightUnit != "" {
		return &measurement.Weight{Amount: *productTypeWeightAmount, Unit: productTypeWeightUnit}, nil
	}

	return nil, errors.Errorf("weight for product variant with id=%s is not set", productVariantID)
}

// GetByOrderLineID finds and returns a product variant by given orderLineID
func (vs *SqlProductVariantStore) GetByOrderLineID(orderLineID string) (*model.ProductVariant, error) {
	var res model.ProductVariant

	query := "SELECT " +
		vs.ModelFields(store.ProductVariantTableName+".").Join(",") +
		" FROM " + store.ProductVariantTableName +
		" INNER JOIN " + store.OrderLineTableName +
		" ON ProductVariants.Id = Orderlines.VariantID WHERE Orderlines.Id = ?"

	err := vs.GetReplicaX().Get(&res, query, orderLineID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantTableName, "orderLineID="+orderLineID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant with order line id=%s", orderLineID)
	}

	return &res, nil
}

// FilterByOption finds and returns product variants based on given option
func (vs *SqlProductVariantStore) FilterByOption(option *model.ProductVariantFilterOption) ([]*model.ProductVariant, error) {
	selectFields := vs.ModelFields(store.ProductVariantTableName + ".")
	if option.SelectRelatedDigitalContent {
		selectFields = append(selectFields, vs.DigitalContent().ModelFields(store.DigitalContentTableName+".")...)
	}

	query := vs.GetQueryBuilder().
		Select(selectFields...).
		From(store.ProductVariantTableName)

	// parse option
	for _, opt := range []squirrel.Sqlizer{
		option.Id,
		option.Name,
		option.ProductID,
		option.ProductVariantChannelListingPriceAmount,
		option.ProductVariantChannelListingChannelID,
		option.ProductVariantChannelListingChannelSlug,
		option.WishlistItemID,
		option.WishlistID,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if option.Distinct {
		query = query.Distinct()
	}

	// catch all inner join
	if option.ProductVariantChannelListingPriceAmount != nil ||
		option.ProductVariantChannelListingChannelID != nil ||
		option.ProductVariantChannelListingChannelSlug != nil {
		query = query.InnerJoin(store.ProductVariantChannelListingTableName + " ON ProductVariantChannelListings.VariantID = ProductVariants.Id")

		if option.ProductVariantChannelListingChannelSlug != nil {
			query = query.InnerJoin(store.ChannelTableName + " ON Channels.Id = ProductVariantChannelListings.ChannelID")
		}
	}

	if option.WishlistItemID != nil ||
		option.WishlistID != nil {
		query = query.InnerJoin(store.WishlistItemProductVariantTableName + " ON WishlistItemProductVariants.ProductVariantID = ProductVariants.Id")

		if option.WishlistID != nil {
			query = query.InnerJoin(store.WishlistItemTableName + " ON WishlistItemProductVariants.WishlistItemID = WishlistItems.Id")
		}
	}

	if option.SelectRelatedDigitalContent {
		query = query.InnerJoin(store.DigitalContentTableName + " ON ProductVariants.Id = DigitalContents.ProductVariantID")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := vs.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variants by options")
	}
	defer rows.Close()

	var res model.ProductVariants

	for rows.Next() {
		var (
			variant        model.ProductVariant
			digitalContent model.DigitalContent
			scanFields     = vs.ScanFields(&variant)
		)
		if option.SelectRelatedDigitalContent {
			scanFields = append(scanFields, vs.DigitalContent().ScanFields(&digitalContent)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product")
		}

		if option.SelectRelatedDigitalContent {
			variant.SetDigitalContent(&digitalContent)
		}
		res = append(res, &variant)
	}

	return res, nil
}
