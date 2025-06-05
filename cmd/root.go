package cmd

import (
	"github.com/spf13/cobra"

	"github.com/redjax/go-b2cleaner/internal/commands/clean_commands"
	"github.com/redjax/go-b2cleaner/internal/commands/list_commands"
)

var (
	bucket     string
	path       string
	appKey     string
	keyID      string
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "b2clean",
	Short: "A tool for cleaning up a Backblaze B2 bucket",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(list_commands.ListCmd)
	rootCmd.AddCommand(clean_commands.CleanCmd)

	rootCmd.PersistentFlags().StringVar(&appKey, "app-key", "", "B2 application key")
	rootCmd.PersistentFlags().StringVar(&keyID, "key-id", "", "B2 key ID")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "Path to config file")
}
