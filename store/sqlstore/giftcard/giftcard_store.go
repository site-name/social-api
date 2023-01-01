package giftcard

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlGiftCardStore struct {
	store.Store
}

func NewSqlGiftCardStore(sqlStore store.Store) store.GiftCardStore {
	return &SqlGiftCardStore{sqlStore}
}

func (s *SqlGiftCardStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
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
func (gcs *SqlGiftCardStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, giftCards ...*model.GiftCard) ([]*model.GiftCard, error) {
	var executor store_iface.SqlxExecutor = gcs.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	for _, giftCard := range giftCards {
		saving := false

		if giftCard.Id == "" {
			giftCard.PreSave()
			saving = true
		} else {
			giftCard.PreUpdate()
		}

		if err := giftCard.IsValid(); err != nil {
			return nil, err
		}

		var (
			err error

			numUpdated int64
		)
		if saving {
			query := "INSERT INTO " + store.GiftcardTableName + "(" + gcs.ModelFields("").Join(",") + ") VALUES (" + gcs.ModelFields(":").Join(",") + ")"
			_, err = executor.NamedExec(query, giftCard)

		} else {
			var oldGiftcard model.GiftCard
			err = executor.Get(&oldGiftcard, "SELECT * FROM "+store.GiftcardTableName+" WHERE Id = ?", giftCard.Id)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(store.GiftcardTableName, giftCard.Id)
				}
				return nil, err
			}

			giftCard.CreateAt = oldGiftcard.CreateAt
			giftCard.Code = oldGiftcard.Code

			query := "UPDATE " + store.GiftcardTableName + " SET " + gcs.
				ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"

			var result sql.Result
			result, err = executor.NamedExec(query, giftCard)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if gcs.IsUniqueConstraintError(err, []string{"Code", "giftcards_code_key", "idx_giftcards_code_unique"}) {
				return nil, store.NewErrInvalidInput(store.GiftcardTableName, "Code", giftCard.Code)
			}
			return nil, errors.Wrapf(err, "failed to upsert giftcard with id=%s", giftCard.Id)
		}

		if numUpdated != 1 {
			return nil, errors.Errorf("%d giftcard(s) were updated instead of 1", numUpdated)
		}
	}

	return giftCards, nil
}

func (gcs *SqlGiftCardStore) GetById(id string) (*model.GiftCard, error) {
	var res model.GiftCard
	if err := gcs.GetReplicaX().Get(&res, "SELECT * FROM "+store.GiftcardTableName+" WHERE Id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard with id=%s", id)
	} else {
		return &res, nil
	}
}

// FilterByOption finds giftcards wth option
func (gs *SqlGiftCardStore) FilterByOption(transaction store_iface.SqlxTxExecutor, option *model.GiftCardFilterOption) ([]*model.GiftCard, error) {
	var selector store_iface.SqlxExecutor = gs.GetReplicaX()
	if transaction != nil {
		selector = transaction
	}

	query := gs.
		GetQueryBuilder().
		Select("*").
		From(store.GiftcardTableName)

	if option.OrderBy != "" {
		query = query.OrderBy(option.OrderBy)
	} else {
		// defaut to code
		query = query.OrderBy(store.TableOrderingMap[store.GiftcardTableName])
	}

	// check code
	if option.Distinct {
		query = query.Distinct()
	}
	if option.CreatedByID != nil {
		query = query.Where(option.CreatedByID)
	}
	if option.UsedByID != nil {
		query = query.Where(option.UsedByID)
	}
	if option.Code != nil {
		query = query.Where(option.Code)
	}
	if option.Currency != nil {
		query = query.Where(option.Currency)
	}
	if option.ExpiryDate != nil {
		query = query.Where(option.ExpiryDate)
	}
	if option.StartDate != nil {
		query = query.Where(option.StartDate)
	}
	if option.CheckoutToken != nil {
		subSelect := gs.GetQueryBuilder().
			Select("GiftcardID").
			From(store.GiftcardCheckoutTableName).
			Where(option.CheckoutToken)

		query = query.Where(squirrel.Expr("Id IN ?", subSelect))
	}
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"IsActive": *option.IsActive})
	}
	if option.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query_toSql")
	}

	var giftcards []*model.GiftCard
	err = selector.Select(&giftcards, queryString, args...)
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
	productTypeQuery := gs.GetQueryBuilder().
		Select(`(1) AS "a"`).
		From(store.ProductTypeTableName).
		Where("ProductTypes.Kind = ?", model.GIFT_CARD).
		Where("ProductTypes.Id = Products.ProductTypeID").
		Limit(1)

	productQuery := gs.GetQueryBuilder().
		Select(`(1) AS "a"`).
		From(store.ProductTableName).
		Where(squirrel.Expr("EXISTS(?)", productTypeQuery)).
		Where("Products.Id = ProductVariants.ProductID").
		Limit(1)

	productVariantQuery := gs.GetQueryBuilder().
		Select(`(1) AS "a"`).
		From(store.ProductVariantTableName).
		Where(squirrel.Expr("EXISTS(?)", productQuery)).
		Where("ProductVariants.Id = Orderlines.VariantID").
		Limit(1)

	orderLineQuery := gs.GetQueryBuilder().
		Select("*").
		From(store.OrderLineTableName).
		Where(squirrel.Eq{"Orderlines.Id": orderLineIDs}).
		Where(squirrel.Expr("EXISTS(?)", productVariantQuery)).
		OrderBy(store.TableOrderingMap[store.OrderLineTableName])

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
		From(store.GiftcardTableName).
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
			model.BOUGHT,
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
		Update(store.GiftcardTableName).
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
