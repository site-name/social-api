package cmd

import (
	"fmt"
	"strings"

	"errors"

	"github.com/sitename/sitename/models"
	"github.com/sitename/sitename/modules/setting"
	"github.com/sitename/sitename/modules/util"
	"github.com/urfave/cli"
)

// argsSet checks that all the required arguments are set. args is a list of
// arguments that must be set in the passed Context.
func argsSet(c *cli.Context, args ...string) error {
	for _, a := range args {
		if !c.IsSet(a) {
			return errors.New(a + " is not set")
		}

		if util.IsEmptyString(a) {
			return errors.New(a + " is required")
		}
	}
	return nil
}

// confirm waits for user input which confirms an action
func confirm() (bool, error) {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, errors.New(response + " isn't a correct confirmation string")
	}
}

func initDB() error {
	return initDBDisableConsole(false)
}

func initDBDisableConsole(disableConsole bool) error {
	setting.NewContext()
	setting.InitDBConfig()

	setting.NewXORMLogService(disableConsole)
	if err := models.SetEngine(); err != nil {
		return fmt.Errorf("models.SetEngine: %v", err)
	}
	return nil
}
