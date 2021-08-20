package shop

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/store"
)

// ShopStaffRelationByShopIDAndStaffID finds if there is a relationship betwwen given user and given shop
func (a *AppShop) ShopStaffRelationByShopIDAndStaffID(shopID string, staffID string) (*shop.ShopStaffRelation, *model.AppError) {
	relation, err := a.app.Srv().Store.ShopStaff().FilterByShopAndStaff(shopID, staffID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ShopStaffRelationByShopIDAndStaffID", "app.shop.shop_staff_by_shopId_and_staff_id_missing.app_error", err)
	}

	return relation, nil
}
