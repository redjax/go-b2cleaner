package schedule_commands

import (
	"fmt"

	"github.com/redjax/go-b2cleaner/internal/services/scheduler_services"
	"github.com/spf13/cobra"
)

func NewScheduleCommand(rootCmd *cobra.Command) *cobra.Command {
	var cronExpr string
	var logFile string

	cmd := &cobra.Command{
		Use:   "schedule [flags] <command> [command flags]",
		Short: "Schedule any command to run on a cron schedule",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cronExpr == "" {
				return fmt.Errorf("cron expression required")
			}
			// args[0] is the subcommand (e.g. "clean" or "list")
			// args[1:] are the flags for that subcommand

			// Prepare the command string
			commandLine := append([]string{args[0]}, args[1:]...)

			return scheduler_services.ScheduleCommandJob(
				cronExpr,
				logFile,
				rootCmd,
				commandLine,
			)
		},
	}

	cmd.Flags().StringVar(&cronExpr, "cron", "", "Cron schedule expression")
	cmd.Flags().StringVar(&logFile, "log", "scheduled-jobs.log", "Log file path")
	return cmd
}
