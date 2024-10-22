package model

import (
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/util"
)

// scopes for permissions
const (
	PermissionScopeSystem = "system_scope"
	PermissionScopeShop   = "shop_scope"
)

// Permission type
type Permission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
}

var (
	PermissionInviteUser *Permission

	PermissionAssignSystemAdminRole *Permission
	PermissionManageRoles           *Permission
	PermissionEditOtherUsers        *Permission

	PermissionAddReaction    *Permission // every users can react to posts
	PermissionRemoveReaction *Permission // every users can react to posts

	PermissionUploadFile                   *Permission
	PermissionGetPublicLink                *Permission
	PermissionManageIncomingWebhooks       *Permission
	PermissionManageOutgoingWebhooks       *Permission
	PermissionManageOthersIncomingWebhooks *Permission
	PermissionManageOthersOutgoingWebhooks *Permission
	PermissionManageOAuth                  *Permission
	PermissionManageSystemWideOAuth        *Permission

	PermissionReadJobs              *Permission
	PermissionManageJobs            *Permission
	PermissionCreateUserAccessToken *Permission
	PermissionReadUserAccessToken   *Permission
	PermissionRevokeUserAccessToken *Permission

	PermissionViewMembers *Permission // view a shop's member

	PermissionManageRemoteClusters           *Permission
	PermissionDownloadComplianceExportResult *Permission
	PermissionCreateDataRetentionJob         *Permission
	PermissionReadDataRetentionJob           *Permission
	PermissionCreateComplianceExportJob      *Permission
	PermissionReadComplianceExportJob        *Permission
	PermissionReadAudits                     *Permission
	PermissionTestSiteUrl                    *Permission
	PermissionTestS3                         *Permission
	PermissionReloadConfig                   *Permission
	PermissionInvalidateCaches               *Permission
	PermissionRecycleDatabaseConnections     *Permission
	PermissionTestEmail                      *Permission

	PermissionTestElasticsearch                     *Permission
	PermissionPurgeElasticsearchIndexes             *Permission
	PermissionCreateElasticsearchPostIndexingJob    *Permission
	PermissionCreateElasticsearchPostAggregationJob *Permission
	PermissionReadElasticsearchPostIndexingJob      *Permission
	PermissionReadElasticsearchPostAggregationJob   *Permission

	PermissionPurgeBleveIndexes         *Permission
	PermissionCreatePostBleveIndexesJob *Permission

	PermissionCreateLdapSyncJob      *Permission
	PermissionReadLdapSyncJob        *Permission
	PermissionTestLdap               *Permission
	PermissionInvalidateEmailInvite  *Permission
	PermissionGetSamlMetadataFromIdp *Permission
	PermissionAddSamlPublicCert      *Permission
	PermissionAddSamlPrivateCert     *Permission
	PermissionAddSamlIdpCert         *Permission
	PermissionRemoveSamlPublicCert   *Permission
	PermissionRemoveSamlPrivateCert  *Permission
	PermissionRemoveSamlIdpCert      *Permission
	PermissionGetSamlCertStatus      *Permission
	PermissionAddLdapPublicCert      *Permission
	PermissionAddLdapPrivateCert     *Permission
	PermissionRemoveLdapPublicCert   *Permission
	PermissionRemoveLdapPrivateCert  *Permission
	PermissionGetLogs                *Permission
	PermissionGetAnalytics           *Permission

	PermissionSysconsoleReadBilling  *Permission
	PermissionSysconsoleWriteBilling *Permission

	PermissionSysconsoleReadReportingSiteStatistics  *Permission
	PermissionSysconsoleWriteReportingSiteStatistics *Permission

	PermissionSysconsoleReadReportingServerLogs  *Permission
	PermissionSysconsoleWriteReportingServerLogs *Permission

	PermissionSysconsoleReadUserManagementUsers  *Permission
	PermissionSysconsoleWriteUserManagementUsers *Permission

	PermissionSysconsoleReadUserManagementGroups  *Permission
	PermissionSysconsoleWriteUserManagementGroups *Permission

	PermissionSysconsoleReadUserManagementPermissions  *Permission
	PermissionSysconsoleWriteUserManagementPermissions *Permission

	PermissionSysconsoleReadUserManagementSystemRoles  *Permission
	PermissionSysconsoleWriteUserManagementSystemRoles *Permission

	PermissionSysconsoleReadEnvironmentWebServer  *Permission
	PermissionSysconsoleWriteEnvironmentWebServer *Permission

	PermissionSysconsoleReadEnvironmentDatabase  *Permission
	PermissionSysconsoleWriteEnvironmentDatabase *Permission

	PermissionSysconsoleReadEnvironmentElasticsearch  *Permission
	PermissionSysconsoleWriteEnvironmentElasticsearch *Permission

	PermissionSysconsoleReadEnvironmentFileStorage  *Permission
	PermissionSysconsoleWriteEnvironmentFileStorage *Permission

	PermissionSysconsoleReadEnvironmentImageProxy  *Permission
	PermissionSysconsoleWriteEnvironmentImageProxy *Permission

	PermissionSysconsoleReadEnvironmentSmtp  *Permission
	PermissionSysconsoleWriteEnvironmentSmtp *Permission

	PermissionSysconsoleReadEnvironmentPushNotificationServer  *Permission
	PermissionSysconsoleWriteEnvironmentPushNotificationServer *Permission

	PermissionSysconsoleReadEnvironmentHighAvailability  *Permission
	PermissionSysconsoleWriteEnvironmentHighAvailability *Permission

	PermissionSysconsoleReadEnvironmentRateLimiting  *Permission
	PermissionSysconsoleWriteEnvironmentRateLimiting *Permission

	PermissionSysconsoleReadEnvironmentLogging  *Permission
	PermissionSysconsoleWriteEnvironmentLogging *Permission

	PermissionSysconsoleReadEnvironmentSessionLengths  *Permission
	PermissionSysconsoleWriteEnvironmentSessionLengths *Permission

	PermissionSysconsoleReadEnvironmentPerformanceMonitoring  *Permission
	PermissionSysconsoleWriteEnvironmentPerformanceMonitoring *Permission

	PermissionSysconsoleReadEnvironmentDeveloper  *Permission
	PermissionSysconsoleWriteEnvironmentDeveloper *Permission

	PermissionSysconsoleReadSiteCustomization  *Permission
	PermissionSysconsoleWriteSiteCustomization *Permission

	PermissionSysconsoleReadSiteLocalization  *Permission
	PermissionSysconsoleWriteSiteLocalization *Permission

	PermissionSysconsoleReadSiteNotifications  *Permission
	PermissionSysconsoleWriteSiteNotifications *Permission

	PermissionSysconsoleReadSiteAnnouncementBanner  *Permission
	PermissionSysconsoleWriteSiteAnnouncementBanner *Permission

	PermissionSysconsoleReadSitePosts  *Permission
	PermissionSysconsoleWriteSitePosts *Permission

	PermissionSysconsoleReadSiteFileSharingAndDownloads  *Permission
	PermissionSysconsoleWriteSiteFileSharingAndDownloads *Permission

	PermissionSysconsoleReadSitePublicLinks  *Permission
	PermissionSysconsoleWriteSitePublicLinks *Permission

	PermissionSysconsoleReadSiteNotices  *Permission
	PermissionSysconsoleWriteSiteNotices *Permission

	PermissionSysconsoleReadAuthenticationSignup  *Permission
	PermissionSysconsoleWriteAuthenticationSignup *Permission

	PermissionSysconsoleReadAuthenticationEmail  *Permission
	PermissionSysconsoleWriteAuthenticationEmail *Permission

	PermissionSysconsoleReadAuthenticationPassword  *Permission
	PermissionSysconsoleWriteAuthenticationPassword *Permission

	PermissionSysconsoleReadAuthenticationMfa  *Permission
	PermissionSysconsoleWriteAuthenticationMfa *Permission

	PermissionSysconsoleReadAuthenticationLdap  *Permission
	PermissionSysconsoleWriteAuthenticationLdap *Permission

	PermissionSysconsoleReadAuthenticationSaml  *Permission
	PermissionSysconsoleWriteAuthenticationSaml *Permission

	PermissionSysconsoleReadAuthenticationOpenid  *Permission
	PermissionSysconsoleWriteAuthenticationOpenid *Permission

	PermissionSysconsoleReadAuthenticationGuestAccess  *Permission
	PermissionSysconsoleWriteAuthenticationGuestAccess *Permission

	PermissionSysconsoleReadPlugins  *Permission
	PermissionSysconsoleWritePlugins *Permission

	PermissionSysconsoleReadIntegrationsIntegrationManagement  *Permission
	PermissionSysconsoleWriteIntegrationsIntegrationManagement *Permission

	PermissionSysconsoleReadIntegrationsBotAccounts  *Permission
	PermissionSysconsoleWriteIntegrationsBotAccounts *Permission

	PermissionSysconsoleReadIntegrationsGif  *Permission
	PermissionSysconsoleWriteIntegrationsGif *Permission

	PermissionSysconsoleReadIntegrationsCors  *Permission
	PermissionSysconsoleWriteIntegrationsCors *Permission

	PermissionSysconsoleReadComplianceDataRetentionPolicy  *Permission
	PermissionSysconsoleWriteComplianceDataRetentionPolicy *Permission

	PermissionSysconsoleReadComplianceComplianceExport  *Permission
	PermissionSysconsoleWriteComplianceComplianceExport *Permission

	PermissionSysconsoleReadComplianceComplianceMonitoring  *Permission
	PermissionSysconsoleWriteComplianceComplianceMonitoring *Permission

	PermissionSysconsoleReadComplianceCustomTermsOfService  *Permission
	PermissionSysconsoleWriteComplianceCustomTermsOfService *Permission

	PermissionSysconsoleReadExperimentalFeatures  *Permission
	PermissionSysconsoleWriteExperimentalFeatures *Permission

	PermissionSysconsoleReadExperimentalFeatureFlags  *Permission
	PermissionSysconsoleWriteExperimentalFeatureFlags *Permission

	PermissionSysconsoleReadExperimentalBleve  *Permission
	PermissionSysconsoleWriteExperimentalBleve *Permission

	// General permission that encompasses all system admin functions
	// in the future this could be broken up to allow access to some
	// admin functions but not others
	PermissionManageSystem *Permission
)

type Permissions []*Permission

func (ps Permissions) IDs() util.AnyArray[string] {
	return lo.Map(ps, func(item *Permission, _ int) string { return item.Id })
}

// AllSystemScopedPermissions contains all system's permissions
var AllSystemScopedPermissions Permissions
var SysconsoleReadPermissions Permissions
var SysconsoleWritePermissions Permissions
var SystemScopedPermissionsMinusSysconsole Permissions

func initializeSystemScopedPermissions() {
	PermissionInviteUser = &Permission{"invite_user", "authentication.permissions.team_invite_user.name", "authentication.permissions.team_invite_user.description", PermissionScopeSystem}
	PermissionAssignSystemAdminRole = &Permission{"assign_system_admin_role", "authentication.permissions.assign_system_admin_role.name", "authentication.permissions.assign_system_admin_role.description", PermissionScopeSystem}
	PermissionManageRoles = &Permission{"manage_roles", "authentication.permissions.manage_roles.name", "authentication.permissions.manage_roles.description", PermissionScopeSystem}
	PermissionManageSystem = &Permission{"manage_system", "authentication.permissions.manage_system.name", "authentication.permissions.manage_system.description", PermissionScopeSystem}
	PermissionEditOtherUsers = &Permission{"edit_other_users", "authentication.permissions.edit_other_users.name", "authentication.permissions.edit_other_users.description", PermissionScopeSystem}

	PermissionGetPublicLink = &Permission{"get_public_link", "authentication.permissions.get_public_link.name", "authentication.permissions.get_public_link.description", PermissionScopeSystem}
	PermissionManageIncomingWebhooks = &Permission{"manage_incoming_webhooks", "authentication.permissions.manage_incoming_webhooks.name", "authentication.permissions.manage_incoming_webhooks.description", PermissionScopeSystem}
	PermissionManageOutgoingWebhooks = &Permission{"manage_outgoing_webhooks", "authentication.permissions.manage_outgoing_webhooks.name", "authentication.permissions.manage_outgoing_webhooks.description", PermissionScopeSystem}
	PermissionManageOthersIncomingWebhooks = &Permission{"manage_others_incoming_webhooks", "authentication.permissions.manage_others_incoming_webhooks.name", "authentication.permissions.manage_others_incoming_webhooks.description", PermissionScopeSystem}
	PermissionManageOthersOutgoingWebhooks = &Permission{"manage_others_outgoing_webhooks", "authentication.permissions.manage_others_outgoing_webhooks.name", "authentication.permissions.manage_others_outgoing_webhooks.description", PermissionScopeSystem}
	PermissionManageOAuth = &Permission{"manage_oauth", "authentication.permissions.manage_oauth.name", "authentication.permissions.manage_oauth.description", PermissionScopeSystem}
	PermissionManageSystemWideOAuth = &Permission{"manage_system_wide_oauth", "authentication.permissions.manage_system_wide_oauth.name", "authentication.permissions.manage_system_wide_oauth.description", PermissionScopeSystem}
	PermissionManageRemoteClusters = &Permission{"manage_remote_clusters", "authentication.permissions.manage_remote_clusters.name", "authentication.permissions.manage_remote_clusters.description", PermissionScopeSystem}
	PermissionCreateDataRetentionJob = &Permission{"create_data_retention_job", "", "", PermissionScopeSystem}
	PermissionReadDataRetentionJob = &Permission{"read_data_retention_job", "", "", PermissionScopeSystem}
	PermissionCreateComplianceExportJob = &Permission{"create_compliance_export_job", "", "", PermissionScopeSystem}
	PermissionReadComplianceExportJob = &Permission{"read_compliance_export_job", "", "", PermissionScopeSystem}
	PermissionReadAudits = &Permission{"read_audits", "", "", PermissionScopeSystem}
	PermissionPurgeBleveIndexes = &Permission{"purge_bleve_indexes", "", "", PermissionScopeSystem}
	PermissionCreatePostBleveIndexesJob = &Permission{"create_post_bleve_indexes_job", "", "", PermissionScopeSystem}
	PermissionCreateLdapSyncJob = &Permission{"create_ldap_sync_job", "", "", PermissionScopeSystem}
	PermissionReadLdapSyncJob = &Permission{"read_ldap_sync_job", "", "", PermissionScopeSystem}
	PermissionTestLdap = &Permission{"test_ldap", "", "", PermissionScopeSystem}
	PermissionInvalidateEmailInvite = &Permission{"invalidate_email_invite", "", "", PermissionScopeSystem}
	PermissionGetSamlMetadataFromIdp = &Permission{"get_saml_metadata_from_idp", "", "", PermissionScopeSystem}
	PermissionAddSamlPublicCert = &Permission{"add_saml_public_cert", "", "", PermissionScopeSystem}
	PermissionAddSamlPrivateCert = &Permission{"add_saml_private_cert", "", "", PermissionScopeSystem}
	PermissionAddSamlIdpCert = &Permission{"add_saml_idp_cert", "", "", PermissionScopeSystem}
	PermissionRemoveSamlPublicCert = &Permission{"remove_saml_public_cert", "", "", PermissionScopeSystem}
	PermissionRemoveSamlPrivateCert = &Permission{"remove_saml_private_cert", "", "", PermissionScopeSystem}
	PermissionRemoveSamlIdpCert = &Permission{"remove_saml_idp_cert", "", "", PermissionScopeSystem}
	PermissionGetSamlCertStatus = &Permission{"get_saml_cert_status", "", "", PermissionScopeSystem}
	PermissionAddLdapPublicCert = &Permission{"add_ldap_public_cert", "", "", PermissionScopeSystem}
	PermissionAddLdapPrivateCert = &Permission{"add_ldap_private_cert", "", "", PermissionScopeSystem}
	PermissionRemoveLdapPublicCert = &Permission{"remove_ldap_public_cert", "", "", PermissionScopeSystem}
	PermissionRemoveLdapPrivateCert = &Permission{"remove_ldap_private_cert", "", "", PermissionScopeSystem}
	PermissionGetLogs = &Permission{"get_logs", "", "", PermissionScopeSystem}
	PermissionGetAnalytics = &Permission{"get_analytics", "", "", PermissionScopeSystem}
	PermissionDownloadComplianceExportResult = &Permission{"download_compliance_export_result", "authentication.permissions.download_compliance_export_result.name", "authentication.permissions.download_compliance_export_result.description", PermissionScopeSystem}
	PermissionTestSiteUrl = &Permission{"test_site_url", "", "", PermissionScopeSystem}
	PermissionTestElasticsearch = &Permission{"test_elasticsearch", "", "", PermissionScopeSystem}
	PermissionTestS3 = &Permission{"test_s3", "", "", PermissionScopeSystem}
	PermissionReloadConfig = &Permission{"reload_config", "", "", PermissionScopeSystem}
	PermissionInvalidateCaches = &Permission{"invalidate_caches", "", "", PermissionScopeSystem}
	PermissionRecycleDatabaseConnections = &Permission{"recycle_database_connections", "", "", PermissionScopeSystem}
	PermissionPurgeElasticsearchIndexes = &Permission{"purge_elasticsearch_indexes", "", "", PermissionScopeSystem}
	PermissionTestEmail = &Permission{"test_email", "", "", PermissionScopeSystem}
	PermissionCreateElasticsearchPostIndexingJob = &Permission{"create_elasticsearch_post_indexing_job", "", "", PermissionScopeSystem}
	PermissionCreateElasticsearchPostAggregationJob = &Permission{"create_elasticsearch_post_aggregation_job", "", "", PermissionScopeSystem}
	PermissionReadElasticsearchPostIndexingJob = &Permission{"read_elasticsearch_post_indexing_job", "", "", PermissionScopeSystem}
	PermissionReadElasticsearchPostAggregationJob = &Permission{"read_elasticsearch_post_aggregation_job", "", "", PermissionScopeSystem}
	PermissionCreateUserAccessToken = &Permission{"create_user_access_token", "authentication.permissions.create_user_access_token.name", "authentication.permissions.create_user_access_token.description", PermissionScopeSystem}
	PermissionReadUserAccessToken = &Permission{"read_user_access_token", "authentication.permissions.read_user_access_token.name", "authentication.permissions.read_user_access_token.description", PermissionScopeSystem}
	PermissionRevokeUserAccessToken = &Permission{"revoke_user_access_token", "authentication.permissions.revoke_user_access_token.name", "authentication.permissions.revoke_user_access_token.description", PermissionScopeSystem}
	PermissionReadJobs = &Permission{"read_jobs", "authentication.permisssions.read_jobs.name", "authentication.permisssions.read_jobs.description", PermissionScopeSystem}
	PermissionManageJobs = &Permission{"manage_jobs", "authentication.permisssions.manage_jobs.name", "authentication.permisssions.manage_jobs.description", PermissionScopeSystem}
	PermissionViewMembers = &Permission{"view_members", "authentication.permisssions.view_members.name", "authentication.permisssions.view_members.description", PermissionScopeSystem}
	PermissionSysconsoleReadBilling = &Permission{"sysconsole_read_billing", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteBilling = &Permission{"sysconsole_write_billing", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadReportingSiteStatistics = &Permission{"sysconsole_read_reporting_site_statistics", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteReportingSiteStatistics = &Permission{"sysconsole_write_reporting_site_statistics", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadReportingServerLogs = &Permission{"sysconsole_read_reporting_server_logs", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteReportingServerLogs = &Permission{"sysconsole_write_reporting_server_logs", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadUserManagementUsers = &Permission{"sysconsole_read_user_management_users", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleWriteUserManagementUsers = &Permission{"sysconsole_write_user_management_users", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleReadUserManagementGroups = &Permission{"sysconsole_read_user_management_groups", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleWriteUserManagementGroups = &Permission{"sysconsole_write_user_management_groups", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleReadUserManagementPermissions = &Permission{"sysconsole_read_user_management_permissions", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleWriteUserManagementPermissions = &Permission{"sysconsole_write_user_management_permissions", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleReadUserManagementSystemRoles = &Permission{"sysconsole_read_user_management_system_roles", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleWriteUserManagementSystemRoles = &Permission{"sysconsole_write_user_management_system_roles", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentWebServer = &Permission{"sysconsole_read_environment_web_server", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentWebServer = &Permission{"sysconsole_write_environment_web_server", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentDatabase = &Permission{"sysconsole_read_environment_database", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentDatabase = &Permission{"sysconsole_write_environment_database", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentElasticsearch = &Permission{"sysconsole_read_environment_elasticsearch", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentElasticsearch = &Permission{"sysconsole_write_environment_elasticsearch", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentFileStorage = &Permission{"sysconsole_read_environment_file_storage", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentFileStorage = &Permission{"sysconsole_write_environment_file_storage", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentImageProxy = &Permission{"sysconsole_read_environment_image_proxy", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentImageProxy = &Permission{"sysconsole_write_environment_image_proxy", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentSmtp = &Permission{"sysconsole_read_environment_smtp", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentSmtp = &Permission{"sysconsole_write_environment_smtp", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentPushNotificationServer = &Permission{"sysconsole_read_environment_push_notification_server", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentPushNotificationServer = &Permission{"sysconsole_write_environment_push_notification_server", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentHighAvailability = &Permission{"sysconsole_read_environment_high_availability", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentHighAvailability = &Permission{"sysconsole_write_environment_high_availability", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentRateLimiting = &Permission{"sysconsole_read_environment_rate_limiting", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentRateLimiting = &Permission{"sysconsole_write_environment_rate_limiting", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentLogging = &Permission{"sysconsole_read_environment_logging", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentLogging = &Permission{"sysconsole_write_environment_logging", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentSessionLengths = &Permission{"sysconsole_read_environment_session_lengths", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentSessionLengths = &Permission{"sysconsole_write_environment_session_lengths", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentPerformanceMonitoring = &Permission{"sysconsole_read_environment_performance_monitoring", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentPerformanceMonitoring = &Permission{"sysconsole_write_environment_performance_monitoring", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadEnvironmentDeveloper = &Permission{"sysconsole_read_environment_developer", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteEnvironmentDeveloper = &Permission{"sysconsole_write_environment_developer", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSiteCustomization = &Permission{"sysconsole_read_site_customization", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSiteCustomization = &Permission{"sysconsole_write_site_customization", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSiteLocalization = &Permission{"sysconsole_read_site_localization", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSiteLocalization = &Permission{"sysconsole_write_site_localization", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSiteNotifications = &Permission{"sysconsole_read_site_notifications", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSiteNotifications = &Permission{"sysconsole_write_site_notifications", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSiteAnnouncementBanner = &Permission{"sysconsole_read_site_announcement_banner", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSiteAnnouncementBanner = &Permission{"sysconsole_write_site_announcement_banner", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSitePosts = &Permission{"sysconsole_read_site_posts", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSitePosts = &Permission{"sysconsole_write_site_posts", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSiteFileSharingAndDownloads = &Permission{"sysconsole_read_site_file_sharing_and_downloads", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSiteFileSharingAndDownloads = &Permission{"sysconsole_write_site_file_sharing_and_downloads", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSitePublicLinks = &Permission{"sysconsole_read_site_public_links", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSitePublicLinks = &Permission{"sysconsole_write_site_public_links", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadSiteNotices = &Permission{"sysconsole_read_site_notices", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteSiteNotices = &Permission{"sysconsole_write_site_notices", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationSignup = &Permission{"sysconsole_read_authentication_signup", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationSignup = &Permission{"sysconsole_write_authentication_signup", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationEmail = &Permission{"sysconsole_read_authentication_email", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationEmail = &Permission{"sysconsole_write_authentication_email", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationPassword = &Permission{"sysconsole_read_authentication_password", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationPassword = &Permission{"sysconsole_write_authentication_password", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationMfa = &Permission{"sysconsole_read_authentication_mfa", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationMfa = &Permission{"sysconsole_write_authentication_mfa", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationLdap = &Permission{"sysconsole_read_authentication_ldap", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationLdap = &Permission{"sysconsole_write_authentication_ldap", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationSaml = &Permission{"sysconsole_read_authentication_saml", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationSaml = &Permission{"sysconsole_write_authentication_saml", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationOpenid = &Permission{"sysconsole_read_authentication_openid", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationOpenid = &Permission{"sysconsole_write_authentication_openid", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadAuthenticationGuestAccess = &Permission{"sysconsole_read_authentication_guest_access", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteAuthenticationGuestAccess = &Permission{"sysconsole_write_authentication_guest_access", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadPlugins = &Permission{"sysconsole_read_plugins", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleWritePlugins = &Permission{"sysconsole_write_plugins", "authentication.permissions.use_group_mentions.name", "authentication.permissions.use_group_mentions.description", PermissionScopeSystem}
	PermissionSysconsoleReadIntegrationsIntegrationManagement = &Permission{"sysconsole_read_integrations_integration_management", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteIntegrationsIntegrationManagement = &Permission{"sysconsole_write_integrations_integration_management", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadIntegrationsBotAccounts = &Permission{"sysconsole_read_integrations_bot_accounts", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteIntegrationsBotAccounts = &Permission{"sysconsole_write_integrations_bot_accounts", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadIntegrationsGif = &Permission{"sysconsole_read_integrations_gif", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteIntegrationsGif = &Permission{"sysconsole_write_integrations_gif", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadIntegrationsCors = &Permission{"sysconsole_read_integrations_cors", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteIntegrationsCors = &Permission{"sysconsole_write_integrations_cors", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadComplianceDataRetentionPolicy = &Permission{"sysconsole_read_compliance_data_retention_policy", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteComplianceDataRetentionPolicy = &Permission{"sysconsole_write_compliance_data_retention_policy", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadComplianceComplianceExport = &Permission{"sysconsole_read_compliance_compliance_export", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteComplianceComplianceExport = &Permission{"sysconsole_write_compliance_compliance_export", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadComplianceComplianceMonitoring = &Permission{"sysconsole_read_compliance_compliance_monitoring", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteComplianceComplianceMonitoring = &Permission{"sysconsole_write_compliance_compliance_monitoring", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadComplianceCustomTermsOfService = &Permission{"sysconsole_read_compliance_custom_terms_of_service", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteComplianceCustomTermsOfService = &Permission{"sysconsole_write_compliance_custom_terms_of_service", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadExperimentalFeatures = &Permission{"sysconsole_read_experimental_features", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteExperimentalFeatures = &Permission{"sysconsole_write_experimental_features", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadExperimentalFeatureFlags = &Permission{"sysconsole_read_experimental_feature_flags", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteExperimentalFeatureFlags = &Permission{"sysconsole_write_experimental_feature_flags", "", "", PermissionScopeSystem}
	PermissionSysconsoleReadExperimentalBleve = &Permission{"sysconsole_read_experimental_bleve", "", "", PermissionScopeSystem}
	PermissionSysconsoleWriteExperimentalBleve = &Permission{"sysconsole_write_experimental_bleve", "", "", PermissionScopeSystem}

	PermissionCreateShop = &Permission{"create_shop", "", "", PermissionScopeSystem}
	PermissionCreateCategoryTranslation = &Permission{"create_categorytranslation", "", "", PermissionScopeSystem}
	PermissionReadCategoryTranslation = &Permission{"read_categorytranslation", "", "", PermissionScopeSystem}
	PermissionUpdateCategoryTranslation = &Permission{"update_categorytranslation", "", "", PermissionScopeSystem}
	PermissionDeleteCategoryTranslation = &Permission{"delete_categorytranslation", "", "", PermissionScopeSystem}

	PermissionCreateOpenExchangeRate = &Permission{"create_openexchangerate", "", "", PermissionScopeSystem}
	PermissionReadOpenExchangeRate = &Permission{"read_openexchangerate", "", "", PermissionScopeSystem}
	PermissionUpdateOpenExchangeRate = &Permission{"update_openexchangerate", "", "", PermissionScopeSystem}
	PermissionDeleteOpenExchangeRate = &Permission{"delete_openexchangerate", "", "", PermissionScopeSystem}

	PermissionCreateToken = &Permission{"create_token", "", "", PermissionScopeSystem}
	PermissionReadToken = &Permission{"read_token", "", "", PermissionScopeSystem}
	// PermissionUpdateToken = &Permission{"update_token", "", "", PermissionScopeSystem}
	PermissionDeleteToken = &Permission{"delete_token", "", "", PermissionScopeSystem}
	PermissionCreateTermsOfService = &Permission{"create_termsofservice", "", "", PermissionScopeSystem}
	PermissionReadTermsOfService = &Permission{"read_termsofservice", "", "", PermissionScopeSystem}
	PermissionUpdateTermsOfService = &Permission{"update_termsofservice", "", "", PermissionScopeSystem}
	PermissionDeleteTermsOfService = &Permission{"delete_termsofservice", "", "", PermissionScopeSystem}
	PermissionCreateChannel = &Permission{"create_channel", "", "", PermissionScopeSystem}
	PermissionReadChannel = &Permission{"read_channel", "", "", PermissionScopeSystem}
	PermissionUpdateChannel = &Permission{"update_channel", "", "", PermissionScopeSystem}
	PermissionDeleteChannel = &Permission{"delete_channel", "", "", PermissionScopeSystem}
	PermissionCreateRole = &Permission{"create_role", "", "", PermissionScopeSystem}
	PermissionReadRole = &Permission{"read_role", "", "", PermissionScopeSystem}
	PermissionUpdateRole = &Permission{"update_role", "", "", PermissionScopeSystem}
	PermissionDeleteRole = &Permission{"delete_role", "", "", PermissionScopeSystem}

	PermissionCreateCompliance = &Permission{"create_compliance", "", "", PermissionScopeSystem}
	PermissionReadCompliance = &Permission{"read_compliance", "", "", PermissionScopeSystem}
	PermissionUpdateCompliance = &Permission{"update_compliance", "", "", PermissionScopeSystem}
	PermissionDeleteCompliance = &Permission{"delete_compliance", "", "", PermissionScopeSystem}
	PermissionCreateClusterDiscovery = &Permission{"create_clusterdiscovery", "", "", PermissionScopeSystem}
	PermissionReadClusterDiscovery = &Permission{"read_clusterdiscovery", "", "", PermissionScopeSystem}
	PermissionUpdateClusterDiscovery = &Permission{"update_clusterdiscovery", "", "", PermissionScopeSystem}
	PermissionDeleteClusterDiscovery = &Permission{"delete_clusterdiscovery", "", "", PermissionScopeSystem}
	PermissionCreateAudit = &Permission{"create_audit", "", "", PermissionScopeSystem}
	PermissionReadAudit = &Permission{"read_audit", "", "", PermissionScopeSystem}
	PermissionUpdateAudit = &Permission{"update_audit", "", "", PermissionScopeSystem}
	PermissionDeleteAudit = &Permission{"delete_audit", "", "", PermissionScopeSystem}
	PermissionCreateTransaction = &Permission{"create_transaction", "", "", PermissionScopeSystem}
	PermissionReadTransaction = &Permission{"read_transaction", "", "", PermissionScopeSystem}
	PermissionUpdateTransaction = &Permission{"update_transaction", "", "", PermissionScopeSystem}
	PermissionDeleteTransaction = &Permission{"delete_transaction", "", "", PermissionScopeSystem}
	PermissionUpdateAttributeValue = &Permission{"update_attributevalue", "", "", PermissionScopeSystem}
	PermissionDeleteAttributeValue = &Permission{"delete_attributevalue", "", "", PermissionScopeSystem}

	PermissionCreatePreference = &Permission{"create_preference", "", "", PermissionScopeSystem}
	PermissionReadPreference = &Permission{"read_preference", "", "", PermissionScopeSystem}
	PermissionUpdatePreference = &Permission{"update_preference", "", "", PermissionScopeSystem}
	PermissionDeletePreference = &Permission{"delete_preference", "", "", PermissionScopeSystem}
	PermissionCreateCategory = &Permission{"create_category", "", "", PermissionScopeSystem}
	PermissionReadCategory = &Permission{"read_category", "", "", PermissionScopeSystem}
	PermissionUpdateCategory = &Permission{"update_category", "", "", PermissionScopeSystem}
	PermissionDeleteCategory = &Permission{"delete_category", "", "", PermissionScopeSystem}
	PermissionCreateAddress = &Permission{"create_address", "", "", PermissionScopeSystem}
	PermissionReadAddress = &Permission{"read_address", "", "", PermissionScopeSystem}
	PermissionUpdateAddress = &Permission{"update_address", "", "", PermissionScopeSystem}
	PermissionDeleteAddress = &Permission{"delete_address", "", "", PermissionScopeSystem}
	PermissionCreateUser = &Permission{"create_user", "", "", PermissionScopeSystem}
	PermissionReadUser = &Permission{"read_user", "", "", PermissionScopeSystem}
	PermissionDeleteUser = &Permission{"delete_user", "", "", PermissionScopeSystem}
}

func initSystemPermissionGroups() {
	SysconsoleReadPermissions = Permissions{
		PermissionSysconsoleReadBilling,
		PermissionSysconsoleReadReportingSiteStatistics,
		PermissionSysconsoleReadReportingServerLogs,
		PermissionSysconsoleReadUserManagementUsers,
		PermissionSysconsoleReadUserManagementGroups,
		PermissionSysconsoleReadUserManagementPermissions,
		PermissionSysconsoleReadUserManagementSystemRoles,
		PermissionSysconsoleReadEnvironmentWebServer,
		PermissionSysconsoleReadEnvironmentDatabase,
		PermissionSysconsoleReadEnvironmentElasticsearch,
		PermissionSysconsoleReadEnvironmentFileStorage,
		PermissionSysconsoleReadEnvironmentImageProxy,
		PermissionSysconsoleReadEnvironmentSmtp,
		PermissionSysconsoleReadEnvironmentPushNotificationServer,
		PermissionSysconsoleReadEnvironmentHighAvailability,
		PermissionSysconsoleReadEnvironmentRateLimiting,
		PermissionSysconsoleReadEnvironmentLogging,
		PermissionSysconsoleReadEnvironmentSessionLengths,
		PermissionSysconsoleReadEnvironmentPerformanceMonitoring,
		PermissionSysconsoleReadEnvironmentDeveloper,
		PermissionSysconsoleReadSiteCustomization,
		PermissionSysconsoleReadSiteLocalization,
		PermissionSysconsoleReadSiteNotifications,
		PermissionSysconsoleReadSiteAnnouncementBanner,
		PermissionSysconsoleReadSitePosts,
		PermissionSysconsoleReadSiteFileSharingAndDownloads,
		PermissionSysconsoleReadSitePublicLinks,
		PermissionSysconsoleReadSiteNotices,
		PermissionSysconsoleReadAuthenticationSignup,
		PermissionSysconsoleReadAuthenticationEmail,
		PermissionSysconsoleReadAuthenticationPassword,
		PermissionSysconsoleReadAuthenticationMfa,
		PermissionSysconsoleReadAuthenticationLdap,
		PermissionSysconsoleReadAuthenticationSaml,
		PermissionSysconsoleReadAuthenticationOpenid,
		PermissionSysconsoleReadAuthenticationGuestAccess,
		PermissionSysconsoleReadPlugins,
		PermissionSysconsoleReadIntegrationsIntegrationManagement,
		PermissionSysconsoleReadIntegrationsBotAccounts,
		PermissionSysconsoleReadIntegrationsGif,
		PermissionSysconsoleReadIntegrationsCors,
		PermissionSysconsoleReadComplianceDataRetentionPolicy,
		PermissionSysconsoleReadComplianceComplianceExport,
		PermissionSysconsoleReadComplianceComplianceMonitoring,
		PermissionSysconsoleReadComplianceCustomTermsOfService,
		PermissionSysconsoleReadExperimentalFeatures,
		PermissionSysconsoleReadExperimentalFeatureFlags,
		PermissionSysconsoleReadExperimentalBleve,
	}

	SysconsoleWritePermissions = Permissions{
		PermissionSysconsoleWriteBilling,
		PermissionSysconsoleWriteReportingSiteStatistics,
		PermissionSysconsoleWriteReportingServerLogs,
		PermissionSysconsoleWriteUserManagementUsers,
		PermissionSysconsoleWriteUserManagementGroups,
		PermissionSysconsoleWriteUserManagementPermissions,
		PermissionSysconsoleWriteUserManagementSystemRoles,
		PermissionSysconsoleWriteEnvironmentWebServer,
		PermissionSysconsoleWriteEnvironmentDatabase,
		PermissionSysconsoleWriteEnvironmentElasticsearch,
		PermissionSysconsoleWriteEnvironmentFileStorage,
		PermissionSysconsoleWriteEnvironmentImageProxy,
		PermissionSysconsoleWriteEnvironmentSmtp,
		PermissionSysconsoleWriteEnvironmentPushNotificationServer,
		PermissionSysconsoleWriteEnvironmentHighAvailability,
		PermissionSysconsoleWriteEnvironmentRateLimiting,
		PermissionSysconsoleWriteEnvironmentLogging,
		PermissionSysconsoleWriteEnvironmentSessionLengths,
		PermissionSysconsoleWriteEnvironmentPerformanceMonitoring,
		PermissionSysconsoleWriteEnvironmentDeveloper,
		PermissionSysconsoleWriteSiteCustomization,
		PermissionSysconsoleWriteSiteLocalization,
		PermissionSysconsoleWriteSiteNotifications,
		PermissionSysconsoleWriteSiteAnnouncementBanner,
		PermissionSysconsoleWriteSitePosts,
		PermissionSysconsoleWriteSiteFileSharingAndDownloads,
		PermissionSysconsoleWriteSitePublicLinks,
		PermissionSysconsoleWriteSiteNotices,
		PermissionSysconsoleWriteAuthenticationSignup,
		PermissionSysconsoleWriteAuthenticationEmail,
		PermissionSysconsoleWriteAuthenticationPassword,
		PermissionSysconsoleWriteAuthenticationMfa,
		PermissionSysconsoleWriteAuthenticationLdap,
		PermissionSysconsoleWriteAuthenticationSaml,
		PermissionSysconsoleWriteAuthenticationOpenid,
		PermissionSysconsoleWriteAuthenticationGuestAccess,
		PermissionSysconsoleWritePlugins,
		PermissionSysconsoleWriteIntegrationsIntegrationManagement,
		PermissionSysconsoleWriteIntegrationsBotAccounts,
		PermissionSysconsoleWriteIntegrationsGif,
		PermissionSysconsoleWriteIntegrationsCors,
		PermissionSysconsoleWriteComplianceDataRetentionPolicy,
		PermissionSysconsoleWriteComplianceComplianceExport,
		PermissionSysconsoleWriteComplianceComplianceMonitoring,
		PermissionSysconsoleWriteComplianceCustomTermsOfService,
		PermissionSysconsoleWriteExperimentalFeatures,
		PermissionSysconsoleWriteExperimentalFeatureFlags,
		PermissionSysconsoleWriteExperimentalBleve,
	}

	SystemScopedPermissionsMinusSysconsole = Permissions{
		PermissionAddReaction,
		PermissionRemoveReaction,
		PermissionAssignSystemAdminRole,
		PermissionManageRoles,
		PermissionManageSystem,
		PermissionEditOtherUsers,
		PermissionGetPublicLink,
		PermissionManageOAuth,
		PermissionManageSystemWideOAuth,
		PermissionCreateUserAccessToken,
		PermissionReadUserAccessToken,
		PermissionRevokeUserAccessToken,
		PermissionReadJobs,
		PermissionManageJobs,
		PermissionManageRemoteClusters,
		PermissionDownloadComplianceExportResult,
		PermissionCreateDataRetentionJob,
		PermissionReadDataRetentionJob,
		PermissionCreateComplianceExportJob,
		PermissionReadComplianceExportJob,
		PermissionReadAudits,
		PermissionTestSiteUrl,
		PermissionTestElasticsearch,
		PermissionTestS3,
		PermissionReloadConfig,
		PermissionInvalidateCaches,
		PermissionRecycleDatabaseConnections,
		PermissionPurgeElasticsearchIndexes,
		PermissionTestEmail,
		PermissionCreateElasticsearchPostIndexingJob,
		PermissionCreateElasticsearchPostAggregationJob,
		PermissionReadElasticsearchPostIndexingJob,
		PermissionReadElasticsearchPostAggregationJob,
		PermissionPurgeBleveIndexes,
		PermissionCreatePostBleveIndexesJob,
		PermissionCreateLdapSyncJob,
		PermissionReadLdapSyncJob,
		PermissionTestLdap,
		PermissionInvalidateEmailInvite,
		PermissionGetSamlMetadataFromIdp,
		PermissionAddSamlPublicCert,
		PermissionAddSamlPrivateCert,
		PermissionAddSamlIdpCert,
		PermissionRemoveSamlPublicCert,
		PermissionRemoveSamlPrivateCert,
		PermissionRemoveSamlIdpCert,
		PermissionGetSamlCertStatus,
		PermissionAddLdapPublicCert,
		PermissionAddLdapPrivateCert,
		PermissionRemoveLdapPublicCert,
		PermissionRemoveLdapPrivateCert,
		PermissionGetAnalytics,
		PermissionGetLogs,
		PermissionManageIncomingWebhooks,
		PermissionManageOutgoingWebhooks,
		PermissionManageOthersIncomingWebhooks,
		PermissionManageOthersOutgoingWebhooks,
		PermissionCreateCategoryTranslation,
		PermissionReadCategoryTranslation,
		PermissionUpdateCategoryTranslation,
		PermissionDeleteCategoryTranslation,
		PermissionCreateOpenExchangeRate,
		PermissionReadOpenExchangeRate,
		PermissionUpdateOpenExchangeRate,
		PermissionDeleteOpenExchangeRate,
		PermissionCreateAudit,
		PermissionReadAudit,
		PermissionUpdateAudit,
		PermissionDeleteAudit,
		PermissionCreateClusterDiscovery,
		PermissionReadClusterDiscovery,
		PermissionUpdateClusterDiscovery,
		PermissionDeleteClusterDiscovery,
		PermissionCreateRole, // important permissions
		PermissionReadRole,
		PermissionUpdateRole,
		PermissionDeleteRole,
		PermissionCreateCompliance,
		PermissionReadCompliance,
		PermissionUpdateCompliance,
		PermissionDeleteCompliance,
		PermissionCreateChannel,
		PermissionReadChannel,
		PermissionUpdateChannel,
		PermissionDeleteChannel,
		PermissionCreateTermsOfService,
		PermissionReadTermsOfService,
		PermissionUpdateTermsOfService,
		PermissionDeleteTermsOfService,
		PermissionCreateToken,
		PermissionReadToken,
		PermissionDeleteToken,
		PermissionCreateUser,
		PermissionReadUser,
		PermissionDeleteUser,
		PermissionCreateCategory,
		PermissionReadCategory,
		PermissionUpdateCategory,
		PermissionDeleteCategory,

		PermissionCreateAttribute,
		PermissionUpdateAttribute,
		PermissionDeleteAttribute,
		PermissionUpdateAttributeValue,
		PermissionDeleteAttributeValue,
		PermissionCreateAttributeValue,
	}

	AllSystemScopedPermissions = append(SystemScopedPermissionsMinusSysconsole, SysconsoleReadPermissions...)
	AllSystemScopedPermissions = append(AllSystemScopedPermissions, SysconsoleWritePermissions...)
}

// initPermissions is called inside init() function located in /model/init.go
func initPermissions() {
	initializeShopScopedPermissions()
	initializeSystemScopedPermissions()

	initSystemPermissionGroups()
	initShopPermissionGroups()
}
