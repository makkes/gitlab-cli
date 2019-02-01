package project

import (
	"fmt"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "project PROJECT",
		Short: "List details about a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := client.FindProject(args[0])
			if err != nil {
				fmt.Printf("Error finding project: %s\n", err)
				return
			}
			table.PrintProjects(projects)
		},
	}
}
