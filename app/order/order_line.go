package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

// AllDigitalOrderLinesOfOrder finds all order lines belong to given order, and are digital products
func (a *AppOrder) AllDigitalOrderLinesOfOrder(orderID string) ([]*order.OrderLine, *model.AppError) {
	orderLines, appErr := a.GetAllOrderLinesByOrderId(orderID)
	if appErr != nil {
		return nil, appErr
	}

	var (
		digitalOrderLines []*order.OrderLine
		appError          *model.AppError
	)
	setAppError := func(err *model.AppError) {
		a.mutex.Lock()
		if err != nil && appError == nil {
			appError = err
		}
		a.mutex.Unlock()
	}

	a.wg.Add(len(orderLines))

	for _, orderLine := range orderLines {
		go func(line *order.OrderLine) {
			orderLineIsDigital, appErr := a.OrderLineIsDiagital(line)
			if appErr != nil {
				setAppError(appErr)
			} else {
				a.mutex.Lock()
				if orderLineIsDigital {
					digitalOrderLines = append(digitalOrderLines, orderLine)
				}
				a.mutex.Unlock()
			}

			a.wg.Done()
		}(orderLine)
	}

	a.wg.Wait()

	if appError != nil {
		return nil, appError
	}

	return digitalOrderLines, nil
}

// GetAllOrderLinesByOrderId finds all order lines belong to given order
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
