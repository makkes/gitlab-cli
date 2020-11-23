package inspect

import (
	"github.com/makkes/gitlab-cli/cmd/inspect/issue"
	"github.com/makkes/gitlab-cli/cmd/inspect/pipeline"
	"github.com/makkes/gitlab-cli/cmd/inspect/project"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Show details of a specific object",
	}

	cmd.AddCommand(issue.NewCommand(client))
	cmd.AddCommand(pipeline.NewCommand(client))
	cmd.AddCommand(project.NewCommand(client))

	return cmd
}
