package issue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
)

func inspectCommand(client api.Client, args []string, out io.Writer) error {
	ids := strings.Split(args[0], ":")
	if len(ids) < 2 || len(ids[0]) == 0 || len(ids[1]) == 0 {
		return fmt.Errorf("ID must be of the form PROJECT_ID:ISSUE_ID")
	}

	resp, status, err := client.Get("/projects/" + ids[0] + "/issues/" + ids[1])
	if err != nil {
		if status == 404 {
			return fmt.Errorf("issue %s not found", args[0])
		}
		return err
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, resp, "", "    "); err != nil {
		return err
	}
	if _, err := buf.WriteTo(out); err != nil {
		return err
	}
	_, err = out.Write([]byte("\n"))
	return err
}

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue ID",
		Short: "Display detailed information on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return inspectCommand(client, args, os.Stdout)
		},
	}
	return cmd
}
