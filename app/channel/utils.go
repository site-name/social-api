package channel

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *ServiceChannel) GetDefaultChannel() (*model.Channel, *model_helper.AppError) {
	channel, appErr := a.ChannelByOption(model_helper.ChannelFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.ChannelWhere.IsActive.EQ(true),
			qm.Limit(1),
		),
	})
	if appErr != nil {
		return nil, appErr
	}

	return channel, nil
}

func (a *ServiceChannel) GetDefaultChannelSlugOrGraphqlError() (string, *model_helper.AppError) {
	channel, appErr := a.GetDefaultChannel()
	if appErr != nil {
		return "", appErr
	}
	return channel.Slug, nil
}
