package clean

import (
	"fmt"

	"github.com/redjax/go-b2cleaner/internal/b2Ops"
	"github.com/redjax/go-b2cleaner/internal/config"
	"github.com/spf13/cobra"
)

var (
	ageStr  string
	dryRun  bool
	bucket  string
	path    string
	recurse bool
)

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
		return client.CleanObjects(bucket, path, ageStr, dryRun, recurse)
	},
}

func init() {
	CleanCmd.Flags().StringVar(&ageStr, "age", "", "Delete files older than this (e.g. 30d, 2m, 3y)")
	CleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be deleted, but don't delete")
	CleanCmd.Flags().StringVar(&bucket, "bucket", "", "Bucket name (overrides config)")
	CleanCmd.Flags().StringVar(&path, "path", "", "Path to clean (overrides config)")
	CleanCmd.Flags().BoolVar(&recurse, "recurse", false, "Recurse into subdirectories")
}
