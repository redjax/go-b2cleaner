package scheduler_services

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func ScheduleCommandJob(cronExpr, logFile string, rootCmd *cobra.Command, commandLine []string) error {
	if cronExpr == "" {
		return ErrMissingCronExpr
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logger := log.New(f, "", log.LstdFlags)
	c := cron.New()

	_, err = c.AddFunc(cronExpr, func() {
		start := time.Now()
		logger.Printf("Starting scheduled job: %v at %s", commandLine, start.Format(time.RFC3339))

		// Clone the root command to avoid flag state issues
		cmd := &cobra.Command{
			Use:   rootCmd.Use,
			Short: rootCmd.Short,
			Run:   rootCmd.Run,
		}
		// Add all subcommands
		for _, sub := range rootCmd.Commands() {
			cmd.AddCommand(sub)
		}
		cmd.SetArgs(commandLine)
		if err := cmd.Execute(); err != nil {
			logger.Printf("Job failed: %v", err)
		} else {
			logger.Printf("Completed in %s", time.Since(start))
		}
	})

	if err != nil {
		return err
	}

	c.Start()
	select {} // Block indefinitely
}

var ErrMissingCronExpr = &MissingCronExprError{}

type MissingCronExprError struct{}

func (e *MissingCronExprError) Error() string {
	return "cron expression required"
}
