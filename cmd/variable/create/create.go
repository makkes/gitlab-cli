package create

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/makkes/gitlab-cli/v3/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "create PROJECT KEY VALUE [ENVIRONMENT_SCOPE]",
		Short: "Create a project-level variable",
		Long:  "Create a project-level variable. The KEY may only contain the characters A-Z, a-z, 0-9, and _ and must be no longer than 255 characters.",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := client.FindProject(args[0])
			if err != nil {
				return fmt.Errorf("Cannot create variable: %s", err)
			}

			body := fmt.Sprintf("key=%s&value=%s", url.QueryEscape(args[1]), url.QueryEscape(args[2]))
			if len(args) == 4 {
				body += fmt.Sprintf("&environment_scope=%s", url.QueryEscape(args[3]))
			}
			_, _, err = client.Post("/projects/"+strconv.Itoa(project.ID)+"/variables",
				strings.NewReader(body))
			if err != nil {
				return fmt.Errorf("Cannot create variable: %s", err)
			}
			return nil
		},
	}
}
