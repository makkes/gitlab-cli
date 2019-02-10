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
		Run: func(cmd *cobra.Command, args []string) {
			project, err := client.FindProjectDetails(args[0])
			if err != nil {
				fmt.Printf("Error finding project: %s\n", err)
				return
			}
			var out bytes.Buffer
			json.Indent(&out, project, "", "    ")
			out.WriteTo(os.Stdout)
			fmt.Println()
		},
	}

	cmd.AddCommand(create.NewCommand(client))

	return cmd
}
