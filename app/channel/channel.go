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

type AppChannel struct {
	app.AppIface
}

func init() {
	app.RegisterChannelApp(func(a app.AppIface) sub_app_iface.ChannelApp {
		return &AppChannel{a}
	})
}

func (a *AppChannel) GetChannelBySlug(slug string) (*channel.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetChannelBySlug", channelMissingErrorId, err)
	}

	return channel, nil
}

func (a *AppChannel) GetDefaultActiveChannel() (*channel.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().GetRandomActiveChannel()
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetDefaultActiveChannel", channelMissingErrorId, err)
	}

	return channel, nil
}

func (a *AppChannel) ValidateChannel(channelSlug string) (*channel.Channel, *model.AppError) {
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

func (a *AppChannel) CleanChannel(channelSlug *string) (*channel.Channel, *model.AppError) {
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
func (a *AppChannel) ChannelsByOption(option *channel.ChannelFilterOption) ([]*channel.Channel, *model.AppError) {
	channels, err := a.Srv().Store.Channel().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ChannelsByOption", "app.channel.error_finding_channels_by_option.app_error", err)
	}

	return channels, nil
}
