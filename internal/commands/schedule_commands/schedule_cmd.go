package schedule_commands

import (
	"fmt"

	"github.com/redjax/go-b2cleaner/internal/commands/clean_commands"
	"github.com/redjax/go-b2cleaner/internal/services/scheduler_services"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewScheduleCommand creates the "schedule" parent command with subcommands for each schedule-able action.
func NewScheduleCommand(rootCmd *cobra.Command) *cobra.Command {
	var cronExpr string
	var logFile string

	scheduleCmd := &cobra.Command{
		Use:   "schedule",
		Short: "Schedule any command to run on a cron schedule",
	}

	// Schedule clean
	cleanSchedule := &cobra.Command{
		Use:   "clean",
		Short: "Schedule a clean operation",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cronExpr == "" {
				return fmt.Errorf("cron expression required")
			}

			// Collect persistent flags from the root command
			persistentFlags := []string{}
			root := cmd.Root()
			if f := root.Flag("config-file"); f != nil && f.Value.String() != "" {
				persistentFlags = append(persistentFlags, "-c", f.Value.String())
			}
			if f := root.Flag("app-key"); f != nil && f.Value.String() != "" {
				persistentFlags = append(persistentFlags, "--app-key", f.Value.String())
			}
			if f := root.Flag("key-id"); f != nil && f.Value.String() != "" {
				persistentFlags = append(persistentFlags, "--key-id", f.Value.String())
			}

			// Gather clean-specific flags and args
			cleanArgs := []string{"clean"}
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if f.Changed {
					if f.Value.Type() == "bool" && f.Value.String() == "true" {
						cleanArgs = append(cleanArgs, "--"+f.Name)
					} else if f.Value.Type() == "filetypes" {
						// Handle repeated filetype flags
						for _, val := range f.Value.String() {
							cleanArgs = append(cleanArgs, "--"+f.Name, string(val))
						}
					} else {
						cleanArgs = append(cleanArgs, "--"+f.Name, f.Value.String())
					}
				}
			})

			return scheduler_services.ScheduleCommandJob(
				cronExpr,
				logFile,
				rootCmd,
				cleanArgs,
				persistentFlags,
			)
		},
	}
	// Inherit all flags from the clean command
	cleanSchedule.Flags().AddFlagSet(clean_commands.CleanCmd.Flags())

	// Add schedule-specific flags to the parent schedule command
	scheduleCmd.PersistentFlags().StringVar(&cronExpr, "cron", "", "Cron schedule expression")
	scheduleCmd.PersistentFlags().StringVar(&logFile, "log", "scheduled-jobs.log", "Log file path")

	// Add subcommands
	scheduleCmd.AddCommand(cleanSchedule)
	// You can add more scheduled subcommands (e.g. list) here similarly

	return scheduleCmd
}
