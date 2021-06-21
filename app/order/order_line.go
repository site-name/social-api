package order

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

func (a *AppOrder) GetAllOrderLinesByOrderId(orderID string) ([]*order.OrderLine, *model.AppError) {
	lines, err := a.Srv().Store.OrderLine().GetAllByOrderID(orderID)
	if err != nil {
		var statusCode int = http.StatusInternalServerError
		var nfErr *store.ErrNotFound
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetAllOrderLinesByOrderId", "app.order.get_child_order_lines.app_error", nil, err.Error(), statusCode)
	}

	return lines, nil
}
