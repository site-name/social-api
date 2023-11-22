package model_helper

import (
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

type RolePatch struct {
	Permissions *[]string `json:"permissions"`
}

type RolePermissions struct {
	RoleID      string
	Permissions []string
}

func Rolepatching(role *model.Role, patch *RolePatch) {
	if patch.Permissions != nil {
		role.Permissions = strings.Join(*patch.Permissions, " ")
	}
}

// Returns an array of permissions that are in either role.Permissions
// or patch.Permissions, but not both.
func PermissionsChangedByPatch(role *model.Role, patch *RolePatch) []string {
	var result []string

	if patch.Permissions == nil {
		return result
	}

	roleMap := make(map[string]bool)
	patchMap := make(map[string]bool)

	rolePermissions := strings.Fields(role.Permissions)

	for _, permission := range rolePermissions {
		roleMap[permission] = true
	}
	for _, permission := range *patch.Permissions {
		patchMap[permission] = true
	}
	for _, permission := range rolePermissions {
		if !patchMap[permission] {
			result = append(result, permission)
		}
	}
	for _, permission := range *patch.Permissions {
		if !roleMap[permission] {
			result = append(result, permission)
		}
	}

	return result
}

func RoleIsValid(r model.Role) bool {
	return RoleIsValidWithoutId(r)
}

// IsValidWithoutId check if current role is valid without checking its Id
func RoleIsValidWithoutId(r model.Role) bool {
	if !IsValidRoleName(r.Name) {
		return false
	}

	// check checks if permissionId is included in perms
	var check = func(perms []*Permission, permissionId string) bool {
		for _, p := range perms {
			if permissionId == p.Id {
				return true
			}
		}
		return false
	}

	rolePermissions := strings.Fields(r.Permissions)

	for _, permissionId := range rolePermissions {
		if !check(AllSystemScopedPermissions, permissionId) {
			return false
		}
	}

	return true
}

// IsValidRoleName checks if roleName's length is valid ( > 0 && <= 64) AND
//
// contains no character other than English's 26 letters, 10 digits and underscore
func IsValidRoleName(roleName string) bool {
	if roleName == "" {
		return false
	}

	if strings.TrimLeft(roleName, "abcdefghijklmnopqrstuvwxyz0123456789_") != "" {
		return false
	}

	return true
}

// CleanRoleNames iterates through given roleNames.
// checks if each name is valid.
func CleanRoleNames(roleNames []string) ([]string, bool) {
	var cleanedRoleNames []string
	for _, roleName := range roleNames {
		if strings.TrimSpace(roleName) == "" {
			continue
		}

		if !IsValidRoleName(roleName) {
			return roleNames, false
		}

		cleanedRoleNames = append(cleanedRoleNames, roleName)
	}

	return cleanedRoleNames, true
}

// MakeDefaultRoles make an map with values are default roles
func MakeDefaultRoles() map[string]*model.Role {
	roles := make(map[string]*model.Role)

	roles[SystemUserRoleId] = &model.Role{
		Name:          SystemUserRoleId,
		DisplayName:   "authentication.roles.global_user.name",
		Description:   "authentication.roles.global_user.description",
		Permissions:   SystemUserPermissions.IDs().Join(" "),
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[SystemUserAccessTokenRoleId] = &model.Role{
		Name:        SystemUserAccessTokenRoleId,
		DisplayName: "authentication.roles.system_user_access_token.name",
		Description: "authentication.roles.system_user_access_token.description",
		Permissions: util.AnyArray[string]{
			PermissionCreateUserAccessToken.Id,
			PermissionReadUserAccessToken.Id,
			PermissionRevokeUserAccessToken.Id,
		}.Join(" "),
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemUserManagerRoleId] = &model.Role{
		Name:          SystemUserManagerRoleId,
		DisplayName:   "authentication.roles.system_user_manager.name",
		Description:   "authentication.roles.system_user_manager.description",
		Permissions:   SystemUserManagerDefaultPermissions.Join(" "),
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemReadOnlyAdminRoleId] = &model.Role{
		Name:          SystemReadOnlyAdminRoleId,
		DisplayName:   "authentication.roles.system_read_only_admin.name",
		Description:   "authentication.roles.system_read_only_admin.description",
		Permissions:   SystemReadOnlyAdminDefaultPermissions.Join(" "),
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemManagerRoleId] = &model.Role{
		Name:          SystemManagerRoleId,
		DisplayName:   "authentication.roles.system_manager.name",
		Description:   "authentication.roles.system_manager.description",
		Permissions:   SystemManagerDefaultPermissions.Join(" "),
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemAdminRoleId] = &model.Role{
		Name:          SystemAdminRoleId,
		DisplayName:   "authentication.roles.global_admin.name",
		Description:   "authentication.roles.global_admin.description",
		Permissions:   AllSystemScopedPermissions.IDs().Join(" "),
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[ShopAdminRoleId] = &model.Role{
		Name:          ShopAdminRoleId,
		DisplayName:   "authentication.roles.shop_admin.name",
		Description:   "authentication.roles.shop_admin.description",
		Permissions:   ShopScopedAllPermissions.IDs().Join(" "),
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[ShopStaffRoleId] = &model.Role{
		Name:          ShopStaffRoleId,
		DisplayName:   "authentication.roles.shop_staff.name",
		Description:   "authentication.roles.shop_staff.description",
		Permissions:   ShopStaffPermissions.IDs().Join(" "),
		SchemeManaged: true,
		BuiltIn:       true,
	}

	return roles
}

func AddAncillaryPermissions(permissions []string) []string {
	for _, permission := range permissions {
		if ancillaryPermissions, ok := SysconsoleAncillaryPermissions[permission]; ok {
			permissions = append(permissions, ancillaryPermissions...)
		}
	}
	return permissions
}

// initRoles is called be the init() function located in /model.init.go file
func initRoles() {
	NewSystemRoleIDs = []string{
		SystemUserManagerRoleId,
		SystemReadOnlyAdminRoleId,
		SystemManagerRoleId,
	}

	BuiltInSchemeManagedRoleIDs = append([]string{
		SystemUserRoleId,
		SystemAdminRoleId,
		SystemUserAccessTokenRoleId,
	}, NewSystemRoleIDs...)

	// When updating the values here, the values in mattermost-redux must also be updated.
	SysconsoleAncillaryPermissions = map[string][]string{
		PermissionSysconsoleReadUserManagementUsers.Id: {
			PermissionGetAnalytics.Id,
		},
		PermissionSysconsoleReadEnvironmentElasticsearch.Id: {
			PermissionReadElasticsearchPostIndexingJob.Id,
			PermissionReadElasticsearchPostAggregationJob.Id,
		},
		PermissionSysconsoleWriteEnvironmentElasticsearch.Id: {
			PermissionTestElasticsearch.Id,
			PermissionCreateElasticsearchPostIndexingJob.Id,
			PermissionCreateElasticsearchPostAggregationJob.Id,
			PermissionPurgeElasticsearchIndexes.Id,
		},
		PermissionSysconsoleWriteEnvironmentWebServer.Id: {
			PermissionTestSiteUrl.Id,
			PermissionReloadConfig.Id,
			PermissionInvalidateCaches.Id,
		},
		PermissionSysconsoleWriteEnvironmentDatabase.Id: {
			PermissionRecycleDatabaseConnections.Id,
		},
		PermissionSysconsoleWriteEnvironmentFileStorage.Id: {
			PermissionTestS3.Id,
		},
		PermissionSysconsoleWriteEnvironmentSmtp.Id: {
			PermissionTestEmail.Id,
		},
		PermissionSysconsoleReadReportingServerLogs.Id: {
			PermissionGetLogs.Id,
		},
		PermissionSysconsoleReadReportingSiteStatistics.Id: {
			PermissionGetAnalytics.Id,
		},
		PermissionSysconsoleWriteUserManagementUsers.Id: {
			PermissionEditOtherUsers.Id,
		},
		PermissionSysconsoleWriteSiteCustomization.Id: {},
		PermissionSysconsoleWriteComplianceDataRetentionPolicy.Id: {
			PermissionCreateDataRetentionJob.Id,
		},
		PermissionSysconsoleReadComplianceDataRetentionPolicy.Id: {
			PermissionReadDataRetentionJob.Id,
		},
		PermissionSysconsoleWriteComplianceComplianceExport.Id: {
			PermissionCreateComplianceExportJob.Id,
			PermissionDownloadComplianceExportResult.Id,
		},
		PermissionSysconsoleReadComplianceComplianceExport.Id: {
			PermissionReadComplianceExportJob.Id,
			PermissionDownloadComplianceExportResult.Id,
		},
		PermissionSysconsoleReadComplianceCustomTermsOfService.Id: {
			PermissionReadAudits.Id,
		},
		PermissionSysconsoleWriteExperimentalBleve.Id: {
			PermissionCreatePostBleveIndexesJob.Id,
			PermissionPurgeBleveIndexes.Id,
		},
		PermissionSysconsoleWriteAuthenticationLdap.Id: {
			PermissionCreateLdapSyncJob.Id,
			PermissionAddLdapPublicCert.Id,
			PermissionRemoveLdapPublicCert.Id,
			PermissionAddLdapPrivateCert.Id,
			PermissionRemoveLdapPrivateCert.Id,
		},
		PermissionSysconsoleReadAuthenticationLdap.Id: {
			PermissionTestLdap.Id,
			PermissionReadLdapSyncJob.Id,
		},
		PermissionSysconsoleWriteAuthenticationEmail.Id: {
			PermissionInvalidateEmailInvite.Id,
		},
		PermissionSysconsoleWriteAuthenticationSaml.Id: {
			PermissionGetSamlMetadataFromIdp.Id,
			PermissionAddSamlPublicCert.Id,
			PermissionAddSamlPrivateCert.Id,
			PermissionAddSamlIdpCert.Id,
			PermissionRemoveSamlPublicCert.Id,
			PermissionRemoveSamlPrivateCert.Id,
			PermissionRemoveSamlIdpCert.Id,
			PermissionGetSamlCertStatus.Id,
		},
	}

	SystemUserManagerDefaultPermissions = []string{
		PermissionSysconsoleReadUserManagementGroups.Id,
		PermissionSysconsoleReadUserManagementPermissions.Id,
		PermissionSysconsoleWriteUserManagementGroups.Id,
		PermissionSysconsoleReadAuthenticationSignup.Id,
		PermissionSysconsoleReadAuthenticationEmail.Id,
		PermissionSysconsoleReadAuthenticationPassword.Id,
		PermissionSysconsoleReadAuthenticationMfa.Id,
		PermissionSysconsoleReadAuthenticationLdap.Id,
		PermissionSysconsoleReadAuthenticationSaml.Id,
		PermissionSysconsoleReadAuthenticationOpenid.Id,
		PermissionSysconsoleReadAuthenticationGuestAccess.Id,
		PermissionDeleteUser.Id,
		PermissionReadUser.Id,
	}

	SystemReadOnlyAdminDefaultPermissions = []string{
		PermissionSysconsoleReadReportingSiteStatistics.Id,
		PermissionSysconsoleReadReportingServerLogs.Id,
		PermissionSysconsoleReadUserManagementUsers.Id,
		PermissionSysconsoleReadUserManagementGroups.Id,
		PermissionSysconsoleReadUserManagementPermissions.Id,
		PermissionSysconsoleReadEnvironmentWebServer.Id,
		PermissionSysconsoleReadEnvironmentDatabase.Id,
		PermissionSysconsoleReadEnvironmentElasticsearch.Id,
		PermissionSysconsoleReadEnvironmentFileStorage.Id,
		PermissionSysconsoleReadEnvironmentImageProxy.Id,
		PermissionSysconsoleReadEnvironmentSmtp.Id,
		PermissionSysconsoleReadEnvironmentPushNotificationServer.Id,
		PermissionSysconsoleReadEnvironmentHighAvailability.Id,
		PermissionSysconsoleReadEnvironmentRateLimiting.Id,
		PermissionSysconsoleReadEnvironmentLogging.Id,
		PermissionSysconsoleReadEnvironmentSessionLengths.Id,
		PermissionSysconsoleReadEnvironmentPerformanceMonitoring.Id,
		PermissionSysconsoleReadEnvironmentDeveloper.Id,
		PermissionSysconsoleReadSiteCustomization.Id,
		PermissionSysconsoleReadSiteLocalization.Id,
		PermissionSysconsoleReadSiteNotifications.Id,
		PermissionSysconsoleReadSiteAnnouncementBanner.Id,
		PermissionSysconsoleReadSitePosts.Id,
		PermissionSysconsoleReadSiteFileSharingAndDownloads.Id,
		PermissionSysconsoleReadSitePublicLinks.Id,
		PermissionSysconsoleReadSiteNotices.Id,
		PermissionSysconsoleReadAuthenticationSignup.Id,
		PermissionSysconsoleReadAuthenticationEmail.Id,
		PermissionSysconsoleReadAuthenticationPassword.Id,
		PermissionSysconsoleReadAuthenticationMfa.Id,
		PermissionSysconsoleReadAuthenticationLdap.Id,
		PermissionSysconsoleReadAuthenticationSaml.Id,
		PermissionSysconsoleReadAuthenticationOpenid.Id,
		PermissionSysconsoleReadAuthenticationGuestAccess.Id,
		PermissionSysconsoleReadPlugins.Id,
		PermissionSysconsoleReadIntegrationsIntegrationManagement.Id,
		PermissionSysconsoleReadIntegrationsBotAccounts.Id,
		PermissionSysconsoleReadIntegrationsGif.Id,
		PermissionSysconsoleReadIntegrationsCors.Id,
		PermissionSysconsoleReadComplianceDataRetentionPolicy.Id,
		PermissionSysconsoleReadComplianceComplianceExport.Id,
		PermissionSysconsoleReadComplianceComplianceMonitoring.Id,
		PermissionSysconsoleReadComplianceCustomTermsOfService.Id,
		PermissionSysconsoleReadExperimentalFeatures.Id,
		PermissionSysconsoleReadExperimentalFeatureFlags.Id,
		PermissionSysconsoleReadExperimentalBleve.Id,
	}

	SystemManagerDefaultPermissions = []string{
		PermissionSysconsoleReadReportingSiteStatistics.Id,
		PermissionSysconsoleReadReportingServerLogs.Id,
		PermissionSysconsoleReadUserManagementGroups.Id,
		PermissionSysconsoleReadUserManagementPermissions.Id,
		PermissionSysconsoleWriteUserManagementGroups.Id,
		PermissionSysconsoleWriteUserManagementPermissions.Id,
		PermissionSysconsoleReadEnvironmentWebServer.Id,
		PermissionSysconsoleReadEnvironmentDatabase.Id,
		PermissionSysconsoleReadEnvironmentElasticsearch.Id,
		PermissionSysconsoleReadEnvironmentFileStorage.Id,
		PermissionSysconsoleReadEnvironmentImageProxy.Id,
		PermissionSysconsoleReadEnvironmentSmtp.Id,
		PermissionSysconsoleReadEnvironmentPushNotificationServer.Id,
		PermissionSysconsoleReadEnvironmentHighAvailability.Id,
		PermissionSysconsoleReadEnvironmentRateLimiting.Id,
		PermissionSysconsoleReadEnvironmentLogging.Id,
		PermissionSysconsoleReadEnvironmentSessionLengths.Id,
		PermissionSysconsoleReadEnvironmentPerformanceMonitoring.Id,
		PermissionSysconsoleReadEnvironmentDeveloper.Id,
		PermissionSysconsoleWriteEnvironmentWebServer.Id,
		PermissionSysconsoleWriteEnvironmentDatabase.Id,
		PermissionSysconsoleWriteEnvironmentElasticsearch.Id,
		PermissionSysconsoleWriteEnvironmentFileStorage.Id,
		PermissionSysconsoleWriteEnvironmentImageProxy.Id,
		PermissionSysconsoleWriteEnvironmentSmtp.Id,
		PermissionSysconsoleWriteEnvironmentPushNotificationServer.Id,
		PermissionSysconsoleWriteEnvironmentHighAvailability.Id,
		PermissionSysconsoleWriteEnvironmentRateLimiting.Id,
		PermissionSysconsoleWriteEnvironmentLogging.Id,
		PermissionSysconsoleWriteEnvironmentSessionLengths.Id,
		PermissionSysconsoleWriteEnvironmentPerformanceMonitoring.Id,
		PermissionSysconsoleWriteEnvironmentDeveloper.Id,
		PermissionSysconsoleReadSiteCustomization.Id,
		PermissionSysconsoleWriteSiteCustomization.Id,
		PermissionSysconsoleReadSiteLocalization.Id,
		PermissionSysconsoleWriteSiteLocalization.Id,
		PermissionSysconsoleReadSiteNotifications.Id,
		PermissionSysconsoleWriteSiteNotifications.Id,
		PermissionSysconsoleReadSiteAnnouncementBanner.Id,
		PermissionSysconsoleWriteSiteAnnouncementBanner.Id,
		PermissionSysconsoleReadSitePosts.Id,
		PermissionSysconsoleWriteSitePosts.Id,
		PermissionSysconsoleReadSiteFileSharingAndDownloads.Id,
		PermissionSysconsoleWriteSiteFileSharingAndDownloads.Id,
		PermissionSysconsoleReadSitePublicLinks.Id,
		PermissionSysconsoleWriteSitePublicLinks.Id,
		PermissionSysconsoleReadSiteNotices.Id,
		PermissionSysconsoleWriteSiteNotices.Id,
		PermissionSysconsoleReadAuthenticationSignup.Id,
		PermissionSysconsoleReadAuthenticationEmail.Id,
		PermissionSysconsoleReadAuthenticationPassword.Id,
		PermissionSysconsoleReadAuthenticationMfa.Id,
		PermissionSysconsoleReadAuthenticationLdap.Id,
		PermissionSysconsoleReadAuthenticationSaml.Id,
		PermissionSysconsoleReadAuthenticationOpenid.Id,
		PermissionSysconsoleReadAuthenticationGuestAccess.Id,
		PermissionSysconsoleReadPlugins.Id,
		PermissionSysconsoleReadIntegrationsIntegrationManagement.Id,
		PermissionSysconsoleReadIntegrationsBotAccounts.Id,
		PermissionSysconsoleReadIntegrationsGif.Id,
		PermissionSysconsoleReadIntegrationsCors.Id,
		PermissionSysconsoleWriteIntegrationsIntegrationManagement.Id,
		PermissionSysconsoleWriteIntegrationsBotAccounts.Id,
		PermissionSysconsoleWriteIntegrationsGif.Id,
		PermissionSysconsoleWriteIntegrationsCors.Id,

		PermissionCreateAttribute.Id,
		PermissionReadAttribute.Id,
		PermissionDeleteAttribute.Id,
		PermissionUpdateAttributeValue.Id,
		PermissionDeleteAttributeValue.Id,
		PermissionCreateAttributeValue.Id,
	}

	// Add the ancillary permissions to each system role
	SystemUserManagerDefaultPermissions = AddAncillaryPermissions(SystemUserManagerDefaultPermissions)
	SystemReadOnlyAdminDefaultPermissions = AddAncillaryPermissions(SystemReadOnlyAdminDefaultPermissions)
	SystemManagerDefaultPermissions = AddAncillaryPermissions(SystemManagerDefaultPermissions)
}
