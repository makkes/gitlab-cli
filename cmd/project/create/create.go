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
		Short: "Display detailed information on an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			res, err := client.Post("/projects", strings.NewReader(fmt.Sprintf("name=%s", url.QueryEscape(args[0]))))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating project: %s\n", err)
				return
			}
			createdProject := make(map[string]interface{})
			err = json.Unmarshal(res, &createdProject)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unexpected result from GitLab API: %s\n", err)
				return
			}
			project, err := client.FindProjectDetails(strconv.Itoa(int(createdProject["id"].(float64))))
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
}
