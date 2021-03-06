package cmd

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/keel-pkg/app/synchronizer/action"
)

func NewSumCommand(done chan error) *cobra.Command {
	c := &cobra.Command{
		Use:   "sum [options]",
		Short: "List all check sum",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			a := action.NewSumAction(args[0])
			go receiver(a, done)
			return action.Execute(a)
		},
	}
	c.DisableFlagsInUseLine = true
	return c
}
