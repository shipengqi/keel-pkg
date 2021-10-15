package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	baseName := filepath.Base(os.Args[0])
	c := &cobra.Command{
		Use:   baseName + " [options]",
		Short: "keel image synchronizer",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	c.AddCommand(NewSyncCommand())
	cobra.EnableCommandSorting = false
	c.CompletionOptions.DisableDefaultCmd = true
	c.DisableFlagsInUseLine = true
	return c
}
