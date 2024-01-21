package model_helper

const (
	PREFERENCE_CATEGORY_DIRECT_CHANNEL_SHOW       = "direct_channel_show"
	PREFERENCE_CATEGORY_GROUP_CHANNEL_SHOW        = "group_channel_show"
	PREFERENCE_CATEGORY_TUTORIAL_STEPS            = "tutorial_step"
	PREFERENCE_CATEGORY_ADVANCED_SETTINGS         = "advanced_settings"
	PREFERENCE_CATEGORY_FLAGGED_POST              = "flagged_post"
	PREFERENCE_CATEGORY_FAVORITE_CHANNEL          = "favorite_channel"
	PREFERENCE_CATEGORY_SIDEBAR_SETTINGS          = "sidebar_settings"
	PREFERENCE_CATEGORY_DISPLAY_SETTINGS          = "display_settings"
	PREFERENCE_NAME_COLLAPSED_THREADS_ENABLED     = "collapsed_reply_threads"
	PREFERENCE_NAME_CHANNEL_DISPLAY_MODE          = "channel_display_mode"
	PREFERENCE_NAME_COLLAPSE_SETTING              = "collapse_previews"
	PREFERENCE_NAME_MESSAGE_DISPLAY               = "message_display"
	PREFERENCE_NAME_NAME_FORMAT                   = "name_format"
	PREFERENCE_NAME_USE_MILITARY_TIME             = "use_military_time"
	PREFERENCE_CATEGORY_THEME                     = "theme"     // the name for theme props is the team id
	PREFERENCE_CATEGORY_AUTHORIZED_OAUTH_APP      = "oauth_app" // the name for oauth_app is the client_id and value is the current scope
	PREFERENCE_CATEGORY_LAST                      = "last"
	PREFERENCE_NAME_LAST_CHANNEL                  = "channel"
	PREFERENCE_NAME_LAST_TEAM                     = "team"
	PREFERENCE_CATEGORY_CUSTOM_STATUS             = "custom_status"
	PREFERENCE_NAME_RECENT_CUSTOM_STATUSES        = "recent_custom_statuses"
	PREFERENCE_NAME_CUSTOM_STATUS_TUTORIAL_STATE  = "custom_status_tutorial_state"
	PREFERENCE_CUSTOM_STATUS_MODAL_VIEWED         = "custom_status_modal_viewed"
	PREFERENCE_CATEGORY_NOTIFICATIONS             = "notifications"
	PREFERENCE_NAME_EMAIL_INTERVAL                = "email_interval"
	PREFERENCE_EMAIL_INTERVAL_NO_BATCHING_SECONDS = "30"  // the "immediate" setting is actually 30s
	PREFERENCE_EMAIL_INTERVAL_BATCHING_SECONDS    = "900" // fifteen minutes is 900 seconds
	PREFERENCE_EMAIL_INTERVAL_IMMEDIATELY         = "immediately"
	PREFERENCE_EMAIL_INTERVAL_FIFTEEN             = "fifteen"
	PREFERENCE_EMAIL_INTERVAL_FIFTEEN_AS_SECONDS  = "900"
	PREFERENCE_EMAIL_INTERVAL_HOUR                = "hour"
	PREFERENCE_EMAIL_INTERVAL_HOUR_AS_SECONDS     = "3600"
)
