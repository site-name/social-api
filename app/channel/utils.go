package channel

import (
	"github.com/sitename/sitename/model"
)

func (a *ServiceChannel) GetDefaultChannel() (*model.Channel, *model.AppError) {
	panic("not implemented")
}

func (a *ServiceChannel) GetDefaultChannelSlugOrGraphqlError() (string, *model.AppError) {
	panic("not implemented")
}
