package app

import (
	"context"
	"reflect"

	"github.com/sitename/sitename/model"
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
	if _, err := s.Store.System().GetByName(model.AdvancedPermissionsMigrationKey); err == nil {
		return
	}

	slog.Info("Migrating roles to database.")
	roles := model.MakeDefaultRoles()

	allSucceeded := true
	for _, role := range roles {
		_, err := s.Store.Role().Save(role)
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
			role.Id = fetchedRole.Id
			if _, err = s.Store.Role().Save(role); err != nil {
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
		Name:  model.AdvancedPermissionsMigrationKey,
		Value: "true",
	}

	if err := s.Store.System().Save(&system); err != nil {
		slog.Critical("Failed to mark advanced permissions migration as completed.", slog.Err(err))
	}
}

func (a *App) SetPhase2PermissionsMigrationStatus(isComplete bool) error {
	if !isComplete {
		if _, err := a.Srv().Store.System().PermanentDeleteByName(model.MigrationKeyAdvancedPermissionsPhase2); err != nil {
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

	roles := model.MakeDefaultRoles()

	allSucceeded := true
	if _, err := s.Store.Role().GetByName(context.Background(), model.SystemManagerRoleId); err != nil {
		if _, err := s.Store.Role().Save(roles[model.SystemManagerRoleId]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model.SystemManagerRoleId))
			allSucceeded = false
		}
	}
	if _, err := s.Store.Role().GetByName(context.Background(), model.SystemReadOnlyAdminRoleId); err != nil {
		if _, err := s.Store.Role().Save(roles[model.SystemReadOnlyAdminRoleId]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model.SystemReadOnlyAdminRoleId))
			allSucceeded = false
		}
	}
	if _, err := s.Store.Role().GetByName(context.Background(), model.SystemUserManagerRoleId); err != nil {
		if _, err := s.Store.Role().Save(roles[model.SystemUserManagerRoleId]); err != nil {
			slog.Critical("Failed to create new role.", slog.Err(err), slog.String("role", model.SystemUserManagerRoleId))
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

	if err := s.Store.System().Save(&system); err != nil {
		slog.Critical("Failed to mark system console roles creation migration as completed.", slog.Err(err))
	}
}

func (s *Server) doContentExtractionConfigDefaultTrueMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := s.Store.System().GetByName(ContentExtractionConfigDefaultTrueMigrationKey); err == nil {
		return
	}

	s.UpdateConfig(func(config *model.Config) {
		config.FileSettings.ExtractContent = model.NewBool(true)
	})

	system := model.System{
		Name:  ContentExtractionConfigDefaultTrueMigrationKey,
		Value: "true",
	}

	if err := s.Store.System().Save(&system); err != nil {
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
	// This migration always must be the last, because can be based on previous
	// migrations. For example, it needs the guest roles migration.
	err := s.doPermissionsMigrations()
	if err != nil {
		slog.Critical("(app.App).DoPermissionsMigrations failed", slog.Err(err))
	}
	s.doContentExtractionConfigDefaultTrueMigration()
}

func SetRolePermissionsFromConfig(roles map[string]*model.Role, cfg *model.Config) map[string]*model.Role {
	if !*cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations {
		roles[model.SystemUserRoleId].Permissions = append(
			roles[model.SystemUserRoleId].Permissions,
			model.PermissionManageOAuth.Id,
		)
	}

	switch *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_AllowEditPost {
	case model.ALLOW_EDIT_POST_ALWAYS, model.ALLOW_EDIT_POST_TIME_LIMIT:
		roles[model.SystemAdminRoleId].Permissions = append(
			roles[model.SystemAdminRoleId].Permissions,
			model.PermissionEditPost.Id,
		)
	}

	return roles
}
