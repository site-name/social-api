package einterfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type AccountMigrationInterface interface {
	MigrateToLdap(fromAuthService string, forignUserFieldNameToMatch string, force bool, dryRun bool) *model_helper.AppError
	MigrateToSaml(fromAuthService string, usersMap map[string]string, auto bool, dryRun bool) *model_helper.AppError
}
