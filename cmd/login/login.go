package login

import (
	"fmt"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient, cfg config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "login TOKEN [URL]",
		Short: "Login to GitLab. If URL is omitted then https://gitlab.com is used.",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			url := "https://gitlab.com"
			if len(args) == 2 {
				url = args[1]
			}
			err, username := client.Login(args[0], url)
			if err != nil {
				fmt.Printf("Error logging you in: %s\n", err)
				return
			}
			fmt.Printf("Logged in as %s\n", username)
			cfg.Write()
		},
	}
}
