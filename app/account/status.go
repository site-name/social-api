package account

import (
	"encoding/json"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAccount) AddStatusCacheSkipClusterSend(status model.Status) {
	a.statusCache.Set(status.UserID, &status)
}

func (a *ServiceAccount) AddStatusCache(status model.Status) {
	a.AddStatusCacheSkipClusterSend(status)

	if a.cluster != nil {
		data, _ := json.Marshal(status)
		msg := &model_helper.ClusterMessage{
			Event:    model_helper.ClusterEventUpdateStatus,
			SendType: model_helper.ClusterSendBestEffort,
			Data:     data,
		}
		a.cluster.SendClusterMessage(msg)
	}
}

func (a *ServiceAccount) StatusByID(statusID string) (*model.Status, *model_helper.AppError) {
	status, err := a.srv.Store.Status().Get(statusID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("StatusByID", "app.user.status_by_id.app_error", nil, err.Error(), statusCode)
	}

	return status, nil
}

func (a *ServiceAccount) StatusesByIDs(statusIDs []string) (model.StatusSlice, *model_helper.AppError) {
	statuses, err := a.srv.Store.Status().GetByIds(statusIDs)
	if err != nil {
		return nil, model_helper.NewAppError("StatusesByIDs", "app.user.statuses_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return statuses, nil
}

func (a *ServiceAccount) GetUserStatusesByIds(userIDs []string) (model.StatusSlice, *model_helper.AppError) {
	if !*a.srv.Config().ServiceSettings.EnableUserStatuses {
		return model.StatusSlice{}, nil
	}

	var statusMap model.StatusSlice

	missingUserIds := []string{}
	for _, userID := range userIDs {
		var status *model.Status
		if err := a.statusCache.Get(userID, &status); err == nil {
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
			a.AddStatusCacheSkipClusterSend(*s)
		}

		statusMap = append(statusMap, statuses...)
	}

	// For the case where the user does not have a row in the Status table and cache
	// remove the existing ids from missingUserIds and then create a offline state for the missing ones
	// This also return the status offline for the non-existing Ids in the system
	for i := 0; i < len(missingUserIds); i++ {
		missingUserId := missingUserIds[i]
		for _, userMap := range statusMap {
			if missingUserId == userMap.UserID {
				missingUserIds = append(missingUserIds[:i], missingUserIds[i+1:]...)
				i--
				break
			}
		}
	}

	for _, userID := range missingUserIds {
		statusMap = append(statusMap, &model.Status{UserID: userID, Status: model_helper.STATUS_OFFLINE})
	}

	return statusMap, nil
}

func (a *ServiceAccount) SetStatusOnline(userID string, manual bool) {
	if !*a.srv.Config().ServiceSettings.EnableUserStatuses {
		return
	}

	broadcast := false

	var oldStatus string = model_helper.STATUS_OFFLINE
	var oldTime int64
	var oldManual bool
	var status *model.Status
	var err error

	status, err = a.GetStatus(userID)
	if err != nil {
		status = &model.Status{UserID: userID, Status: model_helper.STATUS_ONLINE, Manual: false, LastActivityAt: model_helper.GetMillis()}
		broadcast = true
	} else {
		if status.Manual && !manual {
			return // manually set status always overrides non-manual one
		}

		if status.Status != model_helper.STATUS_ONLINE {
			broadcast = true
		}

		oldStatus = status.Status
		oldTime = status.LastActivityAt
		oldManual = status.Manual

		status.Status = model_helper.STATUS_ONLINE
		status.Manual = false // for "online" there's no manual setting
		status.LastActivityAt = model_helper.GetMillis()
	}

	a.AddStatusCache(*status)

	// Only update the database if the status has changed, the status has been manually set,
	// or enough time has passed since the previous action
	if status.Status != oldStatus || status.Manual != oldManual || status.LastActivityAt-oldTime > model_helper.STATUS_MIN_UPDATE_TIME {
		if broadcast {
			if status, err = a.srv.Store.Status().Upsert(*status); err != nil {
				slog.Warn("Failed to save status", slog.String("user_id", userID), slog.Err(err), slog.String("user_id", userID))
			}
		} else {
			if err := a.srv.Store.Status().UpdateLastActivityAt(status.UserID, status.LastActivityAt); err != nil {
				slog.Error("Failed to save status", slog.String("user_id", userID), slog.Err(err), slog.String("user_id", userID))
			}
		}
	}

	if broadcast {
		a.BroadcastStatus(*status)
	}
}

func (a *ServiceAccount) SetStatusOffline(userID string, manual bool) {
	if !*a.srv.Config().ServiceSettings.EnableUserStatuses {
		return
	}

	stt, err := a.GetStatus(userID)
	if err == nil && stt.Manual && !manual {
		return // manually set status always overrides non-manual one
	}

	status := model.Status{UserID: userID, Status: model_helper.STATUS_OFFLINE, Manual: manual, LastActivityAt: model_helper.GetMillis()}

	a.SaveAndBroadcastStatus(status)
}

func (a *ServiceAccount) SaveAndBroadcastStatus(status model.Status) {
	a.AddStatusCache(status)

	savedStatus, err := a.srv.Store.Status().Upsert(status)
	if err != nil {
		slog.Warn("Failed to save status", slog.String("user_id", status.UserID), slog.Err(err))
	}

	a.BroadcastStatus(*savedStatus)
}

func (a *ServiceAccount) BroadcastStatus(status model.Status) {
	if a.srv.Busy.IsBusy() {
		// this is considered a non-critical service and will be disabled when server busy.
		return
	}
	event := model_helper.NewWebSocketEvent(model_helper.WebsocketEventStatusChange, status.UserID, nil)
	event.Add("status", status.Status)
	event.Add("user_id", status.UserID)
	a.srv.Publish(event)
}
