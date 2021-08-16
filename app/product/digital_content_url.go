package product

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// CreateDigitalContentURL create a digital content url then returns it
func (a *AppProduct) CreateDigitalContentURL(contentURL *product_and_discount.DigitalContentUrl) (*product_and_discount.DigitalContentUrl, *model.AppError) {
	contentURL, err := a.Srv().Store.DigitalContentUrl().Save(contentURL)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		if invalidInputErr, ok := err.(*store.ErrInvalidInput); ok { // this happens when duplicate Line`ID
			return nil, model.NewAppError("CreateDigitalContentURL", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": invalidInputErr.Field}, invalidInputErr.Error(), http.StatusBadRequest)
		}
	}

	return contentURL, nil
}
