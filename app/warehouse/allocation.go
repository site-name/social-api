package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
)

// AllocationsByOption returns all warehouse allocations filtered based on given option
func (a *AppWarehouse) AllocationsByOption(option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, *model.AppError) {
	allocations, err := a.Srv().Store.Allocation().FilterByOption(option)
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
