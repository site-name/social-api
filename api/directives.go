package api

import (
	"context"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

type Role string

const (
	RoleSystemUser            Role = "system_user"
	RoleSystemAdmin           Role = "system_admin"
	RoleSystemUserAccessToken Role = "system_user_access_token"
	RoleSystemUserManager     Role = "system_user_manager"
	RoleSystemReadOnlyAdmin   Role = "system_read_only_admin"
	RoleSystemManager         Role = "system_manager"
	RoleShopAdmin             Role = "shop_admin"
	RoleShopStaff             Role = "shop_staff"
)

// HasRolesDirective helps check user is authenticated and has all given roles
type HasRolesDirective struct {
	Roles []Role
}

func (h *HasRolesDirective) ImplementsDirective() string {
	return "hasRoles"
}

func (h *HasRolesDirective) Validate(ctx context.Context, _ interface{}) error {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	strRoles := lo.Map(h.Roles, func(r Role, _ int) string { return string(r) })
	embedCtx.CheckAuthenticatedAndHasRoles("HasRolesDirective.Validate", strRoles...)
	return embedCtx.Err
}

type HasRoleAnyDirective struct {
	Roles []Role
}

func (h *HasRoleAnyDirective) ImplementsDirective() string {
	return "hasRoleAny"
}

func (h *HasRoleAnyDirective) Validate(ctx context.Context, _ interface{}) error {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	strRoles := lo.Map(h.Roles, func(r Role, _ int) string { return string(r) })
	embedCtx.CheckAuthenticatedAndHasRoleAny("HasRoleAnyDirective.Validate", strRoles...)
	return embedCtx.Err
}

// AuthenticatedDirective checks if user is authenticated or not
type AuthenticatedDirective struct{}

func (h *AuthenticatedDirective) ImplementsDirective() string {
	return "authenticated"
}

func (h *AuthenticatedDirective) Validate(ctx context.Context, _ interface{}) error {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	return embedCtx.Err
}

// HasPermissionsDirective checks if user is authenticated and has all given permissions
type HasPermissionsDirective struct {
	Permissions []PermissionEnum // permission ids
}

func (h *HasPermissionsDirective) ImplementsDirective() string {
	return "hasPermissions"
}

func (h *HasPermissionsDirective) Validate(ctx context.Context, _ interface{}) error {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	permissions := lo.Map(h.Permissions, func(p PermissionEnum, _ int) *model.Permission { return &model.Permission{Id: string(p)} })
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(permissions...)
	return embedCtx.Err
}

// HasPermissionAnyDirective checks if user is authenticated and has any of given permissions
type HasPermissionAnyDirective struct {
	Permissions []PermissionEnum // permission ids
}

func (h *HasPermissionAnyDirective) ImplementsDirective() string {
	return "hasPermissionAny"
}

func (h *HasPermissionAnyDirective) Validate(ctx context.Context, _ interface{}) error {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	permissions := lo.Map(h.Permissions, func(p PermissionEnum, _ int) *model.Permission { return &model.Permission{Id: string(p)} })
	embedCtx.CheckAuthenticatedAndHasPermissionToAny(permissions...)
	return embedCtx.Err
}
