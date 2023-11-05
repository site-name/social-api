package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// AllocationsByOption returns all warehouse allocations filtered based on given option
func (a *ServiceWarehouse) AllocationsByOption(option *model.AllocationFilterOption) (model.Allocations, *model.AppError) {
	allocations, err := a.srv.Store.Allocation().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("AllocationByOption", "app.warehouse.error_finding_allocations_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return allocations, nil
}

// BulkUpsertAllocations upserts or inserts given allocations into database then returns them
func (a *ServiceWarehouse) BulkUpsertAllocations(transaction *gorm.DB, allocations []*model.Allocation) ([]*model.Allocation, *model.AppError) {
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
func (a *ServiceWarehouse) BulkDeleteAllocations(transaction *gorm.DB, allocationIDs []string) *model.AppError {
	err := a.srv.Store.Allocation().BulkDelete(transaction, allocationIDs)
	if err != nil {
		return model.NewAppError("BulkDeleteAllocations", "app.warehouse.error_deleting_allocations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
