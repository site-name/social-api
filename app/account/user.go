package account

import (
	"context"
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

func (a *AppAccount) UserById(ctx context.Context, userID string) (*account.User, *model.AppError) {
	user, err := a.Srv().Store.User().Get(ctx, userID)
	if err != nil {
		var nfErr *store.ErrNotFound
		statusCode := http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetUserById", "app.account.missing_user.app_error", nil, "", statusCode)
	}

	return user, nil
}

func (a *AppAccount) UserSetDefaultAddress(userID, addressID, addressType string) (*account.User, *model.AppError) {
	// check if address is owned by user
	addresses, appErr := a.AddressesByUserId(userID)
	if appErr != nil {
		return nil, appErr
	}

	addressBelongToUser := false
	for _, addr := range addresses {
		if addr.Id == addressID {
			addressBelongToUser = true
		}
	}

	if !addressBelongToUser {
		return nil, model.NewAppError("UserSetDefaultAddress", userNotOwnAddress, nil, "", http.StatusForbidden)
	}

	// get user with given id
	user, appErr := a.UserById(context.Background(), userID)
	if appErr != nil {
		return nil, appErr
	}

	// set new address accordingly
	if addressType == account.ADDRESS_TYPE_BILLING {
		user.DefaultBillingAddressID = &addressID
	} else if addressType == account.ADDRESS_TYPE_SHIPPING {
		user.DefaultShippingAddressID = &addressID
	}

	// update
	userUpdate, err := a.Srv().Store.User().Update(user, false)
	if err != nil {
		if appErr, ok := (err).(*model.AppError); ok {
			return nil, appErr
		} else if errInput, ok := (err).(*store.ErrInvalidInput); ok {
			return nil, model.NewAppError(
				"UserSetDefaultAddress",
				"app.account.invalid_input.app_error",
				map[string]interface{}{
					"field": errInput.Field,
					"value": errInput.Value}, "",
				http.StatusBadRequest,
			)
		} else {
			return nil, model.NewAppError(
				"UserSetDefaultAddress",
				"app.account.update_error.app_error",
				nil, "",
				http.StatusInternalServerError,
			)
		}
	}

	return userUpdate.New, nil
}
