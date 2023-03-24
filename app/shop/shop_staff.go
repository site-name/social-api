package shop

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// ShopStaffRelationByShopIDAndStaffID finds if there is a relationship betwwen given user and given shop
func (a *ServiceShop) ShopStaffRelationByShopIDAndStaffID(shopID string, staffID string) (*model.ShopStaffRelation, *model.AppError) {
	relations, err := a.srv.Store.ShopStaff().FilterByOptions(&model.ShopStaffRelationFilterOptions{
		ShopID:  squirrel.Eq{store.ShopStaffTableName + ".ShopID": shopID},
		StaffID: squirrel.Eq{store.ShopStaffTableName + ".StaffID": staffID},
	})
	var statusCode int
	errMsg := ""

	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(relations) == 0 {
		statusCode = http.StatusNotFound
	}
	if statusCode != 0 {
		return nil, model.NewAppError("ShopStaffRelationByShopIDAndStaffID", "app.shop.shop_staffs_by_options.app_error", nil, errMsg, statusCode)
	}

	return relations[0], nil
}
