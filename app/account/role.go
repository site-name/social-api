package account

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// GetRole get 1 model.Role from database, returns nil and concret error if a problem occur
func (a *ServiceAccount) GetRole(id string) (*model.Role, *model.AppError) {
	role, err := a.srv.Store.Role().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetRole", "app.role.get.app_error", err)
	}

	// appErr := a.srv.mergeChannelHigherScopedPermissions([]*model.Role{role})
	// if appErr != nil {
	// 	return nil, appErr
	// }

	return role, nil
}

// GetRoleByName gets a model.Role from database with given name, returns nil and concret error if a problem occur
func (s *ServiceAccount) GetRoleByName(ctx context.Context, name string) (*model.Role, *model.AppError) {
	role, nErr := s.srv.Store.Role().GetByName(ctx, name)
	if nErr != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetRoleByName", "app.role.get_by_name.app_error", nErr)
	}

	// err := s.mergeChannelHigherScopedPermissions([]*model.Role{role})
	// if err != nil {
	// 	return nil, err
	// }

	return role, nil
}

// GetRolesByNames returns a slice of model.Role by given names
func (a *ServiceAccount) GetRolesByNames(names []string) ([]*model.Role, *model.AppError) {
	roles, nErr := a.srv.Store.Role().GetByNames(names)
	if nErr != nil {
		return nil, model.NewAppError("GetRolesByNames", "app.role.get_by_names.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}

	// TODO: fixme
	// err := a.mergeChannelHigherScopedPermissions(roles)
	// if err != nil {
	// 	return nil, err
	// }

	return roles, nil
}

// mergeChannelHigherScopedPermissions updates the permissions based on the role type, whether the permission is
// moderated, and the value of the permission on the higher-scoped scheme.
// func (s *Server) mergeChannelHigherScopedPermissions(roles []*model.Role) *model.AppError {
// 	var higherScopeNamesToQuery []string

// 	for _, role := range roles {
// 		if role.SchemeManaged {
// 			higherScopeNamesToQuery = append(higherScopeNamesToQuery, role.Name)
// 		}
// 	}

// 	if len(higherScopeNamesToQuery) == 0 {
// 		return nil
// 	}

// 	higherScopedPermissionsMap, err := s.Store.Role().ChannelHigherScopedPermissions(higherScopeNamesToQuery)
// 	if err != nil {
// 		return model.NewAppError("mergeChannelHigherScopedPermissions", "app.role.get_by_names.app_error", nil, err.Error(), http.StatusInternalServerError)
// 	}

// 	for _, role := range roles {
// 		if role.SchemeManaged {
// 			if higherScopedPermissions, ok := higherScopedPermissionsMap[role.Name]; ok {
// 				role.MergeChannelHigherScopedPermissions(higherScopedPermissions)
// 			}
// 		}
// 	}

// 	return nil
// }

// mergeChannelHigherScopedPermissions updates the permissions based on the role type, whether the permission is
// moderated, and the value of the permission on the higher-scoped scheme.
// func (a *ServiceAccount) mergeChannelHigherScopedPermissions(roles []*model.Role) *model.AppError {
// 	return a.srv.mergeChannelHigherScopedPermissions(roles)
// }

func (a *ServiceAccount) PatchRole(role *model.Role, patch *model.RolePatch) (*model.Role, *model.AppError) {
	// If patch is a no-op then short-circuit the store.
	if patch.Permissions != nil && reflect.DeepEqual(*patch.Permissions, role.Permissions) {
		return role, nil
	}

	role.Patch(patch)
	role, err := a.UpdateRole(role)
	if err != nil {
		return nil, err
	}

	return role, err
}

// CreateRole takes a role struct and save it to database
func (a *ServiceAccount) CreateRole(role *model.Role) (*model.Role, *model.AppError) {
	role.Id = ""
	role.CreateAt = 0
	role.UpdateAt = 0
	role.DeleteAt = 0
	role.BuiltIn = false
	role.SchemeManaged = false

	var err error
	role, err = a.srv.Store.Role().Save(role)
	if err != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &invErr):
			return nil, model.NewAppError("CreateRole", "app.role.save.invalid_role.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model.NewAppError("CreateRole", "app.role.save.insert.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return role, nil
}

func (a *ServiceAccount) UpdateRole(role *model.Role) (*model.Role, *model.AppError) {
	savedRole, err := a.srv.Store.Role().Save(role)
	if err != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &invErr):
			return nil, model.NewAppError("UpdateRole", "app.role.save.invalid_role.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model.NewAppError("UpdateRole", "app.role.save.insert.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// TODO: need investigation

	// builtInChannelRoles := []string{
	// 	model.CHANNEL_GUEST_ROLE_ID,
	// 	model.CHANNEL_USER_ROLE_ID,
	// 	model.CHANNEL_ADMIN_ROLE_ID,
	// }

	// builtInRolesMinusChannelRoles := append(
	// 	util.RemoveStringsFromSlice(model.BuiltInSchemeManagedRoleIDs, builtInChannelRoles...),
	// 	model.NewSystemRoleIDs...,
	// )

	// if util.StringInSlice(savedRole.Name, builtInRolesMinusChannelRoles) {
	// 	return savedRole, nil
	// }

	// var roleRetrievalFunc func() ([]*model.Role, *model.AppError)

	// if util.StringInSlice(savedRole.Name, builtInChannelRoles) {
	// 	roleRetrievalFunc = func() ([]*model.Role, *model.AppError) {
	// 		roles, nErr := a.srv.Store.Role().AllChannelSchemeRoles()
	// 		if nErr != nil {
	// 			return nil, model.NewAppError("UpdateRole", "app.role.get.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	// 		}

	// 		return roles, nil
	// 	}
	// } else {
	// 	roleRetrievalFunc = func() ([]*model.Role, *model.AppError) {
	// 		roles, nErr := a.srv.Store.Role().ChannelRolesUnderTeamRole(savedRole.Name)
	// 		if nErr != nil {
	// 			return nil, model.NewAppError("UpdateRole", "app.role.get.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	// 		}

	// 		return roles, nil
	// 	}
	// }

	// impactedRoles, appErr := roleRetrievalFunc()
	// if appErr != nil {
	// 	return nil, appErr
	// }
	// impactedRoles = append(impactedRoles, role)

	// appErr = a.mergeChannelHigherScopedPermissions(impactedRoles)
	// if appErr != nil {
	// 	return nil, appErr
	// }

	// for _, ir := range impactedRoles {
	// 	if ir.Name != role.Name {
	// 		a.sendUpdatedRoleEvent(ir)
	// 	}
	// }

	return savedRole, nil
}

// CheckRolesExist get role model instances with given roleNames,
// checks if at least one db role has name contained in given roleNames.
func (a *ServiceAccount) CheckRolesExist(roleNames []string) *model.AppError {
	roles, err := a.GetRolesByNames(roleNames)
	if err != nil {
		return err
	}

	for _, name := range roleNames {
		nameFound := false
		for _, role := range roles {
			if name == role.Name {
				nameFound = true
				break
			}
		}
		if !nameFound {
			return model.NewAppError("CheckRolesExist", "app.role.check_roles_exist.role_not_found", nil, "role="+name, http.StatusBadRequest)
		}
	}

	return nil
}

func RemoveRoles(rolesToRemove []string, roles string) string {
	roleList := strings.Fields(roles)
	newRoles := make([]string, 0)

	for _, role := range roleList {
		shouldRemove := false
		for _, roleToRemove := range rolesToRemove {
			if role == roleToRemove {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			newRoles = append(newRoles, role)
		}
	}

	return strings.Join(newRoles, " ")
}