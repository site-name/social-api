package warehouse

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

// AllocationsByOption returns all warehouse allocations filtered based on given option
func (a *AppWarehouse) AllocationsByOption(option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, *model.AppError) {
	allocations, err := a.Srv().Store.Allocation().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AllocationByOption", "app.warehouse.error_finding_allocations_by_option.app_error", err)
	}

	return allocations, nil
}
