package account

import (
	"context"
	"errors"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// CreateSession try saving given session to the database. If success then add that session to cache.
func (a *ServiceAccount) CreateSession(session model.Session) (*model.Session, *model_helper.AppError) {
	session.Token = ""
	savedSession, err := a.srv.Store.Session().Save(session)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*model_helper.AppError); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("CreateSession", "app.session.save_session.app_error", nil, err.Error(), statusCode)
	}

	a.AddSessionToCache(savedSession)

	return savedSession, nil
}

func (a *ServiceAccount) ReturnSessionToPool(session *model.Session) {
	if session != nil {
		session.ID = ""
		a.sessionPool.Put(session)
	}
}

// AddSessionToCache add given session `s` to server's sessionCache, key is session's Token, expiry time as in config
func (a *ServiceAccount) AddSessionToCache(s *model.Session) {
	a.sessionCache.SetWithExpiry(s.Token, s, time.Duration(*a.srv.Config().ServiceSettings.SessionCacheInMinutes))
}

// GetSessions get session from database with UserID attribute of given `userID`
func (a *ServiceAccount) GetSessions(userID string) ([]*model.Session, *model_helper.AppError) {

	sessions, err := a.srv.Store.Session().GetSessions(userID)
	if err != nil {
		return nil, model_helper.NewAppError("GetSessions", "app.session.get_sessions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return sessions, nil
}

// RevokeAllSessions get sessions from database that has UserID of given userID, then removes them
func (a *ServiceAccount) RevokeAllSessions(userID string) *model_helper.AppError {
	sessions, err := a.srv.Store.Session().GetSessions(userID)
	if err != nil {
		return model_helper.NewAppError("RevokeAllSessions", "app.session.get_sessions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, session := range sessions {
		if session.IsOauth {
			// TODO: fixme
			// a.RevokeAccessToken(session.Token)
		} else {
			if err := a.srv.Store.Session().Remove(session.ID); err != nil {
				return model_helper.NewAppError("RevokeAllSessions", "app.session.remove.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	a.ClearSessionCacheForUser(userID)

	return nil
}

// ClearSessionCacheForUser clears all sessions that have `UserID` attribute of given `userID` in server's `sessionCache`
func (a *ServiceAccount) ClearSessionCacheForUser(userID string) {
	a.ClearSessionCacheForUserSkipClusterSend(userID)

	if a.cluster != nil {
		msg := &model_helper.ClusterMessage{
			Event:    model_helper.ClusterEventClearSessionCacheForUser,
			SendType: model_helper.ClusterSendReliable,
			Data:     []byte(userID),
		}
		a.cluster.SendClusterMessage(msg)
	}
}

func (a *ServiceAccount) ClearUserSessionCacheLocal(userID string) {
	if keys, err := a.sessionCache.Keys(); err == nil {
		var session *model.Session
		for _, key := range keys {
			if err := a.sessionCache.Get(key, &session); err == nil {
				if session.UserID == userID {
					a.sessionCache.Remove(key)
					if a.metrics != nil {
						a.metrics.IncrementMemCacheInvalidationCounterSession()
					}
				}
			}
		}
	}
}

// ClearSessionCacheForUserSkipClusterSend iterates through server's sessionCache, if it finds any session belong to given userID, removes that session.
func (a *ServiceAccount) ClearSessionCacheForUserSkipClusterSend(userID string) {
	a.ClearUserSessionCacheLocal(userID)
}

func (a *ServiceAccount) GetSession(token string) (*model.Session, *model_helper.AppError) {
	var session = a.sessionPool.Get().(*model.Session)

	var err *model_helper.AppError
	if err := a.sessionCache.Get(token, session); err == nil {
		if a.metrics != nil {
			a.metrics.IncrementMemCacheHitCounterSession()
		}
	} else {
		if a.metrics != nil {
			a.metrics.IncrementMemCacheMissCounterSession()
		}
	}

	if session.ID == "" {
		var nErr error
		session, nErr = a.srv.Store.Session().Get(sqlstore.WithMaster(context.Background()), token)
		if nErr == nil {
			if session != nil {
				if session.Token != token {
					return nil, model_helper.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "session token is different from the one in DB", http.StatusUnauthorized)
				}

				if !model_helper.SessionIsExpired(*session) {
					a.AddSessionToCache(session)
				}
			}
		} else if _, ok := nErr.(*store.ErrNotFound); !ok {
			return nil, model_helper.NewAppError("GetSession", "app.session.get.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	if session == nil || session.ID == "" {
		session, err = a.createSessionForUserAccessToken(token)
		if err != nil {
			var (
				detailedError string
				statusCode    = http.StatusUnauthorized
			)

			if err.Id != "app.user_access_token.invalid_or_missing" {
				detailedError = err.Error()
				statusCode = err.StatusCode
			} else {
				slog.Warn("Error while creating session for user access token", slog.Err(err))
			}
			return nil, model_helper.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": detailedError}, "", statusCode)
		}
	}

	if session.ID == "" || model_helper.SessionIsExpired(*session) {
		return nil, model_helper.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "session is either nil or expired", http.StatusUnauthorized)
	}

	if *a.srv.Config().ServiceSettings.SessionIdleTimeoutInMinutes > 0 &&
		!session.IsOauth &&
		!model_helper.SessionIsMobileApp(*session) &&
		session.Props[model_helper.SESSION_PROP_TYPE] != model_helper.SESSION_TYPE_USER_ACCESS_TOKEN &&
		!*a.srv.Config().ServiceSettings.ExtendSessionLengthWithActivity {

		timeout := int64(*a.srv.Config().ServiceSettings.SessionIdleTimeoutInMinutes) * 1000 * 60
		if (model_helper.GetMillis() - session.LastActivityAt) > timeout {
			// Revoking the session is an asynchronous task anyways since we are not checking
			// for the return value of the call before returning the error.
			// So moving this to a goroutine has 2 advantages:
			// 1. We are treating this as a proper asynchronous task.
			// 2. This also fixes a race condition in the web hub, where GetSession
			// gets called from (*WebConn).isMemberOfTeam and revoking a session involves
			// clearing the webconn cache, which needs the hub again.
			a.srv.Go(func() {
				err := a.RevokeSessionById(session.ID)
				if err != nil {
					slog.Warn("Error while revoking session", slog.Err(err))
				}
			})
			return nil, model_helper.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "idle timeout", http.StatusUnauthorized)
		}
	}

	return session, nil
}

// RevokeSessionById gets session with given sessionID then revokes it
func (a *ServiceAccount) RevokeSessionById(sessionID string) *model_helper.AppError {
	session, err := a.srv.Store.Session().Get(context.Background(), sessionID)
	if err != nil {
		return model_helper.NewAppError("RevokeSessionById", "app.session.get.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	return a.RevokeSession(session)
}

// func (a *ServiceAccount) ClearSessionCacheForAllUsersSkipClusterSend() {
// 	a.srv.clearSessionCacheForAllUsersSkipClusterSend()
// }

// RevokeSession removes session from database
func (a *ServiceAccount) RevokeSession(session *model.Session) *model_helper.AppError {
	if session.IsOauth {
		// TODO: fixme
		// if err := a.RevokeAccessToken(session.Token); err != nil {
		// 	return err
		// }
	} else {
		if err := a.srv.Store.Session().Remove(session.ID); err != nil {
			return model_helper.NewAppError("RevokeSession", "app.session.remove.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	a.ClearSessionCacheForUser(session.UserID)

	return nil
}

func (a *ServiceAccount) createSessionForUserAccessToken(tokenString string) (*model.Session, *model_helper.AppError) {
	token, nErr := a.srv.Store.UserAccessToken().GetByToken(tokenString)
	if nErr != nil {
		return nil, model_helper.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, nErr.Error(), http.StatusUnauthorized)
	}

	if !token.IsActive {
		return nil, model_helper.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "inactive_token", http.StatusUnauthorized)
	}

	user, nErr := a.srv.Store.User().Get(context.Background(), token.UserID)
	if nErr != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(nErr, &nfErr):
			return nil, model_helper.NewAppError("createSessionForUserAccessToken", MissingAccountError, nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model_helper.NewAppError("createSessionForUserAccessToken", "app.user.get.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	if !*a.srv.Config().ServiceSettings.EnableUserAccessTokens {
		return nil, model_helper.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "EnableUserAccessTokens=false", http.StatusUnauthorized)
	}

	if user.DeleteAt != 0 {
		return nil, model_helper.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "inactive_user_id="+user.ID, http.StatusUnauthorized)
	}

	session := model.Session{
		Token:   token.Token,
		UserID:  user.ID,
		Roles:   model_helper.UserGetRawRoles(*user),
		IsOauth: false,
		Props: model_types.JSONString{
			model_helper.SESSION_PROP_USER_ACCESS_TOKEN_ID: token.ID,
			model_helper.SESSION_PROP_TYPE:                 model_helper.SESSION_TYPE_USER_ACCESS_TOKEN,
		},
	}

	a.SetSessionExpireInDays(&session, model_helper.SESSION_USER_ACCESS_TOKEN_EXPIRY)

	savedSession, nErr := a.srv.Store.Session().Save(session)
	if nErr != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(nErr, &invErr):
			return nil, model_helper.NewAppError("CreateSession", "app.session.save.existing.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("CreateSession", "app.session.save.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	a.AddSessionToCache(savedSession)

	return savedSession, nil
}

// SetSessionExpireInDays sets the session's expiry the specified number of days
// relative to either the session creation date or the current time, depending
// on the `ExtendSessionOnActivity` config setting.
func (a *ServiceAccount) SetSessionExpireInDays(session *model.Session, days int) {
	if session.CreatedAt == 0 || *a.srv.Config().ServiceSettings.ExtendSessionLengthWithActivity {
		session.ExpiresAt = model_helper.GetMillis() + (1000 * 60 * 60 * 24 * int64(days))
	} else {
		session.ExpiresAt = session.CreatedAt + (1000 * 60 * 60 * 24 * int64(days))
	}
}

func (a *ServiceAccount) SessionCacheLength() int {
	if l, err := a.sessionCache.Len(); err == nil {
		return l
	}

	return 0
}

func (a *ServiceAccount) RevokeSessionsForDeviceId(userID string, deviceID string, currentSessionId string) *model_helper.AppError {
	sessions, err := a.srv.Store.Session().GetSessions(userID)
	if err != nil {
		return model_helper.NewAppError("RevokeSessionsForDeviceId", "app.session.get_sessions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, session := range sessions {
		if session.DeviceID == deviceID && session.ID != currentSessionId {
			slog.Debug("Revoking sessionId for userId. Re-login with the same devide Id", slog.String("session_id", session.ID), slog.String("user_id", userID))
			if err := a.RevokeSession(session); err != nil {
				slog.Warn("Could not revoke session for device", slog.String("device_id", deviceID), slog.Err(err))
			}
		}
	}

	return nil
}

func (a *ServiceAccount) GetSessionById(sessionID string) (*model.Session, *model_helper.AppError) {
	session, err := a.srv.Store.Session().Get(context.Background(), sessionID)
	if err != nil {
		return nil, model_helper.NewAppError("GetSessionById", "app.session.get.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	return session, nil
}

func (a *ServiceAccount) AttachDeviceId(sessionID string, deviceID string, expiresAt int64) *model_helper.AppError {
	_, err := a.srv.Store.Session().UpdateDeviceId(sessionID, deviceID, expiresAt)
	if err != nil {
		return model_helper.NewAppError("AttachDeviceId", "app.session.update_device_id.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceAccount) UpdateLastActivityAtIfNeeded(session model.Session) {
	now := model_helper.GetMillis()

	// TODO: studyme
	// a.UpdateWebConnUserActivity(session, now)

	if now-session.LastActivityAt < model_helper.SESSION_ACTIVITY_TIMEOUT {
		return
	}

	if err := a.srv.Store.Session().UpdateLastActivityAt(session.ID, now); err != nil {
		slog.Warn("Failed to update LastActivityAt", slog.String("user_id", session.UserID), slog.String("session_id", session.ID), slog.Err(err))
	}

	session.LastActivityAt = now
	a.AddSessionToCache(&session)
}

// GetSessionLengthInMillis returns the session length, in milliseconds,
// based on the type of session (Mobile, SSO, Web/LDAP).
func (a *ServiceAccount) GetSessionLengthInMillis(session *model.Session) int64 {
	if session == nil {
		return 0
	}
	var days int
	switch {
	case model_helper.SessionIsMobileApp(*session):
		days = *a.srv.Config().ServiceSettings.SessionLengthMobileInDays
	case model_helper.SessionIsSSOLogin(*session):
		days = *a.srv.Config().ServiceSettings.SessionLengthSSOInDays
	default:
		days = *a.srv.Config().ServiceSettings.SessionLengthWebInDays
	}
	return int64(days * 24 * 60 * 60 * 1000)
}

func (a *ServiceAccount) CreateUserAccessToken(token model.UserAccessToken) (*model.UserAccessToken, *model_helper.AppError) {
	user, appErr := a.UserById(context.Background(), token.UserID)
	if appErr != nil {
		return nil, appErr
	}

	if !*a.srv.Config().ServiceSettings.EnableUserAccessTokens {
		return nil, model_helper.NewAppError("CreateUserAccessToken", "app.user_access_token.disabled", nil, "", http.StatusNotImplemented)
	}

	token.Token = model_helper.NewId()

	savedToken, nErr := a.srv.Store.UserAccessToken().Save(token)
	if nErr != nil {
		var appErr *model_helper.AppError
		switch {
		case errors.As(nErr, &appErr):
			return nil, appErr
		default:
			return nil, model_helper.NewAppError("CreateUserAccessToken", "app.user_access_token.save.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	// Don't send emails to bot users.
	if err := a.srv.EmailService.SendUserAccessTokenAddedEmail(user.Email, user.Locale.String(), a.srv.GetSiteURL()); err != nil {
		a.srv.Log.Error("Unable to send user access token added email", slog.Err(err), slog.String("user_id", user.ID))
	}

	return savedToken, nil
}

func (a *ServiceAccount) RevokeUserAccessToken(token *model.UserAccessToken) *model_helper.AppError {
	var session *model.Session
	session, _ = a.srv.Store.Session().Get(context.Background(), token.Token)

	if err := a.srv.Store.UserAccessToken().Delete(token.ID); err != nil {
		return model_helper.NewAppError("RevokeUserAccessToken", "app.user_access_token.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if session == nil {
		return nil
	}

	return a.RevokeSession(session)
}

func (a *ServiceAccount) DisableUserAccessToken(token *model.UserAccessToken) *model_helper.AppError {
	var session *model.Session
	session, _ = a.srv.Store.Session().Get(context.Background(), token.Token)

	if err := a.srv.Store.UserAccessToken().UpdateTokenDisable(token.ID); err != nil {
		return model_helper.NewAppError("DisableUserAccessToken", "app.user_access_token.update_token_disable.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if session == nil {
		return nil
	}

	return a.RevokeSession(session)
}

func (a *ServiceAccount) EnableUserAccessToken(token *model.UserAccessToken) *model_helper.AppError {
	var session *model.Session
	session, _ = a.srv.Store.Session().Get(context.Background(), token.Token)

	err := a.srv.Store.UserAccessToken().UpdateTokenEnable(token.ID)
	if err != nil {
		return model_helper.NewAppError("EnableUserAccessToken", "app.user_access_token.update_token_enable.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if session == nil {
		return nil
	}

	return nil
}

func (a *ServiceAccount) GetUserAccessTokens(page, perPage int) (model.UserAccessTokenSlice, *model_helper.AppError) {
	tokens, err := a.srv.Store.
		UserAccessToken().
		GetAll(qm.Limit(perPage), qm.Offset(page*perPage))
	if err != nil {
		return nil, model_helper.NewAppError("GetUserAccessTokens", "app.user_access_token.get_all.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, token := range tokens {
		token.Token = ""
	}

	return tokens, nil
}

func (a *ServiceAccount) GetUserAccessTokensForUser(userID string, page, perPage int) (model.UserAccessTokenSlice, *model_helper.AppError) {
	tokens, err := a.srv.Store.
		UserAccessToken().
		GetAll(qm.Limit(perPage), qm.Offset(page*perPage), model.UserAccessTokenWhere.UserID.EQ(userID))
	if err != nil {
		return nil, model_helper.NewAppError("GetUserAccessTokensForUser", "app.user_access_token.get_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, token := range tokens {
		token.Token = ""
	}

	return tokens, nil
}

func (a *ServiceAccount) GetUserAccessToken(tokenID string, sanitize bool) (*model.UserAccessToken, *model_helper.AppError) {
	token, err := a.srv.Store.UserAccessToken().Get(tokenID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("GetUserAccessToken", "app.user.accesstoken_missing.app_error", nil, err.Error(), statusCode)
	}

	if sanitize {
		token.Token = ""
	}
	return token, nil
}

func (a *ServiceAccount) SearchUserAccessTokens(term string) (model.UserAccessTokenSlice, *model_helper.AppError) {
	tokens, err := a.srv.Store.UserAccessToken().Search(term)
	if err != nil {
		return nil, model_helper.NewAppError("SearchUserAccessTokens", "app.user_access_token.search.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, token := range tokens {
		token.Token = ""
	}
	return tokens, nil
}

// ExtendSessionExpiryIfNeeded extends Session.ExpiresAt based on session lengths in config.
// A new ExpiresAt is only written if enough time has elapsed since last update.
// Returns true only if the session was extended.
func (a *ServiceAccount) ExtendSessionExpiryIfNeeded(session *model.Session) bool {
	if !*a.srv.Config().ServiceSettings.ExtendSessionLengthWithActivity {
		return false
	}

	if session == nil || model_helper.SessionIsExpired(*session) {
		return false
	}

	sessionLength := a.GetSessionLengthInMillis(session)

	// Only extend the expiry if the lessor of 1% or 1 day has elapsed within the
	// current session duration.
	threshold := int64(math.Min(float64(sessionLength)*0.01, float64(24*60*60*1000)))
	// Minimum session length is 1 day as of this writing, therefore a minimum ~14 minutes threshold.
	// However we'll add a sanity check here in case that changes. Minimum 5 minute threshold,
	// meaning we won't write a new expiry more than every 5 minutes.
	if threshold < 5*60*1000 {
		threshold = 5 * 60 * 1000
	}

	now := model_helper.GetMillis()
	elapsed := now - (session.ExpiresAt - sessionLength)
	if elapsed < threshold {
		return false
	}

	// auditRec := a.MakeAuditRecord("extendSessionExpiry", audit.Fail)
	// defer a.LogAuditRec(auditRec, nil)
	// auditRec.AddMeta("session", session)

	newExpiry := now + sessionLength
	if err := a.srv.Store.Session().UpdateExpiresAt(session.ID, newExpiry); err != nil {
		slog.Error("Failed to update ExpiresAt", slog.String("user_id", session.UserID), slog.String("session_id", session.ID), slog.Err(err))
		// auditRec.AddMeta("err", err.Error())
		return false
	}

	// Update local cache. No need to invalidate cache for cluster as the session cache timeout
	// ensures each node will get an extended expiry within the next 10 minutes.
	// Worst case is another node may generate a redundant expiry update.
	session.ExpiresAt = newExpiry
	a.AddSessionToCache(session)

	slog.Debug(
		"Session extended",
		slog.String("user_id", session.UserID),
		slog.String("session_id", session.ID),
		slog.Int64("newExpiry", newExpiry),
		slog.Int64("session_length", sessionLength),
	)

	// auditRec.Success()
	// auditRec.AddMeta("extended_session", session)
	return true
}

func (a *ServiceAccount) GetCloudSession(token string) (*model.Session, *model_helper.AppError) {
	apiKey := os.Getenv("SN_CLOUD_API_KEY")
	if apiKey != "" && apiKey == token {
		// Need a bare-bones session object for later checks
		session := &model.Session{
			Token:   token,
			IsOauth: false,
			Props: model_types.JSONString{
				model_helper.SESSION_PROP_TYPE: model_helper.SESSION_TYPE_CLOUD_KEY,
			},
		}
		return session, nil
	}

	return nil, model_helper.NewAppError("GetCloudSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "The provided token is invalid", http.StatusUnauthorized)
}
