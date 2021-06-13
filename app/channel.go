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
func (a *App) GetDefaultChannel() (*channel.Channel, *model.AppError) {
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

	channel, err := a.GetDefaultChannel()
	if err != nil {
		return nil, err
	}
	if !channel.IsActive {
		return nil, model.NewAppError("CleanChannel", "app.channel.channel_inactive.app_error", nil, "", 0)
	}
	return channel, nil
}
