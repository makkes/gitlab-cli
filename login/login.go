package login

import (
	"fmt"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "login TOKEN USER",
		Short: "Login to Gitlab.com",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Logging in using %s %s\n", args[0], args[1])
			config.WriteLoginCredentials(args[0], args[1])
		},
	}
}
