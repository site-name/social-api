package product

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

func (a *ServiceProduct) GetDefaultDigitalContentSettings(aShop model.ShopSettings) *model.ShopDefaultDigitalContentSettings {
	return &model.ShopDefaultDigitalContentSettings{
		AutomaticFulfillmentDigitalProducts: aShop.AutomaticFulfillmentDigitalProducts,
		DefaultDigitalMaxDownloads:          aShop.DefaultDigitalMaxDownloads,
		DefaultDigitalUrlValidDays:          aShop.DefaultDigitalUrlValidDays,
	}
}

// DigitalContentUrlIsValid Check if digital url is still valid for customer.
//
// It takes default settings or digital product's settings
// to check if url is still valid.
func (a *ServiceProduct) DigitalContentUrlIsValid(contentURL *model.DigitalContentUrl) (bool, *model_helper.AppError) {
	digitalContent, appErr := a.DigitalContentbyOption(&model.DigitalContentFilterOption{
		Conditions: squirrel.Eq{model.DigitalContentTableName + ".Id": contentURL.ContentID},
	})
	if appErr != nil {
		return false, appErr
	}

	var (
		urlValidDays *int
		maxDownloads *int
	)
	if *digitalContent.UseDefaultSettings {
		urlValidDays = a.srv.Config().ShopSettings.DefaultDigitalUrlValidDays
		maxDownloads = a.srv.Config().ShopSettings.DefaultDigitalMaxDownloads
	} else {
		urlValidDays = digitalContent.UrlValidDays
		maxDownloads = digitalContent.MaxDownloads
	}

	if urlValidDays != nil {
		if util.
			TimeFromMillis(contentURL.CreateAt).
			Add(time.Hour * 24 * time.Duration(*urlValidDays)).
			Before(time.Now().UTC()) {
			return false, nil
		}
	}

	if maxDownloads != nil && *maxDownloads <= contentURL.DownloadNum {
		return false, nil
	}

	return true, nil
}

func (a *ServiceProduct) IncrementDownloadCount(contentURL model.DigitalContentUrl) (*model.DigitalContentUrl, *model_helper.AppError) {
	contentURL.DownloadNum++
	updatedContentUrl, appErr := a.UpsertDigitalContentURL(&contentURL)
	if appErr != nil {
		return nil, appErr
	}

	// create order event for this download
	if contentURL.LineID != nil {
		orderLines, appErr := a.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
			Conditions: squirrel.Eq{model.OrderLineTableName + "." + model.OrderLineColumnId: *contentURL.LineID},
			Preload:    []string{"Order"},
		})
		if appErr != nil {
			return nil, appErr
		}
		if len(orderLines) == 0 {
			return updatedContentUrl, nil
		}

		orderLine := orderLines[0]

		_, appErr = a.srv.AccountService().CommonCustomerCreateEvent(
			nil,
			orderLine.Order.UserID,
			&orderLine.OrderID,
			model.CUSTOMER_EVENT_TYPE_DIGITAL_LINK_DOWNLOADED,
			model.StringInterface{"order_line_pk": orderLine.Id},
		)
		if appErr != nil {
			return nil, appErr
		}

	}

	return updatedContentUrl, nil
}
