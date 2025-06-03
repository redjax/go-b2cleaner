package list

import (
	"fmt"

	"github.com/redjax/go-b2cleaner/internal/b2Ops"
	"github.com/redjax/go-b2cleaner/internal/config"
	"github.com/spf13/cobra"
)

var (
	bucket  string
	path    string
	recurse bool
	sortBy  string
	order   string
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List objects in a B2 bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig(cmd)

		validSorts := map[string]bool{
			"name":    true,
			"size":    true,
			"created": true,
		}
		if !validSorts[cfg.SortBy] {
			return fmt.Errorf("invalid sort option: %s (must be one of: name, size, created)", cfg.SortBy)
		}

		client := b2Ops.NewClient(cfg)

		return client.ListObjects(cfg.Bucket, cfg.Path, cfg.Recurse)
	},
}

func init() {
	ListCmd.Flags().StringVar(&bucket, "bucket", "", "Bucket name (required)")
	ListCmd.Flags().StringVar(&path, "path", "", "Path to list objects in")
	ListCmd.Flags().BoolVar(&recurse, "recurse", false, "Recurse into subdirectories")
	ListCmd.Flags().StringVar(&sortBy, "sort", "name", "Sort by: name, size, or created")
	ListCmd.Flags().StringVar(&order, "order", "asc", "Sort order: asc or desc (default: asc)")

	// ListCmd.MarkFlagRequired("bucket")
}
