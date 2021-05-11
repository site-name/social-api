package commands

import (
	"errors"
	"strings"

	// "github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Management of users",
}

// var UserActivateCmd = &cobra.Command{
// 	Use:   "activate [emails, usernames, userIds]",
// 	Short: "Activate users",
// 	Long:  "Activate users that have been deactivated.",
// 	Example: `  user activate user@example.com
//   user activate username`,
// 	RunE: userActivateCmdF,
// }

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

// var ResetUserPasswordCmd = &cobra.Command{
// 	Use:     "password [user] [password]",
// 	Short:   "Set a user's password",
// 	Long:    "Set a user's password",
// 	Example: "  user password user@example.com Password1",
// 	RunE:    resetUserPasswordCmdF,
// }

// var updateUserEmailCmd = &cobra.Command{
// 	Use:     "email [user] [new email]",
// 	Short:   "Change email of the user",
// 	Long:    "Change email of the user.",
// 	Example: "  user email testuser user@example.com",
// 	RunE:    updateUserEmailCmdF,
// }

// var ResetUserMfaCmd = &cobra.Command{
// 	Use:   "resetmfa [users]",
// 	Short: "Turn off MFA",
// 	Long: `Turn off multi-factor authentication for a user.
// If MFA enforcement is enabled, the user will be forced to re-enable MFA as soon as they login.`,
// 	Example: "  user resetmfa user@example.com",
// 	RunE:    resetUserMfaCmdF,
// }

// var DeleteUserCmd = &cobra.Command{
// 	Use:     "delete [users]",
// 	Short:   "Delete users and all posts",
// 	Long:    "Permanently delete user and all related information including posts.",
// 	Example: "  user delete user@example.com",
// 	RunE:    deleteUserCmdF,
// }

// var DeleteAllUsersCmd = &cobra.Command{
// 	Use:     "deleteall",
// 	Short:   "Delete all users and all posts",
// 	Long:    "Permanently delete all users and all related information including posts.",
// 	Example: "  user deleteall",
// 	RunE:    deleteAllUsersCommandF,
// }

// var MigrateAuthCmd = &cobra.Command{
// 	Use:     "migrate_auth [from_auth] [to_auth] [migration-options]",
// 	Short:   "Mass migrate user accounts authentication type",
// 	Long:    `Migrates accounts from one authentication provider to another. For example, you can upgrade your authentication provider from email to ldap.`,
// 	Example: "  user migrate_auth email saml users.json",
// 	Args: func(command *cobra.Command, args []string) error {
// 		if len(args) < 2 {
// 			return errors.New("Auth migration requires at least 2 arguments.")
// 		}

// 		toAuth := args[1]

// 		if toAuth != "ldap" && toAuth != "saml" {
// 			return errors.New("Invalid to_auth parameter, must be saml or ldap.")
// 		}

// 		if toAuth == "ldap" && len(args) != 3 {
// 			return errors.New("Ldap migration requires 3 arguments.")
// 		}

// 		autoFlag, _ := command.Flags().GetBool("auto")

// 		if toAuth == "saml" && autoFlag {
// 			if len(args) != 2 {
// 				return errors.New("Saml migration requires two arguments when using the --auto flag. See help text for details.")
// 			}
// 		}

// 		if toAuth == "saml" && !autoFlag {
// 			if len(args) != 3 {
// 				return errors.New("Saml migration requires three arguments when not using the --auto flag. See help text for details.")
// 			}
// 		}
// 		return nil
// 	},
// 	RunE: migrateAuthCmdF,
// }

// var VerifyUserCmd = &cobra.Command{
// 	Use:     "verify [users]",
// 	Short:   "Verify email of users",
// 	Long:    "Verify the emails of some users.",
// 	Example: "  user verify user1",
// 	RunE:    verifyUserCmdF,
// }

// var SearchUserCmd = &cobra.Command{
// 	Use:     "search [users]",
// 	Short:   "Search for users",
// 	Long:    "Search for users based on username, email, or user ID.",
// 	Example: "  user search user1@mail.com user2@mail.com",
// 	RunE:    searchUserCmdF,
// }

func init() {
	UserCreateCmd.Flags().String("username", "", "Required. Username for the new user account.")
	UserCreateCmd.Flags().String("email", "", "Required. The email address for the new user account.")
	UserCreateCmd.Flags().String("password", "", "Required. The password for the new user account.")
	UserCreateCmd.Flags().String("nickname", "", "Optional. The nickname for the new user account.")
	UserCreateCmd.Flags().String("firstname", "", "Optional. The first name for the new user account.")
	UserCreateCmd.Flags().String("lastname", "", "Optional. The last name for the new user account.")
	UserCreateCmd.Flags().String("locale", "", "Optional. The locale (ex: en, fr) for the new user account.")
	UserCreateCmd.Flags().Bool("system_admin", false, "Optional. If supplied, the new user will be a system administrator. Defaults to false.")

	UserCmd.AddCommand(UserCreateCmd)
	RootCmd.AddCommand(UserCmd)
}

// func init() {
// 	UserCreateCmd.Flags().String("username", "", "Required. Username for the new user account.")
// 	UserCreateCmd.Flags().String("email", "", "Required. The email address for the new user account.")
// 	UserCreateCmd.Flags().String("password", "", "Required. The password for the new user account.")
// 	UserCreateCmd.Flags().String("nickname", "", "Optional. The nickname for the new user account.")
// 	UserCreateCmd.Flags().String("firstname", "", "Optional. The first name for the new user account.")
// 	UserCreateCmd.Flags().String("lastname", "", "Optional. The last name for the new user account.")
// 	UserCreateCmd.Flags().String("locale", "", "Optional. The locale (ex: en, fr) for the new user account.")
// 	UserCreateCmd.Flags().Bool("system_admin", false, "Optional. If supplied, the new user will be a system administrator. Defaults to false.")

// 	DeleteUserCmd.Flags().Bool("confirm", false, "Confirm you really want to delete the user and a DB backup has been performed.")

// 	DeleteAllUsersCmd.Flags().Bool("confirm", false, "Confirm you really want to delete the user and a DB backup has been performed.")

// 	MigrateAuthCmd.Flags().Bool("force", false, "Force the migration to occur even if there are duplicates on the LDAP server. Duplicates will not be migrated. (ldap only)")
// 	MigrateAuthCmd.Flags().Bool("auto", false, "Automatically migrate all users. Assumes the usernames and emails are identical between Mattermost and SAML services. (saml only)")
// 	MigrateAuthCmd.Flags().Bool("dryRun", false, "Run a simulation of the migration process without changing the database.")
// 	MigrateAuthCmd.SetUsageTemplate(`Usage:
//   mattermost user migrate_auth [from_auth] [to_auth] [migration-options] [flags]

// Examples:
// {{.Example}}

// Arguments:
//   from_auth:
//     The authentication service to migrate users accounts from.
//     Supported options: email, gitlab, ldap, saml.

//   to_auth:
//     The authentication service to migrate users to.
//     Supported options: ldap, saml.

//   migration-options:
//     Migration specific options, full command help for more information.

// Flags:
// {{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

// Global Flags:
// {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
// `)
// 	MigrateAuthCmd.SetHelpTemplate(`Usage:
//   mattermost user migrate_auth [from_auth] [to_auth] [migration-options] [flags]

// Examples:
// {{.Example}}

// Arguments:
//   from_auth:
//     The authentication service to migrate users accounts from.
//     Supported options: email, gitlab, ldap, saml.

//   to_auth:
//     The authentication service to migrate users to.
//     Supported options: ldap, saml.

//   migration-options (ldap):
//     match_field:
//       The field that is guaranteed to be the same in both authentication services. For example, if the users emails are consistent set to email.
//       Supported options: email, username.

//   migration-options (saml):
//     users_file:
//       The path of a json file with the usernames and emails of all users to migrate to SAML. The username and email must be the same that the SAML service provider store. And the email must match with the email in mattermost database.

//       Example json content:
//         {
//           "usr1@email.com": "usr.one",
//           "usr2@email.com": "usr.two"
//         }

// Flags:
// {{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

// Global Flags:
// {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
// `)

// 	UserCmd.AddCommand(
// 		// UserActivateCmd,
// 		// UserDeactivateCmd,
// 		UserCreateCmd,
// 		// UserConvertCmd,
// 		// UserInviteCmd,
// 		ResetUserPasswordCmd,
// 		updateUserEmailCmd,
// 		ResetUserMfaCmd,
// 		DeleteUserCmd,
// 		DeleteAllUsersCmd,
// 		MigrateAuthCmd,
// 		VerifyUserCmd,
// 		SearchUserCmd,
// 	)

// 	RootCmd.AddCommand(UserCmd)
// }

// // func userActivateCmdF(command *cobra.Command, args []string) error {
// // 	a, err := InitDBCommandContextCobra(command)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer a.Srv().Shutdown()

// // 	if len(args) < 1 {
// // 		return errors.New("Expected at least one argument. See help text for details.")
// // 	}

// // 	changeUsersActiveStatus(a, args, true)

// // 	return nil
// // }

// // func changeUsersActiveStatus(a *app.App, userArgs []string, active bool) {
// // 	users := getUsersFromUserArgs(a, userArgs)
// // 	for i, user := range users {
// // 		err := changeUserActiveStatus(a, user, userArgs[i], active)

// // 		if err != nil {
// // 			CommandPrintErrorln(err.Error())
// // 		}
// // 	}
// // }

// // func changeUserActiveStatus(a *app.App, user *model.User, userArg string, activate bool) error {
// // 	if user == nil {
// // 		return fmt.Errorf("Can't find user '%v'", userArg)
// // 	}
// // 	if user.IsSSOUser() {
// // 		fmt.Println("You must also deactivate this user in the SSO provider or they will be reactivated on next login or sync.")
// // 	}
// // 	updatedUser, err := a.UpdateActive(user, activate)
// // 	if err != nil {
// // 		return fmt.Errorf("Unable to change activation status of user: %v", userArg)
// // 	}

// // 	auditRec := a.MakeAuditRecord("changeActiveUserStatus", audit.Success)
// // 	auditRec.AddMeta("user", updatedUser)
// // 	auditRec.AddMeta("activate", activate)
// // 	a.LogAuditRec(auditRec, nil)

// // 	return nil
// // }

func userCreateCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	username, erru := command.Flags().GetString("username")
	if erru != nil || username == "" {
		return errors.New("Username is required")
	}
	email, erre := command.Flags().GetString("email")
	if erre != nil || email == "" {
		return errors.New("Email is required")
	}
	email = strings.ToLower((email))
	password, errp := command.Flags().GetString("password")
	if errp != nil || password == "" {
		return errors.New("Password is required")
	}
	nickname, _ := command.Flags().GetString("nickname")
	firstname, _ := command.Flags().GetString("firstname")
	lastname, _ := command.Flags().GetString("lastname")
	locale, _ := command.Flags().GetString("locale")
	systemAdmin, _ := command.Flags().GetBool("system_admin")

	user := &account.User{
		Username:  username,
		Email:     email,
		Password:  password,
		Nickname:  nickname,
		FirstName: firstname,
		LastName:  lastname,
		Locale:    locale,
	}

	ruser, err := a.CreateUser(user)
	if ruser == nil {
		return errors.New("Unable to create user. Error: " + err.Error())
	}

	if systemAdmin {
		if _, err := a.UpdateUserRolesWithUser(ruser, "system_user system_admin", false); err != nil {
			return errors.New("Unable to make user system admin. Error: " + err.Error())
		}
	} else {
		// This else case exists to prevent the first user created from being
		// created as a system admin unless explicitly specified.
		if _, err := a.UpdateUserRolesWithUser(ruser, "system_user", false); err != nil {
			return errors.New("If this is the first user: Unable to prevent user from being system admin. Error: " + err.Error())
		}
	}

	CommandPrettyPrintln("id: " + ruser.Id)
	CommandPrettyPrintln("username: " + ruser.Username)
	CommandPrettyPrintln("nickname: " + ruser.Nickname)
	// CommandPrettyPrintln("position: " + ruser.Position)
	CommandPrettyPrintln("first_name: " + ruser.FirstName)
	CommandPrettyPrintln("last_name: " + ruser.LastName)
	CommandPrettyPrintln("email: " + ruser.Email)
	CommandPrettyPrintln("auth_service: " + ruser.AuthService)

	// auditRec := a.MakeAuditRecord("userCreate", audit.Success)
	// auditRec.AddMeta("user", ruser)
	// auditRec.AddMeta("system_admin", systemAdmin)
	// a.LogAuditRec(auditRec, nil)

	return nil
}

// func resetUserPasswordCmdF(command *cobra.Command, args []string) error {
// 	a, err := InitDBCommandContextCobra(command)
// 	if err != nil {
// 		return err
// 	}
// 	defer a.Srv().Shutdown()

// 	if len(args) != 2 {
// 		return errors.New("Expected two arguments. See help text for details.")
// 	}

// 	user := getUserFromUserArg(a, args[0])
// 	if user == nil {
// 		return errors.New("Unable to find user '" + args[0] + "'")
// 	}
// 	password := args[1]

// 	if err := a.Srv().Store.User().UpdatePassword(user.Id, model.HashPassword(password)); err != nil {
// 		return err
// 	}

// 	// auditRec := a.MakeAuditRecord("resetUserPassword", audit.Success)
// 	// auditRec.AddMeta("user", user)
// 	// a.LogAuditRec(auditRec, nil)

// 	return nil
// }

// func updateUserEmailCmdF(command *cobra.Command, args []string) error {
// 	a, err := InitDBCommandContextCobra(command)
// 	if err != nil {
// 		return err
// 	}
// 	defer a.Srv().Shutdown()

// 	if len(args) != 2 {
// 		return errors.New("Expected two arguments. See help text for details.")
// 	}

// 	newEmail := args[1]
// 	newEmail = strings.ToLower(newEmail)
// 	if !model.IsValidEmail(newEmail) {
// 		return errors.New("Invalid email: '" + newEmail + "'")
// 	}

// 	if len(args) != 2 {
// 		return errors.New("Expected two arguments. See help text for details.")
// 	}

// 	user := getUserFromUserArg(a, args[0])
// 	if user == nil {
// 		return errors.New("Unable to find user '" + args[0] + "'")
// 	}

// 	user.Email = newEmail
// 	_, errUpdate := a.UpdateUser(user, true)
// 	if errUpdate != nil {
// 		return errors.New(errUpdate.Message)
// 	}

// 	// auditRec := a.MakeAuditRecord("updateUserEmail", audit.Success)
// 	// auditRec.AddMeta("user", user)
// 	// auditRec.AddMeta("email", newEmail)
// 	// a.LogAuditRec(auditRec, nil)

// 	return nil
// }
