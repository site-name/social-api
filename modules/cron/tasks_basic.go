package cron

// import (
// 	"context"

// 	"github.com/sitename/sitename/models"
// )

// func registerSyncExternalUsers() {
// 	RegisterTaskFatal("sync_external_users", &UpdateExistingConfig{
// 		BaseConfig: BaseConfig{
// 			Enabled:    true,
// 			RunAtStart: false,
// 			Schedule:   "@every 24h",
// 		},
// 		UpdateExisting: true,
// 	}, func(ctx context.Context, _ *models.User, config Config) error {
// 		realConfig := config.(*UpdateExistingConfig)
// 		return models.SyncExternalUsers(ctx, realConfig.UpdateExisting)
// 	})
// }
