package account

import (
	"context"
	"net/http"
	"strings"

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
	return a.RolesGrantPermission(session.GetUserRoles(), permission.Id)
}

// SessionHasPermissionToAny checks if current user has atleast one of given permissions
func (a *ServiceAccount) SessionHasPermissionToAny(session *model.Session, permissions []*model.Permission) bool {
	for _, perm := range permissions {
		if a.SessionHasPermissionTo(session, perm) {
			return true
		}
	}
	return false
}

// SessionHasPermissionToAll checks if given session has all given permissions
func (a *ServiceAccount) SessionHasPermissionToAll(session *model.Session, permissions []*model.Permission) bool {
	for _, perm := range permissions {
		if !a.SessionHasPermissionTo(session, perm) {
			return false
		}
	}

	return true
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

	return a.RolesGrantPermission(user.GetRoles(), permission.Id)
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

// RolesGrantPermission gets all model.Role with given roleNames.
// Then checks if one of these model.Role satisfies:
//
// 1) Not deleted
//
// 2) one item in the role's Permissions is equal to given permissionId
func (a *ServiceAccount) RolesGrantPermission(roleNames []string, permissionId string) bool {
	roles, err := a.GetRolesByNames(roleNames)
	if err != nil {
		// This should only happen if something is very broken. We can't realistically
		// recover the situation, so deny permission and log an error.
		slog.Error("Failed to get roles from database with role names: "+strings.Join(roleNames, ",")+" ", slog.Err(err))
		return false
	}

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
