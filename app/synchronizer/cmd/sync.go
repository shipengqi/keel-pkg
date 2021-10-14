package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/keel-pkg/app/synchronizer/action"
	"github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"
)

func NewSyncCommand() *cobra.Command {
	o := &action.SyncOptions{
		Options: &client.Options{},
	}

	c := &cobra.Command{
		Use:   "sync",
		Short: "Sync images",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cobra.EnableCommandSorting = false

	flags := c.Flags()
	flags.SortFlags = false

	addSyncFlags(flags, o)
	addRegistryClientFlags(flags, o)
	return c
}

func addSyncFlags(f *pflag.FlagSet, o *action.SyncOptions) {
	f.StringVar(
		&o.Db, "db", "bolt.db",
		"The location of boltdb file",
	)
	f.IntVar(
		&o.QueryLimit, "query-limit", 2,
		"Set http query limit",
	)
	f.IntVar(
		&o.Limit, "limit", 2,
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
		&o.Retry, "retry", 5,
		"Retry count.",
	)
	f.DurationVar(
		&o.RetryInterval, "retry-interval", 5*time.Second,
		"Retry interval",
	)
	f.StringSliceVar(
		&o.AdditionalNS, "addition-ns", nil,
		"Additional namespaces to sync")
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
		&o.PushToRepo, "push-to", "registry.cn-hangzhou.aliyuncs.com",
		"The registry to be pushed",
	)
	f.StringVar(
		&o.PushToNS, "push-ns", "keel",
		"The namespace of the registry to be pushed",
	)
}
