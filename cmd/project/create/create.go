package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "create NAME",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			res, _, err := client.Post("/projects", strings.NewReader(fmt.Sprintf("name=%s", url.QueryEscape(args[0]))))
			if err != nil {
				return fmt.Errorf("Cannot create project: %s", err)
			}
			createdProject := make(map[string]interface{})
			err = json.Unmarshal(res, &createdProject)
			if err != nil {
				return fmt.Errorf("Unexpected result from GitLab API: %s", err)
			}
			project, err := client.FindProjectDetails(strconv.Itoa(int(createdProject["id"].(float64))))
			if err != nil {
				return fmt.Errorf("Project has been created but cannot be shown: %s", err)
			}
			var out bytes.Buffer
			json.Indent(&out, project, "", "    ")
			out.WriteTo(os.Stdout)
			fmt.Println()
			return nil
		},
	}
}
