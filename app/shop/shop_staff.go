package shop

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (a *ServiceShop) StaffsByOptions(options *model.ShopStaffFilterOptions) ([]*model.ShopStaff, *model_helper.AppError) {
	staffs, err := a.srv.Store.ShopStaff().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("StaffsByOptions", "app.shop.staffs_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return staffs, nil
}

func (a *ServiceShop) ShopStaffByOptions(options *model.ShopStaffFilterOptions) (*model.ShopStaff, *model_helper.AppError) {
	relation, err := a.srv.Store.ShopStaff().GetByOptions(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("ShopStaffByOptions", "app.shop.shop_staff_by_options.app_error", nil, err.Error(), statusCode)
	}
	return relation, nil
}

func (s *ServiceShop) UserIsStaffOfShop(userID string) bool {
	relation, appErr := s.ShopStaffByOptions(&model.ShopStaffFilterOptions{
		Conditions: squirrel.Eq{model.ShopStaffTableName + ".StaffID": userID},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			return false
		}
		slog.Error("failed to find shop-staff", slog.Err(appErr))
		return false
	}

	return relation != nil && relation.EndAt == nil
}
