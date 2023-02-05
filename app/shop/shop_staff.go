package shop

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// ShopStaffRelationByShopIDAndStaffID finds if there is a relationship betwwen given user and given shop
func (a *ServiceShop) ShopStaffRelationByShopIDAndStaffID(shopID string, staffID string) (*model.ShopStaffRelation, *model.AppError) {
	relation, err := a.srv.Store.ShopStaff().FilterByShopAndStaff(shopID, staffID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ShopStaffRelationByShopIDAndStaffID", "app.shop.shop_staff_by_shopId_and_staff_id_missing.app_error", nil, err.Error(), statusCode)
	}

	return relation, nil
}
