/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package channel

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/store"
)

const (
	channelMissingErrorId  = "app.channel.missing_channel.app_error"
	channelInactiveErrorId = "app.channel.channel_inactive.app_error"
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

func (a *ServiceChannel) GetChannelBySlug(slug string) (*channel.Channel, *model.AppError) {
	channel, err := a.srv.Store.Channel().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetChannelBySlug", channelMissingErrorId, err)
	}

	return channel, nil
}

func (a *ServiceChannel) GetDefaultActiveChannel() (*channel.Channel, *model.AppError) {
	channel, err := a.srv.Store.Channel().GetRandomActiveChannel()
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetDefaultActiveChannel", channelMissingErrorId, err)
	}

	return channel, nil
}

func (a *ServiceChannel) ValidateChannel(channelSlug string) (*channel.Channel, *model.AppError) {
	channel, appErr := a.GetChannelBySlug(channelSlug)
	if appErr != nil {
		return nil, appErr
	}

	// check if channel is active
	if !channel.IsActive {
		return nil, model.NewAppError("CleanChannel", channelInactiveErrorId, nil, "", http.StatusNotModified)
	}

	return channel, nil
}

func (a *ServiceChannel) CleanChannel(channelSlug *string) (*channel.Channel, *model.AppError) {
	var (
		channel *channel.Channel
		appErr  *model.AppError
	)

	if channelSlug != nil && *channelSlug != "" {
		channel, appErr = a.ValidateChannel(*channelSlug)
	} else {
		channel, appErr = a.GetDefaultActiveChannel()
	}

	if appErr != nil {
		return nil, appErr
	}

	return channel, nil
}

// ChannelsByOption returns a list of channels by given options
func (a *ServiceChannel) ChannelsByOption(option *channel.ChannelFilterOption) ([]*channel.Channel, *model.AppError) {
	channels, err := a.srv.Store.Channel().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ChannelsByOption", "app.channel.error_finding_channels_by_option.app_error", err)
	}

	return channels, nil
}
