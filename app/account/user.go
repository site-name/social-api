package account

import (
	"context"
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

func (a *AppAccount) GetUserById(ctx context.Context, userID string) (*account.User, *model.AppError) {
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
