package app

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

// GetAddressById returns address with given id. If not found returns nil and concret error
func (a *App) GetAddressById(id string) (*account.Address, *model.AppError) {
	address, err := a.srv.Store.Address().Get(id)
	if err != nil {
		var nfErr *store.ErrNotFound
		var statusCode int = http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetAddressById", "app.address.address_by_id.app_error", nil, err.Error(), statusCode)
	}

	return address, nil
}
