package sqlstore

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

const (
	SessionsCleanupDelayMilliseconds = 100
)

type SqlSessionStore struct {
	*SqlStore
}

func newSqlSessionStore(sqlStore *SqlStore) store.SessionStore {
	us := &SqlSessionStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Session{}, "Sessions").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(26)
		table.ColMap("UserId").SetMaxSize(26)
		table.ColMap("DeviceId").SetMaxSize(512)
		table.ColMap("Roles").SetMaxSize(64)
		table.ColMap("Props").SetMaxSize(1000)
	}

	return us
}

func (me SqlSessionStore) createIndexesIfNotExists() {
	me.CreateIndexIfNotExists("idx_sessions_user_id", "Sessions", "UserId")
	me.CreateIndexIfNotExists("idx_sessions_token", "Sessions", "Token")
	me.CreateIndexIfNotExists("idx_sessions_expires_at", "Sessions", "ExpiresAt")
	me.CreateIndexIfNotExists("idx_sessions_create_at", "Sessions", "CreateAt")
	me.CreateIndexIfNotExists("idx_sessions_last_activity_at", "Sessions", "LastActivityAt")
}

// insert new session to the database
func (me *SqlSessionStore) Save(session *model.Session) (*model.Session, error) {
	if session.Id != "" {
		return nil, store.NewErrInvalidInput("Session", "id", session.Id)
	}
	session.PreSave()

	if err := me.GetMaster().Insert(session); err != nil {
		return nil, errors.Wrapf(err, "failed to save Session with id=%s", session.Id)
	}

	return session, nil
}

func (me *SqlSessionStore) Get(ctx context.Context, sessionIdOrToken string) (*model.Session, error) {
	var sessions []*model.Session

	if _, err := me.DBFromContext(ctx).Select(&sessions, "SELECT * FROM Sessions WHERE Token = :Token OR Id = :Id LIMIT 1", map[string]interface{}{"Token": sessionIdOrToken, "Id": sessionIdOrToken}); err != nil {
		return nil, errors.Wrapf(err, "failed to find Sessions with sessionIdOrToken=%s", sessionIdOrToken)
	} else if len(sessions) == 0 {
		return nil, store.NewErrNotFound("Session", fmt.Sprintf("sessionIdOrToken=%s", sessionIdOrToken))
	}

	return sessions[0], nil
}

func (me *SqlSessionStore) GetSessions(userId string) ([]*model.Session, error) {
	var sessions []*model.Session

	if _, err := me.GetReplica().Select(&sessions, "SELECT * FROM Sessions WHERE UserId = :UserId ORDER BY LastActivityAt DESC", map[string]interface{}{"UserId": userId}); err != nil {
		return nil, errors.Wrapf(err, "failed to find Sessions with userId=%s", userId)
	}

	return sessions, nil
}

func (me *SqlSessionStore) GetSessionsWithActiveDeviceIds(userId string) ([]*model.Session, error) {
	query :=
		`SELECT *
	FROM
		Sessions
	WHERE
		UserId = :UserId AND
		ExpiresAt != 0 AND
		:ExpiresAt <= ExpiresAt AND
		DeviceId != ''`

	var sessions []*model.Session
	_, err := me.GetReplica().Select(&sessions, query, map[string]interface{}{"UserId": userId, "ExpiresAt": model.GetMillis()})
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

	queryStr, args, err := me.getQueryBuilder().Select("*").From("Sessions").Where(cond).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sessions_tosql")
	}
	var sessions []*model.Session

	_, err = me.GetReplica().Select(&sessions, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find Sessions")
	}
	return sessions, nil
}

func (me *SqlSessionStore) UpdateExpiredNotify(sessionId string, notified bool) error {
	query, args, err := me.
		getQueryBuilder().
		Update("Sessions").
		Set("ExpiredNotify", notified).
		Where(squirrel.Eq{"Id": sessionId}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "sessions_tosql")
	}

	_, err = me.GetMaster().Exec(query, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with id=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) Remove(sessionIdOrToken string) error {
	_, err := me.GetMaster().Exec("DELETE FROM Sessions WHERE Id = :Id Or Token = :Token", map[string]interface{}{"Id": sessionIdOrToken, "Token": sessionIdOrToken})
	if err != nil {
		return errors.Wrapf(err, "failed to delete Session with sessionIdOrToken=%s", sessionIdOrToken)
	}
	return nil
}

func (me *SqlSessionStore) RemoveAllSessions() error {
	_, err := me.GetMaster().Exec("DELETE FROM Sessions")
	if err != nil {
		return errors.Wrap(err, "failed to delete all Sessions")
	}
	return nil
}

func (me *SqlSessionStore) PermanentDeleteSessionsByUser(userId string) error {
	_, err := me.GetMaster().Exec("DELETE FROM Sessions WHERE UserId = :UserId", map[string]interface{}{"UserId": userId})
	if err != nil {
		return errors.Wrapf(err, "failed to delete Session with userId=%s", userId)
	}

	return nil
}

func (me *SqlSessionStore) UpdateExpiresAt(sessionId string, time int64) error {
	_, err := me.GetMaster().Exec("UPDATE Sessions SET ExpiresAt = :ExpiresAt, ExpiredNotify = false WHERE Id = :Id", map[string]interface{}{"ExpiresAt": time, "Id": sessionId})
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with sessionId=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) UpdateLastActivityAt(sessionId string, time int64) error {
	_, err := me.GetMaster().Exec("UPDATE Sessions SET LastActivityAt = :LastActivityAt WHERE Id = :Id", map[string]interface{}{"LastActivityAt": time, "Id": sessionId})
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with id=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) UpdateRoles(userId, roles string) (string, error) {
	query := "UPDATE Sessions SET Roles = :Roles WHERE UserId = :UserId"

	_, err := me.GetMaster().Exec(query, map[string]interface{}{"Roles": roles, "UserId": userId})
	if err != nil {
		return "", errors.Wrapf(err, "failed to update Session with userId=%s and roles=%s", userId, roles)
	}
	return userId, nil
}

func (me *SqlSessionStore) UpdateDeviceId(id string, deviceId string, expiresAt int64) (string, error) {

	query := "UPDATE Sessions SET DeviceId = :DeviceId, ExpiresAt = :ExpiresAt, ExpiredNotify = false WHERE Id = :Id"

	_, err := me.GetMaster().Exec(query, map[string]interface{}{"DeviceId": deviceId, "Id": id, "ExpiresAt": expiresAt})
	if err != nil {
		return "", errors.Wrapf(err, "failed to update Session with id=%s", id)
	}
	return deviceId, nil
}

func (me *SqlSessionStore) UpdateProps(session *model.Session) error {
	query := "UPDATE Sessions SET Props = :Props WHERE Id = :Id"
	res, err := me.GetMaster().Exec(query, map[string]interface{}{"Props": session.Props, "Id": session.Id})
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
	query := "SELECT COUNT(*) FROM Sessions WHERE ExpiresAt > :Time"
	count, err := me.GetReplica().SelectInt(query, map[string]interface{}{"Time": model.GetMillis()})
	if err != nil {
		return 0, errors.Wrap(err, "failed to count Sessions")
	}

	return count, nil
}

func (me *SqlSessionStore) Cleanup(expiryTime int64, batchSize int64) {
	slog.Debug("Cleaning up session store.")

	query := "DELETE FROM Sessions WHERE Id = any (array (SELECT Id FROM Sessions WHERE ExpiresAt != 0 AND :ExpiresAt > ExpiresAt LIMIT :Limit))"

	var RowsAffected int64 = 1
	for RowsAffected > 0 {
		res, err := me.GetMaster().Exec(query, map[string]interface{}{"ExpiresAt": expiryTime, "Limit": batchSize})
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
