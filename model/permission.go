package model

// scopes for permissions
const (
	PermissionScopeSystem  = "system_scope"
	PermissionScopeTeam    = "team_scope"
	PermissionScopeChannel = "channel_scope"
)

// Permission type
type Permission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
}

// All permissions from saleor
var (
	PERMISSION_MANAGE_USERS                        *Permission
	PERMISSION_MANAGE_STAFF                        *Permission
	PERMISSION_MANAGE_APPS                         *Permission
	PERMISSION_MANAGE_CHANNELS                     *Permission
	PERMISSION_MANAGE_DISCOUNTS                    *Permission
	PERMISSION_MANAGE_PLUGINS                      *Permission
	PERMISSION_MANAGE_GIFT_CARD                    *Permission
	PERMISSION_MANAGE_MENUS                        *Permission
	PERMISSION_MANAGE_CHECKOUTS                    *Permission
	PERMISSION_MANAGE_ORDERS                       *Permission
	PERMISSION_MANAGE_PAGES                        *Permission
	PERMISSION_MANAGE_PAGE_TYPES_AND_ATTRIBUTES    *Permission
	PERMISSION_MANAGE_PRODUCTS                     *Permission
	PERMISSION_MANAGE_PRODUCT_TYPES_AND_ATTRIBUTES *Permission
	PERMISSION_MANAGE_SHIPPING                     *Permission
	PERMISSION_MANAGE_SETTINGS                     *Permission
	PERMISSION_MANAGE_TRANSLATIONS                 *Permission

	SaleorPermissionEnumList []*Permission
	// SaleorPermissionEnumMap has keys of permission ids,
	// values of permission names
	SaleorPermissionEnumMap map[string]string
)

// init all saleor's permissions
func initializeSaleorPermission() {
	PERMISSION_MANAGE_USERS = &Permission{"manage_users", "account.manage_users.name", "account.manage_users.description", PermissionScopeSystem}
	PERMISSION_MANAGE_STAFF = &Permission{"manage_staff", "account.manage_staff.name", "account.manage_staff.description", PermissionScopeSystem}
	PERMISSION_MANAGE_APPS = &Permission{"manage_apps", "app.manage_apps.name", "app.manage_apps.description", PermissionScopeSystem}
	PERMISSION_MANAGE_CHANNELS = &Permission{"manage_channels", "channel.manage_channels.name", "channel.manage_channels.description", PermissionScopeSystem}
	PERMISSION_MANAGE_DISCOUNTS = &Permission{"manage_discounts", "discount.manage_discounts.name", "discount.manage_discounts.description", PermissionScopeSystem}
	PERMISSION_MANAGE_PLUGINS = &Permission{"manage_plugins", "plugins.manage_plugins.name", "plugins.manage_plugins.description", PermissionScopeSystem}
	PERMISSION_MANAGE_GIFT_CARD = &Permission{"manage_gift_card", "giftcard.manage_gift_card.name", "giftcard.manage_gift_card.description", PermissionScopeSystem}
	PERMISSION_MANAGE_MENUS = &Permission{"manage_menus", "menu.manage_menus.name", "menu.manage_menus.description", PermissionScopeSystem}
	PERMISSION_MANAGE_CHECKOUTS = &Permission{"manage_checkouts", "checkout.manage_checkouts.name", "checkout.manage_checkouts.description", PermissionScopeSystem}
	PERMISSION_MANAGE_ORDERS = &Permission{"manage_orders", "orders.manage_orders.name", "orders.manage_orders.description", PermissionScopeSystem}
	PERMISSION_MANAGE_PAGES = &Permission{"manage_pages", "page.manage_pages.name", "page.manage_pages.description", PermissionScopeSystem}
	PERMISSION_MANAGE_PAGE_TYPES_AND_ATTRIBUTES = &Permission{"manage_page_types_and_attributes", "page.manage_page_types_and_attributes.name", "page.manage_page_types_and_attributes.description", PermissionScopeSystem}
	PERMISSION_MANAGE_PRODUCTS = &Permission{"manage_products", "product.manage_products.name", "product.manage_products.description", PermissionScopeSystem}
	PERMISSION_MANAGE_PRODUCT_TYPES_AND_ATTRIBUTES = &Permission{"manage_product_types_and_attributes", "product.manage_product_types_and_attributes.name", "product.manage_product_types_and_attributes.description", PermissionScopeSystem}
	PERMISSION_MANAGE_SHIPPING = &Permission{"manage_shipping", "shipping.manage_shipping.name", "shipping.manage_shipping.description", PermissionScopeSystem}
	PERMISSION_MANAGE_SETTINGS = &Permission{"manage_settings", "site.manage_settings.name", "site.manage_settings.description", PermissionScopeSystem}
	PERMISSION_MANAGE_TRANSLATIONS = &Permission{"manage_translations", "site.manage_translations.name", "site.manage_translations.description", PermissionScopeSystem}

	SaleorPermissionEnumList = []*Permission{
		PERMISSION_MANAGE_USERS,
		PERMISSION_MANAGE_STAFF,
		PERMISSION_MANAGE_APPS,
		PERMISSION_MANAGE_CHANNELS,
		PERMISSION_MANAGE_DISCOUNTS,
		PERMISSION_MANAGE_PLUGINS,
		PERMISSION_MANAGE_GIFT_CARD,
		PERMISSION_MANAGE_MENUS,
		PERMISSION_MANAGE_CHECKOUTS,
		PERMISSION_MANAGE_ORDERS,
		PERMISSION_MANAGE_PAGES,
		PERMISSION_MANAGE_PAGE_TYPES_AND_ATTRIBUTES,
		PERMISSION_MANAGE_PRODUCTS,
		PERMISSION_MANAGE_PRODUCT_TYPES_AND_ATTRIBUTES,
		PERMISSION_MANAGE_SHIPPING,
		PERMISSION_MANAGE_SETTINGS,
		PERMISSION_MANAGE_TRANSLATIONS,
	}

	SaleorPermissionEnumMap = make(map[string]string)
	for _, perm := range SaleorPermissionEnumList {
		SaleorPermissionEnumMap[perm.Id] = perm.Name
	}
}

// deprecated permissions
var (
	PERMISSION_PERMANENT_DELETE_USER           *Permission
	PERMISSION_MANAGE_WEBHOOKS                 *Permission
	PERMISSION_MANAGE_OTHERS_WEBHOOKS          *Permission
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION *Permission
	PERMISSION_SYSCONSOLE_READ_SITE            *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE           *Permission
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT     *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT    *Permission
	PERMISSION_SYSCONSOLE_READ_REPORTING       *Permission
	PERMISSION_SYSCONSOLE_WRITE_REPORTING      *Permission
	PERMISSION_SYSCONSOLE_READ_ABOUT           *Permission
	PERMISSION_SYSCONSOLE_WRITE_ABOUT          *Permission
	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL    *Permission
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL   *Permission
	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS    *Permission
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS   *Permission
	PERMISSION_SYSCONSOLE_READ_COMPLIANCE      *Permission
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE     *Permission
)

var (
	PERMISSION_INVITE_USER                  *Permission // user permissions
	PERMISSION_USE_SLASH_COMMANDS           *Permission // slash commands
	PERMISSION_MANAGE_SLASH_COMMANDS        *Permission //
	PERMISSION_MANAGE_OTHERS_SLASH_COMMANDS *Permission //

	PERMISSION_ASSIGN_SYSTEM_ADMIN_ROLE *Permission
	PERMISSION_MANAGE_ROLES             *Permission
	PERMISSION_EDIT_OTHER_USERS         *Permission

	PERMISSION_ADD_REACTION            *Permission
	PERMISSION_REMOVE_REACTION         *Permission
	PERMISSION_REMOVE_OTHERS_REACTIONS *Permission

	PERMISSION_UPLOAD_FILE                     *Permission
	PERMISSION_GET_PUBLIC_LINK                 *Permission
	PERMISSION_MANAGE_INCOMING_WEBHOOKS        *Permission
	PERMISSION_MANAGE_OUTGOING_WEBHOOKS        *Permission
	PERMISSION_MANAGE_OTHERS_INCOMING_WEBHOOKS *Permission
	PERMISSION_MANAGE_OTHERS_OUTGOING_WEBHOOKS *Permission
	PERMISSION_MANAGE_OAUTH                    *Permission
	PERMISSION_MANAGE_SYSTEM_WIDE_OAUTH        *Permission
	PERMISSION_CREATE_POST                     *Permission
	PERMISSION_CREATE_POST_PUBLIC              *Permission
	PERMISSION_CREATE_POST_EPHEMERAL           *Permission
	PERMISSION_EDIT_POST                       *Permission
	PERMISSION_EDIT_OTHERS_POSTS               *Permission
	PERMISSION_DELETE_POST                     *Permission
	PERMISSION_DELETE_OTHERS_POSTS             *Permission

	PERMISSION_READ_JOBS                *Permission
	PERMISSION_MANAGE_JOBS              *Permission
	PERMISSION_CREATE_USER_ACCESS_TOKEN *Permission
	PERMISSION_READ_USER_ACCESS_TOKEN   *Permission
	PERMISSION_REVOKE_USER_ACCESS_TOKEN *Permission

	PERMISSION_VIEW_MEMBERS                              *Permission
	PERMISSION_EDIT_BRAND                                *Permission
	PERMISSION_MANAGE_REMOTE_CLUSTERS                    *Permission
	PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT         *Permission
	PERMISSION_CREATE_DATA_RETENTION_JOB                 *Permission
	PERMISSION_READ_DATA_RETENTION_JOB                   *Permission
	PERMISSION_CREATE_COMPLIANCE_EXPORT_JOB              *Permission
	PERMISSION_READ_COMPLIANCE_EXPORT_JOB                *Permission
	PERMISSION_READ_AUDITS                               *Permission
	PERMISSION_TEST_ELASTICSEARCH                        *Permission
	PERMISSION_TEST_SITE_URL                             *Permission
	PERMISSION_TEST_S3                                   *Permission
	PERMISSION_RELOAD_CONFIG                             *Permission
	PERMISSION_INVALIDATE_CACHES                         *Permission
	PERMISSION_RECYCLE_DATABASE_CONNECTIONS              *Permission
	PERMISSION_PURGE_ELASTICSEARCH_INDEXES               *Permission
	PERMISSION_TEST_EMAIL                                *Permission
	PERMISSION_CREATE_ELASTICSEARCH_POST_INDEXING_JOB    *Permission
	PERMISSION_CREATE_ELASTICSEARCH_POST_AGGREGATION_JOB *Permission
	PERMISSION_READ_ELASTICSEARCH_POST_INDEXING_JOB      *Permission
	PERMISSION_READ_ELASTICSEARCH_POST_AGGREGATION_JOB   *Permission
	PERMISSION_PURGE_BLEVE_INDEXES                       *Permission
	PERMISSION_CREATE_POST_BLEVE_INDEXES_JOB             *Permission
	PERMISSION_CREATE_LDAP_SYNC_JOB                      *Permission
	PERMISSION_READ_LDAP_SYNC_JOB                        *Permission
	PERMISSION_TEST_LDAP                                 *Permission
	PERMISSION_INVALIDATE_EMAIL_INVITE                   *Permission
	PERMISSION_GET_SAML_METADATA_FROM_IDP                *Permission
	PERMISSION_ADD_SAML_PUBLIC_CERT                      *Permission
	PERMISSION_ADD_SAML_PRIVATE_CERT                     *Permission
	PERMISSION_ADD_SAML_IDP_CERT                         *Permission
	PERMISSION_REMOVE_SAML_PUBLIC_CERT                   *Permission
	PERMISSION_REMOVE_SAML_PRIVATE_CERT                  *Permission
	PERMISSION_REMOVE_SAML_IDP_CERT                      *Permission
	PERMISSION_GET_SAML_CERT_STATUS                      *Permission
	PERMISSION_ADD_LDAP_PUBLIC_CERT                      *Permission
	PERMISSION_ADD_LDAP_PRIVATE_CERT                     *Permission
	PERMISSION_REMOVE_LDAP_PUBLIC_CERT                   *Permission
	PERMISSION_REMOVE_LDAP_PRIVATE_CERT                  *Permission
	PERMISSION_GET_LOGS                                  *Permission
	PERMISSION_GET_ANALYTICS                             *Permission

	PERMISSION_SYSCONSOLE_READ_BILLING  *Permission
	PERMISSION_SYSCONSOLE_WRITE_BILLING *Permission

	PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_REPORTING_SITE_STATISTICS *Permission

	PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_REPORTING_SERVER_LOGS *Permission

	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_USERS *Permission

	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_GROUPS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_GROUPS *Permission

	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_PERMISSIONS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_PERMISSIONS *Permission

	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_SYSTEM_ROLES  *Permission
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_SYSTEM_ROLES *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_WEB_SERVER  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_WEB_SERVER *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DATABASE  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DATABASE *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_ELASTICSEARCH *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_FILE_STORAGE  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_FILE_STORAGE *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_IMAGE_PROXY  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_IMAGE_PROXY *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SMTP  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SMTP *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PUSH_NOTIFICATION_SERVER  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PUSH_NOTIFICATION_SERVER *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_HIGH_AVAILABILITY  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_HIGH_AVAILABILITY *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_RATE_LIMITING  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_RATE_LIMITING *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_LOGGING  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_LOGGING *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SESSION_LENGTHS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SESSION_LENGTHS *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PERFORMANCE_MONITORING  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PERFORMANCE_MONITORING *Permission

	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DEVELOPER  *Permission
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DEVELOPER *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_CUSTOMIZATION  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_CUSTOMIZATION *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_LOCALIZATION  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_LOCALIZATION *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_NOTIFICATIONS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_NOTIFICATIONS *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_ANNOUNCEMENT_BANNER  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_ANNOUNCEMENT_BANNER *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_POSTS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_POSTS *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_FILE_SHARING_AND_DOWNLOADS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_FILE_SHARING_AND_DOWNLOADS *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_PUBLIC_LINKS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_PUBLIC_LINKS *Permission

	PERMISSION_SYSCONSOLE_READ_SITE_NOTICES  *Permission
	PERMISSION_SYSCONSOLE_WRITE_SITE_NOTICES *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SIGNUP  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SIGNUP *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_EMAIL  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_EMAIL *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_PASSWORD  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_PASSWORD *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_MFA  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_MFA *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_LDAP *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SAML  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SAML *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_OPENID  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_OPENID *Permission

	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_GUEST_ACCESS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_GUEST_ACCESS *Permission

	PERMISSION_SYSCONSOLE_READ_PLUGINS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_PLUGINS *Permission

	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_INTEGRATION_MANAGEMENT  *Permission
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_INTEGRATION_MANAGEMENT *Permission

	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_BOT_ACCOUNTS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_BOT_ACCOUNTS *Permission

	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_GIF  *Permission
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_GIF *Permission

	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_CORS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_CORS *Permission

	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_DATA_RETENTION_POLICY  *Permission
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_DATA_RETENTION_POLICY *Permission

	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_EXPORT  *Permission
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_EXPORT *Permission

	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_MONITORING  *Permission
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_MONITORING *Permission

	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE  *Permission
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE *Permission

	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURES  *Permission
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURES *Permission

	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURE_FLAGS  *Permission
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURE_FLAGS *Permission

	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_BLEVE  *Permission
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_BLEVE *Permission

	// General permission that encompasses all system admin functions
	// in the future this could be broken up to allow access to some
	// admin functions but not others
	PERMISSION_MANAGE_SYSTEM *Permission
)

// AllPermissions contains all system's permissions
var AllPermissions []*Permission

// DeprecatedPermissions contains permissions that are deprecated
var DeprecatedPermissions []*Permission

var SysconsoleReadPermissions []*Permission
var SysconsoleWritePermissions []*Permission

func initializePermissions() {
	PERMISSION_INVITE_USER = &Permission{
		"invite_user",
		"authentication.permissions.team_invite_user.name",
		"authentication.permissions.team_invite_user.description",
		PermissionScopeTeam,
	}
	PERMISSION_USE_SLASH_COMMANDS = &Permission{
		"use_slash_commands",
		"authentication.permissions.team_use_slash_commands.name",
		"authentication.permissions.team_use_slash_commands.description",
		PermissionScopeChannel,
	}
	PERMISSION_MANAGE_SLASH_COMMANDS = &Permission{
		"manage_slash_commands",
		"authentication.permissions.manage_slash_commands.name",
		"authentication.permissions.manage_slash_commands.description",
		PermissionScopeTeam,
	}
	PERMISSION_MANAGE_OTHERS_SLASH_COMMANDS = &Permission{
		"manage_others_slash_commands",
		"authentication.permissions.manage_others_slash_commands.name",
		"authentication.permissions.manage_others_slash_commands.description",
		PermissionScopeTeam,
	}
	PERMISSION_ASSIGN_SYSTEM_ADMIN_ROLE = &Permission{
		"assign_system_admin_role",
		"authentication.permissions.assign_system_admin_role.name",
		"authentication.permissions.assign_system_admin_role.description",
		PermissionScopeSystem,
	}
	PERMISSION_MANAGE_ROLES = &Permission{
		"manage_roles",
		"authentication.permissions.manage_roles.name",
		"authentication.permissions.manage_roles.description",
		PermissionScopeSystem,
	}
	PERMISSION_MANAGE_SYSTEM = &Permission{
		"manage_system",
		"authentication.permissions.manage_system.name",
		"authentication.permissions.manage_system.description",
		PermissionScopeSystem,
	}
	PERMISSION_EDIT_OTHER_USERS = &Permission{
		"edit_other_users",
		"authentication.permissions.edit_other_users.name",
		"authentication.permissions.edit_other_users.description",
		PermissionScopeSystem,
	}
	PERMISSION_ADD_REACTION = &Permission{
		"add_reaction",
		"authentication.permissions.add_reaction.name",
		"authentication.permissions.add_reaction.description",
		PermissionScopeChannel,
	}
	PERMISSION_REMOVE_REACTION = &Permission{
		"remove_reaction",
		"authentication.permissions.remove_reaction.name",
		"authentication.permissions.remove_reaction.description",
		PermissionScopeChannel,
	}
	PERMISSION_REMOVE_OTHERS_REACTIONS = &Permission{
		"remove_others_reactions",
		"authentication.permissions.remove_others_reactions.name",
		"authentication.permissions.remove_others_reactions.description",
		PermissionScopeChannel,
	}

	PERMISSION_UPLOAD_FILE = &Permission{
		"upload_file",
		"authentication.permissions.upload_file.name",
		"authentication.permissions.upload_file.description",
		PermissionScopeChannel,
	}
	PERMISSION_GET_PUBLIC_LINK = &Permission{
		"get_public_link",
		"authentication.permissions.get_public_link.name",
		"authentication.permissions.get_public_link.description",
		PermissionScopeSystem,
	}

	PERMISSION_MANAGE_INCOMING_WEBHOOKS = &Permission{
		"manage_incoming_webhooks",
		"authentication.permissions.manage_incoming_webhooks.name",
		"authentication.permissions.manage_incoming_webhooks.description",
		PermissionScopeTeam,
	}
	PERMISSION_MANAGE_OUTGOING_WEBHOOKS = &Permission{
		"manage_outgoing_webhooks",
		"authentication.permissions.manage_outgoing_webhooks.name",
		"authentication.permissions.manage_outgoing_webhooks.description",
		PermissionScopeTeam,
	}
	PERMISSION_MANAGE_OTHERS_INCOMING_WEBHOOKS = &Permission{
		"manage_others_incoming_webhooks",
		"authentication.permissions.manage_others_incoming_webhooks.name",
		"authentication.permissions.manage_others_incoming_webhooks.description",
		PermissionScopeTeam,
	}
	PERMISSION_MANAGE_OTHERS_OUTGOING_WEBHOOKS = &Permission{
		"manage_others_outgoing_webhooks",
		"authentication.permissions.manage_others_outgoing_webhooks.name",
		"authentication.permissions.manage_others_outgoing_webhooks.description",
		PermissionScopeTeam,
	}
	PERMISSION_MANAGE_OAUTH = &Permission{
		"manage_oauth",
		"authentication.permissions.manage_oauth.name",
		"authentication.permissions.manage_oauth.description",
		PermissionScopeSystem,
	}
	PERMISSION_MANAGE_SYSTEM_WIDE_OAUTH = &Permission{
		"manage_system_wide_oauth",
		"authentication.permissions.manage_system_wide_oauth.name",
		"authentication.permissions.manage_system_wide_oauth.description",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_POST = &Permission{
		"create_post",
		"authentication.permissions.create_post.name",
		"authentication.permissions.create_post.description",
		PermissionScopeChannel,
	}
	PERMISSION_CREATE_POST_PUBLIC = &Permission{
		"create_post_public",
		"authentication.permissions.create_post_public.name",
		"authentication.permissions.create_post_public.description",
		PermissionScopeChannel,
	}
	PERMISSION_CREATE_POST_EPHEMERAL = &Permission{
		"create_post_ephemeral",
		"authentication.permissions.create_post_ephemeral.name",
		"authentication.permissions.create_post_ephemeral.description",
		PermissionScopeChannel,
	}
	PERMISSION_EDIT_POST = &Permission{
		"edit_post",
		"authentication.permissions.edit_post.name",
		"authentication.permissions.edit_post.description",
		PermissionScopeChannel,
	}
	PERMISSION_EDIT_OTHERS_POSTS = &Permission{
		"edit_others_posts",
		"authentication.permissions.edit_others_posts.name",
		"authentication.permissions.edit_others_posts.description",
		PermissionScopeChannel,
	}
	PERMISSION_DELETE_POST = &Permission{
		"delete_post",
		"authentication.permissions.delete_post.name",
		"authentication.permissions.delete_post.description",
		PermissionScopeChannel,
	}
	PERMISSION_DELETE_OTHERS_POSTS = &Permission{
		"delete_others_posts",
		"authentication.permissions.delete_others_posts.name",
		"authentication.permissions.delete_others_posts.description",
		PermissionScopeChannel,
	}
	PERMISSION_MANAGE_REMOTE_CLUSTERS = &Permission{
		"manage_remote_clusters",
		"authentication.permissions.manage_remote_clusters.name",
		"authentication.permissions.manage_remote_clusters.description",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_DATA_RETENTION_JOB = &Permission{"create_data_retention_job", "", "", PermissionScopeSystem}
	PERMISSION_READ_DATA_RETENTION_JOB = &Permission{
		"read_data_retention_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_COMPLIANCE_EXPORT_JOB = &Permission{
		"create_compliance_export_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_READ_COMPLIANCE_EXPORT_JOB = &Permission{
		"read_compliance_export_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_READ_AUDITS = &Permission{
		"read_audits",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_PURGE_BLEVE_INDEXES = &Permission{
		"purge_bleve_indexes",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_POST_BLEVE_INDEXES_JOB = &Permission{
		"create_post_bleve_indexes_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_LDAP_SYNC_JOB = &Permission{
		"create_ldap_sync_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_READ_LDAP_SYNC_JOB = &Permission{
		"read_ldap_sync_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_TEST_LDAP = &Permission{
		"test_ldap",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_INVALIDATE_EMAIL_INVITE = &Permission{
		"invalidate_email_invite",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_GET_SAML_METADATA_FROM_IDP = &Permission{
		"get_saml_metadata_from_idp",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_ADD_SAML_PUBLIC_CERT = &Permission{
		"add_saml_public_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_ADD_SAML_PRIVATE_CERT = &Permission{
		"add_saml_private_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_ADD_SAML_IDP_CERT = &Permission{
		"add_saml_idp_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_REMOVE_SAML_PUBLIC_CERT = &Permission{
		"remove_saml_public_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_REMOVE_SAML_PRIVATE_CERT = &Permission{
		"remove_saml_private_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_REMOVE_SAML_IDP_CERT = &Permission{
		"remove_saml_idp_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_GET_SAML_CERT_STATUS = &Permission{
		"get_saml_cert_status",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_ADD_LDAP_PUBLIC_CERT = &Permission{
		"add_ldap_public_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_ADD_LDAP_PRIVATE_CERT = &Permission{
		"add_ldap_private_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_REMOVE_LDAP_PUBLIC_CERT = &Permission{
		"remove_ldap_public_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_REMOVE_LDAP_PRIVATE_CERT = &Permission{
		"remove_ldap_private_cert",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_GET_LOGS = &Permission{"get_logs", "", "", PermissionScopeSystem}
	PERMISSION_GET_ANALYTICS = &Permission{
		"get_analytics",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT = &Permission{
		"download_compliance_export_result",
		"authentication.permissions.download_compliance_export_result.name",
		"authentication.permissions.download_compliance_export_result.description",
		PermissionScopeSystem,
	}
	PERMISSION_TEST_SITE_URL = &Permission{
		"test_site_url",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_TEST_ELASTICSEARCH = &Permission{
		"test_elasticsearch",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_TEST_S3 = &Permission{
		"test_s3",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_RELOAD_CONFIG = &Permission{
		"reload_config",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_INVALIDATE_CACHES = &Permission{
		"invalidate_caches",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_RECYCLE_DATABASE_CONNECTIONS = &Permission{
		"recycle_database_connections",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_PURGE_ELASTICSEARCH_INDEXES = &Permission{
		"purge_elasticsearch_indexes",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_TEST_EMAIL = &Permission{
		"test_email",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_ELASTICSEARCH_POST_INDEXING_JOB = &Permission{
		"create_elasticsearch_post_indexing_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_ELASTICSEARCH_POST_AGGREGATION_JOB = &Permission{
		"create_elasticsearch_post_aggregation_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_READ_ELASTICSEARCH_POST_INDEXING_JOB = &Permission{
		"read_elasticsearch_post_indexing_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_READ_ELASTICSEARCH_POST_AGGREGATION_JOB = &Permission{
		"read_elasticsearch_post_aggregation_job",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_CREATE_USER_ACCESS_TOKEN = &Permission{
		"create_user_access_token",
		"authentication.permissions.create_user_access_token.name",
		"authentication.permissions.create_user_access_token.description",
		PermissionScopeSystem,
	}
	PERMISSION_READ_USER_ACCESS_TOKEN = &Permission{
		"read_user_access_token",
		"authentication.permissions.read_user_access_token.name",
		"authentication.permissions.read_user_access_token.description",
		PermissionScopeSystem,
	}
	PERMISSION_REVOKE_USER_ACCESS_TOKEN = &Permission{
		"revoke_user_access_token",
		"authentication.permissions.revoke_user_access_token.name",
		"authentication.permissions.revoke_user_access_token.description",
		PermissionScopeSystem,
	}
	PERMISSION_READ_JOBS = &Permission{
		"read_jobs",
		"authentication.permisssions.read_jobs.name",
		"authentication.permisssions.read_jobs.description",
		PermissionScopeSystem,
	}
	PERMISSION_MANAGE_JOBS = &Permission{
		"manage_jobs",
		"authentication.permisssions.manage_jobs.name",
		"authentication.permisssions.manage_jobs.description",
		PermissionScopeSystem,
	}
	PERMISSION_VIEW_MEMBERS = &Permission{
		"view_members",
		"authentication.permisssions.view_members.name",
		"authentication.permisssions.view_members.description",
		PermissionScopeTeam,
	}
	PERMISSION_EDIT_BRAND = &Permission{
		"edit_brand",
		"authentication.permissions.edit_brand.name",
		"authentication.permissions.edit_brand.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_BILLING = &Permission{
		"sysconsole_read_billing",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_BILLING = &Permission{
		"sysconsole_write_billing",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS = &Permission{
		"sysconsole_read_reporting_site_statistics",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_REPORTING_SITE_STATISTICS = &Permission{
		"sysconsole_write_reporting_site_statistics",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS = &Permission{
		"sysconsole_read_reporting_server_logs",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_REPORTING_SERVER_LOGS = &Permission{
		"sysconsole_write_reporting_server_logs",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS = &Permission{
		"sysconsole_read_user_management_users",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_USERS = &Permission{
		"sysconsole_write_user_management_users",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_GROUPS = &Permission{
		"sysconsole_read_user_management_groups",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_GROUPS = &Permission{
		"sysconsole_write_user_management_groups",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_PERMISSIONS = &Permission{
		"sysconsole_read_user_management_permissions",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_PERMISSIONS = &Permission{
		"sysconsole_write_user_management_permissions",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_SYSTEM_ROLES = &Permission{
		"sysconsole_read_user_management_system_roles",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_SYSTEM_ROLES = &Permission{
		"sysconsole_write_user_management_system_roles",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_WEB_SERVER = &Permission{
		"sysconsole_read_environment_web_server",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_WEB_SERVER = &Permission{
		"sysconsole_write_environment_web_server",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DATABASE = &Permission{
		"sysconsole_read_environment_database",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DATABASE = &Permission{
		"sysconsole_write_environment_database",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH = &Permission{
		"sysconsole_read_environment_elasticsearch",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_ELASTICSEARCH = &Permission{
		"sysconsole_write_environment_elasticsearch",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_FILE_STORAGE = &Permission{
		"sysconsole_read_environment_file_storage",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_FILE_STORAGE = &Permission{
		"sysconsole_write_environment_file_storage",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_IMAGE_PROXY = &Permission{
		"sysconsole_read_environment_image_proxy",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_IMAGE_PROXY = &Permission{
		"sysconsole_write_environment_image_proxy",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SMTP = &Permission{
		"sysconsole_read_environment_smtp",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SMTP = &Permission{
		"sysconsole_write_environment_smtp",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PUSH_NOTIFICATION_SERVER = &Permission{
		"sysconsole_read_environment_push_notification_server",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PUSH_NOTIFICATION_SERVER = &Permission{
		"sysconsole_write_environment_push_notification_server",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_HIGH_AVAILABILITY = &Permission{
		"sysconsole_read_environment_high_availability",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_HIGH_AVAILABILITY = &Permission{
		"sysconsole_write_environment_high_availability",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_RATE_LIMITING = &Permission{
		"sysconsole_read_environment_rate_limiting",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_RATE_LIMITING = &Permission{
		"sysconsole_write_environment_rate_limiting",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_LOGGING = &Permission{
		"sysconsole_read_environment_logging",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_LOGGING = &Permission{
		"sysconsole_write_environment_logging",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SESSION_LENGTHS = &Permission{
		"sysconsole_read_environment_session_lengths",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SESSION_LENGTHS = &Permission{
		"sysconsole_write_environment_session_lengths",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PERFORMANCE_MONITORING = &Permission{
		"sysconsole_read_environment_performance_monitoring",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PERFORMANCE_MONITORING = &Permission{
		"sysconsole_write_environment_performance_monitoring",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DEVELOPER = &Permission{
		"sysconsole_read_environment_developer",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DEVELOPER = &Permission{
		"sysconsole_write_environment_developer",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_CUSTOMIZATION = &Permission{
		"sysconsole_read_site_customization",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_CUSTOMIZATION = &Permission{
		"sysconsole_write_site_customization",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_LOCALIZATION = &Permission{
		"sysconsole_read_site_localization",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_LOCALIZATION = &Permission{
		"sysconsole_write_site_localization",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_NOTIFICATIONS = &Permission{
		"sysconsole_read_site_notifications",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_NOTIFICATIONS = &Permission{
		"sysconsole_write_site_notifications",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_ANNOUNCEMENT_BANNER = &Permission{
		"sysconsole_read_site_announcement_banner",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_ANNOUNCEMENT_BANNER = &Permission{
		"sysconsole_write_site_announcement_banner",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_POSTS = &Permission{
		"sysconsole_read_site_posts",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_POSTS = &Permission{
		"sysconsole_write_site_posts",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_FILE_SHARING_AND_DOWNLOADS = &Permission{
		"sysconsole_read_site_file_sharing_and_downloads",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_FILE_SHARING_AND_DOWNLOADS = &Permission{
		"sysconsole_write_site_file_sharing_and_downloads",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_PUBLIC_LINKS = &Permission{
		"sysconsole_read_site_public_links",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_PUBLIC_LINKS = &Permission{
		"sysconsole_write_site_public_links",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE_NOTICES = &Permission{
		"sysconsole_read_site_notices",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE_NOTICES = &Permission{
		"sysconsole_write_site_notices",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SIGNUP = &Permission{
		"sysconsole_read_authentication_signup",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SIGNUP = &Permission{
		"sysconsole_write_authentication_signup",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_EMAIL = &Permission{
		"sysconsole_read_authentication_email",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_EMAIL = &Permission{
		"sysconsole_write_authentication_email",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_PASSWORD = &Permission{
		"sysconsole_read_authentication_password",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_PASSWORD = &Permission{
		"sysconsole_write_authentication_password",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_MFA = &Permission{
		"sysconsole_read_authentication_mfa",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_MFA = &Permission{
		"sysconsole_write_authentication_mfa",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP = &Permission{
		"sysconsole_read_authentication_ldap",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_LDAP = &Permission{
		"sysconsole_write_authentication_ldap",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SAML = &Permission{
		"sysconsole_read_authentication_saml",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SAML = &Permission{
		"sysconsole_write_authentication_saml",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_OPENID = &Permission{
		"sysconsole_read_authentication_openid",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_OPENID = &Permission{
		"sysconsole_write_authentication_openid",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_GUEST_ACCESS = &Permission{
		"sysconsole_read_authentication_guest_access",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_GUEST_ACCESS = &Permission{
		"sysconsole_write_authentication_guest_access",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_PLUGINS = &Permission{
		"sysconsole_read_plugins",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_PLUGINS = &Permission{
		"sysconsole_write_plugins",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_INTEGRATION_MANAGEMENT = &Permission{
		"sysconsole_read_integrations_integration_management",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_INTEGRATION_MANAGEMENT = &Permission{
		"sysconsole_write_integrations_integration_management",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_BOT_ACCOUNTS = &Permission{
		"sysconsole_read_integrations_bot_accounts",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_BOT_ACCOUNTS = &Permission{
		"sysconsole_write_integrations_bot_accounts",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_GIF = &Permission{
		"sysconsole_read_integrations_gif",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_GIF = &Permission{
		"sysconsole_write_integrations_gif",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_CORS = &Permission{
		"sysconsole_read_integrations_cors",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_CORS = &Permission{
		"sysconsole_write_integrations_cors",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_DATA_RETENTION_POLICY = &Permission{
		"sysconsole_read_compliance_data_retention_policy",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_DATA_RETENTION_POLICY = &Permission{
		"sysconsole_write_compliance_data_retention_policy",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_EXPORT = &Permission{
		"sysconsole_read_compliance_compliance_export",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_EXPORT = &Permission{
		"sysconsole_write_compliance_compliance_export",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_MONITORING = &Permission{
		"sysconsole_read_compliance_compliance_monitoring",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_MONITORING = &Permission{
		"sysconsole_write_compliance_compliance_monitoring",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE = &Permission{
		"sysconsole_read_compliance_custom_terms_of_service",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE = &Permission{
		"sysconsole_write_compliance_custom_terms_of_service",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURES = &Permission{
		"sysconsole_read_experimental_features",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURES = &Permission{
		"sysconsole_write_experimental_features",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURE_FLAGS = &Permission{
		"sysconsole_read_experimental_feature_flags",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURE_FLAGS = &Permission{
		"sysconsole_write_experimental_feature_flags",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_BLEVE = &Permission{
		"sysconsole_read_experimental_bleve",
		"",
		"",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_BLEVE = &Permission{
		"sysconsole_write_experimental_bleve",
		"",
		"",
		PermissionScopeSystem,
	}

	// deprecated-------------------------------
	PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL = &Permission{
		"sysconsole_read_experimental",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL = &Permission{
		"sysconsole_write_experimental",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_COMPLIANCE = &Permission{
		"sysconsole_read_compliance",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE = &Permission{
		"sysconsole_write_compliance",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_INTEGRATIONS = &Permission{
		"sysconsole_read_integrations",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS = &Permission{
		"sysconsole_write_integrations",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_AUTHENTICATION = &Permission{
		"sysconsole_read_authentication",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION = &Permission{
		"sysconsole_write_authentication",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_SITE = &Permission{
		"sysconsole_read_site",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_SITE = &Permission{
		"sysconsole_write_site",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ENVIRONMENT = &Permission{
		"sysconsole_read_environment",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT = &Permission{
		"sysconsole_write_environment",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_REPORTING = &Permission{
		"sysconsole_read_reporting",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_REPORTING = &Permission{
		"sysconsole_write_reporting",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_READ_ABOUT = &Permission{
		"sysconsole_read_about",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_SYSCONSOLE_WRITE_ABOUT = &Permission{
		"sysconsole_write_about",
		"authentication.permissions.use_group_mentions.name",
		"authentication.permissions.use_group_mentions.description",
		PermissionScopeSystem,
	}
	PERMISSION_MANAGE_WEBHOOKS = &Permission{
		"manage_webhooks",
		"authentication.permissions.manage_webhooks.name",
		"authentication.permissions.manage_webhooks.description",
		PermissionScopeTeam,
	}
	PERMISSION_PERMANENT_DELETE_USER = &Permission{
		"permanent_delete_user",
		"authentication.permissions.permanent_delete_user.name",
		"authentication.permissions.permanent_delete_user.description",
		PermissionScopeSystem,
	}
	PERMISSION_MANAGE_OTHERS_WEBHOOKS = &Permission{
		"manage_others_webhooks",
		"authentication.permissions.manage_others_webhooks.name",
		"authentication.permissions.manage_others_webhooks.description",
		PermissionScopeTeam,
	}
	DeprecatedPermissions = []*Permission{
		PERMISSION_PERMANENT_DELETE_USER,
		PERMISSION_MANAGE_WEBHOOKS,
		PERMISSION_MANAGE_OTHERS_WEBHOOKS,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION,
		PERMISSION_SYSCONSOLE_READ_SITE,
		PERMISSION_SYSCONSOLE_WRITE_SITE,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT,
		PERMISSION_SYSCONSOLE_READ_REPORTING,
		PERMISSION_SYSCONSOLE_WRITE_REPORTING,
		PERMISSION_SYSCONSOLE_READ_ABOUT,
		PERMISSION_SYSCONSOLE_WRITE_ABOUT,
		PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL,
		PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE,
		PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE,
	}

	SysconsoleReadPermissions = []*Permission{
		PERMISSION_SYSCONSOLE_READ_BILLING,
		PERMISSION_SYSCONSOLE_READ_REPORTING_SITE_STATISTICS,
		PERMISSION_SYSCONSOLE_READ_REPORTING_SERVER_LOGS,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_GROUPS,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_PERMISSIONS,
		PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_SYSTEM_ROLES,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_WEB_SERVER,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DATABASE,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_ELASTICSEARCH,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_FILE_STORAGE,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_IMAGE_PROXY,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SMTP,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PUSH_NOTIFICATION_SERVER,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_HIGH_AVAILABILITY,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_RATE_LIMITING,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_LOGGING,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_SESSION_LENGTHS,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_PERFORMANCE_MONITORING,
		PERMISSION_SYSCONSOLE_READ_ENVIRONMENT_DEVELOPER,
		PERMISSION_SYSCONSOLE_READ_SITE_CUSTOMIZATION,
		PERMISSION_SYSCONSOLE_READ_SITE_LOCALIZATION,
		PERMISSION_SYSCONSOLE_READ_SITE_NOTIFICATIONS,
		PERMISSION_SYSCONSOLE_READ_SITE_ANNOUNCEMENT_BANNER,
		PERMISSION_SYSCONSOLE_READ_SITE_POSTS,
		PERMISSION_SYSCONSOLE_READ_SITE_FILE_SHARING_AND_DOWNLOADS,
		PERMISSION_SYSCONSOLE_READ_SITE_PUBLIC_LINKS,
		PERMISSION_SYSCONSOLE_READ_SITE_NOTICES,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SIGNUP,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_EMAIL,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_PASSWORD,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_MFA,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_LDAP,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_SAML,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_OPENID,
		PERMISSION_SYSCONSOLE_READ_AUTHENTICATION_GUEST_ACCESS,
		PERMISSION_SYSCONSOLE_READ_PLUGINS,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_INTEGRATION_MANAGEMENT,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_BOT_ACCOUNTS,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_GIF,
		PERMISSION_SYSCONSOLE_READ_INTEGRATIONS_CORS,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_DATA_RETENTION_POLICY,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_EXPORT,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_COMPLIANCE_MONITORING,
		PERMISSION_SYSCONSOLE_READ_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE,
		PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURES,
		PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_FEATURE_FLAGS,
		PERMISSION_SYSCONSOLE_READ_EXPERIMENTAL_BLEVE,
	}

	SysconsoleWritePermissions = []*Permission{
		PERMISSION_SYSCONSOLE_WRITE_BILLING,
		PERMISSION_SYSCONSOLE_WRITE_REPORTING_SITE_STATISTICS,
		PERMISSION_SYSCONSOLE_WRITE_REPORTING_SERVER_LOGS,
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_USERS,
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_GROUPS,
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_PERMISSIONS,
		PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_SYSTEM_ROLES,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_WEB_SERVER,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DATABASE,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_ELASTICSEARCH,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_FILE_STORAGE,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_IMAGE_PROXY,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SMTP,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PUSH_NOTIFICATION_SERVER,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_HIGH_AVAILABILITY,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_RATE_LIMITING,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_LOGGING,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_SESSION_LENGTHS,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_PERFORMANCE_MONITORING,
		PERMISSION_SYSCONSOLE_WRITE_ENVIRONMENT_DEVELOPER,
		PERMISSION_SYSCONSOLE_WRITE_SITE_CUSTOMIZATION,
		PERMISSION_SYSCONSOLE_WRITE_SITE_LOCALIZATION,
		PERMISSION_SYSCONSOLE_WRITE_SITE_NOTIFICATIONS,
		PERMISSION_SYSCONSOLE_WRITE_SITE_ANNOUNCEMENT_BANNER,
		PERMISSION_SYSCONSOLE_WRITE_SITE_POSTS,
		PERMISSION_SYSCONSOLE_WRITE_SITE_FILE_SHARING_AND_DOWNLOADS,
		PERMISSION_SYSCONSOLE_WRITE_SITE_PUBLIC_LINKS,
		PERMISSION_SYSCONSOLE_WRITE_SITE_NOTICES,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SIGNUP,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_EMAIL,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_PASSWORD,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_MFA,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_LDAP,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_SAML,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_OPENID,
		PERMISSION_SYSCONSOLE_WRITE_AUTHENTICATION_GUEST_ACCESS,
		PERMISSION_SYSCONSOLE_WRITE_PLUGINS,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_INTEGRATION_MANAGEMENT,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_BOT_ACCOUNTS,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_GIF,
		PERMISSION_SYSCONSOLE_WRITE_INTEGRATIONS_CORS,
		PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_DATA_RETENTION_POLICY,
		PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_EXPORT,
		PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_COMPLIANCE_MONITORING,
		PERMISSION_SYSCONSOLE_WRITE_COMPLIANCE_CUSTOM_TERMS_OF_SERVICE,
		PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURES,
		PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_FEATURE_FLAGS,
		PERMISSION_SYSCONSOLE_WRITE_EXPERIMENTAL_BLEVE,
	}

	SystemScopedPermissionsMinusSysconsole := []*Permission{
		PERMISSION_ASSIGN_SYSTEM_ADMIN_ROLE,
		PERMISSION_MANAGE_ROLES,
		PERMISSION_MANAGE_SYSTEM,
		PERMISSION_EDIT_OTHER_USERS,
		PERMISSION_GET_PUBLIC_LINK,
		PERMISSION_MANAGE_OAUTH,
		PERMISSION_MANAGE_SYSTEM_WIDE_OAUTH,
		PERMISSION_CREATE_USER_ACCESS_TOKEN,
		PERMISSION_READ_USER_ACCESS_TOKEN,
		PERMISSION_REVOKE_USER_ACCESS_TOKEN,
		PERMISSION_READ_JOBS,
		PERMISSION_MANAGE_JOBS,
		PERMISSION_EDIT_BRAND,
		PERMISSION_MANAGE_REMOTE_CLUSTERS,
		PERMISSION_DOWNLOAD_COMPLIANCE_EXPORT_RESULT,
		PERMISSION_CREATE_DATA_RETENTION_JOB,
		PERMISSION_READ_DATA_RETENTION_JOB,
		PERMISSION_CREATE_COMPLIANCE_EXPORT_JOB,
		PERMISSION_READ_COMPLIANCE_EXPORT_JOB,
		PERMISSION_READ_AUDITS,
		PERMISSION_TEST_SITE_URL,
		PERMISSION_TEST_ELASTICSEARCH,
		PERMISSION_TEST_S3,
		PERMISSION_RELOAD_CONFIG,
		PERMISSION_INVALIDATE_CACHES,
		PERMISSION_RECYCLE_DATABASE_CONNECTIONS,
		PERMISSION_PURGE_ELASTICSEARCH_INDEXES,
		PERMISSION_TEST_EMAIL,
		PERMISSION_CREATE_ELASTICSEARCH_POST_INDEXING_JOB,
		PERMISSION_CREATE_ELASTICSEARCH_POST_AGGREGATION_JOB,
		PERMISSION_READ_ELASTICSEARCH_POST_INDEXING_JOB,
		PERMISSION_READ_ELASTICSEARCH_POST_AGGREGATION_JOB,
		PERMISSION_PURGE_BLEVE_INDEXES,
		PERMISSION_CREATE_POST_BLEVE_INDEXES_JOB,
		PERMISSION_CREATE_LDAP_SYNC_JOB,
		PERMISSION_READ_LDAP_SYNC_JOB,
		PERMISSION_TEST_LDAP,
		PERMISSION_INVALIDATE_EMAIL_INVITE,
		PERMISSION_GET_SAML_METADATA_FROM_IDP,
		PERMISSION_ADD_SAML_PUBLIC_CERT,
		PERMISSION_ADD_SAML_PRIVATE_CERT,
		PERMISSION_ADD_SAML_IDP_CERT,
		PERMISSION_REMOVE_SAML_PUBLIC_CERT,
		PERMISSION_REMOVE_SAML_PRIVATE_CERT,
		PERMISSION_REMOVE_SAML_IDP_CERT,
		PERMISSION_GET_SAML_CERT_STATUS,
		PERMISSION_ADD_LDAP_PUBLIC_CERT,
		PERMISSION_ADD_LDAP_PRIVATE_CERT,
		PERMISSION_REMOVE_LDAP_PUBLIC_CERT,
		PERMISSION_REMOVE_LDAP_PRIVATE_CERT,
		PERMISSION_GET_ANALYTICS,
		PERMISSION_GET_LOGS,
	}

	TeamScopedPermissions := []*Permission{
		PERMISSION_INVITE_USER,
		PERMISSION_MANAGE_SLASH_COMMANDS,
		PERMISSION_MANAGE_OTHERS_SLASH_COMMANDS,
		PERMISSION_MANAGE_INCOMING_WEBHOOKS,
		PERMISSION_MANAGE_OUTGOING_WEBHOOKS,
		PERMISSION_MANAGE_OTHERS_INCOMING_WEBHOOKS,
		PERMISSION_MANAGE_OTHERS_OUTGOING_WEBHOOKS,
		PERMISSION_VIEW_MEMBERS,
	}

	ChannelScopedPermissions := []*Permission{
		PERMISSION_USE_SLASH_COMMANDS,
		PERMISSION_ADD_REACTION,
		PERMISSION_REMOVE_REACTION,
		PERMISSION_REMOVE_OTHERS_REACTIONS,
		PERMISSION_UPLOAD_FILE,
		PERMISSION_CREATE_POST,
		PERMISSION_CREATE_POST_PUBLIC,
		PERMISSION_CREATE_POST_EPHEMERAL,
		PERMISSION_EDIT_POST,
		PERMISSION_EDIT_OTHERS_POSTS,
		PERMISSION_DELETE_POST,
		PERMISSION_DELETE_OTHERS_POSTS,
	}

	AllPermissions = []*Permission{}
	AllPermissions = append(AllPermissions, SystemScopedPermissionsMinusSysconsole...)
	AllPermissions = append(AllPermissions, TeamScopedPermissions...)
	AllPermissions = append(AllPermissions, ChannelScopedPermissions...)
	AllPermissions = append(AllPermissions, SysconsoleReadPermissions...)
	AllPermissions = append(AllPermissions, SysconsoleWritePermissions...)
}

func init() {
	initializePermissions()
	initializeSaleorPermission()
}
