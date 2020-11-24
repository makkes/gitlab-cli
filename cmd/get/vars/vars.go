package vars

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/makkes/gitlab-cli/cmd/get/output"
	"github.com/makkes/gitlab-cli/table"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, project *string, format *string) *cobra.Command {
	return &cobra.Command{
		Use:   "vars",
		Short: "List variables in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == nil || *project == "" {
				return fmt.Errorf("please provide a project scope")
			}
			project, err := client.FindProject(*project)
			if err != nil {
				return fmt.Errorf("cannot list variables: %s", err)
			}
			resp, _, err := client.Get("/projects/" + strconv.Itoa(project.ID) + "/variables")
			if err != nil {
				return fmt.Errorf("cannot list variables: %s", err)
			}

			vars := make([]api.Var, 0)
			err = json.Unmarshal(resp, &vars)
			if err != nil {
				return fmt.Errorf("cannot list variables: %s", err)
			}

			return output.Print(resp, *format, os.Stdout, func() error {
				table.PrintVars(os.Stdout, vars)
				return nil
			}, func() error {
				for _, v := range vars {
					fmt.Fprintf(os.Stdout, "%s\n", v.Key)
				}
				return nil
			}, vars)
		},
	}
}
