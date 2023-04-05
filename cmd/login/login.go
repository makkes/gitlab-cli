package login

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
)

func NewCommand(client api.Client, cfg config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "login [URL]",
		Short: "Login to GitLab. If URL is omitted then https://gitlab.com is used.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := "https://gitlab.com"
			if len(args) == 2 {
				url = args[1]
			}

			token, err := readPasswordFromStdin("Please enter your GitLab personal access token (PAT): ")
			if err != nil {
				return fmt.Errorf("failed reading token from stdin: %w", err)
			}

			username, err := client.Login(token, url)
			if err != nil {
				return fmt.Errorf("cannot login to %s: %s", url, err)
			}
			fmt.Printf("Logged in as %s\n", username)
			cfg.Write()
			return nil
		},
	}
}
