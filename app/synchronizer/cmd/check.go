package cmd

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/keel-pkg/app/synchronizer/action"
	gcrc "github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"
)

func NewCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check if the image needs to be synchronized",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			action.NewCheckAction(gcrc.NewDefaultOptions())
		},
	}
}