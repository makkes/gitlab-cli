package remove

import (
	"fmt"
	"net/url"
	"os"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "remove PROJECT KEY",
		Short: "Remove a project-level variable",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			project, err := client.FindProject(args[0])
			if err != nil {
				fmt.Printf("Error finding project: %s\n", err)
				return
			}
			err = client.Delete(fmt.Sprintf("/projects/%d/variables/%s", project.ID, url.PathEscape(args[1])))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating variable: %s\n", err)
				return
			}
		},
	}
}
