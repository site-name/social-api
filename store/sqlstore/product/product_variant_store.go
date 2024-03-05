package product

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlProductVariantStore struct {
	store.Store
}

func NewSqlProductVariantStore(s store.Store) store.ProductVariantStore {
	return &SqlProductVariantStore{s}
}

func (ps *SqlProductVariantStore) Upsert(tx boil.ContextTransactor, variant model.ProductVariant) (*model.ProductVariant, error) {
	if tx == nil {
		tx = ps.GetMaster()
	}

	isSaving := variant.ID == ""
	if isSaving {
		model_helper.ProductVariantPreSave(&variant)
	} else {
		model_helper.ProductVariantCommonPre(&variant)
	}

	if err := model_helper.ProductVariantIsValid(variant); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = variant.Insert(tx, boil.Infer())
	} else {
		_, err = variant.Update(tx, boil.Blacklist(
			model.ProductVariantColumns.ProductID,
		))
	}

	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{model.ProductVariantColumns.Sku, "product_variants_sku_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.ProductVariants, model.ProductVariantColumns.Sku, variant.Sku)
		}
		return nil, err
	}

	return &variant, nil
}

func (ps *SqlProductVariantStore) Get(id string) (*model.ProductVariant, error) {
	variant, err := model.FindProductVariant(ps.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ProductVariants, id)
		}
		return nil, err
	}

	return variant, nil
}

func (ps *SqlProductVariantStore) GetWeight(productVariantID string) (*measurement.Weight, error) {
	var (
		variantWeightAmount, productWeightAmount *float64
		variantWeightUnit, productWeightUnit     measurement.WeightUnit
	)

	err := model.ProductVariants(
		qm.Select(
			model.ProductVariantTableColumns.Weight,
			model.ProductVariantTableColumns.WeightUnit,
			model.ProductTableColumns.Weight,
			model.ProductTableColumns.WeightUnit,
		),
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Products, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID)),
		model.ProductVariantWhere.ID.EQ(productVariantID),
	).
		Query.
		QueryRow(ps.GetReplica()).
		Scan(
			&variantWeightAmount,
			&variantWeightUnit,
			&productWeightAmount,
			&productWeightUnit,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ProductVariants, productVariantID)
		}
		return nil, errors.Wrapf(err, "failed to scan result for productVariantId=%s", productVariantID)
	}

	if variantWeightAmount != nil && variantWeightUnit != "" {
		return &measurement.Weight{Amount: *variantWeightAmount, Unit: variantWeightUnit}, nil
	}
	if productWeightAmount != nil && productWeightUnit != "" {
		return &measurement.Weight{Amount: *productWeightAmount, Unit: productWeightUnit}, nil
	}
	return nil, errors.Errorf("weight for product variant with id=%s is not set", productVariantID)
}

func (vs *SqlProductVariantStore) GetByOrderLineID(orderLineID string) (*model.ProductVariant, error) {
	variant, err := model.ProductVariants(
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.OrderLines, model.ProductVariantTableColumns.ID, model.OrderLineTableColumns.VariantID)),
		model.OrderLineWhere.ID.EQ(orderLineID),
	).One(vs.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ProductVariants, "orderLineID="+orderLineID)
		}
		return nil, err
	}

	return variant, nil
}

func (vs *SqlProductVariantStore) FilterByOption(option model_helper.ProductVariantFilterOption) ([]*model.ProductVariant, error) {
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

// func (s *SqlProductVariantStore) ToggleProductVariantRelations(
// 	tx boil.ContextTransactor,
// 	variants model.ProductVariants,
// 	medias model.ProductMedias,
// 	sales model.Sales,
// 	vouchers model.Vouchers,
// 	wishlistItems model.WishlistItems,
// 	isDelete bool,
// ) error {
// 	if tx == nil {
// 		tx = s.GetMaster()
// 	}

// 	/*
// 		Sales                  Sales                         `json:"-" gorm:"many2many:SaleProductVariants"`
// 		Vouchers               Vouchers                      `json:"-" gorm:"many2many:VoucherVariants"`
// 		ProductMedias          ProductMedias                 `json:"-" gorm:"many2many:VariantMedias"`
// 		WishlistItems          WishlistItems                 `json:"-" gorm:"many2many:WishlistItemProductVariants"`
// 	*/

// 	for _, variant := range variants {
// 		if variant == nil {
// 			continue
// 		}

// 		for assocName, relations := range map[string]any{
// 			"ProductMedias": medias,
// 			"Sales":         sales,
// 			"Vouchers":      vouchers,
// 			"WishlistItems": wishlistItems,
// 		} {
// 			if relations != nil {
// 				var err error
// 				if isDelete {
// 					err = tx.Model(variant).Association(assocName).Delete(relations)
// 				} else {
// 					err = tx.Model(variant).Association(assocName).Append(relations)
// 				}
// 				if err != nil {
// 					return errors.Wrap(err, "failed to toggle "+assocName+" product variant relations")
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

func (s *SqlProductVariantStore) FindVariantsAvailableForPurchase(variantIds []string, channelID string) (model.ProductVariantSlice, error) {
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

func (s *SqlProductVariantStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = s.GetMaster()
	}

	_, err := s.GetQueryBuilder().
		Delete(model.TableNames.OrderLines).
		Where(squirrel.Expr(
			fmt.Sprintf("%s IN ?", model.OrderLineColumns.ID),

			s.GetQueryBuilder(squirrel.Question).
				Select(model.OrderLineTableColumns.ID).
				Prefix("(").
				Suffix(")").
				From(model.TableNames.OrderLines).
				InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Orders, model.OrderLineTableColumns.OrderID, model.OrderTableColumns.ID)).
				Where(squirrel.Eq{
					model.OrderLineTableColumns.VariantID: ids,
					model.OrderTableColumns.Status:        model.OrderStatusDraft,
				}),
		)).
		RunWith(tx).
		Exec()

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
