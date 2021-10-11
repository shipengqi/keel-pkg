package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "packer",
		Short: "keel packer for kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	return rootCmd
}
