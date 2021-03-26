package zsh

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "generate zsh completion",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenZshCompletion(os.Stdout)
		},
	}
}
