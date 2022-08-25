package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// PreOrderAllocationsByOptions returns a list of preorder allocations filtered using given options
func (s *ServiceWarehouse) PreOrderAllocationsByOptions(options *warehouse.PreorderAllocationFilterOption) (warehouse.PreorderAllocations, *model.AppError) {
	allocations, err := s.srv.Store.PreorderAllocation().FilterByOption(options)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(allocations) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("PreOrderAllocationsByOptions", "app.warehouse.error_finding_preorder_allocations_by_options.app_error", nil, errMessage, statusCode)
	}

	return allocations, nil
}

// DeletePreorderAllocations tells store to delete given preorder allocations
func (s *ServiceWarehouse) DeletePreorderAllocations(transaction store_iface.SqlxTxExecutor, preorderAllocationIDs ...string) *model.AppError {
	err := s.srv.Store.PreorderAllocation().Delete(transaction, preorderAllocationIDs...)
	if err != nil {
		return model.NewAppError("DeletePreorderAllocations", "app.warehouse.error_deleting_preorder_allocations_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// BulkCreate tells store to insert given preorder allocations into database then returns them
func (s *ServiceWarehouse) BulkCreate(transaction store_iface.SqlxTxExecutor, preorderAllocations []*warehouse.PreorderAllocation) ([]*warehouse.PreorderAllocation, *model.AppError) {
	allocations, err := s.srv.Store.PreorderAllocation().BulkCreate(transaction, preorderAllocations)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("BulkCreate", "app.warehouse.error_bulk_creating_preorder_allocations.app_error", nil, err.Error(), statusCode)
	}

	return allocations, nil
}
