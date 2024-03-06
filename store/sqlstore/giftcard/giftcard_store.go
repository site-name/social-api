package giftcard

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlGiftCardStore struct {
	store.Store
}

func NewSqlGiftCardStore(sqlStore store.Store) store.GiftCardStore {
	return &SqlGiftCardStore{sqlStore}
}

func (gcs *SqlGiftCardStore) BulkUpsert(transaction boil.ContextTransactor, giftCards model.GiftcardSlice) (model.GiftcardSlice, error) {
	if transaction == nil {
		transaction = gcs.GetMaster()
	}

	for _, giftCard := range giftCards {
		if giftCard == nil {
			continue
		}

		isSaving := giftCard.ID == ""
		if isSaving {
			model_helper.GiftcardPreSave(giftCard)
		} else {
			model_helper.GiftcardCommonPre(giftCard)
		}

		if err := model_helper.GiftcardIsValid(*giftCard); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = giftCard.Insert(transaction, boil.Infer())
		} else {
			_, err = giftCard.Update(transaction, boil.Blacklist(model.GiftcardColumns.CreatedAt))
		}

		if err != nil {
			if gcs.IsUniqueConstraintError(err, []string{model.GiftcardColumns.Code, "giftcards_code_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.Giftcards, model.GiftcardColumns.Code, giftCard.Code)
			}
			return nil, err
		}
	}

	return giftCards, nil
}

func (gcs *SqlGiftCardStore) GetById(id string) (*model.Giftcard, error) {
	giftcard, err := model.FindGiftcard(gcs.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Giftcards, id)
		}
		return nil, err
	}

	return giftcard, nil
}

func (gs *SqlGiftCardStore) commonQueryBuilder(option model_helper.GiftcardFilterOption) []qm.QueryMod {
	conds := option.Conditions

	if option.OrderID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.OrderGiftcards, model.OrderGiftcardTableColumns.GiftcardID, model.GiftcardTableColumns.ID)),
			option.OrderID,
		)
	}
	if option.CheckoutToken != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.GiftcardCheckouts, model.GiftcardCheckoutTableColumns.GiftcardID, model.GiftcardTableColumns.ID)),
			option.CheckoutToken,
		)
	}

	var annotations = model_helper.AnnotationAggregator{}
	if option.AnnotateRelatedProductNameAndSlug {
		conds = append(conds, qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Products, model.ProductTableColumns.ID, model.GiftcardTableColumns.ProductID)))

		annotations[model_helper.GiftcardAnnotationKeys.RelatedProductName] = model.ProductTableColumns.Name
		annotations[model_helper.GiftcardAnnotationKeys.RelatedProductSlug] = model.ProductTableColumns.Slug
	}
	if option.AnnotateUsedByFirstNameAndLastNames {
		conds = append(conds, qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Users, model.UserTableColumns.ID, model.GiftcardTableColumns.UsedByID)))

		annotations[model_helper.GiftcardAnnotationKeys.RelatedUsedByFirstName] = model.UserTableColumns.FirstName
		annotations[model_helper.GiftcardAnnotationKeys.RelatedUsedBylastName] = model.UserTableColumns.LastName
	}

	return append(conds, annotations)
}

func (gs *SqlGiftCardStore) FilterByOption(option model_helper.GiftcardFilterOption) (model.GiftcardSlice, error) {
	return model.Giftcards(gs.commonQueryBuilder(option)...).All(gs.GetReplica())
}

func (gs *SqlGiftCardStore) GetGiftcardLines(orderLineIDs []string) (model.OrderLineSlice, error) {
	/*
	   -- sample query for demonstration (Produced with django)

	   SELECT
	     "*"
	   FROM
	     "order_orderline"
	   WHERE
	     (
	       "order_orderline"."id" IN (1, 2, 3)
	       AND EXISTS(
	         SELECT
	           (1) AS "a"
	         FROM
	           "product_productvariant" W0
	         WHERE
	           (
	             EXISTS(
	               SELECT
	                 (1) AS "a"
	               FROM
	                 "product_product" V0
	               WHERE
	                 (
	                   EXISTS(
	                     SELECT
	                       (1) AS "a"
	                     FROM
	                       "product_producttype" U0
	                     WHERE
	                       (
	                         U0."kind" = 'gift_card'
	                         AND U0."id" = V0."product_type_id"
	                       )
	                     LIMIT
	                       1
	                   )
	                   AND V0."id" = W0."product_id"
	                 )
	               LIMIT
	                 1
	             )
	             AND W0."id" = "order_orderline"."variant_id"
	           )
	         LIMIT
	           1
	       )
	     )
	   ORDER BY
	     "order_orderline"."id" ASC
	*/

	productQuery := gs.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(model.TableNames.Products).
		Where(squirrel.Eq{
			model.ProductTableColumns.ID: model.ProductVariantTableColumns.ProductID,
		}).
		Prefix("EXISTS (").Suffix(")").
		Limit(1)

	productVariantQuery := gs.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(model.TableNames.ProductVariants).
		Where(productQuery).
		Where(squirrel.Eq{
			model.ProductVariantTableColumns.ID: model.OrderLineTableColumns.VariantID,
		}).
		Prefix("EXISTS (").Suffix(")").
		Limit(1)

	orderLineQuery := gs.GetQueryBuilder().
		Select("*").
		From(model.TableNames.OrderLines).
		Where(squirrel.Eq{model.OrderLineTableColumns.ID: orderLineIDs}).
		Where(productVariantQuery)

	queryString, args, err := orderLineQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetGiftcardLines_ToSql")
	}

	var res model.OrderLineSlice
	err = queries.Raw(queryString, args...).Bind(context.Background(), gs.GetReplica(), &res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given ids")
	}

	return res, nil
}

// DeactivateOrderGiftcards update giftcards
// which have giftcard events with type == 'bought', parameters.order_id == given order id
// by setting their IsActive attribute to false
func (gs *SqlGiftCardStore) DeactivateOrderGiftcards(tx boil.ContextTransactor, orderID string) ([]string, error) {
	query := fmt.Sprintf(
		`UPDATE %[1]s SET
			%[2]s = false
		WHERE (
			EXISTS (
				SELECT
					(1) AS "a"
				FROM
					%[3]s
				WHERE (
					%[4]s ->> '%[5]s' = ?
					AND %[6]s = ?
					AND %[7]s = %[8]s
				)
				LIMIT 1
			)
		) RETURNING %[8]s`,
		model.TableNames.Giftcards,                 // 1
		model.GiftcardTableColumns.IsActive,        // 2
		model.TableNames.GiftcardEvents,            // 3
		model.GiftcardEventTableColumns.Parameters, // 4
		"order_id",                                 // 5
		model.GiftcardEventTableColumns.Type,       // 6
		model.GiftcardEventTableColumns.GiftcardID, // 7
		model.GiftcardColumns.ID,                   // 8
	)

	var giftcardIDs []string
	err := queries.Raw(query, orderID, model.GiftcardEventTypeBought).Bind(context.Background(), gs.GetReplica(), &giftcardIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update giftcards")
	}

	return giftcardIDs, nil
}

func (s *SqlGiftCardStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.Giftcards(model.GiftcardWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
