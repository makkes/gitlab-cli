package pipelines

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/get/output"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, format *string) *cobra.Command {
	var all *bool
	var recent *bool
	var projectFlag *string

	cmd := &cobra.Command{
		Use:   "pipelines",
		Short: "List pipelines in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectFlag == nil || *projectFlag == "" {
				return fmt.Errorf("please provide a project scope")
			}
			project, err := client.FindProject(*projectFlag)
			if err != nil {
				return fmt.Errorf("cannot list pipelines: %s", err)
			}
			resp, _, err := client.Get("/projects/" + strconv.Itoa(project.ID) + "/pipelines")
			if err != nil {
				return err
			}

			var respSlice []map[string]interface{}
			err = json.Unmarshal(resp, &respSlice)
			if err != nil {
				return err
			}

			var pipelines api.Pipelines
			err = json.Unmarshal(resp, &pipelines)
			if err != nil {
				return err
			}
			if len(pipelines) <= 0 {
				return nil
			}
			if *recent {
				pipelines = []api.Pipeline{pipelines[0]}
				respSlice = respSlice[0:1]
			}
			if len(pipelines) <= 0 {
				return fmt.Errorf("no pipelines found for project '%s'", *projectFlag)
			}

			filteredRespSlice := make([]map[string]interface{}, 0)
			filteredPipelines := pipelines.Filter(func(idx int, p api.Pipeline) bool {
				include := *all || p.Status == "running" || p.Status == "pending"
				if include {
					filteredRespSlice = append(filteredRespSlice, respSlice[idx])
				}
				return include
			})

			pds := make([]api.PipelineDetails, 0)
			for _, p := range filteredPipelines {
				resp, _, err = client.Get("/projects/" + strconv.Itoa(project.ID) + "/pipelines/" + strconv.Itoa(p.ID))
				if err != nil {
					return fmt.Errorf("Error retrieving pipeline %d: %s", p.ID, err)
				}
				var pd api.PipelineDetails
				err = json.Unmarshal(resp, &pd)
				if err != nil {
					fmt.Println(err)
				} else {
					pd.ProjectID = project.ID
					pds = append(pds, pd)
				}
			}

			resp, err = json.Marshal(filteredRespSlice)
			if err != nil {
				return err
			}

			return output.Print(resp, *format, os.Stdout, func() error {
				table.PrintPipelines(pds)
				return nil
			}, func() error {
				for _, p := range pds {
					fmt.Printf("%d:%d\n", p.ProjectID, p.ID)
				}
				return nil
			}, pds)
		},
	}

	projectFlag = cmd.PersistentFlags().StringP("project", "p", "", "If present, the project scope for this CLI request")
	all = cmd.Flags().BoolP("all", "a", false, "Show all pipelines (default shows just running/pending.)")
	recent = cmd.Flags().BoolP("recent", "r", false, "Show only the most recent pipeline")

	return cmd
}
