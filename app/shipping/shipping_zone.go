package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// ShippingZonesByOption returns all shipping zones that satisfy given options
func (a *ServiceShipping) ShippingZonesByOption(option *model.ShippingZoneFilterOption) ([]*model.ShippingZone, *model_helper.AppError) {
	shippingZones, err := a.srv.Store.ShippingZone().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("ShippingZonesByOption", "app.shipping.shipping_zones_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return shippingZones, nil
}

func (s *ServiceShipping) UpsertShippingZone(transaction boil.ContextTransactor, zone *model.ShippingZone) (*model.ShippingZone, *model_helper.AppError) {
	zone, err := s.srv.Store.ShippingZone().Upsert(transaction, zone)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("UpsertShippingZone", "app.shipping.upsert_shipping_zone.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return zone, nil
}

func (s *ServiceShipping) DeleteShippingZones(transaction boil.ContextTransactor, conditions *model.ShippingZoneFilterOption) (int64, *model_helper.AppError) {
	numDeleted, err := s.srv.Store.ShippingZone().Delete(transaction, conditions)
	if err != nil {
		return 0, model_helper.NewAppError("DeleteShippingZones", "app.shipping.delete_shipping_zones.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return numDeleted, nil
}

func (s *ServiceShipping) ToggleShippingZoneRelations(transaction boil.ContextTransactor, zones model.ShippingZones, warehouseIds, channelIds []string, delete bool) *model_helper.AppError {
	err := s.srv.Store.ShippingZone().ToggleRelations(transaction, zones, warehouseIds, channelIds, delete)
	if err != nil {
		if _, ok := err.(*store.ErrInvalidInput); ok {
			return model_helper.NewAppError("ToggleShippingZoneRelations", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "relations"}, err.Error(), http.StatusBadRequest)
		}
		return model_helper.NewAppError("ToggleShippingZoneRelations", "app.channel.add_channel_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
