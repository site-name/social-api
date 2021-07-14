package account

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/cluster"
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
