package app

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

const EmojisPermissionsMigrationKey = "EmojisPermissionsMigrationComplete"
const GuestRolesCreationMigrationKey = "GuestRolesCreationMigrationComplete"
const SystemConsoleRolesCreationMigrationKey = "SystemConsoleRolesCreationMigrationComplete"
const ContentExtractionConfigDefaultTrueMigrationKey = "ContentExtractionConfigDefaultTrueMigrationComplete"

// This function migrates the default built in roles from code/config to the database.
func (a *App) DoAdvancedPermissionsMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := a.Srv().Store.System().GetByName(model.ADVANCED_PERMISSIONS_MIGRATION_KEY); err == nil {
		return
	}

	slog.Info("Migrating roles to database.")
	roles := model.MakeDefaultRoles()
	fmt.Println("roles:-------------", roles)

	roles = util.SetRolePermissionsFromConfig(roles, a.Config(), true)

	allSucceeded := true

	for _, role := range roles {
		_, err := a.Srv().Store.Role().Save(role)
		if err == nil {
			continue
		}

		// If this failed for reasons other than the role already existing, don't mark the migration as done.
		fetchedRole, err := a.Srv().Store.Role().GetByName(context.Background(), role.Name)
		if err != nil {
			slog.Critical("Failed to migrate role to database.", slog.Err(err))
			allSucceeded = false
			continue
		}

		// if the role already existed, check it is the same and update if not.
		if !reflect.DeepEqual(fetchedRole.Permissions, role.Permissions) ||
			fetchedRole.DisplayName != role.DisplayName ||
			fetchedRole.Description != role.Description ||
			fetchedRole.SchemeManaged != role.SchemeManaged {
			role.Id = fetchedRole.Id
			if _, err := a.Srv().Store.Role().Save(role); err != nil {
				// Role is not the same but failed to update
				slog.Critical("Failed to migrate role to database.", slog.Err(err))
				allSucceeded = false
			}
		}
	}

	if !allSucceeded {
		return
	}

	config := a.Config()
	if *config.ServiceSettings.DEPRECATED_DO_NOT_USE_AllowEditPost == model.ALLOW_EDIT_POST_ALWAYS {
		*config.ServiceSettings.PostEditTimeLimit = -1
		if err := a.SaveConfig(config, true); err != nil {
			slog.Error("Failed to update config in Advanced Permissions Phase 1 Migration.", slog.Err(err))
		}
	}

	system := model.System{
		Name:  model.ADVANCED_PERMISSIONS_MIGRATION_KEY,
		Value: "true",
	}

	if err := a.Srv().Store.System().Save(&system); err != nil {
		slog.Critical("Failed to mark advanced permissions migration as completed.", slog.Err(err))
	}
}

// func (a *App) DoGuestRolesCreationMigration() {
// 	// If the migration is already marked as completed, don't do it again.
// 	if _, err := a.Srv().Store.System().GetByName(GuestRolesCreationMigrationKey); err == nil {
// 		return
// 	}

// 	roles := model.MakeDefaultRoles()

// 	allSucceeded := true
// 	if _, err := a.Srv().Store.Role().GetByName(context.Background(), model.SYSTEM_GUEST_ROLE_ID); err != nil {
// 		if _, err := a.Srv().Store.Role().Save(roles[model.SYSTEM_GUEST_ROLE_ID]); err != nil {
// 			slog.Critical("Failed to create new guest tole to database.", slog.Err(err))
// 			allSucceeded = false
// 		}
// 	}

// 	schemes, err := a.Srv().Store.Scheme()
// }

func (a *App) DoSystemConsoleRolesCreationMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := a.Srv().Store.System().GetByName(SystemConsoleRolesCreationMigrationKey); err == nil {
		return
	}

	roles := model.MakeDefaultRoles()

	allSucceeded := true
	if _, err := a.Srv().Store.Role().GetByName(context.Background(), model.SYSTEM_MANAGER_ROLE_ID); err != nil {
		if _, err := a.Srv().Store.Role().Save(roles[model.SYSTEM_MANAGER_ROLE_ID]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model.SYSTEM_MANAGER_ROLE_ID))
			allSucceeded = false
		}
	}
	if _, err := a.Srv().Store.Role().GetByName(context.Background(), model.SYSTEM_READ_ONLY_ADMIN_ROLE_ID); err != nil {
		if _, err := a.Srv().Store.Role().Save(roles[model.SYSTEM_READ_ONLY_ADMIN_ROLE_ID]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model.SYSTEM_READ_ONLY_ADMIN_ROLE_ID))
			allSucceeded = false
		}
	}
	if _, err := a.Srv().Store.Role().GetByName(context.Background(), model.SYSTEM_USER_MANAGER_ROLE_ID); err != nil {
		if _, err := a.Srv().Store.Role().Save(roles[model.SYSTEM_USER_MANAGER_ROLE_ID]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model.SYSTEM_USER_MANAGER_ROLE_ID))
			allSucceeded = false
		}
	}

	if !allSucceeded {
		return
	}

	system := model.System{
		Name:  SystemConsoleRolesCreationMigrationKey,
		Value: "true",
	}

	if err := a.Srv().Store.System().Save(&system); err != nil {
		slog.Critical("Failed to mark system console roles creation migration as completed.", slog.Err(err))
	}
}

func (a *App) doContentExtractionConfigDefaultTrueMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := a.Srv().Store.System().GetByName(ContentExtractionConfigDefaultTrueMigrationKey); err == nil {
		return
	}

	a.UpdateConfig(func(config *model.Config) {
		config.FileSettings.ExtractContent = model.NewBool(true)
	})

	system := model.System{
		Name:  ContentExtractionConfigDefaultTrueMigrationKey,
		Value: "true",
	}

	if err := a.Srv().Store.System().Save(&system); err != nil {
		slog.Critical("Failed to mark content extraction config migration as completed.", slog.Err(err))
	}
}
func (a *App) DoAppMigrations() {
	a.DoAdvancedPermissionsMigration()
	// a.DoEmojisPermissionsMigration() // NOTE: need investigating
	// a.DoGuestRolesCreationMigration()
	a.DoSystemConsoleRolesCreationMigration()
	// This migration always must be the last, because can be based on previous
	// migrations. For example, it needs the guest roles migration.
	err := a.DoPermissionsMigrations()
	if err != nil {
		slog.Critical("(app.App).DoPermissionsMigrations failed", slog.Err(err))
	}
	a.doContentExtractionConfigDefaultTrueMigration()
}
