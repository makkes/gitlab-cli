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

func NewCommand(client api.Client, format *string) *cobra.Command {
	var project *string

	cmd := &cobra.Command{
		Use:   "vars [VARNAME]",
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
			if args[0] != "" {
				found := false
				for _, variable := range vars {
					if variable.Key == args[0] {
						found = true
						vars = []api.Var{variable}
						break
					}
				}
				if !found {
					// var requested that doesn't exist
					vars = []api.Var{}
				}
			}

			return output.NewPrinter(
				output.NoListWithSingleEntry(),
			).Print(*format, os.Stdout, func() error {
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

	project = cmd.PersistentFlags().StringP("project", "p", "", "If present, the project scope for this CLI request")

	return cmd
}
