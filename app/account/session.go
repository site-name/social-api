package account

import (
	"context"
	"errors"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/modules/audit"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore"
)

// CreateSession try saving given session to the database. If success then add that session to cache.
func (a *AppAccount) CreateSession(session *model.Session) (*model.Session, *model.AppError) {
	session.Token = ""
	session, err := a.Srv().Store.Session().Save(session)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*model.AppError); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("CreateSession", "app.session.save_session.app_error", nil, err.Error(), statusCode)
	}

	a.AddSessionToCache(session)

	return session, nil
}

func (a *AppAccount) ReturnSessionToPool(session *model.Session) {
	if session != nil {
		session.Id = ""
		a.sessionPool.Put(session)
	}
}

// AddSessionToCache add given session `s` to server's sessionCache, key is session's Token, expiry time as in config
func (a *AppAccount) AddSessionToCache(s *model.Session) {
	a.Srv().SessionCache.SetWithExpiry(s.Token, s, time.Duration(*a.Config().ServiceSettings.SessionCacheInMinutes))
}

// GetSessions get session from database with UserID attribute of given `userID`
func (a *AppAccount) GetSessions(userID string) ([]*model.Session, *model.AppError) {

	sessions, err := a.Srv().Store.Session().GetSessions(userID)
	if err != nil {
		return nil, model.NewAppError("GetSessions", "app.session.get_sessions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return sessions, nil
}

// RevokeAllSessions get session from database that has UserID of given userID, then removes it
func (a *AppAccount) RevokeAllSessions(userID string) *model.AppError {
	sessions, err := a.Srv().Store.Session().GetSessions(userID)
	if err != nil {
		return model.NewAppError("RevokeAllSessions", "app.session.get_sessions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, session := range sessions {
		if session.IsOAuth {
			// TODO: fixme
			// a.RevokeAccessToken(session.Token)
		} else {
			if err := a.Srv().Store.Session().Remove(session.Id); err != nil {
				return model.NewAppError("RevokeAllSessions", "app.session.remove.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	a.ClearSessionCacheForUser(userID)

	return nil
}

// ClearSessionCacheForUser clears all sessions that have `UserID` attribute of given `userID` in server's `sessionCache`
func (a *AppAccount) ClearSessionCacheForUser(userID string) {
	a.ClearSessionCacheForUserSkipClusterSend(userID)

	if a.Cluster() != nil {
		msg := &cluster.ClusterMessage{
			Event:    cluster.CLUSTER_EVENT_CLEAR_SESSION_CACHE_FOR_USER,
			SendType: cluster.CLUSTER_SEND_RELIABLE,
			Data:     userID,
		}
		a.Cluster().SendClusterMessage(msg)
	}
}

// ClearSessionCacheForUserSkipClusterSend iterates through server's sessionCache, if it finds any session belong to given userID, removes that session.
func (a *AppAccount) ClearSessionCacheForUserSkipClusterSend(userID string) {
	a.Srv().ClearSessionCacheForUserSkipClusterSend(userID)
}

func (a *AppAccount) GetSession(token string) (*model.Session, *model.AppError) {
	metrics := a.Metrics()

	var session = a.sessionPool.Get().(*model.Session)

	var err *model.AppError
	if err := a.Srv().SessionCache.Get(token, session); err == nil {
		if metrics != nil {
			metrics.IncrementMemCacheHitCounterSession()
		}
	} else {
		if metrics != nil {
			metrics.IncrementMemCacheMissCounterSession()
		}
	}

	if session.Id == "" {
		var nErr error
		if session, nErr = a.Srv().Store.Session().Get(sqlstore.WithMaster(context.Background()), token); nErr == nil {
			if session != nil {
				if session.Token != token {
					return nil, model.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "session token is different from the one in DB", http.StatusUnauthorized)
				}

				if !session.IsExpired() {
					a.AddSessionToCache(session)
				}
			}
		} else if nfErr := new(store.ErrNotFound); !errors.As(nErr, &nfErr) {
			return nil, model.NewAppError("GetSession", "app.session.get.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	if session == nil || session.Id == "" {
		session, err = a.createSessionForUserAccessToken(token)
		if err != nil {
			detailedError := ""
			statusCode := http.StatusUnauthorized
			if err.Id != "app.user_access_token.invalid_or_missing" {
				detailedError = err.Error()
				statusCode = err.StatusCode
			} else {
				slog.Warn("Error while creating session for user access token", slog.Err(err))
			}
			return nil, model.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": detailedError}, "", statusCode)
		}
	}

	if session.Id == "" || session.IsExpired() {
		return nil, model.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "session is either nil or expired", http.StatusUnauthorized)
	}

	if *a.Config().ServiceSettings.SessionIdleTimeoutInMinutes > 0 &&
		!session.IsOAuth && !session.IsMobileApp() &&
		session.Props[model.SESSION_PROP_TYPE] != model.SESSION_TYPE_USER_ACCESS_TOKEN &&
		!*a.Config().ServiceSettings.ExtendSessionLengthWithActivity {

		timeout := int64(*a.Config().ServiceSettings.SessionIdleTimeoutInMinutes) * 1000 * 60
		if (model.GetMillis() - session.LastActivityAt) > timeout {
			// Revoking the session is an asynchronous task anyways since we are not checking
			// for the return value of the call before returning the error.
			// So moving this to a goroutine has 2 advantages:
			// 1. We are treating this as a proper asynchronous task.
			// 2. This also fixes a race condition in the web hub, where GetSession
			// gets called from (*WebConn).isMemberOfTeam and revoking a session involves
			// clearing the webconn cache, which needs the hub again.
			a.Srv().Go(func() {
				err := a.RevokeSessionById(session.Id)
				if err != nil {
					slog.Warn("Error while revoking session", slog.Err(err))
				}
			})
			return nil, model.NewAppError("GetSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "idle timeout", http.StatusUnauthorized)
		}
	}

	return session, nil
}

// RevokeSessionById gets session with given sessionID then revokes it
func (a *AppAccount) RevokeSessionById(sessionID string) *model.AppError {
	session, err := a.Srv().Store.Session().Get(context.Background(), sessionID)
	if err != nil {
		return model.NewAppError("RevokeSessionById", "app.session.get.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	return a.RevokeSession(session)
}

// func (a *AppAccount) ClearSessionCacheForAllUsersSkipClusterSend() {
// 	a.Srv().clearSessionCacheForAllUsersSkipClusterSend()
// }

// RevokeSession removes session from database
func (a *AppAccount) RevokeSession(session *model.Session) *model.AppError {
	if session.IsOAuth {
		// TODO: fixme
		// if err := a.RevokeAccessToken(session.Token); err != nil {
		// 	return err
		// }
	} else {
		if err := a.Srv().Store.Session().Remove(session.Id); err != nil {
			return model.NewAppError("RevokeSession", "app.session.remove.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	a.ClearSessionCacheForUser(session.UserId)

	return nil
}

func (a *AppAccount) createSessionForUserAccessToken(tokenString string) (*model.Session, *model.AppError) {
	token, nErr := a.Srv().Store.UserAccessToken().GetByToken(tokenString)
	if nErr != nil {
		return nil, model.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, nErr.Error(), http.StatusUnauthorized)
	}

	if !token.IsActive {
		return nil, model.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "inactive_token", http.StatusUnauthorized)
	}

	user, nErr := a.Srv().Store.User().Get(context.Background(), token.UserId)
	if nErr != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(nErr, &nfErr):
			return nil, model.NewAppError("createSessionForUserAccessToken", MissingAccountError, nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("createSessionForUserAccessToken", "app.user.get.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	if !*a.Config().ServiceSettings.EnableUserAccessTokens {
		return nil, model.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "EnableUserAccessTokens=false", http.StatusUnauthorized)
	}

	if user.DeleteAt != 0 {
		return nil, model.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "inactive_user_id="+user.Id, http.StatusUnauthorized)
	}

	session := &model.Session{
		Token:   token.Token,
		UserId:  user.Id,
		Roles:   user.GetRawRoles(),
		IsOAuth: false,
	}

	session.AddProp(model.SESSION_PROP_USER_ACCESS_TOKEN_ID, token.Id)
	session.AddProp(model.SESSION_PROP_TYPE, model.SESSION_TYPE_USER_ACCESS_TOKEN)
	// if user.IsBot {
	// 	session.AddProp(model.SESSION_PROP_IS_BOT, model.SESSION_PROP_IS_BOT_VALUE)
	// }
	if user.IsGuest() {
		session.AddProp(model.SESSION_PROP_IS_GUEST, "true")
	} else {
		session.AddProp(model.SESSION_PROP_IS_GUEST, "false")
	}
	a.SetSessionExpireInDays(session, model.SESSION_USER_ACCESS_TOKEN_EXPIRY)

	session, nErr = a.Srv().Store.Session().Save(session)
	if nErr != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(nErr, &invErr):
			return nil, model.NewAppError("CreateSession", "app.session.save.existing.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model.NewAppError("CreateSession", "app.session.save.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	a.AddSessionToCache(session)

	return session, nil
}

// SetSessionExpireInDays sets the session's expiry the specified number of days
// relative to either the session creation date or the current time, depending
// on the `ExtendSessionOnActivity` config setting.
func (a *AppAccount) SetSessionExpireInDays(session *model.Session, days int) {
	if session.CreateAt == 0 || *a.Config().ServiceSettings.ExtendSessionLengthWithActivity {
		session.ExpiresAt = model.GetMillis() + (1000 * 60 * 60 * 24 * int64(days))
	} else {
		session.ExpiresAt = session.CreateAt + (1000 * 60 * 60 * 24 * int64(days))
	}
}

func (a *AppAccount) SessionCacheLength() int {
	if l, err := a.Srv().SessionCache.Len(); err == nil {
		return l
	}

	return 0
}

func (a *AppAccount) RevokeSessionsForDeviceId(userID string, deviceID string, currentSessionId string) *model.AppError {
	sessions, err := a.Srv().Store.Session().GetSessions(userID)
	if err != nil {
		return model.NewAppError("RevokeSessionsForDeviceId", "app.session.get_sessions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, session := range sessions {
		if session.DeviceId == deviceID && session.Id != currentSessionId {
			slog.Debug("Revoking sessionId for userId. Re-login with the same devide Id", slog.String("session_id", session.Id), slog.String("user_id", userID))
			if err := a.RevokeSession(session); err != nil {
				slog.Warn("Could not revoke session for device", slog.String("device_id", deviceID), slog.Err(err))
			}
		}
	}

	return nil
}

func (a *AppAccount) GetSessionById(sessionID string) (*model.Session, *model.AppError) {
	session, err := a.Srv().Store.Session().Get(context.Background(), sessionID)
	if err != nil {
		return nil, model.NewAppError("GetSessionById", "app.session.get.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	return session, nil
}

func (a *AppAccount) AttachDeviceId(sessionID string, deviceID string, expiresAt int64) *model.AppError {
	_, err := a.Srv().Store.Session().UpdateDeviceId(sessionID, deviceID, expiresAt)
	if err != nil {
		return model.NewAppError("AttachDeviceId", "app.session.update_device_id.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *AppAccount) UpdateLastActivityAtIfNeeded(session model.Session) {
	now := model.GetMillis()

	// TODO: studyme
	// a.UpdateWebConnUserActivity(session, now)

	if now-session.LastActivityAt < model.SESSION_ACTIVITY_TIMEOUT {
		return
	}

	if err := a.Srv().Store.Session().UpdateLastActivityAt(session.Id, now); err != nil {
		slog.Warn("Failed to update LastActivityAt", slog.String("user_id", session.UserId), slog.String("session_id", session.Id), slog.Err(err))
	}

	session.LastActivityAt = now
	a.AddSessionToCache(&session)
}

// GetSessionLengthInMillis returns the session length, in milliseconds,
// based on the type of session (Mobile, SSO, Web/LDAP).
func (a *AppAccount) GetSessionLengthInMillis(session *model.Session) int64 {
	if session == nil {
		return 0
	}

	var days int
	if session.IsMobileApp() {
		days = *a.Config().ServiceSettings.SessionLengthMobileInDays
	} else if session.IsSSOLogin() {
		days = *a.Config().ServiceSettings.SessionLengthSSOInDays
	} else {
		days = *a.Config().ServiceSettings.SessionLengthWebInDays
	}
	return int64(days * 24 * 60 * 60 * 1000)
}

func (a *AppAccount) CreateUserAccessToken(token *account.UserAccessToken) (*account.UserAccessToken, *model.AppError) {
	user, appErr := a.UserById(context.Background(), token.UserId)
	if appErr != nil {
		return nil, appErr
	}

	if !*a.Config().ServiceSettings.EnableUserAccessTokens /*&& !user.IsBot*/ {
		return nil, model.NewAppError("CreateUserAccessToken", "app.user_access_token.disabled", nil, "", http.StatusNotImplemented)
	}

	token.Token = model.NewId()

	token, nErr := a.Srv().Store.UserAccessToken().Save(token)
	if nErr != nil {
		var appErr *model.AppError
		switch {
		case errors.As(nErr, &appErr):
			return nil, appErr
		default:
			return nil, model.NewAppError("CreateUserAccessToken", "app.user_access_token.save.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	// Don't send emails to bot users.
	if err := a.Srv().EmailService.SendUserAccessTokenAddedEmail(user.Email, user.Locale, a.GetSiteURL()); err != nil {
		a.Log().Error("Unable to send user access token added email", slog.Err(err), slog.String("user_id", user.Id))
	}

	return token, nil
}

func (a *AppAccount) RevokeUserAccessToken(token *account.UserAccessToken) *model.AppError {
	var session *model.Session
	session, _ = a.Srv().Store.Session().Get(context.Background(), token.Token)

	if err := a.Srv().Store.UserAccessToken().Delete(token.Id); err != nil {
		return model.NewAppError("RevokeUserAccessToken", "app.user_access_token.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if session == nil {
		return nil
	}

	return a.RevokeSession(session)
}

func (a *AppAccount) DisableUserAccessToken(token *account.UserAccessToken) *model.AppError {
	var session *model.Session
	session, _ = a.Srv().Store.Session().Get(context.Background(), token.Token)

	if err := a.Srv().Store.UserAccessToken().UpdateTokenDisable(token.Id); err != nil {
		return model.NewAppError("DisableUserAccessToken", "app.user_access_token.update_token_disable.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if session == nil {
		return nil
	}

	return a.RevokeSession(session)
}

func (a *AppAccount) EnableUserAccessToken(token *account.UserAccessToken) *model.AppError {
	var session *model.Session
	session, _ = a.Srv().Store.Session().Get(context.Background(), token.Token)

	err := a.Srv().Store.UserAccessToken().UpdateTokenEnable(token.Id)
	if err != nil {
		return model.NewAppError("EnableUserAccessToken", "app.user_access_token.update_token_enable.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if session == nil {
		return nil
	}

	return nil
}

func (a *AppAccount) GetUserAccessTokens(page, perPage int) ([]*account.UserAccessToken, *model.AppError) {
	tokens, err := a.Srv().Store.UserAccessToken().GetAll(page*perPage, perPage)
	if err != nil {
		return nil, model.NewAppError("GetUserAccessTokens", "app.user_access_token.get_all.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, token := range tokens {
		token.Token = ""
	}

	return tokens, nil
}

func (a *AppAccount) GetUserAccessTokensForUser(userID string, page, perPage int) ([]*account.UserAccessToken, *model.AppError) {
	tokens, err := a.Srv().Store.UserAccessToken().GetByUser(userID, page*perPage, perPage)
	if err != nil {
		return nil, model.NewAppError("GetUserAccessTokensForUser", "app.user_access_token.get_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, token := range tokens {
		token.Token = ""
	}

	return tokens, nil
}

func (a *AppAccount) GetUserAccessToken(tokenID string, sanitize bool) (*account.UserAccessToken, *model.AppError) {
	token, err := a.Srv().Store.UserAccessToken().Get(tokenID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetUserAccessToken", "app.user.accesstoken_missing.app_error", err)
	}

	if sanitize {
		token.Token = ""
	}
	return token, nil
}

func (a *AppAccount) SearchUserAccessTokens(term string) ([]*account.UserAccessToken, *model.AppError) {
	tokens, err := a.Srv().Store.UserAccessToken().Search(term)
	if err != nil {
		return nil, model.NewAppError("SearchUserAccessTokens", "app.user_access_token.search.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, token := range tokens {
		token.Token = ""
	}
	return tokens, nil
}

// ExtendSessionExpiryIfNeeded extends Session.ExpiresAt based on session lengths in config.
// A new ExpiresAt is only written if enough time has elapsed since last update.
// Returns true only if the session was extended.
func (a *AppAccount) ExtendSessionExpiryIfNeeded(session *model.Session) bool {
	if !*a.Srv().Config().ServiceSettings.ExtendSessionLengthWithActivity {
		return false
	}

	if session == nil || session.IsExpired() {
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

	now := model.GetMillis()
	elapsed := now - (session.ExpiresAt - sessionLength)
	if elapsed < threshold {
		return false
	}

	auditRec := a.MakeAuditRecord("extendSessionExpiry", audit.Fail)
	defer a.LogAuditRec(auditRec, nil)
	auditRec.AddMeta("session", session)

	newExpiry := now + sessionLength
	if err := a.Srv().Store.Session().UpdateExpiresAt(session.Id, newExpiry); err != nil {
		slog.Error("Failed to update ExpiresAt", slog.String("user_id", session.UserId), slog.String("session_id", session.Id), slog.Err(err))
		auditRec.AddMeta("err", err.Error())
		return false
	}

	// Update local cache. No need to invalidate cache for cluster as the session cache timeout
	// ensures each node will get an extended expiry within the next 10 minutes.
	// Worst case is another node may generate a redundant expiry update.
	session.ExpiresAt = newExpiry
	a.AddSessionToCache(session)

	slog.Debug("Session extended", slog.String("user_id", session.UserId), slog.String("session_id", session.Id),
		slog.Int64("newExpiry", newExpiry), slog.Int64("session_length", sessionLength))

	auditRec.Success()
	auditRec.AddMeta("extended_session", session)
	return true
}

func (a *AppAccount) GetCloudSession(token string) (*model.Session, *model.AppError) {
	apiKey := os.Getenv("SN_CLOUD_API_KEY")
	if apiKey != "" && apiKey == token {
		// Need a bare-bones session object for later checks
		session := &model.Session{
			Token:   token,
			IsOAuth: false,
		}

		session.AddProp(model.SESSION_PROP_TYPE, model.SESSION_TYPE_CLOUD_KEY)
		return session, nil
	}

	return nil, model.NewAppError("GetCloudSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "The provided token is invalid", http.StatusUnauthorized)
}

// func (a *AppAccount) GetRemoteClusterSession(token string, remoteId string) (*model.Session, *model.AppError) {
// 	rc, appErr := a.GetRemoteCluster(remoteId)
// 	if appErr == nil && rc.Token == token {
// 		// Need a bare-bones session object for later checks
// 		session := &model.Session{
// 			Token:   token,
// 			IsOAuth: false,
// 		}

// 		session.AddProp(model.SESSION_PROP_TYPE, model.SESSION_TYPE_REMOTECLUSTER_TOKEN)
// 		return session, nil
// 	}
// 	return nil, model.NewAppError("GetRemoteClusterSession", "api.context.invalid_token.error", map[string]interface{}{"Token": token, "Error": ""}, "The provided token is invalid", http.StatusUnauthorized)
// }
