package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/shipengqi/keel-pkg/lib/utils/fmtutil"
)

const (
	DefaultPushToRegistry = "registry.cn-hangzhou.aliyuncs.com"
	DefaultPushToNs       = "keel"
	DefaultDBFile         = "sync.bolt.db"
	DefaultQueryLimit     = 10
	DefaultLimit          = 5
	DefaultRetryCount     = 5
)

var (
	Version   string
	GitCommit string
	BuildTime = "1970-01-01T00:00:00Z"
)

func New() *cobra.Command {
	baseName := filepath.Base(os.Args[0])
	c := &cobra.Command{
		Use:     baseName + " [options]",
		Short:   "keel image synchronizer",
		Version: Version,
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
	c.SetVersionTemplate(fmtutil.VersionTmpl("Keel Synctl", Version, GitCommit, BuildTime))
	return c
}
