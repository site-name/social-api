package giftcard

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
	"gorm.io/gorm"
)

type SqlGiftCardStore struct {
	store.Store
}

func NewSqlGiftCardStore(sqlStore store.Store) store.GiftCardStore {
	return &SqlGiftCardStore{sqlStore}
}

func (s *SqlGiftCardStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Code",
		"CreatedByID",
		"UsedByID",
		"CreatedByEmail",
		"UsedByEmail",
		"CreateAt",
		"StartDate",
		"ExpiryDate",
		"Tag",
		"ProductID",
		"LastUsedOn",
		"IsActive",
		"Currency",
		"InitialBalanceAmount",
		"CurrentBalanceAmount",
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

// BulkUpsert depends on given giftcards's Id properties then perform according operation
func (gcs *SqlGiftCardStore) BulkUpsert(transaction store_iface.SqlxExecutor, giftCards ...*model.GiftCard) ([]*model.GiftCard, error) {
	var executor store_iface.SqlxExecutor = gcs.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	for _, giftCard := range giftCards {
		saving := false

		if !model.IsValidId(giftCard.Id) {
			giftCard.Id = ""
			giftCard.PreSave()
			saving = true
		} else {
			giftCard.PreUpdate()
		}

		if err := giftCard.IsValid(); err != nil {
			return nil, err
		}

		var err error
		if saving {
			query := "INSERT INTO " + model.GiftcardTableName + "(" + gcs.ModelFields("").Join(",") + ") VALUES (" + gcs.ModelFields(":").Join(",") + ")"
			_, err = executor.NamedExec(query, giftCard)

		} else {
			var oldGiftcard model.GiftCard
			err = executor.Get(&oldGiftcard, "SELECT * FROM "+model.GiftcardTableName+" WHERE Id = ?", giftCard.Id)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, store.NewErrNotFound(model.GiftcardTableName, giftCard.Id)
				}
				return nil, err
			}

			giftCard.CreateAt = oldGiftcard.CreateAt
			giftCard.Code = oldGiftcard.Code

			query := "UPDATE " + model.GiftcardTableName + " SET " + gcs.
				ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"

			_, err = executor.NamedExec(query, giftCard)
		}

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
	if err := gcs.GetReplicaX().Get(&res, "SELECT * FROM "+model.GiftcardTableName+" WHERE Id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.GiftcardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard with id=%s", id)
	}
	return &res, nil
}

// FilterByOption finds giftcards wth option
func (gs *SqlGiftCardStore) FilterByOption(option *model.GiftCardFilterOption) ([]*model.GiftCard, error) {
	query := gs.
		GetQueryBuilder().
		Select(model.GiftcardTableName + ".").
		From(model.GiftcardTableName)

	if option.OrderBy != "" {
		query = query.OrderBy(option.OrderBy)
	} else {
		// defaut to code
		query = query.OrderBy("Code ASC")
	}

	for _, opt := range []squirrel.Sqlizer{
		option.Id, option.CreatedByID,
		option.UsedByID, option.Code,
		option.Currency, option.ExpiryDate,
		option.StartDate, option.IsActive,
		option.Tag, option.ProductID,
		option.CurrentBalanceAmount, option.InitialBalanceAmount,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if option.OrderID != nil {
		query = query.InnerJoin(model.OrderGiftCardTableName + " ON OrderGiftCards.GiftCardID = GiftCards.Id").
			Where(option.OrderID)
	}
	if option.CheckoutToken != nil {
		subSelect := gs.GetQueryBuilder(squirrel.Question).
			Select("GiftcardID").
			From(model.GiftcardCheckoutTableName).
			Where(option.CheckoutToken)

		query = query.Where(squirrel.Expr("GiftCards.Id IN ?", subSelect))
	}
	if option.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if option.Distinct {
		query = query.Distinct()
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query_toSql")
	}

	var giftcards []*model.GiftCard
	err = gs.GetReplicaX().Select(&giftcards, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds giftcards with code")
	}

	return giftcards, nil
}

// GetGiftcardLines returns a list of order lines
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
	err = gs.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given ids")
	}

	return res, nil
}

// DeactivateOrderGiftcards update giftcards
// which have giftcard events with type == 'bought', parameters.order_id == given order id
// by setting their IsActive attribute to false
func (gs *SqlGiftCardStore) DeactivateOrderGiftcards(orderID string) ([]string, error) {
	query, args, err := gs.GetQueryBuilder().
		Select("*").
		From(model.GiftcardTableName).
		Where(
			`EXISTS (
				SELECT
					(1) AS "a"
				FROM
					GiftcardEvents
				WHERE (
					GiftcardEvents.Parameters -> 'order_id' = ?
					AND GiftcardEvents.Type = ?
					AND GiftcardEvents.GiftcardID = GiftCards.Id
				)
				LIMIT 1
			)`,
			orderID,
			model.GIFT_CARD_EVENT_TYPE_BOUGHT,
		).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "DeactivateOrderGiftcards_ToSql")
	}

	var giftcards model.Giftcards
	err = gs.GetReplicaX().Select(&giftcards, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find giftcards with Parameters.order_id = %s", orderID)
	}

	giftcardIDs := giftcards.IDs()
	query, args, _ = gs.GetQueryBuilder().
		Update(model.GiftcardTableName).
		Set("IsActive", false).
		Where(squirrel.Eq{"Id": giftcardIDs}).
		ToSql()
	res, err := gs.GetMasterX().Exec(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update giftcards")
	}
	if num, _ := res.RowsAffected(); int(num) != len(giftcards) {
		return nil, errors.Errorf("%d giftcards updated instead of %d", num, len(giftcards))
	}

	return giftcardIDs, nil
}

func (s *SqlGiftCardStore) DeleteGiftcards(transaction store_iface.SqlxExecutor, ids []string) error {
	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	query, args, err := s.GetQueryBuilder().Delete(model.GiftcardTableName).Where(squirrel.Eq{"Id": ids}).ToSql()
	if err != nil {
		return errors.Wrap(err, "DeleteGiftcards_ToSql")
	}

	result, err := runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete gift cards")
	}
	numDeleted, _ := result.RowsAffected()
	if int(numDeleted) != len(ids) {
		return errors.Errorf("%d gift card(s) was/were deleted instead of %d", numDeleted, len(ids))
	}

	return nil
}
