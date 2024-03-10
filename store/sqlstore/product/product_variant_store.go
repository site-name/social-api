package product

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/model_types"
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

	err := model.ProductVariantSlice(
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
	variant, err := model.ProductVariantSlice(
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

func (vs *SqlProductVariantStore) commonQueryBuilder(option model_helper.ProductVariantFilterOptions) []qm.QueryMod {
	conds := option.Conditions
	for _, load := range option.Preloads {
		conds = append(conds, qm.Load(load))
	}
	if option.VoucherID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherProductVariants, model.VoucherProductVariantTableColumns.ProductVariantID, model.ProductVariantTableColumns.ID)),
			option.VoucherID,
		)
	}
	if option.SaleID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.SaleProductVariants, model.SaleProductVariantTableColumns.ProductVariantID, model.ProductVariantTableColumns.ID)),
			option.SaleID,
		)
	}
	if option.RelatedProductVariantChannelListingConds != nil ||
		option.ProductVariantChannelListingChannelSlug != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariantChannelListings, model.ProductVariantChannelListingTableColumns.VariantID, model.ProductVariantTableColumns.ID)),
		)

		if option.RelatedProductVariantChannelListingConds != nil {
			conds = append(conds, option.RelatedProductVariantChannelListingConds)
		}

		if option.ProductVariantChannelListingChannelSlug != nil {
			conds = append(
				conds,
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ProductVariantChannelListingTableColumns.ChannelID, model.ChannelTableColumns.ID)),
				option.ProductVariantChannelListingChannelSlug,
			)
		}
	}
	if option.WishlistID != nil ||
		option.WishlistItemID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.WishlistItemProductVariants, model.WishlistItemProductVariantTableColumns.ProductVariantID, model.ProductVariantTableColumns.ID)),
		)

		if option.WishlistItemID != nil {
			conds = append(conds, option.WishlistItemID)
		}

		if option.WishlistID != nil {
			conds = append(
				conds,
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.WishlistItems, model.WishlistItemTableColumns.ID, model.WishlistItemProductVariantTableColumns.WishlistItemID)),
				option.WishlistID,
			)
		}
	}

	return conds
}

func (vs *SqlProductVariantStore) FilterByOption(option model_helper.ProductVariantFilterOptions) (model.ProductVariantSlice, error) {
	conds := vs.commonQueryBuilder(option)
	return model.ProductVariantSlice(conds...).All(vs.GetReplica())
}

func (s *SqlProductVariantStore) FindVariantsAvailableForPurchase(variantIds []string, channelID string) (model.ProductVariantSlice, error) {
	startOfDay := util.MillisFromTime(util.StartOfDay(time.Now().UTC()))

	return model.ProductVariantSlice(
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Products, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID)),
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductChannelListings, model.ProductChannelListingTableColumns.ProductID, model.ProductTableColumns.ID)),
		model.ProductChannelListingWhere.ChannelID.EQ(channelID),
		model.ProductChannelListingWhere.AvailableForPurchase.LTE(model_types.NewNullInt64(startOfDay)),
		model.ProductVariantWhere.ID.IN(variantIds),
	).All(s.GetReplica())
}

func (s *SqlProductVariantStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = s.GetMaster()
	}

	// delete related order lines with status draft
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
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete draft order lines related to given variants")
	}

	// TODO: consider what to deleted along with variants

	// delete assigned attribute values
	// attributeValuesQuery := fmt.Sprintf(
	// 	`DELETE FROM %[1]s WHERE %[2]s IN (
	// 		SELECT
	// 			%[1]s.%[2]s
	// 		FROM
	// 			%[1]s
	// 		INNER JOIN
	// 			%[3]s ON %[3]s.%[4]s = %[1]s.%[5]s
	// 		INNER JOIN
	// 			%[6]s ON %[6]s.%[7]s = %[1]s.%[2]s
	// 		INNER JOIN
	// 			%[8]s ON %[8]s.%[9]s = %[6]s.%[10]s
	// 		WHERE (
	// 			%[3]s.%[11]s IN ?
	// 			AND %[8]s.%[12]s IN ?
	// 		)
	// 	)`,
	// 	model.AttributeValueTableName,                         // 1
	// 	model.AttributeValueColumnId,                          // 2
	// 	model.AttributeTableName,                              // 3
	// 	model.AttributeColumnId,                               // 4
	// 	model.AttributeValueColumnAttributeID,                 // 5
	// 	model.AssignedVariantAttributeValueTableName,          // 6
	// 	model.AssignedVariantAttributeValueColumnValueID,      // 7
	// 	model.AssignedVariantAttributeTableName,               // 8
	// 	model.AssignedVariantAttributeColumnId,                // 9
	// 	model.AssignedVariantAttributeValueColumnAssignmentID, // 10
	// 	model.AttributeColumnInputType,                        // 11
	// 	model.AssignedVariantAttributeColumnVariantID,         // 12
	// )

	// err = tx.Exec(attributeValuesQuery, model.TYPES_WITH_UNIQUE_VALUES, ids).Error
	// if err != nil {
	// 	return 0, errors.Wrap(err, "failed to delete attribute values related to given variants")
	// }

	// delete variants
	return model.ProductVariantSlice(
		model.ProductVariantWhere.ID.IN(ids),
	).DeleteAll(tx)
}
