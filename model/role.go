package model

import (
	"strings"
)

var (
	SysconsoleAncillaryPermissions        map[string][]string // SysconsoleAncillaryPermissions maps the non-sysconsole permissions required by each sysconsole view.
	SystemManagerDefaultPermissions       []string
	SystemUserManagerDefaultPermissions   []string
	SystemReadOnlyAdminDefaultPermissions []string
	BuiltInSchemeManagedRoleIDs           []string
	NewSystemRoleIDs                      []string
)

const (
	SystemUserRoleId            = "system_user"
	SystemAdminRoleId           = "system_admin"
	SystemUserAccessTokenRoleId = "system_user_access_token"
	SystemUserManagerRoleId     = "system_user_manager"
	SystemReadOnlyAdminRoleId   = "system_read_only_admin"
	SystemManagerRoleId         = "system_manager"
	ShopAdminRoleId             = "shop_admin"
	ShopStaffRoleId             = "shop_staff"

	RoleNameMaxLength        = 64
	RoleDisplayNameMaxLength = 128
	RoleDescriptionMaxLength = 1024
)

func init() {
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

type Role struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	DisplayName   string   `json:"display_name"`
	Description   string   `json:"description"`
	CreateAt      int64    `json:"create_at"`
	UpdateAt      int64    `json:"update_at"`
	DeleteAt      int64    `json:"delete_at"`
	Permissions   []string `json:"permissions"`
	SchemeManaged bool     `json:"scheme_managed"`
	BuiltIn       bool     `json:"built_in"`
}

type RolePatch struct {
	Permissions *[]string `json:"permissions"`
}

type RolePermissions struct {
	RoleID      string
	Permissions []string
}

// PreSave set `Id`, `CreateAt` and `PermissionsStr`
func (r *Role) PreSave() {
	if r.Id == "" {
		r.Id = NewId()
	}
	r.CreateAt = GetMillis()
	r.UpdateAt = r.CreateAt
}

// PreUpdate set `UpdateAt` and `PermissionsStr`
func (r *Role) PreUpdate() {
	r.UpdateAt = GetMillis()
}

func (r *Role) ToJSON() string {
	return ModelToJson(r)
}

func RoleListToJson(r []*Role) string {
	return ModelToJson(r)
}

func (r *RolePatch) ToJSON() string {
	return ModelToJson(r)
}

func (r *Role) Patch(patch *RolePatch) {
	if patch.Permissions != nil {
		r.Permissions = *patch.Permissions
	}
}

// Returns an array of permissions that are in either role.Permissions
// or patch.Permissions, but not both.
func PermissionsChangedByPatch(role *Role, patch *RolePatch) []string {
	var result []string

	if patch.Permissions == nil {
		return result
	}

	roleMap := make(map[string]bool)
	patchMap := make(map[string]bool)

	for _, permission := range role.Permissions {
		roleMap[permission] = true
	}
	for _, permission := range *patch.Permissions {
		patchMap[permission] = true
	}
	for _, permission := range role.Permissions {
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

func (r *Role) IsValid() bool {
	if !IsValidId(r.Id) {
		return false
	}

	return r.IsValidWithoutId()
}

// IsValidWithoutId check if current role is valid without checking its Id
func (r *Role) IsValidWithoutId() bool {
	if !IsValidRoleName(r.Name) {
		return false
	}

	if r.DisplayName == "" || len(r.DisplayName) > RoleDisplayNameMaxLength {
		return false
	}

	if len(r.Description) > RoleDescriptionMaxLength {
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

	for _, permissionId := range r.Permissions {
		if !check(AllSystemScopedPermissions, permissionId) {
			return false
		}
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

// IsValidRoleName checks if roleName's length is valid ( > 0 && <= 64) AND
//
// contains no character other than English's 26 letters, 10 digits and underscore
func IsValidRoleName(roleName string) bool {
	if roleName == "" || len(roleName) > RoleNameMaxLength {
		return false
	}

	if strings.TrimLeft(roleName, "abcdefghijklmnopqrstuvwxyz0123456789_") != "" {
		return false
	}

	return true
}

// MakeDefaultRoles make an map with values are default roles
func MakeDefaultRoles() map[string]*Role {
	roles := make(map[string]*Role)

	roles[SystemUserRoleId] = &Role{
		Name:          SystemUserRoleId,
		DisplayName:   "authentication.roles.global_user.name",
		Description:   "authentication.roles.global_user.description",
		Permissions:   SystemUserPermissions.IDs(),
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[SystemUserAccessTokenRoleId] = &Role{
		Name:        SystemUserAccessTokenRoleId,
		DisplayName: "authentication.roles.system_user_access_token.name",
		Description: "authentication.roles.system_user_access_token.description",
		Permissions: []string{
			PermissionCreateUserAccessToken.Id,
			PermissionReadUserAccessToken.Id,
			PermissionRevokeUserAccessToken.Id,
		},
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemUserManagerRoleId] = &Role{
		Name:          SystemUserManagerRoleId,
		DisplayName:   "authentication.roles.system_user_manager.name",
		Description:   "authentication.roles.system_user_manager.description",
		Permissions:   SystemUserManagerDefaultPermissions,
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemReadOnlyAdminRoleId] = &Role{
		Name:          SystemReadOnlyAdminRoleId,
		DisplayName:   "authentication.roles.system_read_only_admin.name",
		Description:   "authentication.roles.system_read_only_admin.description",
		Permissions:   SystemReadOnlyAdminDefaultPermissions,
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemManagerRoleId] = &Role{
		Name:          SystemManagerRoleId,
		DisplayName:   "authentication.roles.system_manager.name",
		Description:   "authentication.roles.system_manager.description",
		Permissions:   SystemManagerDefaultPermissions,
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemAdminRoleId] = &Role{
		Name:          SystemAdminRoleId,
		DisplayName:   "authentication.roles.global_admin.name",
		Description:   "authentication.roles.global_admin.description",
		Permissions:   AllSystemScopedPermissions.IDs(),
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[ShopAdminRoleId] = &Role{
		Name:          ShopAdminRoleId,
		DisplayName:   "authentication.roles.shop_admin.name",
		Description:   "authentication.roles.shop_admin.description",
		Permissions:   ShopScopedAllPermissions.IDs(),
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[ShopStaffRoleId] = &Role{
		Name:          ShopStaffRoleId,
		DisplayName:   "authentication.roles.shop_staff.name",
		Description:   "authentication.roles.shop_staff.description",
		Permissions:   ShopStaffPermissions.IDs(),
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
