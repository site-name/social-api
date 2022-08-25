package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// AllocationsByOption returns all warehouse allocations filtered based on given option
func (a *ServiceWarehouse) AllocationsByOption(transaction store_iface.SqlxTxExecutor, option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, *model.AppError) {
	allocations, err := a.srv.Store.Allocation().FilterByOption(transaction, option)
	var (
		statusCode   int
		errorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	} else if len(allocations) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("AllocationByOption", "app.warehouse.error_finding_allocations_by_option.app_error", nil, errorMessage, statusCode)
	}

	return allocations, nil
}

// BulkUpsertAllocations upserts or inserts given allocations into database then returns them
func (a *ServiceWarehouse) BulkUpsertAllocations(transaction store_iface.SqlxTxExecutor, allocations []*warehouse.Allocation) ([]*warehouse.Allocation, *model.AppError) {
	allocations, err := a.srv.Store.Allocation().BulkUpsert(transaction, allocations)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		var statusCode int = http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("BulkUpsertAllocations", "app.warehouse.error_upserting_allocations.app_error", nil, err.Error(), statusCode)
	}

	return allocations, nil
}

// BulkDeleteAllocations performs bulk delete given allocations.
// If non-nil transaction is provided, perform bulk delete operation within it.
func (a *ServiceWarehouse) BulkDeleteAllocations(transaction store_iface.SqlxTxExecutor, allocationIDs []string) *model.AppError {
	err := a.srv.Store.Allocation().BulkDelete(transaction, allocationIDs)
	if err != nil {
		return model.NewAppError("BulkDeleteAllocations", "app.warehouse.error_deleting_allocations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
