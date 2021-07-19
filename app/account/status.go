package account

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (a *AppAccount) AddStatusCacheSkipClusterSend(status *account.Status) {
	a.Srv().StatusCache.Set(status.UserId, status)
}

func (a *AppAccount) AddStatusCache(status *account.Status) {
	a.AddStatusCacheSkipClusterSend(status)

	if a.Cluster() != nil {
		msg := &cluster.ClusterMessage{
			Event:    cluster.CLUSTER_EVENT_UPDATE_STATUS,
			SendType: cluster.CLUSTER_SEND_BEST_EFFORT,
			Data:     status.ToClusterJson(),
		}
		a.Cluster().SendClusterMessage(msg)
	}
}

func (a *AppAccount) StatusByID(statusID string) (*account.Status, *model.AppError) {
	status, err := a.Srv().Store.Status().Get(statusID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("StatusByID", "app.account.status_missing.app_error", err)
	}

	return status, nil
}

func (a *AppAccount) StatusesByIDs(statusIDs []string) ([]*account.Status, *model.AppError) {
	statuses, err := a.Srv().Store.Status().GetByIds(statusIDs)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("StatusesByIDs", "app.account.statuses_missing.app_error", err)
	}

	return statuses, nil
}

func (a *AppAccount) GetUserStatusesByIds(userIDs []string) ([]*account.Status, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableUserStatuses {
		return []*account.Status{}, nil
	}

	var statusMap []*account.Status

	missingUserIds := []string{}
	for _, userID := range userIDs {
		var status *account.Status
		if err := a.Srv().StatusCache.Get(userID, &status); err == nil {
			statusMap = append(statusMap, status)
			if a.Metrics() != nil {
				a.Metrics().IncrementMemCacheHitCounter("Status")
			}
		} else {
			missingUserIds = append(missingUserIds, userID)
			if a.Metrics() != nil {
				a.Metrics().IncrementMemCacheMissCounter("Status")
			}
		}
	}

	if len(missingUserIds) > 0 {
		statuses, appErr := a.StatusesByIDs(missingUserIds)
		if appErr != nil {
			return nil, appErr
		}

		for _, s := range statuses {
			a.AddStatusCacheSkipClusterSend(s)
		}

		statusMap = append(statusMap, statuses...)
	}

	// For the case where the user does not have a row in the Status table and cache
	// remove the existing ids from missingUserIds and then create a offline state for the missing ones
	// This also return the status offline for the non-existing Ids in the system
	for i := 0; i < len(missingUserIds); i++ {
		missingUserId := missingUserIds[i]
		for _, userMap := range statusMap {
			if missingUserId == userMap.UserId {
				missingUserIds = append(missingUserIds[:i], missingUserIds[i+1:]...)
				i--
				break
			}
		}
	}

	for _, userID := range missingUserIds {
		statusMap = append(statusMap, &account.Status{UserId: userID, Status: account.STATUS_OFFLINE})
	}

	return statusMap, nil
}

func (a *AppAccount) SetStatusOnline(userID string, manual bool) {
	if !*a.Config().ServiceSettings.EnableUserStatuses {
		return
	}

	broadcast := false

	var oldStatus string = account.STATUS_OFFLINE
	var oldTime int64
	var oldManual bool
	var status *account.Status
	var err *model.AppError

	if status, err = a.GetStatus(userID); err != nil {
		status = &account.Status{UserId: userID, Status: account.STATUS_ONLINE, Manual: false, LastActivityAt: model.GetMillis(), ActiveChannel: ""}
		broadcast = true
	} else {
		if status.Manual && !manual {
			return // manually set status always overrides non-manual one
		}

		if status.Status != account.STATUS_ONLINE {
			broadcast = true
		}

		oldStatus = status.Status
		oldTime = status.LastActivityAt
		oldManual = status.Manual

		status.Status = account.STATUS_ONLINE
		status.Manual = false // for "online" there's no manual setting
		status.LastActivityAt = model.GetMillis()
	}

	a.AddStatusCache(status)

	// Only update the database if the status has changed, the status has been manually set,
	// or enough time has passed since the previous action
	if status.Status != oldStatus || status.Manual != oldManual || status.LastActivityAt-oldTime > account.STATUS_MIN_UPDATE_TIME {
		if broadcast {
			if err := a.Srv().Store.Status().SaveOrUpdate(status); err != nil {
				slog.Warn("Failed to save status", slog.String("user_id", userID), slog.Err(err), slog.String("user_id", userID))
			}
		} else {
			if err := a.Srv().Store.Status().UpdateLastActivityAt(status.UserId, status.LastActivityAt); err != nil {
				slog.Error("Failed to save status", slog.String("user_id", userID), slog.Err(err), slog.String("user_id", userID))
			}
		}
	}

	if broadcast {
		a.BroadcastStatus(status)
	}
}

func (a *AppAccount) SetStatusOffline(userID string, manual bool) {
	if !*a.Config().ServiceSettings.EnableUserStatuses {
		return
	}

	status, err := a.GetStatus(userID)
	if err == nil && status.Manual && !manual {
		return // manually set status always overrides non-manual one
	}

	status = &account.Status{UserId: userID, Status: account.STATUS_OFFLINE, Manual: manual, LastActivityAt: model.GetMillis(), ActiveChannel: ""}

	a.SaveAndBroadcastStatus(status)
}

func (a *AppAccount) SaveAndBroadcastStatus(status *account.Status) {
	a.AddStatusCache(status)

	if err := a.Srv().Store.Status().SaveOrUpdate(status); err != nil {
		slog.Warn("Failed to save status", slog.String("user_id", status.UserId), slog.Err(err))
	}

	a.BroadcastStatus(status)
}

func (a *AppAccount) BroadcastStatus(status *account.Status) {
	if a.Srv().Busy.IsBusy() {
		// this is considered a non-critical service and will be disabled when server busy.
		return
	}
	event := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_STATUS_CHANGE, status.UserId, nil)
	event.Add("status", status.Status)
	event.Add("user_id", status.UserId)
	a.Publish(event)
}
