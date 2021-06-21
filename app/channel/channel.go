package channel

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/store"
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
		var nfErr *store.ErrNotFound
		statusCode := http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetChannelBySlug", "app.channel.missing_channel.app_error", nil, err.Error(), statusCode)
	}

	return channel, nil
}

// GetDefaultChannel get random channel that is active
func (a *AppChannel) GetDefaultActiveChannel() (*channel.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().GetRandomActiveChannel()
	if err != nil {
		var nfErr *store.ErrNotFound
		statusCode := http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetDefaultChannel", "app.channel.missing_channel.app_error", nil, err.Error(), statusCode)
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
			return nil, model.NewAppError("CleanChannel", "app.channel.channel_inactive.app_error", nil, "", 0)
		}
		return channel, nil
	}

	channel, err := a.GetDefaultActiveChannel()
	if err != nil {
		return nil, err
	}
	if !channel.IsActive {
		return nil, model.NewAppError("CleanChannel", "app.channel.channel_inactive.app_error", nil, "", 0)
	}
	return channel, nil
}
