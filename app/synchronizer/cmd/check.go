package cmd

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/keel-pkg/app/synchronizer/action"
	gcrc "github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"
)

func NewCheckCommand(done chan error) *cobra.Command {
	o := &action.CheckOptions{
		Options: gcrc.NewDefaultOptions(),
	}
	c := &cobra.Command{
		Use:   "check [options]",
		Short: "Check if the image needs to be synchronized",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.CheckSum = args[0]
			a := action.NewCheckAction(o)
			go receiver(a, done)
			return action.Execute(a)
		},
	}
	c.DisableFlagsInUseLine = true

	c.Flags().StringVar(&o.Db, "db", DefaultDBFile, "The location of boltdb file")

	return c
}
