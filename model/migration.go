package model

// keys for saving in System to checks if a specific migration is completed ot not
const (
	AdvancedPermissionsMigrationKey       = "AdvancedPermissionsMigrationComplete"
	MigrationKeyAdvancedPermissionsPhase2 = "migration_advanced_permissions_phase_2"
	PopulateCategoriesForTheFirstTimeKey  = "populated_categories_for_first_time"

	MigrationKeyAddSystemRolesPermissions              = "add_system_roles_permissions"
	MigrationKeyAddBillingPermissions                  = "add_billing_permissions"
	MigrationKeyAddSystemConsolePermissions            = "add_system_console_permissions"
	MigrationKeySidebarCategoriesPhase2                = "migration_sidebar_categories_phase_2"
	MigrationKeyViewMembersNewPermission               = "view_members_new_permission"
	MigrationKeyRemovePermanentDeleteUser              = "remove_permanent_delete_user"
	MigrationKeyWebhookPermissionsSplit                = "webhook_permissions_split"
	MigrationKeyAddManageSecureConnectionsPermissions  = "manage_secure_connections_permissions"
	MigrationKeyAddDownloadComplianceExportResults     = "download_compliance_export_results"
	MigrationKeyAddComplianceSubsectionPermissions     = "compliance_subsection_permissions"
	MigrationKeyAddExperimentalSubsectionPermissions   = "experimental_subsection_permissions"
	MigrationKeyAddAuthenticationSubsectionPermissions = "authentication_subsection_permissions"
	MigrationKeyAddSiteSubsectionPermissions           = "site_subsection_permissions"
	MigrationKeyAddEnvironmentSubsectionPermissions    = "environment_subsection_permissions"
	MigrationKeyAddReportingSubsectionPermissions      = "reporting_subsection_permissions"
	MigrationKeyAddAboutSubsectionPermissions          = "about_subsection_permissions"
	MigrationKeyAddIntegrationsSubsectionPermissions   = "integrations_subsection_permissions"
	MigrationKeyAddTestEmailAncillaryPermission        = "test_email_ancillary_permission"
)
