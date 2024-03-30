/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package channel

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ServiceChannel struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Channel = &ServiceChannel{s}
		return nil
	})
}

func (s *ServiceChannel) ChannelByOption(option model_helper.ChannelFilterOptions) (*model.Channel, *model_helper.AppError) {
	channels, appErr := s.ChannelsByOption(option)
	if appErr != nil {
		return nil, appErr
	}
	if len(channels) == 0 {
		return nil, model_helper.NewAppError("ChannelByOption", "app.channel.channel_by_options.app_error", nil, "no channel exist", http.StatusNotFound)
	}

	return channels[0], nil
}

func (a *ServiceChannel) ValidateChannel(channelID string) (*model.Channel, *model_helper.AppError) {
	channel, appErr := a.ChannelByOption(model_helper.ChannelFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(model.ChannelWhere.ID.EQ(channelID)),
	})
	if appErr != nil {
		return nil, appErr
	}

	if !channel.IsActive {
		return nil, model_helper.NewAppError("CleanChannel", "app.channel.channel_inactive.app_error", nil, "", http.StatusNotModified)
	}

	return channel, nil
}

func (a *ServiceChannel) CleanChannel(channelID *string) (*model.Channel, *model_helper.AppError) {
	if channelID != nil {
		return a.ValidateChannel(*channelID)
	}
	return a.ChannelByOption(model_helper.ChannelFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(model.ChannelWhere.IsActive.EQ(true)),
	})
}

func (a *ServiceChannel) ChannelsByOption(option model_helper.ChannelFilterOptions) (model.ChannelSlice, *model_helper.AppError) {
	channels, err := a.srv.Store.Channel().FilterByOptions(option)
	if err != nil {
		return nil, model_helper.NewAppError("ChannelsByOptions", "app.channel.error_finding_channels_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return channels, nil
}

func (a *ServiceChannel) UpsertChannel(tx boil.ContextTransactor, channel model.Channel) (*model.Channel, *model_helper.AppError) {
	savedChannel, err := a.srv.Store.Channel().Upsert(tx, channel)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("UpsertChannel", "app.channel.upsert_channel.app_error", nil, err.Error(), statusCode)
	}
	return savedChannel, nil
}

func (s *ServiceChannel) DeleteChannels(transaction boil.ContextTransactor, ids []string) *model_helper.AppError {
	err := s.srv.Store.Channel().DeleteChannels(transaction, ids)
	if err != nil {
		return model_helper.NewAppError("DeleteChannels", "app.channel.channel_delete_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
