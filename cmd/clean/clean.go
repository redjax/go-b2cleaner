package clean

import (
	"fmt"
	"strings"

	"github.com/redjax/go-b2cleaner/internal/b2Ops"
	"github.com/redjax/go-b2cleaner/internal/config"
	"github.com/spf13/cobra"
)

var (
	ageStr     string
	dryRun     bool
	bucket     string
	path       string
	recurse    bool
	outputPath string
	filetypes  filetypesFlag
)

type filetypesFlag []string

func (f *filetypesFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *filetypesFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *filetypesFlag) Type() string {
	return "filetypes"
}

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up a B2 bucket",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig(cmd)

		// Use flag if set, otherwise use config value
		if bucket == "" {
			bucket = cfg.Bucket
		}
		if path == "" {
			path = cfg.Path
		}
		if bucket == "" || path == "" {
			return fmt.Errorf("bucket and path must be set (either in config or via flags)")
		}

		client := b2Ops.NewClient(cfg)
		return client.CleanObjects(bucket, path, ageStr, dryRun, recurse, outputPath, filetypes)
	},
}

func init() {
	CleanCmd.Flags().StringVar(&ageStr, "age", "", "Delete files older than this (e.g. 30d, 2m, 3y)")
	CleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be deleted, but don't delete. Backblaze's B2 API causes this flag to sometimes be inconsistent, additional files may be deleted that do not show up in a dry run.")
	CleanCmd.Flags().StringVar(&bucket, "bucket", "", "Bucket name (overrides config)")
	CleanCmd.Flags().StringVar(&path, "path", "", "Path to clean (overrides config)")
	CleanCmd.Flags().BoolVar(&recurse, "recurse", false, "Recurse into subdirectories")
	CleanCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Write deleted objects to CSV file")
	CleanCmd.Flags().Var(&filetypes, "filetype", "Only delete files with these extensions (can be specified multiple times, e.g. --filetype=backup --filetype=.jpg)")

}
