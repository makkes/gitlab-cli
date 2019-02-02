package projects

import (
	"encoding/json"
	"fmt"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	var quiet *bool
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List all your projects",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := client.Get("/users/${user}/projects")
			if err != nil {
				fmt.Println(err)
				return
			}
			projects := make([]api.Project, 0)
			err = json.Unmarshal(resp, &projects)
			if err != nil {
				fmt.Println(err)
				return
			}

			if *quiet {
				for _, p := range projects {
					fmt.Println(p.ID)
				}
				return
			}
			table.PrintProjects(projects)
		},
	}

	quiet = cmd.Flags().BoolP("quiet", "q", false, "Only display numeric IDs")

	return cmd
}
