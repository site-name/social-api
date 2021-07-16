package app

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type permissionTransformation struct {
	On     func(*model.Role, map[string]map[string]bool) bool
	Add    []string
	Remove []string
}
type permissionsMap []permissionTransformation

const (
	PermissionManageSystem                 = "manage_system"
	PermissionManageWebhooks               = "manage_webhooks"
	PermissionManageOthersWebhooks         = "manage_others_webhooks"
	PermissionManageIncomingWebhooks       = "manage_incoming_webhooks"
	PermissionManageOthersIncomingWebhooks = "manage_others_incoming_webhooks"
	PermissionManageOutgoingWebhooks       = "manage_outgoing_webhooks"
	PermissionManageOthersOutgoingWebhooks = "manage_others_outgoing_webhooks"
	PermissionPermanentDeleteUser          = "permanent_delete_user"
	PermissionViewMembers                  = "view_members"
	PermissionInviteUser                   = "invite_user"
	PermissionCreatePost                   = "create_post"
	PermissionCreatePost_PUBLIC            = "create_post_public"
	PermissionAddReaction                  = "add_reaction"
	PermissionRemoveReaction               = "remove_reaction"
	PermissionReadJobs                     = "read_jobs"
	PermissionManageJobs                   = "manage_jobs"
	PermissionEditOtherUsers               = "edit_other_users"
	PermissionEditBrand                    = "edit_brand"
	PermissionManageSecureConnections      = "manage_secure_connections"
)

// func isNotRole(roleName string) func(*model.Role, map[string]map[string]bool) bool {
// 	return func(role *model.Role, _ map[string]map[string]bool) bool {
// 		return role.Name != roleName
// 	}
// }
// func isNotSchemeRole(roleName string) func(*model.Role, map[string]map[string]bool) bool {
// 	return func(role *model.Role, _ map[string]map[string]bool) bool {
// 		return !strings.Contains(role.DisplayName, roleName)
// 	}
// }
// func permissionNotExists(permission string) func(*model.Role, map[string]map[string]bool) bool {
// 	return func(role *model.Role, permissionsMap map[string]map[string]bool) bool {
// 		val, ok := permissionsMap[role.Name][permission]
// 		return !(ok && val)
// 	}
// }
// func onOtherRole(otherRole string, function func(*model.Role, map[string]map[string]bool) bool) func(*model.Role, map[string]map[string]bool) bool {
// 	return func(role *model.Role, permissionsMap map[string]map[string]bool) bool {
// 		return function(&model.Role{Name: otherRole}, permissionsMap)
// 	}
// }
// func permissionAnd(funcs ...func(*model.Role, map[string]map[string]bool) bool) func(*model.Role, map[string]map[string]bool) bool {
// 	return func(role *model.Role, permissionsMap map[string]map[string]bool) bool {
// 		for _, f := range funcs {
// 			if !f(role, permissionsMap) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// }

func permissionExists(permissionID string) func(*model.Role, map[string]map[string]bool) bool {
	return func(role *model.Role, permissionsMap map[string]map[string]bool) bool {
		val, ok := permissionsMap[role.Name][permissionID]
		return ok && val
	}
}

func permissionOr(funcs ...func(*model.Role, map[string]map[string]bool) bool) func(*model.Role, map[string]map[string]bool) bool {
	return func(role *model.Role, permissionsMap map[string]map[string]bool) bool {
		for _, f := range funcs {
			if f(role, permissionsMap) {
				return true
			}
		}
		return false
	}
}

func isRole(roleName string) func(*model.Role, map[string]map[string]bool) bool {
	return func(role *model.Role, _ map[string]map[string]bool) bool {
		return role.Name == roleName
	}
}

// applyPermissionsMap
func applyPermissionsMap(role *model.Role, roleMap map[string]map[string]bool, migrationMap permissionsMap) *[]string {
	var result []string

	roleName := role.Name
	for _, transformation := range migrationMap {
		if transformation.On(role, roleMap) {
			for _, permission := range transformation.Add {
				roleMap[roleName][permission] = true
			}
			for _, permission := range transformation.Remove {
				roleMap[roleName][permission] = false
			}
		}
	}

	for key, active := range roleMap[roleName] {
		if active {
			result = append(result, key)
		}
	}
	return &result
}

// doPermissionsMigration
//
// roles: all roles available in system
func (s *Server) doPermissionsMigration(key string, migrationMap permissionsMap, roles []*model.Role) *model.AppError {
	if _, err := s.Store.System().GetByName(key); err == nil {
		return nil
	}

	roleMap := make(map[string]map[string]bool)
	for _, role := range roles {
		roleMap[role.Name] = make(map[string]bool)
		for _, permission := range role.Permissions {
			roleMap[role.Name][permission] = true
		}
	}

	for _, role := range roles {
		role.Permissions = *applyPermissionsMap(role, roleMap, migrationMap)
		if _, err := s.Store.Role().Save(role); err != nil {
			var invErr *store.ErrInvalidInput
			switch {
			case errors.As(err, &invErr):
				return model.NewAppError("doPermissionsMigration", "app.role.save.invalid_role.app_error", nil, invErr.Error(), http.StatusBadRequest)
			default:
				return model.NewAppError("doPermissionsMigration", "app.role.save.insert.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	if err := s.Store.System().Save(&model.System{Name: key, Value: "true"}); err != nil {
		return model.NewAppError("doPermissionsMigration", "app.system.save.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *App) getWebhooksPermissionsSplitMigration() (permissionsMap, error) {
	return permissionsMap{
		permissionTransformation{
			On:     permissionExists(PermissionManageWebhooks),
			Add:    []string{PermissionManageIncomingWebhooks, PermissionManageOutgoingWebhooks},
			Remove: []string{PermissionManageWebhooks},
		},
		permissionTransformation{
			On:     permissionExists(PermissionManageOthersWebhooks),
			Add:    []string{PermissionManageOthersIncomingWebhooks, PermissionManageOthersOutgoingWebhooks},
			Remove: []string{PermissionManageOthersWebhooks},
		},
	}, nil
}

func (a *App) removePermanentDeleteUserMigration() (permissionsMap, error) {
	return permissionsMap{
		permissionTransformation{
			On:     permissionExists(PermissionPermanentDeleteUser),
			Remove: []string{PermissionPermanentDeleteUser},
		},
	}, nil
}

func (a *App) getViewMembersPermissionMigration() (permissionsMap, error) {
	return permissionsMap{
		permissionTransformation{
			On:  isRole(model.SYSTEM_USER_ROLE_ID),
			Add: []string{PermissionViewMembers},
		},
		permissionTransformation{
			On:  isRole(model.SYSTEM_ADMIN_ROLE_ID),
			Add: []string{PermissionViewMembers},
		},
	}, nil
}

func (a *App) getAddSystemConsolePermissionsMigration() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	permissionsToAdd := []string{}
	for _, permission := range append(model.SysconsoleReadPermissions, model.SysconsoleWritePermissions...) {
		permissionsToAdd = append(permissionsToAdd, permission.Id)
	}

	// add the new permissions to system admin
	transformations = append(transformations,
		permissionTransformation{
			On:  isRole(model.SYSTEM_ADMIN_ROLE_ID),
			Add: permissionsToAdd,
		})

	// add read_jobs to all roles with manage_jobs
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(PermissionManageJobs),
		Add: []string{PermissionReadJobs},
	})

	// add read_other_users_teams to all roles with edit_other_users
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(PermissionEditOtherUsers),
		Add: []string{},
	})

	// add edit_brand to all roles with manage_system
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(PermissionManageSystem),
		Add: []string{PermissionEditBrand},
	})

	return transformations, nil
}

func (a *App) getSystemRolesPermissionsMigration() (permissionsMap, error) {
	return permissionsMap{
		permissionTransformation{
			On:  isRole(model.SYSTEM_ADMIN_ROLE_ID),
			Add: []string{model.PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_SYSTEM_ROLES.Id, model.PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_SYSTEM_ROLES.Id},
		},
	}, nil
}

func (a *App) getBillingPermissionsMigration() (permissionsMap, error) {
	return []permissionTransformation{
		{
			On:  isRole(model.SYSTEM_ADMIN_ROLE_ID),
			Add: []string{model.PERMISSION_SYSCONSOLE_READ_BILLING.Id, model.PERMISSION_SYSCONSOLE_WRITE_BILLING.Id},
		},
	}, nil
}

func (a *App) getAddManageSecureConnectionsPermissionsMigration() (permissionsMap, error) {
	transformations := []permissionTransformation{
		{ // add the new permission to system admin
			On:  isRole(model.SYSTEM_ADMIN_ROLE_ID),
			Add: []string{PermissionManageSecureConnections},
		},
		{ // remote the decprecated permission from system admin
			On:     isRole(model.SYSTEM_ADMIN_ROLE_ID),
			Remove: []string{PermissionManageSecureConnections},
		},
	}
	return transformations, nil
}

func (a *App) getAddDownloadComplianceExportResult() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	permissionsToAddComplianceRead := []string{model.PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT.Id, model.PERMISSION_READ_DATA_RETENTION_JOB.Id}
	permissionsToAddComplianceWrite := []string{model.PERMISSION_MANAGE_JOBS.Id}

	// add the new permissions to system admin
	transformations = append(transformations,
		permissionTransformation{
			On:  isRole(model.SYSTEM_ADMIN_ROLE_ID),
			Add: []string{model.PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT.Id},
		})

	// add Download Compliance Export Result and Read Jobs to all roles with sysconsole_read_compliance
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE.Id),
		Add: permissionsToAddComplianceRead,
	})

	// add manage_jobs to all roles with sysconsole_write_compliance
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE.Id),
		Add: permissionsToAddComplianceWrite,
	})

	return transformations, nil
}

func (a *App) getAddExperimentalSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	permissionsExperimentalRead := []string{model.PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_BLEVE.Id, model.PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURES.Id, model.PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURE_FLAGS.Id}
	permissionsExperimentalWrite := []string{model.PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_BLEVE.Id, model.PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURES.Id, model.PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURE_FLAGS.Id}

	// Give the new subsection READ permissions to any user with READ_EXPERIMENTAL
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL.Id),
		Add: permissionsExperimentalRead,
	})

	// Give the new subsection WRITE permissions to any user with WRITE_EXPERIMENTAL
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL.Id),
		Add: permissionsExperimentalWrite,
	})

	// Give the ancillary permissions MANAGE_JOBS and PURGE_BLEVE_INDEXES to anyone with WRITE_EXPERIMENTAL_BLEVE
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_BLEVE.Id),
		Add: []string{model.PERMISSION_CREATE_POST_BLEVE_INDEXES_JOB.Id, model.PERMISSION_PURGE_BLEVE_INDEXES.Id},
	})

	return transformations, nil
}

func (a *App) getAddIntegrationsSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	permissionsIntegrationsRead := []string{model.PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_INTEGRATION_MANAGEMENT.Id, model.PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_BOT_ACCOUNTS.Id, model.PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_GIF.Id, model.PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_CORS.Id}
	permissionsIntegrationsWrite := []string{model.PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_INTEGRATION_MANAGEMENT.Id, model.PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_BOT_ACCOUNTS.Id, model.PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_GIF.Id, model.PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_CORS.Id}

	// Give the new subsection READ permissions to any user with READ_INTEGRATIONS
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_INTEGRATIONS.Id),
		Add: permissionsIntegrationsRead,
	})

	// Give the new subsection WRITE permissions to any user with WRITE_EXPERIMENTAL
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS.Id),
		Add: permissionsIntegrationsWrite,
	})

	return transformations, nil
}

func (a *App) getAddSiteSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_SITE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_READ_SITE.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_READ_SITE_CUSTOMIZATION.Id,
			model.PERMISSION_SYSCONSOLE_READ_SITE_LOCALIZATION.Id,
			model.PERMISSION_SYSCONSOLE_READ_SITE_NOTIFICATIONS.Id,
			model.PERMISSION_SYSCONSOLE_READ_SITE_ANNOUNCEMENT_BANNER.Id,
			model.PERMISSION_SYSCONSOLE_READ_SITE_POSTS.Id,
			model.PERMISSION_SYSCONSOLE_READ_SITE_FILE_SHARING_AND_DOWNLOADS.Id,
			model.PERMISSION_SYSCONSOLE_READ_SITE_PUBLIC_LINKS.Id,
			model.PERMISSION_SYSCONSOLE_READ_SITE_NOTICES.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_SITE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_SITE.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_CUSTOMIZATION.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_LOCALIZATION.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_NOTIFICATIONS.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_ANNOUNCEMENT_BANNER.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_POSTS.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_FILE_SHARING_AND_DOWNLOADS.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_PUBLIC_LINKS.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_SITE_NOTICES.Id,
		},
	})

	// Give the ancillary permissions EDIT_BRAND to anyone with WRITE_SITE_CUSTOMIZATION
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_SITE_CUSTOMIZATION.Id),
		Add: []string{model.PERMISSION_EDIT_BRAND.Id},
	})

	return transformations, nil
}

func (a *App) getAddComplianceSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_COMPLIANCE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE_DATA_RETENTION_POLICY.Id,
			model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_EXPORT.Id,
			model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_MONITORING.Id,
			model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_COMPLIANCE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_DATA_RETENTION_POLICY.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_EXPORT.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_MONITORING.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE.Id,
		},
	})

	// Ancilary permissions
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_DATA_RETENTION_POLICY.Id),
		Add: []string{model.PERMISSION_CREATE_DATA_RETENTION_JOB.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE_DATA_RETENTION_POLICY.Id),
		Add: []string{model.PERMISSION_READ_DATA_RETENTION_JOB.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_EXPORT.Id),
		Add: []string{model.PERMISSION_CREATE_COMPLIANCE_EXPORT_JOB.Id, model.PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_EXPORT.Id),
		Add: []string{model.PERMISSION_READ_COMPLIANCE_EXPORT_JOB.Id, model.PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE.Id),
		Add: []string{model.PERMISSION_READ_AUDITS.Id},
	})

	return transformations, nil
}

func (a *App) getAddEnvironmentSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_ENVIRONMENT
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_WEB_SERVER.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DATABASE.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_FILE_STORAGE.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_IMAGE_PROXY.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SMTP.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PUSH_NOTIFICATION_SERVER.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_HIGH_AVAILABILITY.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_RATE_LIMITING.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_LOGGING.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SESSION_LENGTHS.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PERFORMANCE_MONITORING.Id,
			model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DEVELOPER.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_ENVIRONMENT
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_WEB_SERVER.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DATABASE.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_ELASTICSEARCH.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_FILE_STORAGE.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_IMAGE_PROXY.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SMTP.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PUSH_NOTIFICATION_SERVER.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_HIGH_AVAILABILITY.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_RATE_LIMITING.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_LOGGING.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SESSION_LENGTHS.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PERFORMANCE_MONITORING.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DEVELOPER.Id,
		},
	})

	// Give these ancillary permissions to anyone with READ_ENVIRONMENT_ELASTICSEARCH
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH.Id),
		Add: []string{
			model.PERMISSION_READ_ELASTICSEARCH_POST_INDEXING_JOB.Id,
			model.PERMISSION_READ_ELASTICSEARCH_POST_AGGREGATION_JOB.Id,
		},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_WEB_SERVER
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_WEB_SERVER.Id),
		Add: []string{
			model.PERMISSION_TEST_SITE_URL.Id,
			model.PERMISSION_RELOAD_CONFIG.Id,
			model.PERMISSION_INVALIDATE_CACHES.Id,
		},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_DATABASE
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DATABASE.Id),
		Add: []string{model.PERMISSION_RECYCLE_DATABASE_CONNECTIONS.Id},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_ELASTICSEARCH
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_ELASTICSEARCH.Id),
		Add: []string{
			model.PERMISSION_TEST_ELASTICSEARCH.Id,
			model.PERMISSION_CREATE_ELASTICSEARCH_POST_INDEXING_JOB.Id,
			model.PERMISSION_CREATE_ELASTICSEARCH_POST_AGGREGATION_JOB.Id,
			model.PERMISSION_PURGE_ELASTICSEARCH_INDEXES.Id,
		},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_FILE_STORAGE
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_FILE_STORAGE.Id),
		Add: []string{model.PERMISSION_TEST_S3.Id},
	})

	return transformations, nil
}

func (a *App) getAddAboutSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{
		permissionTransformation{
			On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_ABOUT.Id),
			Add: []string{},
		},
		permissionTransformation{
			On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_ABOUT.Id),
			Add: []string{},
		},
	}

	return transformations, nil
}

func (a *App) getAddReportingSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_REPORTING
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_READ_REPORTING.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS.Id,
			model.PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_REPORTING
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_REPORTING.Id),
		Add: []string{
			model.PERMISSION_SYSCONSOLE_WRITE_REPORTING_SITE_STATISTICS.Id,
			model.PERMISSION_SYSCONSOLE_WRITE_REPORTING_SERVER_LOGS.Id,
		},
	})

	// Give the ancillary permissions PERMISSION_GET_ANALYTICS to anyone with PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS or PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS
	transformations = append(transformations, permissionTransformation{
		On: permissionOr(
			permissionExists(model.PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS.Id),
			permissionExists(model.PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS.Id),
		),
		Add: []string{model.PERMISSION_GET_ANALYTICS.Id},
	})

	// Give the ancillary permissions PERMISSION_GET_LOGS to anyone with PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS.Id),
		Add: []string{model.PERMISSION_GET_LOGS.Id},
	})

	return transformations, nil
}

func (a *App) getAddAuthenticationSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{
		permissionTransformation{ // Give the new subsection READ permissions to any user with READ_AUTHENTICATION
			On: permissionExists(model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION.Id),
			Add: []string{
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SIGNUP.Id,
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_EMAIL.Id,
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_PASSWORD.Id,
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_MFA.Id,
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP.Id,
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SAML.Id,
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_OPENID.Id,
				model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_GUEST_ACCESS.Id,
			},
		},
		permissionTransformation{ // Give the new subsection WRITE permissions to any user with WRITE_AUTHENTICATION
			On: permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION.Id),
			Add: []string{
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SIGNUP.Id,
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_EMAIL.Id,
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_PASSWORD.Id,
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_MFA.Id,
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_LDAP.Id,
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SAML.Id,
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_OPENID.Id,
				model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_GUEST_ACCESS.Id,
			},
		},
		permissionTransformation{ // Give the ancillary permissions for LDAP to anyone with WRITE_AUTHENTICATION_LDAP
			On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_LDAP.Id),
			Add: []string{model.PERMISSION_CREATE_LDAP_SYNC_JOB.Id, model.PERMISSION_TEST_LDAP.Id, model.PERMISSION_ADD_LDAP_PUBLIC_CERT.Id, model.PERMISSION_ADD_LDAP_PRIVATE_CERT.Id, model.PERMISSION_REMOVE_LDAP_PUBLIC_CERT.Id, model.PERMISSION_REMOVE_LDAP_PRIVATE_CERT.Id},
		},
		permissionTransformation{ // Give the ancillary permissions PERMISSION_TEST_LDAP to anyone with READ_AUTHENTICATION_LDAP
			On:  permissionExists(model.PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP.Id),
			Add: []string{model.PERMISSION_READ_LDAP_SYNC_JOB.Id},
		},
		permissionTransformation{
			On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_EMAIL.Id),
			Add: []string{model.PERMISSION_INVALIDATE_EMAIL_INVITE.Id},
		},
		permissionTransformation{ // Give the ancillary permissions for SAML to anyone with WRITE_AUTHENTICATION_SAML
			On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SAML.Id),
			Add: []string{model.PERMISSION_GET_SAML_METADATA_FROM_IDP.Id, model.PERMISSION_ADD_SAML_PUBLIC_CERT.Id, model.PERMISSION_ADD_SAML_PRIVATE_CERT.Id, model.PERMISSION_ADD_SAML_IDP_CERT.Id, model.PERMISSION_REMOVE_SAML_PUBLIC_CERT.Id, model.PERMISSION_REMOVE_SAML_PRIVATE_CERT.Id, model.PERMISSION_REMOVE_SAML_IDP_CERT.Id, model.PERMISSION_GET_SAML_CERT_STATUS.Id},
		},
	}

	return transformations, nil
}

// This migration fixes https://github.com/mattermost/mattermost-server/issues/17642 where this particular ancillary permission was forgotten during the initial migrations
func (a *App) getAddTestEmailAncillaryPermission() (permissionsMap, error) {
	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_SMTP
	return []permissionTransformation{
		{
			On:  permissionExists(model.PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SMTP.Id),
			Add: []string{model.PERMISSION_TEST_EMAIL.Id},
		},
	}, nil
}

// DoPermissionsMigrations execute all the permissions migrations need by the current version.
func (a *App) DoPermissionsMigrations() error {
	return a.Srv().doPermissionsMigrations()
}

func (s *Server) doPermissionsMigrations() error {
	a := New(ServerConnector(s))

	PermissionsMigrations := []struct {
		Key       string
		Migration func() (permissionsMap, error)
	}{
		{Key: model.MIGRATION_KEY_WEBHOOK_PERMISSIONS_SPLIT, Migration: a.getWebhooksPermissionsSplitMigration},
		{Key: model.MIGRATION_KEY_REMOVE_PERMANENT_DELETE_USER, Migration: a.removePermanentDeleteUserMigration},
		{Key: model.MIGRATION_KEY_VIEW_MEMBERS_NEW_PERMISSION, Migration: a.getViewMembersPermissionMigration},
		{Key: model.MIGRATION_KEY_ADD_SYSTEM_CONSOLE_PERMISSIONS, Migration: a.getAddSystemConsolePermissionsMigration},
		{Key: model.MIGRATION_KEY_ADD_MANAGE_SECURE_CONNECTIONS_PERMISSIONS, Migration: a.getAddManageSecureConnectionsPermissionsMigration},
		{Key: model.MIGRATION_KEY_ADD_SYSTEM_ROLES_PERMISSIONS, Migration: a.getSystemRolesPermissionsMigration},
		{Key: model.MIGRATION_KEY_ADD_BILLING_PERMISSIONS, Migration: a.getBillingPermissionsMigration},
		{Key: model.MIGRATION_KEY_ADD_DOWNLOAD_COMPLIANCE_EXPORT_RESULTS, Migration: a.getAddDownloadComplianceExportResult},
		{Key: model.MIGRATION_KEY_ADD_EXPERIMENTAL_SUBSECTION_PERMISSIONS, Migration: a.getAddExperimentalSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_AUTHENTICATION_SUBSECTION_PERMISSIONS, Migration: a.getAddAuthenticationSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_INTEGRATIONS_SUBSECTION_PERMISSIONS, Migration: a.getAddIntegrationsSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_SITE_SUBSECTION_PERMISSIONS, Migration: a.getAddSiteSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_COMPLIANCE_SUBSECTION_PERMISSIONS, Migration: a.getAddComplianceSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_ENVIRONMENT_SUBSECTION_PERMISSIONS, Migration: a.getAddEnvironmentSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_ABOUT_SUBSECTION_PERMISSIONS, Migration: a.getAddAboutSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_REPORTING_SUBSECTION_PERMISSIONS, Migration: a.getAddReportingSubsectionPermissions},
		{Key: model.MIGRATION_KEY_ADD_TEST_EMAIL_ANCILLARY_PERMISSION, Migration: a.getAddTestEmailAncillaryPermission},
	}

	roles, err := s.Store.Role().GetAll()
	if err != nil {
		return err
	}

	for _, migration := range PermissionsMigrations {
		permissionMap, err := migration.Migration()
		if err != nil {
			return err
		}
		if err := s.doPermissionsMigration(migration.Key, permissionMap, roles); err != nil {
			return err
		}
	}
	return nil
}
