package list

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/makkes/gitlab-cli/table"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list PROJECT",
		Short: "List a project's variables",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			project, err := client.FindProject(args[0])
			if err != nil {
				fmt.Printf("Error finding project: %s\n", err)
				return
			}
			res, err := client.Get("/projects/" + strconv.Itoa(project.ID) + "/variables")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating variable: %s\n", err)
				return
			}

			vars := make([]api.Var, 0)
			err = json.Unmarshal(res, &vars)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error querying variables: %s\n", err)
			}

			table.PrintVars(os.Stdout, vars)
		},
	}
}
