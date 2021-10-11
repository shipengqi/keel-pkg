package cmd

import (
	"github.com/spf13/cobra"
)

func NewSyncCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync images",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	return rootCmd
}
