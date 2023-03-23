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
func applyPermissionsMap(role *model.Role, roleMap map[string]map[string]bool, migrationMap permissionsMap) []string {
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
	return result
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
		role.Permissions = applyPermissionsMap(role, roleMap, migrationMap)
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

	if err := s.Store.System().SaveOrUpdate(&model.System{Name: key, Value: "true"}); err != nil {
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
			On:  isRole(model.SystemUserRoleId),
			Add: []string{PermissionViewMembers},
		},
		permissionTransformation{
			On:  isRole(model.SystemAdminRoleId),
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
			On:  isRole(model.SystemAdminRoleId),
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
			On:  isRole(model.SystemAdminRoleId),
			Add: []string{model.PermissionSysconsoleReadUserManagementSystemRoles.Id, model.PermissionSysconsoleWriteUserManagementSystemRoles.Id},
		},
	}, nil
}

func (a *App) getBillingPermissionsMigration() (permissionsMap, error) {
	return []permissionTransformation{
		{
			On:  isRole(model.SystemAdminRoleId),
			Add: []string{model.PermissionSysconsoleReadBilling.Id, model.PermissionSysconsoleWriteBilling.Id},
		},
	}, nil
}

func (a *App) getAddManageSecureConnectionsPermissionsMigration() (permissionsMap, error) {
	transformations := []permissionTransformation{
		{
			On:  isRole(model.SystemAdminRoleId),
			Add: []string{PermissionManageSecureConnections},
		},
		{
			On:     isRole(model.SystemAdminRoleId),
			Remove: []string{PermissionManageSecureConnections},
		},
	}
	return transformations, nil
}

func (a *App) getAddDownloadComplianceExportResult() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	permissionsToAddComplianceRead := []string{model.PermissionDownloadComplianceExportResult.Id, model.PermissionReadDataRetentionJob.Id}
	permissionsToAddComplianceWrite := []string{model.PermissionManageJobs.Id}

	// add the new permissions to system admin
	transformations = append(transformations,
		permissionTransformation{
			On:  isRole(model.SystemAdminRoleId),
			Add: []string{model.PermissionDownloadComplianceExportResult.Id},
		})

	// add Download Compliance Export Result and Read Jobs to all roles with sysconsole_read_compliance
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleReadCompliance.Id),
		Add: permissionsToAddComplianceRead,
	})

	// add manage_jobs to all roles with sysconsole_write_compliance
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteCompliance.Id),
		Add: permissionsToAddComplianceWrite,
	})

	return transformations, nil
}

func (a *App) getAddExperimentalSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	permissionsExperimentalRead := []string{
		model.PermissionSysconsoleReadExperimentalBleve.Id,
		model.PermissionSysconsoleReadExperimentalFeatures.Id,
		model.PermissionSysconsoleReadExperimentalFeatureFlags.Id,
	}
	permissionsExperimentalWrite := []string{
		model.PermissionSysconsoleWriteExperimentalBleve.Id,
		model.PermissionSysconsoleWriteExperimentalFeatures.Id,
		model.PermissionSysconsoleWriteExperimentalFeatureFlags.Id,
	}

	// Give the new subsection READ permissions to any user with READ_EXPERIMENTAL
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleReadExperimental.Id),
		Add: permissionsExperimentalRead,
	})

	// Give the new subsection WRITE permissions to any user with WRITE_EXPERIMENTAL
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteExperimental.Id),
		Add: permissionsExperimentalWrite,
	})

	// Give the ancillary permissions MANAGE_JOBS and PURGE_BLEVE_INDEXES to anyone with WRITE_EXPERIMENTAL_BLEVE
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteExperimentalBleve.Id),
		Add: []string{model.PermissionCreatePostBleveIndexesJob.Id, model.PermissionPurgeBleveIndexes.Id},
	})

	return transformations, nil
}

func (a *App) getAddIntegrationsSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	permissionsIntegrationsRead := []string{
		model.PermissionSysconsoleReadIntegrationsIntegrationManagement.Id,
		model.PermissionSysconsoleReadIntegrationsBotAccounts.Id,
		model.PermissionSysconsoleReadIntegrationsGif.Id,
		model.PermissionSysconsoleReadIntegrationsCors.Id,
	}
	permissionsIntegrationsWrite := []string{
		model.PermissionSysconsoleWriteIntegrationsIntegrationManagement.Id,
		model.PermissionSysconsoleWriteIntegrationsBotAccounts.Id,
		model.PermissionSysconsoleWriteIntegrationsGif.Id,
		model.PermissionSysconsoleWriteIntegrationsCors.Id,
	}

	// Give the new subsection READ permissions to any user with READ_INTEGRATIONS
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleReadIntegrations.Id),
		Add: permissionsIntegrationsRead,
	})

	// Give the new subsection WRITE permissions to any user with WRITE_EXPERIMENTAL
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteIntegrations.Id),
		Add: permissionsIntegrationsWrite,
	})

	return transformations, nil
}

func (a *App) getAddSiteSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_SITE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleReadSite.Id),
		Add: []string{
			model.PermissionSysconsoleReadSiteCustomization.Id,
			model.PermissionSysconsoleReadSiteLocalization.Id,
			model.PermissionSysconsoleReadSiteNotifications.Id,
			model.PermissionSysconsoleReadSiteAnnouncementBanner.Id,
			model.PermissionSysconsoleReadSitePosts.Id,
			model.PermissionSysconsoleReadSiteFileSharingAndDownloads.Id,
			model.PermissionSysconsoleReadSitePublicLinks.Id,
			model.PermissionSysconsoleReadSiteNotices.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_SITE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleWriteSite.Id),
		Add: []string{
			model.PermissionSysconsoleWriteSiteCustomization.Id,
			model.PermissionSysconsoleWriteSiteLocalization.Id,
			model.PermissionSysconsoleWriteSiteNotifications.Id,
			model.PermissionSysconsoleWriteSiteAnnouncementBanner.Id,
			model.PermissionSysconsoleWriteSitePosts.Id,
			model.PermissionSysconsoleWriteSiteFileSharingAndDownloads.Id,
			model.PermissionSysconsoleWriteSitePublicLinks.Id,
			model.PermissionSysconsoleWriteSiteNotices.Id,
		},
	})

	// Give the ancillary permissions EDIT_BRAND to anyone with WRITE_SITE_CUSTOMIZATION
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteSiteCustomization.Id),
		Add: []string{model.PermissionEditBrand.Id},
	})

	return transformations, nil
}

func (a *App) getAddComplianceSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_COMPLIANCE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleReadCompliance.Id),
		Add: []string{
			model.PermissionSysconsoleReadComplianceDataRetentionPolicy.Id,
			model.PermissionSysconsoleReadComplianceComplianceExport.Id,
			model.PermissionSysconsoleReadComplianceComplianceMonitoring.Id,
			model.PermissionSysconsoleReadComplianceCustomTermsOfService.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_COMPLIANCE
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleReadCompliance.Id),
		Add: []string{
			model.PermissionSysconsoleWriteComplianceDataRetentionPolicy.Id,
			model.PermissionSysconsoleWriteComplianceComplianceExport.Id,
			model.PermissionSysconsoleWriteComplianceComplianceMonitoring.Id,
			model.PermissionSysconsoleWriteComplianceCustomTermsOfService.Id,
		},
	})

	// Ancilary permissions
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteComplianceDataRetentionPolicy.Id),
		Add: []string{model.PermissionCreateDataRetentionJob.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleReadComplianceDataRetentionPolicy.Id),
		Add: []string{model.PermissionReadDataRetentionJob.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteComplianceComplianceExport.Id),
		Add: []string{model.PermissionCreateComplianceExportJob.Id, model.PermissionDownloadComplianceExportResult.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleReadComplianceComplianceExport.Id),
		Add: []string{model.PermissionReadComplianceExportJob.Id, model.PermissionDownloadComplianceExportResult.Id},
	})

	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleReadComplianceCustomTermsOfService.Id),
		Add: []string{model.PermissionReadAudits.Id},
	})

	return transformations, nil
}

func (a *App) getAddEnvironmentSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_ENVIRONMENT
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleReadEnvironment.Id),
		Add: []string{
			model.PermissionSysconsoleReadEnvironmentWebServer.Id,
			model.PermissionSysconsoleReadEnvironmentDatabase.Id,
			model.PermissionSysconsoleReadEnvironmentElasticsearch.Id,
			model.PermissionSysconsoleReadEnvironmentFileStorage.Id,
			model.PermissionSysconsoleReadEnvironmentImageProxy.Id,
			model.PermissionSysconsoleReadEnvironmentSmtp.Id,
			model.PermissionSysconsoleReadEnvironmentPushNotificationServer.Id,
			model.PermissionSysconsoleReadEnvironmentHighAvailability.Id,
			model.PermissionSysconsoleReadEnvironmentRateLimiting.Id,
			model.PermissionSysconsoleReadEnvironmentLogging.Id,
			model.PermissionSysconsoleReadEnvironmentSessionLengths.Id,
			model.PermissionSysconsoleReadEnvironmentPerformanceMonitoring.Id,
			model.PermissionSysconsoleReadEnvironmentDeveloper.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_ENVIRONMENT
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleWriteEnvironment.Id),
		Add: []string{
			model.PermissionSysconsoleWriteEnvironmentWebServer.Id,
			model.PermissionSysconsoleWriteEnvironmentDatabase.Id,
			model.PermissionSysconsoleWriteEnvironmentElasticsearch.Id,
			model.PermissionSysconsoleWriteEnvironmentFileStorage.Id,
			model.PermissionSysconsoleWriteEnvironmentImageProxy.Id,
			model.PermissionSysconsoleWriteEnvironmentSmtp.Id,
			model.PermissionSysconsoleWriteEnvironmentPushNotificationServer.Id,
			model.PermissionSysconsoleWriteEnvironmentHighAvailability.Id,
			model.PermissionSysconsoleWriteEnvironmentRateLimiting.Id,
			model.PermissionSysconsoleWriteEnvironmentLogging.Id,
			model.PermissionSysconsoleWriteEnvironmentSessionLengths.Id,
			model.PermissionSysconsoleWriteEnvironmentPerformanceMonitoring.Id,
			model.PermissionSysconsoleWriteEnvironmentDeveloper.Id,
		},
	})

	// Give these ancillary permissions to anyone with READ_ENVIRONMENT_ELASTICSEARCH
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleReadEnvironmentElasticsearch.Id),
		Add: []string{
			model.PermissionReadElasticsearchPostIndexingJob.Id,
			model.PermissionReadElasticsearchPostAggregationJob.Id,
		},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_WEB_SERVER
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleWriteEnvironmentWebServer.Id),
		Add: []string{
			model.PermissionTestSiteUrl.Id,
			model.PermissionReloadConfig.Id,
			model.PermissionInvalidateCaches.Id,
		},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_DATABASE
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteEnvironmentDatabase.Id),
		Add: []string{model.PermissionRecycleDatabaseConnections.Id},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_ELASTICSEARCH
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleWriteEnvironmentElasticsearch.Id),
		Add: []string{
			model.PermissionTestElasticsearch.Id,
			model.PermissionCreateElasticsearchPostIndexingJob.Id,
			model.PermissionCreateElasticsearchPostAggregationJob.Id,
			model.PermissionPurgeElasticsearchIndexes.Id,
		},
	})

	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_FILE_STORAGE
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleWriteEnvironmentFileStorage.Id),
		Add: []string{model.PermissionTestS3.Id},
	})

	return transformations, nil
}

func (a *App) getAddReportingSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{}

	// Give the new subsection READ permissions to any user with READ_REPORTING
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleReadReporting.Id),
		Add: []string{
			model.PermissionSysconsoleReadReportingSiteStatistics.Id,
			model.PermissionSysconsoleReadReportingServerLogs.Id,
		},
	})

	// Give the new subsection WRITE permissions to any user with WRITE_REPORTING
	transformations = append(transformations, permissionTransformation{
		On: permissionExists(model.PermissionSysconsoleWriteReporting.Id),
		Add: []string{
			model.PermissionSysconsoleWriteReportingSiteStatistics.Id,
			model.PermissionSysconsoleWriteReportingServerLogs.Id,
		},
	})

	// Give the ancillary permissions PERMISSION_GET_ANALYTICS to anyone with PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS or PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS
	transformations = append(transformations, permissionTransformation{
		On: permissionOr(
			permissionExists(model.PermissionSysconsoleReadUserManagementUsers.Id),
			permissionExists(model.PermissionSysconsoleReadReportingSiteStatistics.Id),
		),
		Add: []string{model.PermissionGetAnalytics.Id},
	})

	// Give the ancillary permissions PERMISSION_GET_LOGS to anyone with PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS
	transformations = append(transformations, permissionTransformation{
		On:  permissionExists(model.PermissionSysconsoleReadReportingServerLogs.Id),
		Add: []string{model.PermissionGetLogs.Id},
	})

	return transformations, nil
}

func (a *App) getAddAuthenticationSubsectionPermissions() (permissionsMap, error) {
	transformations := []permissionTransformation{
		{ // Give the new subsection READ permissions to any user with READ_AUTHENTICATION
			On: permissionExists(model.PermissionSysconsoleReadAuthentication.Id),
			Add: []string{
				model.PermissionSysconsoleReadAuthenticationSignup.Id,
				model.PermissionSysconsoleReadAuthenticationEmail.Id,
				model.PermissionSysconsoleReadAuthenticationPassword.Id,
				model.PermissionSysconsoleReadAuthenticationMfa.Id,
				model.PermissionSysconsoleReadAuthenticationLdap.Id,
				model.PermissionSysconsoleReadAuthenticationSaml.Id,
				model.PermissionSysconsoleReadAuthenticationOpenid.Id,
				model.PermissionSysconsoleReadAuthenticationGuestAccess.Id,
			},
		},
		{ // Give the new subsection WRITE permissions to any user with WRITE_AUTHENTICATION
			On: permissionExists(model.PermissionSysconsoleWriteAuthentication.Id),
			Add: []string{
				model.PermissionSysconsoleWriteAuthenticationSignup.Id,
				model.PermissionSysconsoleWriteAuthenticationEmail.Id,
				model.PermissionSysconsoleWriteAuthenticationPassword.Id,
				model.PermissionSysconsoleWriteAuthenticationMfa.Id,
				model.PermissionSysconsoleWriteAuthenticationLdap.Id,
				model.PermissionSysconsoleWriteAuthenticationSaml.Id,
				model.PermissionSysconsoleWriteAuthenticationOpenid.Id,
				model.PermissionSysconsoleWriteAuthenticationGuestAccess.Id,
			},
		},
		{ // Give the ancillary permissions for LDAP to anyone with WRITE_AUTHENTICATION_LDAP
			On:  permissionExists(model.PermissionSysconsoleWriteAuthenticationLdap.Id),
			Add: []string{model.PermissionCreateLdapSyncJob.Id, model.PermissionTestLdap.Id, model.PermissionAddLdapPublicCert.Id, model.PermissionAddLdapPrivateCert.Id, model.PermissionRemoveLdapPublicCert.Id, model.PermissionRemoveLdapPrivateCert.Id},
		},
		{ // Give the ancillary permissions PERMISSION_TEST_LDAP to anyone with READ_AUTHENTICATION_LDAP
			On:  permissionExists(model.PermissionSysconsoleReadAuthenticationLdap.Id),
			Add: []string{model.PermissionReadLdapSyncJob.Id},
		},
		{
			On:  permissionExists(model.PermissionSysconsoleWriteAuthenticationEmail.Id),
			Add: []string{model.PermissionInvalidateEmailInvite.Id},
		},
		{ // Give the ancillary permissions for SAML to anyone with WRITE_AUTHENTICATION_SAML
			On:  permissionExists(model.PermissionSysconsoleWriteAuthenticationSaml.Id),
			Add: []string{model.PermissionGetSamlMetadataFromIdp.Id, model.PermissionAddSamlPublicCert.Id, model.PermissionAddSamlPrivateCert.Id, model.PermissionAddSamlIdpCert.Id, model.PermissionRemoveSamlPublicCert.Id, model.PermissionRemoveSamlPrivateCert.Id, model.PermissionRemoveSamlIdpCert.Id, model.PermissionGetSamlCertStatus.Id},
		},
	}

	return transformations, nil
}

// This migration fixes https://github.com/mattermost/mattermost-server/issues/17642 where this particular ancillary permission was forgotten during the initial migrations
func (a *App) getAddTestEmailAncillaryPermission() (permissionsMap, error) {
	// Give these ancillary permissions to anyone with WRITE_ENVIRONMENT_SMTP
	return []permissionTransformation{
		{
			On:  permissionExists(model.PermissionSysconsoleWriteEnvironmentSmtp.Id),
			Add: []string{model.PermissionTestEmail.Id},
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
		{Key: model.MigrationKeyWebhookPermissionsSplit, Migration: a.getWebhooksPermissionsSplitMigration},
		{Key: model.MigrationKeyRemovePermanentDeleteUser, Migration: a.removePermanentDeleteUserMigration},
		{Key: model.MigrationKeyViewMembersNewPermission, Migration: a.getViewMembersPermissionMigration},
		{Key: model.MigrationKeyAddSystemConsolePermissions, Migration: a.getAddSystemConsolePermissionsMigration},
		{Key: model.MigrationKeyAddManageSecureConnectionsPermissions, Migration: a.getAddManageSecureConnectionsPermissionsMigration},
		{Key: model.MigrationKeyAddSystemRolesPermissions, Migration: a.getSystemRolesPermissionsMigration},
		{Key: model.MigrationKeyAddBillingPermissions, Migration: a.getBillingPermissionsMigration},
		{Key: model.MigrationKeyAddDownloadComplianceExportResults, Migration: a.getAddDownloadComplianceExportResult},
		{Key: model.MigrationKeyAddExperimentalSubsectionPermissions, Migration: a.getAddExperimentalSubsectionPermissions},
		{Key: model.MigrationKeyAddAuthenticationSubsectionPermissions, Migration: a.getAddAuthenticationSubsectionPermissions},
		{Key: model.MigrationKeyAddIntegrationsSubsectionPermissions, Migration: a.getAddIntegrationsSubsectionPermissions},
		{Key: model.MigrationKeyAddSiteSubsectionPermissions, Migration: a.getAddSiteSubsectionPermissions},
		{Key: model.MigrationKeyAddComplianceSubsectionPermissions, Migration: a.getAddComplianceSubsectionPermissions},
		{Key: model.MigrationKeyAddEnvironmentSubsectionPermissions, Migration: a.getAddEnvironmentSubsectionPermissions},
		{Key: model.MigrationKeyAddReportingSubsectionPermissions, Migration: a.getAddReportingSubsectionPermissions},
		{Key: model.MigrationKeyAddTestEmailAncillaryPermission, Migration: a.getAddTestEmailAncillaryPermission},
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
