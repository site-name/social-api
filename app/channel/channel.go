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

// GetChannelBySlug get a channel from database with given slug
func (a *AppChannel) GetChannelBySlug(slug string) (*channel.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetChannelBySlug", channelMissingErrorId, err)
	}

	return channel, nil
}

// GetDefaultChannel get random channel that is active
func (a *AppChannel) GetDefaultActiveChannel() (*channel.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().GetRandomActiveChannel()
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetDefaultActiveChannel", channelMissingErrorId, err)
	}

	return channel, nil
}

// CleanChannel performs:
//
// 1) If given slug is not nil, try getting a channel with that slug.
//   +) if found, check if channel is active
//
// 2) If given slug if nil, it try
func (a *AppChannel) CleanChannel(channelSlug *string) (*channel.Channel, *model.AppError) {
	var channel *channel.Channel

	if channelSlug != nil {
		channel, err := a.GetChannelBySlug(*channelSlug)
		if err != nil {
			return nil, err
		}
		if !channel.IsActive {
			return nil, model.NewAppError("CleanChannel", channelInactiveErrorId, nil, "", http.StatusNotModified)
		}
		return channel, nil
	}

	channel, err := a.GetDefaultActiveChannel()
	if err != nil {
		return nil, err
	}
	if !channel.IsActive {
		return nil, model.NewAppError("CleanChannel", channelInactiveErrorId, nil, "", http.StatusNotModified)
	}
	return channel, nil
}
