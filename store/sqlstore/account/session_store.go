package account

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
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
	err := me.GetMaster().Create(session).Error
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (me *SqlSessionStore) Get(ctx context.Context, sessionIdOrToken string) (*model.Session, error) {
	var session model.Session

	err := me.DBXFromContext(ctx).First(&session, "Token = ? OR Id = ?", sessionIdOrToken, sessionIdOrToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.SessionTableName, sessionIdOrToken)
		}
		return nil, errors.Wrapf(err, "failed to find Sessions with sessionIdOrToken=%s", sessionIdOrToken)
	}

	return &session, nil
}

func (me *SqlSessionStore) GetSessions(userId string) ([]*model.Session, error) {
	var sessions []*model.Session

	if err := me.GetReplica().Order("LastActivityAt DESC").Find(&sessions, "UserId = ?", userId).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find Sessions with userId=%s", userId)
	}

	return sessions, nil
}

func (me *SqlSessionStore) GetSessionsWithActiveDeviceIds(userId string) ([]*model.Session, error) {
	var sessions []*model.Session
	err := me.GetReplica().Find(&sessions, `UserId = ? AND ExpiresAt != 0 AND ExpiresAt >= ? AND DeviceId != ''`, userId, model.GetMillis()).Error
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

	var sessions []*model.Session
	err := me.GetReplica().Find(&sessions, store.BuildSqlizer(cond)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find Sessions")
	}
	return sessions, nil
}

func (me *SqlSessionStore) UpdateExpiredNotify(sessionId string, notified bool) error {
	err := me.GetMaster().Raw("UPDATE Sessions SET ExpiredNotify = ? WHERE Id = ?", notified, sessionId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with id=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) Remove(sessionIdOrToken string) error {
	err := me.GetMaster().Raw("DELETE FROM Sessions WHERE Id = ? Or Token = ?", sessionIdOrToken, sessionIdOrToken).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete Session with sessionIdOrToken=%s", sessionIdOrToken)
	}
	return nil
}

func (me *SqlSessionStore) RemoveAllSessions() error {
	err := me.GetMaster().Raw("DELETE FROM Sessions").Error
	if err != nil {
		return errors.Wrap(err, "failed to delete all Sessions")
	}
	return nil
}

func (me *SqlSessionStore) PermanentDeleteSessionsByUser(userId string) error {
	err := me.GetMaster().Raw("DELETE FROM Sessions WHERE UserId = ?", userId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete Session with userId=%s", userId)
	}

	return nil
}

func (me *SqlSessionStore) UpdateExpiresAt(sessionId string, time int64) error {
	err := me.GetMaster().Raw("UPDATE Sessions SET ExpiresAt = ?, ExpiredNotify = false WHERE Id = ?", time, sessionId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with sessionId=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) UpdateLastActivityAt(sessionId string, time int64) error {
	err := me.GetMaster().Raw("UPDATE Sessions SET LastActivityAt = ? WHERE Id = ?", time, sessionId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update Session with id=%s", sessionId)
	}
	return nil
}

func (me *SqlSessionStore) UpdateRoles(userId, roles string) (string, error) {
	err := me.GetMaster().Raw(
		"UPDATE Sessions SET Roles = ? WHERE UserId = ?",
		roles,
		userId,
	).Error
	if err != nil {
		return "", errors.Wrapf(err, "failed to update Session with userId=%s and roles=%s", userId, roles)
	}
	return userId, nil
}

func (me *SqlSessionStore) UpdateDeviceId(id string, deviceId string, expiresAt int64) (string, error) {
	err := me.GetMaster().Raw(
		"UPDATE Sessions SET DeviceId = ?, ExpiresAt = ?, ExpiredNotify = false WHERE Id = ?",
		deviceId,
		expiresAt,
		id,
	).Error
	if err != nil {
		return "", errors.Wrapf(err, "failed to update Session with id=%s", id)
	}
	return deviceId, nil
}

func (me *SqlSessionStore) UpdateProps(session *model.Session) error {
	result := me.GetMaster().Raw("UPDATE Sessions SET Props = ? WHERE Id = ?", session.Props, session.Id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update session")
	}

	return nil
}

func (me *SqlSessionStore) AnalyticsSessionCount() (int64, error) {
	var count int64
	err := me.GetMaster().Raw("SELECT COUNT(*) FROM Sessions WHERE ExpiresAt > ?", model.GetMillis()).Scan(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to count Sessions")
	}

	return count, nil
}

func (me *SqlSessionStore) Cleanup(expiryTime int64, batchSize int64) {
	slog.Debug("Cleaning up session store.")

	var RowsAffected int64 = 1
	for RowsAffected > 0 {
		result := me.GetMaster().Raw(
			"DELETE FROM Sessions WHERE Id = any (array (SELECT Id FROM Sessions WHERE ExpiresAt != 0 AND ExpiresAt < ? LIMIT ?))",
			expiryTime, batchSize)
		if result.Error != nil {
			slog.Error("Unable to cleanup session store", slog.Err(result.Error))
			return
		}
		RowsAffected = result.RowsAffected

		time.Sleep(SessionsCleanupDelayMilliseconds * time.Second)
	}
}
