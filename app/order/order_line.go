package order

import (
	"net/http"
	"sync/atomic"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gorm.io/gorm"
)

// UpsertOrderLine depends on given orderLine's Id property to decide update order save it
func (a *ServiceOrder) UpsertOrderLine(transaction boil.ContextTransactor, orderLine *model.OrderLine) (*model.OrderLine, *model_helper.AppError) {
	orderLine, err := a.srv.Store.OrderLine().Upsert(transaction, orderLine)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok { // this not found error is caused by Get method
			status = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("UpsertOrderLine", "app.order.error_upserting_order_line.app_error", nil, err.Error(), status)
	}

	return orderLine, nil
}

// DeleteOrderLines perform bulk delete given order lines
func (a *ServiceOrder) DeleteOrderLines(tx *gorm.DB, orderLineIDs []string) *model_helper.AppError {
	err := a.srv.Store.OrderLine().BulkDelete(tx, orderLineIDs)
	if err != nil {
		return model_helper.NewAppError("DeleteOrderLines", "app.order.error_deleting_order_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// OrderLinesByOption returns a list of order lines by given option
func (a *ServiceOrder) OrderLinesByOption(option *model.OrderLineFilterOption) (model.OrderLineSlice, *model_helper.AppError) {
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
		return nil, model_helper.NewAppError("OrderLinesByOption", "app.order.error_finding_order_lines_by_option.app_error", nil, errMessage, statusCode)
	}

	return orderLines, nil
}

// AllDigitalOrderLinesOfOrder finds all order lines belong to given order, and are digital products
func (a *ServiceOrder) AllDigitalOrderLinesOfOrder(orderID string) (model.OrderLineSlice, *model_helper.AppError) {
	orderLines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": orderID},
	})
	if appErr != nil {
		return nil, appErr
	}

	var (
		digitalOrderLines model.OrderLineSlice
		atomicValue       atomic.Int32
		appErrChan        = make(chan *model_helper.AppError)
		dititalLineChan   = make(chan *model.OrderLine) // every digital orderlines are sent to this channel
	)
	defer close(appErrChan)
	defer close(dititalLineChan)
	atomicValue.Add(int32(len(orderLines)))

	for _, orderLine := range orderLines {
		go func(anOrderLine *model.OrderLine) {
			defer atomicValue.Add(-1)

			orderLineIsDigital, appErr := a.OrderLineIsDigital(anOrderLine)
			if appErr != nil {
				appErrChan <- appErr
				return
			}

			if orderLineIsDigital {
				dititalLineChan <- anOrderLine
			}

		}(orderLine)
	}

	for atomicValue.Load() != 0 {
		select {
		case appErr := <-appErrChan:
			return nil, appErr
		case line := <-dititalLineChan:
			digitalOrderLines = append(digitalOrderLines, line)
		}
	}

	return digitalOrderLines, nil
}

// OrderLineById returns an order line byt given orderLineID
func (a *ServiceOrder) OrderLineById(orderLineID string) (*model.OrderLine, *model_helper.AppError) {
	orderLine, err := a.srv.Store.OrderLine().Get(orderLineID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("OrderLineById", "app.order.missing_order_line.app_error", nil, err.Error(), statusCode)
	}

	return orderLine, nil
}

// OrderLineIsDigital Check if a variant is digital and contains digital content.
func (a *ServiceOrder) OrderLineIsDigital(orderLine *model.OrderLine) (bool, *model_helper.AppError) {
	if orderLine.VariantID == nil {
		return false, nil
	}

	// check if the related product type is digital does not require shipping:
	productVariantIsDigital, appErr := a.srv.ProductService().ProductVariantIsDigital(*orderLine.VariantID)
	if appErr != nil {
		return false, appErr
	}

	// check if there is a digital content accompanies order line's product variant:
	digitalContent, appErr := a.srv.ProductService().DigitalContentbyOption(&model.DigitalContentFilterOption{
		Conditions: squirrel.Eq{model.DigitalContentTableName + ".ProductVariantID": *orderLine.VariantID},
	})
	if appErr != nil {
		return false, appErr
	}

	return productVariantIsDigital && digitalContent != nil, nil
}

// BulkUpsertOrderLines perform bulk upsert given order lines
func (a *ServiceOrder) BulkUpsertOrderLines(transaction boil.ContextTransactor, orderLines model.OrderLineSlice) (model.OrderLineSlice, *model_helper.AppError) {
	orderLines, err := a.srv.Store.OrderLine().BulkUpsert(transaction, orderLines)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok { // this error is caused by Get() method
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("BulkUpsertOrderLines", "app.order.error_bulk_update_order_lines.app_error", nil, err.Error(), statusCode)
	}

	return orderLines, nil
}
