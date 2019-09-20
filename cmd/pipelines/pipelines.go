package pipelines

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	var all *bool
	var recent *bool
	var quiet *bool
	cmd := &cobra.Command{
		Use:   "pipelines PROJECT",
		Short: "List pipelines of a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := client.FindProject(args[0])
			if err != nil {
				return fmt.Errorf("Cannot list pipelines: %s", err)
			}
			resp, _, err := client.Get("/projects/" + strconv.Itoa(project.ID) + "/pipelines")
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
			}
			if len(pipelines) <= 0 {
				return fmt.Errorf("No pipelines found for project '%s'", args[0])
			}

			filteredPipelines := pipelines.Filter(func(p api.Pipeline) bool {
				return *all || p.Status == "running" || p.Status == "pending"
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

			if *quiet {
				for _, p := range pds {
					fmt.Printf("%d:%d\n", p.ProjectID, p.ID)
				}
				return nil
			}
			table.PrintPipelines(pds)
			return nil
		},
	}

	all = cmd.Flags().BoolP("all", "a", false, "Show all pipelines (default shows just running/pending.)")
	recent = cmd.Flags().BoolP("recent", "r", false, "Show only the most recent pipeline")
	quiet = cmd.Flags().BoolP("quiet", "q", false, "Only display numeric IDs")

	return cmd
}
