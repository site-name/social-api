package model

import (
	"io"
	"strings"
)

var (
	SysconsoleAncillaryPermissions        map[string][]*Permission // SysconsoleAncillaryPermissions maps the non-sysconsole permissions required by each sysconsole view.
	SystemManagerDefaultPermissions       []string
	SystemUserManagerDefaultPermissions   []string
	SystemReadOnlyAdminDefaultPermissions []string
	BuiltInSchemeManagedRoleIDs           []string
	NewSystemRoleIDs                      []string
)

type RoleType string
type RoleScope string

const (
	SystemGuestRoleId           = "system_guest"
	SystemUserRoleId            = "system_user"
	SystemAdminRoleId           = "system_admin"
	SystemPostAllRoleId         = "system_post_all"
	SystemPostAllPublicRoleId   = "system_post_all_public"
	SystemUserAccessTokenRoleId = "system_user_access_token"
	SystemUserManagerRoleId     = "system_user_manager"
	SystemReadOnlyAdminRoleId   = "system_read_only_admin"
	SystemManagerRoleId         = "system_manager"

	RoleNameMaxLength        = 64
	RoleDisplayNameMaxLength = 128
	RoleDescriptionMaxLength = 1024

	RoleScopeSystem RoleScope = "System"

	RoleTypeUser  RoleType = "User"
	RoleTypeAdmin RoleType = "Admin"
)

func init() {
	NewSystemRoleIDs = []string{
		SystemUserManagerRoleId,
		SystemReadOnlyAdminRoleId,
		SystemManagerRoleId,
	}

	BuiltInSchemeManagedRoleIDs = append([]string{
		SystemGuestRoleId,
		SystemUserRoleId,
		SystemAdminRoleId,
		SystemPostAllRoleId,
		SystemPostAllPublicRoleId,
		SystemUserAccessTokenRoleId,
	}, NewSystemRoleIDs...)

	// When updating the values here, the values in mattermost-redux must also be updated.
	SysconsoleAncillaryPermissions = map[string][]*Permission{
		PermissionSysconsoleReadUserManagementUsers.Id: {
			PermissionGetAnalytics,
		},
		PermissionSysconsoleReadEnvironmentElasticsearch.Id: {
			PermissionReadElasticsearchPostIndexingJob,
			PermissionReadElasticsearchPostAggregationJob,
		},
		PermissionSysconsoleWriteEnvironmentWebServer.Id: {
			PermissionTestSiteUrl,
			PermissionReloadConfig,
			PermissionInvalidateCaches,
		},
		PermissionSysconsoleWriteEnvironmentDatabase.Id: {
			PermissionRecycleDatabaseConnections,
		},
		PermissionSysconsoleWriteEnvironmentElasticsearch.Id: {
			PermissionTestElasticsearch,
			PermissionCreateElasticsearchPostIndexingJob,
			PermissionCreateElasticsearchPostAggregationJob,
			PermissionPurgeElasticsearchIndexes,
		},
		PermissionSysconsoleWriteEnvironmentFileStorage.Id: {
			PermissionTestS3,
		},
		PermissionSysconsoleWriteEnvironmentSmtp.Id: {
			PermissionTestEmail,
		},
		PermissionSysconsoleReadReportingServerLogs.Id: {
			PermissionGetLogs,
		},
		PermissionSysconsoleReadReportingSiteStatistics.Id: {
			PermissionGetAnalytics,
		},
		PermissionSysconsoleWriteUserManagementUsers.Id: {
			PermissionEditOtherUsers,
		},
		PermissionSysconsoleWriteSiteCustomization.Id: {
			PermissionEditBrand,
		},
		PermissionSysconsoleWriteComplianceDataRetentionPolicy.Id: {
			PermissionCreateDataRetentionJob,
		},
		PermissionSysconsoleReadComplianceDataRetentionPolicy.Id: {
			PermissionReadDataRetentionJob,
		},
		PermissionSysconsoleWriteComplianceComplianceExport.Id: {
			PermissionCreateComplianceExportJob,
			PermissionDownloadComplianceExportResult,
		},
		PermissionSysconsoleReadComplianceComplianceExport.Id: {
			PermissionReadComplianceExportJob,
			PermissionDownloadComplianceExportResult,
		},
		PermissionSysconsoleReadComplianceCustomTermsOfService.Id: {
			PermissionReadAudits,
		},
		PermissionSysconsoleWriteExperimentalBleve.Id: {
			PermissionCreatePostBleveIndexesJob,
			PermissionPurgeBleveIndexes,
		},
		PermissionSysconsoleWriteAuthenticationLdap.Id: {
			PermissionCreateLdapSyncJob,
			PermissionAddLdapPublicCert,
			PermissionRemoveLdapPublicCert,
			PermissionAddLdapPrivateCert,
			PermissionRemoveLdapPrivateCert,
		},
		PermissionSysconsoleReadAuthenticationLdap.Id: {
			PermissionTestLdap,
			PermissionReadLdapSyncJob,
		},
		PermissionSysconsoleWriteAuthenticationEmail.Id: {
			PermissionInvalidateEmailInvite,
		},
		PermissionSysconsoleWriteAuthenticationSaml.Id: {
			PermissionGetSamlMetadataFromIdp,
			PermissionAddSamlPublicCert,
			PermissionAddSamlPrivateCert,
			PermissionAddSamlIdpCert,
			PermissionRemoveSamlPublicCert,
			PermissionRemoveSamlPrivateCert,
			PermissionRemoveSamlIdpCert,
			PermissionGetSamlCertStatus,
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
	}

	// Add the ancillary permissions to each system role
	SystemUserManagerDefaultPermissions = AddAncillaryPermissions(SystemUserManagerDefaultPermissions)
	SystemReadOnlyAdminDefaultPermissions = AddAncillaryPermissions(SystemReadOnlyAdminDefaultPermissions)
	SystemManagerDefaultPermissions = AddAncillaryPermissions(SystemManagerDefaultPermissions)
}

type Role struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	DisplayName    string   `json:"display_name"`
	Description    string   `json:"description"`
	CreateAt       int64    `json:"create_at"`
	UpdateAt       int64    `json:"update_at"`
	DeleteAt       int64    `json:"delete_at"`
	Permissions    []string `json:"permissions" db:"-"`       // NOT save to database, populate by `PermissionsStr`
	PermissionsStr string   `json:"permission_str,omitempty"` // save to database
	SchemeManaged  bool     `json:"scheme_managed"`
	BuiltIn        bool     `json:"built_in"`
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
	r.PermissionsStr = strings.Join(r.Permissions, " ")
}

// PreUpdate set `UpdateAt` and `PermissionsStr`
func (r *Role) PreUpdate() {
	r.UpdateAt = GetMillis()
	r.PermissionsStr = strings.Join(r.Permissions, " ")
}

// PopulatePermissionSlice populates role's Permissions slice
func (r *Role) PopulatePermissionSlice() {
	r.Permissions = strings.Fields(r.PermissionsStr)
	r.PermissionsStr = ""
}

func (r *Role) ToJSON() string {
	return ModelToJson(r)
}

func RoleFromJson(data io.Reader) *Role {
	var r *Role
	ModelFromJson(&r, data)
	return r
}

func RoleListToJson(r []*Role) string {
	return ModelToJson(r)
}

func RoleListFromJson(data io.Reader) []*Role {
	var roles []*Role
	ModelFromJson(&roles, data)
	return roles
}

func (r *RolePatch) ToJSON() string {
	return ModelToJson(r)
}

func RolePatchFromJson(data io.Reader) *RolePatch {
	var rolePatch *RolePatch
	ModelFromJson(&rolePatch, data)
	return rolePatch
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

	ok, _, _ := r.IsValidWithoutId()
	return ok
}

// IsValidWithoutId check if current role is valid without checking its Id
func (r *Role) IsValidWithoutId() (bool, string, interface{}) {
	if !IsValidRoleName(r.Name) {
		return false, "Name", r.Name
	}

	if r.DisplayName == "" || len(r.DisplayName) > RoleDisplayNameMaxLength {
		return false, "DisplayName", r.DisplayName
	}

	if len(r.Description) > RoleDescriptionMaxLength {
		return false, "Description", r.Description
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
		permissionValidated := check(AllPermissions, permissionId) || check(DeprecatedPermissions, permissionId)
		if !permissionValidated {
			return false, "Permissions", r.Permissions
		}
	}

	return true, "nil", nil
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

	roles[SystemGuestRoleId] = &Role{
		Name:          SystemGuestRoleId,
		DisplayName:   "authentication.roles.global_guest.name",
		Description:   "authentication.roles.global_guest.description",
		Permissions:   []string{},
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[SystemUserRoleId] = &Role{
		Name:        SystemUserRoleId,
		DisplayName: "authentication.roles.global_user.name",
		Description: "authentication.roles.global_user.description",
		Permissions: []string{
			PermissionViewMembers.Id,
		},
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[SystemPostAllRoleId] = &Role{
		Name:        SystemPostAllRoleId,
		DisplayName: "authentication.roles.system_post_all.name",
		Description: "authentication.roles.system_post_all.description",
		Permissions: []string{
			PermissionCreatePost.Id,
		},
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SystemPostAllPublicRoleId] = &Role{
		Name:        SystemPostAllPublicRoleId,
		DisplayName: "authentication.roles.system_post_all_public.name",
		Description: "authentication.roles.system_post_all_public.description",
		Permissions: []string{
			PermissionCreatePostPublic.Id,
		},
		SchemeManaged: false,
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

	allPermissionIDs := []string{}
	for _, permission := range AllPermissions {
		allPermissionIDs = append(allPermissionIDs, permission.Id)
	}

	roles[SystemAdminRoleId] = &Role{
		Name:          SystemAdminRoleId,
		DisplayName:   "authentication.roles.global_admin.name",
		Description:   "authentication.roles.global_admin.description",
		Permissions:   allPermissionIDs,
		SchemeManaged: true,
		BuiltIn:       true,
	}

	return roles
}

func AddAncillaryPermissions(permissions []string) []string {
	for _, permission := range permissions {
		if ancillaryPermissions, ok := SysconsoleAncillaryPermissions[permission]; ok {
			for _, ancillaryPermission := range ancillaryPermissions {
				permissions = append(permissions, ancillaryPermission.Id)
			}
		}
	}
	return permissions
}
