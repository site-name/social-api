package model_helper

import (
	"net/http"
	"os"

	"github.com/sitename/sitename/model"
)

const (
	CDS_OFFLINE_AFTER_MILLIS = 1000 * 60 * 30 // 30 minutes
	CDS_TYPE_APP             = "sitename_app"
)

type ClusterEvent string

const (
	ClusterEventPublish                                     ClusterEvent = "publish"
	ClusterEventUpdateStatus                                ClusterEvent = "update_status"
	ClusterEventInvalidateAllCaches                         ClusterEvent = "inv_all_caches"
	ClusterEventInvalidateCacheForReactions                 ClusterEvent = "inv_reactions"
	ClusterEventInvalidateCacheForChannelMembersNotifyProps ClusterEvent = "inv_channel_members_notify_props"
	ClusterEventInvalidateCacheForChannelByName             ClusterEvent = "inv_channel_name"
	ClusterEventInvalidateCacheForChannel                   ClusterEvent = "inv_channel"
	ClusterEventInvalidateCacheForChannelGuestCount         ClusterEvent = "inv_channel_guest_count"
	ClusterEventInvalidateCacheForUser                      ClusterEvent = "inv_user"
	ClusterEventInvalidateCacheForUserTeams                 ClusterEvent = "inv_user_teams"
	ClusterEventClearSessionCacheForUser                    ClusterEvent = "clear_session_user"
	ClusterEventInvalidateCacheForRoles                     ClusterEvent = "inv_roles"
	ClusterEventInvalidateCacheForRolePermissions           ClusterEvent = "inv_role_permissions"
	ClusterEventInvalidateCacheForProfileByIds              ClusterEvent = "inv_profile_ids"
	ClusterEventInvalidateCacheForProfileInChannel          ClusterEvent = "inv_profile_in_channel"
	ClusterEventInvalidateCacheForSchemes                   ClusterEvent = "inv_schemes"
	ClusterEventInvalidateCacheForFileInfos                 ClusterEvent = "inv_file_infos"
	ClusterEventInvalidateCacheForWebhooks                  ClusterEvent = "inv_webhooks"
	ClusterEventInvalidateCacheForEmojisById                ClusterEvent = "inv_emojis_by_id"
	ClusterEventInvalidateCacheForEmojisIdByName            ClusterEvent = "inv_emojis_id_by_name"
	ClusterEventInvalidateCacheForChannelPinnedpostsCounts  ClusterEvent = "inv_channel_pinnedposts_counts"
	ClusterEventInvalidateCacheForChannelMemberCounts       ClusterEvent = "inv_channel_member_counts"
	ClusterEventInvalidateCacheForLastPosts                 ClusterEvent = "inv_last_posts"
	ClusterEventInvalidateCacheForLastPostTime              ClusterEvent = "inv_last_post_time"
	ClusterEventInvalidateCacheForTeams                     ClusterEvent = "inv_teams"
	ClusterEventClearSessionCacheForAllUsers                ClusterEvent = "inv_all_user_sessions"
	ClusterEventInstallPlugin                               ClusterEvent = "install_plugin"
	ClusterEventRemovePlugin                                ClusterEvent = "remove_plugin"
	ClusterEventPluginEvent                                 ClusterEvent = "plugin_event"
	ClusterEventInvalidateCacheForTermsOfService            ClusterEvent = "inv_terms_of_service"
	ClusterEventBusyStateChanged                            ClusterEvent = "busy_state_change"
	ClusterEventInvalidateCacheForCategoryByIds             ClusterEvent = "inv_category_ids"

	// Gossip communication
	ClusterGossipEventRequestGetLogs            = "gossip_request_get_logs"
	ClusterGossipEventResponseGetLogs           = "gossip_response_get_logs"
	ClusterGossipEventRequestGetClusterStats    = "gossip_request_cluster_stats"
	ClusterGossipEventResponseGetClusterStats   = "gossip_response_cluster_stats"
	ClusterGossipEventRequestGetPluginStatuses  = "gossip_request_plugin_statuses"
	ClusterGossipEventResponseGetPluginStatuses = "gossip_response_plugin_statuses"
	ClusterGossipEventRequestSaveConfig         = "gossip_request_save_config"
	ClusterGossipEventResponseSaveConfig        = "gossip_response_save_config"

	// SendTypes for ClusterMessage.
	ClusterSendBestEffort = "best_effort"
	ClusterSendReliable   = "reliable"
)

type ClusterMessage struct {
	Event            ClusterEvent      `json:"event"`
	SendType         string            `json:"-"`
	WaitForAllToSend bool              `json:"-"`
	Data             []byte            `json:"data,omitempty"`
	Props            map[string]string `json:"props,omitempty"`
}

type ClusterStats struct {
	Id                        string `json:"id"`
	TotalWebsocketConnections int    `json:"total_websocket_connections"`
	TotalReadDbConnections    int    `json:"total_read_db_connections"`
	TotalMasterDbConnections  int    `json:"total_master_db_connections"`
}

type ClusterInfo struct {
	Id         string `json:"id"`
	Version    string `json:"version"`
	ConfigHash string `json:"config_hash"`
	IpAddress  string `json:"ipaddress"`
	Hostname   string `json:"hostname"`
}

func ClusterDiscoveryAutoFillHostname(c *model.ClusterDiscovery) {
	if c.HostName == "" {
		if hn, err := os.Hostname(); err == nil {
			c.HostName = hn
		}
	}
}

func ClusterDiscoveryAutoFillIpAddress(c *model.ClusterDiscovery, iface, ipAddress string) {
	if c.HostName == "" {
		if ipAddress != "" {
			c.HostName = ipAddress
		} else {
			c.HostName = GetServerIpAddress(iface)
		}
	}
}

func ClusterDiscoveriesAreEqual(c1 *model.ClusterDiscovery, c2 *model.ClusterDiscovery) bool {
	return c1 != nil &&
		c2 != nil &&
		c1.Type == c2.Type &&
		c1.ClusterName == c2.ClusterName &&
		c1.HostName == c2.HostName
}

func ClusterDiscoveryIsValid(c *model.ClusterDiscovery) *AppError {
	if c.ClusterName == "" {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.cluster_name.app_error", nil, "please provide cluster name", http.StatusBadRequest)
	}
	if c.Type == "" {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.type.app_error", nil, "please provide cluster type", http.StatusBadRequest)
	}
	if c.HostName == "" {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.host_name.app_error", nil, "please provide host name", http.StatusBadRequest)
	}

	return nil
}
