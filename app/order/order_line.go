package order

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

func (a *AppOrder) GetAllOrderLinesByOrderId(orderID string) ([]*order.OrderLine, *model.AppError) {
	// validate orderID
	if !model.IsValidId(orderID) {
		return nil, model.NewAppError("GetAllOrderLinesByOrderId", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderID"}, "", http.StatusBadRequest)
	}

	lines, err := a.Srv().Store.OrderLine().GetAllByOrderID(orderID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetAllOrderLinesByOrderId", "app.order.error_finding_child_order_lines.app_error", err)
	}

	return lines, nil
}

func (a *AppOrder) OrderLineById(orderLineID string) (*order.OrderLine, *model.AppError) {
	// validate orderID
	if !model.IsValidId(orderLineID) {
		return nil, model.NewAppError("GetAllOrderLinesByOrderId", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderLineID"}, "", http.StatusBadRequest)
	}

	odrLine, err := a.Srv().Store.OrderLine().Get(orderLineID)
	if err != nil {
		var nfErr *store.ErrNotFound
		statusCode := http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("OrderLineById", "app.order.missing_order_line.app_error", nil, err.Error(), statusCode)
	}

	return odrLine, nil
}

// OrderLineIsDiagital Check if a variant is digital and contains digital content.
func (a *AppOrder) OrderLineIsDiagital(orderLine *order.OrderLine) (bool, *model.AppError) {
	if orderLine.VariantID == nil {
		return false, nil
	}

	// check if the related product type is digital does not require shipping:
	productVariantOfOrderLineIsDigital, appErr := a.ProductApp().ProductVariantIsDigital(*orderLine.VariantID)
	if appErr != nil {
		return false, appErr
	}

	var orderLineProductVariantHasDigitalContent bool

	// check if there is a digital content accompanies order line's product variant:
	digitalContent, appErr := a.ProductApp().DigitalContentByProductVariantID(*orderLine.VariantID)
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			orderLineProductVariantHasDigitalContent = false
		} else {
			return false, appErr
		}
	}

	if digitalContent != nil {
		orderLineProductVariantHasDigitalContent = true
	}

	return productVariantOfOrderLineIsDigital && orderLineProductVariantHasDigitalContent, nil
}
