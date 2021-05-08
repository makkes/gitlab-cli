package jobs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/get/output"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, format *string) *cobra.Command {
	var pipeline *string
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "List jobs of a pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pipeline == nil || *pipeline == "" {
				return fmt.Errorf("please provide a pipeline scope")
			}
			ids := strings.Split(*pipeline, ":")
			if len(ids) < 2 || ids[0] == "" || ids[1] == "" {
				return fmt.Errorf("ID must be of the form PROJECT_ID:PIPELINE_ID")
			}
			resp, _, err := client.Get(fmt.Sprintf("/projects/%s/pipelines/%s/jobs", url.PathEscape(ids[0]), url.PathEscape(ids[1])))
			if err != nil {
				return fmt.Errorf("error retrieving jobs: %w", err)
			}

			var respSlice []map[string]interface{}
			err = json.Unmarshal(resp, &respSlice)
			if err != nil {
				return err
			}

			var jobs api.Jobs
			err = json.Unmarshal(resp, &jobs)
			if err != nil {
				return err
			}
			if len(jobs) <= 0 {
				return nil
			}

			projectID, err := strconv.Atoi(ids[0])
			if err != nil {
				return fmt.Errorf("error converting project ID '%s' to integer: %w", ids[0], err)
			}
			for idx := range jobs {
				jobs[idx].ProjectID = projectID
			}

			return output.Print(resp, *format, os.Stdout, func() error {
				table.PrintJobs(jobs)
				return nil
			}, func() error {
				for _, j := range jobs {
					fmt.Printf("%d\n", j.ID)
				}
				return nil
			}, jobs)
		},
	}

	pipeline = cmd.PersistentFlags().StringP("pipeline", "p", "", "The pipeline to show jobs from")

	return cmd
}
