package vars

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, project *string) *cobra.Command {
	return &cobra.Command{
		Use:   "var KEY VALUE [ENVIRONMENT_SCOPE]",
		Short: "Create a project-level variable",
		Long:  "Create a project-level variable. The KEY may only contain the characters A-Z, a-z, 0-9, and _ and must be no longer than 255 characters.",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == nil || *project == "" {
				return fmt.Errorf("please provide a project scope")
			}
			project, err := client.FindProject(*project)
			if err != nil {
				return fmt.Errorf("Cannot create variable: %s", err)
			}

			body := fmt.Sprintf("key=%s&value=%s", url.QueryEscape(args[0]), url.QueryEscape(args[1]))
			if len(args) == 3 {
				body += fmt.Sprintf("&environment_scope=%s", url.QueryEscape(args[2]))
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
