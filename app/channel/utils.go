package channel

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
)

func (a *ServiceChannel) GetDefaultChannel() (*channel.Channel, *model.AppError) {
	panic("not implemented")
}

func (a *ServiceChannel) GetDefaultChannelSlugOrGraphqlError() (string, *model.AppError) {
	panic("not implemented")
}
