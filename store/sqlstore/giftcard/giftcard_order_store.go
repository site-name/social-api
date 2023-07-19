package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlGiftCardOrderStore struct {
	store.Store
}

func NewSqlGiftCardOrderStore(s store.Store) store.GiftCardOrderStore {
	return &SqlGiftCardOrderStore{s}
}

func (s *SqlGiftCardOrderStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{"Id", "GiftCardID", "OrderID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (gs *SqlGiftCardOrderStore) Save(giftCardOrder *model.OrderGiftCard) (*model.OrderGiftCard, error) {
	giftCardOrder.PreSave()
	if err := giftCardOrder.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.OrderGiftCardTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"
	if _, err := gs.GetMaster().NamedExec(query, giftCardOrder); err != nil {
		if gs.IsUniqueConstraintError(err, []string{"GiftCardID", "OrderID", "ordergiftcards_giftcardid_orderid_key"}) {
			return nil, store.NewErrInvalidInput(model.OrderGiftCardTableName, "GiftCardID/OrderID", giftCardOrder.GiftCardID+"/"+giftCardOrder.OrderID)
		}
		return nil, errors.Wrapf(err, "failed to save giftcard-order relation with id=%s", giftCardOrder.Id)
	}

	return giftCardOrder, nil
}

func (gs *SqlGiftCardOrderStore) Get(id string) (*model.OrderGiftCard, error) {
	var res model.OrderGiftCard
	err := gs.GetReplica().Get(&res, "SELECT * FROM "+model.OrderGiftCardTableName+" WHERE Id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderGiftCardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get order-giftcard with id=%s", id)
	}

	return &res, nil
}

// BulkUpsert upserts given order-giftcard relations and returns it
func (gs *SqlGiftCardOrderStore) BulkUpsert(transaction *gorm.DB, orderGiftcards ...*model.OrderGiftCard) ([]*model.OrderGiftCard, error) {
	var executor *gorm.DB = gs.GetMaster()
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
			query := "INSERT INTO " + model.OrderGiftCardTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"
			_, err = executor.NamedExec(query, relation)

		} else {
			query := "UPDATE " + model.OrderGiftCardTableName + " SET " + gs.
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
				return nil, store.NewErrInvalidInput(model.OrderGiftCardTableName, "GiftCardID/OrderID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert order-giftcard relation with id=%s", relation.Id)
		}

		if numUpdated != 1 {
			return nil, errors.Errorf("%d relation(s) were updated, expected 1", numUpdated)
		}
	}

	return orderGiftcards, nil
}

func (s *SqlGiftCardOrderStore) FilterByOptions(options *model.OrderGiftCardFilterOptions) ([]*model.OrderGiftCard, error) {
	query := s.GetQueryBuilder().Select("Id", "GiftCardID", "OrderID").From(model.OrderGiftCardTableName)

	if options.GiftCardID != nil {
		query = query.Where(options.GiftCardID)
	}
	if options.OrderID != nil {
		query = query.Where(options.OrderID)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOptions_ToSql")
	}

	var res []*model.OrderGiftCard
	err = s.GetReplica().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order-giftcard relations with given options")
	}

	return res, nil
}
