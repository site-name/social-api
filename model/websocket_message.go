package model

import (
	"encoding/json"
	"fmt"
)

const (
	WEBSOCKET_EVENT_TYPING                                   = "typing"
	WEBSOCKET_EVENT_POSTED                                   = "posted"
	WEBSOCKET_EVENT_POST_EDITED                              = "post_edited"
	WEBSOCKET_EVENT_POST_DELETED                             = "post_deleted"
	WEBSOCKET_EVENT_POST_UNREAD                              = "post_unread"
	WEBSOCKET_EVENT_CHANNEL_CONVERTED                        = "channel_converted"
	WEBSOCKET_EVENT_CHANNEL_CREATED                          = "channel_created"
	WEBSOCKET_EVENT_CHANNEL_DELETED                          = "channel_deleted"
	WEBSOCKET_EVENT_CHANNEL_RESTORED                         = "channel_restored"
	WEBSOCKET_EVENT_CHANNEL_UPDATED                          = "channel_updated"
	WEBSOCKET_EVENT_CHANNEL_MEMBER_UPDATED                   = "channel_member_updated"
	WEBSOCKET_EVENT_CHANNEL_SCHEME_UPDATED                   = "channel_scheme_updated"
	WEBSOCKET_EVENT_DIRECT_ADDED                             = "direct_added"
	WEBSOCKET_EVENT_GROUP_ADDED                              = "group_added"
	WEBSOCKET_EVENT_NEW_USER                                 = "new_user"
	WEBSOCKET_EVENT_ADDED_TO_TEAM                            = "added_to_team"
	WEBSOCKET_EVENT_LEAVE_TEAM                               = "leave_team"
	WEBSOCKET_EVENT_UPDATE_TEAM                              = "update_team"
	WEBSOCKET_EVENT_DELETE_TEAM                              = "delete_team"
	WEBSOCKET_EVENT_RESTORE_TEAM                             = "restore_team"
	WEBSOCKET_EVENT_UPDATE_TEAM_SCHEME                       = "update_team_scheme"
	WEBSOCKET_EVENT_USER_ADDED                               = "user_added"
	WEBSOCKET_EVENT_USER_UPDATED                             = "user_updated"
	WEBSOCKET_EVENT_USER_ROLE_UPDATED                        = "user_role_updated"
	WEBSOCKET_EVENT_MEMBERROLE_UPDATED                       = "memberrole_updated"
	WEBSOCKET_EVENT_USER_REMOVED                             = "user_removed"
	WEBSOCKET_EVENT_PREFERENCE_CHANGED                       = "preference_changed"
	WEBSOCKET_EVENT_PREFERENCES_CHANGED                      = "preferences_changed"
	WEBSOCKET_EVENT_PREFERENCES_DELETED                      = "preferences_deleted"
	WEBSOCKET_EVENT_EPHEMERAL_MESSAGE                        = "ephemeral_message"
	WEBSOCKET_EVENT_STATUS_CHANGE                            = "status_change"
	WEBSOCKET_EVENT_HELLO                                    = "hello"
	WEBSOCKET_AUTHENTICATION_CHALLENGE                       = "authentication_challenge"
	WEBSOCKET_EVENT_REACTION_ADDED                           = "reaction_added"
	WEBSOCKET_EVENT_REACTION_REMOVED                         = "reaction_removed"
	WEBSOCKET_EVENT_RESPONSE                                 = "response"
	WEBSOCKET_EVENT_EMOJI_ADDED                              = "emoji_added"
	WEBSOCKET_EVENT_CHANNEL_VIEWED                           = "channel_viewed"
	WEBSOCKET_EVENT_PLUGIN_STATUSES_CHANGED                  = "plugin_statuses_changed"
	WEBSOCKET_EVENT_PLUGIN_ENABLED                           = "plugin_enabled"
	WEBSOCKET_EVENT_PLUGIN_DISABLED                          = "plugin_disabled"
	WEBSOCKET_EVENT_ROLE_UPDATED                             = "role_updated"
	WEBSOCKET_EVENT_LICENSE_CHANGED                          = "license_changed"
	WEBSOCKET_EVENT_CONFIG_CHANGED                           = "config_changed"
	WEBSOCKET_EVENT_OPEN_DIALOG                              = "open_dialog"
	WEBSOCKET_EVENT_GUESTS_DEACTIVATED                       = "guests_deactivated"
	WEBSOCKET_EVENT_USER_ACTIVATION_STATUS_CHANGE            = "user_activation_status_change"
	WEBSOCKET_EVENT_RECEIVED_GROUP                           = "received_group"
	WEBSOCKET_EVENT_RECEIVED_GROUP_ASSOCIATED_TO_TEAM        = "received_group_associated_to_team"
	WEBSOCKET_EVENT_RECEIVED_GROUP_NOT_ASSOCIATED_TO_TEAM    = "received_group_not_associated_to_team"
	WEBSOCKET_EVENT_RECEIVED_GROUP_ASSOCIATED_TO_CHANNEL     = "received_group_associated_to_channel"
	WEBSOCKET_EVENT_RECEIVED_GROUP_NOT_ASSOCIATED_TO_CHANNEL = "received_group_not_associated_to_channel"
	WEBSOCKET_EVENT_SIDEBAR_CATEGORY_CREATED                 = "sidebar_category_created"
	WEBSOCKET_EVENT_SIDEBAR_CATEGORY_UPDATED                 = "sidebar_category_updated"
	WEBSOCKET_EVENT_SIDEBAR_CATEGORY_DELETED                 = "sidebar_category_deleted"
	WEBSOCKET_EVENT_SIDEBAR_CATEGORY_ORDER_UPDATED           = "sidebar_category_order_updated"
	WEBSOCKET_WARN_METRIC_STATUS_RECEIVED                    = "warn_metric_status_received"
	WEBSOCKET_WARN_METRIC_STATUS_REMOVED                     = "warn_metric_status_removed"
	WEBSOCKET_EVENT_CLOUD_PAYMENT_STATUS_UPDATED             = "cloud_payment_status_updated"
	WEBSOCKET_EVENT_THREAD_UPDATED                           = "thread_updated"
	WEBSOCKET_EVENT_THREAD_FOLLOW_CHANGED                    = "thread_follow_changed"
	WEBSOCKET_EVENT_THREAD_READ_CHANGED                      = "thread_read_changed"
	WEBSOCKET_FIRST_ADMIN_VISIT_MARKETPLACE_STATUS_RECEIVED  = "first_admin_visit_marketplace_status_received"
)

type WebSocketMessage interface {
	ToJson() string
	IsValid() bool
	EventType() string
}

type WebsocketBroadcast struct {
	OmitUsers             map[string]bool `json:"omit_users"` // broadcast is omitted for users listed here
	UserId                string          `json:"user_id"`    // broadcast only occurs for this user
	ContainsSanitizedData bool            `json:"-"`
	ContainsSensitiveData bool            `json:"-"`
	// ChannelId             string          `json:"channel_id"` // broadcast only occurs for users in this channel
	// TeamId                string          `json:"team_id"`    // broadcast only occurs for users in this team
}

type precomputedWebSocketEventJSON struct {
	Event     json.RawMessage
	Data      json.RawMessage
	Broadcast json.RawMessage
}

// **NOTE**: Direct access to WebSocketEvent fields is deprecated. They will be
// made unexported in next major version release. Provided getter functions should be used instead.
type WebSocketEvent struct {
	Event           string                 // Deprecated: use EventType()
	Data            map[string]interface{} // Deprecated: use GetData()
	Broadcast       *WebsocketBroadcast    // Deprecated: use GetBroadcast()
	Sequence        int64                  // Deprecated: use GetSequence()
	precomputedJSON *precomputedWebSocketEventJSON
}

// webSocketEventJSON mirrors WebSocketEvent to make some of its unexported fields serializable
type webSocketEventJSON struct {
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Broadcast *WebsocketBroadcast    `json:"broadcast"`
	Sequence  int64                  `json:"seq"`
}

func NewWebSocketEvent(event, userId string, omitUsers map[string]bool) *WebSocketEvent {
	return &WebSocketEvent{
		Event: event,
		Data:  make(map[string]interface{}),
		Broadcast: &WebsocketBroadcast{
			UserId:    userId,
			OmitUsers: omitUsers,
		},
	}
}

func (ev *WebSocketEvent) IsValid() bool {
	return ev.Event != ""
}

func (ev *WebSocketEvent) Add(key string, value interface{}) {
	ev.Data[key] = value
}

func (ev *WebSocketEvent) EventType() string {
	return ev.Event
}

func (ev *WebSocketEvent) ToJson() string {
	if ev.precomputedJSON != nil {
		return fmt.Sprintf(`{"event": %s, "data": %s, "broadcast": %s, "seq": %d}`, ev.precomputedJSON.Event, ev.precomputedJSON.Data, ev.precomputedJSON.Broadcast, ev.Sequence)
	}

	b, _ := json.Marshal(webSocketEventJSON{
		ev.Event,
		ev.Data,
		ev.Broadcast,
		ev.Sequence,
	})
	return string(b)
}

func (ev *WebSocketEvent) GetBroadcast() *WebsocketBroadcast {
	return ev.Broadcast
}
