package account

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type AppAccount struct {
	app.AppIface
}

func init() {
	app.RegisterAccountApp(func(a app.AppIface) sub_app_iface.AccountApp {
		return &AppAccount{a}
	})
}

func (a *AppAccount) GetAddressById(id string) (*account.Address, *model.AppError) {
	address, err := a.Srv().Store.Address().Get(id)
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
