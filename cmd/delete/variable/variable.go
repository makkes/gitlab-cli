package variable

import (
	"fmt"
	"net/url"
	"os"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, project *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vars KEY [KEY...]",
		Short: "Delete project-level variables",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if project == nil || *project == "" {
				return fmt.Errorf("please provide a project scope")
			}
			// silence errors now since we already log an error for every single variable
			cmd.SilenceErrors = true
			project, err := client.FindProject(*project)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error removing variables: %s\n", err)
				return
			}
			for _, key := range args[0:] {
				status, err := client.Delete(fmt.Sprintf("/projects/%d/variables/%s", project.ID, url.PathEscape(key)))
				if err != nil {
					if status == 404 {
						fmt.Fprintf(os.Stderr, "Error: no such variable: %s\n", key)
					} else {
						fmt.Fprintf(os.Stderr, "Error removing variable %s: %s\n", key, err)
					}
				} else {
					fmt.Printf("Removed variable %s\n", key)
				}
			}
			return
		},
	}

	return cmd
}
