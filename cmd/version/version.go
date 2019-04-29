package version

import (
	"fmt"

	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the version of GitLab CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("GitLab CLI %s-%s\n", config.Version, config.Commit)
			return nil
		},
	}
}
