package list

import (
	"github.com/redjax/go-b2cleaner/internal/b2Ops"
	"github.com/redjax/go-b2cleaner/internal/config"
	"github.com/spf13/cobra"
)

var (
	bucket  string
	path    string
	recurse bool
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List objects in a B2 bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig(cmd)
		client := b2Ops.NewClient(cfg)

		return client.ListObjects(cfg.Bucket, cfg.Path, cfg.Recurse)
	},
}

func init() {
	ListCmd.Flags().StringVar(&bucket, "bucket", "", "Bucket name (required)")
	ListCmd.Flags().StringVar(&path, "path", "", "Path to list objects in")
	ListCmd.Flags().BoolVar(&recurse, "recurse", false, "Recurse into subdirectories")

	// ListCmd.MarkFlagRequired("bucket")
}
