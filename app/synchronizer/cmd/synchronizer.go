package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	DefaultPushToRegistry = "registry.cn-hangzhou.aliyuncs.com"
	DefaultPushToNs       = "keel"
	DefaultDBFile         = "sync.bolt.db"
	DefaultQueryLimit     = 10
	DefaultLimit          = 5
	DefaultRetryCount     = 5
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

	c.AddCommand(
		NewSyncCommand(),
		NewCheckCommand(),
		NewSumCommand(),
	)
	cobra.EnableCommandSorting = false
	c.CompletionOptions.DisableDefaultCmd = true
	c.DisableFlagsInUseLine = true
	return c
}
