package status

import (
	"fmt"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, cfg config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Display the current configuration of GitLab CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := cfg.Get("url")
			if url == "" {
				return fmt.Errorf("GitLab CLI is not configured, yet. Run »gitlab-cli login« first")
			}
			fmt.Printf("Logged in at %s as %s\n", url, cfg.Get("user"))
			return nil
		},
	}
}
