package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
)

// PreOrderAllocationsByOptions returns a list of preorder allocations filtered using given options
func (s *ServiceWarehouse) PreOrderAllocationsByOptions(options *warehouse.PreorderAllocationFilterOption) ([]*warehouse.PreorderAllocation, *model.AppError) {
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
