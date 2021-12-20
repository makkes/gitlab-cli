package project

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func createCommand(client api.Client, args []string, out io.Writer) error {
	res, _, err := client.Post("/projects", map[string]interface{}{
		"name": url.QueryEscape(args[0]),
	})
	if err != nil {
		return fmt.Errorf("Cannot create project: %s", err)
	}
	createdProject := make(map[string]interface{})
	err = json.Unmarshal(res, &createdProject)
	if err != nil {
		return fmt.Errorf("Unexpected result from GitLab API: %s", err)
	}
	fmt.Fprintf(out, `Project '%s' created
Clone via SSH: %s
Clone via HTTP: %s
`,
		createdProject["name"], createdProject["ssh_url_to_repo"], createdProject["http_url_to_repo"])

	return nil

}

func NewCommand(client api.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "project NAME",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createCommand(client, args, os.Stdout)
		},
	}
}
