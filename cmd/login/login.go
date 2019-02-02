package login

import (
	"fmt"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "login TOKEN",
		Short: "Login to Gitlab.com",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err, username := client.Login(args[0])
			if err != nil {
				fmt.Printf("Error logging you in: %s\n", err)
				return
			}
			fmt.Printf("Logged in as %s\n", username)
			cfg.Write()
		},
	}
}
