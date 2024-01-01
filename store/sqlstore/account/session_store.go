package account

import (
	"context"
	"database/sql"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

func (me *SqlSessionStore) Save(session model.Session) (*model.Session, error) {
	err := session.Insert(me.GetMaster(), boil.Infer())
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (me *SqlSessionStore) Get(ctx context.Context, sessionIdOrToken string) (*model.Session, error) {
	session, err := model.
		Sessions(model.SessionWhere.Token.EQ(sessionIdOrToken), qm.Or(model.SessionColumns.ID+" = ?", sessionIdOrToken)).
		One(me.DBXFromContext(ctx))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Sessions, sessionIdOrToken)
		}
		return nil, err
	}

	return session, nil
}

func (me *SqlSessionStore) GetSessions(userId string) (model.SessionSlice, error) {
	return model.
		Sessions(model.SessionWhere.UserID.EQ(userId), qm.OrderBy(model.SessionColumns.LastActivityAt+" DESC")).
		All(me.GetReplica())
}

func (me *SqlSessionStore) GetSessionsWithActiveDeviceIds(userId string) (model.SessionSlice, error) {
	return model.
		Sessions(
			model.SessionWhere.UserID.EQ(userId),
			model.SessionWhere.ExpiresAt.NEQ(0),
			model.SessionWhere.ExpiresAt.GTE(model_helper.GetMillis()),
			model.SessionWhere.DeviceID.NEQ(""),
		).
		All(me.GetReplica())
}

func (me *SqlSessionStore) GetSessionsExpired(thresholdMillis int64, mobileOnly bool, unnotifiedOnly bool) (model.SessionSlice, error) {
	now := model_helper.GetMillis()

	queryMods := []qm.QueryMod{
		model.SessionWhere.ExpiresAt.NEQ(0),
		model.SessionWhere.ExpiresAt.GT(now - thresholdMillis),
		model.SessionWhere.ExpiresAt.LT(now),
	}
	if mobileOnly {
		queryMods = append(queryMods, model.SessionWhere.DeviceID.NEQ(""))
	}
	if unnotifiedOnly {
		queryMods = append(queryMods, model.SessionWhere.ExpiredNotify.NEQ(true))
	}

	return model.Sessions(queryMods...).All(me.GetReplica())
}

func (me *SqlSessionStore) UpdateExpiredNotify(sessionId string, notified bool) error {
	_, err := model.
		Sessions(model.SessionWhere.ID.EQ(sessionId)).
		UpdateAll(me.GetMaster(), model.M{model.SessionColumns.ExpiredNotify: notified})
	return err
}

func (me *SqlSessionStore) Remove(sessionIdOrToken string) error {
	_, err := model.
		Sessions(model.SessionWhere.ID.EQ(sessionIdOrToken), qm.Or(model.SessionColumns.Token+" = ?", sessionIdOrToken)).
		DeleteAll(me.GetMaster())
	return err
}

func (me *SqlSessionStore) RemoveAllSessions() error {
	_, err := model.Sessions().DeleteAll(me.GetMaster())
	return err
}

func (me *SqlSessionStore) PermanentDeleteSessionsByUser(userId string) error {
	_, err := model.Sessions(model.SessionWhere.UserID.EQ(userId)).DeleteAll(me.GetMaster())
	return err
}

func (me *SqlSessionStore) UpdateExpiresAt(sessionId string, expireTime int64) error {
	_, err := model.Sessions(model.SessionWhere.ID.EQ(sessionId)).
		UpdateAll(me.GetMaster(), model.M{
			model.SessionColumns.ExpiresAt:     expireTime,
			model.SessionColumns.ExpiredNotify: false,
		})
	return err
}

func (me *SqlSessionStore) UpdateLastActivityAt(sessionId string, lastActivityAt int64) error {
	_, err := model.Sessions(model.SessionWhere.ID.EQ(sessionId)).
		UpdateAll(me.GetMaster(), model.M{
			model.SessionColumns.LastActivityAt: lastActivityAt,
		})
	return err
}

func (me *SqlSessionStore) UpdateRoles(userId, roles string) (string, error) {
	_, err := model.Sessions(model.SessionWhere.UserID.EQ(userId)).
		UpdateAll(me.GetMaster(), model.M{
			model.SessionColumns.Roles: roles,
		})
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (me *SqlSessionStore) UpdateDeviceId(id string, deviceId string, expiresAt int64) (string, error) {
	_, err := model.Sessions(model.SessionWhere.ID.EQ(id)).
		UpdateAll(me.GetMaster(), model.M{
			model.SessionColumns.DeviceID:      deviceId,
			model.SessionColumns.ExpiredNotify: false,
			model.SessionColumns.ExpiresAt:     expiresAt,
		})
	if err != nil {
		return "", err
	}
	return deviceId, err
}

func (me *SqlSessionStore) UpdateProps(session model.Session) error {
	_, err := model.Sessions(model.SessionWhere.ID.EQ(session.ID)).
		UpdateAll(me.GetMaster(), model.M{
			model.SessionColumns.Props: session.Props,
		})
	return err
}

func (me *SqlSessionStore) AnalyticsSessionCount() (int64, error) {
	return model.Sessions(model.SessionWhere.ExpiresAt.GT(model_helper.GetMillis())).Count(me.GetReplica())
}

func (me *SqlSessionStore) Cleanup(expiryTime int64, batchSize int64) {
	slog.Debug("Cleaning up session store.")

	var RowsAffected int64 = 1
	for RowsAffected > 0 {
		result, err := queries.
			Raw("DELETE FROM Sessions WHERE Id IN (SELECT Id FROM Sessions WHERE ExpiresAt != 0 AND ? > ExpiresAt LIMIT ?)", expiryTime, batchSize).
			Exec(me.GetMaster())
		if err != nil {
			slog.Error("Unable to cleanup session store", slog.Err(err))
			return
		}
		RowsAffected, _ = result.RowsAffected()

		time.Sleep(SessionsCleanupDelayMilliseconds * time.Second)
	}

}
