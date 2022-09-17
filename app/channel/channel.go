/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package channel

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type ServiceChannel struct {
	srv *app.Server
}

func init() {
	app.RegisterChannelService(func(s *app.Server) (sub_app_iface.ChannelService, error) {
		return &ServiceChannel{
			srv: s,
		}, nil
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

// ValidateChannel check if a channel with given slug is active
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

// CleanChannel
func (a *ServiceChannel) CleanChannel(channelID *string) (*model.Channel, *model.AppError) {
	var (
		needChannel *model.Channel
		appErr      *model.AppError
	)

	if channelID != nil && *channelID != "" {
		needChannel, appErr = a.ValidateChannel(*channelID)
	} else {
		needChannel, appErr = a.ChannelByOption(&model.ChannelFilterOption{
			IsActive: model.NewBool(true),
		})
	}

	if appErr != nil {
		return nil, appErr
	}

	return needChannel, nil
}

// ChannelsByOption returns a list of channels by given options
func (a *ServiceChannel) ChannelsByOption(option *model.ChannelFilterOption) (model.Channels, *model.AppError) {
	channels, err := a.srv.Store.Channel().FilterByOption(option)
	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	}
	if len(channels) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ChannelsByOptions", "app.channel.error_finding_channels_by_options.app_error", nil, errMsg, statusCode)
	}

	return channels, nil
}
