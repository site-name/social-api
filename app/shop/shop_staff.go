package shop

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (a *ServiceShop) ShopStaffByOptions(options *model.ShopStaffRelationFilterOptions) (*model.ShopStaff, *model.AppError) {
	relation, err := a.srv.Store.ShopStaff().GetByOptions(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("ShopStaffByOptions", "app.shop.shop_staff_by_options.app_error", nil, err.Error(), statusCode)
	}
	return relation, nil
}

func (s *ServiceShop) UserIsStaffOfShop(userID string) bool {
	relation, appErr := s.ShopStaffByOptions(&model.ShopStaffRelationFilterOptions{
		StaffID: squirrel.Eq{store.ShopStaffTableName + ".StaffID": userID},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			return false
		}
		slog.Error("failed to find shop-staff relation and staffId", slog.Err(appErr))
		return false
	}

	return relation != nil && relation.EndAt == nil
}
