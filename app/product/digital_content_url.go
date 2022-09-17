package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// UpsertDigitalContentURL create a digital content url then returns it
func (a *ServiceProduct) UpsertDigitalContentURL(contentURL *model.DigitalContentUrl) (*model.DigitalContentUrl, *model.AppError) {
	contentURL, err := a.srv.Store.DigitalContentUrl().Upsert(contentURL)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("UpsertDigitalContentURL", "app.product.error_upserting_content_url.app_error", nil, err.Error(), statusCode)
	}

	return contentURL, nil
}
