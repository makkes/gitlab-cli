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
		Run: func(cmd *cobra.Command, args []string) {
			ids := strings.Split(args[0], ":")
			if len(ids) < 2 {
				fmt.Printf("ID must be of the form PROJECT_ID:ISSUE_ID\n")
				return
			}

			resp, err := client.Get("/projects/" + ids[0] + "/issues/" + ids[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			var out bytes.Buffer
			json.Indent(&out, resp, "", "    ")
			out.WriteTo(os.Stdout)
			fmt.Println()
		},
	}
	return cmd
}
