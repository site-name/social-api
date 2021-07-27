package channel

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
)

func (a *AppChannel) GetDefaultChannel() (*channel.Channel, *model.AppError) {
	panic("not implemented")
}

func (a *AppChannel) GetDefaultChannelSlugOrGraphqlError() (string, *model.AppError) {
	panic("not implemented")
}
