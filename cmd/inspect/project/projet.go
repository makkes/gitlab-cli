package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
)

func inspectCommand(client api.Client, args []string, out io.Writer) error {
	project, err := client.FindProjectDetails(args[0])
	if err != nil {
		return fmt.Errorf("cannot show project: %s", err)
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, project, "", "    "); err != nil {
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
		Use:   "project ID",
		Short: "Display detailed information on a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return inspectCommand(client, args, os.Stdout)
		},
	}
	return cmd
}
