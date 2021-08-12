package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

func (a *AppOrder) GetAllOrderLinesByOrderId(orderID string) ([]*order.OrderLine, *model.AppError) {
	lines, err := a.Srv().Store.OrderLine().GetAllByOrderID(orderID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetAllOrderLinesByOrderId", "app.order.error_finding_child_order_lines.app_error", err)
	}

	return lines, nil
}

func (a *AppOrder) OrderLineById(orderLineID string) (*order.OrderLine, *model.AppError) {
	odrLine, err := a.Srv().Store.OrderLine().Get(orderLineID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("OrderLineById", "app.order.missing_order_line.app_error", err)
	}

	return odrLine, nil
}

// OrderLineIsDiagital Check if a variant is digital and contains digital content.
func (a *AppOrder) OrderLineIsDiagital(orderLine *order.OrderLine) (bool, *model.AppError) {
	if orderLine.VariantID == nil {
		return false, nil
	}

	// check if the related product type is digital does not require shipping:
	productVariantIsDigital, appErr := a.ProductApp().ProductVariantIsDigital(*orderLine.VariantID)
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

	return productVariantIsDigital && orderLineProductVariantHasDigitalContent, nil
}
