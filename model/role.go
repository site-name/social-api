package model

import (
	"io"
	"strings"
)

// SysconsoleAncillaryPermissions maps the non-sysconsole permissions required by each sysconsole view.
var SysconsoleAncillaryPermissions map[string][]*Permission
var SystemManagerDefaultPermissions []string
var SystemUserManagerDefaultPermissions []string
var SystemReadOnlyAdminDefaultPermissions []string

var BuiltInSchemeManagedRoleIDs []string

var NewSystemRoleIDs []string

func init() {
	NewSystemRoleIDs = []string{
		SYSTEM_USER_MANAGER_ROLE_ID,
		SYSTEM_READ_ONLY_ADMIN_ROLE_ID,
		SYSTEM_MANAGER_ROLE_ID,
	}

	BuiltInSchemeManagedRoleIDs = append([]string{
		SYSTEM_GUEST_ROLE_ID,
		SYSTEM_USER_ROLE_ID,
		SYSTEM_ADMIN_ROLE_ID,
		SYSTEM_POST_ALL_ROLE_ID,
		SYSTEM_POST_ALL_PUBLIC_ROLE_ID,
		SYSTEM_USER_ACCESS_TOKEN_ROLE_ID,
	}, NewSystemRoleIDs...)

	// When updating the values here, the values in mattermost-redux must also be updated.
	SysconsoleAncillaryPermissions = map[string][]*Permission{
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS.Id: {
			PERMISSION_GET_ANALYTICS,
		},
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH.Id: {
			PERMISSION_READ_ELASTICSEARCH_POST_INDEXING_JOB,
			PERMISSION_READ_ELASTICSEARCH_POST_AGGREGATION_JOB,
		},
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_WEB_SERVER.Id: {
			PERMISSION_TEST_SITE_URL,
			PERMISSION_RELOAD_CONFIG,
			PERMISSION_INVALIDATE_CACHES,
		},
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DATABASE.Id: {
			PERMISSION_RECYCLE_DATABASE_CONNECTIONS,
		},
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_ELASTICSEARCH.Id: {
			PERMISSION_TEST_ELASTICSEARCH,
			PERMISSION_CREATE_ELASTICSEARCH_POST_INDEXING_JOB,
			PERMISSION_CREATE_ELASTICSEARCH_POST_AGGREGATION_JOB,
			PERMISSION_PURGE_ELASTICSEARCH_INDEXES,
		},
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_FILE_STORAGE.Id: {
			PERMISSION_TEST_S3,
		},
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SMTP.Id: {
			PERMISSION_TEST_EMAIL,
		},
		PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS.Id: {
			PERMISSION_GET_LOGS,
		},
		PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS.Id: {
			PERMISSION_GET_ANALYTICS,
		},
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_USERS.Id: {
			PERMISSION_EDIT_OTHER_USERS,
		},
		PERMISSION_SYSCONSOLE_WRITE_SITE_CUSTOMIZATION.Id: {
			PERMISSION_EDIT_BRAND,
		},
		PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_DATA_RETENTION_POLICY.Id: {
			PERMISSION_CREATE_DATA_RETENTION_JOB,
		},
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_DATA_RETENTION_POLICY.Id: {
			PERMISSION_READ_DATA_RETENTION_JOB,
		},
		PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_EXPORT.Id: {
			PERMISSION_CREATE_COMPLIANCE_EXPORT_JOB,
			PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT,
		},
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_EXPORT.Id: {
			PERMISSION_READ_COMPLIANCE_EXPORT_JOB,
			PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT,
		},
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE.Id: {
			PERMISSION_READ_AUDITS,
		},
		PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_BLEVE.Id: {
			PERMISSION_CREATE_POST_BLEVE_INDEXES_JOB,
			PERMISSION_PURGE_BLEVE_INDEXES,
		},
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_LDAP.Id: {
			PERMISSION_CREATE_LDAP_SYNC_JOB,
			PERMISSION_ADD_LDAP_PUBLIC_CERT,
			PERMISSION_REMOVE_LDAP_PUBLIC_CERT,
			PERMISSION_ADD_LDAP_PRIVATE_CERT,
			PERMISSION_REMOVE_LDAP_PRIVATE_CERT,
		},
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP.Id: {
			PERMISSION_TEST_LDAP,
			PERMISSION_READ_LDAP_SYNC_JOB,
		},
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_EMAIL.Id: {
			PERMISSION_INVALIDATE_EMAIL_INVITE,
		},
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SAML.Id: {
			PERMISSION_GET_SAML_METADATA_FROM_IDP,
			PERMISSION_ADD_SAML_PUBLIC_CERT,
			PERMISSION_ADD_SAML_PRIVATE_CERT,
			PERMISSION_ADD_SAML_IDP_CERT,
			PERMISSION_REMOVE_SAML_PUBLIC_CERT,
			PERMISSION_REMOVE_SAML_PRIVATE_CERT,
			PERMISSION_REMOVE_SAML_IDP_CERT,
			PERMISSION_GET_SAML_CERT_STATUS,
		},
	}

	SystemUserManagerDefaultPermissions = []string{
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_GROUPS.Id,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_PERMISSIONS.Id,
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_GROUPS.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SIGNUP.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_EMAIL.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_PASSWORD.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_MFA.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SAML.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_OPENID.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_GUEST_ACCESS.Id,
	}

	SystemReadOnlyAdminDefaultPermissions = []string{
		PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS.Id,
		PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS.Id,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS.Id,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_GROUPS.Id,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_PERMISSIONS.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_WEB_SERVER.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DATABASE.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_FILE_STORAGE.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_IMAGE_PROXY.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SMTP.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PUSH_NOTIFICATION_SERVER.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_HIGH_AVAILABILITY.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_RATE_LIMITING.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_LOGGING.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SESSION_LENGTHS.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PERFORMANCE_MONITORING.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DEVELOPER.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_CUSTOMIZATION.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_LOCALIZATION.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_NOTIFICATIONS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_ANNOUNCEMENT_BANNER.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_POSTS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_FILE_SHARING_AND_DOWNLOADS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_PUBLIC_LINKS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_NOTICES.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SIGNUP.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_EMAIL.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_PASSWORD.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_MFA.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SAML.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_OPENID.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_GUEST_ACCESS.Id,
		PERMISSION_SYSCONSOLE_READ_PLUGINS.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_INTEGRATION_MANAGEMENT.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_BOT_ACCOUNTS.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_GIF.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_CORS.Id,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_DATA_RETENTION_POLICY.Id,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_EXPORT.Id,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_MONITORING.Id,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE.Id,
		PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURES.Id,
		PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURE_FLAGS.Id,
		PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_BLEVE.Id,
	}

	SystemManagerDefaultPermissions = []string{
		PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS.Id,
		PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS.Id,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_GROUPS.Id,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_PERMISSIONS.Id,
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_GROUPS.Id,
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_PERMISSIONS.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_WEB_SERVER.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DATABASE.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_FILE_STORAGE.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_IMAGE_PROXY.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SMTP.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PUSH_NOTIFICATION_SERVER.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_HIGH_AVAILABILITY.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_RATE_LIMITING.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_LOGGING.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SESSION_LENGTHS.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PERFORMANCE_MONITORING.Id,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DEVELOPER.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_WEB_SERVER.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DATABASE.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_ELASTICSEARCH.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_FILE_STORAGE.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_IMAGE_PROXY.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SMTP.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PUSH_NOTIFICATION_SERVER.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_HIGH_AVAILABILITY.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_RATE_LIMITING.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_LOGGING.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SESSION_LENGTHS.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PERFORMANCE_MONITORING.Id,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DEVELOPER.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_CUSTOMIZATION.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_CUSTOMIZATION.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_LOCALIZATION.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_LOCALIZATION.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_NOTIFICATIONS.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_NOTIFICATIONS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_ANNOUNCEMENT_BANNER.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_ANNOUNCEMENT_BANNER.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_POSTS.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_POSTS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_FILE_SHARING_AND_DOWNLOADS.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_FILE_SHARING_AND_DOWNLOADS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_PUBLIC_LINKS.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_PUBLIC_LINKS.Id,
		PERMISSION_SYSCONSOLE_READ_SITE_NOTICES.Id,
		PERMISSION_SYSCONSOLE_WRITE_SITE_NOTICES.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SIGNUP.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_EMAIL.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_PASSWORD.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_MFA.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SAML.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_OPENID.Id,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_GUEST_ACCESS.Id,
		PERMISSION_SYSCONSOLE_READ_PLUGINS.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_INTEGRATION_MANAGEMENT.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_BOT_ACCOUNTS.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_GIF.Id,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_CORS.Id,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_INTEGRATION_MANAGEMENT.Id,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_BOT_ACCOUNTS.Id,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_GIF.Id,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_CORS.Id,
	}

	// Add the ancillary permissions to each system role
	SystemUserManagerDefaultPermissions = AddAncillaryPermissions(SystemUserManagerDefaultPermissions)
	SystemReadOnlyAdminDefaultPermissions = AddAncillaryPermissions(SystemReadOnlyAdminDefaultPermissions)
	SystemManagerDefaultPermissions = AddAncillaryPermissions(SystemManagerDefaultPermissions)
}

type RoleType string
type RoleScope string

const (
	SYSTEM_GUEST_ROLE_ID             = "system_guest"
	SYSTEM_USER_ROLE_ID              = "system_user"
	SYSTEM_ADMIN_ROLE_ID             = "system_admin"
	SYSTEM_POST_ALL_ROLE_ID          = "system_post_all"
	SYSTEM_POST_ALL_PUBLIC_ROLE_ID   = "system_post_all_public"
	SYSTEM_USER_ACCESS_TOKEN_ROLE_ID = "system_user_access_token"
	SYSTEM_USER_MANAGER_ROLE_ID      = "system_user_manager"
	SYSTEM_READ_ONLY_ADMIN_ROLE_ID   = "system_read_only_admin"
	SYSTEM_MANAGER_ROLE_ID           = "system_manager"

	ROLE_NAME_MAX_LENGTH         = 64
	ROLE_DISPLAY_NAME_MAX_LENGTH = 128
	ROLE_DESCRIPTION_MAX_LENGTH  = 1024

	RoleScopeSystem RoleScope = "System"

	RoleTypeUser  RoleType = "User"
	RoleTypeAdmin RoleType = "Admin"
)

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

func (r *Role) ToJson() string {
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

func (r *RolePatch) ToJson() string {
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

// MergeChannelHigherScopedPermissions is meant to be invoked on a channel scheme's role and merges the higher-scoped
// channel role's permissions.
// func (r *Role) MergeChannelHigherScopedPermissions(higherScopedPermissions *RolePermissions) {
// 	mergedPermissions := []string{}

// 	higherScopedPermissionsMap := AsStringBoolMap(higherScopedPermissions.Permissions)
// 	rolePermissionsMap := AsStringBoolMap(r.Permissions)

// 	for _, cp := range AllPermissions {
// 		if cp.Scope != PermissionScopeChannel {
// 			continue
// 		}

// 		_, presentOnHigherScope := higherScopedPermissionsMap[cp.Id]

// 		// For the channel admin role always look to the higher scope to determine if the role has their permission.
// 		// The channel admin is a special case because they're not part of the UI to be "channel moderated", only
// 		// channel members and channel guests are.
// 		if higherScopedPermissions.RoleID == CHANNEL_ADMIN_ROLE_ID && presentOnHigherScope {
// 			mergedPermissions = append(mergedPermissions, cp.Id)
// 			continue
// 		}

// 		_, permissionIsModerated := ChannelModeratedPermissionsMap[cp.Id]
// 		if permissionIsModerated {
// 			_, presentOnRole := rolePermissionsMap[cp.Id]
// 			if presentOnRole && presentOnHigherScope {
// 				mergedPermissions = append(mergedPermissions, cp.Id)
// 			}
// 		} else {
// 			if presentOnHigherScope {
// 				mergedPermissions = append(mergedPermissions, cp.Id)
// 			}
// 		}
// 	}

// 	r.Permissions = mergedPermissions
// }

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

	if r.DisplayName == "" || len(r.DisplayName) > ROLE_DISPLAY_NAME_MAX_LENGTH {
		return false, "DisplayName", r.DisplayName
	}

	if len(r.Description) > ROLE_DESCRIPTION_MAX_LENGTH {
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
	if roleName == "" || len(roleName) > ROLE_NAME_MAX_LENGTH {
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

	roles[SYSTEM_GUEST_ROLE_ID] = &Role{
		Name:          "system_guest",
		DisplayName:   "authentication.roles.global_guest.name",
		Description:   "authentication.roles.global_guest.description",
		Permissions:   []string{},
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[SYSTEM_USER_ROLE_ID] = &Role{
		Name:        "system_user",
		DisplayName: "authentication.roles.global_user.name",
		Description: "authentication.roles.global_user.description",
		Permissions: []string{
			PERMISSION_VIEW_MEMBERS.Id,
		},
		SchemeManaged: true,
		BuiltIn:       true,
	}

	roles[SYSTEM_POST_ALL_ROLE_ID] = &Role{
		Name:        "system_post_all",
		DisplayName: "authentication.roles.system_post_all.name",
		Description: "authentication.roles.system_post_all.description",
		Permissions: []string{
			PERMISSION_CREATE_POST.Id,
		},
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SYSTEM_POST_ALL_PUBLIC_ROLE_ID] = &Role{
		Name:        "system_post_all_public",
		DisplayName: "authentication.roles.system_post_all_public.name",
		Description: "authentication.roles.system_post_all_public.description",
		Permissions: []string{
			PERMISSION_CREATE_POST_PUBLIC.Id,
		},
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SYSTEM_USER_ACCESS_TOKEN_ROLE_ID] = &Role{
		Name:        "system_user_access_token",
		DisplayName: "authentication.roles.system_user_access_token.name",
		Description: "authentication.roles.system_user_access_token.description",
		Permissions: []string{
			PERMISSION_CREATE_USER_ACCESS_TOKEN.Id,
			PERMISSION_READ_USER_ACCESS_TOKEN.Id,
			PERMISSION_REVOKE_USER_ACCESS_TOKEN.Id,
		},
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SYSTEM_USER_MANAGER_ROLE_ID] = &Role{
		Name:          "system_user_manager",
		DisplayName:   "authentication.roles.system_user_manager.name",
		Description:   "authentication.roles.system_user_manager.description",
		Permissions:   SystemUserManagerDefaultPermissions,
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SYSTEM_READ_ONLY_ADMIN_ROLE_ID] = &Role{
		Name:          "system_read_only_admin",
		DisplayName:   "authentication.roles.system_read_only_admin.name",
		Description:   "authentication.roles.system_read_only_admin.description",
		Permissions:   SystemReadOnlyAdminDefaultPermissions,
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SYSTEM_MANAGER_ROLE_ID] = &Role{
		Name:          "system_manager",
		DisplayName:   "authentication.roles.system_manager.name",
		Description:   "authentication.roles.system_manager.description",
		Permissions:   SystemManagerDefaultPermissions,
		SchemeManaged: false,
		BuiltIn:       true,
	}

	roles[SYSTEM_ADMIN_ROLE_ID] = &Role{
		Name:          "system_admin",
		DisplayName:   "authentication.roles.global_admin.name",
		Description:   "authentication.roles.global_admin.description",
		Permissions:   []string{}, // system admin can do every thing
		SchemeManaged: true,
		BuiltIn:       true,
	}

	for _, permission := range AllPermissions {
		roles[SYSTEM_ADMIN_ROLE_ID].Permissions = append(roles[SYSTEM_ADMIN_ROLE_ID].Permissions, permission.Id)
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
