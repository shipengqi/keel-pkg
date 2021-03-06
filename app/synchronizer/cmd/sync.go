package cmd

import (
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/keel-pkg/app/synchronizer/action"
	"github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"
	"github.com/shipengqi/keel-pkg/lib/deps"
)

func NewSyncCommand(done chan error) *cobra.Command {
	o := &action.SyncOptions{
		Options: client.NewDefaultOptions(),
	}
	c := &cobra.Command{
		Use:   "sync [options]",
		Short: "Sync images",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if o.ImageSetFile == "" {
				return errors.New("invalid flag --image-set")
			}
			setBytes, err := os.ReadFile(o.ImageSetFile)
			if err != nil {
				return err
			}
			set := &deps.ImageSet{}
			err = jsoniter.Unmarshal(setBytes, set)
			if err != nil {
				return errors.Wrap(err, "unmarshal")
			}
			o.ImageSet = set
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			a := action.NewSyncAction(o)
			go receiver(a, done)
			return action.Execute(a)
		},
	}

	cobra.EnableCommandSorting = false

	flags := c.Flags()
	flags.SortFlags = false
	c.DisableFlagsInUseLine = true
	c.SilenceUsage = true

	addRegistryClientFlags(flags, o)
	addSyncFlags(flags, o)
	return c
}

func addSyncFlags(f *pflag.FlagSet, o *action.SyncOptions) {
	f.StringVar(
		&o.Db, "db", DefaultDBFile,
		"The location of boltdb file",
	)
	f.IntVar(
		&o.QueryLimit, "query-limit", DefaultQueryLimit,
		"Set http query limit",
	)
	f.IntVar(
		&o.Limit, "limit", DefaultLimit,
		"Set sync limit",
	)
	f.DurationVar(
		&o.CmdTimeout, "command-timeout", 0,
		"Set timeout for the command execution",
	)
	f.DurationVar(
		&o.PushTimeout, "push-timeout", 15*time.Minute,
		"Set timeout for pushing a image",
	)
	f.IntVar(
		&o.Retry, "retry", DefaultRetryCount,
		"Retry count.",
	)
	f.DurationVar(
		&o.RetryInterval, "retry-interval", 5*time.Second,
		"Retry interval",
	)
	f.StringSliceVar(
		&o.AdditionalNS, "addition-ns", nil,
		"Additional namespaces to sync")
	f.StringVar(
		&o.ImageSetFile, "image-set", "image_set.json",
		"The location of image-set file",
	)
}

func addRegistryClientFlags(f *pflag.FlagSet, o *action.SyncOptions) {
	f.StringVarP(
		&o.Username, "username", "u", "",
		"The username of the registry to be pushed",
	)
	f.StringVarP(
		&o.Password, "password", "p", "",
		"The password of the registry to be pushed",
	)
	f.StringVar(
		&o.PushToRepo, "push-to", DefaultPushToRegistry,
		"The registry to be pushed",
	)
	f.StringVar(
		&o.PushToNS, "push-ns", DefaultPushToNs,
		"The namespace of the registry to be pushed",
	)
}
