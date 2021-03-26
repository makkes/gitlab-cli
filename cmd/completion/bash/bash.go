package bash

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use: "bash",
		Short: "generate bash completion",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
}
