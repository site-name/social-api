package product

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductVariantStore struct {
	store.Store
}

func NewSqlProductVariantStore(s store.Store) store.ProductVariantStore {
	return &SqlProductVariantStore{s}
}

func (ps *SqlProductVariantStore) Save(tx *gorm.DB, variant *model.ProductVariant) (*model.ProductVariant, error) {
	if tx == nil {
		tx = ps.GetMaster()
	}

	if err := tx.Save(variant).Error; err != nil {
		if ps.IsUniqueConstraintError(err, []string{"sku", "idx_productvariants_sku_unique", "productvariants_sku_key"}) {
			return nil, store.NewErrInvalidInput(model.ProductVariantTableName, "Sku", variant.Sku)
		}
		return nil, errors.Wrapf(err, "failed to save product variant with id=%s", variant.Id)
	}

	return variant, nil
}

func (ps *SqlProductVariantStore) Get(id string) (*model.ProductVariant, error) {
	var variant model.ProductVariant
	err := ps.GetReplica().First(&variant, model.ProductVariantColumnId+" = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductVariantTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find product variant with id=%s", id)
	}

	return &variant, nil
}

// GetWeight returns either given variant's weight or the accompany product's weight or product type of accompany product's weight
func (ps *SqlProductVariantStore) GetWeight(productVariantID string) (*measurement.Weight, error) {
	queryString := fmt.Sprintf(`SELECT 
	%[1]s.%[2]s,
	%[1]s.%[3]s,

	%[4]s.%[5]s,
	%[4]s.%[6]s,

	%[7]s.%[8]s,
	%[7]s.%[9]s
	FROM
		%[1]s
	INNER JOIN %[4]s
	ON
		%[4]s.%[10]s = %[1]s.%[11]s
	INNER JOIN %[7]s
	ON
		%[7]s.%[12]s = %[4]s.%[13]s
	WHERE
		%[1]s.%[14]s = ?
	`,
		model.ProductVariantTableName,        // 1
		model.ProductVariantColumnWeight,     // 2
		model.ProductVariantColumnWeightUnit, // 3

		model.ProductTableName,    // 4
		model.ProductColumnWeight, // 5
		model.ProductColumnWeight, // 6

		model.ProductTypeTableName,        // 7
		model.ProductTypeColumnWeight,     // 8
		model.ProductTypeColumnWeightUnit, // 9

		model.ProductColumnId,               // 10
		model.ProductVariantColumnProductID, // 11

		model.ProductTypeColumnId,        // 12
		model.ProductColumnProductTypeID, // 13

		model.ProductVariantColumnId, // 14
	)

	var (
		variantWeightAmount *float32
		variantWeightUnit   measurement.WeightUnit

		productWeightAmount *float32
		productWeightUnit   measurement.WeightUnit

		productTypeWeightAmount *float32
		productTypeWeightUnit   measurement.WeightUnit
	)
	err := ps.
		GetReplica().
		Raw(queryString, productVariantID).
		Row().
		Scan(
			&variantWeightAmount,
			&variantWeightUnit,
			&productWeightAmount,
			&productWeightUnit,
			&productTypeWeightAmount,
			&productTypeWeightUnit,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound(model.ProductVariantTableName, productVariantID)
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

	query := fmt.Sprintf(
		`SELECT %[1]s.*
	FROM
		%[1]s
	INNER JOIN 
		%[2]s
	ON 
		%[1]s.%[3]s = %[2]s.%[4]s
	WHERE 
		%[2]s.%[5]s = ?`,
		model.ProductVariantTableName,  // 1
		model.OrderLineTableName,       // 2
		model.ProductVariantColumnId,   // 3
		model.OrderLineColumnVariantID, // 4
		model.OrderLineColumnId,        // 5
	)

	err := vs.GetReplica().Raw(query, orderLineID).Scan(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductVariantTableName, "orderLineID="+orderLineID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant with order line id=%s", orderLineID)
	}

	return &res, nil
}

// FilterByOption finds and returns product variants based on given option
func (vs *SqlProductVariantStore) FilterByOption(option *model.ProductVariantFilterOption) ([]*model.ProductVariant, error) {
	db := vs.GetReplica()
	if option.Distinct {
		db = db.Distinct()
	}

	for _, preload := range option.Preloads {
		db = db.Preload(preload)
	}

	conditions := squirrel.And{}
	if option.Conditions != nil {
		conditions = append(conditions, option.Conditions)
	}

	if option.RelatedProductVariantChannelListingConditions != nil ||
		option.ProductVariantChannelListingChannelSlug != nil {

		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.ProductVariantChannelListingTableName,       // 1
				model.ProductVariantTableName,                     // 2
				model.ProductVariantChannelListingColumnVariantID, // 3
				model.ProductVariantColumnId,                      // 4
			),
		)

		if option.RelatedProductVariantChannelListingConditions != nil {
			conditions = append(conditions, option.RelatedProductVariantChannelListingConditions)
		}

		if option.ProductVariantChannelListingChannelSlug != nil {
			conditions = append(conditions, option.ProductVariantChannelListingChannelSlug)
			db = db.Joins(
				fmt.Sprintf(
					"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.ChannelTableName,                            // 1
					model.ProductVariantChannelListingTableName,       // 2
					model.ChannelColumnId,                             // 3
					model.ProductVariantChannelListingColumnChannelID, // 4
				),
			)
		}
	}

	if option.WishlistID != nil ||
		option.WishlistItemID != nil {
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.WishlistItemProductVariantTableName, // 1
				model.ProductVariantTableName,             // 2
				"product_variant_id",                      // 3
				model.ProductVariantColumnId,              // 4
			),
		)

		if option.WishlistItemID != nil {
			conditions = append(conditions, option.WishlistItemID)
		}
		if option.WishlistID != nil {
			db = db.Joins(
				fmt.Sprintf(
					"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.WishlistItemTableName,               // 1
					model.WishlistItemProductVariantTableName, // 2
					"wishlist_item_id",                        // 3
					model.WishlistItemColumnId,                // 4
				),
			)
			conditions = append(conditions, option.WishlistID)
		}
	}

	if option.VoucherID != nil {
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.VoucherProductVariantTableName, // 1
				model.ProductVariantTableName,        // 2
				"product_variant_id",                 // 3
				model.ProductVariantColumnId,         // 4
			),
		)
		conditions = append(conditions, option.VoucherID)
	}

	if option.SaleID != nil {
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.SaleProductVariantTableName, // 1
				model.ProductVariantTableName,     // 2
				"product_variant_id",              // 3
				model.ProductVariantColumnId,      // 4
			),
		)
		conditions = append(conditions, option.SaleID)
	}

	args, err := store.BuildSqlizer(conditions, "Productvariant_FilterByOptions")
	if err != nil {
		return nil, err
	}
	var res model.ProductVariants
	err = db.Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variants iwth given options")
	}
	return res, nil
}

func (s *SqlProductVariantStore) ToggleProductVariantRelations(
	tx *gorm.DB,
	variants model.ProductVariants,
	medias model.ProductMedias,
	sales model.Sales,
	vouchers model.Vouchers,
	wishlistItems model.WishlistItems,
	isDelete bool,
) error {
	if tx == nil {
		tx = s.GetMaster()
	}

	/*
		Sales                  Sales                         `json:"-" gorm:"many2many:SaleProductVariants"`
		Vouchers               Vouchers                      `json:"-" gorm:"many2many:VoucherVariants"`
		ProductMedias          ProductMedias                 `json:"-" gorm:"many2many:VariantMedias"`
		WishlistItems          WishlistItems                 `json:"-" gorm:"many2many:WishlistItemProductVariants"`
	*/

	for _, variant := range variants {
		if variant == nil {
			continue
		}

		for assocName, relations := range map[string]any{
			"ProductMedias": medias,
			"Sales":         sales,
			"Vouchers":      vouchers,
			"WishlistItems": wishlistItems,
		} {
			if relations != nil {
				var err error
				if isDelete {
					err = tx.Model(variant).Association(assocName).Delete(relations)
				} else {
					err = tx.Model(variant).Association(assocName).Append(relations)
				}
				if err != nil {
					return errors.Wrap(err, "failed to toggle "+assocName+" product variant relations")
				}
			}
		}
	}

	return nil
}

func (s *SqlProductVariantStore) FindVariantsAvailableForPurchase(variantIds []string, channelID string) (model.ProductVariants, error) {
	query := fmt.Sprintf(
		`SELECT %[1]s.*
	FROM
		%[1]s
	INNER JOIN
		%[2]s
	ON
		%[1]s.%[3]s = %[2]s.%[4]s
	INNER JOIN
		%[5]s
	ON
		%[5]s.%[6]s = %[2]s.%[4]s
	WHERE
		%[5]s.%[7]s = ?      -- productChannelListing.ChannelId
		AND %[5]s.%[8]s <= ? -- productChannelListing.AvailableForPurchase
		AND %[1]s.%[9]s IN ? -- productVariant.Id`,

		model.ProductVariantTableName,                         // 1
		model.ProductTableName,                                // 2
		model.ProductVariantColumnProductID,                   // 3
		model.ProductColumnId,                                 // 4
		model.ProductChannelListingTableName,                  // 5
		model.ProductChannelListingColumnProductID,            // 6
		model.ProductChannelListingColumnChannelID,            // 7
		model.ProductChannelListingColumnAvailableForPurchase, // 8
		model.ProductVariantColumnId,                          // 9
	)

	now := util.StartOfDay(time.Now())
	var res model.ProductVariants
	err := s.GetReplica().Raw(query, channelID, now, variantIds).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find availabe for purchase product variants")
	}

	return res, nil
}

func (s *SqlProductVariantStore) Delete(tx *gorm.DB, ids []string) (int64, error) {
	if tx == nil {
		tx = s.GetMaster()
	}

	// delete draft order lines of variants
	draftOrderLinesQuery := fmt.Sprintf(
		`DELETE FROM %[1]s WHERE %[2]s IN (
			SELECT
				%[1]s.%[2]s
			FROM
				%[1]s
			INNER JOIN
				%[3]s ON %[3]s.%[4]s = %[1]s.%[5]s
			WHERE (
				%[1]s.%[6]s IN ?
				AND %[3]s.%[7]s = ?
			)
		)`,
		model.OrderLineTableName,       // 1
		model.OrderLineColumnId,        // 2
		model.OrderTableName,           // 3
		model.OrderColumnId,            // 4
		model.OrderLineColumnOrderID,   // 5
		model.OrderLineColumnVariantID, // 6
		model.OrderColumnStatus,        // 7
	)

	err := tx.Exec(draftOrderLinesQuery, ids, model.ORDER_STATUS_DRAFT).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete related draft order lines of given variants")
	}

	// delete assigned attribute values
	attributeValuesQuery := fmt.Sprintf(
		`DELETE FROM %[1]s WHERE %[2]s IN (
			SELECT
				%[1]s.%[2]s
			FROM
				%[1]s
			INNER JOIN
				%[3]s ON %[3]s.%[4]s = %[1]s.%[5]s
			INNER JOIN
				%[6]s ON %[6]s.%[7]s = %[1]s.%[2]s
			INNER JOIN
				%[8]s ON %[8]s.%[9]s = %[6]s.%[10]s
			WHERE (
				%[3]s.%[11]s IN ?
				AND %[8]s.%[12]s IN ?
			)
		)`,
		model.AttributeValueTableName,                         // 1
		model.AttributeValueColumnId,                          // 2
		model.AttributeTableName,                              // 3
		model.AttributeColumnId,                               // 4
		model.AttributeValueColumnAttributeID,                 // 5
		model.AssignedVariantAttributeValueTableName,          // 6
		model.AssignedVariantAttributeValueColumnValueID,      // 7
		model.AssignedVariantAttributeTableName,               // 8
		model.AssignedVariantAttributeColumnId,                // 9
		model.AssignedVariantAttributeValueColumnAssignmentID, // 10
		model.AttributeColumnInputType,                        // 11
		model.AssignedVariantAttributeColumnVariantID,         // 12
	)

	err = tx.Exec(attributeValuesQuery, model.TYPES_WITH_UNIQUE_VALUES, ids).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete attribute values related to given variants")
	}

	// delete variants
	tx = tx.Exec("DELETE FROM "+model.ProductVariantTableName+" WHERE Id IN ?", ids)
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "failed to delete given variants")
	}

	return tx.RowsAffected, nil
}
