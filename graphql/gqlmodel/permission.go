package gqlmodel

import (
	"strings"

	"github.com/sitename/sitename/model"
)

// SaleorGraphqlPermissionsToSystemPermission converts all given graphql permission enums to a slice of system's permissions
func SaleorGraphqlPermissionsToSystemPermission(saleorGraphqlPermissions ...PermissionEnum) []*model.Permission {
	var res = []*model.Permission{}
	for _, perm := range saleorGraphqlPermissions {
		res = append(res, SaleorGraphqlPermissionToSystemPermission(perm))
	}

	return res
}

// SaleorGraphqlPermissionToSystemPermission converts given graphql saleor permission enum to a system permission
func SaleorGraphqlPermissionToSystemPermission(saleorPermission PermissionEnum) *model.Permission {
	for _, perm := range model.SaleorPermissionEnumList {
		if perm.Id == strings.ToLower(string(saleorPermission)) {
			return perm
		}
	}

	return nil
}