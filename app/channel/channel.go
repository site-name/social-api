/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package channel

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
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

// ChannelByOption returns a channel that satisfies given options
func (s *ServiceChannel) ChannelByOption(option *model.ChannelFilterOption) (*model.Channel, *model_helper.AppError) {
	option.Limit = 1
	channels, appErr := s.ChannelsByOption(option)
	if appErr != nil {
		return nil, appErr
	}
	if channels.Len() == 0 {
		return nil, model_helper.NewAppError("ChannelByOption", "app.channel.channel_by_options.app_error", nil, "no channel exist", http.StatusNotFound)
	}

	return channels[0], nil
}

// ValidateChannel check if a channel with given id is active
func (a *ServiceChannel) ValidateChannel(channelID string) (*model.Channel, *model_helper.AppError) {
	channel, appErr := a.ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": channelID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// check if channel is active
	if !channel.IsActive {
		return nil, model_helper.NewAppError("CleanChannel", "app.channel.channel_inactive.app_error", nil, "", http.StatusNotModified)
	}

	return channel, nil
}

func (a *ServiceChannel) CleanChannel(channelID *string) (*model.Channel, *model_helper.AppError) {
	if channelID != nil {
		return a.ValidateChannel(*channelID)
	}
	return a.ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Expr(model.ChannelTableName + ".IsActive"),
	})
}

// ChannelsByOption returns a list of channels by given options
func (a *ServiceChannel) ChannelsByOption(option *model.ChannelFilterOption) (model.Channels, *model_helper.AppError) {
	channels, err := a.srv.Store.Channel().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("ChannelsByOptions", "app.channel.error_finding_channels_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return channels, nil
}

func (a *ServiceChannel) UpsertChannel(transaction *gorm.DB, channel *model.Channel) (*model.Channel, *model_helper.AppError) {
	channel, err := a.srv.Store.Channel().Upsert(transaction, channel)
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
	return channel, nil
}

func (s *ServiceChannel) DeleteChannels(transaction *gorm.DB, ids ...string) *model_helper.AppError {
	err := s.srv.Store.Channel().DeleteChannels(transaction, ids)
	if err != nil {
		return model_helper.NewAppError("DeleteChannels", "app.channel.channel_delete_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
