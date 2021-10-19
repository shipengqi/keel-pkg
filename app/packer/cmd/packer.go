package cmd

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/keel-pkg/lib/utils/fmtutil"
)

var (
	Version      string
	GitCommit    string
	BuildTime    = "1970-01-01T00:00:00Z"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "packer",
		Short: "keel packer for kubernetes",
		Version: Version,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	c.SetVersionTemplate(fmtutil.VersionTmpl("Keel Packer", Version, GitCommit, BuildTime))
	return c
}


