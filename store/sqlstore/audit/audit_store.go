package audit

import (
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlAuditStore struct {
	store.Store
}

func NewSqlAuditStore(sqlStore store.Store) store.AuditStore {
	return &SqlAuditStore{sqlStore}
}

func (s *SqlAuditStore) Save(audit model.Audit) error {
	err := audit.Insert(s.GetMaster(), boil.Infer())
	if err != nil {
		return errors.Wrapf(err, "failed to save Audit with userId=%s and action=%s", audit.UserID, audit.Action)
	}
	return nil
}

func (s *SqlAuditStore) Get(userId string, offset int, limit int) (model.AuditSlice, error) {
	if limit > 1000 || limit <= 0 {
		return nil, store.NewErrOutOfBounds(limit)
	}

	return model.Audits(model.AuditWhere.UserID.EQ(userId), qm.Limit(limit), qm.Offset(offset)).All(s.GetReplica())
}

func (s *SqlAuditStore) PermanentDeleteByUser(userId string) error {
	_, err := model.Audits(model.AuditWhere.UserID.EQ(userId)).DeleteAll(s.GetMaster())
	return err
}
