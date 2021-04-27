package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sitename/sitename/config"
	"github.com/sitename/sitename/modules/slog"
	"github.com/spf13/cobra"
)

var JobserverCmd = &cobra.Command{
	Use:   "jobserver",
	Short: "Start the Sitenname job server",
	RunE:  jobServerCmdF,
}

func init() {
	JobserverCmd.Flags().Bool("nojobs", false, "Do not run jobs in this jobserver.")
	JobserverCmd.Flags().Bool("noschedule", false, "Do not schedule jobs from this jobserver.")

	RootCmd.AddCommand(JobserverCmd)
}

func jobServerCmdF(command *cobra.Command, args []string) error {
	noJobs, _ := command.Flags().GetBool("nojobs")
	noSchedule, _ := command.Flags().GetBool("noschedule")

	// Initialize
	a, err := initDBCommandContext(getConfigDSN(command, config.GetEnvironment()), false)
	if err != nil {
		return err
	}

	defer a.Srv().Shutdown()

	a.InitServer()

	// Run jobs
	slog.Info("Starting Sitename job server")
	defer slog.Info("Stopped Sitename job server")

	if !noJobs {
		a.Srv().Jobs.StartWorkers()
		defer a.Srv().Jobs.StopWorkers()
	}
	if !noSchedule {
		a.Srv().Jobs.StartSchedulers()
		defer a.Srv().Jobs.StopSchedulers()
	}

	// if !noJobs || !noSchedule {
	// 	auditRec := a.MakeAuditRecord("jobServer", audit.Success)
	// 	a.LogAuditRec(auditRec, nil)
	// }

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	slog.Info("Stopping Sitename job server")

	return nil
}
