package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/config"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the version of GitLab CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			var version strings.Builder
			fmt.Fprintf(&version, "GitLab CLI %s", config.Version)
			fmt.Println(version.String())
			return nil
		},
	}
}
