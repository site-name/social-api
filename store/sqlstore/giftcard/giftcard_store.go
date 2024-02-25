package giftcard

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"gorm.io/gorm"
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

func (gs *SqlGiftCardStore) FilterByOption(option model_helper.GiftcardFilterOption) (int64, model.GiftcardSlice, error) {
	query := gs.
		GetQueryBuilder().
		Select(model.GiftcardTableName + ".*").
		From(model.GiftcardTableName).
		Where(option.Conditions)

	if option.OrderID != nil {
		query = query.
			InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.GiftCardID = %[2]s.Id", model.OrderGiftCardTableName, model.GiftcardTableName)).
			Where(option.OrderID)
	}
	if option.CheckoutToken != nil {
		subSelect := gs.GetQueryBuilder(squirrel.Question).
			Select("GiftcardID").
			From(model.GiftcardCheckoutTableName).
			Where(option.CheckoutToken)

		query = query.Where(squirrel.Expr("GiftCards.Id IN ?", subSelect))
	}
	if option.SelectForUpdate && option.Transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}
	if option.Distinct {
		query = query.Distinct()
	}

	// those annotations are used for pagination sorting
	if option.AnnotateRelatedProductName || option.AnnotateRelatedProductSlug {
		query = query.InnerJoin(model.ProductTableName + " ON Products.Id = GiftCards.ProductID")

		if option.AnnotateRelatedProductName {
			query = query.Column(`Products.Name AS "GiftCards.RelatedProductName"`)
		}
		if option.AnnotateRelatedProductSlug {
			query = query.Column(`Products.Slug AS "GiftCards.RelatedProductSlug"`)
		}
	}

	if option.AnnotateUsedByFirstName || option.AnnotateUsedByLastName {
		query = query.InnerJoin(model.UserTableName + " ON GiftCards.UsedByID = Users.Id")

		if option.AnnotateUsedByFirstName {
			query = query.Column(`Users.FirstName AS "GiftCards.RelatedUsedByFirstName"`)
		}
		if option.AnnotateUsedByLastName {
			query = query.Column(`Users.LastName AS "GiftCards.RelatedUsedByLastName"`)
		}
	}

	var totalCount int64
	if option.CountTotal {
		countQuery, args, err := gs.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOption_CountTotal_ToSql")
		}

		err = gs.GetReplica().Raw(countQuery, args...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total giftcards by options")
		}
	}

	// check pagination
	option.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "query_toSql")
	}

	runner := gs.GetReplica()
	if option.Transaction != nil {
		runner = option.Transaction
	}

	rows, err := runner.Raw(queryString, args...).Rows()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to count total number of giftcards by options")
	}
	defer rows.Close()

	var res model.Giftcards
	for rows.Next() {
		var gc model.Giftcard
		var scanFields = gs.ScanFields(&gc)

		if option.AnnotateRelatedProductName {
			scanFields = append(scanFields, &gc.RelatedProductName)
		}
		if option.AnnotateRelatedProductSlug {
			scanFields = append(scanFields, &gc.RelatedProductSlug)
		}
		if option.AnnotateUsedByFirstName {
			scanFields = append(scanFields, &gc.RelatedUsedByFirstName)
		}
		if option.AnnotateUsedByLastName {
			scanFields = append(scanFields, &gc.RelatedUsedByLastName)
		}

		err := rows.Scan(scanFields...)
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to scan a row of giftcard")
		}

		res = append(res, &gc)
	}

	return totalCount, res, nil
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

	// select exists product type with kind == "gift_card":
	productTypeQuery := gs.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(model.ProductTypeTableName).
		Where("ProductTypes.Kind = ?", model.GIFT_CARD).
		Where("ProductTypes.Id = Products.ProductTypeID").
		Prefix("EXISTS (").Suffix(")").
		Limit(1)

	productQuery := gs.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(model.TableNames.Products).
		Where(productTypeQuery).
		Where("Products.Id = ProductVariants.ProductID").
		Prefix("EXISTS (").Suffix(")").
		Limit(1)

	productVariantQuery := gs.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(model.TableNames.ProductVariants).
		Where(productQuery).
		Where("ProductVariants.Id = Orderlines.VariantID").
		Prefix("EXISTS (").Suffix(")").
		Limit(1)

	orderLineQuery := gs.GetQueryBuilder().
		Select("*").
		From(model.TableNames.OrderLines).
		Where(squirrel.Eq{"Orderlines.Id": orderLineIDs}).
		Where(productVariantQuery)

	queryString, args, err := orderLineQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetGiftcardLines_ToSql")
	}

	var res model.OrderLineSlice
	err = gs.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given ids")
	}

	return res, nil
}

// DeactivateOrderGiftcards update giftcards
// which have giftcard events with type == 'bought', parameters.order_id == given order id
// by setting their IsActive attribute to false
func (gs *SqlGiftCardStore) DeactivateOrderGiftcards(tx *gorm.DB, orderID string) ([]string, error) {
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
