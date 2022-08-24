package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlGiftCardOrderStore struct {
	store.Store
}

func NewSqlGiftCardOrderStore(s store.Store) store.GiftCardOrderStore {
	return &SqlGiftCardOrderStore{s}
}

func (s *SqlGiftCardOrderStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{"Id", "GiftCardID", "OrderID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (gs *SqlGiftCardOrderStore) Save(giftCardOrder *giftcard.OrderGiftCard) (*giftcard.OrderGiftCard, error) {
	giftCardOrder.PreSave()
	if err := giftCardOrder.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.OrderGiftCardTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"
	if _, err := gs.GetMasterX().NamedExec(query, giftCardOrder); err != nil {
		if gs.IsUniqueConstraintError(err, []string{"GiftCardID", "OrderID", "ordergiftcards_giftcardid_orderid_key"}) {
			return nil, store.NewErrInvalidInput(store.OrderGiftCardTableName, "GiftCardID/OrderID", giftCardOrder.GiftCardID+"/"+giftCardOrder.OrderID)
		}
		return nil, errors.Wrapf(err, "failed to save giftcard-order relation with id=%s", giftCardOrder.Id)
	}

	return giftCardOrder, nil
}

func (gs *SqlGiftCardOrderStore) Get(id string) (*giftcard.OrderGiftCard, error) {
	var res giftcard.OrderGiftCard
	err := gs.GetReplicaX().Get(&res, "SELECT * FROM "+store.OrderGiftCardTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderGiftCardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get order-giftcard with id=%s", id)
	}

	return &res, nil
}

// BulkUpsert upserts given order-giftcard relations and returns it
func (gs *SqlGiftCardOrderStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, orderGiftcards ...*giftcard.OrderGiftCard) ([]*giftcard.OrderGiftCard, error) {
	var executor store_iface.SqlxExecutor = gs.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	for _, relation := range orderGiftcards {
		isSaving := false

		if !model.IsValidId(relation.Id) {
			relation.PreSave()
			isSaving = true
		}

		if err := relation.IsValid(); err != nil {
			return nil, err
		}

		var (
			err        error
			numUpdated int64
		)
		if isSaving {
			query := "INSERT INTO " + store.OrderGiftCardTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"
			_, err = executor.NamedExec(query, relation)

		} else {
			query := "UPDATE " + store.OrderGiftCardTableName + " SET " + gs.
				ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=?"

			var result sql.Result
			result, err = executor.NamedExec(query, relation)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if gs.IsUniqueConstraintError(err, []string{"GiftCardID", "OrderID", "ordergiftcards_giftcardid_orderid_key"}) {
				return nil, store.NewErrInvalidInput(store.OrderGiftCardTableName, "GiftCardID/OrderID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert order-giftcard relation with id=%s", relation.Id)
		}

		if numUpdated != 1 {
			return nil, errors.Errorf("%d relation(s) were updated, expected 1", numUpdated)
		}
	}

	return orderGiftcards, nil
}
