package checkout

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutLineStore struct {
	store.Store
}

func NewSqlCheckoutLineStore(sqlStore store.Store) store.CheckoutLineStore {
	cls := &SqlCheckoutLineStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.CheckoutLine{}, store.CheckoutLineTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
	}
	return cls
}

func (cls *SqlCheckoutLineStore) CreateIndexesIfNotExists() {
	cls.CreateIndexIfNotExists("idx_checkoutlines_checkout_id", store.CheckoutLineTableName, "CheckoutID")
	cls.CreateIndexIfNotExists("idx_checkoutlines_variant_id", store.CheckoutLineTableName, "VariantID")

	// foreign keys:
	cls.CreateForeignKeyIfNotExists(store.CheckoutLineTableName, "CheckoutID", store.CheckoutTableName, "Token", true)
	cls.CreateForeignKeyIfNotExists(store.CheckoutLineTableName, "VariantID", store.ProductVariantTableName, "Id", true)
}

func (cls *SqlCheckoutLineStore) Upsert(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, error) {
	var isSave bool

	if checkoutLine.Id == "" {
		isSave = true
	}

	checkoutLine.PreSave()
	if err := checkoutLine.IsValid(); err != nil {
		return nil, err
	}
	if isSave {
		if err := cls.GetMaster().Insert(checkoutLine); err != nil {
			return nil, errors.Wrapf(err, "failed to save checkout line with id=%s", checkoutLine.Id)
		}
	}

	if updated, err := cls.GetMaster().Update(checkoutLine); err != nil {
		return nil, errors.Wrapf(err, "failed to update checkout line with id=%s", checkoutLine.Id)
	} else if updated > 1 {
		return nil, errors.Errorf("multiple checkout lines were updated: %d, expected: 1", updated)
	}

	return checkoutLine, nil
}

func (cls *SqlCheckoutLineStore) Get(id string) (*checkout.CheckoutLine, error) {
	res, err := cls.GetReplica().Get(checkout.CheckoutLine{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to to find checkout line with id=%s", id)
	}

	return res.(*checkout.CheckoutLine), nil
}

func (cls *SqlCheckoutLineStore) CheckoutLinesByCheckoutID(checkoutID string) ([]*checkout.CheckoutLine, error) {
	var res []*checkout.CheckoutLine
	_, err := cls.GetReplica().Select(
		&res,
		`SELECT * FROM `+store.CheckoutLineTableName+` AS CkL 
		INNER JOIN `+store.CheckoutTableName+` AS Ck ON (
			Ck.Id = CkL.CheckoutID
		)
		WHERE (
			CkL.CheckoutID = :CheckoutID
		) 
		ORDER BY Ck.CreateAt ASC`,
		map[string]interface{}{"CheckoutID": checkoutID},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutLineTableName, "checkoutID="+checkoutID)
		}
		return nil, errors.Wrapf(err, "failed to get checkout lines belong to checkout with id=%s", checkoutID)
	}

	return res, nil
}

func (cls *SqlCheckoutLineStore) DeleteLines(ids []string) error {
	// validate id list
	for _, id := range ids {
		if !model.IsValidId(id) {
			return store.NewErrInvalidInput(store.CheckoutLineTableName, "ids", ids)
		}
	}

	tx, err := cls.GetMaster().Begin()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	result, err := tx.Exec("DELETE FROM "+store.CheckoutLineTableName+" WHERE Id IN :IDs", map[string]interface{}{"IDs": ids})
	if err != nil {
		return errors.Wrap(err, "failed to delete checkout lines")
	}
	if rows, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "failed to count number of checkout lines deleted")
	} else if rows != int64(len(ids)) {
		return errors.Errorf("expect %d checkout lines to be deleted but got %d", len(ids), rows)
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}

func (cls *SqlCheckoutLineStore) BulkUpdate(lines []*checkout.CheckoutLine) error {
	for _, line := range lines {
		if line == nil || line.IsValid() != nil {
			return store.NewErrInvalidInput(store.CheckoutLineTableName, "lines", "nil value")
		}
	}

	tx, err := cls.GetMaster().Begin()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, line := range lines {
		numUpdated, err := tx.Update(line)
		if err != nil {
			return errors.Wrapf(err, "failed to update checkout line with id=%s", line.Id)
		}
		if numUpdated > 1 {
			return errors.Errorf("multiple checkout lines updated: %d instead of 1", numUpdated)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}

func (cls *SqlCheckoutLineStore) BulkCreate(lines []*checkout.CheckoutLine) ([]*checkout.CheckoutLine, error) {
	for _, line := range lines {
		if line == nil {
			return nil, store.NewErrInvalidInput(store.CheckoutLineTableName, "lines", "nil value")
		}
	}

	tx, err := cls.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, line := range lines {
		line.PreSave()
		if appErr := line.IsValid(); appErr != nil {
			return nil, appErr
		}
		err = tx.Insert(line)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to save checkout line with id=%s", line.Id)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit_transaction")
	}

	return lines, nil
}
