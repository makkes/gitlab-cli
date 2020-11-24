package get

import (
	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/get/issues"
	"github.com/makkes/gitlab-cli/cmd/get/output"
	"github.com/makkes/gitlab-cli/cmd/get/pipelines"
	"github.com/makkes/gitlab-cli/cmd/get/projects"
	"github.com/makkes/gitlab-cli/cmd/get/vars"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, cfg config.Config) *cobra.Command {
	var project *string
	var format *string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more objects",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	project = cmd.PersistentFlags().StringP("project", "p", "", "If present, the project scope for this CLI request")
	format = output.AddFlag(cmd)

	cmd.AddCommand(issues.NewCommand(client, project, format))
	cmd.AddCommand(pipelines.NewCommand(client, project, format))
	cmd.AddCommand(projects.NewCommand(client, cfg, format))
	cmd.AddCommand(vars.NewCommand(client, project, format))

	return cmd
}
