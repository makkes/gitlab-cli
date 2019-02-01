package pipelines

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	var all *bool
	var recent *bool
	cmd := &cobra.Command{
		Use:   "pipelines PROJECT",
		Short: "List pipelines of a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := client.FindProject(args[0])
			if err != nil {
				fmt.Printf("Error finding projects: %s\n", err)
				return
			}
			if projects == nil {
				fmt.Printf("Project '%s' not found\n", args[0])
				return
			}
			// just pick the first one. Room for improvement on the selection algorithm.
			project := projects[0]
			resp, err := client.Get("/projects/" + strconv.Itoa(project.ID) + "/pipelines")
			if err != nil {
				fmt.Println(err)
				return
			}
			var pipelines api.Pipelines
			err = json.Unmarshal(resp, &pipelines)
			if err != nil {
				fmt.Println(err)
				return
			}
			if *recent {
				pipelines = []api.Pipeline{pipelines[0]}
			}
			if len(pipelines) <= 0 {
				fmt.Printf("No pipelines found for project '%s'\n", args[0])
				return
			}

			filteredPipelines := pipelines.Filter(func(p api.Pipeline) bool {
				return *all || p.Status == "running" || p.Status == "pending"
			})

			pds := make([]api.PipelineDetails, 0)
			for _, p := range filteredPipelines {
				resp, err = client.Get("/projects/" + strconv.Itoa(project.ID) + "/pipelines/" + strconv.Itoa(p.ID))
				var pd api.PipelineDetails
				err = json.Unmarshal(resp, &pd)
				if err != nil {
					fmt.Println(err)
				} else {
					pds = append(pds, pd)
				}
			}

			table.PrintPipelines(pds)
		},
	}

	all = cmd.Flags().BoolP("all", "a", false, "Show all pipelines (default shows just running/pending.)")
	recent = cmd.Flags().BoolP("recent", "r", false, "Show only the most recent pipeline")

	return cmd
}
