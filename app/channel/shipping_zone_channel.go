package channel

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

func (s *ServiceChannel) BulkUpsertShippingZoneChannels(transaction store_iface.SqlxExecutor, relations []*model.ShippingZoneChannel) ([]*model.ShippingZoneChannel, *model.AppError) {
	savedRelations, err := s.srv.Store.ShippingZoneChannel().BulkSave(transaction, relations)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("BulkUpsertShippingZoneChannels", "app.channel.bulk_saving_shipping_zone_channels.app_error", nil, err.Error(), statusCode)
	}

	return savedRelations, nil
}

func (s *ServiceChannel) BulkDeleteShippingZoneChannels(transaction store_iface.SqlxExecutor, options *model.ShippingZoneChannelFilterOptions) *model.AppError {
	err := s.srv.Store.ShippingZoneChannel().BulkDelete(transaction, options)
	if err != nil {
		return model.NewAppError("BulkDeleteShippingZoneChannels", "app.channel.bulk_delete_shipping_zone_channels.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
