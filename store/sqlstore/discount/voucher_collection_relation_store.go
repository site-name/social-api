package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlVoucherCollectionStore struct {
	s store.Store
}

func NewSqlVoucherCollectionStore(s store.Store) store.VoucherCollectionStore {
	return &SqlVoucherCollectionStore{s}
}

func (s *SqlVoucherCollectionStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id", "VoucherID", "CollectionID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert saves or updates given voucher collection then returns it with an error
func (vcs *SqlVoucherCollectionStore) Upsert(voucherCollection *model.VoucherCollection) (*model.VoucherCollection, error) {
	var saving bool
	if voucherCollection.Id == "" {
		voucherCollection.PreSave()
		saving = true
	}

	if err := voucherCollection.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		query := "INSERT INTO " + store.VoucherCollectionTableName + "(" + vcs.ModelFields("").Join(",") + ") VALUES (" + vcs.ModelFields(":").Join(",") + ")"
		_, err = vcs.s.GetMasterX().NamedExec(query, voucherCollection)

	} else {
		query := "UPDATE " + store.VoucherCollectionTableName + " SET " + vcs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = vcs.s.GetMasterX().NamedExec(query, voucherCollection)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if vcs.s.IsUniqueConstraintError(err, []string{"VoucherID", "CollectionID", "vouchercollections_voucherid_collectionid_key"}) {
			return nil, store.NewErrInvalidInput(store.VoucherCollectionTableName, "VoucherID/CollectionID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-collection relation with id=%s", voucherCollection.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-collection relations updated: %d instead of 1", numUpdated)
	}

	return voucherCollection, nil
}

// Get finds a voucher collection with given id, then returns it with an error
func (vcs *SqlVoucherCollectionStore) Get(voucherCollectionID string) (*model.VoucherCollection, error) {
	var res model.VoucherCollection
	err := vcs.s.GetReplicaX().Get(&res, "SELECT * FROM "+store.VoucherCollectionTableName+" WHERE Id = ?", voucherCollectionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherCollectionTableName, voucherCollectionID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-collection relation with id=%s", voucherCollectionID)
	}

	return &res, nil
}

func (s *SqlVoucherCollectionStore) FilterByOptions(options *model.VoucherCollectionFilterOptions) ([]*model.VoucherCollection, error) {
	query := s.s.GetQueryBuilder().Select("*").From(store.VoucherCollectionTableName)

	if options.VoucherID != nil {
		query = query.Where(options.VoucherID)
	}
	if options.CollectionID != nil {
		query = query.Where(options.CollectionID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.VoucherCollection
	err = s.s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher collection relations by given options")
	}

	return res, nil
}
