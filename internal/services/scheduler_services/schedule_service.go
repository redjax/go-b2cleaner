package scheduler_services

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func ScheduleCommandJob(
	cronExpr string,
	logFile string,
	rootCmd *cobra.Command,
	commandLine []string,
	persistentFlags []string,
) error {
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
		fullArgs := append(persistentFlags, commandLine...)
		logger.Printf("Starting scheduled job: %v at %s", strings.Join(fullArgs, " "), start.Format(time.RFC3339))

		// SetArgs on the actual rootCmd, then Execute
		rootCmd.SetArgs(fullArgs)
		if err := rootCmd.Execute(); err != nil {
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
