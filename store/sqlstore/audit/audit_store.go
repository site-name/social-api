package audit

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlAuditStore struct {
	store.Store
}

func NewSqlAuditStore(sqlStore store.Store) store.AuditStore {
	return &SqlAuditStore{sqlStore}
}

func (s SqlAuditStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
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
	audit.PreSave()
	if err := audit.IsValid(); err != nil {
		return err
	}

	query := "INSERT INTO " + store.AuditTableName + " (" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
	if _, err := s.GetMasterX().NamedExec(query, audit); err != nil {
		return errors.Wrapf(err, "failed to save Audit with userId=%s and action=%s", audit.UserId, audit.Action)
	}
	return nil
}

func (s *SqlAuditStore) Get(userId string, offset int, limit int) (model.Audits, error) {
	if limit > 1000 {
		return nil, store.NewErrOutOfBounds(limit)
	}

	query := s.GetQueryBuilder().
		Select("*").
		From("Audits").
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if userId != "" {
		query = query.Where(sq.Eq{"UserId": userId})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "audits_tosql")
	}

	var audits model.Audits
	if err := s.GetReplicaX().Select(&audits, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to get Audit list for userId=%s", userId)
	}
	return audits, nil
}

func (s *SqlAuditStore) PermanentDeleteByUser(userId string) error {
	if _, err := s.GetMasterX().Exec("DELETE FROM Audits WHERE UserId = ?", userId); err != nil {
		return errors.Wrapf(err, "failed to delete Audit with userId=%s", userId)
	}
	return nil
}
