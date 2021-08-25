package warehouse

import (
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

// AllocationsByOption returns all warehouse allocations filtered based on given option
func (a *AppWarehouse) AllocationsByOption(transaction *gorp.Transaction, option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, *model.AppError) {
	allocations, err := a.Srv().Store.Allocation().FilterByOption(transaction, option)
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
func (a *AppWarehouse) BulkUpsertAllocations(transaction *gorp.Transaction, allocations []*warehouse.Allocation) ([]*warehouse.Allocation, *model.AppError) {
	allocations, err := a.Srv().Store.Allocation().BulkUpsert(transaction, allocations)
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
