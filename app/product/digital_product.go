package product

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

func (a *ServiceProduct) GetDefaultDigitalContentSettings(aShop *model.Shop) *model.ShopDefaultDigitalContentSettings {
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
func (a *ServiceProduct) DigitalContentUrlIsValid(contentURL *model.DigitalContentUrl) (bool, *model.AppError) {
	digitalContent, appErr := a.DigitalContentbyOption(&model.DigitalContenetFilterOption{
		Id: squirrel.Eq{store.DigitalContentTableName + ".Id": contentURL.ContentID},
	})
	if appErr != nil {
		return false, appErr
	}

	var (
		urlValidDays *uint
		maxDownloads *uint
	)
	if *digitalContent.UseDefaultSettings {
		shop, appErr := a.srv.ShopService().ShopById(digitalContent.ShopID)
		if appErr != nil {
			return false, appErr
		}
		shopDigitalContentSetting := a.GetDefaultDigitalContentSettings(shop)

		urlValidDays = shopDigitalContentSetting.DefaultDigitalUrlValidDays
		maxDownloads = shopDigitalContentSetting.DefaultDigitalMaxDownloads
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

	if maxDownloads != nil && *maxDownloads <= uint(contentURL.DownloadNum) {
		return false, nil
	}

	return true, nil
}

func (a *ServiceProduct) IncrementDownloadCount(contentURL *model.DigitalContentUrl) *model.AppError {
	contentURL.DownloadNum++
	_, appErr := a.UpsertDigitalContentURL(contentURL)
	if appErr != nil {
		return appErr
	}

	if contentURL.LineID != nil {
		orderLine, appErr := a.srv.OrderService().OrderLineById(*contentURL.LineID)
		if appErr != nil {
			return appErr
		}
		userByOrderId, appErr := a.srv.AccountService().UserByOrderId(orderLine.OrderID)
		if appErr != nil {
			return appErr
		}

		if orderLine != nil && userByOrderId != nil {
			_, appErr = a.srv.AccountService().CommonCustomerCreateEvent(
				&userByOrderId.Id,
				&orderLine.OrderID,
				model.DIGITAL_LINK_DOWNLOADED,
				map[string]interface{}{"order_line_pk": orderLine.Id},
			)
			if appErr != nil {
				return appErr
			}
		}
	}

	return nil
}
