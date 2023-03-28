package account

import (
	"context"
	"net/http"
	"strings"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

func (a *ServiceAccount) MakePermissionError(s *model.Session, permissions ...*model.Permission) *model.AppError {
	permissionsStr := "permission="
	for _, permission := range permissions {
		permissionsStr += permission.Id
		permissionsStr += ","
	}
	return model.NewAppError("Permissions", "api.context.permissions.app_error", nil, "userId="+s.UserId+", "+permissionsStr, http.StatusForbidden)
}

// SessionHasPermissionTo checks if this user has given permission to procceed
func (a *ServiceAccount) SessionHasPermissionTo(session *model.Session, permission *model.Permission) bool {
	if session.IsUnrestricted() {
		return true
	}
	roleNames := session.GetUserRoles()
	roles, appErr := a.GetRolesByNames(roleNames)
	if appErr != nil {
		slog.Error("Failed to get roles from database with role names: "+strings.Join(roleNames, ",")+" ", slog.Err(appErr))
		return false
	}
	return a.RolesGrantPermission(roles, permission.Id)
}

// SessionHasPermissionToAny checks if current user has atleast one of given permissions
func (a *ServiceAccount) SessionHasPermissionToAny(session *model.Session, permissions ...*model.Permission) bool {
	if session.IsUnrestricted() {
		return true
	}
	roleNames := session.GetUserRoles()
	roles, appErr := a.GetRolesByNames(roleNames)
	if appErr != nil {
		slog.Error("Failed to get roles from database with role names: "+strings.Join(roleNames, ",")+" ", slog.Err(appErr))
		return false
	}
	return lo.SomeBy(permissions, func(perm *model.Permission) bool {
		return a.RolesGrantPermission(roles, perm.Id)
	})
}

// SessionHasPermissionToAll checks if given session has all given permissions
func (a *ServiceAccount) SessionHasPermissionToAll(session *model.Session, permissions ...*model.Permission) bool {
	if session.IsUnrestricted() {
		return true
	}
	roleNames := session.GetUserRoles()
	roles, appErr := a.GetRolesByNames(roleNames)
	if appErr != nil {
		slog.Error("Failed to get roles from database with role names: "+strings.Join(roleNames, ",")+" ", slog.Err(appErr))
		return false
	}
	return lo.EveryBy(permissions, func(perm *model.Permission) bool {
		return a.RolesGrantPermission(roles, perm.Id)
	})
}

// SessionHasPermissionToUser checks if current user has permission to perform modifications to another user with Id of given userID
func (a *ServiceAccount) SessionHasPermissionToUser(session *model.Session, userID string) bool {
	if userID == "" {
		return false
	}
	if session.IsUnrestricted() {
		return true
	}

	if session.UserId == userID {
		return true
	}

	if a.SessionHasPermissionTo(session, model.PermissionEditOtherUsers) {
		return true
	}

	return false
}

// HasPermissionTo checks if an user with Id of `askingUserId` has permission of given permission
func (a *ServiceAccount) HasPermissionTo(askingUserId string, permission *model.Permission) bool {
	user, err := a.UserById(context.Background(), askingUserId)
	if err != nil {
		return false
	}
	roleNames := user.GetRoles()
	roles, appErr := a.GetRolesByNames(roleNames)
	if appErr != nil {
		slog.Error("Failed to get roles from database with role names: "+strings.Join(roleNames, ",")+" ", slog.Err(appErr))
		return false
	}

	return a.RolesGrantPermission(roles, permission.Id)
}

// HasPermissionToUser checks if an user with Id of `askingUserId` has permission to modify another user with Id of given `userID`
func (a *ServiceAccount) HasPermissionToUser(askingUserId string, userID string) bool {
	if askingUserId == userID {
		return true
	}

	if a.HasPermissionTo(askingUserId, model.PermissionEditOtherUsers) {
		return true
	}

	return false
}

func (a *ServiceAccount) RolesGrantPermission(roles []*model.Role, permissionId string) bool {
	for _, role := range roles {
		if role.DeleteAt != 0 {
			continue
		}

		for _, permission := range role.Permissions {
			if permission == permissionId {
				return true
			}
		}
	}

	return false
}
