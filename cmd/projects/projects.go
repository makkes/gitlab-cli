package projects

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	var quiet *bool
	var format *string

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

			if *format != "" {
				tmpl, err := template.New("").Parse(*format)
				if err != nil {
					fmt.Printf("Template parsing error: %s\n", err)
					return
				}
				for _, p := range projects {
					err = tmpl.Execute(os.Stdout, p)
					if err != nil {
						fmt.Printf("Template parsing error: %s\n", err)
					} else {
						fmt.Println()
					}
				}
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
	format = cmd.Flags().String("format", "", "Pretty-print projects using a Go template")

	return cmd
}
