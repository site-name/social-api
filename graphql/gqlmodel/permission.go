package gqlmodel

import (
	"strings"

	"github.com/sitename/sitename/model"
)

// SaleorGraphqlPermissionsToSystemPermissions converts all given graphql permission enums to a slice of system's permissions
func SaleorGraphqlPermissionsToSystemPermissions(saleorGraphqlPermissions ...PermissionEnum) []*model.Permission {
	var res = []*model.Permission{}
	for _, perm := range saleorGraphqlPermissions {
		res = append(res, SaleorGraphqlPermissionToSystemPermission(perm))
	}

	return res
}

// SaleorGraphqlPermissionToSystemPermission converts given graphql saleor permission enum to a system permission
func SaleorGraphqlPermissionToSystemPermission(saleorPermission PermissionEnum) *model.Permission {
	for _, perm := range model.SaleorPermissions {
		if perm.Id == strings.ToLower(string(saleorPermission)) {
			return perm
		}
	}

	return nil
}
