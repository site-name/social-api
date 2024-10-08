package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// PreOrderAllocationsByOptions returns a list of preorder allocations filtered using given options
func (s *ServiceWarehouse) PreOrderAllocationsByOptions(options *model.PreorderAllocationFilterOption) (model.PreorderAllocations, *model_helper.AppError) {
	allocations, err := s.srv.Store.PreorderAllocation().FilterByOption(options)
	if err != nil {
		return nil, model_helper.NewAppError("PreOrderAllocationsByOptions", "app.warehouse.error_finding_preorder_allocations_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return allocations, nil
}

// DeletePreorderAllocations tells store to delete given preorder allocations
func (s *ServiceWarehouse) DeletePreorderAllocations(transaction boil.ContextTransactor, preorderAllocationIDs ...string) *model_helper.AppError {
	err := s.srv.Store.PreorderAllocation().Delete(transaction, preorderAllocationIDs...)
	if err != nil {
		return model_helper.NewAppError("DeletePreorderAllocations", "app.warehouse.error_deleting_preorder_allocations_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// BulkCreate tells store to insert given preorder allocations into database then returns them
func (s *ServiceWarehouse) BulkCreate(transaction boil.ContextTransactor, preorderAllocations []*model.PreorderAllocation) ([]*model.PreorderAllocation, *model_helper.AppError) {
	allocations, err := s.srv.Store.PreorderAllocation().BulkCreate(transaction, preorderAllocations)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("BulkCreate", "app.warehouse.error_bulk_creating_preorder_allocations.app_error", nil, err.Error(), statusCode)
	}

	return allocations, nil
}
