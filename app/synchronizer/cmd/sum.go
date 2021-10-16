package cmd

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/keel-pkg/app/synchronizer/action"
)

func NewSumCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sum",
		Short: "List all check sum",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			a := action.NewSumAction()
			return action.Execute(a)
		},
	}
}