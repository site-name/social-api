package util

import (
	"github.com/sitename/sitename/model"
)

func SetRolePermissionsFromConfig(roles map[string]*model.Role, cfg *model.Config) map[string]*model.Role {
	if !*cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations {
		roles[model.TEAM_USER_ROLE_ID].Permissions = append(
			roles[model.TEAM_USER_ROLE_ID].Permissions,
			model.PERMISSION_MANAGE_INCOMING_WEBHOOKS.Id,
			model.PERMISSION_MANAGE_OUTGOING_WEBHOOKS.Id,
			model.PERMISSION_MANAGE_SLASH_COMMANDS.Id,
		)
		roles[model.SYSTEM_USER_ROLE_ID].Permissions = append(
			roles[model.SYSTEM_USER_ROLE_ID].Permissions,
			model.PERMISSION_MANAGE_OAUTH.Id,
		)
	}

	switch *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_RestrictPostDelete {
	case model.PERMISSIONS_DELETE_POST_ALL:
		roles[model.CHANNEL_USER_ROLE_ID].Permissions = append(
			roles[model.CHANNEL_USER_ROLE_ID].Permissions,
			model.PERMISSION_DELETE_POST.Id,
		)
		roles[model.TEAM_ADMIN_ROLE_ID].Permissions = append(
			roles[model.TEAM_ADMIN_ROLE_ID].Permissions,
			model.PERMISSION_DELETE_POST.Id,
			model.PERMISSION_DELETE_OTHERS_POSTS.Id,
		)
	case model.PERMISSIONS_DELETE_POST_TEAM_ADMIN:
		roles[model.TEAM_ADMIN_ROLE_ID].Permissions = append(
			roles[model.TEAM_ADMIN_ROLE_ID].Permissions,
			model.PERMISSION_DELETE_POST.Id,
			model.PERMISSION_DELETE_OTHERS_POSTS.Id,
		)
	}

	switch *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_AllowEditPost {
	case model.ALLOW_EDIT_POST_ALWAYS, model.ALLOW_EDIT_POST_TIME_LIMIT:
		roles[model.CHANNEL_USER_ROLE_ID].Permissions = append(
			roles[model.CHANNEL_USER_ROLE_ID].Permissions,
			model.PERMISSION_EDIT_POST.Id,
		)
		roles[model.SYSTEM_ADMIN_ROLE_ID].Permissions = append(
			roles[model.SYSTEM_ADMIN_ROLE_ID].Permissions,
			model.PERMISSION_EDIT_POST.Id,
		)
	}

	return roles
}
