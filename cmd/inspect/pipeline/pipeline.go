package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
)

func NewCommand(client api.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "pipeline ID",
		Short: "List details of a pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := strings.Split(args[0], ":")
			if len(ids) < 2 || ids[0] == "" || ids[1] == "" {
				return fmt.Errorf("ID must be of the form PROJECT_ID:PIPELINE_ID")
			}
			pipeline, err := client.GetPipelineDetails(ids[0], ids[1])
			if err != nil {
				return fmt.Errorf("cannot show pipeline: %s", err)
			}
			var out bytes.Buffer
			if err := json.Indent(&out, pipeline, "", "    "); err != nil {
				return err
			}
			if _, err := out.WriteTo(os.Stdout); err != nil {
				return nil
			}
			fmt.Println()
			return nil
		},
	}
}
