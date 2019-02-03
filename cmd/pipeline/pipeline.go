package pipeline

import (
	"strings"
	"os"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/makkes/gitlab-cli/api"
)

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "pipeline ID",
		Short: "List details of a pipeline",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ids := strings.Split(args[0], ":")
			if len(ids) < 2 {
				fmt.Printf("ID must be of the form PROJECT_ID:PIPELINE_ID\n")
				return
			}
			pipeline, err := client.GetPipelineDetails(ids[0], ids[1])
			if err != nil {
				fmt.Printf("Error finding pipeline: %s\n", err)
				return
			}
			var out bytes.Buffer
			json.Indent(&out, pipeline, "", "    ")
			out.WriteTo(os.Stdout)
			fmt.Println()},
	}
}
