package inspect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect ID",
		Short: "Display detailed information on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := strings.Split(args[0], ":")
			if len(ids) < 2 || len(ids[0]) == 0 || len(ids[1]) == 0 {
				return fmt.Errorf("ID must be of the form PROJECT_ID:ISSUE_ID")
			}

			resp, status, err := client.Get("/projects/" + ids[0] + "/issues/" + ids[1])
			if err != nil {
				if status == 404 {
					return fmt.Errorf("Issue %s not found", args[0])
				}
				return err
			}
			var out bytes.Buffer
			json.Indent(&out, resp, "", "    ")
			out.WriteTo(os.Stdout)
			fmt.Println()
			return nil
		},
	}
	return cmd
}
