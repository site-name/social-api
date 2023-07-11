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
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
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
func (s *ServiceChannel) ChannelByOption(option *model.ChannelFilterOption) (*model.Channel, *model.AppError) {
	foundChannel, err := s.srv.Store.Channel().GetbyOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ChannelByOption", "app.channel.error_finding_channel_by_options.app_error", nil, err.Error(), statusCode)
	}

	return foundChannel, nil
}

// ValidateChannel check if a channel with given id is active
func (a *ServiceChannel) ValidateChannel(channelID string) (*model.Channel, *model.AppError) {
	channel, appErr := a.ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": channelID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// check if channel is active
	if !channel.IsActive {
		return nil, model.NewAppError("CleanChannel", "app.channel.channel_inactive.app_error", nil, "", http.StatusNotModified)
	}

	return channel, nil
}

func (a *ServiceChannel) CleanChannel(channelID *string) (*model.Channel, *model.AppError) {
	var (
		channel *model.Channel
		appErr  *model.AppError
	)

	if channelID != nil {
		channel, appErr = a.ValidateChannel(*channelID)
	} else {
		channel, appErr = a.ChannelByOption(&model.ChannelFilterOption{
			IsActive: squirrel.Eq{store.ChannelTableName + ".IsActive": true},
		})
	}
	if appErr != nil {
		return nil, appErr
	}

	return channel, nil
}

// ChannelsByOption returns a list of channels by given options
func (a *ServiceChannel) ChannelsByOption(option *model.ChannelFilterOption) (model.Channels, *model.AppError) {
	channels, err := a.srv.Store.Channel().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ChannelsByOptions", "app.channel.error_finding_channels_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return channels, nil
}

func (a *ServiceChannel) UpsertChannel(transaction store_iface.SqlxTxExecutor, channel *model.Channel) (*model.Channel, *model.AppError) {
	channel, err := a.srv.Store.Channel().Upsert(transaction, channel)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("UpsertChannel", "app.channel.upsert_channel.app_error", nil, err.Error(), statusCode)
	}
	return channel, nil
}

func (s *ServiceChannel) DeleteChannels(transaction store_iface.SqlxTxExecutor, ids ...string) *model.AppError {
	err := s.srv.Store.Channel().DeleteChannels(transaction, ids)
	if err != nil {
		return model.NewAppError("DeleteChannels", "app.channel.channel_delete_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
