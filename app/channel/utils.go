package channel

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
)

func (a *ServiceChannel) GetDefaultChannel() (*model.Channel, *model.AppError) {
	channel, appErr := a.ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Expr(model.ChannelTableName + ".IsActive"),
		Limit:      1,
	})
	if appErr != nil {
		return nil, appErr
	}

	return channel, nil
}

func (a *ServiceChannel) GetDefaultChannelSlugOrGraphqlError() (string, *model.AppError) {
	panic("not implemented")
}
