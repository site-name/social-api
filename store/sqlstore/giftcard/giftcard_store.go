package giftcard

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlGiftCardStore struct {
	store.Store
}

func NewSqlGiftCardStore(sqlStore store.Store) store.GiftCardStore {
	return &SqlGiftCardStore{sqlStore}
}

func (gcs *SqlGiftCardStore) BulkUpsert(transaction *gorm.DB, giftCards ...*model.GiftCard) ([]*model.GiftCard, error) {
	if transaction == nil {
		transaction = gcs.GetMaster()
	}

	for _, giftCard := range giftCards {
		err := transaction.Save(giftCard).Error
		if err != nil {
			if gcs.IsUniqueConstraintError(err, []string{"Code", "giftcards_code_key", "idx_giftcards_code_unique"}) {
				return nil, store.NewErrInvalidInput(model.GiftcardTableName, "Code", giftCard.Code)
			}
			return nil, errors.Wrapf(err, "failed to upsert giftcard with id=%s", giftCard.Id)
		}
	}

	return giftCards, nil
}

func (gcs *SqlGiftCardStore) GetById(id string) (*model.GiftCard, error) {
	var res model.GiftCard
	if err := gcs.GetReplica().First(&res, "Id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.GiftcardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard with id=%s", id)
	}
	return &res, nil
}

func (gs *SqlGiftCardStore) FilterByOption(option *model.GiftCardFilterOption) (int64, []*model.GiftCard, error) {
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
		var gc model.GiftCard
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

func (gs *SqlGiftCardStore) GetGiftcardLines(orderLineIDs []string) (model.OrderLines, error) {
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
		From(model.ProductTableName).
		Where(productTypeQuery).
		Where("Products.Id = ProductVariants.ProductID").
		Prefix("EXISTS (").Suffix(")").
		Limit(1)

	productVariantQuery := gs.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(model.ProductVariantTableName).
		Where(productQuery).
		Where("ProductVariants.Id = Orderlines.VariantID").
		Prefix("EXISTS (").Suffix(")").
		Limit(1)

	orderLineQuery := gs.GetQueryBuilder().
		Select("*").
		From(model.OrderLineTableName).
		Where(squirrel.Eq{"Orderlines.Id": orderLineIDs}).
		Where(productVariantQuery)

	queryString, args, err := orderLineQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetGiftcardLines_ToSql")
	}

	var res []*model.OrderLine
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
	giftcardIDs := []string{}

	if tx == nil {
		tx = gs.GetMaster()
	}
	err := tx.Raw(`UPDATE `+model.GiftcardTableName+`SET
IsActive = false
WHERE (
	EXISTS (
		SELECT
			(1) AS "a"
		FROM
			GiftcardEvents
		WHERE (
			GiftcardEvents.Parameters ->> 'order_id' = ?
			AND GiftcardEvents.Type = ?
			AND GiftcardEvents.GiftcardID = GiftCards.Id
		)
		LIMIT 1
	)
) RETURNING Id`, orderID, model.GIFT_CARD_EVENT_TYPE_BOUGHT).Scan(&giftcardIDs).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to update giftcards")
	}

	return giftcardIDs, nil
}

func (s *SqlGiftCardStore) DeleteGiftcards(transaction *gorm.DB, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}
	err := transaction.Delete(&model.GiftCard{}, "Id IN ?", ids).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete gift cards")
	}

	return nil
}

func (s *SqlGiftCardStore) AddRelations(transaction *gorm.DB, giftcards model.Giftcards, relations any) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	var association string
	switch relations.(type) {
	case []*model.Order: // model.Orders ok also
		association = "Orders"
	case []*model.Checkout:
		association = "Checkouts"

	default:
		return store.NewErrInvalidInput("Giftcard.AddRelations", "relations", fmt.Sprintf("%T", relations))
	}

	for _, giftcard := range giftcards {
		if giftcard != nil && giftcard.Id != "" {
			err := transaction.Model(giftcard).Association(association).Append(relations)
			if err != nil {
				return errors.Wrapf(err, "failed to create giftcard-%s relations with giftcard id = %s", strings.ToLower(association), giftcard.Id)
			}
		}
	}

	return nil
}

func (s *SqlGiftCardStore) RemoveRelations(transaction *gorm.DB, giftcards model.Giftcards, relations any) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	var association string
	switch relations.(type) {
	case []*model.Order: // model.Orders ok also
		association = "Orders"
	case []*model.Checkout:
		association = "Checkouts"

	default:
		return store.NewErrInvalidInput("Giftcard.AddRelations", "relations", fmt.Sprintf("%T", relations))
	}

	for _, giftcard := range giftcards {
		if giftcard != nil && giftcard.Id != "" {
			err := transaction.Model(giftcard).Association(association).Delete(relations)
			if err != nil {
				return errors.Wrapf(err, "failed to delete giftcard-%s relations with giftcard id = %s", strings.ToLower(association), giftcard.Id)
			}
		}
	}

	return nil
}
