package sqlstore

import (
	"github.com/sitename/sitename/modules/slog"
)

// upgradeDatabase attempts to migrate the schema to the latest supported version.
// The value of model.CurrentVersion is accepted as a parameter for unit testing, but it is not
// used to stop migrations at that version.
func upgradeDatabase(sqlStore *SqlStore, currentModelVersionString string) error {
	// currentModelVersion, err :=
	slog.Warn("Upgrade is not implemented yet.")
	return nil
}

// func shouldPerformUpgrade(sqlStore *SqlStore, currentSchemaVersion string, expectedSchemaVersion string) bool {
// 	if sqlStore.GetCurrentSchemaVersion() == currentSchemaVersion {

// 	}
// }
