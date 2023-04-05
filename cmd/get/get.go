package get

import (
	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/get/accesstokens"
	"github.com/makkes/gitlab-cli/cmd/get/issues"
	"github.com/makkes/gitlab-cli/cmd/get/jobs"
	"github.com/makkes/gitlab-cli/cmd/get/logs"
	"github.com/makkes/gitlab-cli/cmd/get/output"
	"github.com/makkes/gitlab-cli/cmd/get/pipelines"
	"github.com/makkes/gitlab-cli/cmd/get/projects"
	"github.com/makkes/gitlab-cli/cmd/get/vars"
	"github.com/makkes/gitlab-cli/config"
)

func NewCommand(client api.Client, cfg config.Config) *cobra.Command {
	var format *string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more objects",
	}

	format = output.AddFlag(cmd)

	cmd.AddCommand(issues.NewCommand(client, format))
	cmd.AddCommand(pipelines.NewCommand(client, format))
	cmd.AddCommand(projects.NewCommand(client, cfg, format))
	cmd.AddCommand(vars.NewCommand(client, format))
	cmd.AddCommand(jobs.NewCommand(client, format))
	cmd.AddCommand(logs.NewCommand(client))
	cmd.AddCommand(accesstokens.NewCommand(client, format))

	return cmd
}
