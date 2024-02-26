package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/audit"
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Management of users",
}

var UserActivateCmd = &cobra.Command{
	Use:   "activate [emails, usernames, userIds]",
	Short: "Activate users",
	Long:  "Activate users that have been deactivated.",
	Example: `  user activate user@example.com
  user activate username`,
	RunE: userActivateCmdF,
}

// var UserDeactivateCmd = &cobra.Command{
// 	Use:   "deactivate [emails, usernames, userIds]",
// 	Short: "Deactivate users",
// 	Long:  "Deactivate users. Deactivated users are immediately logged out of all sessions and are unable to log back in.",
// 	Example: `  user deactivate user@example.com
//   user deactivate username`,
// 	RunE: userDeactivateCmdF,
// }

var UserCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a user",
	Long:    "Create a user",
	Example: `  user create --email user@example.com --username userexample --password Password1`,
	RunE:    userCreateCmdF,
}

var ResetUserPasswordCmd = &cobra.Command{
	Use:     "password [user] [password]",
	Short:   "Set a user's password",
	Long:    "Set a user's password",
	Example: "  user password user@example.com Password1",
	RunE:    resetUserPasswordCmdF,
}

var updateUserEmailCmd = &cobra.Command{
	Use:     "email [user] [new email]",
	Short:   "Change email of the user",
	Long:    "Change email of the user.",
	Example: "  user email testuser user@example.com",
	RunE:    updateUserEmailCmdF,
}

var ResetUserMfaCmd = &cobra.Command{
	Use:   "resetmfa [users]",
	Short: "Turn off MFA",
	Long: `Turn off multi-factor authentication for a user.
If MFA enforcement is enabled, the user will be forced to re-enable MFA as soon as they login.`,
	Example: "  user resetmfa user@example.com",
	RunE:    resetUserMfaCmdF,
}

var DeleteUserCmd = &cobra.Command{
	Use:     "delete [users]",
	Short:   "Delete users and all posts",
	Long:    "Permanently delete user and all related information including posts.",
	Example: "  user delete user@example.com",
	RunE:    deleteUserCmdF,
}

var DeleteAllUsersCmd = &cobra.Command{
	Use:     "deleteall",
	Short:   "Delete all users and all posts",
	Long:    "Permanently delete all users and all related information including posts.",
	Example: "  user deleteall",
	RunE:    deleteAllUsersCommandF,
}

var MigrateAuthCmd = &cobra.Command{
	Use:     "migrate_auth [from_auth] [to_auth] [migration-options]",
	Short:   "Mass migrate user accounts authentication type",
	Long:    `Migrates accounts from one authentication provider to another. For example, you can upgrade your authentication provider from email to ldap.`,
	Example: "  user migrate_auth email saml users.json",
	Args: func(command *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Auth migration requires at least 2 arguments.")
		}

		toAuth := args[1]

		if toAuth != "ldap" && toAuth != "saml" {
			return errors.New("Invalid to_auth parameter, must be saml or ldap.")
		}

		if toAuth == "ldap" && len(args) != 3 {
			return errors.New("Ldap migration requires 3 arguments.")
		}

		autoFlag, _ := command.Flags().GetBool("auto")

		if toAuth == "saml" && autoFlag {
			if len(args) != 2 {
				return errors.New("Saml migration requires two arguments when using the --auto flag. See help text for details.")
			}
		}

		if toAuth == "saml" && !autoFlag {
			if len(args) != 3 {
				return errors.New("Saml migration requires three arguments when not using the --auto flag. See help text for details.")
			}
		}
		return nil
	},
	RunE: migrateAuthCmdF,
}

var VerifyUserCmd = &cobra.Command{
	Use:     "verify [users]",
	Short:   "Verify email of users",
	Long:    "Verify the emails of some users.",
	Example: "  user verify user1",
	RunE:    verifyUserCmdF,
}

var SearchUserCmd = &cobra.Command{
	Use:     "search [users]",
	Short:   "Search for users",
	Long:    "Search for users based on username, email, or user ID.",
	Example: "  user search user1@mail.com user2@mail.com",
	RunE:    searchUserCmdF,
}

func init() {
	// UserCreateCmd.Flags().String("username", "", "Required. Username for the new user model.")
	// UserCreateCmd.Flags().String("email", "", "Required. The email address for the new user model.")
	// UserCreateCmd.Flags().String("password", "", "Required. The password for the new user model.")
	// UserCreateCmd.Flags().String("nickname", "", "Optional. The nickname for the new user model.")
	// UserCreateCmd.Flags().String("firstname", "", "Optional. The first name for the new user model.")
	// UserCreateCmd.Flags().String("lastname", "", "Optional. The last name for the new user model.")
	// UserCreateCmd.Flags().String("locale", "", "Optional. The locale (ex: en, fr) for the new user model.")
	// UserCreateCmd.Flags().Bool("system_admin", false, "Optional. If supplied, the new user will be a system administrator. Defaults to false.")

	DeleteUserCmd.Flags().Bool("confirm", false, "Confirm you really want to delete the user and a DB backup has been performed.")

	DeleteAllUsersCmd.Flags().Bool("confirm", false, "Confirm you really want to delete the user and a DB backup has been performed.")

	MigrateAuthCmd.Flags().Bool("force", false, "Force the migration to occur even if there are duplicates on the LDAP server. Duplicates will not be migrated. (ldap only)")
	MigrateAuthCmd.Flags().Bool("auto", false, "Automatically migrate all users. Assumes the usernames and emails are identical between Mattermost and SAML services. (saml only)")
	MigrateAuthCmd.Flags().Bool("dryRun", false, "Run a simulation of the migration process without changing the database.")
	MigrateAuthCmd.SetUsageTemplate(`Usage:
  mattermost user migrate_auth [from_auth] [to_auth] [migration-options] [flags]

Examples:
{{.Example}}

Arguments:
  from_auth:
    The authentication service to migrate users accounts from.
    Supported options: email, gitlab, ldap, saml.

  to_auth:
    The authentication service to migrate users to.
    Supported options: ldap, saml.

  migration-options:
    Migration specific options, full command help for more information.

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
`)
	MigrateAuthCmd.SetHelpTemplate(`Usage:
  mattermost user migrate_auth [from_auth] [to_auth] [migration-options] [flags]

Examples:
{{.Example}}

Arguments:
  from_auth:
    The authentication service to migrate users accounts from.
    Supported options: email, gitlab, ldap, saml.

  to_auth:
    The authentication service to migrate users to.
    Supported options: ldap, saml.

  migration-options (ldap):
    match_field:
      The field that is guaranteed to be the same in both authentication services. For example, if the users emails are consistent set to email.
      Supported options: email, username.

  migration-options (saml):
    users_file:
      The path of a json file with the usernames and emails of all users to migrate to SAML. The username and email must be the same that the SAML service provider store. And the email must match with the email in mattermost database.

      Example json content:
        {
          "usr1@email.com": "usr.one",
          "usr2@email.com": "usr.two"
        }

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
`)

	UserCmd.AddCommand(
		UserActivateCmd,
		UserCreateCmd,
		ResetUserPasswordCmd,
		updateUserEmailCmd,
		ResetUserMfaCmd,
		DeleteUserCmd,
		DeleteAllUsersCmd,
		MigrateAuthCmd,
		VerifyUserCmd,
		SearchUserCmd,
	)
	RootCmd.AddCommand(UserCmd)
}

func userCreateCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	type prompt struct {
		name  string
		iface interface{}
	}

	var (
		username        string
		email           string
		password        string
		passwordConfirm string
		nickname        string
		firstname       string
		lastname        string
		locale          model.LanguageCode
		systemAdmin     bool
	)
	fmt.Println("---------------CREATE USER--------------")
	var prompts = []prompt{
		{name: "username", iface: &username},
		{name: "email", iface: &email},
		{name: "password", iface: &password},
		{name: "passwordConfirm", iface: &passwordConfirm},
		{name: "nickname", iface: &nickname},
		{name: "firstname", iface: &firstname},
		{name: "lastname", iface: &lastname},
		{name: "locale", iface: &locale},
		{name: "systemAdmin", iface: &systemAdmin},
	}
	for _, prompt := range prompts {
		fmt.Print(prompt.name + ": ")
		fmt.Scanln(prompt.iface)
	}

	if username == "" {
		return errors.New("username is required")
	}
	if email == "" {
		return errors.New("email is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	if passwordConfirm == "" {
		return errors.New("you have to confirm your password")
	}
	if passwordConfirm != password {
		return errors.New("passwords do not match")
	}

	user := model.User{
		Username:  username,
		Email:     email,
		Password:  password,
		Nickname:  nickname,
		FirstName: firstname,
		LastName:  lastname,
		Locale:    locale,
	}

	ruser, err := a.Srv().AccountService().CreateUser(request.Context{}, user)
	if ruser == nil {
		return errors.New("Unable to create user. Error: " + err.Error())
	}

	if systemAdmin {
		if _, err := a.Srv().AccountService().UpdateUserRolesWithUser(*ruser, "system_user system_admin", false); err != nil {
			return errors.New("Unable to make user system admin. Error: " + err.Error())
		}
	} else {
		// This else case exists to prevent the first user created from being
		// created as a system admin unless explicitly specified.
		if _, err := a.Srv().AccountService().UpdateUserRolesWithUser(*ruser, "system_user", false); err != nil {
			return errors.New("If this is the first user: Unable to prevent user from being system admin. Error: " + err.Error())
		}
	}

	CommandPrettyPrintln("id: " + ruser.ID)
	CommandPrettyPrintln("username: " + ruser.Username)
	CommandPrettyPrintln("nickname: " + ruser.Nickname)
	CommandPrettyPrintln("first_name: " + ruser.FirstName)
	CommandPrettyPrintln("last_name: " + ruser.LastName)
	CommandPrettyPrintln("email: " + ruser.Email)
	CommandPrettyPrintln("auth_service: " + ruser.AuthService)

	auditRec := a.MakeAuditRecord("userCreate", audit.Success)
	auditRec.AddMeta("user", ruser)
	auditRec.AddMeta("system_admin", systemAdmin)
	a.LogAuditRec(auditRec, nil)

	return nil
}

func searchUserCmdF(cmd *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(cmd)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) < 1 {
		return errors.New("Expected at elast one argument. See help text for details.")
	}

	users := getUsersFromUserArgs(a, args)

	for i, user := range users {
		if i > 0 {
			CommandPrettyPrintln("------------------------------")
		}
		if user == nil {
			CommandPrintErrorln("Unable to find user '" + args[i] + "'")
			continue
		}

		CommandPrettyPrintln("id: " + user.ID)
		CommandPrettyPrintln("username: " + user.Username)
		CommandPrettyPrintln("nickname: " + user.Nickname)
		// CommandPrettyPrintln("position: " + user.Position)
		CommandPrettyPrintln("first_name: " + user.FirstName)
		CommandPrettyPrintln("last_name: " + user.LastName)
		CommandPrettyPrintln("email: " + user.Email)
		CommandPrettyPrintln("auth_service: " + user.AuthService)
	}

	return nil
}

func userActivateCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) < 1 {
		return errors.New("Expected at least one argument. See help text for details.")
	}

	changeUsersActiveStatus(a, args, true)

	return nil
}

func changeUsersActiveStatus(a *app.App, userArgs []string, active bool) {
	users := getUsersFromUserArgs(a, userArgs)
	for i, user := range users {
		err := changeUserActiveStatus(a, user, userArgs[i], active)

		if err != nil {
			CommandPrintErrorln(err.Error())
		}
	}
}

func changeUserActiveStatus(a *app.App, user *model.User, userArg string, activate bool) error {
	if user == nil {
		return fmt.Errorf("Can't find user '%v'", userArg)
	}

	if model_helper.UserIsSSO(*user) {
		fmt.Println("You must also deactivate this user in the SSO provider or they will be reactivated on next login or sync.")
	}
	updatedUser, err := a.Srv().AccountService().UpdateActive(&request.Context{}, *user, activate)
	if err != nil {
		return fmt.Errorf("Unable to change activation status of user: %v", userArg)
	}

	auditRec := a.MakeAuditRecord("changeActiveUserStatus", audit.Success)
	auditRec.AddMeta("user", updatedUser)
	auditRec.AddMeta("activate", activate)
	a.LogAuditRec(auditRec, nil)

	return nil
}

func verifyUserCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) < 1 {
		return errors.New("Expected at least one argument. See help text for details.")
	}

	users := getUsersFromUserArgs(a, args)

	for i, user := range users {
		if user == nil {
			CommandPrintErrorln("Unable to find user '" + args[i] + "'")
			continue
		}
		if _, err := a.Srv().Store.User().VerifyEmail(user.ID, user.Email); err != nil {
			CommandPrintErrorln("Unable to verify '" + args[i] + "' email. Error: " + err.Error())
		}
	}

	return nil
}

func migrateAuthCmdF(command *cobra.Command, args []string) error {
	if args[1] == "saml" {
		return migrateAuthToSamlCmdF(command, args)
	}
	return migrateAuthToLdapCmdF(command, args)
}

func migrateAuthToLdapCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	fromAuth := args[0]
	matchField := args[2]

	if fromAuth == "" || (fromAuth != "email" && fromAuth != "gitlab" && fromAuth != "saml") {
		return errors.New("Invalid from_auth argument")
	}

	// Email auth in Mattermost system is represented by ""
	if fromAuth == "email" {
		fromAuth = ""
	}

	if matchField == "" || (matchField != "email" && matchField != "username") {
		return errors.New("Invalid match_field argument")
	}

	forceFlag, _ := command.Flags().GetBool("force")
	dryRunFlag, _ := command.Flags().GetBool("dryRun")

	if migrate := a.AccountMigration(); migrate != nil {
		if err := migrate.MigrateToLdap(fromAuth, matchField, forceFlag, dryRunFlag); err != nil {
			return errors.New("Error while migrating users: " + err.Error())
		}

		CommandPrettyPrintln("Successfully migrated accounts.")

		if !dryRunFlag {
			auditRec := a.MakeAuditRecord("migrateAuthToLdap", audit.Success)
			auditRec.AddMeta("fromAuth", fromAuth)
			auditRec.AddMeta("matchField", matchField)
			auditRec.AddMeta("force", forceFlag)
			a.LogAuditRec(auditRec, nil)
		}
	}
	return nil
}

func migrateAuthToSamlCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	dryRunFlag, _ := command.Flags().GetBool("dryRun")
	autoFlag, _ := command.Flags().GetBool("auto")

	matchesFile := ""
	matches := map[string]string{}
	if !autoFlag {
		matchesFile = args[2]

		file, e := ioutil.ReadFile(matchesFile)
		if e != nil {
			return errors.New("Invalid users file.")
		}
		if json.Unmarshal(file, &matches) != nil {
			return errors.New("Invalid users file.")
		}
	}

	fromAuth := args[0]

	if fromAuth == "" || (fromAuth != "email" && fromAuth != "gitlab" && fromAuth != "ldap") {
		return errors.New("Invalid from_auth argument")
	}

	if autoFlag && !dryRunFlag {
		var confirm string
		CommandPrettyPrintln("You are about to perform an automatic \"" + fromAuth + " to saml\" migration. This must only be done if your current Mattermost users with " + fromAuth + " auth have the same username and email in your SAML service. Otherwise, provide the usernames and emails from your SAML Service using the \"users file\" without the \"--auto\" option.\n\nDo you want to proceed with automatic migration anyway? (YES/NO):")
		fmt.Scanln(&confirm)

		if confirm != "YES" {
			return errors.New("ABORTED: You did not answer YES exactly, in all capitals.")
		}
	}

	// Email auth in Mattermost system is represented by ""
	if fromAuth == "email" {
		fromAuth = ""
	}

	if migrate := a.AccountMigration(); migrate != nil {
		if err := migrate.MigrateToSaml(fromAuth, matches, autoFlag, dryRunFlag); err != nil {
			return errors.New("Error while migrating users: " + err.Error())
		}

		CommandPrettyPrintln("Successfully migrated accounts.")

		if !dryRunFlag {
			auditRec := a.MakeAuditRecord("migrateAuthToSaml", audit.Success)
			auditRec.AddMeta("auto", autoFlag)
			a.LogAuditRec(auditRec, nil)
		}
	}
	return nil
}

func deleteAllUsersCommandF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) > 0 {
		return errors.New("Expected zero arguments.")
	}

	confirmFlag, _ := command.Flags().GetBool("confirm")
	if !confirmFlag {
		var confirm string
		CommandPrettyPrintln("Have you performed a database backup? (YES/NO): ")
		fmt.Scanln(&confirm)

		if confirm != "YES" {
			return errors.New("ABORTED: You did not answer YES exactly, in all capitals.")
		}
		CommandPrettyPrintln("Are you sure you want to permanently delete all user accounts? (YES/NO): ")
		fmt.Scanln(&confirm)
		if confirm != "YES" {
			return errors.New("ABORTED: You did not answer YES exactly, in all capitals.")
		}
	}

	if err := a.Srv().AccountService().PermanentDeleteAllUsers(&request.Context{}); err != nil {
		return err
	}
	CommandPrettyPrintln("All user accounts successfully deleted.")

	auditRec := a.MakeAuditRecord("deleteAllUsers", audit.Success)
	a.LogAuditRec(auditRec, nil)

	return nil
}

func deleteUserCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) < 1 {
		return errors.New("Expected at least one argument. See help text for details.")
	}

	confirmFlag, _ := command.Flags().GetBool("confirm")
	if !confirmFlag {
		var confirm string
		CommandPrettyPrintln("Have you performed a database backup? (YES/NO): ")
		fmt.Scanln(&confirm)

		if confirm != "YES" {
			return errors.New("ABORTED: You did not answer YES exactly, in all capitals.")
		}
		CommandPrettyPrintln("Are you sure you want to permanently delete the specified users? (YES/NO): ")
		fmt.Scanln(&confirm)
		if confirm != "YES" {
			return errors.New("ABORTED: You did not answer YES exactly, in all capitals.")
		}
	}

	users := getUsersFromUserArgs(a, args)

	for i, user := range users {
		if user == nil {
			return errors.New("Unable to find user '" + args[i] + "'")
		}

		if err := a.Srv().AccountService().PermanentDeleteUser(&request.Context{}, *user); err != nil {
			return err
		}

		auditRec := a.MakeAuditRecord("deleteUser", audit.Success)
		auditRec.AddMeta("user", user)
		a.LogAuditRec(auditRec, nil)
	}

	return nil
}

func resetUserMfaCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) < 1 {
		return errors.New("Expected at least one argument. See help text for details.")
	}

	users := getUsersFromUserArgs(a, args)
	for i, user := range users {
		if user == nil {
			return errors.New("Unable to find user '" + args[i] + "'")
		}

		if err := a.Srv().AccountService().DeactivateMfa(user.ID); err != nil {
			return err
		}

		auditRec := a.MakeAuditRecord("resetUserMfa", audit.Success)
		auditRec.AddMeta("user", user)
		a.LogAuditRec(auditRec, nil)
	}

	return nil
}

func updateUserEmailCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) != 2 {
		return errors.New("Expected two arguments. See help text for details.")
	}

	newEmail := args[1]
	newEmail = strings.ToLower(newEmail)
	if !model_helper.IsValidEmail(newEmail) {
		return errors.New("Invalid email: '" + newEmail + "'")
	}

	if len(args) != 2 {
		return errors.New("Expected two arguments. See help text for details.")
	}

	user := getUserFromUserArg(a, args[0])
	if user == nil {
		return errors.New("Unable to find user '" + args[0] + "'")
	}

	user.Email = newEmail
	_, errUpdate := a.Srv().AccountService().UpdateUser(*user, true)
	if errUpdate != nil {
		return errors.New(errUpdate.Message)
	}

	auditRec := a.MakeAuditRecord("updateUserEmail", audit.Success)
	auditRec.AddMeta("user", user)
	auditRec.AddMeta("email", newEmail)
	a.LogAuditRec(auditRec, nil)

	return nil
}

func resetUserPasswordCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	if len(args) != 2 {
		return errors.New("Expected two arguments. See help text for details.")
	}

	user := getUserFromUserArg(a, args[0])
	if user == nil {
		return errors.New("Unable to find user '" + args[0] + "'")
	}
	password := args[1]

	if err := a.Srv().Store.User().UpdatePassword(user.ID, model_helper.HashPassword(password)); err != nil {
		return err
	}

	auditRec := a.MakeAuditRecord("resetUserPassword", audit.Success)
	auditRec.AddMeta("user", user)
	a.LogAuditRec(auditRec, nil)

	return nil
}
