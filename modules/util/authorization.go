package util

// import (
// 	"github.com/sitename/sitename/model"
// )

// func SetRolePermissionsFromConfig(roles map[string]*model.Role, cfg *model.Config) map[string]*model.Role {
// 	if !*cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations {
// 		roles[model.SYSTEM_USER_ROLE_ID].Permissions = append(
// 			roles[model.SYSTEM_USER_ROLE_ID].Permissions,
// 			model.PERMISSION_MANAGE_OAUTH.Id,
// 		)
// 	}

// 	switch *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_AllowEditPost {
// 	case model.ALLOW_EDIT_POST_ALWAYS, model.ALLOW_EDIT_POST_TIME_LIMIT:
// 		roles[model.SYSTEM_ADMIN_ROLE_ID].Permissions = append(
// 			roles[model.SYSTEM_ADMIN_ROLE_ID].Permissions,
// 			model.PERMISSION_EDIT_POST.Id,
// 		)
// 	}

// 	return roles
// }
