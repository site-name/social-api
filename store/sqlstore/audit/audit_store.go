package audit

import (
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlAuditStore struct {
	store.Store
}

func NewSqlAuditStore(sqlStore store.Store) store.AuditStore {
	return &SqlAuditStore{sqlStore}
}

func (s SqlAuditStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"CreateAt",
		"UserId",
		"Action",
		"ExtraInfo",
		"IpAddress",
		"SessionId",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (s *SqlAuditStore) Save(audit *model.Audit) error {
	err := s.GetMaster().Create(audit).Error
	if err != nil {
		return errors.Wrapf(err, "failed to save Audit with userId=%s and action=%s", audit.UserId, audit.Action)
	}
	return nil
}

func (s *SqlAuditStore) Get(userId string, offset int, limit int) (model.Audits, error) {
	if limit > 1000 {
		return nil, store.NewErrOutOfBounds(limit)
	}

	query := s.GetQueryBuilder().Select("*").From(model.AuditTableName)
	if offset > 0 {
		query = query.Offset(uint64(offset))
	}
	if limit > 0 {
		query = query.Limit(uint64(limit))
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Get_ToSql")
	}

	var audits model.Audits
	err = s.GetReplica().Raw(queryStr, args...).Scan(&audits).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get Audit list for userId=%s", userId)
	}
	return audits, nil
}

func (s *SqlAuditStore) PermanentDeleteByUser(userId string) error {
	return s.GetMaster().Raw("DELETE FROM "+model.AuditTableName+" WHERE UserId = ?", userId).Error
}
