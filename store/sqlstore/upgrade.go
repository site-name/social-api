package sqlstore

// upgradeDatabase attempts to migrate the schema to the latest supported version.
// The value of model.CurrentVersion is accepted as a parameter for unit testing, but it is not
// used to stop migrations at that version.
func upgradeDatabase(sqlStore *SqlStore, currentModelVersionString string) error {
	// currentModelVersion, err :=
	return nil
}
