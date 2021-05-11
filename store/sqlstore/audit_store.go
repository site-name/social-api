package sqlstore

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model/audit"
	"github.com/sitename/sitename/store"
)

type SqlAuditStore struct {
	*SqlStore
}

func newSqlAuditStore(sqlStore *SqlStore) store.AuditStore {
	s := &SqlAuditStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(audit.Audit{}, "Audits").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("UserId").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Action").SetMaxSize(512)
		table.ColMap("ExtraInfo").SetMaxSize(1024)
		table.ColMap("IpAddress").SetMaxSize(64)
		table.ColMap("SessionId").SetMaxSize(UUID_MAX_LENGTH)
	}

	return s
}

func (s SqlAuditStore) createIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_audits_user_id", "Audits", "UserId")
}

func (s SqlAuditStore) Save(audit *audit.Audit) error {
	audit.PreSave()
	if err := audit.IsValid(); err != nil {
		return err
	}
	if err := s.GetMaster().Insert(audit); err != nil {
		return errors.Wrapf(err, "failed to save Audit with userId=%s and action=%s", audit.UserId, audit.Action)
	}
	return nil
}

func (s SqlAuditStore) Get(userId string, offset int, limit int) (audit.Audits, error) {
	if limit > 1000 {
		return nil, store.NewErrOutOfBounds(limit)
	}

	query := s.getQueryBuilder().
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

	var audits audit.Audits
	if _, err := s.GetReplica().Select(&audits, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to get Audit list for userId=%s", userId)
	}
	return audits, nil
}

func (s SqlAuditStore) PermanentDeleteByUser(userId string) error {
	if _, err := s.GetMaster().Exec("DELETE FROM Audits WHERE UserId = :userId",
		map[string]interface{}{"userId": userId}); err != nil {
		return errors.Wrapf(err, "failed to delete Audit with userId=%s", userId)
	}
	return nil
}
