package app

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/store"
)

// GetChannelBySlug get a channel from database with given slug
func (a *App) GetChannelBySlug(slug string) (*channel.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().GetBySlug(slug)
	if err != nil {
		var nfErr *store.ErrNotFound
		if errors.As(err, &nfErr) {
			return nil, model.NewAppError("GetChannelBySlug", "app.channel.missing_channel.app_error", nil, err.Error(), http.StatusNotFound)
		}
		return nil, model.NewAppError("GetChannelBySlug", "app.channel.missing_channel.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return channel, nil
}

// GetDefaultChannel get random channel that is active
func (a *App) GetDefaultActiveChannel() (*channel.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().GetRandomActiveChannel()
	if err != nil {
		var nfErr *store.ErrNotFound
		if errors.As(err, &nfErr) {
			return nil, model.NewAppError("GetDefaultChannel", "app.channel.missing_channel.app_error", nil, err.Error(), http.StatusNotFound)
		}
		return nil, model.NewAppError("GetDefaultChannel", "app.channel.missing_channel.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return channel, nil
}

// CleanChannel performs:
//
// 1) If given slug is not nil, try getting a channel with that slug.
//   +) if found, check if channel is active
//
// 2) If given slug if nil, it try
func (a *App) CleanChannel(channelSlug *string) (*channel.Channel, *model.AppError) {
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
