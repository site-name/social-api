package cron

import (
	"context"
	"time"

	"github.com/sitename/sitename/models"
)

func registerDeleteInactiveUsers() {
	RegisterTaskFatal(
		"delete_inactive_accounts",
		&OlderThanConfig{
			BaseConfig: BaseConfig{
				Enabled:    false,
				RunAtStart: false,
				Schedule:   "@annually",
			},
			OlderThan: 0 * time.Second,
		},
		func(ctx context.Context, _ *models.User, config Config) error {
			olderThanConfig := config.(*OlderThanConfig)
			return models.DeleteInactiveUsers(ctx, olderThanConfig.OlderThan)
		},
	)
}

func initExtendedTasks() {
	registerDeleteInactiveUsers()
}
