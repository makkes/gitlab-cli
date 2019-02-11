package create

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "create PROJECT KEY VALUE",
		Short: "Create a project-level variable",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := client.FindProject(args[0])
			if err != nil {
				return fmt.Errorf("Cannot create variable: %s", err)
			}
			_, _, err = client.Post("/projects/"+strconv.Itoa(project.ID)+"/variables", strings.NewReader(fmt.Sprintf("key=%s&value=%s", url.QueryEscape(args[1]), url.QueryEscape(args[2]))))
			if err != nil {
				return fmt.Errorf("Cannot create variable: %s", err)
			}
			return nil
		},
	}
}
