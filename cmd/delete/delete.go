package delete

import (
	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/delete/variable"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client, cfg config.Config) *cobra.Command {
	var project *string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources such as projects or variables",
	}

	project = cmd.PersistentFlags().StringP("project", "p", "", "If present, the project scope for this CLI request")

	cmd.AddCommand(variable.NewCommand(client, project))

	return cmd
}
