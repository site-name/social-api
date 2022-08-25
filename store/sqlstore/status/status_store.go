package status

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlStatusStore struct {
	store.Store
}

func NewSqlStatusStore(sqlStore store.Store) store.StatusStore {
	return &SqlStatusStore{sqlStore}
}

func (s *SqlStatusStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"UserId",
		"Status",
		"Manual",
		"LastActivityAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (s SqlStatusStore) SaveOrUpdate(status *account.Status) error {
	var (
		saveQuery   = "INSERT INTO Status (" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE Status SET " + s.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)

	if err := s.GetReplicaX().Get(&account.Status{}, "SELECT * FROM Status WHERE UserId = :UserId", map[string]interface{}{"UserId": status.UserId}); err == nil {
		if _, err := s.GetMasterX().NamedExec(updateQuery, status); err != nil {
			return errors.Wrap(err, "failed to update Status")
		}
	} else {
		if _, err := s.GetMasterX().NamedExec(saveQuery, status); err != nil {
			if !(strings.Contains(err.Error(), "for key 'PRIMARY'") && strings.Contains(err.Error(), "Duplicate entry")) {
				return errors.Wrap(err, "failed in save Status")
			}
		}
	}
	return nil
}

func (s *SqlStatusStore) Get(userId string) (*account.Status, error) {
	var status account.Status

	if err := s.GetReplicaX().Get(&status, `SELECT	* FROM Status WHERE UserId = ?`, userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Status", fmt.Sprintf("userId=%s", userId))
		}
		return nil, errors.Wrapf(err, "failed to get Status with userId=%s", userId)
	}
	return &status, nil
}

func (s *SqlStatusStore) GetByIds(userIds []string) ([]*account.Status, error) {
	rows, err := s.GetReplicaX().QueryX("SELECT UserId, Status, Manual, LastActivityAt FROM Status WHERE UserId IN ?", userIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find Statuses")
	}
	var statuses []*account.Status
	defer rows.Close()
	for rows.Next() {
		var status account.Status
		if err = rows.Scan(&status.UserId, &status.Status, &status.Manual, &status.LastActivityAt); err != nil {
			return nil, errors.Wrap(err, "unable to scan from rows")
		}
		statuses = append(statuses, &status)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed while iterating over rows")
	}

	return statuses, nil
}

func (s *SqlStatusStore) ResetAll() error {
	if _, err := s.GetMasterX().Exec("UPDATE Status SET Status = ? WHERE Manual = false", account.STATUS_OFFLINE); err != nil {
		return errors.Wrap(err, "failed to update Statuses")
	}
	return nil
}

func (s *SqlStatusStore) GetTotalActiveUsersCount() (int64, error) {

	var (
		time  = model.GetMillis() - (1000 * 60 * 60 * 24)
		count int64
	)
	err := s.GetReplicaX().Get(&count, "SELECT COUNT(UserId) FROM Status WHERE LastActivityAt > ?", time)
	if err != nil {
		return count, errors.Wrap(err, "failed to count active users")
	}
	return count, nil
}

func (s *SqlStatusStore) UpdateLastActivityAt(userId string, lastActivityAt int64) error {
	if _, err := s.GetMasterX().Exec("UPDATE Status SET LastActivityAt = ? WHERE UserId = ?", lastActivityAt, userId); err != nil {
		return errors.Wrapf(err, "failed to update last activity for userId=%s", userId)
	}

	return nil
}
