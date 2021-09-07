package gqlmodel

import "github.com/sitename/sitename/model"

// GraphqlPermissionsToSystemPermissions converts all given graphql permission enums to a slice of system's permissions
func GraphqlPermissionsToSystemPermissions(gqlPermissions ...PermissionEnum) []*model.Permission {
	var res = []*model.Permission{}

	for _, inPerm := range gqlPermissions {
		switch inPerm {
		case PermissionEnumManageUsers:
			res = append(res, model.PermissionManageUsers)
		case PermissionEnumManageStaff:
			res = append(res, model.PermissionManageStaff)
		case PermissionEnumManageApps:
			res = append(res, model.PermissionManageApps)
		case PermissionEnumManageChannels:
			res = append(res, model.PermissionManageChannels)
		case PermissionEnumManageDiscounts:
			res = append(res, model.PermissionManageDiscounts)
		case PermissionEnumManagePlugins:
			res = append(res, model.PermissionManagePlugins)
		case PermissionEnumManageGiftCard:
			res = append(res, model.PermissionManageGiftcard)
		case PermissionEnumManageMenus:
			res = append(res, model.PermissionManageMenus)
		case PermissionEnumManageOrders:
			res = append(res, model.PermissionManageOrders)
		case PermissionEnumManagePages:
			res = append(res, model.PermissionManagePages)
		case PermissionEnumManagePageTypesAndAttributes:
			res = append(res, model.PermissionManagePageTypesAndAttributes)
		case PermissionEnumHandlePayments:
			res = append(res, model.PermissionHandlePayments)
		case PermissionEnumManageProducts:
			res = append(res, model.PermissionManageProducts)
		case PermissionEnumManageProductTypesAndAttributes:
			res = append(res, model.PermissionManageProductTypesAndAttributes)
		case PermissionEnumManageShipping:
			res = append(res, model.PermissionManageShipping)
		case PermissionEnumManageSettings:
			res = append(res, model.PermissionManageSettings)
		case PermissionEnumManageTranslations:
			res = append(res, model.PermissionManageTranslations)
		case PermissionEnumManageCheckouts:
			res = append(res, model.PermissionManageCheckouts)
		}
	}

	return res
}
