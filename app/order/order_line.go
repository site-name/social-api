package order

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

// UpsertOrderLine depends on given orderLine's Id property to decide update order save it
func (a *AppOrder) UpsertOrderLine(orderLine *order.OrderLine) (*order.OrderLine, *model.AppError) {
	orderLine, err := a.Srv().Store.OrderLine().Upsert(orderLine)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok { // this not found error is caused by Get method
			status = http.StatusNotFound
		}
		return nil, model.NewAppError("UpsertOrderLine", "app.order.error_upserting_order_line.app_error", nil, err.Error(), status)
	}

	return orderLine, nil
}

// DeleteOrderLines perform bulk delete given order lines
func (a *AppOrder) DeleteOrderLines(orderLineIDs []string) *model.AppError {
	// validate given ids
	for _, id := range orderLineIDs {
		if !model.IsValidId(id) {
			return model.NewAppError("DeleteOrderLines", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderLineIDs"}, "", http.StatusBadRequest)
		}
	}

	err := a.Srv().Store.OrderLine().BulkDelete(orderLineIDs)
	if err != nil {
		return model.NewAppError("DeleteOrderLines", "app.order.error_deleting_order_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// OrderLinesByOption returns a list of order lines by given option
func (a *AppOrder) OrderLinesByOption(option *order.OrderLineFilterOption) ([]*order.OrderLine, *model.AppError) {
	orderLines, err := a.Srv().Store.OrderLine().FilterbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("OrderLinesByOption", "app.order.error_finding_order_lines_by_option.app_error", err)
	}

	return orderLines, nil
}

// AllDigitalOrderLinesOfOrder finds all order lines belong to given order, and are digital products
func (a *AppOrder) AllDigitalOrderLinesOfOrder(orderID string) ([]*order.OrderLine, *model.AppError) {
	orderLines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: orderID,
			},
		},
	})
	if appErr != nil {
		appErr.Where = "AllDigitalOrderLinesOfOrder"
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
		go func(anOrderLine *order.OrderLine) {
			orderLineIsDigital, appErr := a.OrderLineIsDiagital(anOrderLine)
			if appErr != nil {
				setAppError(appErr)
			} else {
				if orderLineIsDigital {

					a.mutex.Lock()
					digitalOrderLines = append(digitalOrderLines, anOrderLine)
					a.mutex.Unlock()

				}
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

// OrderLineById returns an order line byt given orderLineID
func (a *AppOrder) OrderLineById(orderLineID string) (*order.OrderLine, *model.AppError) {
	orderLine, err := a.Srv().Store.OrderLine().Get(orderLineID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("OrderLineById", "app.order.missing_order_line.app_error", err)
	}

	return orderLine, nil
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

// BulkUpsertOrderLines perform bulk upsert given order lines
func (a *AppOrder) BulkUpsertOrderLines(orderLines []*order.OrderLine) ([]*order.OrderLine, *model.AppError) {
	orderLines, err := a.Srv().Store.OrderLine().BulkUpsert(orderLines)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok { // this error is caused by Get() method
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("BulkUpsertOrderLines", "app.order.error_bulk_update_order_lines.app_error", nil, err.Error(), statusCode)
	}

	return orderLines, nil
}
