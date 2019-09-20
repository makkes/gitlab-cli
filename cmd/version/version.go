package version

import (
	"fmt"
	"strings"

	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the version of GitLab CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			var version strings.Builder
			fmt.Fprintf(&version, "GitLab CLI %s", config.Version)
			if config.Commit != "" {
				fmt.Fprintf(&version, "-%s", config.Commit)
			}
			fmt.Println(version.String())
			return nil
		},
	}
}
