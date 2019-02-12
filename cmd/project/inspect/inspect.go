package inspect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func inspectCommand(client api.Client, args []string, out io.Writer) error {
	project, err := client.FindProjectDetails(args[0])
	if err != nil {
		return fmt.Errorf("Cannot show project: %s", err)
	}
	var buf bytes.Buffer
	json.Indent(&buf, project, "", "    ")
	buf.WriteTo(out)
	out.Write([]byte("\n"))
	return nil
}

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect ID",
		Short: "Display detailed information on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return inspectCommand(client, args, os.Stdout)
		},
	}
	return cmd
}
