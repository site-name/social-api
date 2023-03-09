package session

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

const (
	SessionsCleanupDelayMilliseconds = 100
)

type SqlSessionStore struct {
	store.Store
}

func NewSqlSessionStore(sqlStore store.Store) store.SessionStore {
	return &SqlSessionStore{sqlStore}
}

func (s *SqlSessionStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Token",
		"CreateAt",
		"ExpiresAt ",
		"LastActivityAt",
		"UserId",
		"DeviceId ",
		"Roles",
		"IsOAuth",
		"ExpiredNotify",
		"Props",
		"Local",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (me *SqlSessionStore) Save(session *model.Session) (*model.Session, error) {
	session.PreSave()
	if err := session.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO Sessions (" + me.ModelFields("").Join(",") + ") VALUES (" + me.ModelFields(":").Join(",") + ")"
	if _, err := me.GetMasterX().NamedExec(query, session); err != nil {
		return nil, errors.Wrapf(err, "failed to save Session with id=%s", session.Id)
	}

	return session, nil
}

func (me *SqlSessionStore) Get(ctx context.Context, sessionIdOrToken string) (*model.Session, error) {
	var session model.Session

	if err := me.DBXFromContext(ctx).Get(&session, "SELECT * FROM Sessions WHERE Token = ? OR Id = ? LIMIT 1", sessionIdOrToken, sessionIdOrToken); err != nil {
		return nil, errors.Wrapf(err, "failed to find Sessions with sessionIdOrToken=%s", sessionIdOrToken)
	}

	return &session, nil
}

func (me *SqlSessionStore) GetSessions(userId string) ([]*model.Session, error) {
	var sessions []*model.Session

	if err := me.GetReplicaX().Select(&sessions, "SELECT * FROM Sessions WHERE UserId = ? ORDER BY LastActivityAt DESC", userId); err != nil {
		return nil, errors.Wrapf(err, "failed to find Sessions with userId=%s", userId)
	}

	return sessions, nil
}

func (me *SqlSessionStore) GetSessionsWithActiveDeviceIds(userId string) ([]*model.Session, error) {
	var sessions []*model.Session
	err := me.GetReplicaX().
		Select(&sessions, `SELECT * FROM Sessions WHERE UserId = ? AND ExpiresAt != 0 AND ExpiresAt >= ? AND DeviceId != ''`, userId, model.GetMillis())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find Sessions with userId=%s", userId)
	}

	return sessions, nil
}

func (me *SqlSessionStore) GetSessionsExpired(thresholdMillis int64, mobileOnly bool, unnotifiedOnly bool) ([]*model.Session, error) {
	now := model.GetMillis()
	cond := squirrel.And{
		squirrel.NotEq{"ExpiresAt": 0},
		squirrel.Lt{"ExpiresAt": now},
		squirrel.Gt{"ExpiresAt": now - thresholdMillis},
	}

	if mobileOnly {
		cond = append(cond, squirrel.NotEq{"DeviceId": ""})
	}
	if unnotifiedOnly {
		cond = append(cond, squirrel.NotEq{"ExpiredNotify": true})
	}

	queryStr, args, err := me.GetQueryBuilder().Select("*").From("Sessions").Where(cond).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sessions_tosql")
	}
	var sessions []*model.Session

	err = me.GetReplicaX().Select(&sessions, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find Sessions")
	}
	return sessions, nil
}

func (me *SqlSessionStore) UpdateExpiredNotify(sessionId string, notified bool) error {
	_, err := me.GetMasterX().Exec("UPDATE Sessions SET ExpiredNotify = ? WHERE Id = ?", notified, sessionId)
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with id=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) Remove(sessionIdOrToken string) error {
	_, err := me.GetMasterX().Exec("DELETE FROM Sessions WHERE Id = ? Or Token = ?", sessionIdOrToken, sessionIdOrToken)
	if err != nil {
		return errors.Wrapf(err, "failed to delete Session with sessionIdOrToken=%s", sessionIdOrToken)
	}
	return nil
}

func (me *SqlSessionStore) RemoveAllSessions() error {
	_, err := me.GetMasterX().Exec("DELETE FROM Sessions")
	if err != nil {
		return errors.Wrap(err, "failed to delete all Sessions")
	}
	return nil
}

func (me *SqlSessionStore) PermanentDeleteSessionsByUser(userId string) error {
	_, err := me.GetMasterX().Exec("DELETE FROM Sessions WHERE UserId = ?", userId)
	if err != nil {
		return errors.Wrapf(err, "failed to delete Session with userId=%s", userId)
	}

	return nil
}

func (me *SqlSessionStore) UpdateExpiresAt(sessionId string, time int64) error {
	_, err := me.GetMasterX().Exec("UPDATE Sessions SET ExpiresAt = ?, ExpiredNotify = false WHERE Id = ?", time, sessionId)
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with sessionId=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) UpdateLastActivityAt(sessionId string, time int64) error {
	_, err := me.GetMasterX().Exec("UPDATE Sessions SET LastActivityAt = ? WHERE Id = ?", time, sessionId)
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with id=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) UpdateRoles(userId, roles string) (string, error) {
	_, err := me.GetMasterX().Exec(
		"UPDATE Sessions SET Roles = ? WHERE UserId = ?",
		roles,
		userId,
	)
	if err != nil {
		return "", errors.Wrapf(err, "failed to update Session with userId=%s and roles=%s", userId, roles)
	}
	return userId, nil
}

func (me *SqlSessionStore) UpdateDeviceId(id string, deviceId string, expiresAt int64) (string, error) {
	_, err := me.GetMasterX().Exec(
		"UPDATE Sessions SET DeviceId = ?, ExpiresAt = ?, ExpiredNotify = false WHERE Id = ?",
		deviceId,
		expiresAt,
		id,
	)
	if err != nil {
		return "", errors.Wrapf(err, "failed to update Session with id=%s", id)
	}
	return deviceId, nil
}

func (me *SqlSessionStore) UpdateProps(session *model.Session) error {
	res, err := me.GetMasterX().Exec("UPDATE Sessions SET Props = ? WHERE Id = ?", session.Props, session.Id)
	if err != nil {
		return errors.Wrap(err, "failed to update session")
	}

	n, err := res.RowsAffected()
	if err != nil || n != 1 {
		return errors.Wrap(err, "failed to update session")
	}

	return nil
}

func (me *SqlSessionStore) AnalyticsSessionCount() (int64, error) {
	var count int64
	err := me.GetMasterX().
		Get(&count, "SELECT COUNT(*) FROM Sessions WHERE ExpiresAt > ?", model.GetMillis())
	if err != nil {
		return 0, errors.Wrap(err, "failed to count Sessions")
	}

	return count, nil
}

func (me *SqlSessionStore) Cleanup(expiryTime int64, batchSize int64) {
	slog.Debug("Cleaning up session store.")

	var RowsAffected int64 = 1
	for RowsAffected > 0 {
		res, err := me.GetMasterX().
			Exec(
				"DELETE FROM Sessions WHERE Id = any (array (SELECT Id FROM Sessions WHERE ExpiresAt != 0 AND ExpiresAt < ? LIMIT ?))",
				expiryTime, batchSize)
		if err != nil {
			slog.Error("Unable to cleanup session store", slog.Err(err))
			return
		}
		RowsAffected, err = res.RowsAffected()
		if err != nil {
			slog.Error("Unable to cleanup session store.", slog.Err(err))
			return
		}

		time.Sleep(SessionsCleanupDelayMilliseconds * time.Second)
	}
}
