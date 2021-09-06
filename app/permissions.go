package app

import (
	"io"
	"net/http"

	"github.com/sitename/sitename/model"
)

const permissionsExportBatchSize = 100
const systemSchemeName = "00000000-0000-0000-0000-000000000000" // Prevents collisions with user-created schemes.

// ResetPermissionsSystem reset permission system
func (a *App) ResetPermissionsSystem() *model.AppError {

	// Reset all Custom Role assignments to Users.
	if err := a.Srv().Store.User().ClearAllCustomRoleAssignments(); err != nil {
		return model.NewAppError("ResetPermissionsSystem", "app.user.clear_all_custom_role_assignments.select.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Purge all roles from the database.
	if err := a.Srv().Store.Role().PermanentDeleteAll(); err != nil {
		return model.NewAppError("ResetPermissionsSystem", "app.role.permanent_delete_all.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Remove the "System" table entry that marks the advanced permissions migration as done.
	if _, err := a.Srv().Store.System().PermanentDeleteByName(model.AdvancedPermissionsMigrationKey); err != nil {
		return model.NewAppError("ResetPermissionSystem", "app.system.permanent_delete_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Remove the "System" table entry that marks the emoji permissions migration as done.
	if _, err := a.Srv().Store.System().PermanentDeleteByName(EmojisPermissionsMigrationKey); err != nil {
		return model.NewAppError("ResetPermissionSystem", "app.system.permanent_delete_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Remove the "System" table entry that marks the guest roles permissions migration as done.
	if _, err := a.Srv().Store.System().PermanentDeleteByName(GuestRolesCreationMigrationKey); err != nil {
		return model.NewAppError("ResetPermissionSystem", "app.system.permanent_delete_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Now that the permissions system has been reset, re-run the migration to reinitialise it.
	a.DoAppMigrations()

	return nil
}

func (a *App) ExportPermissions(w io.Writer) error {
	panic("not implemented")
}
