package app

import (
	"context"
	"reflect"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
)

const EmojisPermissionsMigrationKey = "EmojisPermissionsMigrationComplete"
const GuestRolesCreationMigrationKey = "GuestRolesCreationMigrationComplete"
const SystemConsoleRolesCreationMigrationKey = "SystemConsoleRolesCreationMigrationComplete"
const ContentExtractionConfigDefaultTrueMigrationKey = "ContentExtractionConfigDefaultTrueMigrationComplete"

// This function migrates the default built in roles from code/config to the database.
func (a *App) DoAdvancedPermissionsMigration() {
	a.Srv().doAdvancedPermissionsMigration()
}

func (s *Server) doAdvancedPermissionsMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := s.Store.System().GetByName(model_helper.AdvancedPermissionsMigrationKey); err == nil {
		return
	}

	slog.Info("Migrating roles to database.")
	roles := model_helper.MakeDefaultRoles()

	allSucceeded := true
	for _, rawRole := range roles {
		role, err := s.Store.Role().Upsert(*rawRole)
		if err == nil {
			continue
		}

		// If this failed for reasons other than the role already existing, don't mark the migration as done.
		fetchedRole, err := s.Store.Role().GetByName(context.Background(), role.Name)
		if err != nil {
			slog.Critical("Failed to migrate role to database.", slog.Err(err))
			allSucceeded = false
			continue
		}

		// If the role already existed, check it is the same and update if not.
		if !reflect.DeepEqual(fetchedRole.Permissions, role.Permissions) ||
			fetchedRole.DisplayName != role.DisplayName ||
			fetchedRole.Description != role.Description ||
			fetchedRole.SchemeManaged != role.SchemeManaged {
			role.ID = fetchedRole.ID
			if _, err = s.Store.Role().Upsert(*role); err != nil {
				// Role is not the same, but failed to update.
				slog.Critical("Failed to migrate role to database.", slog.Err(err))
				allSucceeded = false
			}
		}
	}

	if !allSucceeded {
		return
	}

	config := s.Config()
	*config.ServiceSettings.PostEditTimeLimit = -1
	if _, _, err := s.SaveConfig(config, true); err != nil {
		slog.Error("failed to update config in Advanced Permissions Phase 1 Migration.", slog.Err(err))
	}

	// create a evidence to prove that all roles have been saved to database.
	system := model.System{
		Name:  model_helper.AdvancedPermissionsMigrationKey,
		Value: "true",
	}

	if err := s.Store.System().Save(system); err != nil {
		slog.Critical("Failed to mark advanced permissions migration as completed.", slog.Err(err))
	}
}

func (a *App) SetPhase2PermissionsMigrationStatus(isComplete bool) error {
	if !isComplete {
		if _, err := a.Srv().Store.System().PermanentDeleteByName(model_helper.MigrationKeyAdvancedPermissionsPhase2); err != nil {
			return err
		}
	}
	a.Srv().phase2PermissionsMigrationComplete = isComplete
	return nil
}

func (a *App) DoSystemConsoleRolesCreationMigration() {
	a.Srv().doSystemConsoleRolesCreationMigration()
}

func (s *Server) doSystemConsoleRolesCreationMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := s.Store.System().GetByName(SystemConsoleRolesCreationMigrationKey); err == nil {
		return
	}

	roles := model_helper.MakeDefaultRoles()

	allSucceeded := true
	if _, err := s.Store.Role().GetByName(context.Background(), model_helper.SystemManagerRoleId); err != nil {
		if _, err := s.Store.Role().Upsert(*roles[model_helper.SystemManagerRoleId]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model_helper.SystemManagerRoleId))
			allSucceeded = false
		}
	}
	if _, err := s.Store.Role().GetByName(context.Background(), model_helper.SystemReadOnlyAdminRoleId); err != nil {
		if _, err := s.Store.Role().Upsert(*roles[model_helper.SystemReadOnlyAdminRoleId]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model_helper.SystemReadOnlyAdminRoleId))
			allSucceeded = false
		}
	}
	if _, err := s.Store.Role().GetByName(context.Background(), model_helper.SystemUserManagerRoleId); err != nil {
		if _, err := s.Store.Role().Upsert(*roles[model_helper.SystemUserManagerRoleId]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model_helper.SystemUserManagerRoleId))
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

	if err := s.Store.System().Save(system); err != nil {
		slog.Critical("Failed to mark system console roles creation migration as completed.", slog.Err(err))
	}
}

func (s *Server) doContentExtractionConfigDefaultTrueMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := s.Store.System().GetByName(ContentExtractionConfigDefaultTrueMigrationKey); err == nil {
		return
	}

	s.UpdateConfig(func(config *model_helper.Config) {
		config.FileSettings.ExtractContent = model_helper.GetPointerOfValue(true)
	})

	system := model.System{
		Name:  ContentExtractionConfigDefaultTrueMigrationKey,
		Value: "true",
	}

	if err := s.Store.System().Save(system); err != nil {
		slog.Critical("Failed to mark content extraction config migration as completed.", slog.Err(err))
	}
}

// DoAppMigrations migrate permissions
func (a *App) DoAppMigrations() {
	a.Srv().doAppMigrations()
}

func (s *Server) doAppMigrations() {
	s.doAdvancedPermissionsMigration()
	s.doSystemConsoleRolesCreationMigration()
	s.doContentExtractionConfigDefaultTrueMigration()
}
