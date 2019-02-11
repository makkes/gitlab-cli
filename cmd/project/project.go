package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/project/create"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project PROJECT",
		Short: "List details about a project by ID or name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := client.FindProjectDetails(args[0])
			if err != nil {
				return fmt.Errorf("Cannot show project: %s", err)
			}
			var out bytes.Buffer
			json.Indent(&out, project, "", "    ")
			out.WriteTo(os.Stdout)
			fmt.Println()
			return nil
		},
	}

	cmd.AddCommand(create.NewCommand(client))

	return cmd
}
