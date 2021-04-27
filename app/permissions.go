package app

import (
	"io"
	"net/http"

	"github.com/sitename/sitename/model"
	// "github.com/sitename/sitename/modules/json"
)

const permissionsExportBatchSize = 100
const systemSchemeName = "00000000-0000-0000-0000-000000000000" // Prevents collisions with user-created schemes.

func (a *App) ResetPermissionsSystem() *model.AppError {
	// Reset all Custom Role assignments to Users.
	if err := a.Srv().Store.User().ClearAllCustomRoleAssignments(); err != nil {
		return model.NewAppError("ResetPermissionsSystem", "app.user.clear_all_custom_role_assignments.select.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Remove the "System" table entry that marks the advanced permissions migration as done.
	if _, err := a.Srv().Store.System().PermanentDeleteByName(model.ADVANCED_PERMISSIONS_MIGRATION_KEY); err != nil {
		return model.NewAppError("ResetPermissionSystem", "app.system.permanent_delete_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Now that the permissions system has been reset, re-run the migration to reinitialise it.
	a.DoAppMigrations()
	return nil
}

func (a *App) ExportPermissions(w io.Writer) error {

	panic("not implemented")
}

// func (a *App) ImportPermissions(jsonl io.Reader) error {
// 	createdSchemeIDs := []string{}

// 	scanner := bufio.NewScanner(jsonl)

// 	for scanner.Scan() {
// 		var schemeConveyor *model.SchemeConveyor
// 		err := json.JSON.Unmarshal(scanner.Bytes(), &schemeConveyor)
// 		if err != nil {

// 		}
// 	}
// }

// func rollback(a *App, createdSchemeIDs []string) {
// 	for _, schemeID := range createdSchemeIDs {
// 		a.DeleteScheme(schemeID)
// 	}
// }
