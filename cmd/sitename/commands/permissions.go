package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var PermissionsCmd = &cobra.Command{
	Use:   "permissions",
	Short: "Management of the Permissions system",
}

var ResetPermissionsCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Reset the permissions system to its default state",
	Long:    "Reset the permissions system to its default state",
	Example: "  permissions reset",
	RunE:    resetPermissionsCmdF,
}

// var ExportPermissionsCmd = &cobra.Command{
// 	Use:     "export",
// 	Short:   "Export permissions data",
// 	Long:    "Export Roles and Schemes to JSONL for use by Mattermost permissions import.",
// 	Example: " permissions export > export.jsonl",
// 	RunE:    exportPermissionsCmdF,
// 	PreRun: func(cmd *cobra.Command, args []string) {
// 		os.Setenv("SN_LOGSETTINGS_CONSOLELEVEL", "error")
// 	},
// }

// var ImportPermissionsCmd = &cobra.Command{
// 	Use:     "import [file]",
// 	Short:   "Import permissions data",
// 	Long:    "Import Roles and Schemes JSONL data as created by the Mattermost permissions export.",
// 	Example: " permissions import export.jsonl",
// 	RunE:    importPermissionsCmdF,
// }

func init() {
	ResetPermissionsCmd.Flags().Bool("confirm", false, "Confirm you really want to reset the permissions system and a database backup has been performed.")

	PermissionsCmd.AddCommand(
		ResetPermissionsCmd,
		// ExportPermissionsCmd,
		// ImportPermissionsCmd,
	)
	RootCmd.AddCommand(PermissionsCmd)
}

func resetPermissionsCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	confirmFlag, _ := command.Flags().GetBool("confirm")
	if !confirmFlag {
		var confirm string
		CommandPrettyPrintln("Have you performed a database backup? (YES/NO): ")
		fmt.Scanln(&confirm)

		if confirm != "YES" {
			return errors.New("ABORTED: You did not answer YES exactly, in all capitals.")
		}

		CommandPrettyPrintln("Are you sure you want to reset the permissions system? All data related to the permissions system will be permanently deleted and all users will revert to having the default permissions. (YES/NO): ")
		fmt.Scanln(&confirm)
		if confirm != "YES" {
			return errors.New("ABORTED: You did not answer YES exactly, in all capitals.")
		}
	}

	if err := a.ResetPermissionsSystem(); err != nil {
		return errors.New(err.Error())
	}

	CommandPrettyPrintln("Permissions system successfully reset.")
	CommandPrettyPrintln("Changes will take effect gradually as the server caches expire.")
	CommandPrettyPrintln("For the changes to take effect immediately, go to the Mattermost System Console > General > Configuration and click \"Purge All Caches\".")

	// auditRec := a.MakeAuditRecord("resetPermissions", audit.Success)
	// a.LogAuditRec(auditRec, nil)

	return nil
}

// func exportPermissionsCmdF(command *cobra.Command, args []string) error {
// 	a, err := InitDBCommandContextCobra(command)
// 	if err != nil {
// 		return err
// 	}
// 	defer a.Srv().Shutdown()

// 	if err = a.ExportPermissions(os.Stdout); err != nil {
// 		return errors.New(err.Error())
// 	}

// 	// auditRec := a.MakeAuditRecord("exportPermissions", audit.Success)
// 	// a.LogAuditRec(auditRec, nil)

// 	return nil
// }

// func importPermissionsCmdF(command *cobra.Command, args []string) error {
// 	a, err := InitDBCommandContextCobra(command)
// 	if err != nil {
// 		return err
// 	}
// 	defer a.Srv().Shutdown()

// 	file, err := os.Open(args[0])
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// auditRec := a.MakeAuditRecord("importPermissions", audit.Success)
// 	// auditRec.AddMeta("file", args[0])
// 	// a.LogAuditRec(auditRec, nil)

// 	return a.ImportPermissions(file)
// }
