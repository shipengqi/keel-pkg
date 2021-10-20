package cmd

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/shipengqi/keel-pkg/lib/deps"
	"github.com/shipengqi/keel-pkg/lib/utils/fsutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"time"

	"github.com/shipengqi/keel-pkg/lib/utils/fmtutil"
)

var (
	Version   = "<unknown>"
	GitCommit = "<unknown>"
	BuildTime = "1970-01-01T00:00:00Z"
)

type packOptions struct {
	RegistryUser string
	RegistryPass string
	VersionsFile string
	ImagesOutput string
	CmdTimeout   time.Duration
}

func New() *cobra.Command {
	o := &packOptions{}
	set := &deps.Versions{}

	c := &cobra.Command{
		Use:     "packer",
		Short:   "keel packer for kubernetes",
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
			return fsutil.MustMkDir(o.ImagesOutput)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pack(o, set)
		},
	}
	flags := c.Flags()
	flags.SortFlags = false
	c.DisableFlagsInUseLine = true

	addPackFlags(flags, o)
	addRegistryClientFlags(flags, o)

	c.SetVersionTemplate(fmtutil.VersionTmpl("Keel Packer",
		Version, GitCommit, BuildTime))
	return c
}

func addPackFlags(f *pflag.FlagSet, o *packOptions) {
	f.DurationVar(
		&o.CmdTimeout, "command-timeout", 0,
		"Set timeout for the command execution",
	)
	f.StringVar(
		&o.VersionsFile, "versions", "image_set.json",
		"The location of image-set file",
	)
	f.StringVarP(
		&o.ImagesOutput, "image-output", "o", DefaultImagesOutput,
		"The location of images output",
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
