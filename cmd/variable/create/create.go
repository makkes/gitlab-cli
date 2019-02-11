package create

import (
	"fmt"
	"net/url"
	"os"
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
		Run: func(cmd *cobra.Command, args []string) {
			project, err := client.FindProject(args[0])
			if err != nil {
				fmt.Printf("Error finding project: %s\n", err)
				return
			}
			_, err = client.Post("/projects/"+strconv.Itoa(project.ID)+"/variables", strings.NewReader(fmt.Sprintf("key=%s&value=%s", url.QueryEscape(args[1]), url.QueryEscape(args[2]))))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating variable: %s\n", err)
				return
			}
		},
	}
}
