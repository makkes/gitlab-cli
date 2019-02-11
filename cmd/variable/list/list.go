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
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := client.FindProject(args[0])
			if err != nil {
				return fmt.Errorf("Cannot list variables: %s", err)
			}
			res, _, err := client.Get("/projects/" + strconv.Itoa(project.ID) + "/variables")
			if err != nil {
				return fmt.Errorf("Cannot list variables: %s", err)
			}

			vars := make([]api.Var, 0)
			err = json.Unmarshal(res, &vars)
			if err != nil {
				return fmt.Errorf("Cannot list variables: %s", err)
			}

			table.PrintVars(os.Stdout, vars)
			return nil
		},
	}
}
