package cmd

import (
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/shipengqi/keel-pkg/lib/deps"
	"github.com/shipengqi/keel-pkg/lib/utils/fsutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/keel-pkg/lib/utils/fmtutil"
)

var (
	Version   = "<unknown>"
	GitCommit = "<unknown>"
	BuildTime = "1970-01-01T00:00:00Z"
)

type packOptions struct {
	RegistryUser   string
	RegistryPass   string
	VersionsFile   string
	ImagesOutput   string
	DownloadOutput string
	CmdTimeout     time.Duration
}

func New() *cobra.Command {
	o := &packOptions{}
	set := &deps.Versions{}

	c := &cobra.Command{
		Use:     "packer [options]",
		Short:   "Pack kubernetes.tar.gz",
		Version: Version,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if o.VersionsFile == "" {
				return errors.New("invalid flag --versions")
			}
			setBytes, err := os.ReadFile(o.VersionsFile)
			if err != nil {
				return err
			}
			err = jsoniter.Unmarshal(setBytes, set)
			if err != nil {
				return errors.Wrap(err, "unmarshal")
			}
			err = fsutil.MustMkDir(o.ImagesOutput)
			if err != nil {
				return err
			}
			return fsutil.MustMkDir(o.DownloadOutput)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pack(o, set)
		},
	}
	flags := c.Flags()
	flags.SortFlags = false
	c.DisableFlagsInUseLine = true
	c.SilenceUsage = true
	c.CompletionOptions.DisableDefaultCmd = true

	addRegistryClientFlags(flags, o)
	addPackFlags(flags, o)

	c.SetVersionTemplate(fmtutil.VersionTmpl("Keel Packer",
		Version, GitCommit, BuildTime))

	c.AddCommand(newPushCommand())

	return c
}

func addPackFlags(f *pflag.FlagSet, o *packOptions) {
	f.DurationVar(
		&o.CmdTimeout, "command-timeout", 0,
		"Set timeout for the command execution",
	)
	f.StringVar(
		&o.VersionsFile, "version-config", "versions.json",
		"The location of versions config file",
	)
	f.StringVarP(
		&o.ImagesOutput, "image-output", "o", DefaultImagesOutput,
		"The location of images output",
	)
	f.StringVarP(
		&o.DownloadOutput, "download-output", "d", DefaultDownloadOutput,
		"The location of download output",
	)
}

func addRegistryClientFlags(f *pflag.FlagSet, o *packOptions) {
	f.StringVarP(
		&o.RegistryUser, "username", "u", "",
		"The username of the registry to be pushed",
	)
	f.StringVarP(
		&o.RegistryPass, "password", "p", "",
		"The password of the registry to be pushed",
	)
}
