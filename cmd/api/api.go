package api

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
)

func NewCommand(client api.Client, cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api PATH",
		Short: "Makes an authenticated HTTP request to the GitLab API and prints the response.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rawRes, sc, err := client.Get(args[0])
			if err != nil {
				return fmt.Errorf("error making API request '%s' (got status code %d): %w", args[0], sc, err)
			}
			fmt.Printf("%s\n", rawRes)
			return nil
		},
	}

	return cmd
}
