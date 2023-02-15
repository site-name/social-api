package model

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/ldap"

	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/slog"
)

const (
	CONN_SECURITY_NONE     = ""
	CONN_SECURITY_PLAIN    = "PLAIN"
	CONN_SECURITY_TLS      = "TLS"
	CONN_SECURITY_STARTTLS = "STARTTLS"

	IMAGE_DRIVER_LOCAL = "local"
	IMAGE_DRIVER_S3    = "amazons3"

	DATABASE_DRIVER_MYSQL    = "mysql"
	DATABASE_DRIVER_POSTGRES = "postgres"

	SEARCHENGINE_ELASTICSEARCH = "elasticsearch"

	MINIO_ACCESS_KEY = "minioaccesskey"
	MINIO_SECRET_KEY = "miniosecretkey"
	MINIO_BUCKET     = "mattermost-test"

	PASSWORD_MAXIMUM_LENGTH = 64
	PASSWORD_MINIMUM_LENGTH = 5

	SERVICE_GITLAB    = "gitlab"
	SERVICE_GOOGLE    = "google"
	SERVICE_OFFICE365 = "office365"
	SERVICE_OPENID    = "openid"

	GENERIC_NO_CHANNEL_NOTIFICATION = "generic_no_channel"
	GENERIC_NOTIFICATION            = "generic"
	GENERIC_NOTIFICATION_SERVER     = "https://push-test.mattermost.com"
	MM_SUPPORT_ADVISOR_ADDRESS      = "support-advisor@mattermost.com"
	FULL_NOTIFICATION               = "full"
	ID_LOADED_NOTIFICATION          = "id_loaded"

	DIRECT_MESSAGE_ANY  = "any"
	DIRECT_MESSAGE_TEAM = "team"

	SHOW_USERNAME          = "username"
	SHOW_NICKNAME_FULLNAME = "nickname_full_name"
	SHOW_FULLNAME          = "full_name"

	PERMISSIONS_ALL           = "all"
	PERMISSIONS_CHANNEL_ADMIN = "channel_admin"
	PERMISSIONS_TEAM_ADMIN    = "team_admin"
	PERMISSIONS_SYSTEM_ADMIN  = "system_admin"

	FAKE_SETTING = "********************************"

	PERMISSIONS_DELETE_POST_ALL          = "all"
	PERMISSIONS_DELETE_POST_TEAM_ADMIN   = "team_admin"
	PERMISSIONS_DELETE_POST_SYSTEM_ADMIN = "system_admin"

	ALLOW_EDIT_POST_ALWAYS     = "always"
	ALLOW_EDIT_POST_NEVER      = "never"
	ALLOW_EDIT_POST_TIME_LIMIT = "time_limit"

	GROUP_UNREAD_CHANNELS_DISABLED    = "disabled"
	GROUP_UNREAD_CHANNELS_DEFAULT_ON  = "default_on"
	GROUP_UNREAD_CHANNELS_DEFAULT_OFF = "default_off"

	COLLAPSED_THREADS_DISABLED    = "disabled"
	COLLAPSED_THREADS_DEFAULT_ON  = "default_on"
	COLLAPSED_THREADS_DEFAULT_OFF = "default_off"

	EMAIL_BATCHING_BUFFER_SIZE = 256
	EMAIL_BATCHING_INTERVAL    = 30

	EMAIL_NOTIFICATION_CONTENTS_FULL    = "full"
	EMAIL_NOTIFICATION_CONTENTS_GENERIC = "generic"

	SITENAME_MAX_LENGTH = 30

	SERVICE_SETTINGS_DEFAULT_SITE_URL           = "http://localhost:8065"
	SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE      = ""
	SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE       = ""
	SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT       = 300
	SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT      = 300
	SERVICE_SETTINGS_DEFAULT_IDLE_TIMEOUT       = 60
	SERVICE_SETTINGS_DEFAULT_MAX_LOGIN_ATTEMPTS = 10
	SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM    = ""
	SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS = ":8065"
	SERVICE_SETTINGS_DEFAULT_GFYCAT_API_KEY     = "2_KtH_W5"
	SERVICE_SETTINGS_DEFAULT_GFYCAT_API_SECRET  = "3wLVZPiswc3DnaiaFoLkDvB4X0IV6CpMkj4tf2inJRsBY6-FnkT08zGmppWFgeof"

	OPEN_EXCHANGE_RATE_API_KEY = "cb3c20ad2a624639806f043e4b5aafa4"

	SQL_SETTINGS_DEFAULT_DATA_SOURCE = "postgres://minh:anhyeuem98@localhost/sitename_test?sslmode=disable&connect_timeout=10"

	FILE_SETTINGS_DEFAULT_DIRECTORY = "./data/"

	IMPORT_SETTINGS_DEFAULT_DIRECTORY      = "./import"
	IMPORT_SETTINGS_DEFAULT_RETENTION_DAYS = 30

	EXPORT_SETTINGS_DEFAULT_DIRECTORY      = "./export"
	EXPORT_SETTINGS_DEFAULT_RETENTION_DAYS = 30

	EMAIL_SETTINGS_DEFAULT_FEEDBACK_ORGANIZATION = ""

	SUPPORT_SETTINGS_DEFAULT_TERMS_OF_SERVICE_LINK = "https://about.mattermost.com/default-terms/"
	SUPPORT_SETTINGS_DEFAULT_PRIVACY_POLICY_LINK   = "https://about.mattermost.com/default-privacy-policy/"
	SUPPORT_SETTINGS_DEFAULT_ABOUT_LINK            = "https://about.mattermost.com/default-about/"
	SUPPORT_SETTINGS_DEFAULT_HELP_LINK             = "https://about.mattermost.com/default-help/"
	SUPPORT_SETTINGS_DEFAULT_REPORT_A_PROBLEM_LINK = "https://about.mattermost.com/default-report-a-problem/"
	SUPPORT_SETTINGS_DEFAULT_SUPPORT_EMAIL         = "feedback@mattermost.com"
	SUPPORT_SETTINGS_DEFAULT_RE_ACCEPTANCE_PERIOD  = 365

	LDAP_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE         = ""
	LDAP_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE          = ""
	LDAP_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE              = ""
	LDAP_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE           = ""
	LDAP_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE           = ""
	LDAP_SETTINGS_DEFAULT_ID_ATTRIBUTE                 = ""
	LDAP_SETTINGS_DEFAULT_POSITION_ATTRIBUTE           = ""
	LDAP_SETTINGS_DEFAULT_LOGIN_FIELD_NAME             = ""
	LDAP_SETTINGS_DEFAULT_GROUP_DISPLAY_NAME_ATTRIBUTE = ""
	LDAP_SETTINGS_DEFAULT_GROUP_ID_ATTRIBUTE           = ""
	LDAP_SETTINGS_DEFAULT_PICTURE_ATTRIBUTE            = ""

	SAML_SETTINGS_DEFAULT_ID_ATTRIBUTE         = ""
	SAML_SETTINGS_DEFAULT_GUEST_ATTRIBUTE      = ""
	SAML_SETTINGS_DEFAULT_ADMIN_ATTRIBUTE      = ""
	SAML_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE = ""
	SAML_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE  = ""
	SAML_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE      = ""
	SAML_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE   = ""
	SAML_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE   = ""
	SAML_SETTINGS_DEFAULT_LOCALE_ATTRIBUTE     = ""
	SAML_SETTINGS_DEFAULT_POSITION_ATTRIBUTE   = ""

	SAML_SETTINGS_SIGNATURE_ALGORITHM_SHA1    = "RSAwithSHA1"
	SAML_SETTINGS_SIGNATURE_ALGORITHM_SHA256  = "RSAwithSHA256"
	SAML_SETTINGS_SIGNATURE_ALGORITHM_SHA512  = "RSAwithSHA512"
	SAML_SETTINGS_DEFAULT_SIGNATURE_ALGORITHM = SAML_SETTINGS_SIGNATURE_ALGORITHM_SHA1

	SAML_SETTINGS_CANONICAL_ALGORITHM_C14N    = "Canonical1.0"
	SAML_SETTINGS_CANONICAL_ALGORITHM_C14N11  = "Canonical1.1"
	SAML_SETTINGS_DEFAULT_CANONICAL_ALGORITHM = SAML_SETTINGS_CANONICAL_ALGORITHM_C14N

	NATIVEAPP_SETTINGS_DEFAULT_APP_DOWNLOAD_LINK         = "https://mattermost.com/download/#mattermostApps"
	NATIVEAPP_SETTINGS_DEFAULT_ANDROID_APP_DOWNLOAD_LINK = "https://about.mattermost.com/mattermost-android-app/"
	NATIVEAPP_SETTINGS_DEFAULT_IOS_APP_DOWNLOAD_LINK     = "https://about.mattermost.com/mattermost-ios-app/"

	EXPERIMENTAL_SETTINGS_DEFAULT_LINK_METADATA_TIMEOUT_MILLISECONDS = 5000

	ANALYTICS_SETTINGS_DEFAULT_MAX_USERS_FOR_STATISTICS = 2500

	ANNOUNCEMENT_SETTINGS_DEFAULT_BANNER_COLOR                    = "#f2a93b"
	ANNOUNCEMENT_SETTINGS_DEFAULT_BANNER_TEXT_COLOR               = "#333333"
	ANNOUNCEMENT_SETTINGS_DEFAULT_NOTICES_JSON_URL                = "https://notices.mattermost.com/"
	ANNOUNCEMENT_SETTINGS_DEFAULT_NOTICES_FETCH_FREQUENCY_SECONDS = 3600

	THEME_SETTING_DEFAULT = "default"

	ELASTICSEARCH_SETTINGS_DEFAULT_CONNECTION_URL                    = "http://localhost:9200"
	ELASTICSEARCH_SETTINGS_DEFAULT_USERNAME                          = "elastic"
	ELASTICSEARCH_SETTINGS_DEFAULT_PASSWORD                          = "changeme"
	ELASTICSEARCH_SETTINGS_DEFAULT_POST_INDEX_REPLICAS               = 1
	ELASTICSEARCH_SETTINGS_DEFAULT_POST_INDEX_SHARDS                 = 1
	ELASTICSEARCH_SETTINGS_DEFAULT_CHANNEL_INDEX_REPLICAS            = 1
	ELASTICSEARCH_SETTINGS_DEFAULT_CHANNEL_INDEX_SHARDS              = 1
	ELASTICSEARCH_SETTINGS_DEFAULT_USER_INDEX_REPLICAS               = 1
	ELASTICSEARCH_SETTINGS_DEFAULT_USER_INDEX_SHARDS                 = 1
	ELASTICSEARCH_SETTINGS_DEFAULT_AGGREGATE_POSTS_AFTER_DAYS        = 365
	ELASTICSEARCH_SETTINGS_DEFAULT_POSTS_AGGREGATOR_JOB_START_TIME   = "03:00"
	ELASTICSEARCH_SETTINGS_DEFAULT_INDEX_PREFIX                      = ""
	ELASTICSEARCH_SETTINGS_DEFAULT_LIVE_INDEXING_BATCH_SIZE          = 1
	ELASTICSEARCH_SETTINGS_DEFAULT_BULK_INDEXING_TIME_WINDOW_SECONDS = 3600
	ELASTICSEARCH_SETTINGS_DEFAULT_REQUEST_TIMEOUT_SECONDS           = 30

	BLEVE_SETTINGS_DEFAULT_INDEX_DIR                         = ""
	BLEVE_SETTINGS_DEFAULT_BULK_INDEXING_TIME_WINDOW_SECONDS = 3600

	DATA_RETENTION_SETTINGS_DEFAULT_MESSAGE_RETENTION_DAYS  = 365
	DATA_RETENTION_SETTINGS_DEFAULT_FILE_RETENTION_DAYS     = 365
	DATA_RETENTION_SETTINGS_DEFAULT_DELETION_JOB_START_TIME = "02:00"

	PLUGIN_SETTINGS_DEFAULT_DIRECTORY          = "./plugins"
	PLUGIN_SETTINGS_DEFAULT_CLIENT_DIRECTORY   = "./client/plugins"
	PLUGIN_SETTINGS_DEFAULT_ENABLE_MARKETPLACE = true
	PLUGIN_SETTINGS_DEFAULT_MARKETPLACE_URL    = "https://api.integrations.mattermost.com"
	PLUGIN_SETTINGS_OLD_MARKETPLACE_URL        = "https://marketplace.integrations.mattermost.com"

	COMPLIANCE_EXPORT_TYPE_CSV             = "csv"
	COMPLIANCE_EXPORT_TYPE_ACTIANCE        = "actiance"
	COMPLIANCE_EXPORT_TYPE_GLOBALRELAY     = "globalrelay"
	COMPLIANCE_EXPORT_TYPE_GLOBALRELAY_ZIP = "globalrelay-zip"
	GLOBALRELAY_CUSTOMER_TYPE_A9           = "A9"
	GLOBALRELAY_CUSTOMER_TYPE_A10          = "A10"

	CLIENT_SIDE_CERT_CHECK_PRIMARY_AUTH   = "primary"
	CLIENT_SIDE_CERT_CHECK_SECONDARY_AUTH = "secondary"

	IMAGE_PROXY_TYPE_LOCAL      = "local"
	IMAGE_PROXY_TYPE_ATMOS_CAMO = "atmos/camo"

	GOOGLE_SETTINGS_DEFAULT_SCOPE             = "profile email"
	GOOGLE_SETTINGS_DEFAULT_AUTH_ENDPOINT     = "https://accounts.google.com/o/oauth2/v2/auth"
	GOOGLE_SETTINGS_DEFAULT_TOKEN_ENDPOINT    = "https://www.googleapis.com/oauth2/v4/token"
	GOOGLE_SETTINGS_DEFAULT_USER_API_ENDPOINT = "https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses,nicknames,metadata"

	OFFICE365_SETTINGS_DEFAULT_SCOPE             = "User.Read"
	OFFICE365_SETTINGS_DEFAULT_AUTH_ENDPOINT     = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	OFFICE365_SETTINGS_DEFAULT_TOKEN_ENDPOINT    = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	OFFICE365_SETTINGS_DEFAULT_USER_API_ENDPOINT = "https://graph.microsoft.com/v1.0/me"

	CLOUD_SETTINGS_DEFAULT_CWS_URL     = "https://customers.mattermost.com"
	CLOUD_SETTINGS_DEFAULT_CWS_API_URL = "https://portal.internal.prod.cloud.mattermost.com"
	OPENID_SETTINGS_DEFAULT_SCOPE      = "profile openid email"

	LOCAL_MODE_SOCKET_PATH = "/var/tmp/sitename_local.socket"
)

func GetDefaultAppCustomURLSchemes() []string {
	return []string{"snauth://", "snauthbeta://"}
}

var ServerTLSSupportedCiphers = map[string]uint16{
	"TLS_RSA_WITH_RC4_128_SHA":                tls.TLS_RSA_WITH_RC4_128_SHA,
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
}

const VAT_LAYER_ACCESS_KEY = "bb98bc77f97b4b060ca39b78e27396ec"

type ThirdPartySettings struct {
	OpenExchangeApiEndpoint           *string
	OpenExchangeRateApiKey            *string `access:"experimental_features"`
	OpenExchangeRecuringDurationHours *int    `access:"experimental_features"`

	VatlayerApi       *string
	VatlayerAccessKey *string
	RateTypeURL       *string
	VateRateURL       *string
}

func (s *ThirdPartySettings) SetDefaults() {
	if s.OpenExchangeRateApiKey == nil {
		s.OpenExchangeRateApiKey = NewPrimitive(OPEN_EXCHANGE_RATE_API_KEY)
	}
	if s.OpenExchangeRecuringDurationHours == nil {
		s.OpenExchangeRecuringDurationHours = NewPrimitive(2)
	}
	if s.OpenExchangeApiEndpoint == nil {
		s.OpenExchangeApiEndpoint = NewPrimitive("http://openexchangerates.org/api/latest.json")
	}

	if s.VatlayerApi == nil {
		s.VatlayerApi = NewPrimitive("http://apilayer.net/api/")
	}
	if s.VatlayerAccessKey == nil {
		s.VatlayerAccessKey = NewPrimitive(VAT_LAYER_ACCESS_KEY)
	}
	if s.RateTypeURL == nil {
		s.RateTypeURL = NewPrimitive(fmt.Sprintf("http://apilayer.net/api/types?access_key=%s", VAT_LAYER_ACCESS_KEY))
	}
	if s.VateRateURL == nil {
		s.VateRateURL = NewPrimitive(fmt.Sprintf("http://apilayer.net/api/rate_list?access_key=%s", VAT_LAYER_ACCESS_KEY))
	}
}

type ServiceSettings struct {
	SiteURL                                           *string  `access:"environment_web_server,authentication_saml,write_restrictable"`
	SiteName                                          *string  `access:"environment_web_server,write_restrictable"`
	WebsocketURL                                      *string  `access:"write_restrictable,cloud_restrictable"`
	LicenseFileLocation                               *string  `access:"write_restrictable,cloud_restrictable"`                        // telemetry: none
	ListenAddress                                     *string  `access:"environment_web_server,write_restrictable,cloud_restrictable"` // telemetry: none
	ConnectionSecurity                                *string  `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	TLSCertFile                                       *string  `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	TLSKeyFile                                        *string  `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	TLSMinVer                                         *string  `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	TLSStrictTransport                                *bool    `access:"write_restrictable,cloud_restrictable"`
	TLSStrictTransportMaxAge                          *int64   `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	TLSOverwriteCiphers                               []string `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	UseLetsEncrypt                                    *bool    `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	LetsEncryptCertificateCacheFile                   *string  `access:"environment_web_server,write_restrictable,cloud_restrictable"` // telemetry: none
	Forward80To443                                    *bool    `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	TrustedProxyIPHeader                              []string `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	ReadTimeout                                       *int     `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	WriteTimeout                                      *int     `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	IdleTimeout                                       *int     `access:"write_restrictable,cloud_restrictable"`
	MaximumLoginAttempts                              *int     `access:"authentication_password,write_restrictable,cloud_restrictable"`
	GoroutineHealthThreshold                          *int     `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	EnableOAuthServiceProvider                        *bool    `access:"integrations_integration_management"`
	EnableIncomingWebhooks                            *bool    `access:"integrations_integration_management"`
	EnableOutgoingWebhooks                            *bool    `access:"integrations_integration_management"`
	EnableCommands                                    *bool    `access:"integrations_integration_management"`
	EnablePostUsernameOverride                        *bool    `access:"integrations_integration_management"`
	EnablePostIconOverride                            *bool    `access:"integrations_integration_management"`
	GoogleDeveloperKey                                *string  `access:"site_posts,write_restrictable,cloud_restrictable"`
	EnableLinkPreviews                                *bool    `access:"site_posts"`
	RestrictLinkPreviews                              *string  `access:"site_posts"`
	EnableTesting                                     *bool    `access:"environment_developer,write_restrictable,cloud_restrictable"`
	EnableDeveloper                                   *bool    `access:"environment_developer,write_restrictable,cloud_restrictable"`
	EnableOpenTracing                                 *bool    `access:"write_restrictable,cloud_restrictable"`
	EnableSecurityFixAlert                            *bool    `access:"environment_smtp,write_restrictable,cloud_restrictable"`
	EnableInsecureOutgoingConnections                 *bool    `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	AllowedUntrustedInternalConnections               *string  `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	EnableMultifactorAuthentication                   *bool    `access:"authentication_mfa"`
	EnforceMultifactorAuthentication                  *bool    `access:"authentication_mfa"`
	EnableUserAccessTokens                            *bool    `access:"integrations_integration_management"`
	AllowCorsFrom                                     *string  `access:"integrations_cors,write_restrictable,cloud_restrictable"`
	CorsExposedHeaders                                *string  `access:"integrations_cors,write_restrictable,cloud_restrictable"`
	CorsAllowCredentials                              *bool    `access:"integrations_cors,write_restrictable,cloud_restrictable"`
	CorsDebug                                         *bool    `access:"integrations_cors,write_restrictable,cloud_restrictable"`
	AllowCookiesForSubdomains                         *bool    `access:"write_restrictable,cloud_restrictable"`
	ExtendSessionLengthWithActivity                   *bool    `access:"environment_session_lengths,write_restrictable,cloud_restrictable"`
	SessionLengthWebInDays                            *int     `access:"environment_session_lengths,write_restrictable,cloud_restrictable"`
	SessionLengthMobileInDays                         *int     `access:"environment_session_lengths,write_restrictable,cloud_restrictable"`
	SessionLengthSSOInDays                            *int     `access:"environment_session_lengths,write_restrictable,cloud_restrictable"`
	SessionCacheInMinutes                             *int     `access:"environment_session_lengths,write_restrictable,cloud_restrictable"`
	SessionIdleTimeoutInMinutes                       *int     `access:"environment_session_lengths,write_restrictable,cloud_restrictable"`
	WebsocketSecurePort                               *int     `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	WebsocketPort                                     *int     `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	WebserverMode                                     *string  `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	EnableGifPicker                                   *bool    `access:"integrations_gif"`
	GfycatApiKey                                      *string  `access:"integrations_gif"`
	GfycatApiSecret                                   *string  `access:"integrations_gif"`
	EnableCustomEmoji                                 *bool    `access:"site_emoji"`
	EnableEmojiPicker                                 *bool    `access:"site_emoji"`
	PostEditTimeLimit                                 *int     `access:"user_management_permissions"`
	TimeBetweenUserTypingUpdatesMilliseconds          *int64   `access:"experimental_features,write_restrictable,cloud_restrictable"`
	EnablePostSearch                                  *bool    `access:"write_restrictable,cloud_restrictable"`
	EnableFileSearch                                  *bool    `access:"write_restrictable"`
	MinimumHashtagLength                              *int     `access:"environment_database,write_restrictable,cloud_restrictable"`
	EnableUserTypingMessages                          *bool    `access:"experimental_features,write_restrictable,cloud_restrictable"`
	EnableChannelViewedMessages                       *bool    `access:"experimental_features,write_restrictable,cloud_restrictable"`
	EnableUserStatuses                                *bool    `access:"write_restrictable,cloud_restrictable"`
	ExperimentalEnableAuthenticationTransfer          *bool    `access:"experimental_features,write_restrictable,cloud_restrictable"`
	ClusterLogTimeoutMilliseconds                     *int     `access:"write_restrictable,cloud_restrictable"`
	CloseUnusedDirectMessages                         *bool    `access:"experimental_features"`
	EnablePreviewFeatures                             *bool    `access:"experimental_features"`
	EnableTutorial                                    *bool    `access:"experimental_features"`
	ExperimentalEnableDefaultChannelLeaveJoinMessages *bool    `access:"experimental_features"`
	ExperimentalGroupUnreadChannels                   *string  `access:"experimental_features"`
	ExperimentalChannelOrganization                   *bool    `access:"experimental_features"`
	EnableAPITeamDeletion                             *bool
	EnableAPIUserDeletion                             *bool
	ExperimentalEnableHardenedMode                    *bool `access:"experimental_features"`
	DisableLegacyMFA                                  *bool `access:"write_restrictable,cloud_restrictable"`
	ExperimentalStrictCSRFEnforcement                 *bool `access:"experimental_features,write_restrictable,cloud_restrictable"`
	EnableEmailInvitations                            *bool `access:"authentication_signup"`
	DisableBotsWhenOwnerIsDeactivated                 *bool `access:"integrations_bot_accounts,write_restrictable,cloud_restrictable"`
	EnableBotAccountCreation                          *bool `access:"integrations_bot_accounts"`
	EnableSVGs                                        *bool `access:"site_posts"`
	EnableLatex                                       *bool `access:"site_posts"`
	EnableAPIChannelDeletion                          *bool
	EnableLocalMode                                   *bool
	LocalModeSocketLocation                           *string // telemetry: none
	EnableAWSMetering                                 *bool   // telemetry: none
	SplitKey                                          *string `access:"experimental_feature_flags,write_restrictable"` // telemetry: none
	FeatureFlagSyncIntervalSeconds                    *int    `access:"experimental_feature_flags,write_restrictable"` // telemetry: none
	DebugSplit                                        *bool   `access:"experimental_feature_flags,write_restrictable"` // telemetry: none
	ThreadAutoFollow                                  *bool   `access:"experimental_features"`
	CollapsedThreads                                  *string `access:"experimental_features"`
	ManagedResourcePaths                              *string `access:"environment_web_server,write_restrictable,cloud_restrictable"`
	EnableLegacySidebar                               *bool   `access:"experimental_features"`
	EnableReliableWebSockets                          *bool   `access:"experimental_features"` // telemetry: none
	EnablePermalinkPreviews                           *bool   `access:"site_posts"`
	EnableInlineLatex                                 *bool   `access:"site_posts"`

	DEPRECATED_DO_NOT_USE_ImageProxyType              *string `json:"ImageProxyType" mapstructure:"ImageProxyType"`                           // Deprecated: do not use
	DEPRECATED_DO_NOT_USE_ImageProxyURL               *string `json:"ImageProxyURL" mapstructure:"ImageProxyURL"`                             // Deprecated: do not use
	DEPRECATED_DO_NOT_USE_ImageProxyOptions           *string `json:"ImageProxyOptions" mapstructure:"ImageProxyOptions"`                     // Deprecated: do not use
	DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations *bool   `json:"EnableOnlyAdminIntegrations" mapstructure:"EnableOnlyAdminIntegrations"` // Deprecated: do not use
	DEPRECATED_DO_NOT_USE_RestrictPostDelete          *string `json:"RestrictPostDelete" mapstructure:"RestrictPostDelete"`                   // Deprecated: do not use
	DEPRECATED_DO_NOT_USE_AllowEditPost               *string `json:"AllowEditPost" mapstructure:"AllowEditPost"`                             // Deprecated: do not use
}

func (s *ServiceSettings) SetDefaults(isUpdate bool) {
	if s.EnableEmailInvitations == nil {
		// If the site URL is also not present then assume this is a clean install
		if s.SiteURL == nil {
			s.EnableEmailInvitations = NewPrimitive(false)
		} else {
			s.EnableEmailInvitations = NewPrimitive(true)
		}
	}
	if s.EnablePermalinkPreviews == nil {
		s.EnablePermalinkPreviews = NewPrimitive(true)
	}
	if s.EnableInlineLatex == nil {
		s.EnableInlineLatex = NewPrimitive(true)
	}

	if s.SiteName == nil {
		s.SiteName = NewPrimitive("SiteName")
	}

	if s.SiteURL == nil {
		if s.EnableDeveloper != nil && *s.EnableDeveloper {
			s.SiteURL = NewPrimitive(SERVICE_SETTINGS_DEFAULT_SITE_URL)
		} else {
			s.SiteURL = NewPrimitive("")
		}
	}

	if s.WebsocketURL == nil {
		s.WebsocketURL = NewPrimitive("")
	}

	if s.LicenseFileLocation == nil {
		s.LicenseFileLocation = NewPrimitive("")
	}

	if s.ListenAddress == nil {
		s.ListenAddress = NewPrimitive(SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS)
	}

	if s.EnableLinkPreviews == nil {
		s.EnableLinkPreviews = NewPrimitive(true)
	}

	if s.RestrictLinkPreviews == nil {
		s.RestrictLinkPreviews = NewPrimitive("")
	}

	if s.EnableTesting == nil {
		s.EnableTesting = NewPrimitive(false)
	}

	if s.EnableDeveloper == nil {
		s.EnableDeveloper = NewPrimitive(false)
	}

	if s.EnableOpenTracing == nil {
		s.EnableOpenTracing = NewPrimitive(false)
	}

	if s.EnableSecurityFixAlert == nil {
		s.EnableSecurityFixAlert = NewPrimitive(true)
	}

	if s.EnableInsecureOutgoingConnections == nil {
		s.EnableInsecureOutgoingConnections = NewPrimitive(false)
	}

	if s.AllowedUntrustedInternalConnections == nil {
		s.AllowedUntrustedInternalConnections = NewPrimitive("")
	}

	if s.EnableMultifactorAuthentication == nil {
		s.EnableMultifactorAuthentication = NewPrimitive(false)
	}

	if s.EnforceMultifactorAuthentication == nil {
		s.EnforceMultifactorAuthentication = NewPrimitive(false)
	}

	if s.EnableUserAccessTokens == nil {
		s.EnableUserAccessTokens = NewPrimitive(false)
	}

	if s.GoroutineHealthThreshold == nil {
		s.GoroutineHealthThreshold = NewPrimitive(-1)
	}

	if s.GoogleDeveloperKey == nil {
		s.GoogleDeveloperKey = NewPrimitive("")
	}

	if s.EnableOAuthServiceProvider == nil {
		s.EnableOAuthServiceProvider = NewPrimitive(false)
	}

	if s.EnableIncomingWebhooks == nil {
		s.EnableIncomingWebhooks = NewPrimitive(true)
	}

	if s.EnableOutgoingWebhooks == nil {
		s.EnableOutgoingWebhooks = NewPrimitive(true)
	}

	if s.ConnectionSecurity == nil {
		s.ConnectionSecurity = NewPrimitive("")
	}

	if s.TLSKeyFile == nil {
		s.TLSKeyFile = NewPrimitive(SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE)
	}

	if s.TLSCertFile == nil {
		s.TLSCertFile = NewPrimitive(SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE)
	}

	if s.TLSMinVer == nil {
		s.TLSMinVer = NewPrimitive("1.2")
	}

	if s.TLSStrictTransport == nil {
		s.TLSStrictTransport = NewPrimitive(false)
	}

	if s.TLSStrictTransportMaxAge == nil {
		s.TLSStrictTransportMaxAge = NewPrimitive[int64](63072000)
	}

	if s.TLSOverwriteCiphers == nil {
		s.TLSOverwriteCiphers = []string{}
	}

	if s.UseLetsEncrypt == nil {
		s.UseLetsEncrypt = NewPrimitive(false)
	}

	if s.LetsEncryptCertificateCacheFile == nil {
		s.LetsEncryptCertificateCacheFile = NewPrimitive("./config/letsencrypt.cache")
	}

	if s.ReadTimeout == nil {
		s.ReadTimeout = NewPrimitive(SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT)
	}

	if s.WriteTimeout == nil {
		s.WriteTimeout = NewPrimitive(SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT)
	}

	if s.IdleTimeout == nil {
		s.IdleTimeout = NewPrimitive(SERVICE_SETTINGS_DEFAULT_IDLE_TIMEOUT)
	}

	if s.MaximumLoginAttempts == nil {
		s.MaximumLoginAttempts = NewPrimitive(SERVICE_SETTINGS_DEFAULT_MAX_LOGIN_ATTEMPTS)
	}

	if s.Forward80To443 == nil {
		s.Forward80To443 = NewPrimitive(false)
	}

	if isUpdate {
		// When updating an existing configuration, ensure that defaults are set.
		if s.TrustedProxyIPHeader == nil {
			s.TrustedProxyIPHeader = []string{HEADER_FORWARDED, HEADER_REAL_IP}
		}
	} else {
		// When generating a blank configuration, leave the list empty.
		s.TrustedProxyIPHeader = []string{}
	}

	if s.TimeBetweenUserTypingUpdatesMilliseconds == nil {
		s.TimeBetweenUserTypingUpdatesMilliseconds = NewPrimitive[int64](5000)
	}

	if s.EnablePostSearch == nil {
		s.EnablePostSearch = NewPrimitive(true)
	}

	if s.EnableFileSearch == nil {
		s.EnableFileSearch = NewPrimitive(true)
	}

	if s.MinimumHashtagLength == nil {
		s.MinimumHashtagLength = NewPrimitive(3)
	}

	if s.EnableUserTypingMessages == nil {
		s.EnableUserTypingMessages = NewPrimitive(true)
	}

	if s.EnableChannelViewedMessages == nil {
		s.EnableChannelViewedMessages = NewPrimitive(true)
	}

	if s.EnableUserStatuses == nil {
		s.EnableUserStatuses = NewPrimitive(true)
	}

	if s.ClusterLogTimeoutMilliseconds == nil {
		s.ClusterLogTimeoutMilliseconds = NewPrimitive(2000)
	}

	if s.CloseUnusedDirectMessages == nil {
		s.CloseUnusedDirectMessages = NewPrimitive(false)
	}

	if s.EnableTutorial == nil {
		s.EnableTutorial = NewPrimitive(true)
	}

	// Must be manually enabled for existing installations.
	if s.ExtendSessionLengthWithActivity == nil {
		s.ExtendSessionLengthWithActivity = NewPrimitive(!isUpdate)
	}

	if s.SessionLengthWebInDays == nil {
		if isUpdate {
			s.SessionLengthWebInDays = NewPrimitive(180)
		} else {
			s.SessionLengthWebInDays = NewPrimitive(30)
		}
	}

	if s.SessionLengthMobileInDays == nil {
		if isUpdate {
			s.SessionLengthMobileInDays = NewPrimitive(180)
		} else {
			s.SessionLengthMobileInDays = NewPrimitive(30)
		}
	}

	if s.SessionLengthSSOInDays == nil {
		s.SessionLengthSSOInDays = NewPrimitive(30)
	}

	if s.SessionCacheInMinutes == nil {
		s.SessionCacheInMinutes = NewPrimitive(10)
	}

	if s.SessionIdleTimeoutInMinutes == nil {
		s.SessionIdleTimeoutInMinutes = NewPrimitive(43200)
	}

	if s.EnableCommands == nil {
		s.EnableCommands = NewPrimitive(true)
	}

	if s.DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations == nil {
		s.DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations = NewPrimitive(true)
	}

	if s.EnablePostUsernameOverride == nil {
		s.EnablePostUsernameOverride = NewPrimitive(false)
	}

	if s.EnablePostIconOverride == nil {
		s.EnablePostIconOverride = NewPrimitive(false)
	}

	if s.WebsocketPort == nil {
		s.WebsocketPort = NewPrimitive(80)
	}

	if s.WebsocketSecurePort == nil {
		s.WebsocketSecurePort = NewPrimitive(443)
	}

	if s.AllowCorsFrom == nil {
		s.AllowCorsFrom = NewPrimitive(SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM)
	}

	if s.CorsExposedHeaders == nil {
		s.CorsExposedHeaders = NewPrimitive("")
	}

	if s.CorsAllowCredentials == nil {
		s.CorsAllowCredentials = NewPrimitive(false)
	}

	if s.CorsDebug == nil {
		s.CorsDebug = NewPrimitive(false)
	}

	if s.AllowCookiesForSubdomains == nil {
		s.AllowCookiesForSubdomains = NewPrimitive(false)
	}

	if s.WebserverMode == nil {
		s.WebserverMode = NewPrimitive("gzip")
	} else if *s.WebserverMode == "regular" {
		*s.WebserverMode = "gzip"
	}

	if s.EnableCustomEmoji == nil {
		s.EnableCustomEmoji = NewPrimitive(true)
	}

	if s.EnableEmojiPicker == nil {
		s.EnableEmojiPicker = NewPrimitive(true)
	}

	if s.EnableGifPicker == nil {
		s.EnableGifPicker = NewPrimitive(true)
	}

	if s.GfycatApiKey == nil || *s.GfycatApiKey == "" {
		s.GfycatApiKey = NewPrimitive(SERVICE_SETTINGS_DEFAULT_GFYCAT_API_KEY)
	}

	if s.GfycatApiSecret == nil || *s.GfycatApiSecret == "" {
		s.GfycatApiSecret = NewPrimitive(SERVICE_SETTINGS_DEFAULT_GFYCAT_API_SECRET)
	}

	if s.DEPRECATED_DO_NOT_USE_RestrictPostDelete == nil {
		s.DEPRECATED_DO_NOT_USE_RestrictPostDelete = NewPrimitive(PERMISSIONS_DELETE_POST_ALL)
	}

	if s.DEPRECATED_DO_NOT_USE_AllowEditPost == nil {
		s.DEPRECATED_DO_NOT_USE_AllowEditPost = NewPrimitive(ALLOW_EDIT_POST_ALWAYS)
	}

	if s.ExperimentalEnableAuthenticationTransfer == nil {
		s.ExperimentalEnableAuthenticationTransfer = NewPrimitive(true)
	}

	if s.PostEditTimeLimit == nil {
		s.PostEditTimeLimit = NewPrimitive(-1)
	}

	if s.EnablePreviewFeatures == nil {
		s.EnablePreviewFeatures = NewPrimitive(true)
	}

	if s.ExperimentalEnableDefaultChannelLeaveJoinMessages == nil {
		s.ExperimentalEnableDefaultChannelLeaveJoinMessages = NewPrimitive(true)
	}

	if s.ExperimentalGroupUnreadChannels == nil {
		s.ExperimentalGroupUnreadChannels = NewPrimitive(GROUP_UNREAD_CHANNELS_DISABLED)
	} else if *s.ExperimentalGroupUnreadChannels == "0" {
		s.ExperimentalGroupUnreadChannels = NewPrimitive(GROUP_UNREAD_CHANNELS_DISABLED)
	} else if *s.ExperimentalGroupUnreadChannels == "1" {
		s.ExperimentalGroupUnreadChannels = NewPrimitive(GROUP_UNREAD_CHANNELS_DEFAULT_ON)
	}

	if s.ExperimentalChannelOrganization == nil {
		experimentalUnreadEnabled := *s.ExperimentalGroupUnreadChannels != GROUP_UNREAD_CHANNELS_DISABLED
		s.ExperimentalChannelOrganization = NewPrimitive(experimentalUnreadEnabled)
	}

	if s.DEPRECATED_DO_NOT_USE_ImageProxyType == nil {
		s.DEPRECATED_DO_NOT_USE_ImageProxyType = NewPrimitive("")
	}

	if s.DEPRECATED_DO_NOT_USE_ImageProxyURL == nil {
		s.DEPRECATED_DO_NOT_USE_ImageProxyURL = NewPrimitive("")
	}

	if s.DEPRECATED_DO_NOT_USE_ImageProxyOptions == nil {
		s.DEPRECATED_DO_NOT_USE_ImageProxyOptions = NewPrimitive("")
	}

	if s.EnableAPITeamDeletion == nil {
		s.EnableAPITeamDeletion = NewPrimitive(false)
	}

	if s.EnableAPIUserDeletion == nil {
		s.EnableAPIUserDeletion = NewPrimitive(false)
	}

	if s.EnableAPIChannelDeletion == nil {
		s.EnableAPIChannelDeletion = NewPrimitive(false)
	}

	if s.ExperimentalEnableHardenedMode == nil {
		s.ExperimentalEnableHardenedMode = NewPrimitive(false)
	}

	if s.DisableLegacyMFA == nil {
		s.DisableLegacyMFA = NewPrimitive(!isUpdate)
	}

	if s.ExperimentalStrictCSRFEnforcement == nil {
		s.ExperimentalStrictCSRFEnforcement = NewPrimitive(false)
	}

	if s.DisableBotsWhenOwnerIsDeactivated == nil {
		s.DisableBotsWhenOwnerIsDeactivated = NewPrimitive(true)
	}

	if s.EnableBotAccountCreation == nil {
		s.EnableBotAccountCreation = NewPrimitive(false)
	}

	if s.EnableSVGs == nil {
		if isUpdate {
			s.EnableSVGs = NewPrimitive(true)
		} else {
			s.EnableSVGs = NewPrimitive(false)
		}
	}

	if s.EnableLatex == nil {
		if isUpdate {
			s.EnableLatex = NewPrimitive(true)
		} else {
			s.EnableLatex = NewPrimitive(false)
		}
	}

	if s.EnableLocalMode == nil {
		s.EnableLocalMode = NewPrimitive(false)
	}

	if s.LocalModeSocketLocation == nil {
		s.LocalModeSocketLocation = NewPrimitive(LOCAL_MODE_SOCKET_PATH)
	}

	if s.EnableAWSMetering == nil {
		s.EnableAWSMetering = NewPrimitive(false)
	}

	if s.SplitKey == nil {
		s.SplitKey = NewPrimitive("")
	}

	if s.FeatureFlagSyncIntervalSeconds == nil {
		s.FeatureFlagSyncIntervalSeconds = NewPrimitive(30)
	}

	if s.DebugSplit == nil {
		s.DebugSplit = NewPrimitive(false)
	}

	if s.ThreadAutoFollow == nil {
		s.ThreadAutoFollow = NewPrimitive(true)
	}

	if s.CollapsedThreads == nil {
		s.CollapsedThreads = NewPrimitive(COLLAPSED_THREADS_DISABLED)
	}

	if s.ManagedResourcePaths == nil {
		s.ManagedResourcePaths = NewPrimitive("")
	}

	if s.EnableLegacySidebar == nil {
		s.EnableLegacySidebar = NewPrimitive(false)
	}

	if s.EnableReliableWebSockets == nil {
		s.EnableReliableWebSockets = NewPrimitive(true)
	}

}

type ClusterSettings struct {
	Enable                             *bool   `access:"environment_high_availability,write_restrictable"`
	ClusterName                        *string `access:"environment_high_availability,write_restrictable,cloud_restrictable"` // telemetry: none
	OverrideHostname                   *string `access:"environment_high_availability,write_restrictable,cloud_restrictable"` // telemetry: none
	NetworkInterface                   *string `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	BindAddress                        *string `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	AdvertiseAddress                   *string `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	UseIpAddress                       *bool   `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	UseExperimentalGossip              *bool   `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	EnableGossipCompression            *bool   `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	EnableExperimentalGossipEncryption *bool   `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	ReadOnlyConfig                     *bool   `access:"environment_high_availability,write_restrictable,cloud_restrictable"`
	GossipPort                         *int    `access:"environment_high_availability,write_restrictable,cloud_restrictable"` // telemetry: none
	StreamingPort                      *int    `access:"environment_high_availability,write_restrictable,cloud_restrictable"` // telemetry: none
	MaxIdleConns                       *int    `access:"environment_high_availability,write_restrictable,cloud_restrictable"` // telemetry: none
	MaxIdleConnsPerHost                *int    `access:"environment_high_availability,write_restrictable,cloud_restrictable"` // telemetry: none
	IdleConnTimeoutMilliseconds        *int    `access:"environment_high_availability,write_restrictable,cloud_restrictable"` // telemetry: none
}

func (s *ClusterSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	if s.ClusterName == nil {
		s.ClusterName = NewPrimitive("")
	}

	if s.OverrideHostname == nil {
		s.OverrideHostname = NewPrimitive("")
	}

	if s.NetworkInterface == nil {
		s.NetworkInterface = NewPrimitive("")
	}

	if s.BindAddress == nil {
		s.BindAddress = NewPrimitive("")
	}

	if s.AdvertiseAddress == nil {
		s.AdvertiseAddress = NewPrimitive("")
	}

	if s.UseIpAddress == nil {
		s.UseIpAddress = NewPrimitive(true)
	}

	if s.UseExperimentalGossip == nil {
		s.UseExperimentalGossip = NewPrimitive(true)
	}

	if s.EnableExperimentalGossipEncryption == nil {
		s.EnableExperimentalGossipEncryption = NewPrimitive(false)
	}

	if s.EnableGossipCompression == nil {
		s.EnableGossipCompression = NewPrimitive(true)
	}

	if s.ReadOnlyConfig == nil {
		s.ReadOnlyConfig = NewPrimitive(true)
	}

	if s.GossipPort == nil {
		s.GossipPort = NewPrimitive(8074)
	}

	if s.StreamingPort == nil {
		s.StreamingPort = NewPrimitive(8075)
	}

	if s.MaxIdleConns == nil {
		s.MaxIdleConns = NewPrimitive(100)
	}

	if s.MaxIdleConnsPerHost == nil {
		s.MaxIdleConnsPerHost = NewPrimitive(128)
	}

	if s.IdleConnTimeoutMilliseconds == nil {
		s.IdleConnTimeoutMilliseconds = NewPrimitive(90000)
	}
}

type MetricsSettings struct {
	Enable           *bool   `access:"environment_performance_monitoring,write_restrictable,cloud_restrictable"`
	BlockProfileRate *int    `access:"environment_performance_monitoring,write_restrictable,cloud_restrictable"`
	ListenAddress    *string `access:"environment_performance_monitoring,write_restrictable,cloud_restrictable"` // telemetry: none
}

func (s *MetricsSettings) SetDefaults() {
	if s.ListenAddress == nil {
		s.ListenAddress = NewPrimitive(":8067")
	}

	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	if s.BlockProfileRate == nil {
		s.BlockProfileRate = NewPrimitive(0)
	}
}

type ExperimentalSettings struct {
	ClientSideCertEnable            *bool   `access:"experimental_features,cloud_restrictable"`
	ClientSideCertCheck             *string `access:"experimental_features,cloud_restrictable"`
	EnableClickToReply              *bool   `access:"experimental_features,write_restrictable,cloud_restrictable"`
	LinkMetadataTimeoutMilliseconds *int64  `access:"experimental_features,write_restrictable,cloud_restrictable"`
	RestrictSystemAdmin             *bool   `access:"experimental_features,write_restrictable"`
	UseNewSAMLLibrary               *bool   `access:"experimental_features,cloud_restrictable"`
	CloudUserLimit                  *int64  `access:"experimental_features,write_restrictable"`
	CloudBilling                    *bool   `access:"experimental_features,write_restrictable"`
	EnableSharedChannels            *bool   `access:"experimental_features"`
	EnableRemoteClusterService      *bool   `access:"experimental_features"`
}

func (s *ExperimentalSettings) SetDefaults() {
	if s.ClientSideCertEnable == nil {
		s.ClientSideCertEnable = NewPrimitive(false)
	}

	if s.ClientSideCertCheck == nil {
		s.ClientSideCertCheck = NewPrimitive(CLIENT_SIDE_CERT_CHECK_SECONDARY_AUTH)
	}

	if s.EnableClickToReply == nil {
		s.EnableClickToReply = NewPrimitive(false)
	}

	if s.LinkMetadataTimeoutMilliseconds == nil {
		s.LinkMetadataTimeoutMilliseconds = NewPrimitive[int64](EXPERIMENTAL_SETTINGS_DEFAULT_LINK_METADATA_TIMEOUT_MILLISECONDS)
	}

	if s.RestrictSystemAdmin == nil {
		s.RestrictSystemAdmin = NewPrimitive(false)
	}

	if s.CloudUserLimit == nil {
		// User limit 0 is treated as no limit
		s.CloudUserLimit = NewPrimitive[int64](0)
	}

	if s.CloudBilling == nil {
		s.CloudBilling = NewPrimitive(false)
	}

	if s.UseNewSAMLLibrary == nil {
		s.UseNewSAMLLibrary = NewPrimitive(false)
	}

	if s.EnableSharedChannels == nil {
		s.EnableSharedChannels = NewPrimitive(false)
	}

	if s.EnableRemoteClusterService == nil {
		s.EnableRemoteClusterService = NewPrimitive(false)
	}
}

type AnalyticsSettings struct {
	MaxUsersForStatistics *int `access:"write_restrictable,cloud_restrictable"`
}

func (s *AnalyticsSettings) SetDefaults() {
	if s.MaxUsersForStatistics == nil {
		s.MaxUsersForStatistics = NewPrimitive(ANALYTICS_SETTINGS_DEFAULT_MAX_USERS_FOR_STATISTICS)
	}
}

type SSOSettings struct {
	Enable            *bool   `access:"authentication_openid"`
	Secret            *string `access:"authentication_openid"` // telemetry: none
	Id                *string `access:"authentication_openid"` // telemetry: none
	Scope             *string `access:"authentication_openid"` // telemetry: none
	AuthEndpoint      *string `access:"authentication_openid"` // telemetry: none
	TokenEndpoint     *string `access:"authentication_openid"` // telemetry: none
	UserApiEndpoint   *string `access:"authentication_openid"` // telemetry: none
	DiscoveryEndpoint *string `access:"authentication_openid"` // telemetry: none
	ButtonText        *string `access:"authentication_openid"` // telemetry: none
	ButtonColor       *string `access:"authentication_openid"` // telemetry: none
}

func (s *SSOSettings) setDefaults(scope, authEndpoint, tokenEndpoint, userApiEndpoint, buttonColor string) {
	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	if s.Secret == nil {
		s.Secret = NewPrimitive("")
	}

	if s.Id == nil {
		s.Id = NewPrimitive("")
	}

	if s.Scope == nil {
		s.Scope = NewPrimitive(scope)
	}

	if s.DiscoveryEndpoint == nil {
		s.DiscoveryEndpoint = NewPrimitive("")
	}

	if s.AuthEndpoint == nil {
		s.AuthEndpoint = NewPrimitive(authEndpoint)
	}

	if s.TokenEndpoint == nil {
		s.TokenEndpoint = NewPrimitive(tokenEndpoint)
	}

	if s.UserApiEndpoint == nil {
		s.UserApiEndpoint = NewPrimitive(userApiEndpoint)
	}

	if s.ButtonText == nil {
		s.ButtonText = NewPrimitive("")
	}

	if s.ButtonColor == nil {
		s.ButtonColor = NewPrimitive(buttonColor)
	}
}

type ReplicaLagSettings struct {
	DataSource       *string `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
	QueryAbsoluteLag *string `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
	QueryTimeLag     *string `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
}

type SqlSettings struct {
	DriverName                        *string               `access:"environment_database,write_restrictable,cloud_restrictable"`
	DataSource                        *string               `access:"environment_database,write_restrictable,cloud_restrictable"` // telemetry: none
	DataSourceReplicas                []string              `access:"environment_database,write_restrictable,cloud_restrictable"`
	DataSourceSearchReplicas          []string              `access:"environment_database,write_restrictable,cloud_restrictable"`
	MaxIdleConns                      *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	ConnMaxLifetimeMilliseconds       *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	ConnMaxIdleTimeMilliseconds       *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	MaxOpenConns                      *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	Trace                             *bool                 `access:"environment_database,write_restrictable,cloud_restrictable"`
	AtRestEncryptKey                  *string               `access:"environment_database,write_restrictable,cloud_restrictable"` // telemetry: none
	QueryTimeout                      *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	DisableDatabaseSearch             *bool                 `access:"environment_database,write_restrictable,cloud_restrictable"`
	MigrationsStatementTimeoutSeconds *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	ReplicaLagSettings                []*ReplicaLagSettings `access:"environment_database,write_restrictable,cloud_restrictable"` // telemetry: none
}

func (s *SqlSettings) SetDefaults(isUpdate bool) {
	if s.DriverName == nil {
		s.DriverName = NewPrimitive(DATABASE_DRIVER_POSTGRES)
	}

	if s.DataSource == nil {
		s.DataSource = NewPrimitive(SQL_SETTINGS_DEFAULT_DATA_SOURCE)
	}

	if s.DataSourceReplicas == nil {
		s.DataSourceReplicas = []string{}
	}

	if s.DataSourceSearchReplicas == nil {
		s.DataSourceSearchReplicas = []string{}
	}

	if isUpdate {
		// When updating an existing configuration, ensure an encryption key has been specified.
		if s.AtRestEncryptKey == nil || *s.AtRestEncryptKey == "" {
			s.AtRestEncryptKey = NewPrimitive(NewRandomString(32))
		}
	} else {
		// When generating a blank configuration, leave this key empty to be generated on server start.
		s.AtRestEncryptKey = NewPrimitive("")
	}

	if s.MaxIdleConns == nil {
		s.MaxIdleConns = NewPrimitive(20)
	}

	if s.MaxOpenConns == nil {
		s.MaxOpenConns = NewPrimitive(300)
	}

	if s.ConnMaxLifetimeMilliseconds == nil {
		s.ConnMaxLifetimeMilliseconds = NewPrimitive(3600000)
	}

	if s.ConnMaxIdleTimeMilliseconds == nil {
		s.ConnMaxIdleTimeMilliseconds = NewPrimitive(300000)
	}

	if s.Trace == nil {
		s.Trace = NewPrimitive(false)
	}

	if s.QueryTimeout == nil {
		s.QueryTimeout = NewPrimitive(30)
	}

	if s.DisableDatabaseSearch == nil {
		s.DisableDatabaseSearch = NewPrimitive(false)
	}

	if s.MigrationsStatementTimeoutSeconds == nil {
		s.MigrationsStatementTimeoutSeconds = NewPrimitive(100000)
	}

	if s.ReplicaLagSettings == nil {
		s.ReplicaLagSettings = []*ReplicaLagSettings{}
	}
}

type LogSettings struct {
	EnableConsole          *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"`
	ConsoleLevel           *string `access:"environment_logging,write_restrictable,cloud_restrictable"`
	ConsoleJson            *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"`
	EnableColor            *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"` // telemetry: none
	EnableFile             *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"`
	FileLevel              *string `access:"environment_logging,write_restrictable,cloud_restrictable"`
	FileJson               *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"`
	FileLocation           *string `access:"environment_logging,write_restrictable,cloud_restrictable"`
	EnableWebhookDebugging *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"`
	EnableDiagnostics      *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"` // telemetry: none
	EnableSentry           *bool   `access:"environment_logging,write_restrictable,cloud_restrictable"` // telemetry: none
	AdvancedLoggingConfig  *string `access:"environment_logging,write_restrictable,cloud_restrictable"`
}

func (s *LogSettings) SetDefaults() {
	if s.EnableConsole == nil {
		s.EnableConsole = NewPrimitive(true)
	}

	if s.ConsoleLevel == nil {
		s.ConsoleLevel = NewPrimitive("DEBUG")
	}

	if s.EnableColor == nil {
		s.EnableColor = NewPrimitive(false)
	}

	if s.EnableFile == nil {
		s.EnableFile = NewPrimitive(true)
	}

	if s.FileLevel == nil {
		s.FileLevel = NewPrimitive("INFO")
	}

	if s.FileLocation == nil {
		s.FileLocation = NewPrimitive("")
	}

	if s.EnableWebhookDebugging == nil {
		s.EnableWebhookDebugging = NewPrimitive(true)
	}

	if s.EnableDiagnostics == nil {
		s.EnableDiagnostics = NewPrimitive(true)
	}

	if s.EnableSentry == nil {
		s.EnableSentry = NewPrimitive(*s.EnableDiagnostics)
	}

	if s.ConsoleJson == nil {
		s.ConsoleJson = NewPrimitive(true)
	}

	if s.FileJson == nil {
		s.FileJson = NewPrimitive(true)
	}

	if s.AdvancedLoggingConfig == nil {
		s.AdvancedLoggingConfig = NewPrimitive("")
	}
}

type ExperimentalAuditSettings struct {
	FileEnabled           *bool   `access:"experimental_features,write_restrictable,cloud_restrictable"`
	FileName              *string `access:"experimental_features,write_restrictable,cloud_restrictable"` // telemetry: none
	FileMaxSizeMB         *int    `access:"experimental_features,write_restrictable,cloud_restrictable"`
	FileMaxAgeDays        *int    `access:"experimental_features,write_restrictable,cloud_restrictable"`
	FileMaxBackups        *int    `access:"experimental_features,write_restrictable,cloud_restrictable"`
	FileCompress          *bool   `access:"experimental_features,write_restrictable,cloud_restrictable"`
	FileMaxQueueSize      *int    `access:"experimental_features,write_restrictable,cloud_restrictable"`
	AdvancedLoggingConfig *string `access:"experimental_features,write_restrictable,cloud_restrictable"`
}

func (s *ExperimentalAuditSettings) SetDefaults() {
	if s.FileEnabled == nil {
		s.FileEnabled = NewPrimitive(false)
	}

	if s.FileName == nil {
		s.FileName = NewPrimitive("")
	}

	if s.FileMaxSizeMB == nil {
		s.FileMaxSizeMB = NewPrimitive(100)
	}

	if s.FileMaxAgeDays == nil {
		s.FileMaxAgeDays = NewPrimitive(0) // no limit on age
	}

	if s.FileMaxBackups == nil { // no limit on number of backups
		s.FileMaxBackups = NewPrimitive(0)
	}

	if s.FileCompress == nil {
		s.FileCompress = NewPrimitive(false)
	}

	if s.FileMaxQueueSize == nil {
		s.FileMaxQueueSize = NewPrimitive(1000)
	}

	if s.AdvancedLoggingConfig == nil {
		s.AdvancedLoggingConfig = NewPrimitive("")
	}
}

type NotificationLogSettings struct {
	EnableConsole         *bool   `access:"write_restrictable,cloud_restrictable"`
	ConsoleLevel          *string `access:"write_restrictable,cloud_restrictable"`
	ConsoleJson           *bool   `access:"write_restrictable,cloud_restrictable"`
	EnableColor           *bool   `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	EnableFile            *bool   `access:"write_restrictable,cloud_restrictable"`
	FileLevel             *string `access:"write_restrictable,cloud_restrictable"`
	FileJson              *bool   `access:"write_restrictable,cloud_restrictable"`
	FileLocation          *string `access:"write_restrictable,cloud_restrictable"`
	AdvancedLoggingConfig *string `access:"write_restrictable,cloud_restrictable"`
}

func (s *NotificationLogSettings) SetDefaults() {
	if s.EnableConsole == nil {
		s.EnableConsole = NewPrimitive(true)
	}

	if s.ConsoleLevel == nil {
		s.ConsoleLevel = NewPrimitive("DEBUG")
	}

	if s.EnableFile == nil {
		s.EnableFile = NewPrimitive(true)
	}

	if s.FileLevel == nil {
		s.FileLevel = NewPrimitive("INFO")
	}

	if s.FileLocation == nil {
		s.FileLocation = NewPrimitive("")
	}

	if s.ConsoleJson == nil {
		s.ConsoleJson = NewPrimitive(true)
	}

	if s.EnableColor == nil {
		s.EnableColor = NewPrimitive(false)
	}

	if s.FileJson == nil {
		s.FileJson = NewPrimitive(true)
	}

	if s.AdvancedLoggingConfig == nil {
		s.AdvancedLoggingConfig = NewPrimitive("")
	}
}

type PasswordSettings struct {
	MinimumLength *int  `access:"authentication_password"`
	Lowercase     *bool `access:"authentication_password"`
	Number        *bool `access:"authentication_password"`
	Uppercase     *bool `access:"authentication_password"`
	Symbol        *bool `access:"authentication_password"`
}

func (s *PasswordSettings) SetDefaults() {
	if s.MinimumLength == nil {
		s.MinimumLength = NewPrimitive(10)
	}

	if s.Lowercase == nil {
		s.Lowercase = NewPrimitive(true)
	}

	if s.Number == nil {
		s.Number = NewPrimitive(true)
	}

	if s.Uppercase == nil {
		s.Uppercase = NewPrimitive(true)
	}

	if s.Symbol == nil {
		s.Symbol = NewPrimitive(true)
	}
}

type FileSettings struct {
	EnableFileAttachments   *bool   `access:"site_file_sharing_and_downloads,cloud_restrictable"`
	EnableMobileUpload      *bool   `access:"site_file_sharing_and_downloads,cloud_restrictable"`
	EnableMobileDownload    *bool   `access:"site_file_sharing_and_downloads,cloud_restrictable"`
	MaxFileSize             *int64  `access:"environment_file_storage,cloud_restrictable"`
	DriverName              *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"`
	Directory               *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"`
	MaxImageResolution      *int64  `access:"environment_file_storage,cloud_restrictable"`
	EnablePublicLink        *bool   `access:"site_public_links,cloud_restrictable"`
	ExtractContent          *bool   `access:"environment_file_storage,write_restrictable"`
	ArchiveRecursion        *bool   `access:"environment_file_storage,write_restrictable"`
	PublicLinkSalt          *string `access:"site_public_links,cloud_restrictable"`                           // telemetry: none
	InitialFont             *string `access:"environment_file_storage,cloud_restrictable"`                    // telemetry: none
	AmazonS3AccessKeyId     *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"` // telemetry: none
	AmazonS3SecretAccessKey *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"` // telemetry: none
	AmazonS3Bucket          *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"` // telemetry: none
	AmazonS3PathPrefix      *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"` // telemetry: none
	AmazonS3Region          *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"` // telemetry: none
	AmazonS3Endpoint        *string `access:"environment_file_storage,write_restrictable,cloud_restrictable"` // telemetry: none
	AmazonS3SSL             *bool   `access:"environment_file_storage,write_restrictable,cloud_restrictable"`
	AmazonS3SignV2          *bool   `access:"environment_file_storage,write_restrictable,cloud_restrictable"`
	AmazonS3SSE             *bool   `access:"environment_file_storage,write_restrictable,cloud_restrictable"`
	AmazonS3Trace           *bool   `access:"environment_file_storage,write_restrictable,cloud_restrictable"`
}

func (s *FileSettings) SetDefaults(isUpdate bool) {
	if s.EnableFileAttachments == nil {
		s.EnableFileAttachments = NewPrimitive(true)
	}

	if s.EnableMobileUpload == nil {
		s.EnableMobileUpload = NewPrimitive(true)
	}

	if s.MaxImageResolution == nil {
		s.MaxImageResolution = NewPrimitive[int64](7680 * 4320) // 8K, ~33MPX
	}

	if s.EnableMobileDownload == nil {
		s.EnableMobileDownload = NewPrimitive(true)
	}

	if s.MaxFileSize == nil {
		s.MaxFileSize = NewPrimitive[int64](MB * 100)
	}

	if s.DriverName == nil {
		s.DriverName = NewPrimitive(IMAGE_DRIVER_LOCAL)
	}

	if s.Directory == nil || *s.Directory == "" {
		s.Directory = NewPrimitive(FILE_SETTINGS_DEFAULT_DIRECTORY)
	}

	if s.EnablePublicLink == nil {
		s.EnablePublicLink = NewPrimitive(false)
	}

	if s.ExtractContent == nil {
		s.ExtractContent = NewPrimitive(true)
	}

	if s.ArchiveRecursion == nil {
		s.ArchiveRecursion = NewPrimitive(false)
	}

	if isUpdate {
		// When updating an existing configuration, ensure link salt has been specified.
		if s.PublicLinkSalt == nil || *s.PublicLinkSalt == "" {
			s.PublicLinkSalt = NewPrimitive(NewRandomString(32))
		}
	} else {
		// When generating a blank configuration, leave link salt empty to be generated on server start.
		s.PublicLinkSalt = NewPrimitive("")
	}

	if s.InitialFont == nil {
		// Defaults to "nunito-bold.ttf"
		s.InitialFont = NewPrimitive("nunito-bold.ttf")
	}

	if s.AmazonS3AccessKeyId == nil {
		s.AmazonS3AccessKeyId = NewPrimitive("")
	}

	if s.AmazonS3SecretAccessKey == nil {
		s.AmazonS3SecretAccessKey = NewPrimitive("")
	}

	if s.AmazonS3Bucket == nil {
		s.AmazonS3Bucket = NewPrimitive("")
	}

	if s.AmazonS3PathPrefix == nil {
		s.AmazonS3PathPrefix = NewPrimitive("")
	}

	if s.AmazonS3Region == nil {
		s.AmazonS3Region = NewPrimitive("")
	}

	if s.AmazonS3Endpoint == nil || *s.AmazonS3Endpoint == "" {
		// Defaults to "s3.amazonaws.com"
		s.AmazonS3Endpoint = NewPrimitive("s3.amazonaws.com")
	}

	if s.AmazonS3SSL == nil {
		s.AmazonS3SSL = NewPrimitive(true) // Secure by default.
	}

	if s.AmazonS3SignV2 == nil {
		s.AmazonS3SignV2 = new(bool)
		// Signature v2 is not enabled by default.
	}

	if s.AmazonS3SSE == nil {
		s.AmazonS3SSE = NewPrimitive(false) // Not Encrypted by default.
	}

	if s.AmazonS3Trace == nil {
		s.AmazonS3Trace = NewPrimitive(false)
	}
}

func (s *FileSettings) ToFileBackendSettings(enableComplianceFeature bool) filestore.FileBackendSettings {
	if *s.DriverName == IMAGE_DRIVER_LOCAL {
		return filestore.FileBackendSettings{
			DriverName: *s.DriverName,
			Directory:  *s.Directory,
		}
	}
	return filestore.FileBackendSettings{
		DriverName:              *s.DriverName,
		AmazonS3AccessKeyId:     *s.AmazonS3AccessKeyId,
		AmazonS3SecretAccessKey: *s.AmazonS3SecretAccessKey,
		AmazonS3Bucket:          *s.AmazonS3Bucket,
		AmazonS3PathPrefix:      *s.AmazonS3PathPrefix,
		AmazonS3Region:          *s.AmazonS3Region,
		AmazonS3Endpoint:        *s.AmazonS3Endpoint,
		AmazonS3SSL:             s.AmazonS3SSL == nil || *s.AmazonS3SSL,
		AmazonS3SignV2:          s.AmazonS3SignV2 != nil && *s.AmazonS3SignV2,
		AmazonS3SSE:             s.AmazonS3SSE != nil && *s.AmazonS3SSE && enableComplianceFeature,
		AmazonS3Trace:           s.AmazonS3Trace != nil && *s.AmazonS3Trace,
	}
}

type EmailSettings struct {
	EnableSignUpWithEmail             *bool   `access:"authentication_email"`
	EnableSignInWithEmail             *bool   `access:"authentication_email"`
	EnableSignInWithUsername          *bool   `access:"authentication_email"`
	SendEmailNotifications            *bool   `access:"site_notifications"`
	UseChannelInEmailNotifications    *bool   `access:"experimental_features"`
	RequireEmailVerification          *bool   `access:"authentication_email"`
	FeedbackName                      *string `access:"site_notifications"`
	FeedbackEmail                     *string `access:"site_notifications,cloud_restrictable"`
	ReplyToAddress                    *string `access:"site_notifications,cloud_restrictable"`
	FeedbackOrganization              *string `access:"site_notifications"`
	EnableSMTPAuth                    *bool   `access:"environment_smtp,write_restrictable,cloud_restrictable"`
	SMTPUsername                      *string `access:"environment_smtp,write_restrictable,cloud_restrictable"` // telemetry: none
	SMTPPassword                      *string `access:"environment_smtp,write_restrictable,cloud_restrictable"` // telemetry: none
	SMTPServer                        *string `access:"environment_smtp,write_restrictable,cloud_restrictable"` // telemetry: none
	SMTPPort                          *string `access:"environment_smtp,write_restrictable,cloud_restrictable"` // telemetry: none
	SMTPServerTimeout                 *int    `access:"cloud_restrictable"`
	ConnectionSecurity                *string `access:"environment_smtp,write_restrictable,cloud_restrictable"`
	SendPushNotifications             *bool   `access:"environment_push_notification_server"`
	PushNotificationServer            *string `access:"environment_push_notification_server"` // telemetry: none
	PushNotificationContents          *string `access:"site_notifications"`
	PushNotificationBuffer            *int    // telemetry: none
	EnableEmailBatching               *bool   `access:"site_notifications"`
	EmailBatchingBufferSize           *int    `access:"experimental_features"`
	EmailBatchingInterval             *int    `access:"experimental_features"`
	EnablePreviewModeBanner           *bool   `access:"site_notifications"`
	SkipServerCertificateVerification *bool   `access:"environment_smtp,write_restrictable,cloud_restrictable"`
	EmailNotificationContentsType     *string `access:"site_notifications"`
	LoginButtonColor                  *string `access:"experimental_features"`
	LoginButtonBorderColor            *string `access:"experimental_features"`
	LoginButtonTextColor              *string `access:"experimental_features"`
}

func (s *EmailSettings) SetDefaults(isUpdate bool) {
	if s.EnableSignUpWithEmail == nil {
		s.EnableSignUpWithEmail = NewPrimitive(true)
	}

	if s.EnableSignInWithEmail == nil {
		s.EnableSignInWithEmail = NewPrimitive(*s.EnableSignUpWithEmail)
	}

	if s.EnableSignInWithUsername == nil {
		s.EnableSignInWithUsername = NewPrimitive(true)
	}

	if s.SendEmailNotifications == nil {
		s.SendEmailNotifications = NewPrimitive(true)
	}

	if s.UseChannelInEmailNotifications == nil {
		s.UseChannelInEmailNotifications = NewPrimitive(false)
	}

	if s.RequireEmailVerification == nil {
		s.RequireEmailVerification = NewPrimitive(false)
	}

	if s.FeedbackName == nil {
		s.FeedbackName = NewPrimitive("")
	}

	if s.FeedbackEmail == nil {
		s.FeedbackEmail = NewPrimitive("test@example.com")
	}

	if s.ReplyToAddress == nil {
		s.ReplyToAddress = NewPrimitive("test@example.com")
	}

	if s.FeedbackOrganization == nil {
		s.FeedbackOrganization = NewPrimitive(EMAIL_SETTINGS_DEFAULT_FEEDBACK_ORGANIZATION)
	}

	if s.EnableSMTPAuth == nil {
		if s.ConnectionSecurity == nil || *s.ConnectionSecurity == CONN_SECURITY_NONE {
			s.EnableSMTPAuth = NewPrimitive(false)
		} else {
			s.EnableSMTPAuth = NewPrimitive(true)
		}
	}

	if s.SMTPUsername == nil {
		s.SMTPUsername = NewPrimitive("")
	}

	if s.SMTPPassword == nil {
		s.SMTPPassword = NewPrimitive("")
	}

	if s.SMTPServer == nil || *s.SMTPServer == "" {
		s.SMTPServer = NewPrimitive("localhost")
	}

	if s.SMTPPort == nil || *s.SMTPPort == "" {
		s.SMTPPort = NewPrimitive("10025")
	}

	if s.SMTPServerTimeout == nil || *s.SMTPServerTimeout == 0 {
		s.SMTPServerTimeout = NewPrimitive(10)
	}

	if s.ConnectionSecurity == nil || *s.ConnectionSecurity == CONN_SECURITY_PLAIN {
		s.ConnectionSecurity = NewPrimitive(CONN_SECURITY_NONE)
	}

	if s.SendPushNotifications == nil {
		s.SendPushNotifications = NewPrimitive(!isUpdate)
	}

	if s.PushNotificationServer == nil {
		if isUpdate {
			s.PushNotificationServer = NewPrimitive("")
		} else {
			s.PushNotificationServer = NewPrimitive(GENERIC_NOTIFICATION_SERVER)
		}
	}

	if s.PushNotificationContents == nil {
		s.PushNotificationContents = NewPrimitive(FULL_NOTIFICATION)
	}

	if s.PushNotificationBuffer == nil {
		s.PushNotificationBuffer = NewPrimitive(1000)
	}

	if s.EnableEmailBatching == nil {
		s.EnableEmailBatching = NewPrimitive(false)
	}

	if s.EmailBatchingBufferSize == nil {
		s.EmailBatchingBufferSize = NewPrimitive(EMAIL_BATCHING_BUFFER_SIZE)
	}

	if s.EmailBatchingInterval == nil {
		s.EmailBatchingInterval = NewPrimitive(EMAIL_BATCHING_INTERVAL)
	}

	if s.EnablePreviewModeBanner == nil {
		s.EnablePreviewModeBanner = NewPrimitive(true)
	}

	if s.EnableSMTPAuth == nil {
		if *s.ConnectionSecurity == CONN_SECURITY_NONE {
			s.EnableSMTPAuth = NewPrimitive(false)
		} else {
			s.EnableSMTPAuth = NewPrimitive(true)
		}
	}

	if *s.ConnectionSecurity == CONN_SECURITY_PLAIN {
		*s.ConnectionSecurity = CONN_SECURITY_NONE
	}

	if s.SkipServerCertificateVerification == nil {
		s.SkipServerCertificateVerification = NewPrimitive(false)
	}

	if s.EmailNotificationContentsType == nil {
		s.EmailNotificationContentsType = NewPrimitive(EMAIL_NOTIFICATION_CONTENTS_FULL)
	}

	if s.LoginButtonColor == nil {
		s.LoginButtonColor = NewPrimitive("#0000")
	}

	if s.LoginButtonBorderColor == nil {
		s.LoginButtonBorderColor = NewPrimitive("#2389D7")
	}

	if s.LoginButtonTextColor == nil {
		s.LoginButtonTextColor = NewPrimitive("#2389D7")
	}
}

type RateLimitSettings struct {
	Enable           *bool  `access:"environment_rate_limiting,write_restrictable,cloud_restrictable"`
	PerSec           *int   `access:"environment_rate_limiting,write_restrictable,cloud_restrictable"`
	MaxBurst         *int   `access:"environment_rate_limiting,write_restrictable,cloud_restrictable"`
	MemoryStoreSize  *int   `access:"environment_rate_limiting,write_restrictable,cloud_restrictable"`
	VaryByRemoteAddr *bool  `access:"environment_rate_limiting,write_restrictable,cloud_restrictable"`
	VaryByUser       *bool  `access:"environment_rate_limiting,write_restrictable,cloud_restrictable"`
	VaryByHeader     string `access:"environment_rate_limiting,write_restrictable,cloud_restrictable"`
}

func (s *RateLimitSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	if s.PerSec == nil {
		s.PerSec = NewPrimitive(10)
	}

	if s.MaxBurst == nil {
		s.MaxBurst = NewPrimitive(100)
	}

	if s.MemoryStoreSize == nil {
		s.MemoryStoreSize = NewPrimitive(10000)
	}

	if s.VaryByRemoteAddr == nil {
		s.VaryByRemoteAddr = NewPrimitive(true)
	}

	if s.VaryByUser == nil {
		s.VaryByUser = NewPrimitive(false)
	}
}

type PrivacySettings struct {
	ShowEmailAddress *bool `access:"site_users_and_teams"`
	ShowFullName     *bool `access:"site_users_and_teams"`
}

func (s *PrivacySettings) setDefaults() {
	if s.ShowEmailAddress == nil {
		s.ShowEmailAddress = NewPrimitive(true)
	}

	if s.ShowFullName == nil {
		s.ShowFullName = NewPrimitive(true)
	}
}

type SupportSettings struct {
	TermsOfServiceLink                     *string `access:"site_customization,write_restrictable,cloud_restrictable"`
	PrivacyPolicyLink                      *string `access:"site_customization,write_restrictable,cloud_restrictable"`
	AboutLink                              *string `access:"site_customization,write_restrictable,cloud_restrictable"`
	HelpLink                               *string `access:"site_customization,write_restrictable,cloud_restrictable"`
	ReportAProblemLink                     *string `access:"site_customization,write_restrictable,cloud_restrictable"`
	SupportEmail                           *string `access:"site_customization"`
	CustomTermsOfServiceEnabled            *bool   `access:"compliance_custom_terms_of_service"`
	CustomTermsOfServiceReAcceptancePeriod *int    `access:"compliance_custom_terms_of_service"`
	EnableAskCommunityLink                 *bool   `access:"site_customization"`
}

func (s *SupportSettings) SetDefaults() {
	if !IsSafeLink(s.TermsOfServiceLink) {
		*s.TermsOfServiceLink = SUPPORT_SETTINGS_DEFAULT_TERMS_OF_SERVICE_LINK
	}

	if s.TermsOfServiceLink == nil {
		s.TermsOfServiceLink = NewPrimitive(SUPPORT_SETTINGS_DEFAULT_TERMS_OF_SERVICE_LINK)
	}

	if !IsSafeLink(s.PrivacyPolicyLink) {
		*s.PrivacyPolicyLink = ""
	}

	if s.PrivacyPolicyLink == nil {
		s.PrivacyPolicyLink = NewPrimitive(SUPPORT_SETTINGS_DEFAULT_PRIVACY_POLICY_LINK)
	}

	if !IsSafeLink(s.AboutLink) {
		*s.AboutLink = ""
	}

	if s.AboutLink == nil {
		s.AboutLink = NewPrimitive(SUPPORT_SETTINGS_DEFAULT_ABOUT_LINK)
	}

	if !IsSafeLink(s.HelpLink) {
		*s.HelpLink = ""
	}

	if s.HelpLink == nil {
		s.HelpLink = NewPrimitive(SUPPORT_SETTINGS_DEFAULT_HELP_LINK)
	}

	if !IsSafeLink(s.ReportAProblemLink) {
		*s.ReportAProblemLink = ""
	}

	if s.ReportAProblemLink == nil {
		s.ReportAProblemLink = NewPrimitive(SUPPORT_SETTINGS_DEFAULT_REPORT_A_PROBLEM_LINK)
	}

	if s.SupportEmail == nil {
		s.SupportEmail = NewPrimitive(SUPPORT_SETTINGS_DEFAULT_SUPPORT_EMAIL)
	}

	if s.CustomTermsOfServiceEnabled == nil {
		s.CustomTermsOfServiceEnabled = NewPrimitive(false)
	}

	if s.CustomTermsOfServiceReAcceptancePeriod == nil {
		s.CustomTermsOfServiceReAcceptancePeriod = NewPrimitive(SUPPORT_SETTINGS_DEFAULT_RE_ACCEPTANCE_PERIOD)
	}

	if s.EnableAskCommunityLink == nil {
		s.EnableAskCommunityLink = NewPrimitive(true)
	}
}

type AnnouncementSettings struct {
	EnableBanner          *bool   `access:"site_announcement_banner"`
	BannerText            *string `access:"site_announcement_banner"` // telemetry: none
	BannerColor           *string `access:"site_announcement_banner"`
	BannerTextColor       *string `access:"site_announcement_banner"`
	AllowBannerDismissal  *bool   `access:"site_announcement_banner"`
	AdminNoticesEnabled   *bool   `access:"site_notices"`
	UserNoticesEnabled    *bool   `access:"site_notices"`
	NoticesURL            *string `access:"site_notices,write_restrictable"` // telemetry: none
	NoticesFetchFrequency *int    `access:"site_notices,write_restrictable"` // telemetry: none
	NoticesSkipCache      *bool   `access:"site_notices,write_restrictable"` // telemetry: none
}

func (s *AnnouncementSettings) SetDefaults() {
	if s.EnableBanner == nil {
		s.EnableBanner = NewPrimitive(false)
	}

	if s.BannerText == nil {
		s.BannerText = NewPrimitive("")
	}

	if s.BannerColor == nil {
		s.BannerColor = NewPrimitive(ANNOUNCEMENT_SETTINGS_DEFAULT_BANNER_COLOR)
	}

	if s.BannerTextColor == nil {
		s.BannerTextColor = NewPrimitive(ANNOUNCEMENT_SETTINGS_DEFAULT_BANNER_TEXT_COLOR)
	}

	if s.AllowBannerDismissal == nil {
		s.AllowBannerDismissal = NewPrimitive(true)
	}

	if s.AdminNoticesEnabled == nil {
		s.AdminNoticesEnabled = NewPrimitive(true)
	}

	if s.UserNoticesEnabled == nil {
		s.UserNoticesEnabled = NewPrimitive(true)
	}
	if s.NoticesURL == nil {
		s.NoticesURL = NewPrimitive(ANNOUNCEMENT_SETTINGS_DEFAULT_NOTICES_JSON_URL)
	}
	if s.NoticesSkipCache == nil {
		s.NoticesSkipCache = NewPrimitive(false)
	}
	if s.NoticesFetchFrequency == nil {
		s.NoticesFetchFrequency = NewPrimitive(ANNOUNCEMENT_SETTINGS_DEFAULT_NOTICES_FETCH_FREQUENCY_SECONDS)
	}

}

type ThemeSettings struct {
	EnableThemeSelection *bool   `access:"experimental_features"`
	DefaultTheme         *string `access:"experimental_features"`
	AllowCustomThemes    *bool   `access:"experimental_features"`
	AllowedThemes        []string
}

func (s *ThemeSettings) SetDefaults() {
	if s.EnableThemeSelection == nil {
		s.EnableThemeSelection = NewPrimitive(true)
	}

	if s.DefaultTheme == nil {
		s.DefaultTheme = NewPrimitive(THEME_SETTING_DEFAULT)
	}

	if s.AllowCustomThemes == nil {
		s.AllowCustomThemes = NewPrimitive(true)
	}

	if s.AllowedThemes == nil {
		s.AllowedThemes = []string{}
	}
}

type ClientRequirements struct {
	AndroidLatestVersion string `access:"write_restrictable,cloud_restrictable"`
	AndroidMinVersion    string `access:"write_restrictable,cloud_restrictable"`
	DesktopLatestVersion string `access:"write_restrictable,cloud_restrictable"`
	DesktopMinVersion    string `access:"write_restrictable,cloud_restrictable"`
	IosLatestVersion     string `access:"write_restrictable,cloud_restrictable"`
	IosMinVersion        string `access:"write_restrictable,cloud_restrictable"`
}

type LdapSettings struct {
	// Basic
	Enable             *bool   `access:"authentication_ldap"`
	EnableSync         *bool   `access:"authentication_ldap"`
	LdapServer         *string `access:"authentication_ldap"` // telemetry: none
	LdapPort           *int    `access:"authentication_ldap"` // telemetry: none
	ConnectionSecurity *string `access:"authentication_ldap"`
	BaseDN             *string `access:"authentication_ldap"` // telemetry: none
	BindUsername       *string `access:"authentication_ldap"` // telemetry: none
	BindPassword       *string `access:"authentication_ldap"` // telemetry: none

	// Filtering
	UserFilter        *string `access:"authentication_ldap"` // telemetry: none
	GroupFilter       *string `access:"authentication_ldap"`
	GuestFilter       *string `access:"authentication_ldap"`
	EnableAdminFilter *bool
	AdminFilter       *string

	// Group Mapping
	GroupDisplayNameAttribute *string `access:"authentication_ldap"`
	GroupIdAttribute          *string `access:"authentication_ldap"`

	// User Mapping
	FirstNameAttribute *string `access:"authentication_ldap"`
	LastNameAttribute  *string `access:"authentication_ldap"`
	EmailAttribute     *string `access:"authentication_ldap"`
	UsernameAttribute  *string `access:"authentication_ldap"`
	NicknameAttribute  *string `access:"authentication_ldap"`
	IdAttribute        *string `access:"authentication_ldap"`
	PositionAttribute  *string `access:"authentication_ldap"`
	LoginIdAttribute   *string `access:"authentication_ldap"`
	PictureAttribute   *string `access:"authentication_ldap"`

	// Synchronization
	SyncIntervalMinutes *int `access:"authentication_ldap"`

	// Advanced
	SkipCertificateVerification *bool   `access:"authentication_ldap"`
	PublicCertificateFile       *string `access:"authentication_ldap"`
	PrivateKeyFile              *string `access:"authentication_ldap"`
	QueryTimeout                *int    `access:"authentication_ldap"`
	MaxPageSize                 *int    `access:"authentication_ldap"`

	// Customization
	LoginFieldName *string `access:"authentication_ldap"`

	LoginButtonColor       *string `access:"experimental_features"`
	LoginButtonBorderColor *string `access:"experimental_features"`
	LoginButtonTextColor   *string `access:"experimental_features"`

	Trace *bool `access:"authentication_ldap"` // telemetry: none
}

func (s *LdapSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	// When unset should default to LDAP Enabled
	if s.EnableSync == nil {
		s.EnableSync = NewPrimitive(*s.Enable)
	}

	if s.EnableAdminFilter == nil {
		s.EnableAdminFilter = NewPrimitive(false)
	}

	if s.LdapServer == nil {
		s.LdapServer = NewPrimitive("")
	}

	if s.LdapPort == nil {
		s.LdapPort = NewPrimitive(389)
	}

	if s.ConnectionSecurity == nil {
		s.ConnectionSecurity = NewPrimitive("")
	}

	if s.PublicCertificateFile == nil {
		s.PublicCertificateFile = NewPrimitive("")
	}

	if s.PrivateKeyFile == nil {
		s.PrivateKeyFile = NewPrimitive("")
	}

	if s.BaseDN == nil {
		s.BaseDN = NewPrimitive("")
	}

	if s.BindUsername == nil {
		s.BindUsername = NewPrimitive("")
	}

	if s.BindPassword == nil {
		s.BindPassword = NewPrimitive("")
	}

	if s.UserFilter == nil {
		s.UserFilter = NewPrimitive("")
	}

	if s.GuestFilter == nil {
		s.GuestFilter = NewPrimitive("")
	}

	if s.AdminFilter == nil {
		s.AdminFilter = NewPrimitive("")
	}

	if s.GroupFilter == nil {
		s.GroupFilter = NewPrimitive("")
	}

	if s.GroupDisplayNameAttribute == nil {
		s.GroupDisplayNameAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_GROUP_DISPLAY_NAME_ATTRIBUTE)
	}

	if s.GroupIdAttribute == nil {
		s.GroupIdAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_GROUP_ID_ATTRIBUTE)
	}

	if s.FirstNameAttribute == nil {
		s.FirstNameAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE)
	}

	if s.LastNameAttribute == nil {
		s.LastNameAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE)
	}

	if s.EmailAttribute == nil {
		s.EmailAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE)
	}

	if s.UsernameAttribute == nil {
		s.UsernameAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE)
	}

	if s.NicknameAttribute == nil {
		s.NicknameAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE)
	}

	if s.IdAttribute == nil {
		s.IdAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_ID_ATTRIBUTE)
	}

	if s.PositionAttribute == nil {
		s.PositionAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_POSITION_ATTRIBUTE)
	}

	if s.PictureAttribute == nil {
		s.PictureAttribute = NewPrimitive(LDAP_SETTINGS_DEFAULT_PICTURE_ATTRIBUTE)
	}

	// For those upgrading to the version when LoginIdAttribute was added
	// they need IdAttribute == LoginIdAttribute not to break
	if s.LoginIdAttribute == nil {
		s.LoginIdAttribute = s.IdAttribute
	}

	if s.SyncIntervalMinutes == nil {
		s.SyncIntervalMinutes = NewPrimitive(60)
	}

	if s.SkipCertificateVerification == nil {
		s.SkipCertificateVerification = NewPrimitive(false)
	}

	if s.QueryTimeout == nil {
		s.QueryTimeout = NewPrimitive(60)
	}

	if s.MaxPageSize == nil {
		s.MaxPageSize = NewPrimitive(0)
	}

	if s.LoginFieldName == nil {
		s.LoginFieldName = NewPrimitive(LDAP_SETTINGS_DEFAULT_LOGIN_FIELD_NAME)
	}

	if s.LoginButtonColor == nil {
		s.LoginButtonColor = NewPrimitive("#0000")
	}

	if s.LoginButtonBorderColor == nil {
		s.LoginButtonBorderColor = NewPrimitive("#2389D7")
	}

	if s.LoginButtonTextColor == nil {
		s.LoginButtonTextColor = NewPrimitive("#2389D7")
	}

	if s.Trace == nil {
		s.Trace = NewPrimitive(false)
	}
}

type ComplianceSettings struct {
	Enable      *bool   `access:"compliance_compliance_monitoring"`
	Directory   *string `access:"compliance_compliance_monitoring"` // telemetry: none
	EnableDaily *bool   `access:"compliance_compliance_monitoring"`
}

func (s *ComplianceSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	if s.Directory == nil {
		s.Directory = NewPrimitive("./data/")
	}

	if s.EnableDaily == nil {
		s.EnableDaily = NewPrimitive(false)
	}
}

type LocalizationSettings struct {
	DefaultServerLocale *string `access:"site_localization"`
	DefaultClientLocale *string `access:"site_localization"`
	AvailableLocales    *string `access:"site_localization"`
	DefaultCountryCode  *string `access:"site_localization"` // added for sitename
}

func (s *LocalizationSettings) SetDefaults() {
	if s.DefaultServerLocale == nil {
		s.DefaultServerLocale = NewPrimitive(DEFAULT_LOCALE)
	}

	if s.DefaultClientLocale == nil {
		s.DefaultClientLocale = NewPrimitive(DEFAULT_LOCALE)
	}

	if s.AvailableLocales == nil {
		s.AvailableLocales = NewPrimitive("")
	}
	if s.DefaultCountryCode == nil {
		s.DefaultCountryCode = NewPrimitive(DEFAULT_COUNTRY)
	}
}

type SamlSettings struct {
	// Basic
	Enable                        *bool `access:"authentication_saml"`
	EnableSyncWithLdap            *bool `access:"authentication_saml"`
	EnableSyncWithLdapIncludeAuth *bool `access:"authentication_saml"`
	IgnoreGuestsLdapSync          *bool `access:"authentication_saml"`

	Verify      *bool `access:"authentication_saml"`
	Encrypt     *bool `access:"authentication_saml"`
	SignRequest *bool `access:"authentication_saml"`

	IdpUrl                      *string `access:"authentication_saml"` // telemetry: none
	IdpDescriptorUrl            *string `access:"authentication_saml"` // telemetry: none
	IdpMetadataUrl              *string `access:"authentication_saml"` // telemetry: none
	ServiceProviderIdentifier   *string `access:"authentication_saml"` // telemetry: none
	AssertionConsumerServiceURL *string `access:"authentication_saml"` // telemetry: none

	SignatureAlgorithm *string `access:"authentication_saml"`
	CanonicalAlgorithm *string `access:"authentication_saml"`

	ScopingIDPProviderId *string `access:"authentication_saml"`
	ScopingIDPName       *string `access:"authentication_saml"`

	IdpCertificateFile    *string `access:"authentication_saml"` // telemetry: none
	PublicCertificateFile *string `access:"authentication_saml"` // telemetry: none
	PrivateKeyFile        *string `access:"authentication_saml"` // telemetry: none

	// User Mapping
	IdAttribute          *string `access:"authentication_saml"`
	GuestAttribute       *string `access:"authentication_saml"`
	EnableAdminAttribute *bool
	AdminAttribute       *string
	FirstNameAttribute   *string `access:"authentication_saml"`
	LastNameAttribute    *string `access:"authentication_saml"`
	EmailAttribute       *string `access:"authentication_saml"`
	UsernameAttribute    *string `access:"authentication_saml"`
	NicknameAttribute    *string `access:"authentication_saml"`
	LocaleAttribute      *string `access:"authentication_saml"`
	PositionAttribute    *string `access:"authentication_saml"`

	LoginButtonText *string `access:"authentication_saml"`

	LoginButtonColor       *string `access:"experimental_features"`
	LoginButtonBorderColor *string `access:"experimental_features"`
	LoginButtonTextColor   *string `access:"experimental_features"`
}

func (s *SamlSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	if s.EnableSyncWithLdap == nil {
		s.EnableSyncWithLdap = NewPrimitive(false)
	}

	if s.EnableSyncWithLdapIncludeAuth == nil {
		s.EnableSyncWithLdapIncludeAuth = NewPrimitive(false)
	}

	if s.IgnoreGuestsLdapSync == nil {
		s.IgnoreGuestsLdapSync = NewPrimitive(false)
	}

	if s.EnableAdminAttribute == nil {
		s.EnableAdminAttribute = NewPrimitive(false)
	}

	if s.Verify == nil {
		s.Verify = NewPrimitive(true)
	}

	if s.Encrypt == nil {
		s.Encrypt = NewPrimitive(true)
	}

	if s.SignRequest == nil {
		s.SignRequest = NewPrimitive(false)
	}

	if s.SignatureAlgorithm == nil {
		s.SignatureAlgorithm = NewPrimitive(SAML_SETTINGS_DEFAULT_SIGNATURE_ALGORITHM)
	}

	if s.CanonicalAlgorithm == nil {
		s.CanonicalAlgorithm = NewPrimitive(SAML_SETTINGS_DEFAULT_CANONICAL_ALGORITHM)
	}

	if s.IdpUrl == nil {
		s.IdpUrl = NewPrimitive("")
	}

	if s.IdpDescriptorUrl == nil {
		s.IdpDescriptorUrl = NewPrimitive("")
	}

	if s.ServiceProviderIdentifier == nil {
		if s.IdpDescriptorUrl != nil {
			s.ServiceProviderIdentifier = NewPrimitive(*s.IdpDescriptorUrl)
		} else {
			s.ServiceProviderIdentifier = NewPrimitive("")
		}
	}

	if s.IdpMetadataUrl == nil {
		s.IdpMetadataUrl = NewPrimitive("")
	}

	if s.IdpCertificateFile == nil {
		s.IdpCertificateFile = NewPrimitive("")
	}

	if s.PublicCertificateFile == nil {
		s.PublicCertificateFile = NewPrimitive("")
	}

	if s.PrivateKeyFile == nil {
		s.PrivateKeyFile = NewPrimitive("")
	}

	if s.AssertionConsumerServiceURL == nil {
		s.AssertionConsumerServiceURL = NewPrimitive("")
	}

	if s.ScopingIDPProviderId == nil {
		s.ScopingIDPProviderId = NewPrimitive("")
	}

	if s.ScopingIDPName == nil {
		s.ScopingIDPName = NewPrimitive("")
	}

	if s.LoginButtonText == nil || *s.LoginButtonText == "" {
		s.LoginButtonText = NewPrimitive(USER_AUTH_SERVICE_SAML_TEXT)
	}

	if s.IdAttribute == nil {
		s.IdAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_ID_ATTRIBUTE)
	}

	if s.GuestAttribute == nil {
		s.GuestAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_GUEST_ATTRIBUTE)
	}
	if s.AdminAttribute == nil {
		s.AdminAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_ADMIN_ATTRIBUTE)
	}
	if s.FirstNameAttribute == nil {
		s.FirstNameAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE)
	}

	if s.LastNameAttribute == nil {
		s.LastNameAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE)
	}

	if s.EmailAttribute == nil {
		s.EmailAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE)
	}

	if s.UsernameAttribute == nil {
		s.UsernameAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE)
	}

	if s.NicknameAttribute == nil {
		s.NicknameAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE)
	}

	if s.PositionAttribute == nil {
		s.PositionAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_POSITION_ATTRIBUTE)
	}

	if s.LocaleAttribute == nil {
		s.LocaleAttribute = NewPrimitive(SAML_SETTINGS_DEFAULT_LOCALE_ATTRIBUTE)
	}

	if s.LoginButtonColor == nil {
		s.LoginButtonColor = NewPrimitive("#34a28b")
	}

	if s.LoginButtonBorderColor == nil {
		s.LoginButtonBorderColor = NewPrimitive("#2389D7")
	}

	if s.LoginButtonTextColor == nil {
		s.LoginButtonTextColor = NewPrimitive("#ffffff")
	}
}

type NativeAppSettings struct {
	AppCustomURLSchemes    []string `access:"site_customization,write_restrictable,cloud_restrictable"` // telemetry: none
	AppDownloadLink        *string  `access:"site_customization,write_restrictable,cloud_restrictable"`
	AndroidAppDownloadLink *string  `access:"site_customization,write_restrictable,cloud_restrictable"`
	IosAppDownloadLink     *string  `access:"site_customization,write_restrictable,cloud_restrictable"`
}

func (s *NativeAppSettings) SetDefaults() {
	if s.AppDownloadLink == nil {
		s.AppDownloadLink = NewPrimitive(NATIVEAPP_SETTINGS_DEFAULT_APP_DOWNLOAD_LINK)
	}

	if s.AndroidAppDownloadLink == nil {
		s.AndroidAppDownloadLink = NewPrimitive(NATIVEAPP_SETTINGS_DEFAULT_ANDROID_APP_DOWNLOAD_LINK)
	}

	if s.IosAppDownloadLink == nil {
		s.IosAppDownloadLink = NewPrimitive(NATIVEAPP_SETTINGS_DEFAULT_IOS_APP_DOWNLOAD_LINK)
	}

	if s.AppCustomURLSchemes == nil {
		s.AppCustomURLSchemes = GetDefaultAppCustomURLSchemes()
	}
}

type ElasticsearchSettings struct {
	ConnectionUrl                 *string `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	Username                      *string `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	Password                      *string `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	EnableIndexing                *bool   `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	EnableSearching               *bool   `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	EnableAutocomplete            *bool   `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	Sniff                         *bool   `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	PostIndexReplicas             *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	PostIndexShards               *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	ChannelIndexReplicas          *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	ChannelIndexShards            *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	UserIndexReplicas             *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	UserIndexShards               *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	AggregatePostsAfterDays       *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"` // telemetry: none
	PostsAggregatorJobStartTime   *string `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"` // telemetry: none
	IndexPrefix                   *string `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	LiveIndexingBatchSize         *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	BulkIndexingTimeWindowSeconds *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	RequestTimeoutSeconds         *int    `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	SkipTLSVerification           *bool   `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
	Trace                         *string `access:"environment_elasticsearch,write_restrictable,cloud_restrictable"`
}

func (s *ElasticsearchSettings) SetDefaults() {
	if s.ConnectionUrl == nil {
		s.ConnectionUrl = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_CONNECTION_URL)
	}

	if s.Username == nil {
		s.Username = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_USERNAME)
	}

	if s.Password == nil {
		s.Password = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_PASSWORD)
	}

	if s.EnableIndexing == nil {
		s.EnableIndexing = NewPrimitive(false)
	}

	if s.EnableSearching == nil {
		s.EnableSearching = NewPrimitive(false)
	}

	if s.EnableAutocomplete == nil {
		s.EnableAutocomplete = NewPrimitive(false)
	}

	if s.Sniff == nil {
		s.Sniff = NewPrimitive(true)
	}

	if s.PostIndexReplicas == nil {
		s.PostIndexReplicas = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_POST_INDEX_REPLICAS)
	}

	if s.PostIndexShards == nil {
		s.PostIndexShards = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_POST_INDEX_SHARDS)
	}

	if s.ChannelIndexReplicas == nil {
		s.ChannelIndexReplicas = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_CHANNEL_INDEX_REPLICAS)
	}

	if s.ChannelIndexShards == nil {
		s.ChannelIndexShards = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_CHANNEL_INDEX_SHARDS)
	}

	if s.UserIndexReplicas == nil {
		s.UserIndexReplicas = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_USER_INDEX_REPLICAS)
	}

	if s.UserIndexShards == nil {
		s.UserIndexShards = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_USER_INDEX_SHARDS)
	}

	if s.AggregatePostsAfterDays == nil {
		s.AggregatePostsAfterDays = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_AGGREGATE_POSTS_AFTER_DAYS)
	}

	if s.PostsAggregatorJobStartTime == nil {
		s.PostsAggregatorJobStartTime = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_POSTS_AGGREGATOR_JOB_START_TIME)
	}

	if s.IndexPrefix == nil {
		s.IndexPrefix = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_INDEX_PREFIX)
	}

	if s.LiveIndexingBatchSize == nil {
		s.LiveIndexingBatchSize = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_LIVE_INDEXING_BATCH_SIZE)
	}

	if s.BulkIndexingTimeWindowSeconds == nil {
		s.BulkIndexingTimeWindowSeconds = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_BULK_INDEXING_TIME_WINDOW_SECONDS)
	}

	if s.RequestTimeoutSeconds == nil {
		s.RequestTimeoutSeconds = NewPrimitive(ELASTICSEARCH_SETTINGS_DEFAULT_REQUEST_TIMEOUT_SECONDS)
	}

	if s.SkipTLSVerification == nil {
		s.SkipTLSVerification = NewPrimitive(false)
	}

	if s.Trace == nil {
		s.Trace = NewPrimitive("")
	}
}

type BleveSettings struct {
	IndexDir                      *string `access:"experimental_bleve"` // telemetry: none
	EnableIndexing                *bool   `access:"experimental_bleve"`
	EnableSearching               *bool   `access:"experimental_bleve"`
	EnableAutocomplete            *bool   `access:"experimental_bleve"`
	BulkIndexingTimeWindowSeconds *int    `access:"experimental_bleve"`
}

func (bs *BleveSettings) SetDefaults() {
	if bs.IndexDir == nil {
		bs.IndexDir = NewPrimitive(BLEVE_SETTINGS_DEFAULT_INDEX_DIR)
	}

	if bs.EnableIndexing == nil {
		bs.EnableIndexing = NewPrimitive(false)
	}

	if bs.EnableSearching == nil {
		bs.EnableSearching = NewPrimitive(false)
	}

	if bs.EnableAutocomplete == nil {
		bs.EnableAutocomplete = NewPrimitive(false)
	}

	if bs.BulkIndexingTimeWindowSeconds == nil {
		bs.BulkIndexingTimeWindowSeconds = NewPrimitive(BLEVE_SETTINGS_DEFAULT_BULK_INDEXING_TIME_WINDOW_SECONDS)
	}
}

type DataRetentionSettings struct {
	EnableMessageDeletion *bool   `access:"compliance_data_retention_policy"`
	EnableFileDeletion    *bool   `access:"compliance_data_retention_policy"`
	MessageRetentionDays  *int    `access:"compliance_data_retention_policy"`
	FileRetentionDays     *int    `access:"compliance_data_retention_policy"`
	DeletionJobStartTime  *string `access:"compliance_data_retention_policy"`
}

func (s *DataRetentionSettings) SetDefaults() {
	if s.EnableMessageDeletion == nil {
		s.EnableMessageDeletion = NewPrimitive(false)
	}

	if s.EnableFileDeletion == nil {
		s.EnableFileDeletion = NewPrimitive(false)
	}

	if s.MessageRetentionDays == nil {
		s.MessageRetentionDays = NewPrimitive(DATA_RETENTION_SETTINGS_DEFAULT_MESSAGE_RETENTION_DAYS)
	}

	if s.FileRetentionDays == nil {
		s.FileRetentionDays = NewPrimitive(DATA_RETENTION_SETTINGS_DEFAULT_FILE_RETENTION_DAYS)
	}

	if s.DeletionJobStartTime == nil {
		s.DeletionJobStartTime = NewPrimitive(DATA_RETENTION_SETTINGS_DEFAULT_DELETION_JOB_START_TIME)
	}
}

type JobSettings struct {
	RunJobs      *bool `access:"write_restrictable,cloud_restrictable"`
	RunScheduler *bool `access:"write_restrictable,cloud_restrictable"`
}

func (s *JobSettings) SetDefaults() {
	if s.RunJobs == nil {
		s.RunJobs = NewPrimitive(true)
	}

	if s.RunScheduler == nil {
		s.RunScheduler = NewPrimitive(true)
	}
}

type CloudSettings struct {
	CWSUrl    *string `access:"write_restrictable"`
	CWSAPIUrl *string `access:"write_restrictable"`
}

func (s *CloudSettings) SetDefaults() {
	if s.CWSUrl == nil {
		s.CWSUrl = NewPrimitive(CLOUD_SETTINGS_DEFAULT_CWS_URL)
	}
	if s.CWSAPIUrl == nil {
		s.CWSAPIUrl = NewPrimitive(CLOUD_SETTINGS_DEFAULT_CWS_API_URL)
	}
}

type PluginState struct {
	Enable bool
}

type PluginSettings struct {
	Enable                      *bool                             `access:"plugins,write_restrictable"`
	EnableUploads               *bool                             `access:"plugins,write_restrictable,cloud_restrictable"`
	AllowInsecureDownloadURL    *bool                             `access:"plugins,write_restrictable,cloud_restrictable"`
	EnableHealthCheck           *bool                             `access:"plugins,write_restrictable,cloud_restrictable"`
	Directory                   *string                           `access:"plugins,write_restrictable,cloud_restrictable"` // telemetry: none
	ClientDirectory             *string                           `access:"plugins,write_restrictable,cloud_restrictable"` // telemetry: none
	Plugins                     map[string]map[string]interface{} `access:"plugins"`                                       // telemetry: none
	PluginStates                map[string]*PluginState           `access:"plugins"`                                       // telemetry: none
	EnableMarketplace           *bool                             `access:"plugins,write_restrictable,cloud_restrictable"`
	EnableRemoteMarketplace     *bool                             `access:"plugins,write_restrictable,cloud_restrictable"`
	AutomaticPrepackagedPlugins *bool                             `access:"plugins,write_restrictable,cloud_restrictable"`
	RequirePluginSignature      *bool                             `access:"plugins,write_restrictable,cloud_restrictable"`
	MarketplaceUrl              *string                           `access:"plugins,write_restrictable,cloud_restrictable"`
	SignaturePublicKeyFiles     []string                          `access:"plugins,write_restrictable,cloud_restrictable"`
	ChimeraOAuthProxyUrl        *string                           `access:"plugins,write_restrictable,cloud_restrictable"`
}

func (s *PluginSettings) SetDefaults(ls LogSettings) {
	if s.Enable == nil {
		s.Enable = NewPrimitive(true)
	}

	if s.EnableUploads == nil {
		s.EnableUploads = NewPrimitive(false)
	}

	if s.AllowInsecureDownloadURL == nil {
		s.AllowInsecureDownloadURL = NewPrimitive(false)
	}

	if s.EnableHealthCheck == nil {
		s.EnableHealthCheck = NewPrimitive(true)
	}

	if s.Directory == nil || *s.Directory == "" {
		s.Directory = NewPrimitive(PLUGIN_SETTINGS_DEFAULT_DIRECTORY)
	}

	if s.ClientDirectory == nil || *s.ClientDirectory == "" {
		s.ClientDirectory = NewPrimitive(PLUGIN_SETTINGS_DEFAULT_CLIENT_DIRECTORY)
	}

	if s.Plugins == nil {
		s.Plugins = make(map[string]map[string]interface{})
	}

	if s.PluginStates == nil {
		s.PluginStates = make(map[string]*PluginState)
	}

	if s.PluginStates["com.mattermost.nps"] == nil {
		// Enable the NPS plugin by default if diagnostics are enabled
		s.PluginStates["com.mattermost.nps"] = &PluginState{Enable: ls.EnableDiagnostics == nil || *ls.EnableDiagnostics}
	}

	if s.PluginStates["com.mattermost.plugin-incident-management"] == nil && BuildEnterpriseReady == "true" {
		// Enable the incident management plugin by default
		s.PluginStates["com.mattermost.plugin-incident-management"] = &PluginState{Enable: true}
	}

	if s.PluginStates["com.mattermost.plugin-channel-export"] == nil && BuildEnterpriseReady == "true" {
		// Enable the channel export plugin by default
		s.PluginStates["com.mattermost.plugin-channel-export"] = &PluginState{Enable: true}
	}

	if s.EnableMarketplace == nil {
		s.EnableMarketplace = NewPrimitive(PLUGIN_SETTINGS_DEFAULT_ENABLE_MARKETPLACE)
	}

	if s.EnableRemoteMarketplace == nil {
		s.EnableRemoteMarketplace = NewPrimitive(true)
	}

	if s.AutomaticPrepackagedPlugins == nil {
		s.AutomaticPrepackagedPlugins = NewPrimitive(true)
	}

	if s.MarketplaceUrl == nil || *s.MarketplaceUrl == "" || *s.MarketplaceUrl == PLUGIN_SETTINGS_OLD_MARKETPLACE_URL {
		s.MarketplaceUrl = NewPrimitive(PLUGIN_SETTINGS_DEFAULT_MARKETPLACE_URL)
	}

	if s.RequirePluginSignature == nil {
		s.RequirePluginSignature = NewPrimitive(false)
	}

	if s.SignaturePublicKeyFiles == nil {
		s.SignaturePublicKeyFiles = []string{}
	}
	if s.ChimeraOAuthProxyUrl == nil {
		s.ChimeraOAuthProxyUrl = NewPrimitive("")
	}
}

type GlobalRelayMessageExportSettings struct {
	CustomerType      *string `access:"compliance_compliance_export"` // must be either A9 or A10, dictates SMTP server url
	SmtpUsername      *string `access:"compliance_compliance_export"`
	SmtpPassword      *string `access:"compliance_compliance_export"`
	EmailAddress      *string `access:"compliance_compliance_export"` // the address to send messages to
	SMTPServerTimeout *int    `access:"compliance_compliance_export"`
}

func (s *GlobalRelayMessageExportSettings) SetDefaults() {
	if s.CustomerType == nil {
		s.CustomerType = NewPrimitive(GLOBALRELAY_CUSTOMER_TYPE_A9)
	}
	if s.SmtpUsername == nil {
		s.SmtpUsername = NewPrimitive("")
	}
	if s.SmtpPassword == nil {
		s.SmtpPassword = NewPrimitive("")
	}
	if s.EmailAddress == nil {
		s.EmailAddress = NewPrimitive("")
	}
	if s.SMTPServerTimeout == nil || *s.SMTPServerTimeout == 0 {
		s.SMTPServerTimeout = NewPrimitive(1800)
	}
}

type MessageExportSettings struct {
	EnableExport          *bool   `access:"compliance_compliance_export"`
	ExportFormat          *string `access:"compliance_compliance_export"`
	DailyRunTime          *string `access:"compliance_compliance_export"`
	ExportFromTimestamp   *int64  `access:"compliance_compliance_export"`
	BatchSize             *int    `access:"compliance_compliance_export"`
	DownloadExportResults *bool   `access:"compliance_compliance_export"`

	// formatter-specific settings - these are only expected to be non-nil if ExportFormat is set to the associated format
	GlobalRelaySettings *GlobalRelayMessageExportSettings `access:"compliance_compliance_export"`
}

func (s *MessageExportSettings) SetDefaults() {
	if s.EnableExport == nil {
		s.EnableExport = NewPrimitive(false)
	}

	if s.DownloadExportResults == nil {
		s.DownloadExportResults = NewPrimitive(false)
	}

	if s.ExportFormat == nil {
		s.ExportFormat = NewPrimitive(COMPLIANCE_EXPORT_TYPE_ACTIANCE)
	}

	if s.DailyRunTime == nil {
		s.DailyRunTime = NewPrimitive("01:00")
	}

	if s.ExportFromTimestamp == nil {
		s.ExportFromTimestamp = NewPrimitive[int64](0)
	}

	if s.BatchSize == nil {
		s.BatchSize = NewPrimitive(10000)
	}

	if s.GlobalRelaySettings == nil {
		s.GlobalRelaySettings = &GlobalRelayMessageExportSettings{}
	}
	s.GlobalRelaySettings.SetDefaults()
}

type DisplaySettings struct {
	CustomUrlSchemes     []string `access:"site_customization"`
	ExperimentalTimezone *bool    `access:"experimental_features"`
}

func (s *DisplaySettings) SetDefaults() {
	if s.CustomUrlSchemes == nil {
		customUrlSchemes := []string{}
		s.CustomUrlSchemes = customUrlSchemes
	}

	if s.ExperimentalTimezone == nil {
		s.ExperimentalTimezone = NewPrimitive(true)
	}
}

type GuestAccountsSettings struct {
	Enable                           *bool   `access:"authentication_guest_access"`
	AllowEmailAccounts               *bool   `access:"authentication_guest_access"`
	EnforceMultifactorAuthentication *bool   `access:"authentication_guest_access"`
	RestrictCreationToDomains        *string `access:"authentication_guest_access"`
}

func (s *GuestAccountsSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewPrimitive(false)
	}

	if s.AllowEmailAccounts == nil {
		s.AllowEmailAccounts = NewPrimitive(true)
	}

	if s.EnforceMultifactorAuthentication == nil {
		s.EnforceMultifactorAuthentication = NewPrimitive(false)
	}

	if s.RestrictCreationToDomains == nil {
		s.RestrictCreationToDomains = NewPrimitive("")
	}
}

type ImageProxySettings struct {
	Enable                  *bool   `access:"environment_image_proxy"`
	ImageProxyType          *string `access:"environment_image_proxy"`
	RemoteImageProxyURL     *string `access:"environment_image_proxy"`
	RemoteImageProxyOptions *string `access:"environment_image_proxy"`
}

func (s *ImageProxySettings) SetDefaults(ss ServiceSettings) {
	if s.Enable == nil {
		if ss.DEPRECATED_DO_NOT_USE_ImageProxyType == nil || *ss.DEPRECATED_DO_NOT_USE_ImageProxyType == "" {
			s.Enable = NewPrimitive(false)
		} else {
			s.Enable = NewPrimitive(true)
		}
	}

	if s.ImageProxyType == nil {
		if ss.DEPRECATED_DO_NOT_USE_ImageProxyType == nil || *ss.DEPRECATED_DO_NOT_USE_ImageProxyType == "" {
			s.ImageProxyType = NewPrimitive(IMAGE_PROXY_TYPE_LOCAL)
		} else {
			s.ImageProxyType = ss.DEPRECATED_DO_NOT_USE_ImageProxyType
		}
	}

	if s.RemoteImageProxyURL == nil {
		if ss.DEPRECATED_DO_NOT_USE_ImageProxyURL == nil {
			s.RemoteImageProxyURL = NewPrimitive("")
		} else {
			s.RemoteImageProxyURL = ss.DEPRECATED_DO_NOT_USE_ImageProxyURL
		}
	}

	if s.RemoteImageProxyOptions == nil {
		if ss.DEPRECATED_DO_NOT_USE_ImageProxyOptions == nil {
			s.RemoteImageProxyOptions = NewPrimitive("")
		} else {
			s.RemoteImageProxyOptions = ss.DEPRECATED_DO_NOT_USE_ImageProxyOptions
		}
	}
}

// ImportSettings defines configuration settings for file imports.
type ImportSettings struct {
	// The directory where to store the imported files.
	Directory *string
	// The number of days to retain the imported files before deleting them.
	RetentionDays *int
}

func (s *ImportSettings) isValid() *AppError {
	if *s.Directory == "" {
		return NewAppError("Config.IsValid", "model.config.is_valid.import.directory.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.RetentionDays <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.import.retention_days_too_low.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

// SetDefaults applies the default settings to the struct.
func (s *ImportSettings) SetDefaults() {
	if s.Directory == nil || *s.Directory == "" {
		s.Directory = NewPrimitive(IMPORT_SETTINGS_DEFAULT_DIRECTORY)
	}

	if s.RetentionDays == nil {
		s.RetentionDays = NewPrimitive(IMPORT_SETTINGS_DEFAULT_RETENTION_DAYS)
	}
}

// ExportSettings defines configuration settings for file exports.
type ExportSettings struct {
	// The directory where to store the exported files.
	Directory *string // telemetry: none
	// The number of days to retain the exported files before deleting them.
	RetentionDays *int
}

// func (s *ExportSettings) isValid() *AppError {
// 	if *s.Directory == "" {
// 		return NewAppError("Config.IsValid", "model.config.is_valid.export.directory.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	if *s.RetentionDays <= 0 {
// 		return NewAppError("Config.IsValid", "model.config.is_valid.export.retention_days_too_low.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	return nil
// }

// SetDefaults applies the default settings to the struct.
func (s *ExportSettings) SetDefaults() {
	if s.Directory == nil || *s.Directory == "" {
		s.Directory = NewPrimitive(EXPORT_SETTINGS_DEFAULT_DIRECTORY)
	}

	if s.RetentionDays == nil {
		s.RetentionDays = NewPrimitive(EXPORT_SETTINGS_DEFAULT_RETENTION_DAYS)
	}
}

type ConfigFunc func() *Config

const ConfigAccessTagType = "access"
const ConfigAccessTagWriteRestrictable = "write_restrictable"
const ConfigAccessTagCloudRestrictable = "cloud_restrictable"

// Allows read access if any PERMISSION_SYSCONSOLE_READ_* is allowed
const ConfigAccessTagAnySysConsoleRead = "*_read"

// Config fields support the 'access' tag with the following values corresponding to the suffix of the associated
// PERMISSION_SYSCONSOLE_*_* permission Id: 'about', 'reporting', 'user_management_users',
// 'user_management_groups', 'user_management_teams', 'user_management_channels',
// 'user_management_permissions', 'environment_web_server', 'environment_database', 'environment_elasticsearch',
// 'environment_file_storage', 'environment_image_proxy', 'environment_smtp', 'environment_push_notification_server',
// 'environment_high_availability', 'environment_rate_limiting', 'environment_logging', 'environment_session_lengths',
// 'environment_performance_monitoring', 'environment_developer', 'site', 'authentication', 'plugins',
// 'integrations', 'compliance', 'plugins', and 'experimental'. They grant read and/or write access to the config field
// to roles without PERMISSION_MANAGE_SYSTEM.
//
// The 'access' tag '*_read' checks for any SYSCONSOLE read permission and grants access if any read permission is allowed.
//
// By default config values can be written with PERMISSION_MANAGE_SYSTEM, but if ExperimentalSettings.RestrictSystemAdmin is true
// and the access tag contains the value 'write_restrictable', then even PERMISSION_MANAGE_SYSTEM does not grant write access.
//
// PERMISSION_MANAGE_SYSTEM always grants read access.
//
// Config values with the access tag 'cloud_restrictable' mean that are marked to be filtered when it's used in a cloud licensed
// environment with ExperimentalSettings.RestrictedSystemAdmin set to true.
//
// Example:
//
//	type HairSettings struct {
//	    // Colour is writeable with either PERMISSION_SYSCONSOLE_WRITE_REPORTING or PERMISSION_SYSCONSOLE_WRITE_USER_MANAGEMENT_GROUPS.
//	    // It is readable by PERMISSION_SYSCONSOLE_READ_REPORTING and PERMISSION_SYSCONSOLE_READ_USER_MANAGEMENT_GROUPS permissions.
//	    // PERMISSION_MANAGE_SYSTEM grants read and write access.
//	    Colour string `access:"reporting,user_management_groups"`
//
//
//	    // Length is only readable and writable via PERMISSION_MANAGE_SYSTEM.
//	    Length string
//
//	    // Product is only writeable by PERMISSION_MANAGE_SYSTEM if ExperimentalSettings.RestrictSystemAdmin is false.
//	    // PERMISSION_MANAGE_SYSTEM can always read the value.
//	    Product bool `access:write_restrictable`
//	}
type Config struct {
	ServiceSettings           ServiceSettings
	ClientRequirements        ClientRequirements
	SqlSettings               SqlSettings
	LogSettings               LogSettings
	ExperimentalAuditSettings ExperimentalAuditSettings
	NotificationLogSettings   NotificationLogSettings
	PasswordSettings          PasswordSettings
	FileSettings              FileSettings
	EmailSettings             EmailSettings
	RateLimitSettings         RateLimitSettings
	PrivacySettings           PrivacySettings
	SupportSettings           SupportSettings
	AnnouncementSettings      AnnouncementSettings
	ThemeSettings             ThemeSettings
	GitLabSettings            SSOSettings
	GoogleSettings            SSOSettings
	OpenIdSettings            SSOSettings
	LdapSettings              LdapSettings
	ComplianceSettings        ComplianceSettings
	LocalizationSettings      LocalizationSettings
	SamlSettings              SamlSettings
	NativeAppSettings         NativeAppSettings
	ClusterSettings           ClusterSettings
	MetricsSettings           MetricsSettings
	ExperimentalSettings      ExperimentalSettings
	AnalyticsSettings         AnalyticsSettings
	ElasticsearchSettings     ElasticsearchSettings
	BleveSettings             BleveSettings
	DataRetentionSettings     DataRetentionSettings
	MessageExportSettings     MessageExportSettings
	JobSettings               JobSettings // telemetry: none
	PluginSettings            PluginSettings
	DisplaySettings           DisplaySettings
	GuestAccountsSettings     GuestAccountsSettings
	ImageProxySettings        ImageProxySettings
	CloudSettings             CloudSettings  // telemetry: none
	FeatureFlags              *FeatureFlags  `access:"*_read" json:",omitempty"` // telemetry: none
	ImportSettings            ImportSettings // telemetry: none
	ExportSettings            ExportSettings
	ThirdPartySettings        ThirdPartySettings
}

func (o *Config) Clone() *Config {
	var ret Config
	if err := json.Unmarshal([]byte(o.ToJSON()), &ret); err != nil {
		panic(err)
	}
	return &ret
}

func (o *Config) ToJSON() string {
	return ModelToJson(o)
}

func (o *Config) ToJsonFiltered(tagType, tagValue string) string {
	filteredConfigMap := structToMapFilteredByTag(*o, tagType, tagValue)
	for key, value := range filteredConfigMap {
		v, ok := value.(map[string]interface{})
		if ok && len(v) == 0 {
			delete(filteredConfigMap, key)
		}
	}
	return ModelToJson(filteredConfigMap)
}

func (o *Config) GetSSOService(service string) *SSOSettings {
	switch service {
	case SERVICE_GITLAB:
		return &o.GitLabSettings
	case SERVICE_GOOGLE:
		return &o.GoogleSettings
	case SERVICE_OPENID:
		return &o.OpenIdSettings
	}

	return nil
}

func ConfigFromJson(data io.Reader) *Config {
	var o *Config
	json.NewDecoder(data).Decode(&o)
	return o
}

// isUpdate detects a pre-existing config based on whether SiteURL has been changed
func (o *Config) isUpdate() bool {
	return o.ServiceSettings.SiteURL != nil
}

func (o *Config) SetDefaults() {
	isUpdate := o.isUpdate()

	o.LdapSettings.SetDefaults()
	o.SamlSettings.SetDefaults()

	o.SqlSettings.SetDefaults(isUpdate)
	o.FileSettings.SetDefaults(isUpdate)
	o.EmailSettings.SetDefaults(isUpdate)
	o.PrivacySettings.setDefaults()
	o.GitLabSettings.setDefaults("", "", "", "", "")
	o.GoogleSettings.setDefaults(GOOGLE_SETTINGS_DEFAULT_SCOPE, GOOGLE_SETTINGS_DEFAULT_AUTH_ENDPOINT, GOOGLE_SETTINGS_DEFAULT_TOKEN_ENDPOINT, GOOGLE_SETTINGS_DEFAULT_USER_API_ENDPOINT, "")
	o.OpenIdSettings.setDefaults(OPENID_SETTINGS_DEFAULT_SCOPE, "", "", "", "#145DBF")
	o.ServiceSettings.SetDefaults(isUpdate)
	o.PasswordSettings.SetDefaults()
	o.MetricsSettings.SetDefaults()
	o.ExperimentalSettings.SetDefaults()
	o.SupportSettings.SetDefaults()
	o.AnnouncementSettings.SetDefaults()
	o.ThemeSettings.SetDefaults()
	o.ClusterSettings.SetDefaults()
	o.PluginSettings.SetDefaults(o.LogSettings)
	o.AnalyticsSettings.SetDefaults()
	o.ComplianceSettings.SetDefaults()
	o.LocalizationSettings.SetDefaults()
	o.ElasticsearchSettings.SetDefaults()
	o.BleveSettings.SetDefaults()
	o.NativeAppSettings.SetDefaults()
	o.DataRetentionSettings.SetDefaults()
	o.RateLimitSettings.SetDefaults()
	o.LogSettings.SetDefaults()
	o.ExperimentalAuditSettings.SetDefaults()
	o.NotificationLogSettings.SetDefaults()
	o.JobSettings.SetDefaults()
	o.MessageExportSettings.SetDefaults()
	o.DisplaySettings.SetDefaults()
	o.GuestAccountsSettings.SetDefaults()
	o.ImageProxySettings.SetDefaults(o.ServiceSettings)
	o.CloudSettings.SetDefaults()
	if o.FeatureFlags == nil {
		o.FeatureFlags = &FeatureFlags{}
		o.FeatureFlags.SetDefaults()
	}
	o.ImportSettings.SetDefaults()
	o.ExportSettings.SetDefaults()
	o.ThirdPartySettings.SetDefaults()
}

func (o *Config) IsValid() *AppError {
	if *o.ServiceSettings.SiteURL == "" && *o.EmailSettings.EnableEmailBatching {
		return NewAppError("Config.IsValid", "model.config.is_valid.site_url_email_batching.app_error", nil, "", http.StatusBadRequest)
	}

	if *o.ClusterSettings.Enable && *o.EmailSettings.EnableEmailBatching {
		return NewAppError("Config.IsValid", "model.config.is_valid.cluster_email_batching.app_error", nil, "", http.StatusBadRequest)
	}

	if *o.ServiceSettings.SiteURL == "" && *o.ServiceSettings.AllowCookiesForSubdomains {
		return NewAppError("Config.IsValid", "model.config.is_valid.allow_cookies_for_subdomains.app_error", nil, "", http.StatusBadRequest)
	}

	if err := o.SqlSettings.isValid(); err != nil {
		return err
	}

	if err := o.FileSettings.isValid(); err != nil {
		return err
	}

	if err := o.EmailSettings.isValid(); err != nil {
		return err
	}

	if err := o.LdapSettings.isValid(); err != nil {
		return err
	}

	if err := o.SamlSettings.isValid(); err != nil {
		return err
	}

	if *o.PasswordSettings.MinimumLength < PASSWORD_MINIMUM_LENGTH || *o.PasswordSettings.MinimumLength > PASSWORD_MAXIMUM_LENGTH {
		return NewAppError("Config.IsValid", "model.config.is_valid.password_length.app_error", map[string]interface{}{"MinLength": PASSWORD_MINIMUM_LENGTH, "MaxLength": PASSWORD_MAXIMUM_LENGTH}, "", http.StatusBadRequest)
	}

	if err := o.RateLimitSettings.isValid(); err != nil {
		return err
	}

	if err := o.ServiceSettings.isValid(); err != nil {
		return err
	}

	if err := o.ElasticsearchSettings.isValid(); err != nil {
		return err
	}

	if err := o.BleveSettings.isValid(); err != nil {
		return err
	}

	if err := o.DataRetentionSettings.isValid(); err != nil {
		return err
	}

	if err := o.LocalizationSettings.isValid(); err != nil {
		return err
	}

	if err := o.MessageExportSettings.isValid(); err != nil {
		return err
	}

	if err := o.DisplaySettings.isValid(); err != nil {
		return err
	}

	if err := o.ImageProxySettings.isValid(); err != nil {
		return err
	}

	if err := o.ImportSettings.isValid(); err != nil {
		return err
	}
	return nil
}

func (s *SqlSettings) isValid() *AppError {
	if *s.AtRestEncryptKey != "" && len(*s.AtRestEncryptKey) < 32 {
		return NewAppError("Config.IsValid", "model.config.is_valid.encrypt_sql.app_error", nil, "", http.StatusBadRequest)
	}

	if !(*s.DriverName == DATABASE_DRIVER_MYSQL || *s.DriverName == DATABASE_DRIVER_POSTGRES) {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_driver.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.MaxIdleConns <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_idle.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.ConnMaxLifetimeMilliseconds < 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_conn_max_lifetime_milliseconds.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.ConnMaxIdleTimeMilliseconds < 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_conn_max_idle_time_milliseconds.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.QueryTimeout <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_query_timeout.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.DataSource == "" {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_data_src.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.MaxOpenConns <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_max_conn.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (s *FileSettings) isValid() *AppError {
	if *s.MaxFileSize <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.max_file_size.app_error", nil, "", http.StatusBadRequest)
	}

	if !(*s.DriverName == IMAGE_DRIVER_LOCAL || *s.DriverName == IMAGE_DRIVER_S3) {
		return NewAppError("Config.IsValid", "model.config.is_valid.file_driver.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.PublicLinkSalt != "" && len(*s.PublicLinkSalt) < 32 {
		return NewAppError("Config.IsValid", "model.config.is_valid.file_salt.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.Directory == "" {
		return NewAppError("Config.IsValid", "model.config.is_valid.directory.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (s *EmailSettings) isValid() *AppError {
	if !(*s.ConnectionSecurity == CONN_SECURITY_NONE || *s.ConnectionSecurity == CONN_SECURITY_TLS || *s.ConnectionSecurity == CONN_SECURITY_STARTTLS || *s.ConnectionSecurity == CONN_SECURITY_PLAIN) {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_security.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.EmailBatchingBufferSize <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_batching_buffer_size.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.EmailBatchingInterval < 30 {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_batching_interval.app_error", nil, "", http.StatusBadRequest)
	}

	if !(*s.EmailNotificationContentsType == EMAIL_NOTIFICATION_CONTENTS_FULL || *s.EmailNotificationContentsType == EMAIL_NOTIFICATION_CONTENTS_GENERIC) {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_notification_contents_type.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (s *RateLimitSettings) isValid() *AppError {
	if *s.MemoryStoreSize <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.rate_mem.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.PerSec <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.rate_sec.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.MaxBurst <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.max_burst.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (s *LdapSettings) isValid() *AppError {
	if !(*s.ConnectionSecurity == CONN_SECURITY_NONE || *s.ConnectionSecurity == CONN_SECURITY_TLS || *s.ConnectionSecurity == CONN_SECURITY_STARTTLS) {
		return NewAppError("Config.IsValid", "model.config.is_valid.ldap_security.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.SyncIntervalMinutes <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.ldap_sync_interval.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.MaxPageSize < 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.ldap_max_page_size.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.Enable {
		if *s.LdapServer == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.ldap_server", nil, "", http.StatusBadRequest)
		}

		if *s.BaseDN == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.ldap_basedn", nil, "", http.StatusBadRequest)
		}

		if *s.EmailAttribute == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.ldap_email", nil, "", http.StatusBadRequest)
		}

		if *s.UsernameAttribute == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.ldap_username", nil, "", http.StatusBadRequest)
		}

		if *s.IdAttribute == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.ldap_id", nil, "", http.StatusBadRequest)
		}

		if *s.LoginIdAttribute == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.ldap_login_id", nil, "", http.StatusBadRequest)
		}

		if *s.UserFilter != "" {
			if _, err := ldap.CompileFilter(*s.UserFilter); err != nil {
				return NewAppError("ValidateFilter", "ent.ldap.validate_filter.app_error", nil, err.Error(), http.StatusBadRequest)
			}
		}

		if *s.GuestFilter != "" {
			if _, err := ldap.CompileFilter(*s.GuestFilter); err != nil {
				return NewAppError("LdapSettings.isValid", "ent.ldap.validate_guest_filter.app_error", nil, err.Error(), http.StatusBadRequest)
			}
		}

		if *s.AdminFilter != "" {
			if _, err := ldap.CompileFilter(*s.AdminFilter); err != nil {
				return NewAppError("LdapSettings.isValid", "ent.ldap.validate_admin_filter.app_error", nil, err.Error(), http.StatusBadRequest)
			}
		}
	}

	return nil
}

func (s *SamlSettings) isValid() *AppError {
	if *s.Enable {
		if *s.IdpUrl == "" || !IsValidHTTPURL(*s.IdpUrl) {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_idp_url.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.IdpDescriptorUrl == "" || !IsValidHTTPURL(*s.IdpDescriptorUrl) {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_idp_descriptor_url.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.IdpCertificateFile == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_idp_cert.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.EmailAttribute == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_email_attribute.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.UsernameAttribute == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_username_attribute.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.ServiceProviderIdentifier == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_spidentifier_attribute.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.Verify {
			if *s.AssertionConsumerServiceURL == "" || !IsValidHTTPURL(*s.AssertionConsumerServiceURL) {
				return NewAppError("Config.IsValid", "model.config.is_valid.saml_assertion_consumer_service_url.app_error", nil, "", http.StatusBadRequest)
			}
		}

		if *s.Encrypt {
			if *s.PrivateKeyFile == "" {
				return NewAppError("Config.IsValid", "model.config.is_valid.saml_private_key.app_error", nil, "", http.StatusBadRequest)
			}

			if *s.PublicCertificateFile == "" {
				return NewAppError("Config.IsValid", "model.config.is_valid.saml_public_cert.app_error", nil, "", http.StatusBadRequest)
			}
		}

		if *s.EmailAttribute == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_email_attribute.app_error", nil, "", http.StatusBadRequest)
		}

		if !(*s.SignatureAlgorithm == SAML_SETTINGS_SIGNATURE_ALGORITHM_SHA1 || *s.SignatureAlgorithm == SAML_SETTINGS_SIGNATURE_ALGORITHM_SHA256 || *s.SignatureAlgorithm == SAML_SETTINGS_SIGNATURE_ALGORITHM_SHA512) {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_signature_algorithm.app_error", nil, "", http.StatusBadRequest)
		}
		if !(*s.CanonicalAlgorithm == SAML_SETTINGS_CANONICAL_ALGORITHM_C14N || *s.CanonicalAlgorithm == SAML_SETTINGS_CANONICAL_ALGORITHM_C14N11) {
			return NewAppError("Config.IsValid", "model.config.is_valid.saml_canonical_algorithm.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.GuestAttribute != "" {
			if !(strings.Contains(*s.GuestAttribute, "=")) {
				return NewAppError("Config.IsValid", "model.config.is_valid.saml_guest_attribute.app_error", nil, "", http.StatusBadRequest)
			}
			if len(strings.Split(*s.GuestAttribute, "=")) != 2 {
				return NewAppError("Config.IsValid", "model.config.is_valid.saml_guest_attribute.app_error", nil, "", http.StatusBadRequest)
			}
		}

		if *s.AdminAttribute != "" {
			if !(strings.Contains(*s.AdminAttribute, "=")) {
				return NewAppError("Config.IsValid", "model.config.is_valid.saml_admin_attribute.app_error", nil, "", http.StatusBadRequest)
			}
			if len(strings.Split(*s.AdminAttribute, "=")) != 2 {
				return NewAppError("Config.IsValid", "model.config.is_valid.saml_admin_attribute.app_error", nil, "", http.StatusBadRequest)
			}
		}
	}

	return nil
}

func (s *ServiceSettings) isValid() *AppError {
	if !(*s.ConnectionSecurity == CONN_SECURITY_NONE || *s.ConnectionSecurity == CONN_SECURITY_TLS) {
		return NewAppError("Config.IsValid", "model.config.is_valid.webserver_security.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.ConnectionSecurity == CONN_SECURITY_TLS && !*s.UseLetsEncrypt {
		appErr := NewAppError("Config.IsValid", "model.config.is_valid.tls_cert_file_missing.app_error", nil, "", http.StatusBadRequest)

		if *s.TLSCertFile == "" {
			return appErr
		} else if _, err := os.Stat(*s.TLSCertFile); os.IsNotExist(err) {
			return appErr
		}

		appErr = NewAppError("Config.IsValid", "model.config.is_valid.tls_key_file_missing.app_error", nil, "", http.StatusBadRequest)

		if *s.TLSKeyFile == "" {
			return appErr
		} else if _, err := os.Stat(*s.TLSKeyFile); os.IsNotExist(err) {
			return appErr
		}
	}

	if len(s.TLSOverwriteCiphers) > 0 {
		for _, cipher := range s.TLSOverwriteCiphers {
			if _, ok := ServerTLSSupportedCiphers[cipher]; !ok {
				return NewAppError("Config.IsValid", "model.config.is_valid.tls_overwrite_cipher.app_error", map[string]interface{}{"name": cipher}, "", http.StatusBadRequest)
			}
		}
	}

	if *s.ReadTimeout <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.read_timeout.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.WriteTimeout <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.write_timeout.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.TimeBetweenUserTypingUpdatesMilliseconds < 1000 {
		return NewAppError("Config.IsValid", "model.config.is_valid.time_between_user_typing.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.MaximumLoginAttempts <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.login_attempts.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.SiteURL != "" {
		if _, err := url.ParseRequestURI(*s.SiteURL); err != nil {
			return NewAppError("Config.IsValid", "model.config.is_valid.site_url.app_error", nil, "", http.StatusBadRequest)
		}
	}

	if *s.WebsocketURL != "" {
		if _, err := url.ParseRequestURI(*s.WebsocketURL); err != nil {
			return NewAppError("Config.IsValid", "model.config.is_valid.websocket_url.app_error", nil, "", http.StatusBadRequest)
		}
	}

	host, port, _ := net.SplitHostPort(*s.ListenAddress)
	var isValidHost bool
	if host == "" {
		isValidHost = true
	} else {
		isValidHost = (net.ParseIP(host) != nil) || IsDomainName(host)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil || !isValidHost || portInt < 0 || portInt > math.MaxUint16 {
		return NewAppError("Config.IsValid", "model.config.is_valid.listen_address.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.ExperimentalGroupUnreadChannels != GROUP_UNREAD_CHANNELS_DISABLED &&
		*s.ExperimentalGroupUnreadChannels != GROUP_UNREAD_CHANNELS_DEFAULT_ON &&
		*s.ExperimentalGroupUnreadChannels != GROUP_UNREAD_CHANNELS_DEFAULT_OFF {
		return NewAppError("Config.IsValid", "model.config.is_valid.group_unread_channels.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.CollapsedThreads != COLLAPSED_THREADS_DISABLED &&
		*s.CollapsedThreads != COLLAPSED_THREADS_DEFAULT_ON &&
		*s.CollapsedThreads != COLLAPSED_THREADS_DEFAULT_OFF {
		return NewAppError("Config.IsValid", "model.config.is_valid.collapsed_threads.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (s *ElasticsearchSettings) isValid() *AppError {
	if *s.EnableIndexing {
		if *s.ConnectionUrl == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.connection_url.app_error", nil, "", http.StatusBadRequest)
		}
	}

	if *s.EnableSearching && !*s.EnableIndexing {
		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.enable_searching.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.EnableAutocomplete && !*s.EnableIndexing {
		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.enable_autocomplete.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.AggregatePostsAfterDays < 1 {
		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.aggregate_posts_after_days.app_error", nil, "", http.StatusBadRequest)
	}

	if _, err := time.Parse("15:04", *s.PostsAggregatorJobStartTime); err != nil {
		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.posts_aggregator_job_start_time.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	if *s.LiveIndexingBatchSize < 1 {
		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.live_indexing_batch_size.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.BulkIndexingTimeWindowSeconds < 1 {
		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.bulk_indexing_time_window_seconds.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.RequestTimeoutSeconds < 1 {
		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.request_timeout_seconds.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (bs *BleveSettings) isValid() *AppError {
	if *bs.EnableIndexing {
		if *bs.IndexDir == "" {
			return NewAppError("Config.IsValid", "model.config.is_valid.bleve_search.filename.app_error", nil, "", http.StatusBadRequest)
		}
	} else {
		if *bs.EnableSearching {
			return NewAppError("Config.IsValid", "model.config.is_valid.bleve_search.enable_searching.app_error", nil, "", http.StatusBadRequest)
		}
		if *bs.EnableAutocomplete {
			return NewAppError("Config.IsValid", "model.config.is_valid.bleve_search.enable_autocomplete.app_error", nil, "", http.StatusBadRequest)
		}
	}
	if *bs.BulkIndexingTimeWindowSeconds < 1 {
		return NewAppError("Config.IsValid", "model.config.is_valid.bleve_search.bulk_indexing_time_window_seconds.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (s *DataRetentionSettings) isValid() *AppError {
	if *s.MessageRetentionDays <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.data_retention.message_retention_days_too_low.app_error", nil, "", http.StatusBadRequest)
	}

	if *s.FileRetentionDays <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.data_retention.file_retention_days_too_low.app_error", nil, "", http.StatusBadRequest)
	}

	if _, err := time.Parse("15:04", *s.DeletionJobStartTime); err != nil {
		return NewAppError("Config.IsValid", "model.config.is_valid.data_retention.deletion_job_start_time.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	return nil
}

func (s *LocalizationSettings) isValid() *AppError {
	if *s.AvailableLocales != "" {
		if !strings.Contains(*s.AvailableLocales, *s.DefaultClientLocale) {
			return NewAppError("Config.IsValid", "model.config.is_valid.localization.available_locales.app_error", nil, "", http.StatusBadRequest)
		}
	}

	return nil
}

func (s *MessageExportSettings) isValid() *AppError {
	if s.EnableExport == nil {
		return NewAppError("Config.IsValid", "model.config.is_valid.message_export.enable.app_error", nil, "", http.StatusBadRequest)
	}
	if *s.EnableExport {
		if s.ExportFromTimestamp == nil || *s.ExportFromTimestamp < 0 || *s.ExportFromTimestamp > GetMillis() {
			return NewAppError("Config.IsValid", "model.config.is_valid.message_export.export_from.app_error", nil, "", http.StatusBadRequest)
		} else if s.DailyRunTime == nil {
			return NewAppError("Config.IsValid", "model.config.is_valid.message_export.daily_runtime.app_error", nil, "", http.StatusBadRequest)
		} else if _, err := time.Parse("15:04", *s.DailyRunTime); err != nil {
			return NewAppError("Config.IsValid", "model.config.is_valid.message_export.daily_runtime.app_error", nil, err.Error(), http.StatusBadRequest)
		} else if s.BatchSize == nil || *s.BatchSize < 0 {
			return NewAppError("Config.IsValid", "model.config.is_valid.message_export.batch_size.app_error", nil, "", http.StatusBadRequest)
		} else if s.ExportFormat == nil || (*s.ExportFormat != COMPLIANCE_EXPORT_TYPE_ACTIANCE && *s.ExportFormat != COMPLIANCE_EXPORT_TYPE_GLOBALRELAY && *s.ExportFormat != COMPLIANCE_EXPORT_TYPE_CSV) {
			return NewAppError("Config.IsValid", "model.config.is_valid.message_export.export_type.app_error", nil, "", http.StatusBadRequest)
		}

		if *s.ExportFormat == COMPLIANCE_EXPORT_TYPE_GLOBALRELAY {
			if s.GlobalRelaySettings == nil {
				return NewAppError("Config.IsValid", "model.config.is_valid.message_export.global_relay.config_missing.app_error", nil, "", http.StatusBadRequest)
			} else if s.GlobalRelaySettings.CustomerType == nil || (*s.GlobalRelaySettings.CustomerType != GLOBALRELAY_CUSTOMER_TYPE_A9 && *s.GlobalRelaySettings.CustomerType != GLOBALRELAY_CUSTOMER_TYPE_A10) {
				return NewAppError("Config.IsValid", "model.config.is_valid.message_export.global_relay.customer_type.app_error", nil, "", http.StatusBadRequest)
			} else if s.GlobalRelaySettings.EmailAddress == nil || !strings.Contains(*s.GlobalRelaySettings.EmailAddress, "@") {
				// validating email addresses is hard - just make sure it contains an '@' sign
				// see https://stackoverflow.com/questions/201323/using-a-regular-expression-to-validate-an-email-address
				return NewAppError("Config.IsValid", "model.config.is_valid.message_export.global_relay.email_address.app_error", nil, "", http.StatusBadRequest)
			} else if s.GlobalRelaySettings.SmtpUsername == nil || *s.GlobalRelaySettings.SmtpUsername == "" {
				return NewAppError("Config.IsValid", "model.config.is_valid.message_export.global_relay.smtp_username.app_error", nil, "", http.StatusBadRequest)
			} else if s.GlobalRelaySettings.SmtpPassword == nil || *s.GlobalRelaySettings.SmtpPassword == "" {
				return NewAppError("Config.IsValid", "model.config.is_valid.message_export.global_relay.smtp_password.app_error", nil, "", http.StatusBadRequest)
			}
		}
	}
	return nil
}

func (s *DisplaySettings) isValid() *AppError {
	if len(s.CustomUrlSchemes) != 0 {
		validProtocolPattern := regexp.MustCompile(`(?i)^\s*[A-Za-z][A-Za-z0-9.+-]*\s*$`)

		for _, scheme := range s.CustomUrlSchemes {
			if !validProtocolPattern.MatchString(scheme) {
				return NewAppError(
					"Config.IsValid",
					"model.config.is_valid.display.custom_url_schemes.app_error",
					map[string]interface{}{"Scheme": scheme},
					"",
					http.StatusBadRequest,
				)
			}
		}
	}

	return nil
}

func (s *ImageProxySettings) isValid() *AppError {
	if *s.Enable {
		switch *s.ImageProxyType {
		case IMAGE_PROXY_TYPE_LOCAL:
			// No other settings to validate
		case IMAGE_PROXY_TYPE_ATMOS_CAMO:
			if *s.RemoteImageProxyURL == "" {
				return NewAppError("Config.IsValid", "model.config.is_valid.atmos_camo_image_proxy_url.app_error", nil, "", http.StatusBadRequest)
			}

			if *s.RemoteImageProxyOptions == "" {
				return NewAppError("Config.IsValid", "model.config.is_valid.atmos_camo_image_proxy_options.app_error", nil, "", http.StatusBadRequest)
			}
		default:
			return NewAppError("Config.IsValid", "model.config.is_valid.image_proxy_type.app_error", nil, "", http.StatusBadRequest)
		}
	}

	return nil
}

// GetSanitizeOptions returns options for User type only
func (o *Config) GetSanitizeOptions() map[string]bool {
	options := map[string]bool{}
	options["fullname"] = *o.PrivacySettings.ShowFullName
	options["email"] = *o.PrivacySettings.ShowEmailAddress

	return options
}

func (o *Config) Sanitize() {
	if o.LdapSettings.BindPassword != nil && *o.LdapSettings.BindPassword != "" {
		*o.LdapSettings.BindPassword = FAKE_SETTING
	}

	*o.FileSettings.PublicLinkSalt = FAKE_SETTING

	if *o.FileSettings.AmazonS3SecretAccessKey != "" {
		*o.FileSettings.AmazonS3SecretAccessKey = FAKE_SETTING
	}

	if o.EmailSettings.SMTPPassword != nil && *o.EmailSettings.SMTPPassword != "" {
		*o.EmailSettings.SMTPPassword = FAKE_SETTING
	}

	if *o.GitLabSettings.Secret != "" {
		*o.GitLabSettings.Secret = FAKE_SETTING
	}

	if o.GoogleSettings.Secret != nil && *o.GoogleSettings.Secret != "" {
		*o.GoogleSettings.Secret = FAKE_SETTING
	}

	if o.OpenIdSettings.Secret != nil && *o.OpenIdSettings.Secret != "" {
		*o.OpenIdSettings.Secret = FAKE_SETTING
	}

	*o.SqlSettings.DataSource = FAKE_SETTING
	*o.SqlSettings.AtRestEncryptKey = FAKE_SETTING

	*o.ElasticsearchSettings.Password = FAKE_SETTING

	for i := range o.SqlSettings.DataSourceReplicas {
		o.SqlSettings.DataSourceReplicas[i] = FAKE_SETTING
	}

	for i := range o.SqlSettings.DataSourceSearchReplicas {
		o.SqlSettings.DataSourceSearchReplicas[i] = FAKE_SETTING
	}

	if o.MessageExportSettings.GlobalRelaySettings.SmtpPassword != nil && *o.MessageExportSettings.GlobalRelaySettings.SmtpPassword != "" {
		*o.MessageExportSettings.GlobalRelaySettings.SmtpPassword = FAKE_SETTING
	}

	if o.ServiceSettings.GfycatApiSecret != nil && *o.ServiceSettings.GfycatApiSecret != "" {
		*o.ServiceSettings.GfycatApiSecret = FAKE_SETTING
	}

	*o.ServiceSettings.SplitKey = FAKE_SETTING
}

// structToMapFilteredByTag converts a struct into a map removing those fields that has the tag passed
// as argument
func structToMapFilteredByTag(t interface{}, typeOfTag, filterTag string) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("Panicked in structToMapFilteredByTag. This should never happen.", slog.Any("recover", r))
		}
	}()

	val := reflect.ValueOf(t)
	elemField := reflect.TypeOf(t)

	if val.Kind() != reflect.Struct {
		return nil
	}

	out := map[string]interface{}{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		structField := elemField.Field(i)
		tagPermissions := strings.Split(structField.Tag.Get(typeOfTag), ",")
		if isTagPresent(filterTag, tagPermissions) {
			continue
		}

		var value interface{}

		switch field.Kind() {
		case reflect.Struct:
			value = structToMapFilteredByTag(field.Interface(), typeOfTag, filterTag)
		case reflect.Ptr:
			indirectType := field.Elem()
			if indirectType.Kind() == reflect.Struct {
				value = structToMapFilteredByTag(indirectType.Interface(), typeOfTag, filterTag)
			} else if indirectType.Kind() != reflect.Invalid {
				value = indirectType.Interface()
			}
		default:
			value = field.Interface()
		}

		out[val.Type().Field(i).Name] = value
	}

	return out
}

func isTagPresent(tag string, tags []string) bool {
	for _, val := range tags {
		tagValue := strings.TrimSpace(val)
		if tagValue != "" && tagValue == tag {
			return true
		}
	}

	return false
}
