package channel

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (a *ServiceChannel) ChannelShopRelationsByOptions(options *model.ChannelShopRelationFilterOptions) ([]*model.ChannelShopRelation, *model.AppError) {
	relations, err := a.srv.Store.ChannelShop().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("ChannelShopRelationsByOptions", "app.channel.channel_shops_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return relations, nil
}

func (a *ServiceChannel) SaveChannelShopRelation(relation *model.ChannelShopRelation) (*model.ChannelShopRelation, *model.AppError) {
	relation, err := a.srv.Store.ChannelShop().Save(relation)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("SaveChannelShopRelation", "app.channel.save_channel_shop.app_error", nil, err.Error(), statusCode)
	}

	return relation, nil
}

func (s *ServiceChannel) ShopSellsInChannel(shopID, channelID string) bool {
	relations, appErr := s.ChannelShopRelationsByOptions(&model.ChannelShopRelationFilterOptions{
		ShopID:    squirrel.Eq{store.ChannelShopRelationTableName + ".ShopID": shopID},
		ChannelID: squirrel.Eq{store.ChannelShopRelationTableName + ".ChannelID": channelID},
	})
	if appErr != nil {
		slog.Error("failed to check if shop sells in channel", slog.Err(appErr))
		return false
	}

	return len(relations) > 0
}
