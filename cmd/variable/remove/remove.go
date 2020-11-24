package remove

import (
	"fmt"
	"net/url"
	"os"

	"github.com/makkes/gitlab-cli/v3/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove PROJECT KEY [KEY...]",
		Short: "Remove project-level variables",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// silence errors since we already log an error for every single variable
			cmd.SilenceErrors = true
			project, err := client.FindProject(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error removing variables: %s\n", err)
				return
			}
			for _, key := range args[1:] {
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
