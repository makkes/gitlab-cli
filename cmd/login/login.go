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
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := "https://gitlab.com"
			if len(args) == 2 {
				url = args[1]
			}
			err, username := client.Login(args[0], url)
			if err != nil {
				return fmt.Errorf("Cannot login to %s: %s", url, err)
			}
			fmt.Printf("Logged in as %s\n", username)
			cfg.Write()
			return nil
		},
	}
}
