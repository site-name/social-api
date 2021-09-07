package gqlmodel

import "github.com/sitename/sitename/model"

// GraphqlPermissionsToSystemPermissions converts all given graphql permission enums to a slice of system's permissions
func GraphqlPermissionsToSystemPermissions(gqlPermissions ...PermissionEnum) []*model.Permission {
	var res = []*model.Permission{}
	for _, inPerm := range gqlPermissions {
		res = append(res, GraphqlPermissionToSystemPermission(inPerm))
	}

	return res
}

// GraphqlPermissionToSystemPermission converts given graphql permission enum to a system's permission
func GraphqlPermissionToSystemPermission(gqlPermission PermissionEnum) *model.Permission {

	switch gqlPermission {
	case PermissionEnumManageUsers:
		return model.PermissionManageUsers
	case PermissionEnumManageStaff:
		return model.PermissionManageStaff
	case PermissionEnumManageApps:
		return model.PermissionManageApps
	case PermissionEnumManageChannels:
		return model.PermissionManageChannels
	case PermissionEnumManageDiscounts:
		return model.PermissionManageDiscounts
	case PermissionEnumManagePlugins:
		return model.PermissionManagePlugins
	case PermissionEnumManageGiftCard:
		return model.PermissionManageGiftcard
	case PermissionEnumManageMenus:
		return model.PermissionManageMenus
	case PermissionEnumManageOrders:
		return model.PermissionManageOrders
	case PermissionEnumManagePages:
		return model.PermissionManagePages
	case PermissionEnumManagePageTypesAndAttributes:
		return model.PermissionManagePageTypesAndAttributes
	case PermissionEnumHandlePayments:
		return model.PermissionHandlePayments
	case PermissionEnumManageProducts:
		return model.PermissionManageProducts
	case PermissionEnumManageProductTypesAndAttributes:
		return model.PermissionManageProductTypesAndAttributes
	case PermissionEnumManageShipping:
		return model.PermissionManageShipping
	case PermissionEnumManageSettings:
		return model.PermissionManageSettings
	case PermissionEnumManageTranslations:
		return model.PermissionManageTranslations
	case PermissionEnumManageCheckouts:
		return model.PermissionManageCheckouts

	default:
		return nil
	}
}
