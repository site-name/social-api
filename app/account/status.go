package account

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAccount) AddStatusCacheSkipClusterSend(status *model.Status) {
	a.srv.StatusCache.Set(status.UserId, status)
}

func (a *ServiceAccount) AddStatusCache(status *model.Status) {
	a.AddStatusCacheSkipClusterSend(status)

	if a.cluster != nil {
		msg := &model.ClusterMessage{
			Event:    model.ClusterEventUpdateStatus,
			SendType: model.ClusterSendBestEffort,
			Data:     []byte(status.ToClusterJson()),
		}
		a.cluster.SendClusterMessage(msg)
	}
}

func (a *ServiceAccount) StatusByID(statusID string) (*model.Status, *model.AppError) {
	status, err := a.srv.Store.Status().Get(statusID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("StatusByID", "app.model.status_missing.app_error", err)
	}

	return status, nil
}

func (a *ServiceAccount) StatusesByIDs(statusIDs []string) ([]*model.Status, *model.AppError) {
	statuses, err := a.srv.Store.Status().GetByIds(statusIDs)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("StatusesByIDs", "app.model.statuses_missing.app_error", err)
	}

	return statuses, nil
}

func (a *ServiceAccount) GetUserStatusesByIds(userIDs []string) ([]*model.Status, *model.AppError) {
	if !*a.srv.Config().ServiceSettings.EnableUserStatuses {
		return []*model.Status{}, nil
	}

	var statusMap []*model.Status

	missingUserIds := []string{}
	for _, userID := range userIDs {
		var status *model.Status
		if err := a.srv.StatusCache.Get(userID, &status); err == nil {
			statusMap = append(statusMap, status)
			if a.metrics != nil {
				a.metrics.IncrementMemCacheHitCounter("Status")
			}
		} else {
			missingUserIds = append(missingUserIds, userID)
			if a.metrics != nil {
				a.metrics.IncrementMemCacheMissCounter("Status")
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
		statusMap = append(statusMap, &model.Status{UserId: userID, Status: model.STATUS_OFFLINE})
	}

	return statusMap, nil
}

func (a *ServiceAccount) SetStatusOnline(userID string, manual bool) {
	if !*a.srv.Config().ServiceSettings.EnableUserStatuses {
		return
	}

	broadcast := false

	var oldStatus string = model.STATUS_OFFLINE
	var oldTime int64
	var oldManual bool
	var status *model.Status
	var err *model.AppError

	if status, err = a.GetStatus(userID); err != nil {
		status = &model.Status{UserId: userID, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: model.GetMillis()}
		broadcast = true
	} else {
		if status.Manual && !manual {
			return // manually set status always overrides non-manual one
		}

		if status.Status != model.STATUS_ONLINE {
			broadcast = true
		}

		oldStatus = status.Status
		oldTime = status.LastActivityAt
		oldManual = status.Manual

		status.Status = model.STATUS_ONLINE
		status.Manual = false // for "online" there's no manual setting
		status.LastActivityAt = model.GetMillis()
	}

	a.AddStatusCache(status)

	// Only update the database if the status has changed, the status has been manually set,
	// or enough time has passed since the previous action
	if status.Status != oldStatus || status.Manual != oldManual || status.LastActivityAt-oldTime > model.STATUS_MIN_UPDATE_TIME {
		if broadcast {
			if err := a.srv.Store.Status().SaveOrUpdate(status); err != nil {
				slog.Warn("Failed to save status", slog.String("user_id", userID), slog.Err(err), slog.String("user_id", userID))
			}
		} else {
			if err := a.srv.Store.Status().UpdateLastActivityAt(status.UserId, status.LastActivityAt); err != nil {
				slog.Error("Failed to save status", slog.String("user_id", userID), slog.Err(err), slog.String("user_id", userID))
			}
		}
	}

	if broadcast {
		a.BroadcastStatus(status)
	}
}

func (a *ServiceAccount) SetStatusOffline(userID string, manual bool) {
	if !*a.srv.Config().ServiceSettings.EnableUserStatuses {
		return
	}

	status, err := a.GetStatus(userID)
	if err == nil && status.Manual && !manual {
		return // manually set status always overrides non-manual one
	}

	status = &model.Status{UserId: userID, Status: model.STATUS_OFFLINE, Manual: manual, LastActivityAt: model.GetMillis()}

	a.SaveAndBroadcastStatus(status)
}

func (a *ServiceAccount) SaveAndBroadcastStatus(status *model.Status) {
	a.AddStatusCache(status)

	if err := a.srv.Store.Status().SaveOrUpdate(status); err != nil {
		slog.Warn("Failed to save status", slog.String("user_id", status.UserId), slog.Err(err))
	}

	a.BroadcastStatus(status)
}

func (a *ServiceAccount) BroadcastStatus(status *model.Status) {
	if a.srv.Busy.IsBusy() {
		// this is considered a non-critical service and will be disabled when server busy.
		return
	}
	event := model.NewWebSocketEvent(model.WebsocketEventStatusChange, status.UserId, nil)
	event.Add("status", status.Status)
	event.Add("user_id", status.UserId)
	a.srv.Publish(event)
}
