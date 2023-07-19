package order

import (
	"net/http"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// UpsertOrderLine depends on given orderLine's Id property to decide update order save it
func (a *ServiceOrder) UpsertOrderLine(transaction *gorm.DB, orderLine *model.OrderLine) (*model.OrderLine, *model.AppError) {
	orderLine, err := a.srv.Store.OrderLine().Upsert(transaction, orderLine)
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
func (a *ServiceOrder) DeleteOrderLines(orderLineIDs []string) *model.AppError {
	err := a.srv.Store.OrderLine().BulkDelete(orderLineIDs)
	if err != nil {
		return model.NewAppError("DeleteOrderLines", "app.order.error_deleting_order_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// OrderLinesByOption returns a list of order lines by given option
func (a *ServiceOrder) OrderLinesByOption(option *model.OrderLineFilterOption) (model.OrderLines, *model.AppError) {
	orderLines, err := a.srv.Store.OrderLine().FilterbyOption(option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(orderLines) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("OrderLinesByOption", "app.order.error_finding_order_lines_by_option.app_error", nil, errMessage, statusCode)
	}

	return orderLines, nil
}

// AllDigitalOrderLinesOfOrder finds all order lines belong to given order, and are digital products
func (a *ServiceOrder) AllDigitalOrderLinesOfOrder(orderID string) ([]*model.OrderLine, *model.AppError) {
	orderLines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		OrderID: squirrel.Eq{model.OrderLineTableName + ".OrderID": orderID},
	})
	if appErr != nil {
		return nil, appErr
	}

	var (
		digitalOrderLines []*model.OrderLine
		appError          *model.AppError
		mut               sync.Mutex
		wg                sync.WaitGroup
	)

	setAppError := func(err *model.AppError) {
		mut.Lock()
		if err != nil && appError == nil {
			appError = err
		}
		mut.Unlock()
	}

	wg.Add(len(orderLines))

	for _, orderLine := range orderLines {
		go func(anOrderLine *model.OrderLine) {
			orderLineIsDigital, appErr := a.OrderLineIsDigital(anOrderLine)
			if appErr != nil {
				setAppError(appErr)
			} else {
				if orderLineIsDigital {

					mut.Lock()
					digitalOrderLines = append(digitalOrderLines, anOrderLine)
					mut.Unlock()

				}
			}

			wg.Done()
		}(orderLine)
	}

	wg.Wait()

	if appError != nil {
		return nil, appError
	}

	return digitalOrderLines, nil
}

// OrderLineById returns an order line byt given orderLineID
func (a *ServiceOrder) OrderLineById(orderLineID string) (*model.OrderLine, *model.AppError) {
	orderLine, err := a.srv.Store.OrderLine().Get(orderLineID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("OrderLineById", "app.order.missing_order_line.app_error", nil, err.Error(), statusCode)
	}

	return orderLine, nil
}

// OrderLineIsDigital Check if a variant is digital and contains digital content.
func (a *ServiceOrder) OrderLineIsDigital(orderLine *model.OrderLine) (bool, *model.AppError) {
	if orderLine.VariantID == nil {
		return false, nil
	}

	// check if the related product type is digital does not require shipping:
	productVariantIsDigital, appErr := a.srv.ProductService().ProductVariantIsDigital(*orderLine.VariantID)
	if appErr != nil {
		return false, appErr
	}

	var orderLineProductVariantHasDigitalContent bool

	// check if there is a digital content accompanies order line's product variant:
	digitalContent, appErr := a.srv.ProductService().DigitalContentbyOption(&model.DigitalContentFilterOption{
		ProductVariantID: squirrel.Eq{model.DigitalContentTableName + ".ProductVariantID": *orderLine.VariantID},
	})
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
func (a *ServiceOrder) BulkUpsertOrderLines(transaction *gorm.DB, orderLines []*model.OrderLine) ([]*model.OrderLine, *model.AppError) {
	orderLines, err := a.srv.Store.OrderLine().BulkUpsert(transaction, orderLines)
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
